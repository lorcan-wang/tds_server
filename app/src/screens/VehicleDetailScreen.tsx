import { useMemo } from 'react';
import {
  ActivityIndicator,
  ScrollView,
  StyleSheet,
  Text,
  View,
  RefreshControl
} from 'react-native';

import { RouteProp, useRoute } from '@react-navigation/native';

import { useVehicleDataQuery } from '@api/vehicle';
import type { RootStackParamList } from '@navigation/index';

const Section = ({ title, children }: { title: string; children: React.ReactNode }) => (
  <View style={styles.section}>
    <Text style={styles.sectionTitle}>{title}</Text>
    {children}
  </View>
);

const VehicleDetailScreen = () => {
  const {
    params: { vehicleTag }
  } = useRoute<RouteProp<RootStackParamList, 'VehicleDetail'>>();

  const { data, isLoading, error, refetch, isRefetching } = useVehicleDataQuery(vehicleTag);

  const metaItems = useMemo(() => {
    if (!data) return [];
    return [
      { label: 'VIN', value: data.vin },
      { label: '状态', value: data.state },
      { label: '访问类型', value: data.access_type ?? '未知' },
      { label: '固件版本', value: data.api_version ?? '未知' }
    ];
  }, [data]);

  if (isLoading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#22D3EE" />
        <Text style={styles.hint}>正在获取车辆数据...</Text>
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.center}>
        <Text style={styles.error}>加载车辆详情失败，请稍后重试。</Text>
        <Text style={styles.errorDetail}>{String(error)}</Text>
      </View>
    );
  }

  if (!data) {
    return (
      <View style={styles.center}>
        <Text style={styles.error}>暂无车辆数据。</Text>
      </View>
    );
  }

  return (
    <ScrollView
      style={styles.container}
      contentContainerStyle={styles.scrollContent}
      refreshControl={<RefreshControl refreshing={isRefetching} onRefresh={() => void refetch()} />}
    >
      <Section title="基础信息">
        {metaItems.map((item) => (
          <Text key={item.label} style={styles.item}>
            {item.label}：{item.value}
          </Text>
        ))}
      </Section>

      <Section title="充电信息">
        <Text style={styles.item}>
          电量：{data.charge_state?.battery_level ?? '--'}
          %
        </Text>
        <Text style={styles.item}>
          状态：{data.charge_state?.charging_state ?? '未知'}
        </Text>
        <Text style={styles.item}>
          预计充满：{data.charge_state?.time_to_full_charge ?? '--'} 小时
        </Text>
      </Section>

      <Section title="空调信息">
        <Text style={styles.item}>
          车内温度：{data.climate_state?.inside_temp ?? '--'}℃
        </Text>
        <Text style={styles.item}>
          车外温度：{data.climate_state?.outside_temp ?? '--'}℃
        </Text>
        <Text style={styles.item}>
          空调：{data.climate_state?.is_climate_on ? '开启' : '关闭'}
        </Text>
      </Section>

      <Section title="位置">
        <Text style={styles.item}>
          经度：{data.drive_state?.longitude ?? '--'}
        </Text>
        <Text style={styles.item}>
          纬度：{data.drive_state?.latitude ?? '--'}
        </Text>
        <Text style={styles.item}>
          朝向：{data.drive_state?.heading ?? '--'}
        </Text>
      </Section>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#0F172A'
  },
  scrollContent: {
    padding: 16,
    paddingBottom: 32
  },
  section: {
    backgroundColor: '#1E293B',
    borderRadius: 16,
    padding: 16,
    marginBottom: 16
  },
  sectionTitle: {
    color: '#F8FAFC',
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 8
  },
  item: {
    color: '#CBD5F5',
    fontSize: 15,
    marginBottom: 4
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
  }
});

export default VehicleDetailScreen;
