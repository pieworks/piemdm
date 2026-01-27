# 部署指南

PieMDM 旨在通過容器化技術輕鬆部署。我們推薦使用 Docker Compose 進行生產環境部署。

## 先決條件

部署環境需要安裝：

- [Docker](https://docs.docker.com/engine/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## 使用 Docker Compose 部署

1. **獲取代碼**

   將項目代碼克隆到服務器：

   ```bash
   git clone https://github.com/pieworks/piemdm.git
   cd piemdm
   ```

2. **配置**

   根據環境需要修改 `docker-compose.yaml` 或相關的配置文件（如數據庫密碼、端口映射等）。

3. **啟動服務**

   在項目根目錄下運行：

   ```bash
   cd deploy
   docker-compose up -d
   ```

   該命令將構建並啟動所有必要的服務容器（API, Web）。

4. **驗證**

   使用 `docker-compose ps` 檢查服務狀態。
   
   - 前端頁面默認端口：`80`
   - API 服務默認端口：`8787`

## 手動構建部署

如果您需要手動構建二進制文件和靜態資源：

### 後端

```bash
cd api
make build
# 運行生成的二進制文件
./bin/server
```

### 前端

```bash
cd web
pnpm install
pnpm build
# 構建產物位於 web/dist 目錄，可使用 Nginx 等 Web 服務器進行託管
```
