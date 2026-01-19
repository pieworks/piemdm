# API 參考

PieMDM 後端 API 是平台的核心，提供主數據管理能力。

## 基礎 URL

預設情況下，API 可透過以下網址存取：

```
http://127.0.0.1:8787/swagger/index.html
```

## 技術棧

後端建構基於：
- **語言**: Go
- **框架**: Gin
- **資料庫**: GORM (MySQL)
- **快取**: Redis

## 端點

> [!NOTE]
> 詳細的 API 文件 (Swagger/OpenAPI) 通常由程式碼生成。請查看後端儲存庫中的 `swagger` 或 `docs` 目錄。

常見模組通常包括：
- **Auth**: 使用者認證與 Token 管理。
- **User**: 使用者資料與管理。
- **System**: 系統配置與狀態。
