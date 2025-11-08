import { useCallback, useState } from 'react';
import { Linking, Modal, Pressable, SafeAreaView, StyleSheet, Text, View } from 'react-native';
import WebView, { type WebViewMessageEvent } from 'react-native-webview';
import { useFocusEffect } from '@react-navigation/native';
import { decode as atob } from 'base-64';

import { useAuthStore } from '@store/authStore';

const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL ?? 'http://localhost:8080';
const LOGIN_URL = `${API_BASE_URL}/api/login`;
const SCHEME_PREFIX = 'tdsclient://auth/callback';

const LoginScreen = () => {
  const setAuthPayload = useAuthStore((state) => state.setAuthPayload);
  const [showWebView, setShowWebView] = useState(false);

  const handlePayload = useCallback(
    (payload: unknown) => {
      if (!payload || typeof payload !== 'object') {
        return;
      }
      try {
        const parsed = payload as Parameters<typeof setAuthPayload>[0];
        setAuthPayload(parsed);
        setShowWebView(false);
      } catch (error) {
        console.warn('处理登录数据失败', error);
      }
    },
    [setAuthPayload]
  );

  const handleDeepLink = useCallback(
    (url: string | null) => {
      if (!url || !url.startsWith(SCHEME_PREFIX)) {
        return;
      }
      const match = url.match(/payload=([^&]+)/);
      if (!match) return;
      const payloadParam = match[1];
      try {
        const decoded = JSON.parse(atob(decodeURIComponent(payloadParam)));
        handlePayload(decoded);
      } catch (error) {
        console.warn('解析登录回调失败', error);
      }
    },
    [handlePayload]
  );

  useFocusEffect(
    useCallback(() => {
      const subscription = Linking.addEventListener('url', (event) => handleDeepLink(event.url));
      Linking.getInitialURL().then(handleDeepLink).catch(console.warn);
      return () => subscription.remove();
    }, [handleDeepLink])
  );

  const handleWebViewMessage = useCallback(
    (event: WebViewMessageEvent) => {
      try {
        const decoded = JSON.parse(event.nativeEvent.data);
        handlePayload(decoded);
      } catch (error) {
        console.warn('解析 WebView 消息失败', error);
      }
    },
    [handlePayload]
  );

  const openLogin = useCallback(() => {
    setShowWebView(true);
  }, []);

  return (
    <SafeAreaView style={styles.safeArea}>
      <View style={styles.card}>
        <Text style={styles.title}>连接 Tesla 账户</Text>
        <Text style={styles.subtitle}>
          点击下方按钮打开特斯拉授权页，完成登录后将自动返回应用。
        </Text>
        <Pressable style={styles.button} onPress={openLogin}>
          <Text style={styles.buttonText}>开始登录</Text>
        </Pressable>
      </View>

      <Modal visible={showWebView} animationType="slide">
        <SafeAreaView style={styles.modalSafeArea}>
          <View style={styles.modalHeader}>
            <Text style={styles.modalTitle}>Tesla 登录</Text>
            <Pressable onPress={() => setShowWebView(false)}>
              <Text style={styles.closeText}>关闭</Text>
            </Pressable>
          </View>
          <WebView
            source={{ uri: LOGIN_URL }}
            onMessage={handleWebViewMessage}
            sharedCookiesEnabled
            thirdPartyCookiesEnabled
          />
        </SafeAreaView>
      </Modal>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
    justifyContent: 'center',
    padding: 24,
    backgroundColor: '#0F172A'
  },
  card: {
    backgroundColor: '#1E293B',
    borderRadius: 16,
    padding: 24
  },
  title: {
    color: '#F8FAFC',
    fontSize: 24,
    fontWeight: '600',
    marginBottom: 12
  },
  subtitle: {
    color: '#CBD5F5',
    fontSize: 16,
    lineHeight: 22,
    marginBottom: 24
  },
  button: {
    backgroundColor: '#22D3EE',
    paddingVertical: 14,
    borderRadius: 12,
    alignItems: 'center'
  },
  buttonText: {
    color: '#0F172A',
    fontSize: 16,
    fontWeight: '600'
  },
  modalSafeArea: {
    flex: 1,
    backgroundColor: '#0F172A'
  },
  modalHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    paddingVertical: 12
  },
  modalTitle: {
    color: '#F8FAFC',
    fontSize: 18,
    fontWeight: '600'
  },
  closeText: {
    color: '#38BDF8',
    fontSize: 16
  }
});

export default LoginScreen;
