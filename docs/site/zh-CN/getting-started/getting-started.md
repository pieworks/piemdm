---
head:
  - - meta
    - name: description
      content: 快速开始
  - - meta
    - name: keywords
      content: 主数据,MDM,PieMDM,开源,Go,Vue,Docker
---

# 快速开始

本指南将帮助您搭建 **PieMDM** 的开发环境。

## 先决条件

请确保您的本地机器已安装以下软件：

- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- [Go](https://go.dev/) (>= 1.20) 用于后端开发
- [Node.js](https://nodejs.org/) (>= 20.10.0) & [pnpm](https://pnpm.io/) (>= 9.0.0) 用于前端开发

## 快速启动 (Docker)

最简单的启动完整技术栈的方法是使用 Docker Compose。这将启动 API、Web、MySQL 和 Redis 服务。

```bash
# 进入部署目录
cd deploy
# 在后台启动所有服务
docker-compose up -d
```

启动后，您可以访问：
- **Web 界面**: [http://localhost:80](http://localhost:80)
- **API 服务器**: [http://localhost:8787](http://localhost:8787)

## 手动搭建

如果您更喜欢单独运行服务进行开发：

### 1. 数据库 & 缓存

您仍然需要 MySQL 和 Redis。可以使用 Docker 启动它们：

```bash
docker-compose up -d mysql redis
```

### 2. 后端 (API)
后端使用 Go 语言开发，基于 Gin 框架。
进入 `api` 目录：

```bash
cd api
```

安装依赖并运行服务器（使用 `air` 启用热重载）：

```bash
make dev
```

常用命令：
- `make test`: 运行单元测试
- `make lint`: 运行代码检查
- `make wire`: 生成依赖注入代码

### 3. 前端 (Web)
前端使用 Vue 3 开发。
进入 `web` 目录：

```bash
cd web
```

安装依赖并启动开发服务器：

```bash
pnpm install
pnpm dev
```

如需从其他设备访问，请传递 `--host` 参数：

```bash
pnpm dev --host
```

## 项目结构

- **`api/`**: 后端应用 (Go, Gin, GORM)
- **`web/`**: 前端应用 (Vue 3, Vite, Bootstrap 5)
- **`docs/`**: 项目文档
