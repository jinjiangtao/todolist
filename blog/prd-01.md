# 博客系统（基础版）产品需求文档

## 1. 项目概述

### 1.1 项目名称
个人博客系统 

### 1.2 项目目标
实现一个最小可用的博客系统，支持用户登录认证和文章发布功能。其他功能（分类、标签、评论等）后续迭代添加。

### 1.3 技术栈
| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + Vue Router + Pinia + Axios + Element Plus |
| 后端 | Go + Gin + GORM |
| 数据库 | SQLite |
| 认证 | JWT |

> 使用 SQLite 无需单独安装数据库服务，数据存储在本地单文件中，适合小型个人博客或本地开发。

---

## 2. 用户角色

| 角色 | 说明 |
|------|------|
| 访客 | 浏览已发布的文章 |
| 管理员 | 登录后发布、编辑、删除文章（单用户场景，无需注册）|

> 基础版只有**一个预设管理员账号**，不开放注册功能。

---

## 3. 功能需求

### 3.1 前台功能（访客可见）

#### 3.1.1 首页
- 展示已发布的文章列表（按发布时间倒序）
- 列表项包含：文章标题、发布时间、摘要
- 分页加载（每页10条）

#### 3.1.2 文章详情页
- 展示文章标题、发布时间、正文内容
- 正文支持Markdown渲染

#### 3.1.3 登录页
- 用户名 + 密码登录
- 登录成功后跳转后台管理页

### 3.2 后台功能（需登录）

#### 3.2.1 文章管理
- 文章列表：展示所有文章（标题、发布时间、状态：草稿/发布）
- 发布文章：
  - 标题（必填）
  - 正文（Markdown编辑器）
  - 状态：保存为草稿 或 直接发布
- 编辑文章：修改标题、正文、状态
- 删除文章：确认后删除

#### 3.2.2 修改密码
- 输入旧密码 + 新密码，修改管理员密码

---

## 4. 非功能需求

### 4.1 基础要求
- 响应式布局（PC端为主，手机端基本可用）
- 前端路由守卫：未登录状态下不能访问后台页面

### 4.2 安全
- 密码加密存储（bcrypt）
- JWT过期时间（默认7天）
- 后端接口统一验证Token（登录接口除外）

### 4.3 部署便利性
- SQLite 数据库文件自动创建在项目根目录 `blog.db`
- 无需额外安装数据库服务，单文件即拷即用

---

## 5. 数据库设计

### 5.1 用户表 `users`
| 字段 | 类型 | 说明 |
|------|------|------|
| id | integer | 主键，自增 |
| username | varchar(50) | 用户名（唯一）|
| password | varchar(255) | bcrypt加密密码 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 5.2 文章表 `articles`
| 字段 | 类型 | 说明 |
|------|------|------|
| id | integer | 主键，自增 |
| title | varchar(200) | 标题（必填）|
| content | text | 正文（Markdown）|
| status | integer | 0-草稿，1-发布，默认0 |
| view_count | integer | 阅读量，默认0 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

> 基础版暂不实现分类、标签、封面图、软删除等功能

### 5.3 SQLite 连接配置
```go
// GORM 连接 SQLite
import "gorm.io/driver/sqlite"

dsn := "blog.db"
db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
```

---

6. API接口设计

6.1 认证相关

方法 路径 说明 是否需要登录
POST /api/v1/login 登录，返回JWT token 否

6.2 文章相关（前台）

方法 路径 说明
GET /api/v1/articles 获取已发布文章列表（分页）
GET /api/v1/articles/:id 获取文章详情（阅读量+1）

6.3 文章相关（后台，需Token）

方法 路径 说明
GET /api/v1/admin/articles 获取所有文章（含草稿，分页）
POST /api/v1/admin/articles 创建文章
PUT /api/v1/admin/articles/:id 更新文章
DELETE /api/v1/admin/articles/:id 删除文章

6.4 用户相关（后台，需Token）

方法 路径 说明
PUT /api/v1/admin/password 修改密码

---

7. 前端页面清单

路由 页面 说明
/ 首页 文章列表
/article/:id 文章详情页 展示文章
/login 登录页 管理员登录
/admin 后台首页 文章管理列表
/admin/article/edit/:id? 编辑/发布文章 无id时为新建

---

8. 预设数据

初始化数据库时自动创建管理员账号：

· 用户名：admin
· 密码：123456（bcrypt加密存储，首次登录后建议修改）

初始化SQLite代码示例：

```go
// 自动迁移
db.AutoMigrate(&models.User{}, &models.Article{})

// 创建默认管理员
var user models.User
result := db.Where("username = ?", "admin").First(&user)
if result.Error != nil {
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
    admin := models.User{Username: "admin", Password: string(hashedPassword)}
    db.Create(&admin)
}
```

### 项目约束
前端代码放到web 文件夹中
后端代码放到 server 文件夹中
前端和后端代码都编写测试用例代码。 交付物可直接使用。

