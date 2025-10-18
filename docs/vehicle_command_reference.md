# Tesla Vehicle Command SDK 指令速查表

本项目在 `internal/service/vehicle_command_service.go` 中集成了官方开源库 `github.com/teslamotors/vehicle-command`，以下整理了当前 SDK 内置的所有可调用指令、参数要求以及回退策略，便于在接入层或调试时快速查阅。

> **执行流程提示**
>
> 1. 服务会先尝试通过 SDK 使用新协议下发指令。
> 2. 若 SDK 返回 `ErrVehicleCommandUseREST`，才会自动回退到旧版 Fleet API REST 接口。
> 3. 部分指令尚未在 SDK 中实现或明确要求使用 REST，见文末“不支持/需回退”小节。

## 参数约定

- `seat_position`（加热）：索引表 `0=前排左, 1=前排右, 2=第二排左, 3=第二排左后, 4=第二排中间, 5=第二排右, 6=第二排右后, 7=第三排左, 8=第三排右`。
- `seat_position`（制冷）：仅支持官方枚举 `0=前排左, 1=前排右`。
- `auto_seat_position`：`0=前排左, 1=前排右`。
- `days_of_week`：大小写不敏感、以逗号分隔，可用别名（如 `SUN`/`Sunday`、`ALL`、`WEEKDAYS` 等）。
- `time` / `departure_time` / `end_off_peak_time`：传入分钟数，服务端会转换成 `time.Duration`。
- `offset_sec`：以秒为单位的延迟时间。

## 指令列表

### 媒体控制

| 指令 | 说明 | 必填参数 | 可选参数 |
| --- | --- | --- | --- |
| `adjust_volume` | 设置音量百分比 | `volume` (float) | - |
| `media_next_fav` / `media_prev_fav` | 下一/上一收藏电台 | - | - |
| `media_next_track` / `media_prev_track` | 下一/上一曲目 | - | - |
| `media_volume_up` / `media_volume_down` | 音量增减一级 | - | - |
| `media_toggle_playback` | 播放/暂停切换 | - | - |
| `remote_boombox` | 暴风音效 | - | ❌ **尚未实现** |

### 空调与舒适

| 指令 | 说明 | 必填参数 | 可选参数 |
| --- | --- | --- | --- |
| `auto_conditioning_start` / `auto_conditioning_stop` | 开/关空调 | - | - |
| `charge_max_range` | 设置电量最大里程模式 | - | - |
| `remote_seat_heater_request` | 座椅加热 | `seat_position` (int), `level` (0-3) | - |
| `remote_seat_cooler_request` | 座椅通风 | `seat_position` (int), `seat_cooler_level` (1-3) | - |
| `remote_auto_seat_climate_request` | 自动座椅空调 | `auto_seat_position` (int), `auto_climate_on` (bool) | - |
| `remote_steering_wheel_heater_request` | 方向盘加热 | `on` (bool) | - |
| `set_bioweapon_mode` | 生化防御模式 | `on` (bool), `manual_override` (bool) | - |
| `set_cabin_overheat_protection` | 过热保护 | `on` (bool) | `fan_only` (bool) |
| `set_climate_keeper_mode` | Cabin Keeper 模式 | `climate_keeper_mode` (0=Off,1=On,2=Dog,3=Camp) | `manual_override` (bool) |
| `set_cop_temp` | 过热保护温度 | `cop_temp` (float) | - |
| `set_preconditioning_max` | 预处理最大功率 | `on` (bool) | `manual_override` (bool) |
| `set_temps` | 空调目标温度 | - | `driver_temp` (float), `passenger_temp` (float) |

### 车身控制

| 指令 | 说明 | 必填参数 | 可选参数 |
| --- | --- | --- | --- |
| `actuate_trunk` | 开启前/后备箱 | - | `which_trunk` (`front`/`rear`，默认后备箱) |
| `charge_port_door_open` / `charge_port_door_close` | 开关充电口盖 | - | - |
| `flash_lights` / `honk_horn` | 闪灯 / 鸣笛 | - | - |
| `remote_start_drive` | 远程启动 | - | - |
| `open_tonneau` / `close_tonneau` / `stop_tonneau` | 货箱盖开合/停止 | - | - |
| `wake_up` | 唤醒车辆 | - | - |
| `window_control` | 车窗开/关 | `command` (`vent`/`close`) | `lat`, `lon`（新协议可省略） |

### 充电相关

| 指令 | 说明 | 必填参数 | 可选参数 |
| --- | --- | --- | --- |
| `charge_standard` / `charge_start` / `charge_stop` | 标准充电 / 开始 / 停止充电 | - | - |
| `set_charging_amps` | 设置充电电流 | `charging_amps` (int) | - |
| `set_charge_limit` | 设置电量上限 | `percent` (int) | - |
| `set_scheduled_charging` | 预约充电 | `enable` (bool) | `time` (分钟) |
| `set_scheduled_departure` | 智能出发 | `enable` (bool) | `off_peak_charging_enabled` / `off_peak_charging_weekdays_only` / `preconditioning_enabled` / `preconditioning_weekdays_only` (bool)，`departure_time` / `end_off_peak_time` (分钟) |

### 充电/空调日程

| 指令 | 说明 | 必填参数 | 可选参数 |
| --- | --- | --- | --- |
| `add_charge_schedule` | 新增充电日程 | `lat`, `lon` (float), `start_enabled`, `end_enabled`, `days_of_week`, `enabled` (bool) | `start_time`, `end_time` (分钟), `id` (int，默认当前时间戳), `one_time` (bool) |
| `add_precondition_schedule` | 新增预处理日程 | `lat`, `lon` (float), `precondition_time` (分钟), `days_of_week`, `enabled` (bool) | `id` (int), `one_time` (bool) |
| `remove_charge_schedule` / `remove_precondition_schedule` | 删除日程 | `id` (int/uint64) | - |

### 安全与防护

| 指令 | 说明 | 必填参数 | 可选参数 |
| --- | --- | --- | --- |
| `set_pin_to_drive` | 启用驾驶 PIN | `on` (bool) | `password` (string) |
| `clear_pin_to_drive_admin` / `reset_pin_to_drive_pin` | 清除/重置驾驶 PIN | - | - |
| `reset_valet_pin` | 重置代客 PIN | - | - |
| `door_lock` / `door_unlock` | 车门锁定/解锁 | - | - |
| `set_valet_mode` | 设置代客模式 | `on` (bool) | `password` (string) |
| `guest_mode` | 访客模式 | `enable` (bool) | - |
| `set_sentry_mode` | 哨兵模式 | `on` (bool) | - |
| `set_vehicle_name` | 修改车辆名称 | `vehicle_name` (string) | - |
| `speed_limit_activate` / `speed_limit_deactivate` / `speed_limit_clear_pin` | 限速启动/关闭/清除 PIN | `pin` (string) | - |
| `speed_limit_clear_pin_admin` | 管理员清除限速 PIN | - | - |
| `speed_limit_set_limit` | 设置限速值 (mph) | `limit_mph` (float) | - |
| `trigger_homelink` | 触发 Homelink | `lat`, `lon` (float) | - |
| `erase_user_data` | 擦除访客数据 | - | - |

### 软件与其他

| 指令 | 说明 | 必填参数 | 可选参数 |
| --- | --- | --- | --- |
| `schedule_software_update` | 预约安装更新 | `offset_sec` (int) | - |
| `cancel_software_update` | 取消软件更新 | - | - |

## 不支持 / 需回退到 REST

| 指令 | 原因 | 处理方式 |
| --- | --- | --- |
| `remote_boombox` | SDK 未实现 | 直接返回错误 |
| `set_managed_charge_current_request` | 官方标记需 REST | 自动回退 REST |
| `set_managed_charger_location` | 官方标记需 REST | 自动回退 REST |
| `set_managed_scheduled_charging_time` | 官方标记需 REST | 自动回退 REST |
| `navigation_request` | 依赖服务器处理 | 自动回退 REST |

> 更新指令列表时，请同步检视 `github.com/teslamotors/vehicle-command/pkg/proxy/command.go` 的 `ExtractCommandAction` 分支，并确保本文档保持一致。
