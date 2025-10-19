# 车辆端点 | Tesla Fleet API

## 端点  
以下为车辆相关的 REST API 端点及其说明：

| 端点 | 方法 | 路径 | 描述 |
|------|------|------|------|
| drivers | GET | `/api/1/vehicles/{vehicle_tag}/drivers` | 返回车辆的所有允许的驾驶员。该端点仅供车主使用。  [oai_citation:0‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| drivers_remove | DELETE | `/api/1/vehicles/{vehicle_tag}/drivers` | 取消驾驶员对车辆的访问。共享用户只能删除自己的访问权限。所有者可以删除共享访问权限或自己的访问权限。  [oai_citation:1‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| eligible_subscriptions | GET | `/api/1/dx/vehicles/subscriptions/eligibility?vin={vin}` | 返回符合条件的车辆订阅。  [oai_citation:2‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| eligible_upgrades | GET | `/api/1/dx/vehicles/upgrades/eligibility?vin={vin}` | 返回符合条件的车辆升级。  [oai_citation:3‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| fleet_status | POST | `/api/1/vehicles/fleet_status` | 提供用于确定车辆状态与应用程序相关信息的必要信息。包括：<br>• `vehicle_command_protocol_required` — 车辆是否需要使用 Vehicle Command Protocol。<br>• `safety_screen_streaming_toggle_enabled` — 用户是否在“安全”界面中启用了“允许第三方应用数据流”开关。<br>• `firmware_version`、`fleet_telemetry_version`、`total_number_of_keys`、`discounted_device_data` 等。  [oai_citation:4‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| fleet_telemetry_config (create) | POST | `/api/1/vehicles/fleet_telemetry_config` | 配置车辆以连接到自托管 fleet-telemetry 服务器。一次调用可配置多辆车辆。若未指定 VIN，响应将包含 `skipped_vehicles`。VIN 可能被拒绝的原因包括：<br>• `missing_key` — 虚拟钥匙尚未添加到车辆中。<br>• `unsupported_hardware` — 2021 年之前的 Model S 和 Model X 不支持。<br>• `unsupported_firmware` — 固件版本早于 2023.20。<br>车辆最多可同时配置向 3 个第三方应用程序传输数据。  [oai_citation:5‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| fleet_telemetry_config (delete) | DELETE | `/api/1/vehicles/{vehicle_tag}/fleet_telemetry_config` | 断开车辆与自托管 fleet-telemetry 服务器的数据流连接。  [oai_citation:6‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| fleet_telemetry_config (get) | GET | `/api/1/vehicles/{vehicle_tag}/fleet_telemetry_config` | 获取车辆的 fleet-telemetry 配置。`synced = true` 表示车辆已采用目标配置。若 `limit_reached = true` 表示车辆已达最大支持应用数，新请求无法添加。  [oai_citation:7‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| fleet_telemetry_config_jws | POST | `/api/1/vehicles/fleet_telemetry_config_jws` | 通过接受签名配置令牌，将车辆配置为连接自托管 fleet-telemetry 服务器。**不建议直接使用此端点。**推荐通过 vehicle-command 代理调用 fleet_telemetry_config create。若直接使用，必须使用 NIST P-256 + SHA-256 的 Schnorr 签名算法创建 JWS 令牌。VIN 可能被拒绝的原因同上。  [oai_citation:8‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| fleet_telemetry_errors | GET | `/api/1/vehicles/{vehicle_tag}/fleet_telemetry_errors` | 返回车辆上最近的车队 Telemetry 错误。  [oai_citation:9‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| list | GET | `/api/1/vehicles` | 返回该账户下车辆的列表。默认页面大小为 100。  [oai_citation:10‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| mobile_enabled | GET | `/api/1/vehicles/{vehicle_tag}/mobile_enabled` | 返回车辆是否启用了移动端设备（如 App）访问。  [oai_citation:11‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| nearby_charging_sites | GET | `/api/1/vehicles/{vehicle_tag}/nearby_charging_sites` | 返回车辆当前位置附近的充电站。  [oai_citation:12‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| options | GET | `/api/1/dx/vehicles/options?vin={vin}` | 返回车辆选项详细信息。（暂未开放，即将推出）  [oai_citation:13‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| recent_alerts | GET | `/api/1/vehicles/{vehicle_tag}/recent_alerts` | 最近警报列表。  [oai_citation:14‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| release_notes | GET | `/api/1/vehicles/{vehicle_tag}/release_notes` | 返回固件版本信息。  [oai_citation:15‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| service_data | GET | `/api/1/vehicles/{vehicle_tag}/service_data` | 获取有关车辆维护状态的信息。  [oai_citation:16‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| share_invites | GET | `/api/1/vehicles/{vehicle_tag}/invitations` | 返回车辆的有效共享邀请（分页，每页最多25条记录）。  [oai_citation:17‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| share_invites (create) | POST | `/api/1/vehicles/{vehicle_tag}/invitations` | 创建共享邀请：<br>• 每个邀请链接仅供一次使用，并在24小时后过期。<br>• 使用邀请的帐户可获得 Tesla 应用对车辆的驾驶员访问权限（包括查看车辆实时位置、发送远程命令、将用户的 Tesla 个人资料下载到车辆）。<br>• 若用户未安装 Tesla 应用，将被跳转至网页获得指引。<br>• 每辆车最多可添加5个驾驶员。该 API 不要求车辆在线。  [oai_citation:18‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| share_invites (redeem) | POST | `/api/1/invitations/redeem` | 兑换共享邀请。兑换后，该帐户可在 Tesla 应用中访问车辆。  [oai_citation:19‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| share_invites (revoke) | POST | `/api/1/vehicles/{vehicle_tag}/invitations/{id}/revoke` | 撤销共享邀请。该操作使链接无效。  [oai_citation:20‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| signed_command | POST | `/api/1/vehicles/{vehicle_tag}/signed_command` | 向车辆发送 Tesla 车辆命令协议。参见 Vehicle Command SDK 了解更多。  [oai_citation:21‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| subscriptions | GET | `/api/1/subscriptions` | 返回此移动设备当前订阅推送通知的车辆列表。  [oai_citation:22‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| subscriptions (set) | POST | `/api/1/subscriptions` | 允许移动设备指定从哪些车辆接收推送通知。调用时仅需提供希望订阅的车辆 ID。  [oai_citation:23‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| vehicle | GET | `/api/1/vehicles/{vehicle_tag}` | 返回车辆信息。  [oai_citation:24‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| vehicle_data | GET | `/api/1/vehicles/{vehicle_tag}/vehicle_data` | 对车辆进行实时呼叫。如果车辆离线，对于运行固件版本 2023.38+ 的车辆，需要 location_data 来获取车辆位置。这将导致位置共享图标显示在车辆 UI 上。  [oai_citation:25‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| vehicle_subscriptions | GET | `/api/1/vehicle_subscriptions` | 返回此移动设备当前订阅推送通知的车辆列表。  [oai_citation:26‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| vehicle_subscriptions (set) | POST | `/api/1/vehicle_subscriptions` | 允许移动设备指定希望接收通知的车辆 ID。  [oai_citation:27‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| wake_up | POST | `/api/1/vehicles/{vehicle_tag}/wake_up` | 将车辆从睡眠状态唤醒。睡眠状态可最大限度减少闲置能耗。  [oai_citation:28‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |
| warranty_details | GET | `/api/1/dx/warranty/details` | 返回车辆的保修信息。  [oai_citation:29‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints) |

---

## 说明  
- 所有请求路径皆以 `/api/1/` 或 `/api/1/dx/` 开头。  
- `vehicle_tag` 指车辆的标识符标签（通常由 Tesla 分配给你的车队车辆）。  
- 某些端点（如 `options`）标注为“暂未开放”。  
- 涉及车队遥测（fleet telemetry）相关端点（如 `fleet_telemetry_config`）允许将车辆配置为连接自托管服务器，需要注意硬件、固件以及钥匙分发等限制。  
- 所有功能应遵从 Tesla 官方文档的使用政策、权限管理和安全要求。  

---

## 📄 来源  
文档整理自 Tesla 官方开发者文档：  
> “车辆端点 | 特斯拉车队 API (Chinese)” —— https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints  [oai_citation:30‡Tesla开发者](https://developer.tesla.cn/docs/fleet-api/endpoints/vehicle-endpoints)  