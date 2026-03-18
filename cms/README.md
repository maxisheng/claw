# CMS 内容管理系统

基于 Golang Gin + React + Ant Design + MySQL 的内容管理系统，采用标准分层架构。

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

## 架构设计

参考流行的 Gin 项目架构（如 gin-vue-admin、go-admin），采用标准分层：

```
backend/
├── main.go              # 应用入口、依赖初始化
├── models/              # 数据模型层
│   ├── admin.go         # 管理员模型
│   ├── article.go       # 文章模型
│   └── category.go      # 分类模型
├── repository/          # 数据访问层 (DAO)
│   ├── admin_repository.go
│   ├── article_repository.go
│   └── category_repository.go
├── service/             # 业务逻辑层
│   ├── admin_service.go
│   ├── article_service.go
│   └── category_service.go
├── handlers/            # HTTP 处理层 (Controller)
│   ├── admin_handler.go
│   ├── article_handler.go
│   └── category_handler.go
├── middleware/          # 中间件层
│   └── auth.go          # JWT 认证、权限验证
└── routes/              # 路由注册层
    └── routes.go
```

### 各层职责

| 层级 | 职责 | 依赖 |
|------|------|------|
| **handlers** | HTTP 请求处理、参数验证、响应格式化 | service |
| **service** | 业务逻辑、事务处理、数据校验 | repository |
| **repository** | 数据库 CRUD 操作 | models |
| **models** | 数据模型定义 | - |
| **middleware** | 认证、授权、日志等横切关注点 | - |
| **routes** | 路由注册、依赖注入 | handlers, middleware |

## 功能特性

- ✅ 管理员登录/登出 (JWT)
- ✅ 修改密码 (bcrypt 加密)
- ✅ 个人资料管理
- ✅ 文章管理 (CRUD)
- ✅ 分类管理 (CRUD)
- ✅ 角色权限 (超级管理员/管理员/编辑)
- ✅ 分层架构设计

## 快速开始

### 1. 准备 MySQL 数据库

```bash
mysql -u root -p123456 -e "CREATE DATABASE cms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### 2. 启动后端

```bash
cd backend
go mod tidy
go run main.go
```

### 3. 启动前端

```bash
cd frontend
npm install
npm run dev
```

## 默认账户

- 用户名：`admin`
- 密码：`admin123`

## API 文档

### 认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/login` | 登录 |
| POST | `/api/admin/logout` | 登出 |
| GET | `/api/admin/me` | 获取当前用户 |
| PUT | `/api/admin/change-password` | 修改密码 |

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

## 扩展开发

### 添加新模块

1. 在 `models/` 定义数据模型
2. 在 `repository/` 实现数据访问方法
3. 在 `service/` 实现业务逻辑
4. 在 `handlers/` 实现 HTTP 处理
5. 在 `routes/` 注册路由

### 示例：添加标签模块

```go
// 1. models/tag.go
type Tag struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"size:50;uniqueIndex"`
}

// 2. repository/tag_repository.go
type TagRepository struct { db *gorm.DB }
func (r *TagRepository) FindAll() ([]Tag, error) { ... }

// 3. service/tag_service.go
type TagService struct { repo *TagRepository }
func (s *TagService) GetAll() ([]Tag, error) { ... }

// 4. handlers/tag_handler.go
type TagHandler struct { service *TagService }
func (h *TagHandler) ListTags(c *gin.Context) { ... }

// 5. routes/routes.go
api.GET("/tags", tagHandler.ListTags)
```

## 项目优势

1. **清晰的分层架构** - 职责分离，易于维护和测试
2. **依赖注入** - 在 routes 层统一初始化依赖
3. **统一响应格式** - handlers 层统一处理响应
4. **业务逻辑复用** - service 层可被多个 handler 调用
5. **易于扩展** - 添加新模块只需遵循相同模式

## 安全说明

1. **JWT Secret**: 生产环境请修改 `middleware/auth.go` 中的 `jwtKey`
2. **默认密码**: 首次登录后立即修改
3. **HTTPS**: 生产环境务必使用 HTTPS
4. **CORS**: 生产环境应限制允许的源
