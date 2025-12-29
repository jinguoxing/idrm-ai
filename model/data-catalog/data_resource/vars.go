package data_resource

import (
	"errors"

	"idrm/pkg/errorx"
)

// ========== 常量定义 ==========

// 资源来源常量
const (
	SourceManual = "manual" // 手动挂接
	SourceSync   = "sync"   // 同步挂接
)

// ========== 错误定义 ==========

var (
	// ErrNotFound 资源不存在
	ErrNotFound = errors.New("data resource not found")

	// ErrDuplicateLogicalView 逻辑视图重复（一个目录只能有一个逻辑视图）
	ErrDuplicateLogicalView = errors.New("duplicate logical view, only one allowed")

	// ErrInvalidResourceType 无效的资源类型
	ErrInvalidResourceType = errors.New("invalid resource type")

	// ErrResourceTypeNotSupported 资源类型不支持
	ErrResourceTypeNotSupported = errors.New("resource type not supported")
)

// ========== 错误构造函数 ==========

// NewResourceNotFoundError 创建资源不存在错误
func NewResourceNotFoundError(id string) error {
	return errorx.NewWithMsg(errorx.ErrCodeDataResourceNotExist, "数据资源已挂载或不存在")
}

// NewResourceTypeNotSupportedError 创建资源类型不支持错误
func NewResourceTypeNotSupportedError(resourceType string) error {
	return errorx.NewWithMsg(errorx.ErrCodeDataResourceTypeNotSupport, "数据资源类型不支持")
}

// NewDuplicateLogicalViewError 创建逻辑视图重复错误
func NewDuplicateLogicalViewError() error {
	return errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "一个目录只能挂接一个逻辑视图")
}
