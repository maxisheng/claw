# CMS - Content Management System

一个基于 Golang + React 的简单内容管理系统。

## 技术栈

**后端:**
- Golang 1.21+
- Gin (Web 框架)
- GORM (ORM)
- MySQL (数据库)
- JWT (认证)

**前端:**
- React 18
- React Router
- Axios
- Vite

## 项目结构

```
cms/
├── backend/           # Go 后端
│   ├── main.go       # 入口文件
│   ├── models/       # 数据模型
│   ├── routes/       # API 路由
│   └── middleware/   # 中间件
└── frontend/         # React 前端
    └── src/
        ├── pages/    # 页面组件
        ├── components/
        └── api/      # API 客户端
```

## 快速开始

### 0. 准备 MySQL 数据库

```sql
CREATE DATABASE cms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

或者复制环境变量文件并修改配置：
```bash
cd backend
cp .env.example .env
# 编辑 .env 文件，修改 MySQL 连接信息
```

### 1. 启动后端

```bash
cd backend
go mod tidy
go run main.go
```

后端将在 `http://localhost:8080` 启动

**环境变量配置:**
- `DB_HOST` - MySQL 主机 (默认：localhost)
- `DB_PORT` - MySQL 端口 (默认：3306)
- `DB_USER` - 数据库用户 (默认：root)
- `DB_PASSWORD` - 数据库密码 (默认：root)
- `DB_NAME` - 数据库名 (默认：cms)
- 或直接设置 `MYSQL_DSN` 指定完整连接字符串

### 2. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端将在 `http://localhost:3000` 启动

## 默认账户

- 用户名：`admin`
- 密码：`admin123`

## API 端点

### 公开接口
- `POST /api/login` - 登录
- `GET /api/articles` - 获取文章列表
- `GET /api/articles/:slug` - 获取单篇文章
- `GET /api/categories` - 获取分类列表

### 需要认证的接口
- `GET /api/user` - 获取当前用户
- `POST /api/articles` - 创建文章
- `PUT /api/articles/:id` - 更新文章
- `DELETE /api/articles/:id` - 删除文章
- `POST /api/categories` - 创建分类
- `PUT /api/categories/:id` - 更新分类
- `DELETE /api/categories/:id` - 删除分类

## 功能特性

- ✅ 用户认证 (JWT)
- ✅ 文章管理 (CRUD)
- ✅ 分类管理
- ✅ 文章状态 (草稿/发布/归档)
- ✅ 响应式管理界面

## 待扩展

- [ ] 富文本编辑器
- [ ] 图片上传
- [ ] 用户管理
- [ ] 角色权限
- [ ] SEO 优化
- [ ] 多数据库支持 (PostgreSQL/MySQL)
