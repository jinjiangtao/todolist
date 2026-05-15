# Todo List 代办清单项目

一个基于 Go + Vue 的现代化代办清单应用。

## 项目结构

```
todolist/
├── server/          # 后端服务 (Go + Gin + SQLite)
│   ├── config/      # 配置管理
│   ├── database/    # 数据库连接
│   ├── handlers/    # API 处理器
│   ├── models/      # 数据模型
│   └── routes/      # 路由配置
└── web/             # 前端应用 (Vue 3 + Vite)
    └── src/         # 源代码
```

## 功能特性

### 用户管理
- ✅ 用户注册（创建用户）
- ✅ 用户列表查询
- ✅ 用户详情查询
- ✅ 用户信息更新
- ✅ 用户删除

### 代办事项管理
- ✅ 创建代办事项
- ✅ 按用户查询代办列表（支持分页）
- ✅ 按状态筛选代办事项
- ✅ 查询单个代办详情
- ✅ 更新代办事项（标题、描述、状态、截止日期）
- ✅ 删除代办事项

## 技术栈

| 组件 | 技术 | 版本 |
|------|------|------|
| 后端框架 | Gin | 1.12.0 |
| ORM | GORM | 1.31.1 |
| 数据库 | SQLite | 嵌入式 |
| 前端框架 | Vue | 3.5.x |
| 构建工具 | Vite | 6.0.x |
| HTTP 客户端 | Axios | 1.7.x |

## 快速开始

### 后端服务

```bash
cd server
go run main.go
```

服务将在 `http://localhost:8080` 启动。

### 前端应用

```bash
cd web
npm install
npm run dev
```

前端将在 `http://localhost:5173` 启动。

## API 接口

### 用户接口

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/users | 创建用户 |
| GET | /api/v1/users | 获取用户列表 |
| GET | /api/v1/users/:id | 获取用户详情 |
| PUT | /api/v1/users/:id | 更新用户 |
| DELETE | /api/v1/users/:id | 删除用户 |

### 代办事项接口

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/todos | 创建代办事项 |
| GET | /api/v1/todos?uid=xxx | 获取代办列表（需指定用户ID） |
| GET | /api/v1/todos/:id | 获取代办详情 |
| PUT | /api/v1/todos/:id | 更新代办事项 |
| DELETE | /api/v1/todos/:id | 删除代办事项 |

## 配置说明

后端服务使用 `config.yaml` 配置：

```yaml
db_name: todo.db        # SQLite 数据库文件名
server_port: "8080"     # 服务端口
```

## Docker 运行

```bash
# 构建后端镜像
cd server
docker build -t todo-server .
docker run -p 8080:8080 todo-server

# 构建前端镜像
cd web
docker build -t todo-web .
docker run -p 5173:5173 todo-web
```

## 许可证

MIT License