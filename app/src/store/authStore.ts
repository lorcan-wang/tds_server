import { create } from 'zustand';

import { deleteItem, getJSONItem, setJSONItem } from '@utils/secureStore';

const AUTH_STORAGE_KEY = 'tds-auth-payload';

export type TeslaTokenPayload = {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type?: string;
  scope?: string;
};

export type LoginPayload = {
  user_id: string;
  jwt: {
    token: string;
    expires_in: number;
    issuer: string;
  };
  tesla_token: TeslaTokenPayload;
};

type AuthState = {
  hydrated: boolean;
  userId?: string;
  jwtToken?: string;
  teslaToken?: TeslaTokenPayload;
  setAuthPayload: (payload: LoginPayload) => void;
  reset: () => void;
  debugClear: () => void;
};

export const useAuthStore = create<AuthState>((set) => ({
  hydrated: false,
  setAuthPayload: (payload) => {
    set({
      userId: payload.user_id,
      jwtToken: payload.jwt.token,
      teslaToken: payload.tesla_token
    });
    void setJSONItem<LoginPayload>(AUTH_STORAGE_KEY, payload);
  },
  reset: () => {
    set({
      userId: undefined,
      jwtToken: undefined,
      teslaToken: undefined
    });
    void deleteItem(AUTH_STORAGE_KEY);
  },
  debugClear: () => {
    set({
      userId: undefined,
      jwtToken: undefined,
      teslaToken: undefined
    });
    void deleteItem(AUTH_STORAGE_KEY);
  }
}));

export async function hydrateAuthStore() {
  const data = await getJSONItem<LoginPayload>(AUTH_STORAGE_KEY);
  if (data) {
    useAuthStore.setState({
      hydrated: true,
      userId: data.user_id,
      jwtToken: data.jwt.token,
      teslaToken: data.tesla_token
    });
    return;
  }
  useAuthStore.setState({ hydrated: true });
}
