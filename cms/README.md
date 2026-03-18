# CMS 内容管理系统

基于 Golang Gin + React + Ant Design + MySQL 的内容管理系统。

## 技术栈

**后端:**
- Golang 1.21+
- Gin (Web 框架)
- GORM (ORM)
- MySQL 8.0 (数据库)
- JWT (认证)
- bcrypt (密码加密)

**前端:**
- React 18
- React Router
- Axios
- Vite
- Ant Design 5 (UI 组件库)

## 功能特性

- ✅ 管理员登录/登出
- ✅ 修改密码
- ✅ 个人资料管理
- ✅ 文章管理 (CRUD)
- ✅ 分类管理 (CRUD)
- ✅ 角色权限 (超级管理员/管理员/编辑)
- ✅ JWT 认证
- ✅ 密码加密存储

## 快速开始

### 1. 准备 MySQL 数据库

```bash
# 登录 MySQL
mysql -u root -p

# 创建数据库
CREATE DATABASE cms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 启动后端

```bash
cd backend
go mod tidy
go run main.go
```

后端将在 `http://localhost:8080` 启动

### 3. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端将在 `http://localhost:3000` 启动

## 默认账户

- 用户名：`admin`
- 密码：`admin123`

**⚠️ 首次登录后请立即修改密码！**

## API 文档

### 认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/login` | 登录 |
| POST | `/api/admin/logout` | 登出 |
| GET | `/api/admin/me` | 获取当前用户 |
| PUT | `/api/admin/change-password` | 修改密码 |
| PUT | `/api/admin/me` | 更新资料 |

### 文章接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/articles` | 文章列表 |
| GET | `/api/articles/:id` | 文章详情 |
| POST | `/api/articles` | 创建文章 |
| PUT | `/api/articles/:id` | 更新文章 |
| DELETE | `/api/articles/:id` | 删除文章 |

### 分类接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/categories` | 分类列表 |
| POST | `/api/categories` | 创建分类 |
| PUT | `/api/categories/:id` | 更新分类 |
| DELETE | `/api/categories/:id` | 删除分类 |

### 管理员接口（仅超级管理员）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admins` | 管理员列表 |
| POST | `/api/admins` | 创建管理员 |
| PUT | `/api/admins/:id` | 更新管理员 |
| DELETE | `/api/admins/:id` | 删除管理员 |

## 项目结构

```
cms/
├── backend/
│   ├── main.go              # 入口文件
│   ├── models/              # 数据模型
│   │   ├── admin.go         # 管理员模型
│   │   ├── article.go       # 文章模型
│   │   └── category.go      # 分类模型
│   ├── middleware/          # 中间件
│   │   └── auth.go          # JWT 认证
│   └── routes/              # 路由处理
│       └── routes.go
└── frontend/
    └── src/
        ├── App.jsx          # 主应用
        ├── main.jsx         # 入口
        ├── api/axios.js     # API 客户端
        └── pages/           # 页面组件
            ├── Login.jsx
            ├── Dashboard.jsx
            ├── Articles.jsx
            ├── ArticleEdit.jsx
            └── Categories.jsx
```

## 安全说明

1. **JWT Secret**: 生产环境请修改 `middleware/auth.go` 中的 `jwtKey`
2. **默认密码**: 首次登录后立即修改默认管理员密码
3. **CORS**: 生产环境应限制允许的源
4. **HTTPS**: 生产环境务必使用 HTTPS

## 开发计划

- [ ] 富文本编辑器
- [ ] 图片上传
- [ ] 标签系统
- [ ] 评论管理
- [ ] SEO 优化
- [ ] 数据备份
