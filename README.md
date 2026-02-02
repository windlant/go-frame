# 用户管理 API 文档

基于 GoFrame 框架的用户批量操作接口，所有接口均位于 `/users` 路由组下，统一使用 POST 方法。

## 统一响应格式

### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 错误响应
```json
{
  "code": 1001,
  "message": "错误描述信息",
  "data": null
}
```

## 错误码说明

| 错误码 | 含义 | 触发场景 |
|--------|------|----------|
| 0 | OK | 操作成功 |
| 1001 | Invalid parameters | 请求参数格式错误或缺失 |
| 1002 | Resource not found | 查询的资源不存在 |
| 1003 | Internal server error | 服务器内部错误 |
| 1004 | User already exists | 用户已存在（邮箱重复） |
| 1005 | Batch size too large | 批量操作数量超出限制 |

## API 接口详情

### 1. 查询所有用户
**Endpoint**: POST /users/list  
**Description**: 获取所有用户列表

**请求示例**:
```bash
curl -X POST http://localhost:8000/users/list
```

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "Alice",
      "email": "alice@example.com"
    }
  ]
}
```

### 2. 批量创建用户
**Endpoint**: POST /users/create  
**Description**: 批量创建新用户

**请求体**:
```json
[
  {
    "name": "Alice",
    "email": "alice@example.com"
  },
  {
    "name": "Bob",
    "email": "bob@example.com"
  }
]
```

**请求示例**:
```bash
curl -X POST http://localhost:8000/users/create \
  -H "Content-Type: application/json" \
  -d '[{"name":"Alice","email":"alice@example.com"}]'
```

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "first_id": 1,
    "count": 1
  }
}
```

### 3. 批量查询用户
**Endpoint**: POST /users/get  
**Description**: 根据 ID 列表或邮箱列表批量查询用户

**请求体** (支持单独或同时查询):
```json
{
  "ids": [1, 2],
  "emails": ["alice@example.com"]
}
```

**请求示例**:
```bash
curl -X POST http://localhost:8000/users/get \
  -H "Content-Type: application/json" \
  -d '{"ids":[1],"emails":["alice@example.com"]}'
```

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "by_id": [
      {"id": 1, "name": "Alice", "email": "alice@example.com"}
    ],
    "by_email": [
      {"id": 1, "name": "Alice", "email": "alice@example.com"}
    ]
  }
}
```

### 4. 批量更新用户
**Endpoint**: POST /users/update  
**Description**: 批量更新用户信息

**请求体**:
```json
[
  {
    "id": 1,
    "name": "Alice Updated",
    "email": "alice.updated@example.com"
  }
]
```

**请求示例**:
```bash
curl -X POST http://localhost:8000/users/update \
  -H "Content-Type: application/json" \
  -d '[{"id":1,"name":"Alice Updated","email":"alice.updated@example.com"}]'
```

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "updated": 1
  }
}
```

### 5. 批量删除用户
**Endpoint**: POST /users/delete  
**Description**: 批量删除用户

**请求体**:
```json
{
  "ids": [1, 2]
}
```

**请求示例**:
```bash
curl -X POST http://localhost:8000/users/delete \
  -H "Content-Type: application/json" \
  -d '{"ids":[1]}'
```

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "deleted": 1
  }
}
```

## 数据模型

用户对象结构：
```json
{
  "id": 1,
  "name": "用户名",
  "email": "user@example.com"
}
```