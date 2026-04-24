# Todo List API 规范文档

## 1. 项目概述

- **项目名称**: Todo List API
- **项目类型**: RESTful Web API
- **核心功能**: 提供代办事项的增删改查（CRUD）接口
- **目标用户**: 前端应用、移动应用等客户端

## 2. 技术栈

- **编程语言**: Go 1.21+
- **Web 框架**: Gin
- **数据库**: MySQL 8.0+
- **ORM**: GORM

## 3. 数据库设计

### 表结构: `todos`

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 唯一标识 |
| title | VARCHAR(255) | NOT NULL | 标题 |
| description | TEXT | NULL | 描述 |
| status | TINYINT | NOT NULL, DEFAULT 0 | 状态：0-待办，1-进行中，2-已完成 |
| due_date | DATE | NULL | 截止日期 |
| created_at | DATETIME | NOT NULL | 创建时间 |
| updated_at | DATETIME | NOT NULL | 更新时间 |

## 4. API 接口设计

### 基础路径: `/api/v1`

### 4.1 创建待办事项
- **方法**: POST
- **路径**: `/todos`
- **请求体**:
```json
{
  "title": "string (必填)",
  "description": "string (可选)",
  "status": "int (可选, 默认0)",
  "due_date": "string (可选, 格式: 2006-01-02)"
}
```
- **响应**: 201 Created
```json
{
  "code": 201,
  "message": "创建成功",
  "data": { /* todo对象 */ }
}
```

### 4.2 获取所有待办事项
- **方法**: GET
- **路径**: `/todos`
- **查询参数**:
  - `status` (可选): 过滤状态
  - `page` (可选): 页码，默认1
  - `page_size` (可选): 每页数量，默认10
- **响应**: 200 OK
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    { /* todo对象 */ }
  ]
}
```

### 4.3 获取单个待办事项
- **方法**: GET
- **路径**: `/todos/:id`
- **响应**: 200 OK
```json
{
  "code": 200,
  "message": "获取成功",
  "data": { /* todo对象 */ }
}
```

### 4.4 更新待办事项
- **方法**: PUT
- **路径**: `/todos/:id`
- **请求体**:
```json
{
  "title": "string (可选)",
  "description": "string (可选)",
  "status": "int (可选)",
  "due_date": "string (可选)"
}
```
- **响应**: 200 OK
```json
{
  "code": 200,
  "message": "更新成功",
  "data": { /* todo对象 */ }
}
```

### 4.5 删除待办事项
- **方法**: DELETE
- **路径**: `/todos/:id`
- **响应**: 200 OK
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

## 5. 项目结构

```
d:\trae\todolist\
├── main.go              # 入口文件
├── config/
│   └── config.go        # 配置管理
├── models/
│   └── todo.go          # 数据模型
├── handlers/
│   └── todo.go          # 处理器函数
├── routes/
│   └── routes.go        # 路由配置
├── database/
│   └── mysql.go         # 数据库连接
└── go.mod               # Go模块文件
```

## 6. 错误处理

### 错误响应格式
```json
{
  "code": 400,
  "message": "错误信息",
  "data": null
}
```

### HTTP 状态码
- 200: 成功
- 201: 创建成功
- 400: 请求参数错误
- 404: 资源不存在
- 500: 服务器内部错误

## 7. 数据库初始化SQL

```sql
CREATE DATABASE IF NOT EXISTS todo_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE todo_db;

CREATE TABLE IF NOT EXISTS todos (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status TINYINT NOT NULL DEFAULT 0,
    due_date DATE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```
