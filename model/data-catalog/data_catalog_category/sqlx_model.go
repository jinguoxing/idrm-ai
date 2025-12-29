package data_catalog_category

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
	Rebind(query string) string
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

func (m *sqlxModel) Insert(ctx context.Context, data []*DataCatalogCategory) error {
	if len(data) == 0 {
		return nil
	}

	query := `
		INSERT INTO t_data_catalog_category (id, catalog_id, category_id, created_at)
		VALUES (?, ?, ?, ?)
	`

	for _, category := range data {
		_, err := m.db.ExecContext(ctx, query,
			category.Id, category.CatalogId, category.CategoryId, category.CreatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to insert data catalog category: %w", err)
		}
	}

	return nil
}

func (m *sqlxModel) FindByCatalogId(ctx context.Context, catalogId string) ([]*DataCatalogCategory, error) {
	var categories []*DataCatalogCategory
	query := `
		SELECT * FROM t_data_catalog_category
		WHERE catalog_id = ?
	`

	err := m.db.SelectContext(ctx, &categories, query, catalogId)
	if err != nil {
		return nil, fmt.Errorf("failed to find categories by catalog id %s: %w", catalogId, err)
	}

	return categories, nil
}

func (m *sqlxModel) DeleteByCatalogId(ctx context.Context, catalogId string) error {
	query := `DELETE FROM t_data_catalog_category WHERE catalog_id = ?`

	_, err := m.db.ExecContext(ctx, query, catalogId)
	if err != nil {
		return fmt.Errorf("failed to delete categories by catalog id %s: %w", catalogId, err)
	}

	return nil
}

func (m *sqlxModel) BatchDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	query, args, err := sqlx.In(`DELETE FROM t_data_catalog_category WHERE id IN (?)`, ids)
	if err != nil {
		return fmt.Errorf("failed to build batch delete query: %w", err)
	}

	query = m.db.Rebind(query)
	_, err = m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to batch delete categories: %w", err)
	}

	return nil
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
