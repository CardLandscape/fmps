# FMPS - 家庭惩戒管理系统

> Family Punishment Management System

本地优先的家庭管理应用，采用 Go (Gin) + Vue 3 + SQLite + Electron 技术栈。

## 技术架构

| 层级 | 技术 |
|------|------|
| 后端 API | Go + Gin + GORM |
| 数据库 | SQLite（极小数据量） |
| 前端 | Vue 3 + Vite + Element Plus |
| 桌面版 | Electron（本地模式，内嵌 Go API） |
| 部署 | Ubuntu ECS + systemd |

## 目录结构

```
server/                 # Go Gin API
  main.go               # 入口
  router.go             # 路由配置
  config.go             # 配置（环境变量）
  db/                   # 数据库初始化
  models/               # 数据模型
  handlers/             # API 处理函数
  middleware/           # 中间件（JWT 鉴权）
  seeds/                # 初始数据
web/                    # Vue 3 前端
  src/
    views/              # 页面（登录、仪表盘、成员、规则、记录、设置）
    utils/api.js        # API 封装
    router/             # 路由
electron/               # Electron 桌面壳
  main.js               # 主进程（启动 Go API）
  preload.js
scripts/
  start-local.sh        # 本地启动脚本
  build-electron.sh     # 打包脚本
deploy/
  fmps.service          # systemd 服务单元
  README.md             # ECS 部署指南
```

## 快速开始（本地开发）

### 启动后端

```bash
cd server
go run .
# 服务启动于 http://localhost:8080
```

### 启动前端开发服务器

```bash
cd web
npm install
npm run dev
# 前端启动于 http://localhost:5173
```

### 一键本地启动

```bash
./scripts/start-local.sh
```

## 鉴权

- 首次启动默认账号：`admin` / `123456`
- 登录接口：`POST /api/login` → 返回 JWT token
- 其他 API 需携带 `Authorization: Bearer <token>`
- 可在设置页面修改用户名/密码

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/login | 登录，获取 token |
| GET | /api/members | 家庭成员列表 |
| POST | /api/members | 添加成员 |
| PUT | /api/members/:id | 修改成员 |
| DELETE | /api/members/:id | 删除成员 |
| GET | /api/rules | 惩戒规则列表 |
| POST | /api/rules | 添加规则 |
| PUT | /api/rules/:id | 修改规则 |
| DELETE | /api/rules/:id | 删除规则 |
| GET | /api/records | 惩戒记录列表 |
| POST | /api/records | 添加记录 |
| DELETE | /api/records/:id | 删除记录 |
| GET | /api/stats | 统计数据 |
| GET | /api/settings | 系统设置 |
| PUT | /api/settings | 修改设置 |

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| FMPS_PORT | 8080 | API 监听端口 |
| FMPS_DB_PATH | ~/.fmps/fmps.db | SQLite 数据库路径 |
| FMPS_JWT_SECRET | fmps-secret-key-2024 | JWT 签名密钥 |

## 桌面版（Electron）

```bash
# 构建完整桌面包
./scripts/build-electron.sh

# 或直接运行（需先构建前端和后端）
cd electron
npm install
npm start
```

## ECS 部署

参见 [deploy/README.md](deploy/README.md)

## 数据库模型

- **settings** - 系统配置（admin_username, admin_password 等）
- **members** - 家庭成员（姓名、角色）
- **rules** - 惩戒规则（名称、描述、分类、分值）
- **records** - 惩戒记录（成员、规则、分值、备注、时间）
