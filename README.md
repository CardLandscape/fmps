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
    migrations/         # SQL 迁移参考文件
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
| PUT | /api/members/:id | 修改成员（需授权密码） |
| POST | /api/members/:id/delete | 删除成员（需授权密码） |
| DELETE | /api/members/:id | 删除成员（向下兼容，无密码校验） |
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
| GET | /api/cases | 案件列表 |
| POST | /api/cases | 新建案件 |
| GET | /api/cases/:id | 案件详情（含惩罚步骤与扣分记录） |
| PUT | /api/cases/:id | 修改案件 |
| DELETE | /api/cases/:id | 删除案件 |
| POST | /api/cases/:id/start | 开始执行惩罚 |
| POST | /api/cases/:id/complete | 结束惩罚 |
| POST | /api/cases/:id/penalty | 扣分 |
| POST | /api/penalty/:id/revoke | 撤回扣分（需授权密码；密码错误则额外扣 100 分）|

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

- **settings** - 系统配置（admin_username、admin_password、authorization_password 等）
- **members** - 家庭成员（角色：家长/小孩；中英文姓名；国籍；证件信息；学籍；外出权限）
- **rules** - 惩戒规则（名称、描述、分类、分值）
- **records** - 惩戒记录（成员、规则、分值、备注、时间）
- **cases** - 案件（家长成员 + 小孩成员 + 惩罚过程文本 + 执行状态）
- **penalty_points** - 案件扣分记录（可撤回；密码错误撤回额外扣 100 分）

参见 [`server/db/migrations/001_initial_schema.sql`](server/db/migrations/001_initial_schema.sql) 了解完整数据库结构。

## 成员规则

### 角色
- 成员类型分为「**家长**」和「**小孩**」，创建后不可变更。
- 同一自然人可同时拥有两条记录（各一条家长/小孩记录）。
- 为已有「小孩」记录的人添加「家长」记录需验证授权密码。
- 系统必须始终保留至少一名有效家长；唯一家长不可删除。

### 证件信息
- 国籍（ISO 3166-1 alpha-3）、出生日期、英文姓名**必填**。
- 国籍为 CHN/HKG/MAC/TWN 时中文姓名**必填**。
- 主证件类型与国籍的约束关系：

  | 证件类型 | 允许国籍 |
  |---|---|
  | 01 / 91 / 04 | CHN |
  | 11 / 02 | HKG 或 MAC |
  | 21 / 03 | TWN |
  | 31 / 05 / 52 | 非 CHN/HKG/MAC/TWN |

- 辅助证件约束（按主证件类型）：

  | 主证件 | 辅助1 | 辅助2 |
  |---|---|---|
  | 01 / 91 | 不允许 | — |
  | 05 | 不允许 | — |
  | 11 | 02（必填） | 90/92/96/97（必填） |
  | 21 | 03（必填） | 93（必填） |
  | 04 | 94（必填）+ proof_doc_type + proof_issue_country | — |
  | 31 | 05（可选） | — |
  | 02 | 90/92/96/97（HKG 限 90/92，MAC 限 96/97）| — |
  | 03 | 93（可选） | — |
  | 52 | 95 或 98（可选） | — |

- 证件号码校验失败仅提示「证件号码无效」，不暴露具体规则。
- 类型 01/91/11/21/31 号码中包含的性别位与出生日期须与录入值一致。
- 类型 93 号码第二位须为 1（男）或 2（女）且与性别一致。
- 90/92/95 证件号码不允许以 W 或 WX 开头。

### 就读学校
- 学校名称须包含「小学」、「中学」、「大学」或「学院」之一。

## 案件与惩罚

- 每个案件需指定一名**家长**和一名**小孩**（不能是同一自然人）。
- 惩罚过程文本格式：`开始时间|持续分钟|惩罚内容|要求1~5|扣分规则|扣分值`（每行一条，可从 .txt 文件导入）。
- 执行中高亮当前步骤并倒计时显示剩余时长。
- 每条步骤可点击「扣分」按钮执行扣分。
- 撤回扣分需输入授权密码；密码错误额外扣 100 分（原因：作弊）。
