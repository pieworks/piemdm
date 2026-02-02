# workflow-vue

A structured workflow builder component for Vue 3.

## 安装

```sh
pnpm add workflow-vue
# 或
npm install workflow-vue
```

## 使用方式

### 1. 作为Vue插件使用

```javascript
// main.js
import { createApp } from 'vue';
import App from './App.vue';
import { WorkflowVue } from 'workflow-vue';
import 'workflow-vue/style.css';

const app = createApp(App);
app.use(WorkflowVue);
app.mount('#app');
```

### 2. 按需导入组件

```javascript
// 在Vue组件中
import { WorkflowBuilder } from 'workflow-vue';

export default {
  components: {
    WorkflowBuilder
  }
};
```

### 3. 使用服务和工具

```javascript
import { WorkflowService, NodeService } from 'workflow-vue';

// 创建工作流
const workflow = WorkflowService.createWorkflow('工作流名称');

// 创建节点
const node = NodeService.createNode('APPROVAL', '审批节点');
```

### 4. 使用常量

```javascript
import { NODE_TYPES } from 'workflow-vue';

console.log(NODE_TYPES.START); // 开始节点配置
```

## 开发

```sh
# 安装依赖
pnpm install

# 开发模式
pnpm dev

# 构建
pnpm build

# 测试
pnpm test
```

## 本地开发说明

### Monorepo 内部引用

在 piemdm monorepo 中，`frontend` 通过 `"workflow-vue": "workspace:*"` 引用本地包时：

- ✅ **无需任何修改**：本地开发完全不受影响
- ✅ **Vite 直接处理源码**：在 pnpm workspace 模式下，Vite 会直接处理 `src/lib/index.js` 源码，不会使用 `package.json` 中的 `exports` 配置
- ✅ **热更新正常**：修改源码后，前端应用会自动热更新
- ✅ **无需预先构建**：本地开发时不需要先执行 `pnpm build`

### 为什么指向构建产物不影响本地开发？

`package.json` 中的 `main`、`module`、`types`、`exports` 字段主要用于：
- **发布到 npm 后**：外部项目通过 `npm install workflow-vue` 安装时使用
- **类型定义**：TypeScript 项目可能需要类型定义文件

在 monorepo 内部，pnpm workspace 会直接链接到源码目录，构建工具（如 Vite）会直接处理源码文件，因此这些字段的配置不会影响本地开发体验。

## 发布

`workflow-vue` 作为 monorepo 的子包，可以直接在 `packages/web/workflow-vue` 目录下发布到 npm 或私有 registry，无需移出 monorepo。

### 发布前检查

确保 `package.json` 配置正确：

- ✅ `name`: 包名（当前为 `workflow-vue`，如需 scoped 包可改为 `@pieteams/workflow-vue`）
- ✅ `private`: 不能为 `true`（否则无法发布）
- ✅ `version`: 按语义化版本递增（例如 `1.0.1`）
- ✅ `main`/`module`/`types`/`exports`: 已正确配置构建产物路径
- ✅ `files`: 已指定要发布的文件（`dist`、`README.md`、`LICENSE` 等）
- ✅ `peerDependencies`: `vue`、`bootstrap`、`bootstrap-icons` 已配置为 peer 依赖
- ✅ `dependencies`: 仅包含库自身需要的运行时依赖

### 发布到 npm（公开包）

在 monorepo 根目录执行：

```sh
# 方式一：使用 -C 指定目录
pnpm -C packages/web/workflow-vue install
pnpm -C packages/web/workflow-vue build
pnpm -C packages/web/workflow-vue test
pnpm -C packages/web/workflow-vue publish --access public

# 方式二：使用 filter（推荐，更符合 monorepo 操作习惯）
pnpm -r --filter workflow-vue install
pnpm -r --filter workflow-vue build
pnpm -r --filter workflow-vue test
pnpm -r --filter workflow-vue publish --access public
```

> **注意**：`package.json` 中已配置 `prepublishOnly` 脚本，发布前会自动执行构建。

### 版本管理

发布前记得更新版本号：

```sh
# 在 packages/web/workflow-vue 目录下
pnpm version patch  # 1.0.0 -> 1.0.1
pnpm version minor  # 1.0.0 -> 1.1.0
pnpm version major  # 1.0.0 -> 2.0.0
```

或直接编辑 `package.json` 中的 `version` 字段。

### 何时考虑拆独立仓库

只有在以下情况下才考虑将包拆分为独立仓库：

- 需要完全独立的权限/可见性控制
- 需要完全独立的发布流程
- 不希望消费者项目获取到 monorepo 相关元信息

否则，在 monorepo 子目录直接发布是最省事且可维护的方案。

## 依赖要求

- Vue 3.4+
- Bootstrap 5.3+
- Bootstrap Icons 1.11+

## 许可证

MIT
