# 特斯拉 Fleet API 调用备忘

## 总览
- Fleet API 是基于 REST 的 HTTPS 接口，常见车辆资源路径位于 `/api/1/vehicles` 前缀下。
- 本项目通过可配置的 `TESLA_AUTH_URL`、`TESLA_TOKEN_URL`、`TESLA_API_URL` 注入目标环境，请在 `.env` 中维护并避免硬编码。
- 所有请求均需附带 OAuth 2.0 Bearer Token，未授权会返回 `401`。

### JWT 环境变量
- `JWT_SECRET`：用于 HS256 签名的对称密钥，需确保在各环境中保持机密且足够复杂；建议长度至少 32 字节，可使用随机生成工具。
- `JWT_ISSUER`：JWT 的 `iss` 字段，默认值为 `tds_server`，如需跨服务校验可设为域名或服务 ID。
- `JWT_EXPIRATION`：JWT 有效期，采用 Go 时长语法（如 `24h`、`72h`）。默认 24 小时，生产环境建议依据业务安全策略调整。

## 鉴权流程
- **授权地址**：`BuildAuthURL` 使用 `response_type=code` 构建登录链接，关键参数：`client_id`、`redirect_uri`、`scope`。默认 scope 覆盖 `openid offline_access user_data vehicle_device_data vehicle_cmds vehicle_charging_cmds`，如需新增权限可扩展。
- **换取令牌**：调用 `POST TESLA_TOKEN_URL`，请求体示例：
  ```json
  {
    "grant_type": "authorization_code",
    "client_id": "<CLIENT_ID>",
    "client_secret": "<CLIENT_SECRET>",
    "audience": "<TESLA_API_URL>",
    "code": "<AUTH_CODE>",
    "redirect_uri": "<REDIRECT_URI>"
  }
  ```
- **刷新令牌**：官方文档提供 `grant_type=refresh_token`，需提交刷新令牌与 `client_id`、`client_secret`；项目已在 `service.RefreshToken` 中封装调用。
- 建议持久化 `access_token`、`refresh_token`、`expires_in`，并在 5 分钟前主动刷新或捕获 `401` 后自动刷新。

## 请求样例
```bash
curl -X GET \
  -H "Authorization: Bearer ${TESLA_ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  "${TESLA_API_URL}/api/1/vehicles"
```
```bash
curl -X POST \
  -H "Authorization: Bearer ${TESLA_ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"command":"wake_up"}' \
  "${TESLA_API_URL}/api/1/vehicles/${VEHICLE_ID}/wake_up"
```

## 开发注意事项
- **令牌刷新**：当返回 `401` 时代表访问令牌失效，可调用 `grant_type=refresh_token` 刷新；代码中已在 `handler` 里自动刷新并重试，保证调用无感知。
- **限流**：官方建议保持调用窗口在每辆车每分钟 < 30 次，命令类接口需处理 `429` 重试。
- **状态同步**：命令下发成功后仍需轮询车辆状态确认；可结合数据库记录下发 ID。
- **错误处理**：接口响应通常形如 `{"response": {...}, "error": "", "error_description": ""}`，应解析 `response` 内部字段。
- **日志与审计**：将请求 ID、车辆 ID、命令类型记录到日志，方便定位问题；敏感字段（如位置）谨慎输出。

## 后续扩展
- 若需 Webhook/订阅数据，请关注官方的 Streaming/Fleet Telemetry 方案。
- CI/CD 中可通过服务账号自动刷新令牌并注入到部署环境，避免人工干预。
