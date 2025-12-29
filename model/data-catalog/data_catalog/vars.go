package data_catalog

import (
	"errors"

	"idrm/pkg/errorx"
)

// ========== 共享类型常量 ==========

const (
	SharedTypeUnconditional = 1 // 无条件共享
	SharedTypeConditional   = 2 // 有条件共享
	SharedTypeNoShared      = 3 // 不予共享
)

// ========== 开放类型常量 ==========

const (
	OpenTypePublic    = 1 // 向公众开放
	OpenTypeNotPublic = 2 // 不向公众开放
)

// ========== 数据归集机制常量 ==========

const (
	SyncMechanismIncremental = 1 // 增量
	SyncMechanismFull        = 2 // 全量
)

// ========== 错误定义 ==========

var (
	// ErrNotFound 目录不存在
	ErrNotFound = errors.New("data catalog not found")

	// ErrDuplicateName 目录名称重复
	ErrDuplicateName = errors.New("catalog name already exists")

	// ErrInvalidStatus 无效的状态
	ErrInvalidStatus = errors.New("invalid status")

	// ErrDraftAlreadyExists 草稿已存在
	ErrDraftAlreadyExists = errors.New("draft already exists")

	// ErrCannotPublishDraft 草稿不能发布
	ErrCannotPublishDraft = errors.New("cannot publish draft")

	// ErrCatalogPublished 目录已发布，不能修改
	ErrCatalogPublished = errors.New("catalog is published, cannot modify")
)

// ========== 错误构造函数 ==========

// NewCatalogNotFoundError 创建目录不存在错误
func NewCatalogNotFoundError(id string) error {
	return errorx.NewWithMsg(errorx.ErrCodeDataCatalogNotFound, "数据资源目录不存在")
}

// NewDepartmentNotFoundError 创建部门不存在错误
func NewDepartmentNotFoundError(deptId string) error {
	return errorx.NewWithMsg(errorx.ErrCodeDataCatalogDepartmentNotFound, "数据资源目录关联部门不存在")
}

// NewDuplicateNameError 创建名称重复错误
func NewDuplicateNameError(title string) error {
	return errorx.NewWithMsg(errorx.ErrCodeCatalogNameRepeat, "数据目录名称重复")
}

// NewInvalidParameterError 创建参数错误
func NewInvalidParameterError(param string) error {
	return errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "参数"+param+"不合法")
}
