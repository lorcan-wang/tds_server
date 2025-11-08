import * as SecureStore from 'expo-secure-store';

export async function setJSONItem<T>(key: string, value: T) {
  try {
    await SecureStore.setItemAsync(key, JSON.stringify(value));
  } catch (error) {
    console.warn('SecureStore set error', error);
  }
}

export async function getJSONItem<T>(key: string): Promise<T | null> {
  try {
    const raw = await SecureStore.getItemAsync(key);
    if (!raw) return null;
    return JSON.parse(raw) as T;
  } catch (error) {
    console.warn('SecureStore get error', error);
    return null;
  }
}

export async function deleteItem(key: string) {
  try {
    await SecureStore.deleteItemAsync(key);
  } catch (error) {
    console.warn('SecureStore delete error', error);
  }
}
