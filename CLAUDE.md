# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码库中工作提供指导。

## 开发命令

### 构建和运行
- **构建**: `go build -o ./bin/server ./cmd/server`
- **本地运行**: `go run ./cmd/server`
- **Docker 构建**: `docker build -t tds_server .`
- **Docker 运行**: `docker run -p 8080:8080 tds_server`

### 依赖管理
- **安装/更新依赖**: `go mod tidy`
- **下载依赖**: `go mod download`

## 架构概览

这是一个基于 Go 的特斯拉数据服务 (TDS) 服务器，提供特斯拉 API 集成和 OAuth 认证流程。应用程序遵循清洁架构模式，具有清晰的关注点分离。

### 核心组件

**入口点**: `cmd/server/main.go` - 应用程序引导、依赖注入和服务器启动

**配置管理**: `internal/config/config.go` - 使用环境变量和 .env 文件的集中配置管理。处理特斯拉 API 凭据、数据库连接和服务器设置。

**数据层**: 
- `internal/data/db.go` - 使用 GORM 和 PostgreSQL 的数据库初始化和连接管理
- `internal/model/` - 数据库模型和实体
- `internal/repository/` - 使用仓储模式的数据访问层

**API 层**:
- `internal/router/router.go` - Gin HTTP 路由器设置和路由定义
- `internal/handler/` - 特斯拉 OAuth 流程和车辆数据的 HTTP 请求处理器
- `internal/middleware/` - HTTP 中间件组件

**业务逻辑**: `internal/service/` - 特斯拉 API 集成和业务逻辑

**工具类**: `internal/util/` - 共享工具函数

### 主要依赖
- **Gin**: HTTP API 的 Web 框架
- **GORM**: 使用 PostgreSQL 驱动的数据库 ORM
- **Resty**: 特斯拉 API 调用的 HTTP 客户端
- **godotenv**: 从 .env 文件加载环境变量

### 特斯拉集成
应用程序实现特斯拉的 OAuth 2.0 流程：
1. `/api/login` - 重定向到特斯拉授权页面
2. `/api/login/callback` - 处理 OAuth 回调和令牌交换
3. `/api/list` - 使用存储的令牌获取车辆数据

### 环境配置
必需的环境变量：
- `TESLA_CLIENT_ID`, `TESLA_CLIENT_SECRET`, `TESLA_REDIRECT_URI`
- `TESLA_AUTH_URL`, `TESLA_TOKEN_URL`, `TESLA_API_URL`
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`

### 数据库
使用 PostgreSQL 和 GORM 作为 ORM。令牌仓库管理特斯拉 API 访问的用户认证令牌。

### 部署
使用多阶段构建的 Docker 化应用程序。最终镜像在 Alpine Linux 上运行编译好的 Go 二进制文件，暴露 8080 端口。