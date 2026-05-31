# 个人博客系统

一个简单的个人博客系统，包含前端和后端。

## 技术栈

### 前端
- Vue 3
- Vue Router
- Pinia
- Element Plus
- Axios
- Vite
- Vitest

### 后端
- Go
- Gin
- GORM
- SQLite
- JWT

## 项目结构

```
blog/
├── server/          # 后端代码
│   ├── main.go
│   ├── go.mod
│   ├── models/
│   ├── controllers/
│   ├── middleware/
│   ├── database/
│   ├── routes/
│   └── utils/
└── web/            # 前端代码
    ├── package.json
    ├── vite.config.js
    └── src/
        ├── main.js
        ├── App.vue
        ├── router/
        ├── store/
        ├── views/
        └── utils/
```

## 快速开始

### 后端

1. 进入后端目录：
```bash
cd server
```

2. 安装依赖：
```bash
go mod download
```

3. 运行后端服务：
```bash
go run main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 前端

1. 进入前端目录：
```bash
cd web
```

2. 安装依赖：
```bash
npm install
```

3. 运行开发服务器：
```bash
npm run dev
```

前端服务将在 `http://localhost:5173` 启动。

## 默认账号

- 用户名：`admin`
- 密码：`123456`

## API 接口

### 认证接口

- `POST /api/v1/login` - 登录

### 前台文章接口

- `GET /api/v1/articles` - 获取已发布的文章列表
- `GET /api/v1/articles/:id` - 获取文章详情

### 后台管理接口（需要认证）

- `GET /api/v1/admin/articles` - 获取所有文章
- `POST /api/v1/admin/articles` - 创建文章
- `PUT /api/v1/admin/articles/:id` - 更新文章
- `DELETE /api/v1/admin/articles/:id` - 删除文章
- `PUT /api/v1/admin/password` - 修改密码

## 测试

### 后端测试

```bash
cd server
go test -v
```

### 前端测试

```bash
cd web
npm run test
```
