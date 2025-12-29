# Specification: 数据资源目录暂存模块

## 1. 业务目标

实现数据资源目录的暂存功能，支持未发布目录的创建暂存和已发布目录的变更暂存，通过草稿机制实现版本控制，确保数据一致性。

## 2. 核心功能

### 2.1 暂存场景
- **创建暂存**: 首次创建全新的目录草稿（CatalogID 为空）
- **变更暂存**: 对已发布目录进行编辑，创建或更新草稿副本
- **草稿更新**: 对现有草稿进行内容更新

### 2.2 数据字段
- **基本属性**: 目录名称、描述、来源部门、所属领域
- **分类信息**: 资源属性分类节点ID（多选）
- **共享开放**: 共享类型（无条件/有条件/不予共享）、开放类型
- **信息项**: 字段列表（名称、类型、长度、精度、描述）
- **挂接资源**: 逻辑视图、数据表、API接口等（最多一个逻辑视图）
- **更多信息**: 同步机制、更新频率、编码规则

### 2.3 状态管理
- **草稿状态** (unpublished): 可直接编辑，无需审核
- **已发布状态** (published): 修改时创建草稿副本，原目录不变
- **草稿副本机制**: 通过 DraftID 关联原目录和草稿副本

## 3. 技术约束

### 3.1 存储与框架
- **存储**: MySQL 8.0+
- **框架**: Go-Zero (go-zero v1.9+)
- **架构**: 四层架构 (Handler → Logic → Model)

### 3.2 API规范
- **接口路径**: `/api/v1/data-catalog/draft`
- **请求方法**: POST
- **请求格式**: JSON
- **响应格式**: 统一JSON响应

### 3.3 数据模型

#### 表结构概览
- **主表**: `t_data_catalog` - 数据资源目录表
- **关联表**: `t_data_catalog_column` - 目录信息项表
- **关联表**: `t_data_catalog_category` - 分类关联表
- **关联表**: `t_data_resource` - 数据资源表

#### 状态字段
- `publish_status`: 发布状态 (unpublished/pub-auditing/published/pub-reject/change-auditing/change-reject)
- `online_status`: 上线状态 (notline/online/offline/up-auditing/down-auditing/up-reject/down-reject)
- `draft_id`: 草稿ID（用于关联草稿副本）

#### ID生成
- 雪花算法生成唯一ID

---

### 3.4 数据库表结构

#### 3.4.1 主表: t_data_catalog

```sql
CREATE TABLE IF NOT EXISTS `t_data_catalog` (
    `id` bigint(20) unsigned NOT NULL COMMENT '唯一id，雪花算法',
    `code` varchar(50) NOT NULL COMMENT '目录编码',
    `title` varchar(500) NOT NULL COMMENT '目录名称',
    `group_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '数据资源目录分类ID',
    `group_name` varchar(128) NOT NULL DEFAULT '' COMMENT '数据资源目录分类名称',
    `theme_id` bigint(20) unsigned DEFAULT NULL COMMENT '主题分类ID',
    `theme_name` varchar(100) DEFAULT NULL COMMENT '主题分类名称',
    `forward_version_id` bigint(20) unsigned DEFAULT NULL COMMENT '当前目录前一版本目录ID',
    `description` varchar(1000) DEFAULT NULL COMMENT '资源目录描述',
    `data_range` tinyint(2) DEFAULT NULL COMMENT '数据范围：字典DM_DATA_SJFW，01全市 02市直 03区县',
    `update_cycle` tinyint(2) DEFAULT NULL COMMENT '更新频率 参考数据字典：GXZQ，1不定时 2实时 3每日 4每周 5每月 6每季度 7每半年 8每年 9其他',
    `other_update_cycle` varchar(100) DEFAULT NULL COMMENT '其他更新频率',
    `data_kind` tinyint(2) NOT NULL COMMENT '基础信息分类 1 人 2 地 4 事 8 物 16 组织 32 其他 可组合，如 人和地 即 1|2 = 3',
    `shared_type` tinyint(1) NOT NULL COMMENT '共享属性 1 无条件共享 2 有条件共享 3 不予共享',
    `shared_condition` varchar(255) DEFAULT NULL COMMENT '共享条件',
    `column_unshared` tinyint(1) NOT NULL COMMENT '信息项不予共享',
    `open_type` tinyint(1) NOT NULL COMMENT '开放属性 1 向公众开放 2 不向公众开放',
    `open_condition` varchar(255) DEFAULT NULL COMMENT '开放条件',
    `shared_mode` tinyint(1) NOT NULL COMMENT '共享方式 1 共享平台方式 2 邮件方式 3 介质方式',
    `physical_deletion` tinyint(1) DEFAULT NULL COMMENT '挂接实体资源是否存在物理删除(1 是 ; 0 否)',
    `sync_mechanism` tinyint(1) DEFAULT NULL COMMENT '数据归集机制(1 增量 ; 2 全量) ----归集到数据中台',
    `sync_frequency` varchar(128) DEFAULT NULL COMMENT '数据归集频率 ----归集到数据中台',
    `view_count` smallint(4) NOT NULL DEFAULT 0 COMMENT '挂接逻辑视图数量',
    `api_count` smallint(4) NOT NULL DEFAULT 0 COMMENT '挂接接口数量',
    `file_count` smallint(4) NOT NULL DEFAULT 0 COMMENT '挂接文件资源数量',
    `flow_node_id` varchar(50) DEFAULT NULL COMMENT '目录当前所处审核流程结点ID',
    `flow_node_name` varchar(200) DEFAULT NULL COMMENT '目录当前所处审核流程结点名称',
    `flow_id` varchar(50) DEFAULT NULL COMMENT '审批流程实例ID',
    `flow_name` varchar(200) DEFAULT NULL COMMENT '审批流程名称',
    `flow_version` varchar(10) DEFAULT NULL COMMENT '审批流程版本',
    `department_id` char(36) NOT NULL COMMENT '所属部门ID',
    `created_at` datetime(3) NOT NULL DEFAULT current_timestamp(3) COMMENT '创建时间',
    `creator_uid` varchar(50) NOT NULL DEFAULT '' COMMENT '创建用户ID',
    `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
    `updater_uid` varchar(50) DEFAULT NULL COMMENT '更新用户ID',
    `source` tinyint(1) NOT NULL DEFAULT 1 COMMENT '数据来源 1 认知平台自动创建 2 人工创建',
    `table_type` tinyint(1) DEFAULT NULL COMMENT '库表类型 1 贴源表 2 标准表',
    `current_version` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否现行版本 0 否 1 是',
    `publish_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否发布到超市 (1 是 ; 0 否)',
    `data_kind_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '基础信息分类是否智能推荐 (1 是 ; 0 否)',
    `label_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '标签是否智能推荐 (1 是 ; 0 否)',
    `src_event_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '来源业务场景是否智能推荐 (1 是 ; 0 否)',
    `rel_event_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '关联业务场景是否智能推荐 (1 是 ; 0 否)',
    `system_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '关联信息系统是否智能推荐 (1 是 ; 0 否)',
    `rel_catalog_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '关联目录是否智能推荐 (1 是 ; 0 否)',
    `published_at` datetime(3) DEFAULT NULL COMMENT '上线发布时间',
    `is_indexed` tinyint(1) DEFAULT NULL COMMENT '是否已建ES索引，0 否 1 是',
    `audit_apply_sn` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '发起审核申请序号',
    `audit_advice` text DEFAULT NULL COMMENT '审核意见，仅驳回时有用',
    `owner_id` varchar(50) NOT NULL DEFAULT '' COMMENT '目录数据owner的用户ID',
    `owner_name` varchar(128) NOT NULL DEFAULT '' COMMENT '目录数据owner的用户名称',
    `is_canceled` tinyint(1) DEFAULT NULL COMMENT '目录下线时针对该目录的关联申请是否已撤销，0 待撤销 1 已撤销',
    `proc_def_key` varchar(128) NOT NULL DEFAULT '' COMMENT '审核流程key',
    `flow_apply_id` varchar(50) DEFAULT '' COMMENT '审核流程ID',
    `online_status` varchar(20) NOT NULL DEFAULT 'notline' COMMENT '接口状态 未上线 notline、已上线 online、已下线offline、上线审核中up-auditing、下线审核中down-auditing、上线审核未通过up-reject、下线审核未通过down-reject、已下线（上线审核中）offline-up-auditing、已下线（上线审核未通过）offline-up-reject',
    `online_time` datetime DEFAULT NULL COMMENT '上线时间',
    `audit_type` varchar(50) NOT NULL DEFAULT 'unpublished' COMMENT '审核类型 unpublished 未发布 af-data-view-online上线审核 af-data-view-offline 下线审核',
    `audit_state` tinyint(1) DEFAULT NULL COMMENT '审核状态，1 审核中  2 通过  3 驳回',
    `publish_status` varchar(20) NOT NULL DEFAULT 'unpublished' COMMENT '发布状态 未发布unpublished 、发布审核中pub-auditing、已发布published、发布审核未通过pub-reject、变更审核中change-auditing、变更审核未通过change-reject',
    `app_scene_classify` tinyint(2) DEFAULT NULL COMMENT '应用场景分类 1 政务服务、2 公共服务、3 监管、4 其他',
    `source_department_id` char(36) NOT NULL COMMENT '数据资源来源部门id',
    `data_related_matters` varchar(255) NOT NULL COMMENT '数据所属事项',
    `business_matters` text NOT NULL COMMENT '业务事项',
    `data_classify` varchar(50) NOT NULL COMMENT '数据分级 标签',
    `data_domain` tinyint(1) COMMENT '数据所在领域',
    `data_level` tinyint(1) COMMENT '数据所在层级',
    `time_range` varchar(100) COMMENT '数据时间范围',
    `provider_channel` tinyint(1) COMMENT '提供渠道',
    `administrative_code` tinyint(1) COMMENT '行政区划代码',
    `central_department_code` tinyint(1) COMMENT '中央业务指导部门代码',
    `processing_level` varchar(100) COMMENT '数据加工程度',
    `catalog_tag` tinyint(1) COMMENT '目录标签',
    `is_electronic_proof` tinyint(1) COMMENT '是否电子证明编码',
    `other_app_scene_classify` varchar(100) COMMENT '其他应用场景分类',
    `draft_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '草稿id',
    `apply_num` int(11) NOT NULL COMMENT '申请量',
    `explore_job_id` VARCHAR(64) DEFAULT NULL COMMENT '探查作业ID',
    `explore_job_version` INT(11) DEFAULT NULL COMMENT '探查作业版本',
    `operation_authorized` tinyint(1) COMMENT '是否可授权运营字段',
    `is_import` TINYINT(1) DEFAULT 0 COMMENT '是否导入',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '目录表';
```

#### 3.4.2 信息项表: t_data_catalog_column

```sql
CREATE TABLE IF NOT EXISTS `t_data_catalog_column` (
    `primary_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `id` bigint(20) unsigned NOT NULL COMMENT '唯一id，雪花算法',
    `catalog_id` bigint(20) unsigned NOT NULL COMMENT '数据资源目录ID',
    `technical_name` varchar(255) NOT NULL COMMENT '技术名称',
    `business_name` varchar(255) DEFAULT NULL COMMENT '业务名称',
    `source_id` char(36) NOT NULL COMMENT '来源id',
    `data_format` TINYINT(2) DEFAULT NULL COMMENT '字段类型 0:数字型 1:字符型 2:日期型 3:日期时间型 5:布尔型 6:其他 7:小数型 8:高精度型 9:时间型',
    `data_length` int(11) DEFAULT NULL COMMENT '字段长度',
    `data_precision` TINYINT(2) DEFAULT NULL COMMENT '数据精度',
    `ranges` varchar(200) DEFAULT NULL COMMENT '字段值域',
    `shared_type` tinyint(1) DEFAULT NULL COMMENT '共享属性 1 无条件共享 2 有条件共享 3 不予共享',
    `open_type` tinyint(1) DEFAULT NULL COMMENT '开放属性 1 向公众开放 2 不向公众开放',
    `timestamp_flag` tinyint(1) DEFAULT NULL COMMENT '是否时间戳(1 是 ; 0 否)',
    `primary_flag` tinyint(1) DEFAULT NULL COMMENT '是否主键(1 是 ; 0 否)',
    `null_flag` tinyint(1) DEFAULT NULL COMMENT '是否为空(1 是 ; 0 否)',
    `classified_flag` tinyint(1) DEFAULT NULL COMMENT '是否涉密属性(1 是 ; 0 否)',
    `sensitive_flag` tinyint(1) DEFAULT NULL COMMENT '是否敏感属性(1 是 ; 0 否)',
    `description` varchar(2048) NOT NULL DEFAULT '' COMMENT '字段描述，对应元数据field_comment',
    `shared_condition` varchar(255) DEFAULT NULL COMMENT '共享条件',
    `open_condition` varchar(255) DEFAULT NULL COMMENT '开放条件',
    `ai_description` varchar(2048) DEFAULT '' COMMENT 'AI数据理解下生成的字段描述',
    `standard_code` varchar(30) NULL COMMENT '数据标准code',
    `code_table_id` varchar(30) NULL COMMENT '码表ID',
    `source_system` varchar(255) NULL COMMENT '来源系统',
    `source_system_level` tinyint(1) NULL COMMENT '来源系统分级 1 自建自用 2 国直 3省直 4市直 5县直',
    `info_item_level` tinyint(1) NULL COMMENT '信息项分级 1级 2级 3级 4级',
    `index` int NOT NULL COMMENT '信息项顺序',
    KEY `id_key` (`id`),
    UNIQUE KEY `t_data_catalog_column_un` (`catalog_id`, `technical_name`),
    PRIMARY KEY (`primary_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '目录关联信息项表';
```

#### 3.4.3 分类关联表: t_data_catalog_category

```sql
CREATE TABLE IF NOT EXISTS `t_data_catalog_category` (
    `id` bigint(20) NOT NULL,
    `category_id` char(36) NOT NULL COMMENT '类目ID',
    `category_type` tinyint(4) NOT NULL COMMENT '类目类型 1:所属部门 2:信息系统 3:所属主题 4 ：自定义',
    `catalog_id` bigint(20) NOT NULL COMMENT '数据目录ID',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '类目目录关联表';
```

#### 3.4.4 数据资源表: t_data_resource

```sql
CREATE TABLE IF NOT EXISTS `t_data_resource` (
    `id` bigint(20) NOT NULL COMMENT '标识',
    `resource_id` varchar(125) NOT NULL COMMENT '数据资源id',
    `name` varchar(255) NOT NULL COMMENT '数据资源名称',
    `code` varchar(255) NOT NULL COMMENT '统一编目编码',
    `type` tinyint(4) NOT NULL COMMENT '数据资源类型 枚举值 1：逻辑视图 2：接口 3:文件资源',
    `view_id` char(36) COMMENT '数据资源类型 为 2：接口 时候类型为接口生成方式来源视图id',
    `interface_count` int COMMENT '生成接口数量',
    `department_id` char(36) NOT NULL DEFAULT '' COMMENT '所属部门',
    `subject_id` char(36) NOT NULL DEFAULT '' COMMENT '所属主题',
    `request_format` varchar(100) COMMENT '请求报文格式',
    `response_format` varchar(100) COMMENT '响应报文格式',
    `scheduling_plan` tinyint(1) COMMENT '调度计划 1 一次性、2按分钟、3按天、4按周、5按月',
    `interval` tinyint(1) COMMENT '间隔',
    `time` varchar(100) COMMENT '时间',
    `catalog_id` bigint(20) NOT NULL COMMENT '数据资源目录ID',
    `publish_at` datetime(3) NOT NULL COMMENT '发布时间',
    `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT ' 视图状态,1正常,2删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `t_data_resource.code.uk` (`code`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '数据资源表';
```

#### 3.4.5 表关系说明

```
t_data_catalog (主表)
    ├── 1:N ── t_data_catalog_column (通过 catalog_id)
    ├── 1:N ── t_data_catalog_category (通过 catalog_id)
    └── 1:N ── t_data_resource (通过 catalog_id)
```

- **t_data_catalog**: 主表，存储数据资源目录的基本信息
- **t_data_catalog_column**: 子表，存储目录下的信息项（字段），通过 `catalog_id` 外键关联
- **t_data_catalog_category**: 关联表，存储目录的分类关联信息，通过 `catalog_id` 外键关联
- **t_data_resource**: 关联表，存储挂载的数据资源（逻辑视图、接口、文件资源），通过 `catalog_id` 外键关联

#### 3.4.6 主要索引

- **t_data_catalog**: 主键索引 `id`
- **t_data_catalog_column**: 主键索引 `primary_id`，唯一索引 `(catalog_id, technical_name)`
- **t_data_catalog_category**: 主键索引 `id`
- **t_data_resource**: 主键索引 `id`，唯一索引 `code`

### 3.5 业务规则
- **名称唯一性**: 同一部门下目录名称不能重复
- **字段验证**: 类型必填、高精度类型需指定长度和精度
- **资源约束**: 逻辑视图最多挂载一个
- **部门验证**: 来源部门必须存在
- **草稿副本**: DraftFlag = 9527 标识草稿副本

## 4. 异常处理

### 4.1 错误码规范

使用**字符串格式**错误码（与原项目保持一致），格式：`idrm.DataCatalog.{错误名称}`

### 4.2 错误码定义

```go
const (
    // 数据资源目录错误
    DataCatalogNotFound           = "idrm.DataCatalog.CatalogNotFound"
    DataCatalogDepartmentNotFound = "idrm.DataCatalog.DataCatalogDepartmentNotFound"
    DataCatalogInfoSystemNotFound = "idrm.DataCatalog.DataCatalogInfoSystemNotFound"
    DataCatalogSubjectNotFound    = "idrm.DataCatalog.DataCatalogSubjectNotFound"
    CategoryNodeIdNotFound        = "idrm.DataCatalog.CategoryNodeIdNotFound"
    CodeTableIDsVerifyFail        = "idrm.DataCatalog.CodeTableIDsVerifyFail"
    StandardCodesVerifyFail       = "idrm.DataCatalog.StandardCodesVerifyFail"
    DeleteDataCatalogFail         = "idrm.DataCatalog.DeleteDataCatalogFail"
    CatalogNameRepeat             = "idrm.DataCatalog.CatalogNameRepeat"
    DataResourceNotExist          = "idrm.DataCatalog.DataResourceNotExist"
    DataResourceTypeNotSupport    = "idrm.DataCatalog.DataResourceTypeNotSupport"
    ImportFileNotExist            = "idrm.DataCatalog.ImportFileNotExist"
    OnlineNeedReport              = "idrm.DataCatalog.OnlineNeedReport"
    ImportInvalidError            = "idrm.DataCatalog.ImportInvalidError"
)
```

### 4.3 错误映射表

```go
var dataCatalogErrorMap = map[string]ErrorInfo{
    DataCatalogNotFound: {
        Description: "数据资源目录不存在",
        Solution:    "请检查",
    },
    DataCatalogDepartmentNotFound: {
        Description: "数据资源目录关联部门不存在",
        Solution:    "请检查",
    },
    DataCatalogInfoSystemNotFound: {
        Description: "数据资源目录关联信息系统不存在",
        Solution:    "请检查",
    },
    DataCatalogSubjectNotFound: {
        Description: "数据资源目录关联主题域不存在",
        Solution:    "请检查",
    },
    CategoryNodeIdNotFound: {
        Description: "资源属性分类不存在",
        Solution:    "请检查",
    },
    CodeTableIDsVerifyFail: {
        Description: "参数值校验不通过:码表不存在",
    },
    StandardCodesVerifyFail: {
        Description: "参数值校验不通过:标准不存在",
    },
    DeleteDataCatalogFail: {
        Description: "删除数据目录失败",
    },
    CatalogNameRepeat: {
        Description: "数据目录名称重复",
    },
    DataResourceNotExist: {
        Description: "数据资源已挂载或不存在",
    },
    DataResourceTypeNotSupport: {
        Description: "数据资源类型不支持",
    },
    ImportFileNotExist: {
        Description: "导入文件不存在",
    },
    OnlineNeedReport: {
        Description: "生成理解报告后，才可以发起上线",
        Solution:    "请检查",
    },
    ImportInvalidError: {
        Description: "导入[sheet]页，第[word]行数据校验失败",
    },
}
```

---

**参考规范**:
- 项目宪章: `constitution.md`
- SPEC-Kit指南: `SPEC_KIT_操作指南.md`
- 错误码定义: `pkg/errorx/errorx.go`
- 原项目错误码: `project_bundle/data-catalog/common/errorcode/data_catalog.go`
