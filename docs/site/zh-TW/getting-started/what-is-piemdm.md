# 什麼是 PieMDM？

PieMDM 是一個全棧主數據管理（Master Data Management）平台，旨在幫助企業集中管理和維護關鍵業務數據。

## 項目概覽

PieMDM 採用現代化的前後端分離架構，提供高效、穩定且易於擴展的主數據管理解決方案。

- **後端**：基於 Go 語言構建，利用 Gin 框架提供高性能的 API 服務。
- **前端**：基於 Vue.js 3 和 Vite 構建，提供響應式且用戶友好的 Web 界面。
- **容器化**：支持 Docker 和 Docker Compose，便於開發和部署。

## 技術棧

### 後端 (Backend)

後端代碼位於 `backend/` 目錄，主要技術選型包括：

- **語言**：Go
- **Web 框架**：Gin
- **數據庫 ORM**：GORM (MySQL)
- **緩存**：Redis
- **依賴注入**：google/wire
- **測試**：testify, go-sqlmock

### 前端 (Frontend)

前端代碼位於 `frontend/` 目錄，主要技術選型包括：

- **框架**：Vue.js 3
- **構建工具**：Vite
- **狀態管理**：Pinia
- **路由**：Vue Router
- **UI 框架**：Bootstrap 5
- **包管理器**：pnpm

## 快速開始

PieMDM 推薦使用 Docker Compose 進行快速啟動：

```bash
# 進入部署目錄
cd deploy
# 在後台啟動所有服務
docker-compose up -d
```

啟動後，您可以訪問：
- **API 服務**：`http://localhost:8787`
- **前端頁面**：`http://localhost:80`

## 獲取更多

- [快速開始指南](./getting-started)
- [API 參考](../reference/open-api)
