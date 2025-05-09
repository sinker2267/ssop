# 智慧海洋数据分析平台 (SSOP)

智慧海洋数据分析平台(Smart Sea Ocean Platform)是一个专为海洋数据管理、分析和预测而设计的系统。该平台提供用户认证、数据管理、分析功能、系统管理等模块，帮助科研人员和学生更高效地开展海洋数据研究工作。

## 项目架构

项目采用典型的分层架构:

- **表示层**: 处理HTTP请求和响应，包括路由和中间件
- **业务逻辑层**: 实现核心业务逻辑
- **数据访问层**: 处理与数据库的交互
- **领域模型层**: 定义系统的核心实体和关系

### 技术栈

- 后端框架: Gin
- 数据库: MySQL
- ORM: GORM
- 认证: JWT (JSON Web Token)
- 缓存: Redis
- 配置管理: godotenv

## 目录结构

```
.
├── cmd/                # 命令行入口
│   └── api/            # API 服务
├── configs/            # 配置文件
├── internal/           # 内部包
│   ├── config/         # 配置加载
│   ├── handlers/       # HTTP 处理器
│   ├── middleware/     # HTTP 中间件
│   ├── models/         # 数据模型
│   ├── repository/     # 数据存储
│   └── services/       # 业务服务
├── migrations/         # 数据库迁移脚本
├── pkg/                # 公共包
│   ├── logger/         # 日志工具
│   ├── redis/          # Redis客户端
│   ├── response/       # 响应格式
│   └── utils/          # 通用工具
├── storage/            # 数据存储目录
│   ├── datasets/       # 数据集文件
│   └── analysis/       # 分析结果
├── scripts/            # 脚本文件
├── .env                # 环境变量
├── .env.example        # 环境变量示例
├── go.mod              # Go 模块定义
└── README.md           # 项目说明
```

## 代码组件和实现说明

### 数据模型层
系统定义了以下核心数据模型：

- **User**: 用户信息模型，包括用户名、密码哈希、邮箱、角色等
- **Dataset**: 数据集模型，包括元数据、存储位置、创建者等
- **AnalysisTask**: 分析任务模型，记录任务类型、参数、状态等
- **AnalysisResult**: 分析结果模型，存储结果数据、图表信息等
- **SystemSetting**: 系统设置模型，管理全局配置选项
- **AuditLog**: 审计日志模型，记录用户操作

### 数据访问层
系统实现了以下仓库接口：

- **UserRepository**: 用户数据存取和查询
- **DatasetRepository**: 数据集元数据管理
- **AnalysisRepository**: 分析任务和结果管理
- **SystemRepository**: 系统设置和日志管理

### 业务逻辑层
系统包含以下核心服务：

- **AuthService**: 用户认证和授权管理
- **TokenService**: JWT令牌生成和验证
- **UserService**: 用户信息管理
- **DatasetService**: 数据集管理和文件处理
- **AnalysisService**: 分析任务处理和结果计算
- **SystemService**: 系统设置和日志记录

### API控制器
系统提供以下HTTP处理器：

- **AuthHandler**: 处理注册、登录、刷新令牌等请求
- **UserHandler**: 处理用户信息相关请求
- **DatasetHandler**: 处理数据集上传、下载等操作
- **AnalysisHandler**: 处理分析任务和结果管理
- **SystemHandler**: 处理系统设置和日志查询

### 工具函数
系统包含以下实用工具：

- **文件操作**: 文件保存、目录创建、文件类型检查等
- **ID生成**: 生成唯一标识符
- **随机值**: 生成随机字符串和令牌
- **JSON响应**: 统一的API响应格式

### 最近更新

最近对系统进行了以下改进：

1. 修复了response.Success函数调用缺少消息参数的问题，确保所有API响应格式一致
2. 优化了存储路径配置，确保数据文件正确保存
3. 增强了文件操作工具函数，提高文件上传和下载的可靠性

## 已实现功能

### 用户认证模块

- 用户注册
  - 接口: `/api/v1/auth/register`
  - 方法: POST
  - 功能: 创建新用户账号

- 用户登录
  - 接口: `/api/v1/auth/login`
  - 方法: POST
  - 功能: 用户登录获取token

- 游客登录
  - 接口: `/api/v1/auth/guest-login`
  - 方法: POST 
  - 功能: 无需注册，以游客身份登录

- 刷新Token
  - 接口: `/api/v1/auth/refresh-token`
  - 方法: POST
  - 功能: 使用刷新令牌获取新的访问令牌，旧令牌自动加入黑名单

- 退出登录
  - 接口: `/api/v1/auth/logout`
  - 方法: POST
  - 功能: 用户退出登录，将Token加入黑名单

- 获取当前用户信息
  - 接口: `/api/v1/users/current`
  - 方法: GET
  - 功能: 获取当前登录用户信息

### 数据管理模块

- 获取数据集列表
  - 接口: `/api/v1/datasets`
  - 方法: GET
  - 功能: 获取数据集列表，支持分页和筛选

- 获取数据集详情
  - 接口: `/api/v1/datasets/{datasetId}`
  - 方法: GET
  - 功能: 获取指定数据集的详细信息

- 上传数据集
  - 接口: `/api/v1/datasets/upload`
  - 方法: POST
  - 功能: 上传新的数据集文件

- 更新数据集
  - 接口: `/api/v1/datasets/{datasetId}`
  - 方法: PUT
  - 功能: 更新数据集信息

- 删除数据集
  - 接口: `/api/v1/datasets/{datasetId}`
  - 方法: DELETE
  - 功能: 删除指定数据集

- 下载数据集
  - 接口: `/api/v1/datasets/{datasetId}/download`
  - 方法: GET
  - 功能: 下载指定数据集文件

### 分析功能模块

- 任务管理
  - 创建分析任务: `POST /api/v1/analysis/tasks`
  - 获取任务列表: `GET /api/v1/analysis/tasks`
  - 获取任务详情: `GET /api/v1/analysis/tasks/{taskId}`
  - 更新任务信息: `PUT /api/v1/analysis/tasks/{taskId}`
  - 删除分析任务: `DELETE /api/v1/analysis/tasks/{taskId}`
  - 获取任务结果: `GET /api/v1/analysis/tasks/{taskId}/results`

- 结果管理
  - 获取结果详情: `GET /api/v1/analysis/results/{resultId}`
  - 删除分析结果: `DELETE /api/v1/analysis/results/{resultId}`

- 温盐分析功能
  - 温盐时间序列: `GET /api/v1/analysis/temperature-salinity/timeseries`
  - 温盐空间分布: `GET /api/v1/analysis/temperature-salinity/spatial`

### 系统管理模块

- 系统设置
  - 获取所有设置: `GET /api/v1/system/settings`
  - 获取分类设置: `GET /api/v1/system/settings/{category}`
  - 更新系统设置: `PUT /api/v1/system/settings/{key}`

- 审计日志
  - 获取操作日志: `GET /api/v1/system/logs`

### 安全特性

- JWT Token认证
  - 支持访问令牌和刷新令牌
  - 令牌基于用户角色包含权限信息
  
- Redis Token黑名单
  - 登出时自动将Token加入黑名单
  - 刷新令牌时自动将旧令牌加入黑名单
  - 黑名单Token会自动过期，避免Redis内存占用过大

- 权限控制
  - 基于用户角色的权限控制
  - 系统管理功能仅对管理员开放
  - 分析任务和结果仅对所有者开放

## 快速开始

### 前置条件

- Go 1.20+
- MySQL 5.7+
- Redis 6.0+

### 本地开发

1. 克隆仓库
```bash
git clone https://github.com/sinker/ssop.git
cd ssop
```

2. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库和Redis连接等信息
```

3. 创建数据库
```sql
CREATE DATABASE IF NOT EXISTS ssop DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

4. 运行应用
```bash
go run cmd/api/main.go
```

5. 访问API
默认情况下，API服务将在 http://localhost:8080/api/v1 上运行。

### 目录说明

- `cmd/api`: 应用入口
- `internal/models`: 数据模型定义
- `internal/repository`: 数据访问层
- `internal/services`: 业务逻辑层
- `internal/handlers`: HTTP处理器
- `pkg`: 公共工具包

## 用户权限

系统定义了以下用户角色:

- `admin`: 系统管理员，拥有所有权限
- `researcher`: 研究人员，可以读取和上传数据，使用分析功能
- `student`: 学生用户，可以读取数据和使用分析功能
- `guest`: 访客，可以读取部分公开数据

## 开发指南

### 添加新功能

1. 在 `internal/models` 中定义模型
2. 在 `internal/repository` 中实现数据访问
3. 在 `internal/services` 中实现业务逻辑
4. 在 `internal/handlers` 中实现HTTP接口
5. 在 `cmd/api/main.go` 中注册新的路由

### 测试与部署

- 单元测试: `go test ./...`
- 构建: `go build -o ssop cmd/api/main.go`
- 部署: 可使用Docker或直接部署二进制文件

## 贡献指南

欢迎贡献代码和提出建议。请fork仓库并提交pull request。 