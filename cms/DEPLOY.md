# CMS 生产环境部署指南

## 服务器信息

- **IP 地址**: 47.85.61.231
- **后端端口**: 8082
- **前端端口**: 3000 (开发) / 80 (生产 via Nginx)

---

## 方式一：一键部署脚本（推荐）

```bash
# 1. 上传代码到服务器
scp -r cms root@47.85.61.231:/home/admin/

# 2. SSH 登录服务器
ssh root@47.85.61.231

# 3. 执行部署脚本
cd /home/admin/cms
chmod +x auto-deploy.sh
./auto-deploy.sh
```

---

## 方式二：手动部署

### 1. 安装依赖

```bash
# 安装 Go
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 安装 Node.js
curl -fsSL https://rpm.nodesource.com/setup_18.x | bash -
yum install -y nodejs

# 安装 MySQL
yum install -y mysql-server
systemctl start mysqld
systemctl enable mysqld
```

### 2. 创建数据库

```bash
mysql -u root -p -e "CREATE DATABASE cms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### 3. 构建后端

```bash
cd backend
go mod tidy
CGO_ENABLED=0 GOOS=linux go build -o cms-backend main.go
```

### 4. 构建前端

```bash
cd frontend
npm install
npm run build
```

### 5. 部署到生产目录

```bash
mkdir -p /home/admin/cms-deploy
cp backend/cms-backend /home/admin/cms-deploy/
cp -r frontend/dist /home/admin/cms-deploy/dist
```

### 6. 配置 systemd 服务

```bash
cat > /etc/systemd/system/cms.service << 'EOF'
[Unit]
Description=CMS Backend Service
After=network.target mysql.service

[Service]
Type=simple
User=root
WorkingDirectory=/home/admin/cms-deploy
ExecStart=/home/admin/cms-deploy/cms-backend
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=cms

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable cms
systemctl start cms
```

### 7. 配置 Nginx（可选）

```bash
# 安装 Nginx
yum install -y nginx

# 复制配置文件
cp nginx.conf /etc/nginx/conf.d/cms.conf

# 启动 Nginx
systemctl start nginx
systemctl enable nginx
```

---

## 访问地址

| 环境 | 地址 | 说明 |
|------|------|------|
| 开发 | http://47.85.61.231:3000 | Vite 开发服务器 |
| 生产 | http://47.85.61.231 | Nginx 反向代理 |
| API | http://47.85.61.231:8082/api | 后端 API 直接访问 |

---

## 默认账号

- 用户名：`admin`
- 密码：`admin123`

**⚠️ 首次登录后请立即修改密码！**

---

## 常用命令

```bash
# 查看服务状态
systemctl status cms

# 查看日志
journalctl -u cms -f

# 重启服务
systemctl restart cms

# 停止服务
systemctl stop cms

# 查看后端进程
ps aux | grep cms-backend

# 查看端口占用
netstat -tlnp | grep 8082
```

---

## 防火墙配置

```bash
# 开放端口
firewall-cmd --permanent --add-port=8082/tcp
firewall-cmd --permanent --add-port=80/tcp
firewall-cmd --reload

# 或使用 iptables
iptables -A INPUT -p tcp --dport 8082 -j ACCEPT
iptables -A INPUT -p tcp --dport 80 -j ACCEPT
```

---

## 安全建议

1. **修改默认密码** - 首次登录后立即修改
2. **配置 HTTPS** - 使用 Let's Encrypt 免费证书
3. **限制 API 访问** - 配置防火墙规则
4. **定期备份** - 备份 MySQL 数据库
5. **更新依赖** - 定期更新 Go 和 Node.js 依赖

---

## 数据库备份

```bash
# 备份
mysqldump -u root -p cms > cms_backup_$(date +%Y%m%d).sql

# 恢复
mysql -u root -p cms < cms_backup_20240101.sql
```

---

## 故障排查

### 后端启动失败

```bash
# 查看日志
journalctl -u cms -f

# 检查端口占用
netstat -tlnp | grep 8082

# 检查 MySQL 连接
mysql -u root -p -e "SHOW DATABASES;"
```

### 前端无法访问

```bash
# 检查 Nginx 配置
nginx -t

# 查看 Nginx 日志
tail -f /var/log/nginx/error.log

# 检查文件权限
ls -la /opt/cms/dist
```
