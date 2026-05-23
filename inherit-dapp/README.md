# Legacy DApp - 去中心化遗产管理

基于 Cosmos SDK 的去中心化遗产管理 DApp，支持心跳监控、Shamir 秘密共享和 IPFS 加密存储。

## 🏗️ 架构

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Frontend   │────▶│   Backend    │────▶│  Chain Node  │
│  React/Vite  │     │  Gin + PG    │     │  Cosmos SDK  │
│  cosmjs      │     │  Heartbeat   │     │  Legacy 模块  │
└─────────────┘     └──────────────┘     └─────────────┘
       │                    │                     │
       │                    ▼                     │
       │            ┌──────────────┐              │
       │            │  PostgreSQL  │              │
       │            └──────────────┘              │
       │                                          │
       └──────────▶ ┌──────────────┐ ◀────────────┘
                    │     IPFS     │
                    │  加密存储     │
                    └──────────────┘
```

## 🚀 快速开始

### 使用 Docker Compose（推荐）

```bash
# 克隆并进入项目目录
cd inherit-dapp

# 复制环境变量
cp .env.example .env

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

服务地址：
- **前端**: http://localhost:3000
- **后端 API**: http://localhost:8080
- **链节点 RPC**: http://localhost:26657
- **IPFS Gateway**: http://localhost:8080

### 本地开发

#### 前置要求

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- IPFS (Kubo)

#### 链节点

```bash
cd chain
go mod tidy
go build -o bin/legacyd ./cmd/legacyd/
./bin/legacyd init legacy-node
./bin/legacyd serve
```

#### 后端服务

```bash
cd backend
go mod tidy
go run main.go
```

#### 前端

```bash
cd frontend
npm install
npm run dev
```

## 📁 项目结构

```
inherit-dapp/
├── chain/                      # Cosmos SDK 区块链
│   ├── cmd/legacyd/           # 链守护进程入口
│   ├── x/legacy/              # 自定义遗产模块
│   │   ├── types/             # 数据类型定义
│   │   ├── keeper/            # 业务逻辑 (Keeper)
│   │   ├── genesis/           # 创世状态
│   │   └── module/            # 模块注册
│   ├── crypto/                # 加密工具
│   │   ├── aes.go             # AES-256-GCM 加密
│   │   ├── shamir.go          # Shamir 秘密共享
│   │   └── ipfs.go            # IPFS 客户端
│   └── Dockerfile
├── backend/                    # 链下后端服务
│   ├── api/                   # Gin REST API
│   ├── monitor/               # 心跳监控器
│   ├── db/                    # PostgreSQL 数据库
│   ├── config/                # 配置管理
│   ├── main.go                # 入口
│   └── Dockerfile
├── frontend/                   # React 前端
│   ├── src/
│   │   ├── components/        # UI 组件
│   │   ├── services/          # API 和区块链服务
│   │   ├── hooks/             # React Hooks
│   │   └── types/             # TypeScript 类型
│   └── Dockerfile
├── docker-compose.yml
├── Makefile
└── .env.example
```

## 🔑 核心功能

### 1. 遗产计划管理
- 创建遗产计划，设置受益人和份额
- 支持多个受益人，份额自定义
- 添加多种资产类型（原生代币、CW20、NFT、IPFS 文档）

### 2. 心跳监控
- 计划创建者定期发送心跳
- 心跳过期自动触发遗产分配
- 后端 Keeper Bot 自动监控和提醒

### 3. 加密安全
- AES-256-GCM 对称加密保护资产数据
- Shamir 秘密共享将密钥分发给受益人
- IPFS 存储加密文档

### 4. Keplr 钱包集成
- 一键连接 Keplr 钱包
- 链上交易签名和广播
- 实时余额查询

## 🔧 API 端点

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/health` | 健康检查 |
| GET | `/api/v1/status` | 系统状态 |
| GET | `/api/v1/plans` | 列出计划 |
| GET | `/api/v1/plans/:id` | 查询计划详情 |
| GET | `/api/v1/plans/:id/assets` | 查询计划资产 |
| POST | `/api/v1/plans/:id/heartbeat` | 发送心跳 |
| POST | `/api/v1/plans/:id/claim` | 认领遗产 |
| GET | `/api/v1/creators/:address/plans` | 查询创建者的计划 |

## 🧪 测试

```bash
# 运行所有测试
make test

# 运行加密工具测试
cd chain && go test ./crypto/... -v

# 运行后端测试
cd backend && go test ./... -v
```

## 📄 许可证

MIT License
