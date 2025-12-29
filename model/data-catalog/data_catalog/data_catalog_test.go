package data_catalog

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestGormModel_Insert 测试 GORM 插入功能
func TestGormModel_Insert(t *testing.T) {
	// 创建测试数据库连接
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      createMockDB(t),
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	model := &gormModel{db: db}
	ctx := context.Background()

	data := &DataCatalog{
		Id:         "1234567890123456789",
		Code:       "TEST001",
		Title:      "测试目录",
		SourceDeptId: "dept001",
		DataKind:   1,
		SharedType: 1,
		OpenType:   1,
		SharedMode: 1,
		PublishStatus: StatusUnpublished,
		OnlineStatus:  OnlineStatusNotline,
		CreatedAt:  time.Now(),
	}

	// 注意: 这个测试需要真实的数据库或者更复杂的 mock
	// 这里仅作为结构示例
	_, err = model.Insert(ctx, data)
	// 实际测试需要 mock 数据库返回
	_ = err
}

// TestGormModel_FindOne 测试 GORM 查询功能
func TestGormModel_FindOne(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	model := &gormModel{db: gormDB}
	ctx := context.Background()

	_, err = model.FindOne(ctx, "1234567890123456789")
	// 实际测试需要设置 mock 期望
	_ = err
}

// TestSqlxModel_Insert 测试 SQLx 插入功能
func TestSqlxModel_Insert(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "mysql")
	model := &sqlxModel{db: sqlxDB}
	ctx := context.Background()

	data := &DataCatalog{
		Id:         "1234567890123456789",
		Code:       "TEST001",
		Title:      "测试目录",
		SourceDeptId: "dept001",
		DataKind:   1,
		SharedType: 1,
		OpenType:   1,
		SharedMode: 1,
		PublishStatus: StatusUnpublished,
		OnlineStatus:  OnlineStatusNotline,
		CreatedAt:  time.Now(),
	}

	_, err = model.Insert(ctx, data)
	// 实际测试需要设置 mock 期望
	_ = err
}

// TestNewModel 测试工厂函数
func TestNewModel(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlDB := db.(*sql.DB)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// 测试 GORM 优先
	RegisterGormFactory(func(db *gorm.DB) Model {
		return &gormModel{db: db}
	})
	RegisterSqlxFactory(func(db *sql.DB) Model {
		return &sqlxModel{db: sqlx.NewDb(db, "mysql")}
	})

	model := NewModel(sqlDB, gormDB)
	assert.NotNil(t, model)

	// 测试仅 SQLx
	model2 := NewModel(sqlDB, nil)
	assert.NotNil(t, model2)

	// 测试 panic
	assert.Panics(t, func() {
		NewModel(nil, nil)
	})
}

// createMockDB 创建 mock 数据库连接
func createMockDB(t *testing.T) *sql.DB {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	return db
}

// TestDataCatalog_TableName 测试表名
func TestDataCatalog_TableName(t *testing.T) {
	catalog := DataCatalog{}
	assert.Equal(t, "t_data_catalog", catalog.TableName())
}

// TestStatusConstants 测试状态常量
func TestStatusConstants(t *testing.T) {
	tests := []struct {
		name   string
		status string
	}{
		{"Unpublished", StatusUnpublished},
		{"PubAuditing", StatusPubAuditing},
		{"Published", StatusPublished},
		{"PubReject", StatusPubReject},
		{"ChangeAuditing", StatusChangeAuditing},
		{"ChangeReject", StatusChangeReject},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.status)
		})
	}
}

// TestOnlineStatusConstants 测试上线状态常量
func TestOnlineStatusConstants(t *testing.T) {
	tests := []struct {
		name   string
		status string
	}{
		{"Notline", OnlineStatusNotline},
		{"Online", OnlineStatusOnline},
		{"Offline", OnlineStatusOffline},
		{"UpAuditing", OnlineStatusUpAuditing},
		{"DownAuditing", OnlineStatusDownAuditing},
		{"UpReject", OnlineStatusUpReject},
		{"DownReject", OnlineStatusDownReject},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.status)
		})
	}
}

// TestErrorConstructors 测试错误构造函数
func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name    string
		fn      func() error
		wantErr bool
	}{
		{
			name: "NotFound",
			fn: func() error {
				return NewCatalogNotFoundError("123")
			},
			wantErr: true,
		},
		{
			name: "DepartmentNotFound",
			fn: func() error {
				return NewDepartmentNotFoundError("dept001")
			},
			wantErr: true,
		},
		{
			name: "DuplicateName",
			fn: func() error {
				return NewDuplicateNameError("测试目录")
			},
			wantErr: true,
		},
		{
			name: "InvalidParameter",
			fn: func() error {
				return NewInvalidParameterError("code")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

// BenchmarkGormModel_Insert GORM 插入性能测试
func BenchmarkGormModel_Insert(b *testing.B) {
	// 需要真实数据库连接才能进行准确的性能测试
	// 这里仅作为结构示例
	b.Skip("需要真实数据库连接")
}

// BenchmarkSqlxModel_Insert SQLx 插入性能测试
func BenchmarkSqlxModel_Insert(b *testing.B) {
	// 需要真实数据库连接才能进行准确的性能测试
	// 这里仅作为结构示例
	b.Skip("需要真实数据库连接")
}
