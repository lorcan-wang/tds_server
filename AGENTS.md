# Repository Guidelines

## 项目结构与模块组织
- 服务入口位于 `cmd/server/main.go`，负责装配 HTTP 路由并启动进程，新增端点时从此处注册。
- 核心逻辑位于 `internal`：`config` 加载环境变量，`data` 建立数据库连接，`handler` 承接 HTTP 请求，`repository` 处理持久化，`service` 封装业务流程，`router` 定义路由，`middleware`、`model` 与 `util` 提供共享能力。
- `internal/model` 保存领域模型，和数据库表结构保持同步；如字段有调整，需同时更新迁移脚本。
- 编译产物输出至 `bin/server`，请通过构建命令重新生成而非直接修改；旧版本可以删除重建。
- 建议在同级目录添加测试文件，例如在 `internal/service/tesla_test.go` 覆盖对应服务逻辑，数据夹具放在包内。

## 构建、测试与开发命令
- `go fmt ./...`：按照 Go 官方格式化规则整理代码，提交前必跑，避免格式差异。
- `go build -o bin/server ./cmd/server`：生成本地或镜像使用的服务二进制，构建失败通常意味着依赖未正确导入。
- `go run ./cmd/server`：在当前工作空间直接运行 API 服务，便于调试，也会加载 `.env` 配置。
- `go test ./...`：执行所有包的单元测试，可通过 `-run TestName` 精确指定用例，建议在 PR 前确保全部通过。
- `docker build -t tds_server .`：基于仓库根目录的 `Dockerfile` 构建容器镜像，CI 也将使用同一流程验证环境。

## 代码风格与命名约定
- 使用 Go 默认缩进（制表符），保持每行一条语句，导入列表通过 `goimports` 排序，避免循环引用。
- 导出标识符采用 PascalCase，内部辅助函数使用 camelCase，包名保持简短小写（如 `handler`、`repository`），文件名突出资源或角色。
- 函数应聚焦单一职责，超出 80 行的函数请拆分子函数；结构体标签使用反引号并与数据库字段对应。
- 新增公共方法时补充注释，说明输入输出及边界情况，方便生成文档与静态分析。

## 测试指引
- 采用 Go 自带 `testing` 框架，推荐使用表驱动测试；函数命名遵循 `Test<单位>_<场景>`，便于阅读。
- 外部依赖请通过仓储或服务层模拟；集成测试需指向可丢弃数据库，依赖环境变量 `DB_*`，可配合 Docker Compose。
- 在测试文件开头记录所需数据或前置步骤，必要时提供伪造数据结构，确保团队成员能快速复现。
- 追求覆盖率稳定在关键服务模块 80% 以上，并在 PR 描述中附上 `go test -cover ./...` 摘要。

## 提交与拉取请求指引
- 遵循 `<type>: <summary>` 提交规范（如 `feat`、`fix`、`upd`、`mod`），摘要不超过 60 字并使用祈使语气，保持历史清晰。
- 如有关联 issue 请在提交或 PR 描述中引用，并标注配置或环境变量变更，方便追踪需求来源。
- PR 描述需包含意图说明、测试结果（如 `go test ./...` 输出）以及必要的人工验证步骤，界面改动可附截图或录屏。
- 遇到重大架构调整请在 PR 中附设计文档或流程图链接，确保评审具备足够上下文。

## 配置与环境
- `.env` 会被 `internal/config` 自动加载，用于配置 Tesla OAuth 凭证、接口地址与数据库参数（`DB_HOST`、`DB_PORT` 等）；示例可参考 `Dockerfile` 中的占位。
- 默认 HTTP 监听地址为 `:8080`，Postgres 端口为 `5432`；请按环境覆盖并避免在源码中硬编码机密，可通过系统变量覆盖。
- 新增环境变量时，请同步更新 `config.Config` 结构体并在 PR 中解释用途，必要时调整部署文档或脚本。
- 发布前核对 `.env` 与部署环境配置，避免遗漏 Client Secret 等关键参数导致服务启动失败。
- `TESLA_API_URL` 默认指向 `https://fleet-api.prd.cn.vn.cloud.tesla.cn`，同时被合作伙伴令牌与注册流程共用；如需切换环境请同时更新相关依赖。
- `TESLA_PARTNER_DOMAIN` 默认值为 `dwdacbj25q.ap-southeast-1.awsapprunner.com`，生产环境必须改为实际绑定公开 `.well-known/appspecific/com.tesla.3p.public-key.pem` 的域名，否则注册会失败。
- `TESLA_COMMAND_KEY_FILE` 配置车辆指令私钥路径，默认读取 `public/.well-known/appspecific/private-key.pem`；需确保车辆已完成对应公钥绑定，否则新协议指令会失败回退。
- 车辆指令调用优先使用官方 `github.com/teslamotors/vehicle-command` SDK（参见 `internal/service/vehicle_command_service.go`），若 SDK 返回 `ErrVehicleCommandUseREST` 才会回退到 REST 调用；若需排查可开启日志并关注 SDK 输出。
- 所有车辆接口（`/list`、`/vehicles/:vehicle_tag/...`、`/command/*`）在进入 Tesla API 前会主动检测 token 是否即将过期（默认 5 分钟），必要时自动刷新并落库；出现 401 仍保留兜底重试，确保调用方无感知。

## 合作伙伴令牌
- `internal/service/partner_token.go` 在服务启动时通过 client_credentials 流程获取合作伙伴令牌，并缓存于内存；刷新会在过期前 30 秒触发。
- 首次成功获取令牌后会自动调用 `/api/1/partner_accounts` 完成官方账号注册，请确保 `TESLA_API_URL` 与域名配置正确。
- 如需在调试环境跳过注册，可临时注释服务调用，但务必在提交前恢复；线上环境依赖该注册确保命令下发权限。***
