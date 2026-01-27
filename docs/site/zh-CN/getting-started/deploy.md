# 部署指南

PieMDM 旨在通过容器化技术轻松部署。我们推荐使用 Docker Compose 进行生产环境部署。

## 先决条件

部署环境需要安装：

- [Docker](https://docs.docker.com/engine/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## 使用 Docker Compose 部署

1. **获取代码**

   将项目代码克隆到服务器：

   ```bash
   git clone https://github.com/pieworks/piemdm.git
   cd piemdm
   ```

2. **配置**

   根据环境需要修改 `docker-compose.yaml` 或相关的配置文件（如数据库密码、端口映射等）。

3. **启动服务**

   在项目根目录下运行：

   ```bash
   cd deploy
   docker-compose up -d
   ```

   该命令将构建并启动所有必要的服务容器（API, Web）。

4. **验证**

   使用 `docker-compose ps` 检查服务状态。
   
   - 前端页面默认端口：`80`
   - API 服务默认端口：`8787`

## 手动构建部署

如果您需要手动构建二进制文件和静态资源：

### 后端

```bash
cd api
make build
# 运行生成的二进制文件
./bin/server
```

### 前端

```bash
cd web
pnpm install
pnpm build
# 构建产物位于 web/dist 目录，可使用 Nginx 等 Web 服务器进行托管
```
