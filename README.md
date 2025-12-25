# IDRM - Intelligent Data Resource Management

> **版本**: v1.0
> **更新日期**: 2024-12-26
> **项目类型**: Go微服务 / AI Native开发

---

## 📖 项目简介

IDRM是一个智能化的数据资源管理平台，采用**AI Native开发理念**，通过Spec文件驱动的开发流程实现高质量、高效率的代码交付。

### 核心特性

- ✅ **AI Native** - Spec文件驱动，AI生成代码，人类审核
- ✅ **四层架构** - Handler → Logic → Model → Database
- ✅ **双ORM模式** - GORM + SQLx，灵活切换
- ✅ **完整测试** - 单元测试覆盖率>85%
- ✅ **高可维护** - 严格遵循架构规范和编码标准

---

## 🚀 快速开始

### 前置要求

- Go 1.21+
- MySQL 8.0+ 或 SQLite
- Redis 7.0+ (可选)
- goctl工具

### 安装goctl

```bash
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

### 运行示例模块（菜单管理）

```bash
# 1. 初始化依赖
go mod tidy

# 2. 编译项目
go build -o bin/menu-api api/menu.go

# 3. 启动服务
./bin/menu-api -f api/etc/menu.yaml

# 4. 测试API
curl http://localhost:8080/api/v1/sys/menu/tree
```

---

## 📚 文档导航

### 核心文档

| 文档 | 说明 | 位置 |
|-----|------|------|
| **项目宪章** | 项目开发的基本原则和规范 | [constitution.md](constitution.md) |
| **Spec-Kit指南** | Spec驱动开发流程完整指南 | [SPEC_KIT_操作指南.md](SPEC_KIT_操作指南.md) |
| **架构文档** | 分层架构、双ORM模式等 | [.specs/architecture/](.specs/architecture/) |
| **编码规范** | 代码风格、命名规范等 | [.specs/coding-standards/](.specs/coding-standards/) |

### 实施案例

- **菜单管理模块** - 完整的Spec-Kit开发示例
  - Spec文件: [specs/specification-menu.md](specs/specification-menu.md)
  - 实施总结: [specs/IMPLEMENTATION_SUMMARY.md](specs/IMPLEMENTATION_SUMMARY.md)
  - 模块文档: [model/system/menu/README.md](model/system/menu/README.md)

---

## 🎯 Spec-Kit开发流程

IDRM采用**四阶段开发流程**：

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Specify   │ => │    Plan     │ => │   Generate  │ => │   Verify    │
│  功能规范    │    │  技术方案    │    │   代码生成   │    │   测试验证   │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

**详细步骤请参考**: [SPEC_KIT_操作指南.md](SPEC_KIT_操作指南.md)

---

## 🏗️ 项目结构

```
idrm-ai/
├── specs/                    # Spec文件（需求定义）
├── api/                      # API层（Go-Zero）
│   ├── internal/
│   │   ├── handler/          # HTTP处理器（goctl生成）
│   │   ├── logic/            # 业务逻辑（AI生成）
│   │   └── svc/              # 服务上下文
├── model/                    # Model层（双ORM）
│   └── {domain}/{module}/
│       ├── interface.go      # 统一接口
│       ├── gorm_dao.go       # GORM实现
│       └── sqlx_model.go     # SQLx实现
├── pkg/                      # 公共包
│   ├── errorx/               # 错误处理
│   ├── middleware/           # 中间件
│   └── response/             # 响应封装
└── deploy/migrations/        # 数据库迁移
```

---

## 🛠️ 常用命令

### Makefile命令

```bash
make help          # 查看所有命令
make build         # 编译项目
make test          # 运行测试
make run           # 运行服务
make check         # 完整检查
make docker-build  # 构建Docker镜像
```

### Go命令

```bash
go mod tidy              # 整理依赖
go build ./...            # 编译检查
go test ./... -v          # 运行测试
go fmt ./...               # 格式化代码
```

### Git工作流

```bash
# 创建功能分支
git checkout -b feature/{module}

# 提交代码
git add .
git commit -m "feat({module}): implement {module}"

# 推送
git push origin feature/{module}
```

---

## 📊 技术栈

### 后端框架
- **Go 1.21+** - 编程语言
- **Go-Zero v1.9+** - 微服务框架
- **GORM v1.31+** - ORM框架
- **SQLx** - 轻量级ORM

### 数据库
- **MySQL 8.0+** - 主数据库
- **SQLite** - 开发/测试数据库
- **Redis 7.0+** - 缓存

### 可观测性
- **OpenTelemetry** - 链路追踪
- **Logx** - 结构化日志
- **Prometheus** - 指标收集

---

## 🎓 开发规范

### 必读文档

1. [项目宪章](constitution.md) - 了解项目基本原则
2. [SPEC_KIT_操作指南](SPEC_KIT_操作指南.md) - 学习开发流程
3. [分层架构](.specs/architecture/layered-architecture.md) - 理解架构设计

### 核心原则

- ✅ **Spec优先** - 所有功能从Spec文件开始
- ✅ **质量优先** - 代码质量重于开发速度
- ✅ **AI生成，人类审核** - AI负责生成，人类负责审查
- ✅ **契约高于一切** - Spec文件是"单一真理来源"

### 禁止操作

- ❌ 跳过Spec直接写代码
- ❌ 修改goctl生成的文件
- ❌ Logic层直接访问数据库
- ❌ 硬编码配置信息
- ❌ 忽略错误处理

---

## 🧪 测试

### 运行测试

```bash
# 所有测试
go test ./... -v

# Model层测试
go test ./model/... -v -cover

# 查看覆盖率
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 测试要求

- 核心业务逻辑：>80%覆盖率
- 工具函数：>90%覆盖率
- 所有代码必须有单元测试

---

## 🐳 Docker支持

### 使用Docker Compose

```bash
# 启动所有服务（MySQL + Redis + API）
docker-compose up -d

# 查看日志
docker-compose logs -f menu-api

# 停止服务
docker-compose down
```

---

## 📞 联系方式

- **项目仓库**: [GitHub](https://github.com/your-org/idrm-ai)
- **问题反馈**: [Issues](https://github.com/your-org/idrm-ai/issues)
- **文档**: [Wiki](https://github.com/your-org/idrm-ai/wiki)

---

## 📝 更新日志

### v1.0 (2024-12-26)
- ✅ 初始版本发布
- ✅ 完成菜单管理模块
- ✅ 建立Spec-Kit开发流程
- ✅ 完善项目文档

---

## 📄 许可证

[MIT License](LICENSE)

---

**IDRM项目遵循AI Native开发理念，通过Spec文件驱动的开发流程，实现高质量、高效率的代码交付。**

**欢迎贡献代码和建议！**
