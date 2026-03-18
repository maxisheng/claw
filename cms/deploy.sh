#!/bin/bash

# CMS 一键部署脚本
# 服务器：47.85.61.231
# 后端端口：8082

set -e

echo "🚀 CMS 部署脚本 - 开始部署"

# 1. 检查 Go
echo "📦 检查 Go 环境..."
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装，请先安装 Go 1.21+"
    exit 1
fi
go version

# 2. 检查 Node.js
echo "📦 检查 Node.js 环境..."
if ! command -v node &> /dev/null; then
    echo "❌ Node.js 未安装，请先安装 Node.js 16+"
    exit 1
fi
node -v

# 3. 检查 MySQL
echo "📦 检查 MySQL 环境..."
if ! command -v mysql &> /dev/null; then
    echo "❌ MySQL 未安装，请先安装 MySQL 8.0"
    exit 1
fi

# 4. 创建数据库
echo "🗄️  创建数据库..."
mysql -u root -p123456 -e "CREATE DATABASE IF NOT EXISTS cms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null || {
    echo "⚠️  数据库可能已存在或密码不正确，请手动检查"
}

# 5. 构建后端
echo "🔨 构建后端..."
cd backend
go mod tidy
CGO_ENABLED=0 GOOS=linux go build -o cms-backend main.go
echo "✅ 后端构建完成"

# 6. 构建前端
echo "🔨 构建前端..."
cd ../frontend
npm install --registry https://registry.npmmirror.com
npm run build
echo "✅ 前端构建完成"

# 7. 创建部署目录
echo "📁 创建部署目录..."
cd ..
mkdir -p /opt/cms
cp backend/cms-backend /opt/cms/
cp -r frontend/dist/* /opt/cms/ 2>/dev/null || true

# 8. 创建 systemd 服务
echo "⚙️  创建 systemd 服务..."
cat > /etc/systemd/system/cms.service << 'EOF'
[Unit]
Description=CMS Backend Service
After=network.target mysql.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/cms
ExecStart=/opt/cms/cms-backend
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=cms

[Install]
WantedBy=multi-user.target
EOF

# 9. 启动服务
echo "🚀 启动服务..."
systemctl daemon-reload
systemctl enable cms
systemctl start cms

# 10. 检查状态
echo "📊 检查服务状态..."
systemctl status cms --no-pager

echo ""
echo "=========================================="
echo "✅ CMS 部署完成！"
echo "=========================================="
echo "📍 后端地址：http://47.85.61.231:8082"
echo "📍 前端地址：http://47.85.61.231:3000 (开发环境)"
echo "📍 生产访问：http://47.85.61.231:8082 (需要配置 Nginx)"
echo ""
echo "🔐 默认账号："
echo "   用户名：admin"
echo "   密码：admin123"
echo ""
echo "📋 常用命令："
echo "   查看日志：journalctl -u cms -f"
echo "   重启服务：systemctl restart cms"
echo "   停止服务：systemctl stop cms"
echo "=========================================="
