# TypeScript API 客户端使用指南

## 快速开始

### 1. 初始化拦截器

在应用入口文件（如 `main.ts`）中初始化拦截器：

```typescript
import { setupInterceptors } from '@/api';

// 设置拦截器
setupInterceptors();
```

### 2. 调用 API

```typescript
import { api } from '@/api';

// 示例：获取用户列表
const response = await api.apiV1AdminUsersGet();
console.log(response.data);

// 示例：创建用户
const newUser = await api.apiV1AdminUsersPost({
  username: 'test',
  email: 'test@example.com',
  password: '123456',
});
```

### 3. 使用类型

```typescript
import type { PiemdmInternalModelUser } from '@/api';

const user: PiemdmInternalModelUser = {
  ID: 1,
  Username: 'admin',
  Email: 'admin@example.com',
};
```

## API 实例

所有 API 方法都通过 `api` 实例调用：

```typescript
import { api } from '@/api';

// GET 请求
await api.apiV1AdminUsersGet();

// POST 请求
await api.apiV1AdminUsersPost({ ...data });

// PUT 请求
await api.apiV1AdminUsersIdPut('user-id', { ...data });

// DELETE 请求
await api.apiV1AdminUsersIdDelete('user-id');
```

## 配置

### 环境变量

在 `.env` 文件中配置 API 基础 URL：

```env
VITE_API_BASE_URL=http://localhost:8787
```

### 动态更新 Token

```typescript
import { updateApiConfig } from '@/api';

// 登录后更新 token
updateApiConfig(userToken);
```

## 错误处理

拦截器会自动处理：
- **401 未授权**：自动登出并跳转到登录页
- **403 禁止访问**：控制台输出错误信息

自定义错误处理：

```typescript
try {
  const response = await api.apiV1AdminUsersGet();
  console.log(response.data);
} catch (error) {
  if (error.response?.status === 404) {
    console.error('User not found');
  }
}
```

## 重新生成客户端

当后端 API 更新后，重新生成客户端：

```bash
# 1. 生成 Swagger 文档
cd backend
make swagger

# 2. 生成 TypeScript 客户端
cd ../frontend
pnpm gen:api

# 或使用根目录命令
cd ..
make gen-api
```

## 注意事项

1. **类型安全**：所有 API 调用都有完整的类型提示
2. **自动完成**：IDE 会提供智能提示
3. **编译检查**：类型错误会在编译时发现
4. **文档同步**：客户端代码与后端 API 自动同步

## 示例

### 完整的用户管理示例

```typescript
import { api } from '@/api';
import type { PiemdmInternalModelUser } from '@/api';

// 获取用户列表
const getUserList = async () => {
  try {
    const response = await api.apiV1AdminUsersGet(1, 15);
    return response.data;
  } catch (error) {
    console.error('Failed to get users:', error);
    throw error;
  }
};

// 创建用户
const createUser = async (userData: any) => {
  try {
    const response = await api.apiV1AdminUsersPost(userData);
    return response.data;
  } catch (error) {
    console.error('Failed to create user:', error);
    throw error;
  }
};

// 更新用户
const updateUser = async (id: string, userData: any) => {
  try {
    const response = await api.apiV1AdminUsersIdPut(id, userData);
    return response.data;
  } catch (error) {
    console.error('Failed to update user:', error);
    throw error;
  }
};

// 删除用户
const deleteUser = async (id: string) => {
  try {
    await api.apiV1AdminUsersIdDelete(id);
  } catch (error) {
    console.error('Failed to delete user:', error);
    throw error;
  }
};
```
