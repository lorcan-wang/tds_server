import { useCallback, useLayoutEffect, useEffect, useRef } from 'react';
import {
  ActivityIndicator,
  FlatList,
  Pressable,
  RefreshControl,
  SafeAreaView,
  StyleSheet,
  Text,
  View
} from 'react-native';

import { useNavigation } from '@react-navigation/native';
import type { NativeStackNavigationProp } from '@react-navigation/native-stack';

import VehicleCard from '@components/VehicleCard';
import { useVehiclesQuery, fetchVehicleData, wakeVehicle } from '@api/vehicle';
import type { RootStackParamList } from '@navigation/index';
import { useAuthStore } from '@store/authStore';
import { useQueryClient } from '@tanstack/react-query';
import type { AxiosError } from 'axios';

const VehicleListScreen = () => {
  const navigation = useNavigation<NativeStackNavigationProp<RootStackParamList>>();
  const { data, isLoading, error, refetch, isRefetching } = useVehiclesQuery();
  const debugClear = useAuthStore((state) => state.debugClear);
  const queryClient = useQueryClient();
  const wokenVehiclesRef = useRef<Set<number>>(new Set());

  const handleRefresh = useCallback(() => {
    void refetch();
  }, [refetch]);

  useLayoutEffect(() => {
    navigation.setOptions({
      headerRight: () => (
        <Pressable onPress={debugClear}>
          <Text style={styles.logout}>退出</Text>
        </Pressable>
      )
    });
  }, [navigation, debugClear]);

  useEffect(() => {
    if (!data) return;
    data.forEach((vehicle) => {
      if (vehicle.state !== 'online' && !wokenVehiclesRef.current.has(vehicle.id)) {
        wokenVehiclesRef.current.add(vehicle.id);
        void wakeVehicle(vehicle.id.toString()).catch((err: AxiosError | Error) => {
          const status =
            'response' in err && err.response
              ? `${(err as AxiosError).response?.status}: ${(err as AxiosError).response?.statusText}`
              : '';
          const body =
            'response' in err && (err as AxiosError).response?.data
              ? JSON.stringify((err as AxiosError).response?.data)
              : err.message;
          console.warn(`唤醒车辆失败 ${vehicle.id} ${status}`, body);
        });
      }
    });
    data.forEach((vehicle) => {
      void queryClient.prefetchQuery({
        queryKey: ['vehicle-data', vehicle.id.toString()],
        queryFn: () => fetchVehicleData(vehicle.id.toString()),
        staleTime: 60 * 1000
      });
    });
  }, [data, queryClient]);

  const getBatteryLevel = useCallback(
    (vehicleId: number) => {
      const cached = queryClient.getQueryData<Awaited<ReturnType<typeof fetchVehicleData>>>([
        'vehicle-data',
        vehicleId.toString()
      ]);
      return cached?.charge_state?.battery_level;
    },
    [queryClient]
  );

  if (isLoading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#22D3EE" />
        <Text style={styles.hint}>正在获取车辆列表...</Text>
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.center}>
        <Text style={styles.error}>加载车辆失败，请稍后重试。</Text>
        <Text style={styles.errorDetail}>{String(error)}</Text>
      </View>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      <FlatList
        contentContainerStyle={styles.list}
        data={data ?? []}
        keyExtractor={(item) => item.id.toString()}
        renderItem={({ item }) => (
          <VehicleCard
            vehicle={item}
            batteryLevelOverride={getBatteryLevel(item.id)}
            onPress={() =>
              navigation.navigate('VehicleDetail', {
                vehicleTag: item.id.toString(),
                displayName: item.display_name || item.vin
              })
            }
          />
        )}
        refreshControl={
          <RefreshControl refreshing={isRefetching} onRefresh={handleRefresh} tintColor="#22D3EE" />
        }
        ListEmptyComponent={<Text style={styles.empty}>暂无车辆，请先完成授权。</Text>}
      />
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#0F172A'
  },
  list: {
    padding: 16
  },
  center: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 24,
    backgroundColor: '#0F172A'
  },
  hint: {
    color: '#CBD5F5',
    marginTop: 8
  },
  error: {
    color: '#FCA5A5',
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 8,
    textAlign: 'center'
  },
  errorDetail: {
    color: '#F87171',
    textAlign: 'center'
  },
  empty: {
    color: '#CBD5F5',
    textAlign: 'center',
    marginTop: 32
  },
  logout: {
    color: '#38BDF8',
    fontSize: 14,
    fontWeight: '600'
  }
});

export default VehicleListScreen;
