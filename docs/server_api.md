# tds_server API 说明 | Server API Reference

本文档汇总当前服务端对外暴露的 HTTP 接口，便于客户端集成与联调。所有示例均默认服务运行在 `http://localhost:8080`，如部署到其他环境请替换域名与端口。

## 基本信息 / Overview
- **Base Path**：所有业务接口位于 `/api` 前缀下。
- **Content-Type**：除特别说明外，响应均为 `application/json`。
- **身份认证**：受保护接口需要在请求头携带 `Authorization: Bearer <JWT>`。

## 鉴权流程 / Authentication Flow
1. **发起登录**：`GET /api/login?state=<optional-uuid>`，服务返回 302 并跳转至 Tesla OAuth 登录页，`state` 将在回调时原样返回，用于识别用户。
2. **处理回调**：Tesla 完成授权后访问 `GET /api/login/callback?code=<auth_code>&state=<uuid>`，服务会：
   - 交换 Tesla OAuth Token，并落库保存 `access_token`/`refresh_token`；
   - 颁发内部 JWT，响应示例：
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
3. **访问受保护接口**：后续请求均需携带回调中返回的 `jwt.token`：
   ```bash
   curl -H "Authorization: Bearer <jwt.token>" http://localhost:8080/api/list
   ```

## 公共端点 / Public Endpoints

| Method | Path | 描述 |
| ------ | ---- | ---- |
| GET | `/` | 健康检查，返回 `{"message":"ping"}`。 |
| GET | `/.well-known/appspecific/com.tesla.3p.public-key.pem` | 静态公钥文件，用于 Tesla 车辆指令签名验证。 |
| GET | `/api/login` | 触发 Tesla OAuth 登录，302 跳转。 |
| GET | `/api/login/callback` | 接收 Tesla 授权回调，返回 JWT 及 Tesla Token。 |

## 受保护端点 / Protected Endpoints

下列接口均需要有效 JWT。服务会自动刷新临近过期的 Tesla Access Token，并在必要时重试一次。

### 车辆清单 / Vehicle Inventory

| Method | Path | 说明 |
| ------ | ---- | ---- |
| GET | `/api/list` | 列出当前用户绑定的车辆（Tesla `/api/1/vehicles`）。 |
| GET | `/api/vehicles/:vehicle_tag` | 获取单辆车的基本信息。 |
| GET | `/api/vehicles/:vehicle_tag/vehicle_data` | 拉取完整 `vehicle_data`。可附加官方支持的查询参数。 |
| GET | `/api/vehicles/:vehicle_tag/states` | 返回车辆可用的状态段列表。 |
| GET | `/api/vehicles/:vehicle_tag/vehicle_data_request/:state` | 获取指定 `state`（如 `charge_state`、`drive_state`）的数据。 |
| POST | `/api/vehicles/:vehicle_tag/wake_up` | 唤醒休眠车辆。 |

### 驾驶员与共享 / Drivers & Sharing

| Method | Path | 说明 |
| ------ | ---- | ---- |
| GET | `/api/vehicles/:vehicle_tag/drivers` | 查询车辆允许的驾驶员。 |
| DELETE | `/api/vehicles/:vehicle_tag/drivers` | 移除驾驶员访问权限。 |
| GET | `/api/vehicles/:vehicle_tag/invitations` | 列出车辆的有效共享邀请。 |
| POST | `/api/vehicles/:vehicle_tag/invitations` | 创建共享邀请。 |
| POST | `/api/vehicles/:vehicle_tag/invitations/:invitation_id/revoke` | 撤销共享邀请。 |
| POST | `/api/invitations/redeem` | 兑换共享邀请链接。 |

### 车辆能力 / Vehicle Capabilities

| Method | Path | 说明 |
| ------ | ---- | ---- |
| GET | `/api/vehicles/:vehicle_tag/mobile_enabled` | 判断是否启用移动端访问。 |
| GET | `/api/vehicles/:vehicle_tag/nearby_charging_sites` | 查询车辆附近充电站。 |
| GET | `/api/vehicles/:vehicle_tag/recent_alerts` | 拉取车辆近期警报。 |
| GET | `/api/vehicles/:vehicle_tag/release_notes` | 获取车辆固件版本信息。 |
| GET | `/api/vehicles/:vehicle_tag/service_data` | 查询车辆维护/服务信息。 |

### 车队遥测 / Fleet Telemetry

| Method | Path | 说明 |
| ------ | ---- | ---- |
| POST | `/api/vehicles/fleet_telemetry_config` | 批量配置自托管遥测服务器。 |
| GET | `/api/vehicles/:vehicle_tag/fleet_telemetry_config` | 查看指定车辆的遥测配置。 |
| DELETE | `/api/vehicles/:vehicle_tag/fleet_telemetry_config` | 移除车辆遥测配置。 |
| POST | `/api/vehicles/fleet_telemetry_config_jws` | 使用签名令牌配置遥测（不推荐直接调用）。 |
| GET | `/api/vehicles/:vehicle_tag/fleet_telemetry_errors` | 查看车辆最新遥测错误。 |

### 通知订阅 / Subscriptions

| Method | Path | 说明 |
| ------ | ---- | ---- |
| GET | `/api/vehicle_subscriptions` | 查看当前设备订阅的车辆。 |
| POST | `/api/vehicle_subscriptions` | 更新设备订阅的车辆列表。 |
| GET | `/api/subscriptions` | 查询移动设备订阅的车辆列表。 |
| POST | `/api/subscriptions` | 设置移动设备订阅的车辆。 |
| GET | `/api/dx/vehicles/subscriptions/eligibility` | 查询车辆可用的订阅资格（需 `vin` 查询参数）。 |
| GET | `/api/dx/vehicles/upgrades/eligibility` | 查询车辆可用的软件升级（需 `vin` 查询参数）。 |

### 其他接口 / Miscellaneous

| Method | Path | 说明 |
| ------ | ---- | ---- |
| POST | `/api/vehicles/fleet_status` | 获取车队状态摘要。 |
| GET | `/api/dx/vehicles/options` | 查询车辆选装配置（需 `vin` 查询参数，官方暂未完全开放）。 |
| GET | `/api/dx/warranty/details` | 查询车辆保修信息。 |
| POST | `/api/vehicles/:vehicle_tag/signed_command` | 发送签名车辆指令（Tesla Vehicle Command Protocol）。 |
| POST | `/api/vehicles/:vehicle_tag/command/*command_path` | 发送 REST 车辆命令；优先使用 SDK，必要时回退至 Tesla REST。 |

### 车辆指令 / Vehicle Commands

- REST 命令接口：`POST /api/vehicles/:vehicle_tag/command/*command_path`，请求体为官方 Fleet API 定义的 JSON；服务会将路径中的 `command_path` 透传（如 `wake_up`、`climate_on`）。
- 签名指令接口：`POST /api/vehicles/:vehicle_tag/signed_command`，需准备 Vehicle Command Protocol 载荷并使用绑定的私钥签名。
- 内置 SDK 指令清单、参数要求及回退策略请参考 `docs/vehicle_command_reference.md`。该文档详列所有可调用指令、必填字段、REST 回退场景以及注意事项。

常见操作示例可结合上述文档与下方示例命令使用。

## 请求示例 / Sample Requests

列出车辆并指定返回字段：
```bash
curl -H "Authorization: Bearer <jwt>" \
     "http://localhost:8080/api/list?fields=vehicle_name,state"
```

唤醒并发送车辆命令：
```bash
curl -X POST -H "Authorization: Bearer <jwt>" \
     "http://localhost:8080/api/vehicles/12345/wake_up"

curl -X POST -H "Authorization: Bearer <jwt>" \
     -H "Content-Type: application/json" \
     -d '{"on": true}' \
     "http://localhost:8080/api/vehicles/12345/command/climate_on"
```

## 注意事项 / Notes
- 服务会在访问 Tesla API 前检查 Access Token 是否在 5 分钟内过期，并自动刷新；若仍遇到 `401` 会重试一次。
- 所有路径参数如 `:vehicle_tag`、`:invitation_id` 需替换为 Tesla 返回的真实值，`vin` 查询参数遵循官方文档要求。
- 某些端点（例如 `options`）可能仍处于内测或未开放状态，调用时请关注 Tesla 官方返回的错误描述。
