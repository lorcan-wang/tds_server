# tds_server API 说明 | Server API Reference

本文档汇总当前服务端对外暴露的 HTTP 接口，便于客户端集成与联调。所有示例均默认服务运行在 `http://localhost:8080`，如部署到其他环境请替换域名与端口。

## 基本信息 / Overview
- **Base Path**：所有业务接口位于 `/api` 前缀下。
- **Content-Type**：除特别说明外，响应均为 `application/json`。
- **身份认证**：受保护接口需要在请求头携带 `Authorization: Bearer <JWT>`。

## 鉴权流程 / Authentication Flow
1. **发起登录**：`GET /api/login?state=<optional-uuid>`，服务返回 302 并跳转至 Tesla OAuth 登录页，`state` 将在回调时原样返回，用于识别用户。
2. **处理回调**：Tesla 完成授权后访问 `GET /api/login/callback?code=<auth_code>&state=<uuid>`，服务会：
   - 如果通过 WebView 打开，回调页面会通过 `window.ReactNativeWebView.postMessage` 推送 JSON，可参考 `docs/rn_login_flow.md` 获取 React Native 集成示例。
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

#### GET `/api/list`
- 返回类型：`VehicleListResponse`。
- 响应示例：
```json
{
  "response": [
    {
      "id": 100021,
      "vehicle_id": 99999,
      "vin": "TEST0000000VIN01",
      "color": null,
      "access_type": "OWNER",
      "display_name": "Owned",
      "option_codes": "TEST0,COUS",
      "granular_access": {
        "hide_private": false
      },
      "tokens": [
        "4f993c5b9e2b937b",
        "7a3153b1bbb48a96"
      ],
      "state": "online",
      "in_service": false,
      "id_s": "100021",
      "calendar_enabled": true,
      "api_version": null,
      "backseat_token": null,
      "backseat_token_updated_at": null
    }
  ],
  "pagination": {
    "previous": null,
    "next": null,
    "current": 1,
    "per_page": 2,
    "count": 2,
    "pages": 1
  },
  "count": 1
}
```

#### GET `/api/vehicles/:vehicle_tag`
- 返回类型：`VehicleResponse`。
- 响应示例：
```json
{
  "response": {
    "id": 100021,
    "vehicle_id": 99999,
    "vin": "TEST0000000VIN01",
    "color": null,
    "access_type": "OWNER",
    "display_name": "Owned",
    "option_codes": "TEST0,COUS",
    "granular_access": {
      "hide_private": false
    },
    "tokens": [
      "4f993c5b9e2b937b",
      "7a3153b1bbb48a96"
    ],
    "state": "online",
    "in_service": false,
    "id_s": "100021",
    "calendar_enabled": true,
    "api_version": null,
    "backseat_token": null,
    "backseat_token_updated_at": null
  }
}
```

#### GET `/api/vehicles/:vehicle_tag/vehicle_data`
- 返回类型：`VehicleDataResponse`。
- 响应示例（字段截取）：
```json
{
  "response": {
    "id": 100021,
    "user_id": 800001,
    "vehicle_id": 99999,
    "vin": "TEST00000000VIN01",
    "access_type": "OWNER",
    "granular_access": {
      "hide_private": false
    },
    "tokens": [
      "4f993c5b9e2b937b",
      "7a3153b1bbb48a96"
    ],
    "state": "online",
    "calendar_enabled": true,
    "api_version": 54,
    "charge_state": {
      "battery_level": 42,
      "charge_limit_soc": 90,
      "charging_state": "Disconnected",
      "time_to_full_charge": 0,
      "timestamp": 1692141038420
    },
    "climate_state": {
      "inside_temp": 38.4,
      "outside_temp": 36.5,
      "is_climate_on": false,
      "is_preconditioning": false,
      "timestamp": 1692141038419
    },
    "drive_state": {
      "latitude": 37.7765494,
      "longitude": -122.4195418,
      "heading": 289,
      "shift_state": null,
      "timestamp": 1692141038420
    },
    "gui_settings": {
      "gui_24_hour_time": false,
      "gui_distance_units": "mi/hr",
      "gui_temperature_units": "F"
    },
    "vehicle_config": {
      "car_type": "modely",
      "exterior_color": "MidnightSilver",
      "wheel_type": "Apollo19",
      "timestamp": 1692141038420
    },
    "vehicle_state": {
      "car_version": "2023.7.20 7910d26d5c64",
      "locked": true,
      "sentry_mode": false,
      "speed_limit_mode": {
        "active": false,
        "current_limit_mph": 85
      },
      "timestamp": 1692141038419
    }
  }
}
```

#### POST `/api/vehicles/:vehicle_tag/wake_up`
- 返回类型：`VehicleWakeUpResponse`。
- 响应示例：
```json
{
  "response": {
    "id": 100021,
    "user_id": 800001,
    "vehicle_id": 99999,
    "vin": "TEST0000000VIN01",
    "color": null,
    "access_type": "OWNER",
    "granular_access": {
      "hide_private": false
    },
    "tokens": [
      "4f993c5b9e2b937b",
      "7a3153b1bbb48a96"
    ],
    "state": "online",
    "in_service": false,
    "id_s": "100021",
    "calendar_enabled": true,
    "api_version": null,
    "backseat_token": null,
    "backseat_token_updated_at": null
  }
}
```

### 驾驶员与共享 / Drivers & Sharing

| Method | Path | 说明 |
| ------ | ---- | ---- |
| GET | `/api/vehicles/:vehicle_tag/drivers` | 查询车辆允许的驾驶员。 |
| DELETE | `/api/vehicles/:vehicle_tag/drivers` | 移除驾驶员访问权限。 |
| GET | `/api/vehicles/:vehicle_tag/invitations` | 列出车辆的有效共享邀请。 |
| POST | `/api/vehicles/:vehicle_tag/invitations` | 创建共享邀请。 |
| POST | `/api/vehicles/:vehicle_tag/invitations/:invitation_id/revoke` | 撤销共享邀请。 |
| POST | `/api/invitations/redeem` | 兑换共享邀请链接。 |

#### GET `/api/vehicles/:vehicle_tag/drivers`
- 返回类型：`DriversResponse`（定义见 `docs/server_api_types.ts`）。
- 响应示例：
```json
{
  "response": [
    {
      "my_tesla_unique_id": 8888888,
      "user_id": 800001,
      "user_id_s": "800001",
      "vault_uuid": "b5c443af-a286-49eb-a4ad-35a97963155d",
      "driver_first_name": "Testy",
      "driver_last_name": "McTesterson",
      "granular_access": {
        "hide_private": false
      },
      "active_pubkeys": [],
      "public_key": ""
    }
  ],
  "count": 1
}
```

#### DELETE `/api/vehicles/:vehicle_tag/drivers`
- 返回类型：`DriverRemoveResponse`。
- 响应示例：
```json
{
  "response": "ok"
}
```

#### GET `/api/vehicles/:vehicle_tag/invitations`
- 返回类型：`ShareInvitesResponse`。
- 响应示例：
```json
{
  "response": [
    {
      "id": 429509621657,
      "owner_id": 429511308124,
      "share_user_id": null,
      "product_id": "TEST0000000VIN01",
      "state": "pending",
      "code": "aqwl4JHU2q4aTeNROz8W9SpngoFvj-ReuDFIJs6-Y0hA",
      "expires_at": "2023-06-29T00:42:00.000Z",
      "revoked_at": null,
      "borrowing_device_id": null,
      "key_id": null,
      "product_type": "vehicle",
      "share_type": "customer",
      "share_user_sso_id": null,
      "active_pubkeys": [
        null
      ],
      "id_s": "429509621657",
      "owner_id_s": "429511308124",
      "share_user_id_s": "",
      "borrowing_key_hash": null,
      "vin": "TEST0000000VIN01",
      "share_link": "https://www.tesla.com/_rs/1/aqwl4JHU2q4aTeNROz8W9SpngoFvj-ReuDFIJs6-Y0hA"
    }
  ],
  "pagination": {
    "previous": null,
    "next": null,
    "current": 1,
    "per_page": 25,
    "count": 1,
    "pages": 1
  },
  "count": 1
}
```

#### POST `/api/vehicles/:vehicle_tag/invitations`
- 返回类型：`ShareInviteCreateResponse`。
- 响应示例：
```json
{
  "response": {
    "id": 429509621657,
    "owner_id": 429511308124,
    "share_user_id": null,
    "product_id": "TEST0000000VIN01",
    "state": "pending",
    "code": "aqwl4JHU2q4aTeNROz8W9SpngoFvj-ReuDFIJs6-Y0hA",
    "expires_at": "2023-06-29T00:42:00.000Z",
    "revoked_at": null,
    "borrowing_device_id": null,
    "key_id": null,
    "product_type": "vehicle",
    "share_type": "customer",
    "share_user_sso_id": null,
    "active_pubkeys": [
      null
    ],
    "id_s": "429509621657",
    "owner_id_s": "429511308124",
    "share_user_id_s": "",
    "borrowing_key_hash": null,
    "vin": "TEST0000000VIN01",
    "share_link": "https://www.tesla.com/_rs/1/aqwl4JHU2q4aTeNROz8W9SpngoFvj-ReuDFIJs6-Y0hA"
  }
}
```

#### POST `/api/invitations/redeem`
- 返回类型：`ShareInviteRedeemResponse`。
- 响应示例：
```json
{
  "response": {
    "vehicle_id_s": "88850",
    "vin": "5YJY000000NEXUS01"
  }
}
```

#### POST `/api/vehicles/:vehicle_tag/invitations/:invitation_id/revoke`
- 返回类型：`ShareInviteRevokeResponse`。
- 响应示例：
```json
{
  "response": true
}
```

### 车辆能力 / Vehicle Capabilities

| Method | Path | 说明 |
| ------ | ---- | ---- |
| GET | `/api/vehicles/:vehicle_tag/mobile_enabled` | 判断是否启用移动端访问。 |
| GET | `/api/vehicles/:vehicle_tag/nearby_charging_sites` | 查询车辆附近充电站。 |
| GET | `/api/vehicles/:vehicle_tag/recent_alerts` | 拉取车辆近期警报。 |
| GET | `/api/vehicles/:vehicle_tag/release_notes` | 获取车辆固件版本信息。 |
| GET | `/api/vehicles/:vehicle_tag/service_data` | 查询车辆维护/服务信息。 |

#### GET `/api/vehicles/:vehicle_tag/nearby_charging_sites`
- 返回类型：`NearbyChargingSitesResponse`。
- 响应示例：
```json
{
  "response": {
    "congestion_sync_time_utc_secs": 1693588513,
    "destination_charging": [
      {
        "location": {
          "lat": 37.409314,
          "long": -122.123068
        },
        "name": "Hilton Garden Inn Palo Alto",
        "type": "destination",
        "distance_miles": 1.35024,
        "amenities": "restrooms,wifi,lodging"
      },
      {
        "location": {
          "lat": 37.407771,
          "long": -122.120076
        },
        "name": "Dinah's Garden Hotel & Poolside Restaurant",
        "type": "destination",
        "distance_miles": 1.534213,
        "amenities": "restrooms,restaurant,wifi,cafe,lodging"
      }
    ],
    "superchargers": [
      {
        "location": {
          "lat": 37.399071,
          "long": -122.111216
        },
        "name": "Los Altos, CA",
        "type": "supercharger",
        "distance_miles": 2.202902,
        "available_stalls": 12,
        "total_stalls": 16,
        "site_closed": false,
        "amenities": "restrooms,restaurant,wifi,cafe,shopping",
        "billing_info": ""
      },
      {
        "location": {
          "lat": 37.441734,
          "long": -122.170202
        },
        "name": "Palo Alto, CA - Stanford Shopping Center",
        "type": "supercharger",
        "distance_miles": 2.339135,
        "available_stalls": 11,
        "total_stalls": 20,
        "site_closed": false,
        "amenities": "restrooms,restaurant,wifi,cafe,shopping",
        "billing_info": ""
      }
    ],
    "timestamp": 1693588576552
  }
}
```

#### GET `/api/vehicles/:vehicle_tag/recent_alerts`
- 返回类型：`RecentAlertsResponse`。
- 响应示例：
```json
{
  "response": {
    "recent_alerts": [
      {
        "name": "Name_Of_The_Alert",
        "time": "2021-03-19T22:01:15.101+00:00",
        "audience": [
          "service-fix",
          "customer"
        ],
        "user_text": "additional description text"
      }
    ]
  }
}
```

#### GET `/api/vehicles/:vehicle_tag/release_notes`
- 返回类型：`ReleaseNotesResponse`。
- 响应示例：
```json
{
  "response": {
    "response": {
      "release_notes": [
        {
          "title": "Minor Fixes",
          "subtitle": "Some more info",
          "description": "This release contains minor fixes and improvements",
          "customer_version": "2022.42.0",
          "icon": "release_notes_icon",
          "image_url": "https://vehicle-files.teslamotors.com/release_notes/d0fa3",
          "light_image_url": "https://vehicle-files.teslamotors.com/release_notes/d0fa3/light"
        }
      ]
    }
  }
}
```

#### GET `/api/vehicles/:vehicle_tag/service_data`
- 返回类型：`ServiceDataResponse`。
- 响应示例：
```json
{
  "response": {
    "service_status": "in_service",
    "service_etc": "2023-05-02T17:10:53-10:00",
    "service_visit_number": "SV12345678",
    "status_id": 8
  }
}
```

### 车队遥测 / Fleet Telemetry

| Method | Path | 说明 |
| ------ | ---- | ---- |
| POST | `/api/vehicles/fleet_telemetry_config` | 批量配置自托管遥测服务器。 |
| GET | `/api/vehicles/:vehicle_tag/fleet_telemetry_config` | 查看指定车辆的遥测配置。 |
| DELETE | `/api/vehicles/:vehicle_tag/fleet_telemetry_config` | 移除车辆遥测配置。 |
| POST | `/api/vehicles/fleet_telemetry_config_jws` | 使用签名令牌配置遥测（不推荐直接调用）。 |
| GET | `/api/vehicles/:vehicle_tag/fleet_telemetry_errors` | 查看车辆最新遥测错误。 |

#### POST `/api/vehicles/fleet_telemetry_config`
- 返回类型：`FleetTelemetryConfigCreateResponse`。
- 响应示例：
```json
{
  "response": {
    "updated_vehicles": 1,
    "skipped_vehicles": {
      "missing_key": [],
      "unsupported_hardware": [
        "5YJ3000000NEXUS02"
      ],
      "unsupported_firmware": [
        "5YJ3000000NEXUS02"
      ],
      "max_configs": []
    }
  }
}
```

#### GET `/api/vehicles/:vehicle_tag/fleet_telemetry_config`
- 返回类型：`FleetTelemetryConfigGetResponse`。
- 响应示例：
```json
{
  "response": {
    "synced": true,
    "config": {
      "hostname": "test-telemetry.com",
      "port": 4443,
      "prefer_typed": true,
      "fields": {
        "DriveRail": {
          "interval_seconds": 1800
        },
        "BmsFullchargecomplete": {
          "interval_seconds": 1800,
          "resend_interval_seconds": 3600
        },
        "ChargerVoltage": {
          "interval_seconds": 1,
          "minimum_delta": 5
        }
      }
    },
    "alert_types": [
      "service"
    ],
    "limit_reached": false,
    "key_paired": false
  }
}
```

#### DELETE `/api/vehicles/:vehicle_tag/fleet_telemetry_config`
- 返回类型：`FleetTelemetryConfigDeleteResponse`。
- 响应示例：
```json
{
  "response": {
    "updated_vehicles": 1
  }
}
```

#### POST `/api/vehicles/fleet_telemetry_config_jws`
- 返回类型：`FleetTelemetryConfigJWSResponse`。
- 响应示例：
```json
{
  "response": {
    "updated_vehicles": 1,
    "skipped_vehicles": {
      "missing_key": [],
      "unsupported_hardware": [
        "5YJ3000000NEXUS02"
      ],
      "unsupported_firmware": [
        "5YJ3000000NEXUS02"
      ],
      "max_configs": []
    }
  }
}
```

#### GET `/api/vehicles/:vehicle_tag/fleet_telemetry_errors`
- 返回类型：`FleetTelemetryErrorsResponse`。
- 响应示例：
```json
{
  "response": {
    "fleet_telemetry_errors": [
      {
        "name": "partner-client-id",
        "error": "msg",
        "vin": "vin"
      },
      {
        "name": "partner-client-id",
        "error": "msg2",
        "vin": "vin"
      }
    ]
  }
}
```

### 通知订阅 / Subscriptions

| Method | Path | 说明 |
| ------ | ---- | ---- |
| GET | `/api/vehicle_subscriptions` | 查看当前设备订阅的车辆。 |
| POST | `/api/vehicle_subscriptions` | 更新设备订阅的车辆列表。 |
| GET | `/api/subscriptions` | 查询移动设备订阅的车辆列表。 |
| POST | `/api/subscriptions` | 设置移动设备订阅的车辆。 |
| GET | `/api/dx/vehicles/subscriptions/eligibility` | 查询车辆可用的订阅资格（需 `vin` 查询参数）。 |
| GET | `/api/dx/vehicles/upgrades/eligibility` | 查询车辆可用的软件升级（需 `vin` 查询参数）。 |

#### GET `/api/dx/vehicles/subscriptions/eligibility`
- 返回类型：`SubscriptionEligibilityResponse`。
- 响应示例：
```json
{
  "response": {
    "country": "US",
    "eligible": [
      {
        "addons": [
          {
            "billingPeriod": "MONTHLY",
            "currencyCode": "USD",
            "optionCode": "APB",
            "price": 99,
            "tax": 0,
            "total": 99
          }
        ],
        "billingOptions": [
          {
            "billingPeriod": "ANNUAL",
            "currencyCode": "USD",
            "optionCode": "APB",
            "price": 999,
            "tax": 0,
            "total": 999
          }
        ],
        "optionCode": "APB",
        "product": "Premium Connectivity",
        "startDate": "2024-01-01"
      }
    ],
    "vin": "LRW3E7FA0MC000002"
  }
}
```

#### GET `/api/dx/vehicles/upgrades/eligibility`
- 返回类型：`UpgradeEligibilityResponse`。
- 响应示例：
```json
{
  "response": {
    "vin": "TEST0000000VIN01",
    "country": "US",
    "type": "VEHICLE",
    "eligible": [
      {
        "optionCode": "$FM3U",
        "optionGroup": "PERF_FIRMWARE",
        "currentOptionCode": "$FM3B",
        "pricing": [
          {
            "price": 2000,
            "total": 2000,
            "currencyCode": "USD",
            "isPrimary": true
          }
        ]
      }
    ]
  }
}
```

### 其他接口 / Miscellaneous

| Method | Path | 说明 |
| ------ | ---- | ---- |
| POST | `/api/vehicles/fleet_status` | 获取车队状态摘要。 |
| GET | `/api/dx/vehicles/options` | 查询车辆选装配置（需 `vin` 查询参数，官方暂未完全开放）。 |
| GET | `/api/dx/warranty/details` | 查询车辆保修信息。 |
| POST | `/api/vehicles/:vehicle_tag/signed_command` | 发送签名车辆指令（Tesla Vehicle Command Protocol）。 |
| POST | `/api/vehicles/:vehicle_tag/command/*command_path` | 发送 REST 车辆命令；优先使用 SDK，必要时回退至 Tesla REST。 |

#### POST `/api/vehicles/fleet_status`
- 返回类型：`FleetStatusResponse`。
- 响应示例：
```json
{
  "response": {
    "key_paired_vins": [],
    "unpaired_vins": [
      "5YJ3000000NEXUS01"
    ],
    "vehicle_info": {
      "5YJ3000000NEXUS01": {
        "firmware_version": "2024.14.30",
        "vehicle_command_protocol_required": true,
        "discounted_device_data": false,
        "fleet_telemetry_version": "1.0.0",
        "total_number_of_keys": 5
      }
    }
  }
}
```

#### GET `/api/dx/warranty/details`
- 返回类型：`WarrantyDetailsResponse`。
- 响应示例：
```json
{
  "response": {
    "activeWarranty": [
      {
        "warrantyType": "NEW_MFG_WARRANTY",
        "warrantyDisplayName": "Basic Vehicle Limited Warranty",
        "expirationDate": "2025-10-21T00:00:00Z",
        "expirationOdometer": 50000,
        "odometerUnit": "MI",
        "warrantyExpiredOn": null,
        "coverageAgeInYears": 4
      },
      {
        "warrantyType": "BATTERY_WARRANTY",
        "warrantyDisplayName": "Battery Limited Warranty",
        "expirationDate": "2029-10-21T00:00:00Z",
        "expirationOdometer": 120000,
        "odometerUnit": "MI",
        "warrantyExpiredOn": null,
        "coverageAgeInYears": 8
      }
    ],
    "upcomingWarranty": [],
    "expiredWarranty": []
  }
}
```

### 车辆指令 / Vehicle Commands

- REST 命令接口：`POST /api/vehicles/:vehicle_tag/command/*command_path`，请求体为官方 Fleet API 定义的 JSON；服务会将路径中的 `command_path` 透传（如 `wake_up`、`climate_on`）。
- 签名指令接口：`POST /api/vehicles/:vehicle_tag/signed_command`，需准备 Vehicle Command Protocol 载荷并使用绑定的私钥签名。
- 内置 SDK 指令清单、参数要求及回退策略请参考 `docs/vehicle_command_reference.md`。该文档详列所有可调用指令、必填字段、REST 回退场景以及注意事项。

常见操作示例可结合上述文档与下方示例命令使用。

#### POST `/api/vehicles/:vehicle_tag/signed_command`
- 返回类型：`SignedCommandResponse`。
- 响应示例：
```json
{
  "response": "base64_response"
}
```

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


## 契约与模拟数据 / Contracts & Mocks

- TypeScript 类型定义位于 `docs/server_api_types.ts`，覆盖登录回调、通用错误、车辆列表、共享邀请、遥测配置等结构，供前端直接引用。
- 常用响应示例存放在 `docs/mock/` 目录：
  - `login_callback.json`：登录回调成功示例。
  - `error_response.json`：统一错误格式。
  - `vehicles_list.json`：`/api/list` 车辆列表响应。
  - `vehicle_data.json`：`/api/vehicles/:vehicle_tag/vehicle_data` 精简示例。
  - `share_invites.json`：共享邀请列表。
  - `fleet_telemetry_config.json`：遥测配置返回。

如需新增端点，请同步更新类型定义与示例，确保前端契约一致。
## 注意事项 / Notes
- 服务会在访问 Tesla API 前检查 Access Token 是否在 5 分钟内过期，并自动刷新；若仍遇到 `401` 会重试一次。
- 所有路径参数如 `:vehicle_tag`、`:invitation_id` 需替换为 Tesla 返回的真实值，`vin` 查询参数遵循官方文档要求。
- 某些端点（例如 `options`）可能仍处于内测或未开放状态，调用时请关注 Tesla 官方返回的错误描述。
