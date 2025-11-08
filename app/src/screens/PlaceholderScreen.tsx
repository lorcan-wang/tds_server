import { SafeAreaView, StyleSheet, Text, Pressable } from 'react-native';

import { useAuthStore } from '@store/authStore';

const PlaceholderScreen = () => {
  const debugClear = useAuthStore((state) => state.debugClear);

  return (
    <SafeAreaView style={styles.container}>
      <Text style={styles.title}>登录成功</Text>
      <Text style={styles.subtitle}>后续阶段将在此展示车辆列表与状态。</Text>
      <Pressable style={styles.button} onPress={debugClear}>
        <Text style={styles.buttonText}>临时：清除登录信息</Text>
      </Pressable>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 24,
    backgroundColor: '#0F172A'
  },
  title: {
    color: '#F8FAFC',
    fontSize: 22,
    fontWeight: '600',
    marginBottom: 8
  },
  subtitle: {
    color: '#CBD5F5',
    fontSize: 16,
    textAlign: 'center'
  },
  button: {
    marginTop: 24,
    paddingVertical: 12,
    paddingHorizontal: 16,
    borderRadius: 8,
    backgroundColor: '#22D3EE'
  },
  buttonText: {
    color: '#0F172A',
    fontSize: 14,
    fontWeight: '600'
  }
});

export default PlaceholderScreen;
