import { useEffect, type ReactNode } from 'react';
import { ActivityIndicator, View, StyleSheet } from 'react-native';
import { NavigationContainer } from '@react-navigation/native';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import RootNavigator from '@navigation/index';
import { hydrateAuthStore, useAuthStore } from '@store/authStore';

const queryClient = new QueryClient();

const BootstrapGate = ({ children }: { children: ReactNode }) => {
  const hydrated = useAuthStore((state) => state.hydrated);

  useEffect(() => {
    hydrateAuthStore();
  }, []);

  if (!hydrated) {
    return (
      <View style={styles.loading}>
        <ActivityIndicator size="large" color="#22D3EE" />
      </View>
    );
  }

  return <>{children}</>;
};

const styles = StyleSheet.create({
  loading: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: '#0F172A'
  }
});

export default function App() {
  return (
    <SafeAreaProvider>
      <QueryClientProvider client={queryClient}>
        <BootstrapGate>
          <NavigationContainer>
            <RootNavigator />
          </NavigationContainer>
        </BootstrapGate>
      </QueryClientProvider>
    </SafeAreaProvider>
  );
}
