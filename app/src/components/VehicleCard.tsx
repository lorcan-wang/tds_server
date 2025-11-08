import { memo } from 'react';
import { Pressable, StyleSheet, Text, View } from 'react-native';

import type { VehicleSummary } from '@types/vehicle';

type Props = {
  vehicle: VehicleSummary;
  onPress: () => void;
  batteryLevelOverride?: number;
};

const VehicleCard = ({ vehicle, onPress, batteryLevelOverride }: Props) => {
  const battery = batteryLevelOverride ?? vehicle.charge_state?.battery_level;

  return (
    <Pressable onPress={onPress} style={styles.card}>
      <View style={styles.header}>
        <Text style={styles.name}>{vehicle.display_name || vehicle.vin}</Text>
        <Text
          style={[
            styles.badge,
            vehicle.state === 'online' ? styles.badgeOnline : styles.badgeOffline
          ]}
        >
          {vehicle.state}
        </Text>
      </View>
      <Text style={styles.vin}>{vehicle.vin}</Text>
      <Text style={styles.meta}>
        车辆 ID：{vehicle.vehicle_id} · 访问类型：{vehicle.access_type ?? '未知'}
      </Text>
      <Text style={styles.battery}>
        电量：{typeof battery === 'number' ? `${battery}%` : '未知'}
      </Text>
    </Pressable>
  );
};

const styles = StyleSheet.create({
  card: {
    backgroundColor: '#1E293B',
    borderRadius: 16,
    padding: 16,
    marginBottom: 12
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 6
  },
  name: {
    color: '#F8FAFC',
    fontSize: 18,
    fontWeight: '600',
    flex: 1,
    marginRight: 8
  },
  badge: {
    paddingHorizontal: 12,
    paddingVertical: 4,
    borderRadius: 999,
    fontSize: 12,
    textTransform: 'uppercase',
    color: '#FFFFFF'
  },
  badgeOnline: {
    backgroundColor: '#16A34A'
  },
  badgeOffline: {
    backgroundColor: '#9CA3AF'
  },
  vin: {
    color: '#CBD5F5',
    marginBottom: 4
  },
  meta: {
    color: '#94A3B8',
    fontSize: 13,
    marginBottom: 8
  },
  battery: {
    color: '#FDE047',
    fontWeight: '600',
    fontSize: 16
  }
});

export default memo(VehicleCard);
