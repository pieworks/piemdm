# 快速开始

要搭建开发环境，您可以使用 **Docker Compose**（推荐用于快速启动）或 **手动** 搭建环境。

## 先决条件

请确保已安装以下软件：
- **Go**: 1.24.12+
- **Node.js**: 20+
- **MySQL**: 8.0+
- **Redis**: 6+

## 快速启动 (Docker Compose)

在后台启动所有服务：

```bash
cp deploy/.env.example .env
# 编辑 .env 以设置 MySQL 连接串和 Redis 地址
vim .env

docker-compose -f deploy/docker-compose.yml up -d
```

这将构建并启动 **后端 API** 和 **前端** 容器。

> [!IMPORTANT]
> 此 Docker Compose 配置 **不包含** MySQL 和 Redis。您必须确保它们在外部运行，并通过环境变量进行配置。

- **前端**: `http://localhost:8081`
- **后端 API**: `http://localhost:8787`
- **API 文档**: `http://localhost:8787/swagger/index.html`

## 手动安装

### 1. 后端

后端是一个位于 `backend/` 目录的 Go 应用程序。

```bash
cd backend

# 安装依赖
go mod tidy

# 复制配置文件
cp config/local.yml.example config/local.yml
# 编辑 config/local.yml 以配置数据库和其他设置

# 运行数据库迁移
go run cmd/migration/main.go

# 启动后端服务
go run cmd/server/main.go
```

### 2. 前端

前端是一个位于 `frontend/` 目录的 Vue.js 应用程序。

```bash
cd frontend

# 安装依赖
npm install
# 或
pnpm install

# 启动开发服务器
npm run dev
# 或
pnpm dev
```

### 3. 数据库设置

如果是手动运行，请确保已创建数据库：

```sql
CREATE DATABASE piemdm CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## 配置

### 环境变量 (.env)

在根目录下创建一个 `.env` 文件用于全局设置：

```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=piemdm

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Application
APP_PORT=8787
```
