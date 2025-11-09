# 车辆接口文档

本文档说明服务中当前开放的特斯拉车辆相关接口、请求参数以及返回字段定义，便于客户端直接对接。

## GET /api/1/vehicles

- **方法**：GET  
- **授权范围**：`vehicle_device_data`  
- **描述**：返回当前授权账号下的车辆列表，分页参数默认 `page=1`、`per_page=100`。

### 查询参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
| ---- | ---- | ---- | ---- | ---- |
| `page` | int | 否 | 请求的页码，默认为 1 | `1` |
| `per_page` | int | 否 | 每页返回的车辆数量，默认 100，最大值由特斯拉平台限制 | `50` |

### 响应体

#### 顶层字段

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `response` | `VehicleSummary[]` | 当前页车辆条目，字段说明见下。 |
| `pagination` | `PaginationMeta` | 分页游标信息，字段说明见下。 |
| `count` | int | 本次响应返回的车辆数量。 |

#### VehicleSummary

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `id` | int64 | 特斯拉 Fleet API 车辆唯一标识。 |
| `vehicle_id` | int64 | 供移动端调用使用的车辆编号。 |
| `vin` | string | 车辆 VIN。 |
| `color` | string\|null | 车辆颜色，未配置时为 `null`。 |
| `access_type` | string | 访问类型，例如 `OWNER`、`DRIVER`。 |
| `display_name` | string | 用户自定义的车辆展示名称。 |
| `option_codes` | string | 车辆选装代码，逗号分隔。 |
| `granular_access.hide_private` | bool | 是否对外隐藏隐私数据（如位置）。 |
| `tokens` | string[] | 旧版车辆控制接口使用的 token。 |
| `state` | string | 车辆当前状态（如 `online`、`asleep`）。 |
| `in_service` | bool | 是否处于服务/维修状态。 |
| `id_s` | string | 字符串形式的车辆 ID。 |
| `calendar_enabled` | bool | 是否启用了日历同步。 |
| `api_version` | int\|null | 车辆当前固件暴露的 API 版本。 |
| `backseat_token` | string\|null | 后排乘客账户 token。 |
| `backseat_token_updated_at` | string\|null | 后排 token 最近更新时间（ISO8601 字符串）。 |

#### PaginationMeta

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `previous` | string\|null | 上一页游标，无上一页时为 `null`。 |
| `next` | string\|null | 下一页游标，无下一页时为 `null`。 |
| `current` | int | 当前页码。 |
| `per_page` | int | 当前请求的每页数量。 |
| `count` | int | 当前页返回的车辆数量。 |
| `pages` | int | 根据条件可取得的总页数。 |

### 示例

```json
{
  "response": [
    {
      "id": 100021,
      "vehicle_id": 99999,
      "vin": "TEST000000VIN01",
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

## GET /api/1/vehicles/{vehicle_tag}/vehicle_data

- **方法**：GET  
- **授权范围**：`vehicle_device_data`、`vehicle_location`（当请求 `location_data` 或 `location_state`）以及 `device_data`  
- **描述**：触发车辆实时唤醒并返回综合状态，可通过 `endpoints` 指定需要的子模块，未指定时返回全部信息。

### 路径参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
| ---- | ---- | ---- | ---- | ---- |
| `vehicle_tag` | string | 是 | 车辆唯一标识，可为 `id`、`id_s` 或 `vin`。 | `100021` |

### 查询参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
| ---- | ---- | ---- | ---- | ---- |
| `endpoints` | string | 否 | 逗号分隔的模块列表，支持：`charge_state`、`climate_state`、`closures_state`、`drive_state`、`gui_settings`、`location_data`、`charge_schedule_data`、`preconditioning_schedule_data`、`vehicle_config`、`vehicle_state`、`vehicle_data_combo`。 | `charge_state,vehicle_state` |

### 响应体

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `response.user_id` | int64 | 车辆所属 Tesla 账户 ID。 |
| `response.user_vehicle_bound_at` | string\|null | 当前用户获得车辆访问权限的时间（ISO8601）。 |
| `response.<VehicleSummary 字段>` | - | 同列表接口返回的车辆概要信息。 |
| `response.charge_state` | object | 充电与电池信息，详见下表。 |
| `response.climate_state` | object | 空调与温控状态，详见下表。 |
| `response.closures_state` | object | 车门、车窗、前后备箱开闭状态（示例未返回）。 |
| `response.drive_state` | object | 行驶及导航状态，详见下表。 |
| `response.gui_settings` | object | 车机界面显示设置，详见下表。 |
| `response.location_data` | object | 精准经纬度（示例未返回，需额外授权）。 |
| `response.charge_schedule_data` | object | 定时充电配置（示例未返回）。 |
| `response.preconditioning_schedule_data` | object | 空调预设计划（示例未返回）。 |
| `response.vehicle_config` | object | 车辆硬件与配置参数，详见下表。 |
| `response.vehicle_state` | object | 车况状态、告警及多媒体信息，详见下表。 |
| `response.vehicle_data_combo` | object | 当请求聚合模块时返回（示例未返回）。 |

#### `charge_state` 字段

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `battery_heater_on` | bool | 电池加热器是否工作。 |
| `battery_level` | int | 当前电量百分比。 |
| `battery_range` | float | 估算续航（英里）。 |
| `charge_amps` | int | 预期充电电流（安培）。 |
| `charge_current_request` | int | 当前请求的充电电流。 |
| `charge_current_request_max` | int | 允许的最大充电电流。 |
| `charge_enable_request` | bool | 是否允许车辆开启充电。 |
| `charge_energy_added` | float | 本次充电累积能量（千瓦时）。 |
| `charge_limit_soc` | int | 当前充电上限 SoC。 |
| `charge_limit_soc_max` | int | 可设置的最大 SoC 上限。 |
| `charge_limit_soc_min` | int | 可设置的最小 SoC 下限。 |
| `charge_limit_soc_std` | int | 标准模式默认 SoC。 |
| `charge_miles_added_ideal` | float | 理想续航里程增量。 |
| `charge_miles_added_rated` | float | 额定续航里程增量。 |
| `charge_port_cold_weather_mode` | bool | 是否启用寒冷天气模式。 |
| `charge_port_color` | string | 车充口指示颜色。 |
| `charge_port_door_open` | bool | 充电口盖是否打开。 |
| `charge_port_latch` | string | 充电口锁止状态。 |
| `charge_rate` | float | 当前充电速率（英里/小时）。 |
| `charger_actual_current` | int | 车辆实际电流。 |
| `charger_phases` | int\|null | 充电相位数。 |
| `charger_pilot_current` | int | 充电桩提供的最大电流。 |
| `charger_power` | int | 充电功率（千瓦）。 |
| `charger_voltage` | int | 充电电压。 |
| `charging_state` | string | 充电状态，如 `Disconnected`。 |
| `conn_charge_cable` | string | 连接的充电电缆类型。 |
| `est_battery_range` | float | 估算续航（英里，基于驾驶行为）。 |
| `fast_charger_brand` | string | 快充品牌。 |
| `fast_charger_present` | bool | 是否连接快充。 |
| `fast_charger_type` | string | 快充类型。 |
| `ideal_battery_range` | float | 理想工况续航。 |
| `managed_charging_active` | bool | 是否处于智能充电模式。 |
| `managed_charging_start_time` | int64\|null | 智能充电计划开始时间（Unix 秒）。 |
| `managed_charging_user_canceled` | bool | 用户是否取消智能充电。 |
| `max_range_charge_counter` | int | 全量充电次数。 |
| `minutes_to_full_charge` | int | 完成充电剩余分钟。 |
| `not_enough_power_to_heat` | bool\|null | 供电不足以开启加热。 |
| `off_peak_charging_enabled` | bool | 是否启用离峰充电。 |
| `off_peak_charging_times` | string | 离峰充电时段。 |
| `off_peak_hours_end_time` | int | 离峰结束时间（分钟）。 |
| `preconditioning_enabled` | bool | 是否启用充电前预热。 |
| `preconditioning_times` | string | 预热计划时段。 |
| `scheduled_charging_mode` | string | 定时充电模式。 |
| `scheduled_charging_pending` | bool | 是否等待定时充电。 |
| `scheduled_charging_start_time` | int64\|null | 定时充电开始时间。 |
| `scheduled_departure_time` | int64 | 计划出发时间（Unix 秒）。 |
| `scheduled_departure_time_minutes` | int | 计划出发时间（分钟）。 |
| `supercharger_session_trip_planner` | bool | 是否来自行程规划快充。 |
| `time_to_full_charge` | float | 完成充电剩余小时。 |
| `timestamp` | int64 | 状态时间戳（毫秒）。 |
| `trip_charging` | bool | 是否处于行程充电模式。 |
| `usable_battery_level` | int | 可用电量百分比。 |
| `user_charge_enable_request` | bool\|null | 用户是否手动请求充电。 |

#### `climate_state` 字段

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `allow_cabin_overheat_protection` | bool | 是否允许座舱过热保护。 |
| `auto_seat_climate_left` | bool | 左座自动座椅加热。 |
| `auto_seat_climate_right` | bool | 右座自动座椅加热。 |
| `auto_steering_wheel_heat` | bool | 方向盘自动加热。 |
| `battery_heater` | bool | 电池加热是否开启。 |
| `battery_heater_no_power` | bool\|null | 电池加热器是否缺少电源。 |
| `bioweapon_mode` | bool | 生化防护模式。 |
| `cabin_overheat_protection` | string | 座舱过热模式。 |
| `cabin_overheat_protection_actively_cooling` | bool | 过热保护是否在制冷。 |
| `climate_keeper_mode` | string | 空调保持模式（off/dog等）。 |
| `cop_activation_temperature` | string | 过热保护触发温度。 |
| `defrost_mode` | int | 除霜模式。 |
| `driver_temp_setting` | float | 驾驶席设定温度。 |
| `fan_status` | int | 风扇档位。 |
| `hvac_auto_request` | string | 自动空调请求状态。 |
| `inside_temp` | float | 车内温度（摄氏）。 |
| `is_auto_conditioning_on` | bool | 自动空调是否运行。 |
| `is_climate_on` | bool | 空调是否开启。 |
| `is_front_defroster_on` | bool | 前挡风除霜是否开启。 |
| `is_preconditioning` | bool | 是否在预处理（预热/预冷）。 |
| `is_rear_defroster_on` | bool | 后挡风除霜是否开启。 |
| `left_temp_direction` | int | 左侧温度调节方向。 |
| `max_avail_temp` | float | 最大可设温度。 |
| `min_avail_temp` | float | 最小可设温度。 |
| `outside_temp` | float | 外部温度。 |
| `passenger_temp_setting` | float | 副驾驶设定温度。 |
| `remote_heater_control_enabled` | bool | 是否支持远程加热。 |
| `right_temp_direction` | int | 右侧温度调节方向。 |
| `seat_heater_left` | int | 驾驶席座椅加热档位。 |
| `seat_heater_rear_center` | int | 后排中间加热档位。 |
| `seat_heater_rear_left` | int | 后排左侧加热档位。 |
| `seat_heater_rear_right` | int | 后排右侧加热档位。 |
| `seat_heater_right` | int | 副驾驶座椅加热档位。 |
| `side_mirror_heaters` | bool | 后视镜加热是否启用。 |
| `steering_wheel_heat_level` | int | 方向盘加热档位。 |
| `steering_wheel_heater` | bool | 方向盘加热是否开启。 |
| `supports_fan_only_cabin_overheat_protection` | bool | 是否支持仅风扇模式过热保护。 |
| `timestamp` | int64 | 状态时间戳（毫秒）。 |
| `wiper_blade_heater` | bool | 雨刷加热是否开启。 |

#### `drive_state` 字段

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `active_route_latitude` | float | 导航路线目的地纬度。 |
| `active_route_longitude` | float | 导航路线目的地经度。 |
| `active_route_traffic_minutes_delay` | int | 导航交通延迟（分钟）。 |
| `gps_as_of` | int64 | GPS 数据时间（秒）。 |
| `heading` | int | 车辆朝向（度）。 |
| `latitude` | float | 当前纬度。 |
| `longitude` | float | 当前经度。 |
| `native_latitude` | float | 本地坐标纬度。 |
| `native_location_supported` | int | 是否支持本地坐标。 |
| `native_longitude` | float | 本地坐标经度。 |
| `native_type` | string | 坐标系类型（如 `wgs`）。 |
| `power` | int | 驱动功率指标。 |
| `shift_state` | string\|null | 档位（P/D/R/N）。 |
| `speed` | float\|null | 速度（英里/小时）。 |
| `timestamp` | int64 | 状态时间戳（毫秒）。 |

#### `gui_settings` 字段

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `gui_24_hour_time` | bool | 是否 24 小时制。 |
| `gui_charge_rate_units` | string | 充电速率单位。 |
| `gui_distance_units` | string | 距离与速度单位。 |
| `gui_range_display` | string | 续航显示模式。 |
| `gui_temperature_units` | string | 温度单位。 |
| `gui_tirepressure_units` | string | 胎压单位。 |
| `show_range_units` | bool | 是否显示续航单位。 |
| `timestamp` | int64 | 状态时间戳（毫秒）。 |

#### `vehicle_config` 字段

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `aux_park_lamps` | string | 辅助泊车灯配置。 |
| `badge_version` | int | 徽章版本。 |
| `can_accept_navigation_requests` | bool | 是否可接受导航下发。 |
| `can_actuate_trunks` | bool | 是否支持电动前后备箱。 |
| `car_special_type` | string | 特殊类型标识。 |
| `car_type` | string | 车型。 |
| `charge_port_type` | string | 充电口类型。 |
| `cop_user_set_temp_supported` | bool | 是否支持用户设置过热保护温度。 |
| `dashcam_clip_save_supported` | bool | 是否支持行车记录仪手动保存。 |
| `default_charge_to_max` | bool | 默认是否充满。 |
| `driver_assist` | string | 驾驶辅助硬件版本。 |
| `ece_restrictions` | bool | 是否启用欧盟限制。 |
| `efficiency_package` | string | 效能包配置。 |
| `eu_vehicle` | bool | 是否欧规车辆。 |
| `exterior_color` | string | 外观颜色。 |
| `exterior_trim` | string | 外部饰件。 |
| `exterior_trim_override` | string | 外饰覆盖配置。 |
| `has_air_suspension` | bool | 是否空气悬挂。 |
| `has_ludicrous_mode` | bool | 是否支持狂暴模式。 |
| `has_seat_cooling` | bool | 是否支持座椅通风。 |
| `headlamp_type` | string | 头灯类型。 |
| `interior_trim_type` | string | 内饰类型。 |
| `key_version` | int | 钥匙版本。 |
| `motorized_charge_port` | bool | 充电口是否电动。 |
| `paint_color_override` | string | 涂装覆盖设置。 |
| `performance_package` | string | 性能包配置。 |
| `plg` | bool | 是否支持电动尾门（power lift gate）。 |
| `pws` | bool | 是否具备行人警示音。 |
| `rear_drive_unit` | string | 后驱动单元型号。 |
| `rear_seat_heaters` | int | 后排加热座椅数量。 |
| `rear_seat_type` | int | 后排座椅类型。 |
| `rhd` | bool | 是否右舵。 |
| `roof_color` | string | 车顶颜色。 |
| `seat_type` | int\|null | 座椅材质/类型。 |
| `spoiler_type` | string | 尾翼配置。 |
| `sun_roof_installed` | int\|null | 天窗配置。 |
| `supports_qr_pairing` | bool | 是否支持二维码配对。 |
| `third_row_seats` | string | 第三排配置。 |
| `timestamp` | int64 | 状态时间戳（毫秒）。 |
| `trim_badging` | string | 装饰徽章。 |
| `use_range_badging` | bool | 是否显示续航徽章。 |
| `utc_offset` | int | 时区偏移（秒）。 |
| `webcam_selfie_supported` | bool | 是否支持车内摄像头自拍。 |
| `webcam_supported` | bool | 是否支持车载摄像头。 |
| `wheel_type` | string | 轮毂类型。 |

#### `vehicle_state` 字段

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `api_version` | int | 车辆 API 版本。 |
| `autopark_state_v3` | string | 自动泊车状态。 |
| `autopark_style` | string | 自动泊车模式。 |
| `calendar_supported` | bool | 是否支持日历同步。 |
| `car_version` | string | 固件版本。 |
| `center_display_state` | int | 中控显示状态。 |
| `dashcam_clip_save_available` | bool | 行车记录仪保存功能可用。 |
| `dashcam_state` | string | 行车记录仪状态。 |
| `df` | int | 驾驶员前门状态（0 关闭）。 |
| `dr` | int | 驾驶员后门状态。 |
| `fd_window` | int | 前左窗状态。 |
| `feature_bitmask` | string | 功能位掩码。 |
| `fp_window` | int | 前右窗状态。 |
| `ft` | int | 前备箱状态。 |
| `homelink_device_count` | int | HomeLink 设备数量。 |
| `homelink_nearby` | bool | 是否检测到 HomeLink。 |
| `is_user_present` | bool | 是否检测到驾乘人员。 |
| `last_autopark_error` | string | 最近自动泊车错误。 |
| `locked` | bool | 车辆是否上锁。 |
| `media_info` | object | 当前媒体播放信息，见下表。 |
| `media_state` | object | 媒体控制状态，见下表。 |
| `notifications_supported` | bool | 是否支持通知推送。 |
| `odometer` | float | 里程（英里）。 |
| `parsed_calendar_supported` | bool | 是否支持日历解析。 |
| `pf` | int | 乘客前门状态。 |
| `pr` | int | 乘客后门状态。 |
| `rd_window` | int | 后左窗状态。 |
| `remote_start` | bool | 是否开启远程启动。 |
| `remote_start_enabled` | bool | 是否支持远程启动。 |
| `remote_start_supported` | bool | 车辆是否具备远程启动功能。 |
| `rp_window` | int | 后右窗状态。 |
| `rt` | int | 后备箱状态。 |
| `santa_mode` | int | 圣诞模式状态。 |
| `sentry_mode` | bool | 哨兵模式是否开启。 |
| `sentry_mode_available` | bool | 是否支持哨兵模式。 |
| `service_mode` | bool | 是否处于维护模式。 |
| `service_mode_plus` | bool | 维护模式增强版。 |
| `smart_summon_available` | bool | 是否支持智能召唤。 |
| `software_update` | object | 固件更新信息，见下表。 |
| `speed_limit_mode` | object | 限速模式设置，见下表。 |
| `summon_standby_mode_enabled` | bool | 是否启用召唤待命模式。 |
| `timestamp` | int64 | 状态时间戳（毫秒）。 |
| `tpms_hard_warning_fl` | bool | 前左胎压硬告警。 |
| `tpms_hard_warning_fr` | bool | 前右胎压硬告警。 |
| `tpms_hard_warning_rl` | bool | 后左胎压硬告警。 |
| `tpms_hard_warning_rr` | bool | 后右胎压硬告警。 |
| `tpms_last_seen_pressure_time_fl` | int64 | 前左胎压最后更新时间（秒）。 |
| `tpms_last_seen_pressure_time_fr` | int64 | 前右胎压最后更新时间。 |
| `tpms_last_seen_pressure_time_rl` | int64 | 后左胎压最后更新时间。 |
| `tpms_last_seen_pressure_time_rr` | int64 | 后右胎压最后更新时间。 |
| `tpms_pressure_fl` | float | 前左胎压（巴）。 |
| `tpms_pressure_fr` | float | 前右胎压。 |
| `tpms_pressure_rl` | float | 后左胎压。 |
| `tpms_pressure_rr` | float | 后右胎压。 |
| `tpms_rcp_front_value` | float | 前轴推荐胎压。 |
| `tpms_rcp_rear_value` | float | 后轴推荐胎压。 |
| `tpms_soft_warning_fl` | bool | 前左胎压软告警。 |
| `tpms_soft_warning_fr` | bool | 前右胎压软告警。 |
| `tpms_soft_warning_rl` | bool | 后左胎压软告警。 |
| `tpms_soft_warning_rr` | bool | 后右胎压软告警。 |
| `valet_mode` | bool | 是否开启代客模式。 |
| `valet_pin_needed` | bool | 是否需要代客 PIN。 |
| `vehicle_name` | string | 车辆昵称。 |
| `vehicle_self_test_progress` | int | 自检进度。 |
| `vehicle_self_test_requested` | bool | 是否请求自检。 |
| `webcam_available` | bool | 车载摄像头是否可用。 |

**媒体相关嵌套对象**

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `media_info.a2dp_source_name` | string | 蓝牙音源名称。 |
| `media_info.audio_volume` | float | 当前音量。 |
| `media_info.audio_volume_increment` | float | 音量步进值。 |
| `media_info.audio_volume_max` | float | 最大音量。 |
| `media_info.media_playback_status` | string | 播放状态。 |
| `media_info.now_playing_album` | string | 当前专辑。 |
| `media_info.now_playing_artist` | string | 当前艺人。 |
| `media_info.now_playing_duration` | int | 当前音频总时长。 |
| `media_info.now_playing_elapsed` | int | 已播放时长。 |
| `media_info.now_playing_source` | string | 音源编号。 |
| `media_info.now_playing_station` | string | 当前电台/频道。 |
| `media_info.now_playing_title` | string | 当前曲目标题。 |
| `media_state.remote_control_enabled` | bool | 是否允许远程媒体控制。 |

**固件更新与限速对象**

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `software_update.download_perc` | int | 下载进度百分比。 |
| `software_update.expected_duration_sec` | int | 预计安装时长（秒）。 |
| `software_update.install_perc` | int | 安装进度百分比。 |
| `software_update.status` | string | 更新状态。 |
| `software_update.version` | string | 目标版本。 |
| `speed_limit_mode.active` | bool | 是否启用限速模式。 |
| `speed_limit_mode.current_limit_mph` | int | 当前限速（英里/小时）。 |
| `speed_limit_mode.max_limit_mph` | int | 最高可设限速。 |
| `speed_limit_mode.min_limit_mph` | int | 最低可设限速。 |
| `speed_limit_mode.pin_code_set` | bool | 是否设定解锁 PIN。 |

### 示例

```json
{
  "response": {
    "id": 100021,
    "user_id": 800001,
    "vehicle_id": 99999,
    "vin": "TEST000000VIN01",
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
      "gui_distance_units": "mi/hr",
      "gui_temperature_units": "F",
      "gui_24_hour_time": false
    },
    "vehicle_config": {
      "car_type": "modely",
      "exterior_color": "MidnightSilver",
      "wheel_type": "Apollo19"
    },
    "vehicle_state": {
      "car_version": "2023.7.20 7910d26d5c64",
      "locked": true,
      "sentry_mode": false
    }
  }
}
```

## GET /api/1/vehicles/{vehicle_tag}/drivers

- **方法**：GET  
- **授权范围**：`vehicle_device_data`  
- **描述**：返回车辆所有已授权驾驶员，仅限车主管理使用。

### 路径参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
| ---- | ---- | ---- | ---- | ---- |
| `vehicle_tag` | string | 是 | 车辆唯一标识，可为 `id`、`id_s` 或 `vin`。 | `100021` |

### 响应体

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `response` | `VehicleDriver[]` | 驾驶员数组，字段见下表。 |
| `count` | int | 返回的驾驶员数量。 |

#### VehicleDriver

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `my_tesla_unique_id` | int64 | 驾驶员在 MyTesla 中的唯一编号。 |
| `user_id` | int64 | 驾驶员 Tesla 账户 ID。 |
| `user_id_s` | string | 字符串形式账户 ID。 |
| `vault_uuid` | string | 与驾驶员授权关联的 Vault 标识。 |
| `driver_first_name` | string | 驾驶员名。 |
| `driver_last_name` | string | 驾驶员姓。 |
| `granular_access.hide_private` | bool | 是否隐藏驾驶员的隐私数据（如位置信息）。 |
| `active_pubkeys` | string[] | 激活的公钥列表，用于密钥交换。 |
| `public_key` | string | 当前绑定的主公钥。 |

### 示例

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

> 注：示例数据来源于官方文档，仅用于字段说明；实际取值以线上返回为准。

## POST /api/1/vehicles/{vehicle_tag}/wake_up

- **方法**：POST  
- **授权范围**：`vehicle_device_data`  
- **描述**：使用 Tesla Fleet API 唤醒指定车辆。接口会立即向 Tesla 后端发送唤醒请求，成功返回车辆概要信息，若车辆仍处于离线/睡眠状态会返回错误。

### 路径参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
| ---- | ---- | ---- | ---- | ---- |
| `vehicle_tag` | string | 是 | 车辆唯一标识，支持 `id`、`id_s` 或 `vin`。 | `100021` |

### 请求体

无。

### 响应体

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `response` | `VehicleSummary` | Tesla 返回的车辆概要信息，字段同上。 |

### 错误示例

- `400 invalid vin`：`vehicle_tag` 无效或未授权访问。
- `500 vehicle unavailable`：车辆仍处于离线/睡眠状态，需稍后重试或通过官方 App/实体操作唤醒。

## GET /api/1/vehicles/{vehicle_tag}

- **方法**：GET  
- **授权范围**：`vehicle_device_data`  
- **描述**：根据 `vehicle_tag`（车辆 id 或 VIN）返回单辆车最新的概要信息。

### 路径参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
| ---- | ---- | ---- | ---- | ---- |
| `vehicle_tag` | string | 是 | 车辆唯一标识，可在列表接口中获取 `id`、`id_s` 或 `vin`。 | `100021` |

### 响应体

| 字段 | 类型 | 说明 |
| ---- | ---- | ---- |
| `response` | `VehicleSummary` | 单辆车概要信息，字段含义与列表接口相同。 |

### 示例

```json
{
  "response": {
    "id": 100021,
    "vehicle_id": 99999,
    "vin": "TEST000000VIN01",
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
