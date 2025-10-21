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
  pagination?: TeslaPagination;
  error?: string;
  error_description?: string;
}

export interface TeslaPagination {
  previous: string | null;
  next: string | null;
  current: number;
  per_page: number;
  count: number;
  pages: number;
}

// Vehicle inventory ---------------------------------------------------------------------------

export interface VehicleListResponse extends TeslaAPIEnvelope<VehicleSummary[]> {}

export interface VehicleResponse extends TeslaAPIEnvelope<VehicleSummary> {}

export interface VehicleWakeUpResponse extends TeslaAPIEnvelope<VehicleSummary> {}

export interface VehicleSummary {
  id: number;
  vehicle_id: number;
  vin: string;
  display_name: string | null;
  option_codes?: string;
  color?: string | null;
  tokens?: string[];
  access_type?: string;
  granular_access?: Record<string, boolean>;
  id_s?: string;
  state: "online" | "offline" | "asleep" | "waking" | "unknown";
  in_service: boolean;
  calendar_enabled: boolean;
  api_version?: number | null;
  backseat_token?: string | null;
  backseat_token_updated_at?: string | null;
  user_id?: number;
}

export interface VehicleDataResponse
  extends TeslaAPIEnvelope<VehicleDataPayload> {}

export interface VehicleDataPayload {
  id: number;
  vehicle_id: number;
  user_id?: number;
  vin: string;
  display_name: string | null;
  color?: string | null;
  access_type?: string;
  granular_access?: Record<string, boolean>;
  tokens?: string[];
  state: string;
  in_service?: boolean;
  id_s?: string;
  calendar_enabled?: boolean;
  api_version?: number | null;
  backseat_token?: string | null;
  backseat_token_updated_at?: string | null;
  charge_state: Record<string, unknown>;
  climate_state: Record<string, unknown>;
  drive_state: Record<string, unknown>;
  gui_settings: Record<string, unknown>;
  vehicle_config: Record<string, unknown>;
  vehicle_state: Record<string, unknown>;
  [key: string]: unknown;
}

// Drivers & sharing ---------------------------------------------------------------------------

export interface DriversResponse extends TeslaAPIEnvelope<DriverProfile[]> {}

export interface DriverProfile {
  my_tesla_unique_id: number;
  user_id: number;
  user_id_s: string;
  vault_uuid: string;
  driver_first_name?: string | null;
  driver_last_name?: string | null;
  granular_access: DriverGranularAccess;
  active_pubkeys: string[];
  public_key: string;
}

export interface DriverGranularAccess {
  hide_private?: boolean;
  [key: string]: boolean | undefined;
}

export interface DriverRemoveResponse extends TeslaAPIEnvelope<string> {}

export interface ShareInvitesResponse
  extends TeslaAPIEnvelope<ShareInvite[]> {}

export interface ShareInvite {
  id: number;
  owner_id: number;
  share_user_id: number | null;
  product_id: string;
  state: string;
  code: string;
  expires_at: string;
  revoked_at?: string | null;
  borrowing_device_id?: string | null;
  key_id?: string | null;
  product_type: string;
  share_type: string;
  share_user_sso_id?: string | null;
  active_pubkeys?: (string | null)[];
  id_s: string;
  owner_id_s: string;
  share_user_id_s: string | null;
  borrowing_key_hash?: string | null;
  vin: string;
  share_link?: string;
}

export interface ShareInviteCreateRequest {
  email?: string;
  role?: "driver" | "owner";
}

export interface ShareInviteCreateResponse
  extends TeslaAPIEnvelope<ShareInvite> {}

export interface ShareInviteRedeemResponse
  extends TeslaAPIEnvelope<ShareInviteRedeemResult> {}

export interface ShareInviteRedeemResult {
  vehicle_id_s: string;
  vin: string;
}

export interface ShareInviteRevokeResponse
  extends TeslaAPIEnvelope<boolean> {}

// Fleet telemetry -----------------------------------------------------------------------------

export interface FleetTelemetryConfigRequest {
  vins: string[];
  endpoint_uri: string;
  streaming_policy?: "LIVE" | "HISTORY";
  client_public_key?: string;
}

export interface FleetTelemetryConfigCreateResponse
  extends TeslaAPIEnvelope<FleetTelemetryConfigCreatePayload> {}

export interface FleetTelemetryConfigCreatePayload {
  updated_vehicles: number;
  skipped_vehicles: FleetTelemetrySkippedVehicles;
}

export interface FleetTelemetryConfigJWSResponse
  extends TeslaAPIEnvelope<FleetTelemetryConfigCreatePayload> {}

export interface FleetTelemetrySkippedVehicles {
  missing_key: string[];
  unsupported_hardware: string[];
  unsupported_firmware: string[];
  max_configs: string[];
}

export interface FleetTelemetryConfigGetResponse
  extends TeslaAPIEnvelope<FleetTelemetryConfigGetPayload> {}

export interface FleetTelemetryConfigGetPayload {
  synced: boolean;
  config: FleetTelemetryConfigDetails;
  alert_types?: string[];
  limit_reached: boolean;
  key_paired: boolean;
}

export interface FleetTelemetryConfigDetails {
  hostname: string;
  ca?: string;
  port?: number;
  prefer_typed?: boolean;
  fields?: Record<string, FleetTelemetryFieldRule>;
}

export interface FleetTelemetryFieldRule {
  interval_seconds: number;
  resend_interval_seconds?: number;
  minimum_delta?: number;
}

export interface FleetTelemetryConfigDeleteResponse
  extends TeslaAPIEnvelope<FleetTelemetryConfigDeletePayload> {}

export interface FleetTelemetryConfigDeletePayload {
  updated_vehicles: number;
}

export interface FleetTelemetryErrorsResponse
  extends TeslaAPIEnvelope<FleetTelemetryErrorsPayload> {}

export interface FleetTelemetryErrorsPayload {
  fleet_telemetry_errors: FleetTelemetryErrorDetail[];
}

export interface FleetTelemetryErrorDetail {
  name: string;
  error: string;
  vin: string;
}

// Vehicle capabilities ------------------------------------------------------------------------

export interface NearbyChargingSitesResponse
  extends TeslaAPIEnvelope<NearbyChargingSitesPayload> {}

export interface NearbyChargingSitesPayload {
  congestion_sync_time_utc_secs: number;
  destination_charging: DestinationChargingSite[];
  superchargers: SuperchargerSite[];
  timestamp: number;
}

export interface ChargingLocation {
  lat: number;
  long: number;
}

export interface ChargingSiteBase {
  location: ChargingLocation;
  name: string;
  type: string;
  distance_miles: number;
  amenities?: string;
}

export interface DestinationChargingSite extends ChargingSiteBase {}

export interface SuperchargerSite extends ChargingSiteBase {
  available_stalls: number;
  total_stalls: number;
  site_closed: boolean;
  billing_info?: string;
}

export interface RecentAlertsResponse
  extends TeslaAPIEnvelope<RecentAlertsPayload> {}

export interface RecentAlertsPayload {
  recent_alerts: RecentAlert[];
}

export interface RecentAlert {
  name: string;
  time: string;
  audience: string[];
  user_text?: string;
}

export interface ReleaseNotesResponse
  extends TeslaAPIEnvelope<ReleaseNotesPayload> {}

export interface ReleaseNotesPayload {
  response: {
    release_notes: ReleaseNote[];
  };
}

export interface ReleaseNote {
  title: string;
  subtitle?: string;
  description?: string;
  customer_version?: string;
  icon?: string;
  image_url?: string;
  light_image_url?: string;
}

export interface ServiceDataResponse
  extends TeslaAPIEnvelope<ServiceDataPayload> {}

export interface ServiceDataPayload {
  service_status: string;
  service_etc?: string;
  service_visit_number?: string;
  status_id?: number;
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

export interface SubscriptionEligibilityResponse
  extends TeslaAPIEnvelope<SubscriptionEligibilityPayload> {}

export interface SubscriptionEligibilityPayload {
  vin: string;
  country: string;
  eligible: SubscriptionEligibilityItem[];
}

export interface SubscriptionEligibilityItem {
  optionCode: string;
  product: string;
  startDate?: string;
  addons: SubscriptionEligibilityPricing[];
  billingOptions: SubscriptionEligibilityPricing[];
}

export interface SubscriptionEligibilityPricing {
  billingPeriod: string;
  currencyCode: string;
  optionCode: string;
  price: number;
  tax: number;
  total: number;
}

export interface UpgradeEligibilityResponse
  extends TeslaAPIEnvelope<UpgradeEligibilityPayload> {}

export interface UpgradeEligibilityPayload {
  vin: string;
  country: string;
  type: string;
  eligible: UpgradeEligibilityItem[];
}

export interface UpgradeEligibilityItem {
  optionCode: string;
  optionGroup: string;
  currentOptionCode?: string;
  pricing: UpgradeEligibilityPricing[];
}

export interface UpgradeEligibilityPricing {
  price: number;
  total: number;
  currencyCode: string;
  isPrimary?: boolean;
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

export interface SignedCommandResponse extends TeslaAPIEnvelope<string> {}

export interface FleetStatusResponse
  extends TeslaAPIEnvelope<FleetStatusPayload> {}

export interface FleetStatusPayload {
  key_paired_vins: string[];
  unpaired_vins: string[];
  vehicle_info: Record<string, FleetStatusVehicleInfo>;
}

export interface FleetStatusVehicleInfo {
  firmware_version?: string;
  vehicle_command_protocol_required?: boolean;
  discounted_device_data?: boolean;
  fleet_telemetry_version?: string;
  total_number_of_keys?: number;
  [key: string]: string | number | boolean | undefined;
}

export interface WarrantyDetailsResponse
  extends TeslaAPIEnvelope<WarrantyDetailsPayload> {}

export interface WarrantyDetailsPayload {
  activeWarranty: WarrantyEntry[];
  upcomingWarranty: WarrantyEntry[];
  expiredWarranty: WarrantyEntry[];
}

export interface WarrantyEntry {
  warrantyType: string;
  warrantyDisplayName: string;
  expirationDate?: string | null;
  expirationOdometer?: number | null;
  odometerUnit?: string | null;
  warrantyExpiredOn?: string | null;
  coverageAgeInYears?: number | null;
}
