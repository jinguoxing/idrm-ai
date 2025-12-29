package data_catalog_column

import (
	"errors"

	"idrm/pkg/errorx"
)

// ========== 常量定义 ==========

// 字段类型常量
const (
	ColumnTypeString   = "string"   // 字符串
	ColumnTypeInteger  = "integer"  // 整数
	ColumnTypeDecimal  = "decimal"  // 小数
	ColumnTypeDate     = "date"     // 日期
	ColumnTypeDateTime = "datetime" // 日期时间
	ColumnTypeBoolean  = "boolean"  // 布尔
	ColumnTypeText     = "text"     // 文本
	ColumnTypeBlob     = "blob"     // 二进制
)

// ========== 错误定义 ==========

var (
	// ErrNotFound 信息项不存在
	ErrNotFound = errors.New("data catalog column not found")

	// ErrInvalidColumnType 无效的字段类型
	ErrInvalidColumnType = errors.New("invalid column type")

	// ErrMissingColumnLength 缺少字段长度
	ErrMissingColumnLength = errors.New("column length is required for this type")

	// ErrMissingColumnScale 缺少字段精度
	ErrMissingColumnScale = errors.New("column scale is required for this type")
)

// ========== 错误构造函数 ==========

// NewColumnNotFoundError 创建信息项不存在错误
func NewColumnNotFoundError(id string) error {
	return errorx.NewWithMsg(errorx.ErrCodeDataCatalogColumnNotFound, "信息项不存在")
}

// NewInvalidColumnTypeError 创建无效字段类型错误
func NewInvalidColumnTypeError(columnType string) error {
	return errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "字段类型"+columnType+"不合法")
}

// NewMissingColumnLengthError 创建缺少字段长度错误
func NewMissingColumnLengthError() error {
	return errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "字段长度不能为空")
}
