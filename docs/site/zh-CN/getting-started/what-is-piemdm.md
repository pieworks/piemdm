---
head:
  - - meta
    - name: description
      content: 什么是 PieMDM？
  - - meta
    - name: keywords
      content: 主数据,MDM,PieMDM,开源,Go,Vue,Docker
---

# 什么是 PieMDM？

PieMDM 是一个全栈主数据管理（Master Data Management）平台，旨在帮助企业集中管理和维护关键业务数据。

## 项目概览

PieMDM 采用现代化的前后端分离架构，提供高效、稳定且易于扩展的主数据管理解决方案。

- **后端**：基于 Go 语言构建，利用 Gin 框架提供高性能的 API 服务。
- **前端**：基于 Vue.js 3 和 Vite 构建，提供响应式且用户友好的 Web 界面。
- **容器化**：支持 Docker 和 Docker Compose，便于开发和部署。

## 技术栈

### 后端 (API)

后端代码位于 `api/` 目录，主要技术选型包括：

- **语言**：Go
- **Web 框架**：Gin
- **数据库 ORM**：GORM (MySQL)
- **缓存**：Redis
- **依赖注入**：google/wire
- **测试**：testify, go-sqlmock

### 前端 (Web)

前端代码位于 `web/` 目录，主要技术选型包括：

- **框架**：Vue.js 3
- **构建工具**：Vite
- **状态管理**：Pinia
- **路由**：Vue Router
- **UI 框架**：Bootstrap 5
- **包管理器**：pnpm

## 快速开始

PieMDM 推荐使用 Docker Compose 进行快速启动：

```bash
# 进入部署目录
cd deploy
# 在后台启动所有服务
docker-compose up -d
```

启动后，您可以访问：
- **API 服务**：`http://localhost:8787`
- **前端页面**：`http://localhost:80`

## 获取更多

- [快速开始指南](./getting-started)
- [API 参考](../reference/open-api)