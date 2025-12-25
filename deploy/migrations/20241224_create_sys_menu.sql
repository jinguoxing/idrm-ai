-- 菜单管理模块 - 数据库表创建脚本
-- 版本: v1.0
-- 日期: 2024-12-24
-- 作者: IDRM Team

-- 如果表已存在，先删除（开发环境）
-- DROP TABLE IF EXISTS sys_menu;

-- 创建系统菜单表
CREATE TABLE IF NOT EXISTS sys_menu (
    -- 主键
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '菜单ID',

    -- 基本信息
    parent_id BIGINT NOT NULL DEFAULT 0 COMMENT '父级菜单ID，0表示根菜单',
    name VARCHAR(50) NOT NULL COMMENT '菜单名称',
    route_path VARCHAR(200) DEFAULT NULL COMMENT '路由路径',
    component_path VARCHAR(200) DEFAULT NULL COMMENT '组件路径',
    icon VARCHAR(100) DEFAULT NULL COMMENT '菜单图标',

    -- 排序和权限
    sort_order INT NOT NULL DEFAULT 0 COMMENT '排序号，数字越小越靠前',
    perm_tag VARCHAR(100) DEFAULT NULL COMMENT '权限标识，用于权限控制',

    -- 状态
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',

    -- 时间戳
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME DEFAULT NULL COMMENT '删除时间（软删除）',

    -- 索引
    UNIQUE KEY uk_id (id),
    KEY idx_parent_id (parent_id),
    KEY idx_perm_tag (perm_tag),
    KEY idx_status (status),
    KEY idx_sort_order (sort_order),
    KEY idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统菜单表';

-- 插入初始数据（示例）
INSERT INTO sys_menu (parent_id, name, route_path, component_path, icon, sort_order, perm_tag, status) VALUES
(0, '系统管理', '/system', 'Layout', 'setting', 100, 'system', 1),
(1, '菜单管理', '/system/menu', 'system/menu/index', 'menu', 1, 'system:menu', 1),
(1, '用户管理', '/system/user', 'system/user/index', 'user', 2, 'system:user', 1),
(1, '角色管理', '/system/role', 'system/role/index', 'team', 3, 'system:role', 1),
(0, '数据资源', '/resource', 'Layout', 'database', 200, 'resource', 1),
(5, '资源目录', '/resource/catalog', 'resource/catalog/index', 'folder', 1, 'resource:catalog', 1),
(5, '数据视图', '/resource/view', 'resource/view/index', 'eye', 2, 'resource:view', 1);

-- 查询验证
-- SELECT * FROM sys_menu WHERE deleted_at IS NULL ORDER BY sort_order ASC, id ASC;
