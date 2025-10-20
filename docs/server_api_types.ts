/**
 * Server API contracts for frontend integration.
 * 用于前端对接的接口类型定义，覆盖登录回调、常见通用错误以及常用 Tesla Fleet API 代理响应。
 */

export interface LoginCallbackResponse {
  user_id: string;
  jwt: JWTEnvelope;
  tesla_token: TeslaTokenResponse;
}

export interface JWTEnvelope {
  token: string;
  expires_in: number;
  issuer: string;
}

export interface TeslaTokenResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  id_token?: string;
  scope?: string;
  token_type?: string;
  state?: string;
}

export interface ErrorResponse {
  error: string;
}

export interface TeslaAPIEnvelope<T> {
  response: T;
  count?: number;
  next_page?: string;
  error?: string;
  error_description?: string;
}

// Vehicle inventory ---------------------------------------------------------------------------

export interface VehicleListResponse extends TeslaAPIEnvelope<VehicleSummary[]> {}

export interface VehicleSummary {
  id: number;
  vehicle_id: number;
  vin: string;
  display_name: string | null;
  option_codes?: string;
  color?: string | null;
  tokens?: string[];
  state: "online" | "offline" | "asleep" | "waking" | "unknown";
  in_service: boolean;
  calendar_enabled: boolean;
  api_version?: number;
  backseat_token?: string | null;
  backseat_token_updated_at?: string | null;
}

export interface VehicleDataResponse
  extends TeslaAPIEnvelope<VehicleDataPayload> {}

export interface VehicleDataPayload {
  id: number;
  vehicle_id: number;
  vin: string;
  display_name: string | null;
  state: string;
  charge_state: Record<string, unknown>;
  climate_state: Record<string, unknown>;
  drive_state: Record<string, unknown>;
  gui_settings: Record<string, unknown>;
  vehicle_config: Record<string, unknown>;
  vehicle_state: Record<string, unknown>;
  [key: string]: unknown;
}

// Drivers & sharing ---------------------------------------------------------------------------

export interface DriversResponse extends TeslaAPIEnvelope<DriverSummary[]> {}

export interface DriverSummary {
  account_type: "OWNER" | "DRIVER" | string;
  display_name: string;
  email: string;
  expires_at?: string | null;
  relationship: string;
}

export interface ShareInvitesResponse
  extends TeslaAPIEnvelope<ShareInvite[]> {}

export interface ShareInvite {
  id: string;
  created_at: string;
  expires_at: string;
  status: "pending" | "redeemed" | "revoked";
  inviter_display_name?: string;
  invitee_email?: string;
  short_url?: string;
}

export interface ShareInviteCreateRequest {
  email?: string;
  role?: "driver" | "owner";
}

// Fleet telemetry -----------------------------------------------------------------------------

export interface FleetTelemetryConfigRequest {
  vins: string[];
  endpoint_uri: string;
  streaming_policy?: "LIVE" | "HISTORY";
  client_public_key?: string;
}

export interface FleetTelemetryConfigResponse
  extends TeslaAPIEnvelope<FleetTelemetryConfigPayload> {}

export interface FleetTelemetryConfigPayload {
  configured_vehicles: FleetTelemetryVehicle[];
  skipped_vehicles: FleetTelemetrySkipReason[];
}

export interface FleetTelemetryVehicle {
  vin: string;
  synced: boolean;
  limit_reached?: boolean;
}

export interface FleetTelemetrySkipReason {
  vin: string;
  reason:
    | "missing_key"
    | "unsupported_hardware"
    | "unsupported_firmware"
    | "limit_reached"
    | string;
}

export interface FleetTelemetryErrorsResponse
  extends TeslaAPIEnvelope<FleetTelemetryError[]> {}

export interface FleetTelemetryError {
  vin: string;
  code: string;
  message: string;
  occurred_at: string;
}

// Subscriptions -------------------------------------------------------------------------------

export interface SubscriptionsRequest {
  vehicle_ids: number[];
}

export interface SubscriptionsResponse
  extends TeslaAPIEnvelope<SubscriptionSummary[]> {}

export interface SubscriptionSummary {
  vehicle_id: number;
  subscribed: boolean;
  notification_types?: string[];
}

export interface EligibilityQuery {
  vin: string;
}

export interface EligibilityResponse
  extends TeslaAPIEnvelope<EligibilityItem[]> {}

export interface EligibilityItem {
  vin: string;
  eligible: boolean;
  reasons?: string[];
}

// Vehicle commands ----------------------------------------------------------------------------

export interface VehicleCommandRESTRequest {
  [key: string]: unknown;
}

export interface VehicleCommandRESTResponse
  extends TeslaAPIEnvelope<Record<string, unknown>> {}

export interface VehicleCommandSignedRequest {
  payload: string;
  signature: string;
}
