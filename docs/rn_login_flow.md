# React Native 登录流程示例

以下示例演示如何在 React Native 中使用 `react-native-webview` 完成 Tesla OAuth 登录，并接收服务端返回的用户凭证。服务端 `/api/login/callback` 已支持 WebView `postMessage` 机制，请按照步骤集成。

## 1. 安装依赖

```bash
yarn add react-native-webview
# 或
npm install react-native-webview
```

## 2. 组件示例

```tsx
import React, {useRef, useState} from 'react';
import {Button, SafeAreaView, Text, View} from 'react-native';
import WebView, {WebViewMessageEvent} from 'react-native-webview';

const TESLA_LOGIN_URL = 'https://your-domain.com/api/login';

export default function TeslaLoginScreen() {
  const webviewRef = useRef<WebView>(null);
  const [isVisible, setVisible] = useState(false);
  const [payload, setPayload] = useState<any>(null);

  const handleMessage = (event: WebViewMessageEvent) => {
    try {
      const data = JSON.parse(event.nativeEvent.data);
      setPayload(data);
      // TODO: 保存 data.jwt.token / data.tesla_token 等信息
    } catch (err) {
      console.warn('Failed to parse login callback payload', err);
    } finally {
      setVisible(false);
    }
  };

  return (
    <SafeAreaView style={{flex: 1}}>
      <Button title="登录 Tesla" onPress={() => setVisible(true)} />

      {isVisible && (
        <WebView
          ref={webviewRef}
          source={{uri: TESLA_LOGIN_URL}}
          onMessage={handleMessage}
          startInLoadingState
          javaScriptEnabled
          domStorageEnabled
          onNavigationStateChange={navState => {
            if (!navState.loading && navState.url.includes('/api/login/callback')) {
              // 可选：在回调页阻止跳转
            }
          }}
        />
      )}

      <View style={{padding: 16}}>
        <Text>登录结果：</Text>
        <Text selectable>{payload ? JSON.stringify(payload, null, 2) : '尚未登录'}</Text>
      </View>
    </SafeAreaView>
  );
}
```

### payload 格式

`payload` 对应后端在回调 HTML 中推送的数据，结构如下（与 `/api/login/callback?format=json` 返回一致）：

```json
{
  "user_id": "9d546e6b-3f5b-4a07-b1a4-4c1bdbd6c9c8",
  "jwt": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600,
    "issuer": "tds_server"
  },
  "tesla_token": {
    "access_token": "...",
    "refresh_token": "...",
    "expires_in": 28800
  }
}
```

## 3. 调试技巧

- 若需直接查看 JSON，可手动访问 `/api/login/callback?format=json`。
- WebView 如需加载自签名证书域名，可在开发环境启用忽略证书的设置或使用代理转发。
- 如果需要在桌面浏览器调试，回调页也支持 `window.opener.postMessage`，可在前端捕获该消息。***
