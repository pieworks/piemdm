# 快速開始

建議使用 **Docker Compose**（推薦用於快速啟動）或 **手動** 搭建環境。

## 先決條件

請確保已安裝以下軟體：
- **Go**: 1.24.12+
- **Node.js**: 20+
- **MySQL**: 8.0+
- **Redis**: 6+

## 快速啟動 (Docker Compose)

在背景啟動所有服務：

```bash
cp deploy/.env.example .env
# 編輯 .env 以設定 MySQL 連線字串和 Redis 地址
vim .env

docker-compose -f deploy/docker-compose.yml up -d
```

這將建置並啟動 **後端 API** 和 **前端** 容器。

> [!IMPORTANT]
> 此 Docker Compose 配置 **不包含** MySQL 和 Redis。您必須確保它們在外部運行，並透過環境變數進行配置。

- **前端**: `http://localhost:8081`
- **後端 API**: `http://localhost:8787`
- **API 文件**: `http://localhost:8787/swagger/index.html`

## 手動安裝

### 1. 後端

後端是一個位於 `backend/` 目錄的 Go 應用程式。

```bash
cd backend

# 安裝依賴
go mod tidy

# 複製設定檔
cp config/local.yml.example config/local.yml
# 編輯 config/local.yml 以配置資料庫和其他設定

# 執行資料庫遷移
go run cmd/migration/main.go

# 啟動後端服務
go run cmd/server/main.go
```

### 2. 前端

前端是一個位於 `frontend/` 目錄的 Vue.js 應用程式。

```bash
cd frontend

# 安裝依賴
npm install
# 或
pnpm install

# 啟動開發伺服器
npm run dev
# 或
pnpm dev
```

### 3. 資料庫設定

如果是手動執行，請確保已建立資料庫：

```sql
CREATE DATABASE piemdm CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## 配置

### 環境變數 (.env)

在根目錄下建立一個 `.env` 檔案用於全域設定：

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
