# IDRM 架构文档

> 本目录包含IDRM项目的架构设计文档，是开发者理解系统设计的必读材料。

## 📚 文档清单

### 核心架构

- **[分层架构](./layered-architecture.md)** - 系统四层架构设计，定义各层职责和依赖关系
- **[双ORM模式](./dual-orm-pattern.md)** - Model层双ORM实现，支持GORM和SQLx灵活切换
- **[微服务模式](./microservices-pattern.md)** - 服务拆分原则和通信机制

### 专项设计

- **[API设计指南](./api-design-guide.md)** - RESTful API规范和最佳实践
- **[可观测性架构](./observability.md)** - 日志、追踪、指标的完整方案
- **[数据流设计](./data-flow.md)** - 请求处理和数据传递流程

## 🎯 架构原则

所有架构设计遵循 [constitution.md](../constitution.md) 中定义的技术原则：

1. **分层清晰** - Handler → Logic → Model → Infrastructure
2. **接口抽象** - 面向接口编程，依赖倒置
3. **高内聚低耦合** - 模块职责单一，减少依赖
4. **可测试性** - 设计便于单元测试和集成测试

## 📖 阅读顺序

### 新人入门
1. 先读 [分层架构](./layered-architecture.md)
2. 再读 [双ORM模式](./dual-orm-pattern.md)
3. 然后读 [API设计指南](./api-design-guide.md)

### 深入理解
4. [微服务模式](./microservices-pattern.md)
5. [可观测性架构](./observability.md)
6. [数据流设计](./data-flow.md)

## 🔧 架构演进

### v1.0 (2024-12-24)
- 确立分层架构
- 实现双ORM模式
- 建立可观测性体系

### 未来规划
- 服务网格集成
- 更细粒度的微服务拆分
- 性能优化方案

## 🤝 贡献指南

架构文档的更新需要：
1. 提出Issue讨论
2. 团队Review
3. 更新文档
4. 同步到项目实现

## 📞 联系方式

架构相关问题请联系架构组或在项目Issue中讨论。
