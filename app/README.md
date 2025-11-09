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
- React Query 封装 `/api/1/vehicles`、`/api/1/vehicles/{vehicle_tag}/vehicle_data` 以及 `/api/1/vehicles/{vehicle_tag}/wake_up`（Fleet 官方唤醒接口），进入列表时自动检测离线车辆并触发唤醒，伙伴接口返回的错误也会同步展示到日志。
- 车辆列表页面展示车辆卡片、在线状态、电量信息，可跳转到详情页查看充电/空调/位置信息。

## 后续规划

1. **命令中心（阶段 3）**  
   - 引入命令库与动态参数表单，优先覆盖常用操作（锁车、空调、充电、哨兵模式等）。  
   - 执行命令前检查车辆状态（离线、充电中等），并对执行结果做 Toast/错误提示。  
   - 将命令执行记录缓存在本地，方便用户查看最近操作。

2. **驾驶员与设置**  
   - 展示 `/api/1/vehicles/{vehicle_tag}/drivers` 数据，支持管理驾驶员（若后端开放写操作）。  
   - Settings 页整合环境切换、日志导出、版本信息等。

3. **体验优化**  
   - 为车辆详情添加地图、图表等可视化组件。  
   - 离线/唤醒失败时在 UI 提示，并允许用户手动重试。  
   - 按需引入 Push 通知、国际化、平板适配等增强功能。
