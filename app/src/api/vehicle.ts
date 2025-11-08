import { useQuery } from '@tanstack/react-query';

import { apiClient } from '@api/client';
import { type VehicleSummary, type VehicleData, type VehiclesResponse } from '@types/vehicle';

export async function fetchVehicles(): Promise<VehicleSummary[]> {
  const { data } = await apiClient.get<VehiclesResponse>('/1/vehicles');
  return data?.response ?? [];
}

export async function fetchVehicleData(vehicleTag: string): Promise<VehicleData> {
  const { data } = await apiClient.get<{ response: VehicleData }>(
    `/1/vehicles/${vehicleTag}/vehicle_data`
  );
  return data?.response;
}

export function useVehiclesQuery() {
  return useQuery({
    queryKey: ['vehicles'],
    queryFn: fetchVehicles,
    staleTime: 60 * 1000
  });
}

export function useVehicleDataQuery(vehicleTag: string) {
  return useQuery({
    queryKey: ['vehicle-data', vehicleTag],
    queryFn: () => fetchVehicleData(vehicleTag),
    enabled: Boolean(vehicleTag)
  });
}
