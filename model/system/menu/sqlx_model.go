package menu

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func init() {
	// 注册SQLx工厂函数
	RegisterSqlxFactory(func(db *sql.DB) Model {
		return &MenuModel{
			db:    db,
			sqlx:  sqlx.NewDb(db, "mysql"),
			table: "sys_menu",
		}
	})
}

// MenuModel SQLx实现的菜单数据访问对象
type MenuModel struct {
	db    *sql.DB
	sqlx  *sqlx.DB
	tx    *sqlx.Tx
	table string
}

// Insert 插入新的菜单记录
func (m *MenuModel) Insert(ctx context.Context, data *Menu) (*Menu, error) {
	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now

	query := `
		INSERT INTO ` + m.table + `
		(parent_id, name, route_path, component_path, icon, sort_order, perm_tag, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var result sql.Result
	var err error

	if m.tx != nil {
		result, err = m.tx.ExecContext(ctx, query,
			data.ParentId, data.Name, data.RoutePath, data.ComponentPath,
			data.Icon, data.SortOrder, data.PermTag, data.Status,
			data.CreatedAt, data.UpdatedAt,
		)
	} else {
		result, err = m.db.ExecContext(ctx, query,
			data.ParentId, data.Name, data.RoutePath, data.ComponentPath,
			data.Icon, data.SortOrder, data.PermTag, data.Status,
			data.CreatedAt, data.UpdatedAt,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to insert menu: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	data.Id = id
	return data, nil
}

// FindOne 根据ID查询单个菜单
func (m *MenuModel) FindOne(ctx context.Context, id int64) (*Menu, error) {
	var menu Menu
	query := `
		SELECT id, parent_id, name, route_path, component_path, icon,
		       sort_order, perm_tag, status, created_at, updated_at
		FROM ` + m.table + `
		WHERE id = ? AND deleted_at IS NULL
	`

	err := m.sqlx.GetContext(ctx, &menu, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find menu by id %d: %w", id, err)
	}

	return &menu, nil
}

// FindByPermTag 根据权限标识查询菜单
func (m *MenuModel) FindByPermTag(ctx context.Context, permTag string) (*Menu, error) {
	var menu Menu
	query := `
		SELECT id, parent_id, name, route_path, component_path, icon,
		       sort_order, perm_tag, status, created_at, updated_at
		FROM ` + m.table + `
		WHERE perm_tag = ? AND deleted_at IS NULL
	`

	err := m.sqlx.GetContext(ctx, &menu, query, permTag)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find menu by perm_tag %s: %w", permTag, err)
	}

	return &menu, nil
}

// Update 更新菜单记录
func (m *MenuModel) Update(ctx context.Context, data *Menu) error {
	data.UpdatedAt = time.Now()

	query := `
		UPDATE ` + m.table + `
		SET parent_id = ?, name = ?, route_path = ?, component_path = ?,
		    icon = ?, sort_order = ?, perm_tag = ?, status = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := m.db.ExecContext(ctx, query,
		data.ParentId, data.Name, data.RoutePath, data.ComponentPath,
		data.Icon, data.SortOrder, data.PermTag, data.Status,
		data.UpdatedAt, data.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to update menu %d: %w", data.Id, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete 删除菜单记录（软删除）
func (m *MenuModel) Delete(ctx context.Context, id int64) error {
	query := `UPDATE ` + m.table + ` SET deleted_at = ? WHERE id = ?`

	result, err := m.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete menu %d: %w", id, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// FindAll 查询所有菜单（支持状态过滤）
func (m *MenuModel) FindAll(ctx context.Context, status int) ([]*Menu, error) {
	var menus []*Menu
	query := `
		SELECT id, parent_id, name, route_path, component_path, icon,
		       sort_order, perm_tag, status, created_at, updated_at
		FROM ` + m.table + `
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	if status != -1 {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY sort_order ASC, id ASC"

	err := m.sqlx.SelectContext(ctx, &menus, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find all menus: %w", err)
	}

	return menus, nil
}

// FindByParentId 根据父级ID查询子菜单列表
func (m *MenuModel) FindByParentId(ctx context.Context, parentId int64) ([]*Menu, error) {
	var menus []*Menu
	query := `
		SELECT id, parent_id, name, route_path, component_path, icon,
		       sort_order, perm_tag, status, created_at, updated_at
		FROM ` + m.table + `
		WHERE parent_id = ? AND deleted_at IS NULL
		ORDER BY sort_order ASC, id ASC
	`

	err := m.sqlx.SelectContext(ctx, &menus, query, parentId)
	if err != nil {
		return nil, fmt.Errorf("failed to find menus by parent_id %d: %w", parentId, err)
	}

	return menus, nil
}

// CountByParentId 统计指定父级下的子菜单数量
func (m *MenuModel) CountByParentId(ctx context.Context, parentId int64) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM ` + m.table + ` WHERE parent_id = ? AND deleted_at IS NULL`

	err := m.db.QueryRowContext(ctx, query, parentId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count menus by parent_id %d: %w", parentId, err)
	}

	return count, nil
}

// HasChildren 检查菜单是否有子菜单
func (m *MenuModel) HasChildren(ctx context.Context, parentId int64) (bool, error) {
	count, err := m.CountByParentId(ctx, parentId)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// WithTx 创建事务副本
func (m *MenuModel) WithTx(tx interface{}) Model {
	sqlxTx, ok := tx.(*sqlx.Tx)
	if !ok {
		panic(fmt.Sprintf("invalid transaction type: expected *sqlx.Tx, got %T", tx))
	}

	return &MenuModel{
		db:    nil, // 事务中不需要原始db
		sqlx:  nil,
		table: m.table,
		tx:    sqlxTx,
	}
}

// Trans 在事务中执行多个操作
func (m *MenuModel) Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error {
	tx, err := m.sqlx.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txModel := &MenuModel{
		db:    nil,
		sqlx:  nil,
		table: m.table,
		tx:    tx,
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(ctx, txModel); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("original error: %w, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// 辅助方法：在事务中执行查询
func (m *MenuModel) queryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	if m.tx != nil {
		return m.tx.QueryxContext(ctx, query, args...)
	}
	return m.sqlx.QueryxContext(ctx, query, args...)
}
