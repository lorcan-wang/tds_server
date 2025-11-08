export type PaginationMeta = {
  previous: string | null;
  next: string | null;
  current: number;
  per_page: number;
  count: number;
  pages: number;
};

export type VehicleSummary = {
  id: number;
  vehicle_id: number;
  vin: string;
  display_name: string;
  state: string;
  color?: string | null;
  access_type?: string;
  in_service?: boolean;
  api_version?: number | null;
  charge_state?: {
    battery_level?: number;
  };
};

export type VehiclesResponse = {
  response: VehicleSummary[];
  pagination?: PaginationMeta;
  count: number;
};

export type ChargeState = {
  battery_level?: number;
  charging_state?: string;
  time_to_full_charge?: number;
};

export type ClimateState = {
  inside_temp?: number;
  outside_temp?: number;
  is_climate_on?: boolean;
};

export type DriveState = {
  latitude?: number;
  longitude?: number;
  heading?: number;
};

export type VehicleData = VehicleSummary & {
  charge_state?: ChargeState;
  climate_state?: ClimateState;
  drive_state?: DriveState;
};
