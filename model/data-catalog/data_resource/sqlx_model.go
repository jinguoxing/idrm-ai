package data_resource

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
	RegisterSqlxFactory(func(db *sql.DB) Model {
		return &sqlxModel{
			db:   sqlx.NewDb(db, "mysql"),
			base: sqlx.NewDb(db, "mysql"),
		}
	})
}

type sqlxModel struct {
	db   sqlxDBer
	base sqlxTxer
}

func (m *sqlxModel) Insert(ctx context.Context, data *DataResource) error {
	query := `
		INSERT INTO t_data_resource (
			id, catalog_id, resource_type, resource_name, resource_code, resource_desc,
			database_id, table_name, table_comment, api_id, api_url, api_method, api_protocol,
			file_id, file_name, file_url, file_size, source, sort_order, created_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	_, err := m.db.ExecContext(ctx, query,
		data.Id, data.CatalogId, data.ResourceType, data.ResourceName, data.ResourceCode,
		data.ResourceDesc, data.DatabaseId, data.TblName, data.TableComment, data.ApiId,
		data.ApiUrl, data.ApiMethod, data.ApiProtocol, data.FileId, data.FileName,
		data.FileUrl, data.FileSize, data.Source, data.SortOrder, data.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert data resource: %w", err)
	}

	return nil
}

func (m *sqlxModel) FindOne(ctx context.Context, id string) (*DataResource, error) {
	var resource DataResource
	query := `SELECT * FROM t_data_resource WHERE id = ? LIMIT 1`

	err := m.db.GetContext(ctx, &resource, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find data resource by id %s: %w", id, err)
	}

	return &resource, nil
}

func (m *sqlxModel) Update(ctx context.Context, data *DataResource) error {
	query := `
		UPDATE t_data_resource SET
			resource_type = ?, resource_name = ?, resource_code = ?, resource_desc = ?,
			database_id = ?, table_name = ?, table_comment = ?, api_id = ?, api_url = ?,
			api_method = ?, api_protocol = ?, file_id = ?, file_name = ?, file_url = ?,
			file_size = ?, source = ?, sort_order = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := m.db.ExecContext(ctx, query,
		data.ResourceType, data.ResourceName, data.ResourceCode, data.ResourceDesc,
		data.DatabaseId, data.TblName, data.TableComment, data.ApiId, data.ApiUrl,
		data.ApiMethod, data.ApiProtocol, data.FileId, data.FileName, data.FileUrl,
		data.FileSize, data.Source, data.SortOrder, data.UpdatedAt, data.Id,
	)

	if err != nil {
		return fmt.Errorf("failed to update data resource %s: %w", data.Id, err)
	}

	return nil
}

func (m *sqlxModel) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM t_data_resource WHERE id = ?`

	result, err := m.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete data resource %s: %w", id, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (m *sqlxModel) FindByCatalogId(ctx context.Context, catalogId string) ([]*DataResource, error) {
	var resources []*DataResource
	query := `
		SELECT * FROM t_data_resource
		WHERE catalog_id = ?
		ORDER BY sort_order ASC, id ASC
	`

	err := m.db.SelectContext(ctx, &resources, query, catalogId)
	if err != nil {
		return nil, fmt.Errorf("failed to find resources by catalog id %s: %w", catalogId, err)
	}

	return resources, nil
}

func (m *sqlxModel) FindByCatalogIdAndType(ctx context.Context, catalogId string, resourceType ResourceType) ([]*DataResource, error) {
	var resources []*DataResource
	query := `
		SELECT * FROM t_data_resource
		WHERE catalog_id = ? AND resource_type = ?
		ORDER BY sort_order ASC, id ASC
	`

	err := m.db.SelectContext(ctx, &resources, query, catalogId, resourceType)
	if err != nil {
		return nil, fmt.Errorf("failed to find resources by catalog id %s and type %s: %w", catalogId, resourceType, err)
	}

	return resources, nil
}

func (m *sqlxModel) DeleteByCatalogId(ctx context.Context, catalogId string) error {
	query := `DELETE FROM t_data_resource WHERE catalog_id = ?`

	_, err := m.db.ExecContext(ctx, query, catalogId)
	if err != nil {
		return fmt.Errorf("failed to delete resources by catalog id %s: %w", catalogId, err)
	}

	return nil
}

func (m *sqlxModel) CountByCatalogId(ctx context.Context, catalogId string) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM t_data_resource WHERE catalog_id = ?`

	err := m.db.GetContext(ctx, &count, query, catalogId)
	if err != nil {
		return 0, fmt.Errorf("failed to count resources by catalog id %s: %w", catalogId, err)
	}

	return count, nil
}

func (m *sqlxModel) CountLogicalViews(ctx context.Context, catalogId string) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM t_data_resource WHERE catalog_id = ? AND resource_type = ?`

	err := m.db.GetContext(ctx, &count, query, catalogId, ResourceTypeLogicalView)
	if err != nil {
		return 0, fmt.Errorf("failed to count logical views by catalog id %s: %w", catalogId, err)
	}

	return count, nil
}

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
