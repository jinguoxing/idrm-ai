package menu

import (
	"errors"

	"idrm/pkg/errorx"
)

// Model层专用错误定义
var (
	// ErrNotFound 菜单不存在
	ErrNotFound = errors.New("menu not found")

	// ErrInvalidParent 无效的父级菜单
	ErrInvalidParent = errors.New("invalid parent menu")

	// ErrHasChildren 存在子菜单
	ErrHasChildren = errors.New("menu has children")

	// ErrParentLoopDetected 检测到循环父级引用
	ErrParentLoopDetected = errors.New("parent loop detected")
)

// NewParentMenuNotFoundError 创建父级菜单不存在错误
func NewParentMenuNotFoundError(parentId int64) error {
	return errorx.NewWithMsg(errorx.ErrCodeParentMenuNotFound, "父级菜单不存在")
}

// NewMenuHasSubEntriesError 创建存在子菜单错误
func NewMenuHasSubEntriesError(menuId int64) error {
	return errorx.NewWithMsg(errorx.ErrCodeMenuHasSubEntries, "存在子菜单无法直接删除")
}
