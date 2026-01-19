# API 参考

PieMDM 后端 API 是平台的核心，提供主数据管理能力。

## 基础 URL

默认情况下，API 可通过以下地址访问：

```
http://127.0.0.1:8787/swagger/index.html
```

## 技术栈

后端构建基于：
- **语言**: Go
- **框架**: Gin
- **数据库**: GORM (MySQL)
- **缓存**: Redis

## 端点

> [!NOTE]
> 详细的 API 文档 (Swagger/OpenAPI) 通常由代码生成。请查看后端仓库中的 `swagger` 或 `docs` 目录。

常见模块通常包括：
- **Auth**: 用户认证与 Token 管理。
- **User**: 用户资料与管理。
- **System**: 系统配置与状态。
