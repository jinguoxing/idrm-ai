package data_catalog

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// sqlxDBer 定义 sqlx.DB 和 sqlx.Tx 的通用接口
type sqlxDBer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// sqlxTxer 定义事务接口
type sqlxTxer interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

func init() {
	// 注册SQLx工厂函数
	RegisterSqlxFactory(func(db *sql.DB) Model {
		return &sqlxModel{
			db:   sqlx.NewDb(db, "mysql"),
			base: sqlx.NewDb(db, "mysql"),
		}
	})
}

// sqlxModel SQLx实现的数据访问对象
type sqlxModel struct {
	db   sqlxDBer      // 当前执行器（可能是DB或Tx）
	base sqlxTxer      // 基础DB（用于开启事务）
}

// Insert 插入新的目录记录
func (m *sqlxModel) Insert(ctx context.Context, data *DataCatalog) (*DataCatalog, error) {
	query := `
		INSERT INTO t_data_catalog (
			id, code, title, description, source_department_id,
			data_domain, shared_type, open_type, publish_status,
			online_status, created_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	_, err := m.db.ExecContext(ctx, query,
		data.Id, data.Code, data.Title, data.Description, data.SourceDeptId,
		data.DataDomain, data.SharedType, data.OpenType, data.PublishStatus,
		data.OnlineStatus, data.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert data catalog: %w", err)
	}

	return data, nil
}

// FindOne 根据ID查询单个目录
func (m *sqlxModel) FindOne(ctx context.Context, id string) (*DataCatalog, error) {
	var catalog DataCatalog
	query := `
		SELECT * FROM t_data_catalog
		WHERE id = ? LIMIT 1
	`

	err := m.db.GetContext(ctx, &catalog, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find data catalog by id %s: %w", id, err)
	}

	return &catalog, nil
}

// Update 更新目录记录
func (m *sqlxModel) Update(ctx context.Context, data *DataCatalog) error {
	query := `
		UPDATE t_data_catalog SET
			code = ?,
			title = ?,
			description = ?,
			data_domain = ?,
			shared_type = ?,
			open_type = ?,
			publish_status = ?,
			online_status = ?,
			updated_at = ?
		WHERE id = ?
	`

	result, err := m.db.ExecContext(ctx, query,
		data.Code, data.Title, data.Description, data.DataDomain,
		data.SharedType, data.OpenType, data.PublishStatus, data.OnlineStatus,
		data.UpdatedAt, data.Id,
	)

	if err != nil {
		return fmt.Errorf("failed to update data catalog %s: %w", data.Id, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete 删除目录记录
func (m *sqlxModel) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM t_data_catalog WHERE id = ?`

	result, err := m.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete data catalog %s: %w", id, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// FindBySourceDept 根据来源部门查询目录列表
func (m *sqlxModel) FindBySourceDept(ctx context.Context, deptId string, status string) ([]*DataCatalog, error) {
	var catalogs []*DataCatalog

	query := `
		SELECT * FROM t_data_catalog
		WHERE source_department_id = ?
	`

	args := []interface{}{deptId}

	if status != "" {
		query += " AND publish_status = ?"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"

	err := m.db.SelectContext(ctx, &catalogs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find catalogs by dept %s: %w", deptId, err)
	}

	return catalogs, nil
}

// FindByDraftId 根据草稿ID查询目录
func (m *sqlxModel) FindByDraftId(ctx context.Context, draftId string) (*DataCatalog, error) {
	var catalog DataCatalog
	query := `
		SELECT * FROM t_data_catalog
		WHERE id = ? LIMIT 1
	`

	err := m.db.GetContext(ctx, &catalog, query, draftId)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find catalog by draft_id %s: %w", draftId, err)
	}

	return &catalog, nil
}

// FindDrafts 查询草稿列表
func (m *sqlxModel) FindDrafts(ctx context.Context, deptId string, status string, page, pageSize int) ([]*DataCatalog, int64, error) {
	var catalogs []*DataCatalog
	var total int64

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if deptId != "" {
		whereClause += " AND source_department_id = ?"
		args = append(args, deptId)
	}

	if status != "" {
		whereClause += " AND publish_status = ?"
		args = append(args, status)
	}

	// 统计总数
	countQuery := "SELECT COUNT(*) FROM t_data_catalog " + whereClause
	if err := m.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count drafts: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	query := `
		SELECT * FROM t_data_catalog
		` + whereClause + `
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	args = append(args, pageSize, offset)

	if err := m.db.SelectContext(ctx, &catalogs, query, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to find drafts: %w", err)
	}

	return catalogs, total, nil
}

// CheckNameExists 检查目录名称是否重复
func (m *sqlxModel) CheckNameExists(ctx context.Context, deptId string, title string, excludeId string) (bool, error) {
	var count int64

	query := `
		SELECT COUNT(*) FROM t_data_catalog
		WHERE source_department_id = ? AND title = ?
	`

	args := []interface{}{deptId, title}

	if excludeId != "" {
		query += " AND id != ?"
		args = append(args, excludeId)
	}

	if err := m.db.GetContext(ctx, &count, query, args...); err != nil {
		return false, fmt.Errorf("failed to check name exists: %w", err)
	}

	return count > 0, nil
}

// ========== 事务方法实现 ==========

// WithTx 创建事务副本
func (m *sqlxModel) WithTx(tx interface{}) Model {
	sqlxTx, ok := tx.(*sqlx.Tx)
	if !ok {
		panic("invalid transaction type, expected *sqlx.Tx")
	}

	return &sqlxModel{
		db:   sqlxTx,
		base: m.base, // 保留基础DB引用
	}
}

// Trans 在事务中执行多个操作
func (m *sqlxModel) Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error {
	tx, err := m.base.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	modelTx := &sqlxModel{
		db:   tx,
		base: m.base,
	}

	if err := fn(ctx, modelTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback error: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
