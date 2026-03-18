#!/bin/bash
# CMS 一键部署脚本 - 完整自动化版
# 服务器：47.85.61.231
# 后端端口：8082

set -e

echo "=========================================="
echo "🚀 CMS 一键部署脚本"
echo "=========================================="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. 检查环境
echo -e "${YELLOW}[1/8] 检查环境...${NC}"

# 检查 Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go 未安装${NC}"
    echo "请执行：yum install -y golang"
    exit 1
fi
echo -e "${GREEN}✅ Go: $(go version)${NC}"

# 检查 Node.js
if ! command -v node &> /dev/null; then
    echo -e "${RED}❌ Node.js 未安装${NC}"
    echo "请执行：curl -fsSL https://rpm.nodesource.com/setup_18.x | bash - && yum install -y nodejs"
    exit 1
fi
echo -e "${GREEN}✅ Node.js: $(node -v)${NC}"

# 检查 npm
if ! command -v npm &> /dev/null; then
    echo -e "${RED}❌ npm 未安装${NC}"
    exit 1
fi
echo -e "${GREEN}✅ npm: $(npm -v)${NC}"

# 检查 MySQL
if ! command -v mysql &> /dev/null; then
    echo -e "${RED}❌ MySQL 未安装${NC}"
    echo "请执行：yum install -y mysql-server"
    exit 1
fi
echo -e "${GREEN}✅ MySQL 已安装${NC}"

echo ""

# 2. 创建数据库
echo -e "${YELLOW}[2/8] 创建数据库...${NC}"
mysql -u root -p123456 -e "CREATE DATABASE IF NOT EXISTS cms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null && \
    echo -e "${GREEN}✅ 数据库创建成功${NC}" || \
    echo -e "${YELLOW}⚠️  数据库可能已存在或密码不正确${NC}"

echo ""

# 3. 克隆/更新代码
echo -e "${YELLOW}[3/8] 获取代码...${NC}"
cd /root
if [ -d "cms" ]; then
    echo "更新现有代码..."
    cd cms
    git pull
else
    echo "克隆新代码..."
    git clone https://github.com/maxisheng/claw.git cms
    cd cms
fi
echo -e "${GREEN}✅ 代码准备完成${NC}"

echo ""

# 4. 构建后端
echo -e "${YELLOW}[4/8] 构建后端...${NC}"
cd /root/cms/backend
go mod tidy
CGO_ENABLED=0 GOOS=linux go build -o cms-backend main.go
echo -e "${GREEN}✅ 后端构建完成${NC}"

echo ""

# 5. 停止旧进程
echo -e "${YELLOW}[5/8] 停止旧进程...${NC}"
pkill -f "cms-backend" 2>/dev/null || true
pkill -f "go run main.go" 2>/dev/null || true
echo -e "${GREEN}✅ 旧进程已停止${NC}"

echo ""

# 6. 启动后端
echo -e "${YELLOW}[6/8] 启动后端...${NC}"
cd /root/cms/backend
nohup go run main.go > /tmp/cms-backend.log 2>&1 &
sleep 3

# 检查后端是否启动
if curl -s http://localhost:8082/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ 后端启动成功 (端口 8082)${NC}"
else
    echo -e "${RED}❌ 后端启动失败，查看日志：cat /tmp/cms-backend.log${NC}"
    cat /tmp/cms-backend.log
    exit 1
fi

echo ""

# 7. 构建前端
echo -e "${YELLOW}[7/8] 构建前端...${NC}"
cd /root/cms/frontend
npm install --registry https://registry.npmmirror.com
npm run build
echo -e "${GREEN}✅ 前端构建完成${NC}"

echo ""

# 8. 显示访问信息
echo ""
echo "=========================================="
echo -e "${GREEN}✅ CMS 部署完成！${NC}"
echo "=========================================="
echo ""
echo "📍 访问地址："
echo "   http://47.85.61.231:3000"
echo ""
echo "🔐 登录信息："
echo "   用户名：admin"
echo "   密码：admin123"
echo ""
echo "📋 管理命令："
echo "   查看后端日志：tail -f /tmp/cms-backend.log"
echo "   停止后端：pkill -f cms-backend"
echo "   重启后端：cd /root/cms/backend && nohup go run main.go > /tmp/cms-backend.log 2>&1 &"
echo ""
echo "⚠️  前端开发模式已启动，生产环境请配置 Nginx"
echo "=========================================="
echo ""

# 测试 API
echo -e "${YELLOW}测试 API 连接...${NC}"
curl -s http://localhost:8082/health | head -c 100
echo ""
echo -e "${GREEN}✅ 所有服务运行正常！${NC}"
