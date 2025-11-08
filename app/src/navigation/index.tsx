import { createNativeStackNavigator } from '@react-navigation/native-stack';

import LoginScreen from '@screens/LoginScreen';
import VehicleListScreen from '@screens/VehicleListScreen';
import VehicleDetailScreen from '@screens/VehicleDetailScreen';
import { useAuthStore } from '@store/authStore';

export type RootStackParamList = {
  Login: undefined;
  VehicleList: undefined;
  VehicleDetail: { vehicleTag: string; displayName?: string };
};

const Stack = createNativeStackNavigator<RootStackParamList>();

const RootNavigator = () => {
  const isAuthenticated = Boolean(useAuthStore((state) => state.jwtToken));

  return (
    <Stack.Navigator
      screenOptions={{
        headerStyle: { backgroundColor: '#0F172A' },
        headerTintColor: '#FFFFFF',
        contentStyle: { backgroundColor: '#0F172A' }
      }}
    >
      {!isAuthenticated ? (
        <Stack.Screen name="Login" component={LoginScreen} options={{ headerShown: false }} />
      ) : (
        <>
          <Stack.Screen
            name="VehicleList"
            component={VehicleListScreen}
            options={{ title: '我的车辆' }}
          />
          <Stack.Screen
            name="VehicleDetail"
            component={VehicleDetailScreen}
            options={({ route }) => ({
              title: route.params.displayName ?? '车辆详情'
            })}
          />
        </>
      )}
    </Stack.Navigator>
  );
};

export default RootNavigator;
