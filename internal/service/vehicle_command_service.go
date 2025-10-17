package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"tds_server/internal/config"

	"github.com/teslamotors/vehicle-command/pkg/account"
	"github.com/teslamotors/vehicle-command/pkg/cache"
	"github.com/teslamotors/vehicle-command/pkg/connector/inet"
	"github.com/teslamotors/vehicle-command/pkg/protocol"
	"github.com/teslamotors/vehicle-command/pkg/proxy"
)

var (
	// ErrVehicleCommandUseREST indicates the request should fall back to the REST API. ErrVehicleCommandUseREST 表示需要回退到 REST 接口。
	ErrVehicleCommandUseREST = errors.New("vehicle command requires REST fallback")
)

const (
	sessionCacheSize      = 1024
	commandUserAgent      = "tds_server"
	defaultCommandTimeout = proxy.DefaultTimeout
	expectedVINLength     = 17
	defaultResponseType   = "application/json"
)

// VehicleCommandService encapsulates command execution with Tesla's vehicle-command SDK. VehicleCommandService 使用 Tesla vehicle-command SDK 封装指令执行逻辑。
type VehicleCommandService struct {
	cfg        *config.Config
	commandKey protocol.ECDHPrivateKey
	sessions   *cache.SessionCache
	timeout    time.Duration
	vinLocks   sync.Map
}

// CommandResult captures a successful command execution. CommandResult 表示一次成功的车辆指令执行结果。
type CommandResult struct {
	Status      int
	Body        []byte
	ContentType string
}

// CommandError wraps failures with HTTP semantics. CommandError 使用 HTTP 语义封装失败信息。
type CommandError struct {
	Status int
	Body   []byte
	Err    error
}

func (e *CommandError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return http.StatusText(e.Status)
}

// NewVehicleCommandService constructs a VehicleCommandService. NewVehicleCommandService 构建一个新的 VehicleCommandService 实例。
func NewVehicleCommandService(cfg *config.Config) (*VehicleCommandService, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	if cfg.TeslaCommandKeyPath == "" {
		return nil, fmt.Errorf("tesla command private key path not configured")
	}

	key, err := protocol.LoadPrivateKey(cfg.TeslaCommandKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load command private key: %w", err)
	}

	return &VehicleCommandService{
		cfg:        cfg,
		commandKey: key,
		sessions:   cache.New(sessionCacheSize),
		timeout:    defaultCommandTimeout,
	}, nil
}

// Execute sends a vehicle command via the new protocol. Execute 会通过新协议发送车辆指令。
func (s *VehicleCommandService) Execute(ctx context.Context, vin, command string, payload []byte, oauthToken string) (*CommandResult, error) {
	if oauthToken == "" {
		return nil, &CommandError{Status: http.StatusUnauthorized, Err: errors.New("missing oauth token")}
	}
	if vin == "" || len(vin) != expectedVINLength {
		return nil, &CommandError{Status: http.StatusBadRequest, Err: fmt.Errorf("invalid vin: %s", vin)}
	}
	if command == "" {
		return nil, &CommandError{Status: http.StatusBadRequest, Err: errors.New("command is required")}
	}
	if ctx == nil {
		ctx = context.Background()
	}

	acct, err := account.New(oauthToken, commandUserAgent)
	if err != nil {
		return nil, &CommandError{Status: http.StatusForbidden, Err: err}
	}

	params := proxy.RequestParameters{}
	if len(payload) > 0 {
		if err := json.Unmarshal(payload, &params); err != nil {
			return nil, &CommandError{Status: http.StatusBadRequest, Err: fmt.Errorf("failed to parse request body: %w", err)}
		}
	}

	execCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	action, err := proxy.ExtractCommandAction(execCtx, command, params)
	if err != nil {
		if protocol.IsNominalError(err) {
			return nominalResult(err.Error()), nil
		}
		switch {
		case errors.Is(err, proxy.ErrCommandUseRESTAPI):
			return nil, ErrVehicleCommandUseREST
		case errors.Is(err, proxy.ErrCommandNotImplemented):
			return nil, &CommandError{Status: http.StatusNotImplemented, Err: err}
		default:
			var httpErr *inet.HTTPError
			if errors.As(err, &httpErr) {
				return nil, &CommandError{Status: httpErr.Code, Body: []byte(httpErr.Message), Err: err}
			}
			return nil, &CommandError{Status: http.StatusBadRequest, Err: err}
		}
	}

	unlock := s.lockVIN(vin)
	defer unlock()

	car, err := acct.GetVehicle(execCtx, vin, s.commandKey, s.sessions)
	if err != nil {
		return nil, &CommandError{Status: http.StatusInternalServerError, Err: err}
	}

	if err := car.Connect(execCtx); err != nil {
		return nil, &CommandError{Status: http.StatusInternalServerError, Err: err}
	}
	defer car.Disconnect()

	if err := car.StartSession(execCtx, nil); err != nil {
		if errors.Is(err, protocol.ErrProtocolNotSupported) {
			return nil, ErrVehicleCommandUseREST
		}
		return nil, &CommandError{Status: http.StatusInternalServerError, Err: err}
	}
	defer func() {
		_ = car.UpdateCachedSessions(s.sessions)
	}()

	if err := action(car); err != nil {
		if protocol.IsNominalError(err) {
			return nominalResult(err.Error()), nil
		}
		if errors.Is(err, proxy.ErrCommandUseRESTAPI) {
			return nil, ErrVehicleCommandUseREST
		}
		return nil, &CommandError{Status: http.StatusInternalServerError, Err: err}
	}

	return successResult(), nil
}

func (s *VehicleCommandService) lockVIN(vin string) func() {
	mutexAny, _ := s.vinLocks.LoadOrStore(vin, &sync.Mutex{})
	mutex := mutexAny.(*sync.Mutex)
	mutex.Lock()
	return func() {
		mutex.Unlock()
	}
}

type commandResponsePayload struct {
	Response struct {
		Result bool   `json:"result"`
		Reason string `json:"reason"`
	} `json:"response"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

func successResult() *CommandResult {
	return buildCommandResult(true, "")
}

func nominalResult(reason string) *CommandResult {
	return buildCommandResult(false, reason)
}

func buildCommandResult(result bool, reason string) *CommandResult {
	payload := commandResponsePayload{}
	payload.Response.Result = result
	payload.Response.Reason = reason

	body, _ := json.Marshal(payload)
	body = append(body, '\n')

	return &CommandResult{
		Status:      http.StatusOK,
		Body:        body,
		ContentType: defaultResponseType,
	}
}
