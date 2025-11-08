# TDS Vehicle App (Stage 1)

React Native / Expo 客户端，用于与 `tds_server` 后端交互。当前阶段聚焦在基础设施与登录流程搭建，确保可以在移动端完成 Tesla OAuth 并将凭证安全持久化。

## 目录结构

```
app/
├── App.tsx                # 入口，挂载 Provider 与导航
├── app.json               # Expo 配置（深链、图标等）
├── assets/                # 图标及启动图占位
├── babel.config.js        # 模块别名配置
├── package.json
├── tsconfig.json
└── src/
    ├── api/client.ts      # axios 单例
    ├── navigation/        # React Navigation Stack
    ├── screens/           # Login / 占位页面
    ├── store/authStore.ts # JWT & Tesla Token 管理 + SecureStore 持久化
    ├── providers/QueryProvider.tsx
    └── utils/secureStore.ts
```

## 快速开始

```bash
cd app
npm install    # 或 pnpm/yarn

# 运行开发服务器
npm run start
# 或安装自定义开发客户端
npx expo run:ios
npx expo run:android
```

环境变量通过 Expo 的 `EXPO_PUBLIC_*` 注入，例如在 shell 中设置：

```bash
export EXPO_PUBLIC_API_BASE_URL=http://localhost:8080/api
```

## 当前能力

- WebView + 深链完成 `/api/login` 授权，并解析回调 payload。
- 使用 `Zustand + expo-secure-store` 在本地持久化 JWT/Tesla token，App 启动时自动恢复登录态。
- axios 拦截器统一注入 Authorization Header，并在 `401` 时清理本地会话。
- 引入 `@tanstack/react-query`（全局 Provider 已就绪），为后续车辆数据缓存打基础。
- 样式暂以 React Native StyleSheet 管理，后续阶段可按需引入 design system。

后续阶段将基于此基础实现车辆列表、详情以及命令中心等业务功能。
