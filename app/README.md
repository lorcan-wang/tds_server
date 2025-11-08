# TDS Vehicle App (Stage 2)

React Native / Expo 客户端，用于与 `tds_server` 后端交互。当前阶段在完整的登录体系之上，已经实现车辆列表与详情浏览，并接入 React Query 进行数据缓存。

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
    ├── api/               # axios 单例 + 车辆请求封装
    ├── components/        # 车辆卡片等 UI 组件
    ├── navigation/        # React Navigation Stack
    ├── screens/           # 登录 / 车辆列表 / 车辆详情
    ├── store/             # JWT & Tesla Token 管理 + SecureStore 落盘
    ├── types/             # Tesla 数据结构
    └── utils/             # SecureStore 工具
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
- 使用 `Zustand + expo-secure-store` 在本地持久化 JWT/Tesla token，App 启动时自动恢复登录态，可通过列表页右上角按钮临时清理登录态。
- axios 拦截器统一注入 Authorization Header，并在 `401` 时清理本地会话。
- React Query 封装 `/api/1/vehicles` 与 `/api/1/vehicles/{vehicle_tag}/vehicle_data`，提供下拉刷新、错误状态提示。
- 车辆列表页面展示车辆卡片、在线状态、电量信息，可跳转到详情页查看充电/空调/位置信息。

下一阶段将基于车辆数据继续实现命令中心、驾驶员管理等高级功能。
