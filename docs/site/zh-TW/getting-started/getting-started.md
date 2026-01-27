# 快速開始

本指南將幫助您搭建 **PieMDM** 的開發環境。

## 先決條件

請確保您的本地機器已安裝以下軟體：

- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- [Go](https://go.dev/) (>= 1.20) 用於後端開發
- [Node.js](https://nodejs.org/) (>= 20.10.0) & [pnpm](https://pnpm.io/) (>= 9.0.0) 用於前端開發

## 快速啟動 (Docker)

最簡單的啟動完整技術棧的方法是使用 Docker Compose。這將啟動 API、Web、MySQL 和 Redis 服務。

```bash
# 進入部署目錄
cd deploy
# 在後台啟動所有服務
docker-compose up -d
```

啟動後，您可以訪問：
- **Web 介面**: [http://localhost:80](http://localhost:80)
- **API 服務器**: [http://localhost:8787](http://localhost:8787)

## 手動搭建

如果您更喜歡單獨運行服務進行開發：

### 1. 數據庫 & 緩存

您仍然需要 MySQL 和 Redis。可以使用 Docker 啟動它們：

```bash
docker-compose up -d mysql redis
```

### 2. 後端 (API)
後端使用 Go 語言開發，基於 Gin 框架。
進入 `api` 目錄：

```bash
cd api
```

安裝依賴並運行服務器（使用 `air` 啟用熱重載）：

```bash
make dev
```

常用命令：
- `make test`: 運行單元測試
- `make lint`: 運行代碼檢查
- `make wire`: 生成依賴注入代碼

### 3. 前端 (Web)
前端使用 Vue 3 開發。
進入 `web` 目錄：

```bash
cd web
```

安裝依賴並啟動開發服務器：

```bash
pnpm install
pnpm dev
```

如需從其他設備訪問，請傳遞 `--host` 參數：

```bash
pnpm dev --host
```

## 專案結構

- **`api/`**: 後端應用 (Go, Gin, GORM)
- **`web/`**: 前端應用 (Vue 3, Vite, Bootstrap 5)
- **`docs/`**: 專案文檔
