-- ============================================================
-- 数据资源目录暂存功能 - 数据库迁移脚本
-- 版本: 1.0.0
-- 日期: 2024-12-26
-- 描述: 创建数据资源目录主表及关联表（草稿功能）
-- ============================================================

-- ============================================================
-- 1. 主表: t_data_catalog
-- 描述: 数据资源目录主表，存储目录的基本信息
-- ============================================================

CREATE TABLE IF NOT EXISTS `t_data_catalog` (
    -- ========== 主键 ==========
    `id` BIGINT(20) UNSIGNED NOT NULL COMMENT '唯一id，雪花算法',

    -- ========== 基础信息 ==========
    `code` VARCHAR(50) NOT NULL COMMENT '目录编码',
    `title` VARCHAR(500) NOT NULL COMMENT '目录名称',

    -- ========== 分类信息 ==========
    `group_id` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '资源属性分类组ID',
    `group_name` VARCHAR(128) DEFAULT '' COMMENT '资源属性分类组名称',
    `theme_id` BIGINT(20) UNSIGNED DEFAULT 0 COMMENT '主题域ID',
    `theme_name` VARCHAR(100) DEFAULT '' COMMENT '主题域名称',

    -- ========== 描述和范围 ==========
    `description` VARCHAR(1000) DEFAULT '' COMMENT '目录描述',
    `data_range` TINYINT DEFAULT 0 COMMENT '数据范围',
    `data_domain` TINYINT DEFAULT 0 COMMENT '数据领域',
    `data_level` TINYINT DEFAULT 0 COMMENT '数据级别',
    `time_range` VARCHAR(100) DEFAULT '' COMMENT '时间范围',

    -- ========== 更新和同步 ==========
    `update_cycle` TINYINT DEFAULT 0 COMMENT '更新周期',
    `other_update_cycle` VARCHAR(100) DEFAULT '' COMMENT '其他更新周期',
    `sync_mechanism` TINYINT DEFAULT 0 COMMENT '数据归集机制：1-增量，2-全量',
    `sync_frequency` VARCHAR(128) DEFAULT '' COMMENT '同步频率',

    -- ========== 基础信息分类 ==========
    `data_kind` INT(32) NOT NULL COMMENT '数据基础信息分类',

    -- ========== 共享开放 ==========
    `shared_type` TINYINT NOT NULL COMMENT '共享类型：1-无条件共享，2-有条件共享，3-不予共享',
    `shared_condition` VARCHAR(255) DEFAULT '' COMMENT '共享条件',
    `column_unshared` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '列是否不共享',
    `open_type` TINYINT NOT NULL COMMENT '开放类型：1-向公众开放，2-不向公众开放',
    `open_condition` VARCHAR(255) DEFAULT '' COMMENT '开放条件',
    `shared_mode` TINYINT NOT NULL COMMENT '共享方式',

    -- ========== 挂接资源统计 ==========
    `view_count` INT(11) DEFAULT 0 COMMENT '逻辑视图数量',
    `api_count` INT(11) DEFAULT 0 COMMENT 'API接口数量',
    `file_count` INT(11) DEFAULT 0 COMMENT '文件数量',
    `physical_deletion` TINYINT DEFAULT NULL COMMENT '物理删除标识',

    -- ========== 流程相关 ==========
    `flow_node_id` VARCHAR(50) DEFAULT '' COMMENT '流程节点ID',
    `flow_node_name` VARCHAR(200) DEFAULT '' COMMENT '流程节点名称',
    `flow_id` VARCHAR(50) DEFAULT '' COMMENT '流程ID',
    `flow_name` VARCHAR(200) DEFAULT '' COMMENT '流程名称',
    `flow_version` VARCHAR(10) DEFAULT '' COMMENT '流程版本',

    -- ========== 部门和用户 ==========
    `department_id` VARCHAR(36) NOT NULL COMMENT '部门ID',
    `source_department_id` VARCHAR(36) NOT NULL COMMENT '来源部门ID',
    `owner_id` VARCHAR(50) NOT NULL COMMENT '所有者ID',
    `owner_name` VARCHAR(128) NOT NULL COMMENT '所有者名称',

    -- ========== 状态字段 ==========
    `publish_status` VARCHAR(20) NOT NULL DEFAULT 'unpublished' COMMENT '发布状态：unpublished-未发布，pub-auditing-发布审核中，published-已发布，pub-reject-发布审核未通过，change-auditing-变更审核中，change-reject-变更审核未通过',
    `online_status` VARCHAR(20) NOT NULL DEFAULT 'notline' COMMENT '上线状态：notline-未上线，online-已上线，offline-已下线，up-auditing-上线审核中，down-auditing-下线审核中，up-reject-上线审核未通过，down-reject-下线审核未通过',
    `audit_type` VARCHAR(50) DEFAULT 'unpublished' COMMENT '审核类型',
    `audit_state` TINYINT DEFAULT NULL COMMENT '审核状态',
    `audit_apply_sn` BIGINT(20) UNSIGNED DEFAULT 0 COMMENT '审核申请流水号',
    `audit_advice` TEXT COMMENT '审核意见',

    -- ========== 草稿相关 ==========
    `draft_id` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '草稿ID，用于关联草稿副本',

    -- ========== 标记字段 ==========
    `source` TINYINT DEFAULT 1 COMMENT '来源',
    `current_version` TINYINT(1) DEFAULT 1 COMMENT '当前版本标识',
    `publish_flag` TINYINT(1) DEFAULT 0 COMMENT '发布标识（9527表示草稿副本）',
    `is_indexed` TINYINT(1) DEFAULT NULL COMMENT '是否已索引',

    -- ========== 更多信息 ==========
    `data_classify` VARCHAR(50) DEFAULT '' COMMENT '数据分类',
    `data_related_matters` VARCHAR(255) DEFAULT '' COMMENT '数据相关事项',
    `business_matters` TEXT COMMENT '业务事项',

    -- ========== 时间戳 ==========
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `published_at` DATETIME DEFAULT NULL COMMENT '发布时间',
    `online_time` DATETIME DEFAULT NULL COMMENT '上线时间',

    -- ========== 更多业务字段 ==========
    `proc_def_key` VARCHAR(128) DEFAULT '' COMMENT '流程定义key',
    `flow_apply_id` VARCHAR(50) DEFAULT '' COMMENT '流程申请ID',
    `apply_num` INT(11) DEFAULT 0 COMMENT '申请数量',
    `explore_job_id` VARCHAR(64) DEFAULT '' COMMENT '探索任务ID',
    `explore_job_version` INT(11) DEFAULT 0 COMMENT '探索任务版本',
    `is_import` TINYINT(1) DEFAULT 0 COMMENT '是否导入',

    PRIMARY KEY (`id`),
    INDEX `idx_source_dept` (`source_department_id`),
    INDEX `idx_publish_status` (`publish_status`),
    INDEX `idx_online_status` (`online_status`),
    INDEX `idx_draft_id` (`draft_id`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据资源目录主表';


-- ============================================================
-- 2. 信息项表: t_data_catalog_column
-- 描述: 数据资源目录信息项（字段定义）
-- ============================================================

CREATE TABLE IF NOT EXISTS `t_data_catalog_column` (
    -- ========== 主键和关联 ==========
    `id` BIGINT(20) UNSIGNED NOT NULL COMMENT '唯一id，雪花算法',
    `catalog_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '目录ID',

    -- ========== 基础信息 ==========
    `column_name` VARCHAR(128) NOT NULL COMMENT '字段名称',
    `column_type` VARCHAR(20) NOT NULL COMMENT '字段类型：string, integer, decimal, date, datetime, boolean, text, blob',
    `column_length` INT(11) DEFAULT 0 COMMENT '字段长度',
    `column_scale` INT(11) DEFAULT 0 COMMENT '字段精度（小数位数）',

    -- ========== 描述信息 ==========
    `column_code` VARCHAR(128) DEFAULT '' COMMENT '字段编码',
    `column_alias` VARCHAR(500) DEFAULT '' COMMENT '字段别名',
    `column_describe` VARCHAR(500) DEFAULT '' COMMENT '字段描述',

    -- ========== 数据属性 ==========
    `data_type` VARCHAR(100) DEFAULT '' COMMENT '数据类型',
    `data_format` VARCHAR(100) DEFAULT '' COMMENT '数据格式',
    `data_unit` VARCHAR(50) DEFAULT '' COMMENT '数据单位',
    `value_range` VARCHAR(500) DEFAULT '' COMMENT '值域范围',
    `is_primary_key` TINYINT(1) DEFAULT 0 COMMENT '是否主键',
    `is_sensitive` TINYINT(1) DEFAULT 0 COMMENT '是否敏感字段',
    `is_required` TINYINT(1) DEFAULT 0 COMMENT '是否必填',
    `is_unique` TINYINT(1) DEFAULT 0 COMMENT '是否唯一',

    -- ========== 扩展信息 ==========
    `dictionary_id` VARCHAR(50) DEFAULT '' COMMENT '字典表ID',
    `standard_code_id` VARCHAR(50) DEFAULT '' COMMENT '标准代码ID',
    `sort_order` INT(11) DEFAULT 0 COMMENT '排序号',

    -- ========== 时间戳 ==========
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',

    PRIMARY KEY (`id`),
    INDEX `idx_catalog_id` (`catalog_id`),
    INDEX `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据资源目录信息项表';


-- ============================================================
-- 3. 分类关联表: t_data_catalog_category
-- 描述: 数据资源目录与资源属性分类的多对多关联
-- ============================================================

CREATE TABLE IF NOT EXISTS `t_data_catalog_category` (
    -- ========== 主键和关联 ==========
    `id` BIGINT(20) UNSIGNED NOT NULL COMMENT '唯一id，雪花算法',
    `catalog_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '目录ID',
    `category_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '资源属性分类节点ID（多选）',

    -- ========== 时间戳 ==========
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    PRIMARY KEY (`id`),
    INDEX `idx_catalog_id` (`catalog_id`),
    INDEX `idx_category_id` (`category_id`),
    UNIQUE KEY `uk_catalog_category` (`catalog_id`, `category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据资源目录分类关联表';


-- ============================================================
-- 4. 挂接资源表: t_data_resource
-- 描述: 数据资源目录挂接的资源（数据表、逻辑视图、API、文件）
-- ============================================================

CREATE TABLE IF NOT EXISTS `t_data_resource` (
    -- ========== 主键和关联 ==========
    `id` BIGINT(20) UNSIGNED NOT NULL COMMENT '唯一id，雪花算法',
    `catalog_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '目录ID',

    -- ========== 资源信息 ==========
    `resource_type` VARCHAR(20) NOT NULL COMMENT '资源类型：database_table-数据表，logical_view-逻辑视图，api_interface-API接口，file-文件',
    `resource_name` VARCHAR(500) NOT NULL COMMENT '资源名称',
    `resource_code` VARCHAR(100) DEFAULT '' COMMENT '资源编码',
    `resource_desc` VARCHAR(1000) DEFAULT '' COMMENT '资源描述',

    -- ========== 资源详情（根据类型不同使用不同字段） ==========
    -- 数据表/逻辑视图
    `database_id` VARCHAR(50) DEFAULT '' COMMENT '数据库ID',
    `table_name` VARCHAR(128) DEFAULT '' COMMENT '表名',
    `table_comment` VARCHAR(500) DEFAULT '' COMMENT '表注释',

    -- API接口
    `api_id` VARCHAR(50) DEFAULT '' COMMENT 'API ID',
    `api_url` VARCHAR(500) DEFAULT '' COMMENT 'API URL',
    `api_method` VARCHAR(10) DEFAULT '' COMMENT 'API方法：GET, POST, PUT, DELETE',
    `api_protocol` VARCHAR(20) DEFAULT '' COMMENT 'API协议：HTTP, HTTPS',

    -- 文件
    `file_id` VARCHAR(50) DEFAULT '' COMMENT '文件ID',
    `file_name` VARCHAR(500) DEFAULT '' COMMENT '文件名',
    `file_url` VARCHAR(500) DEFAULT '' COMMENT '文件URL',
    `file_size` BIGINT(20) DEFAULT 0 COMMENT '文件大小（字节）',

    -- ========== 更多信息 ==========
    `source` VARCHAR(20) DEFAULT 'manual' COMMENT '来源：manual-手动挂接，sync-同步挂接',
    `sort_order` INT(11) DEFAULT 0 COMMENT '排序号',

    -- ========== 时间戳 ==========
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',

    PRIMARY KEY (`id`),
    INDEX `idx_catalog_id` (`catalog_id`),
    INDEX `idx_resource_type` (`resource_type`),
    INDEX `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据资源挂接表';


-- ============================================================
-- 数据插入说明
-- ============================================================
-- 1. 本脚本仅创建表结构，不插入初始数据
-- 2. 草稿标识：publish_flag = 9527 表示该记录为草稿副本
-- 3. 草稿关联：通过 draft_id 字段关联原目录和草稿副本
-- 4. 业务规则：
--    - 同一部门下目录名称不能重复
--    - 高精度类型（decimal）必须指定长度和精度
--    - 逻辑视图最多挂载一个
--    - 已发布目录修改时创建草稿副本
-- ============================================================
