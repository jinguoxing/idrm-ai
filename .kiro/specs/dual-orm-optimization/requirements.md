# 双ORM优化需求文档

## 项目概述

对IDRM项目的双ORM模式进行优化，提升性能、可维护性和可观测性。

## 术语表

- **ORM**: Object-Relational Mapping，对象关系映射
- **GORM**: Go语言的ORM库，功能丰富
- **SQLx**: Go语言的SQL扩展库，轻量高效
- **工厂模式**: 创建对象的设计模式
- **连接池**: 数据库连接的缓存池

## 需求

### 需求1：配置驱动的ORM选择策略

**用户故事**: 作为系统管理员，我希望能够通过配置文件灵活控制ORM选择策略，以便在不同环境下优化性能。

#### 验收标准

1. WHEN 配置strategy为"gorm_first" THEN 系统应优先使用GORM，GORM不可用时降级到SQLx
2. WHEN 配置strategy为"sqlx_first" THEN 系统应优先使用SQLx，SQLx不可用时降级到GORM  
3. WHEN 配置strategy为"gorm_only" THEN 系统应仅使用GORM，不可用时返回错误
4. WHEN 配置strategy为"sqlx_only" THEN 系统应仅使用SQLx，不可用时返回错误
5. WHEN 配置strategy为"auto" THEN 系统应根据操作类型智能选择最优ORM

### 需求2：数据库连接池优化

**用户故事**: 作为系统架构师，我希望能够精确控制数据库连接池参数，以便优化系统资源使用和性能。

#### 验收标准

1. WHEN 配置maxOpenConns THEN 系统应限制最大打开连接数不超过配置值
2. WHEN 配置maxIdleConns THEN 系统应维护指定数量的空闲连接
3. WHEN 配置connMaxLifetime THEN 连接应在指定时间后被回收
4. WHEN 配置connMaxIdleTime THEN 空闲连接应在指定时间后被关闭
5. WHEN 连接池配置无效 THEN 系统应使用合理的默认值并记录警告

### 需求3：统一错误处理

**用户故事**: 作为开发者，我希望GORM和SQLx返回统一的错误类型，以便简化错误处理逻辑。

#### 验收标准

1. WHEN GORM返回RecordNotFound错误 THEN 系统应转换为统一的ErrNotFound
2. WHEN SQLx返回NoRows错误 THEN 系统应转换为统一的ErrNotFound
3. WHEN GORM返回DuplicatedKey错误 THEN 系统应转换为统一的ErrAlreadyExists
4. WHEN 发生数据库连接错误 THEN 系统应返回统一的数据库错误码
5. WHEN 错误转换失败 THEN 系统应保留原始错误信息

### 需求4：性能监控指标

**用户故事**: 作为运维人员，我希望能够监控双ORM的使用情况和性能指标，以便优化系统配置。

#### 验收标准

1. WHEN 执行GORM操作 THEN 系统应记录GORM查询次数和平均延迟
2. WHEN 执行SQLx操作 THEN 系统应记录SQLx查询次数和平均延迟  
3. WHEN 发生ORM切换 THEN 系统应记录切换次数和原因
4. WHEN 查询指标 THEN 系统应返回实时的性能统计数据
5. WHEN 重置指标 THEN 系统应清零所有统计计数器

### 需求5：配置结构增强

**用户故事**: 作为配置管理员，我希望配置文件结构清晰合理，支持所有新增的配置选项。

#### 验收标准

1. WHEN 解析配置文件 THEN 系统应正确读取数据库连接池配置
2. WHEN 解析配置文件 THEN 系统应正确读取ORM策略配置
3. WHEN 配置项缺失 THEN 系统应使用合理的默认值
4. WHEN 配置项格式错误 THEN 系统应返回清晰的错误信息
5. WHEN 配置更新 THEN 系统应支持热重载（可选）

## 技术约束

- 必须保持现有Model接口的向后兼容性
- 必须支持现有的GORM和SQLx实现
- 必须遵循项目的分层架构原则
- 必须包含完整的单元测试
- 配置变更不应影响现有功能

## 非功能需求

- 性能：ORM选择开销应小于1ms
- 可靠性：错误处理覆盖率100%
- 可维护性：代码复杂度不应显著增加
- 可观测性：关键操作必须有日志记录

## 验收标准

- 所有单元测试通过
- 代码覆盖率不低于80%
- 性能测试显示无显著性能退化
- 配置文档完整且准确