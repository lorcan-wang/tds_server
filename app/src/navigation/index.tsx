import { createNativeStackNavigator } from '@react-navigation/native-stack';

import LoginScreen from '@screens/LoginScreen';
import PlaceholderScreen from '@screens/PlaceholderScreen';
import { useAuthStore } from '@store/authStore';

export type RootStackParamList = {
  Login: undefined;
  Home: undefined;
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
        <Stack.Screen
          name="Home"
          component={PlaceholderScreen}
          options={{ title: '车辆概览' }}
        />
      )}
    </Stack.Navigator>
  );
};

export default RootNavigator;
