package data_catalog_category

import (
	"errors"
)

// ========== 错误定义 ==========

var (
	// ErrNotFound 分类关联不存在
	ErrNotFound = errors.New("data catalog category not found")
)
