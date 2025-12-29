package data_catalog_column

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
	db   sqlxDBer
	base sqlxTxer
}

// Insert 批量插入信息项记录
func (m *sqlxModel) Insert(ctx context.Context, data []*DataCatalogColumn) error {
	if len(data) == 0 {
		return nil
	}

	query := `
		INSERT INTO t_data_catalog_column (
			id, catalog_id, column_name, column_type, column_length, column_scale,
			column_code, column_alias, column_describe, data_type, data_format,
			data_unit, value_range, is_primary_key, is_sensitive, is_required,
			is_unique, dictionary_id, standard_code_id, sort_order, created_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	for _, column := range data {
		_, err := m.db.ExecContext(ctx, query,
			column.Id, column.CatalogId, column.ColumnName, column.ColumnType,
			column.ColumnLength, column.ColumnScale, column.ColumnCode,
			column.ColumnAlias, column.ColumnDescribe, column.DataType,
			column.DataFormat, column.DataUnit, column.ValueRange,
			column.IsPrimaryKey, column.IsSensitive, column.IsRequired,
			column.IsUnique, column.DictionaryId, column.StandardCodeId,
			column.SortOrder, column.CreatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to insert data catalog column: %w", err)
		}
	}

	return nil
}

// FindOne 根据ID查询单个信息项
func (m *sqlxModel) FindOne(ctx context.Context, id string) (*DataCatalogColumn, error) {
	var column DataCatalogColumn
	query := `
		SELECT * FROM t_data_catalog_column
		WHERE id = ? LIMIT 1
	`

	err := m.db.GetContext(ctx, &column, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find data catalog column by id %s: %w", id, err)
	}

	return &column, nil
}

// Update 更新信息项记录
func (m *sqlxModel) Update(ctx context.Context, data *DataCatalogColumn) error {
	query := `
		UPDATE t_data_catalog_column SET
			column_name = ?,
			column_type = ?,
			column_length = ?,
			column_scale = ?,
			column_code = ?,
			column_alias = ?,
			column_describe = ?,
			data_type = ?,
			data_format = ?,
			data_unit = ?,
			value_range = ?,
			is_primary_key = ?,
			is_sensitive = ?,
			is_required = ?,
			is_unique = ?,
			dictionary_id = ?,
			standard_code_id = ?,
			sort_order = ?,
			updated_at = ?
		WHERE id = ?
	`

	_, err := m.db.ExecContext(ctx, query,
		data.ColumnName, data.ColumnType, data.ColumnLength, data.ColumnScale,
		data.ColumnCode, data.ColumnAlias, data.ColumnDescribe,
		data.DataType, data.DataFormat, data.DataUnit, data.ValueRange,
		data.IsPrimaryKey, data.IsSensitive, data.IsRequired,
		data.IsUnique, data.DictionaryId, data.StandardCodeId,
		data.SortOrder, data.UpdatedAt, data.Id,
	)

	if err != nil {
		return fmt.Errorf("failed to update data catalog column %s: %w", data.Id, err)
	}

	return nil
}

// Delete 删除信息项记录
func (m *sqlxModel) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM t_data_catalog_column WHERE id = ?`

	result, err := m.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete data catalog column %s: %w", id, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// FindByCatalogId 根据目录ID查询信息项列表
func (m *sqlxModel) FindByCatalogId(ctx context.Context, catalogId string) ([]*DataCatalogColumn, error) {
	var columns []*DataCatalogColumn
	query := `
		SELECT * FROM t_data_catalog_column
		WHERE catalog_id = ?
		ORDER BY sort_order ASC, id ASC
	`

	err := m.db.SelectContext(ctx, &columns, query, catalogId)
	if err != nil {
		return nil, fmt.Errorf("failed to find columns by catalog id %s: %w", catalogId, err)
	}

	return columns, nil
}

// DeleteByCatalogId 根据目录ID删除所有信息项
func (m *sqlxModel) DeleteByCatalogId(ctx context.Context, catalogId string) error {
	query := `DELETE FROM t_data_catalog_column WHERE catalog_id = ?`

	_, err := m.db.ExecContext(ctx, query, catalogId)
	if err != nil {
		return fmt.Errorf("failed to delete columns by catalog id %s: %w", catalogId, err)
	}

	return nil
}

// BatchInsert 批量插入信息项
func (m *sqlxModel) BatchInsert(ctx context.Context, columns []*DataCatalogColumn) error {
	// SQLx 不像 GORM 有内置的批量插入优化，使用普通的 Insert 方法
	return m.Insert(ctx, columns)
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
		base: m.base,
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
