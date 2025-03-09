# SSOP 后端系统

基于 Golang 的后端系统，采用领域驱动设计(DDD)架构。

## 技术栈

- Golang
- Gin Web 框架
- MySQL 数据库
- Redis 缓存
- JWT 认证
- Swagger API 文档

## 项目架构

项目采用领域驱动设计(DDD)架构，主要分为以下几层：

### 应用层 (Application)
- 处理用户请求，协调领域层和基础设施层
- 包含API接口、中间件、DTO等
- 不包含业务逻辑，仅用于协调各层

### 领域层 (Domain)
- 包含核心业务逻辑
- 领域模型、领域服务、领域事件
- 领域对象包含业务规则和行为

### 基础设施层 (Infrastructure)
- 提供技术支持
- 数据库访问、缓存、第三方服务集成
- 实现领域层定义的接口

### 接口层 (Interfaces)
- 处理用户界面和外部系统集成
- 定义API接口和格式

## 目录结构

```
ssop_rear/
├── cmd/                  # 应用入口
│   └── api/              # API服务入口
├── internal/             # 内部代码，不对外暴露
│   ├── application/      # 应用层
│   │   ├── dto/          # 数据传输对象
│   │   └── service/      # 应用服务
│   ├── domain/           # 领域层
│   │   ├── entity/       # 实体
│   │   ├── repository/   # 仓储接口
│   │   ├── service/      # 领域服务
│   │   └── valueobject/  # 值对象
│   ├── infrastructure/   # 基础设施层
│   │   ├── persistence/  # 持久化实现
│   │   ├── cache/        # 缓存实现
│   │   └── middleware/   # 中间件
│   └── interfaces/       # 接口层
│       ├── api/          # API接口
│       └── dto/          # 数据传输对象
├── pkg/                  # 公共包
│   ├── config/           # 配置
│   ├── logger/           # 日志
│   ├── errors/           # 错误处理
│   └── utils/            # 工具函数
├── configs/              # 配置文件
├── scripts/              # 脚本
├── docs/                 # 文档
│   └── swagger/          # Swagger文档
└── Makefile              # 构建脚本
```

## 如何启动

1. 安装依赖
```bash
go mod tidy
```

2. 编译
```bash
make build
```

3. 配置
修改 configs/config.yaml 文件

4. 运行
```bash
./bin/api
```

## API文档

启动服务后，访问 http://localhost:8080/swagger/index.html

## 数据库迁移

```bash
make migrate
```

## 测试

```bash
make test
```

## 代码规范

- 使用 gofmt 格式化代码
- 遵循 Go 官方代码规范
- 使用依赖注入管理依赖
- 编写单元测试确保代码质量

## 项目规划与改进方向

- 实现完整的用户认证和授权系统
- 添加更多单元测试和集成测试
- 优化数据库查询和缓存策略
- 实现分布式部署支持 