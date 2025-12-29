package orm

import (
	"database/sql"
	"errors"
	"strings"

	"idrm/pkg/errorx"

	"gorm.io/gorm"
)

// ORMErrorAdapter ORM错误适配器
// 将不同ORM的错误转换为统一的业务错误
type ORMErrorAdapter struct{}

// NewORMErrorAdapter 创建错误适配器实例
func NewORMErrorAdapter() *ORMErrorAdapter {
	return &ORMErrorAdapter{}
}

// AdaptError 适配错误
// 根据错误类型自动选择适配方法
func (a *ORMErrorAdapter) AdaptError(err error, ormType ORMType) error {
	if err == nil {
		return nil
	}

	switch ormType {
	case ORMTypeGORM:
		return a.AdaptGORMError(err)
	case ORMTypeSQLx:
		return a.AdaptSQLxError(err)
	default:
		return err
	}
}

// AdaptGORMError 适配GORM错误
func (a *ORMErrorAdapter) AdaptGORMError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return errorx.NewWithCode(errorx.ErrCodeNotFound)

	case errors.Is(err, gorm.ErrInvalidTransaction):
		return errorx.NewWithMsg(errorx.ErrCodeDatabase, "无效的事务操作")

	case errors.Is(err, gorm.ErrNotImplemented):
		return errorx.NewWithMsg(errorx.ErrCodeSystem, "功能未实现")

	case errors.Is(err, gorm.ErrMissingWhereClause):
		return errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "缺少WHERE条件")

	case errors.Is(err, gorm.ErrUnsupportedRelation):
		return errorx.NewWithMsg(errorx.ErrCodeSystem, "不支持的关联关系")

	case errors.Is(err, gorm.ErrPrimaryKeyRequired):
		return errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "缺少主键")

	case errors.Is(err, gorm.ErrModelValueRequired):
		return errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "缺少模型值")

	case errors.Is(err, gorm.ErrInvalidData):
		return errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "无效的数据")

	default:
		// 检查是否是数据库层面的错误
		return a.adaptDatabaseError(err)
	}
}

// AdaptSQLxError 适配SQLx错误
func (a *ORMErrorAdapter) AdaptSQLxError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errorx.NewWithCode(errorx.ErrCodeNotFound)

	case errors.Is(err, sql.ErrTxDone):
		return errorx.NewWithMsg(errorx.ErrCodeDatabase, "事务已完成")

	case errors.Is(err, sql.ErrConnDone):
		return errorx.NewWithMsg(errorx.ErrCodeDatabase, "数据库连接已关闭")

	default:
		// 检查是否是数据库层面的错误
		return a.adaptDatabaseError(err)
	}
}

// adaptDatabaseError 适配数据库层面的错误
func (a *ORMErrorAdapter) adaptDatabaseError(err error) error {
	errMsg := strings.ToLower(err.Error())

	switch {
	// MySQL错误
	case strings.Contains(errMsg, "duplicate entry"):
		return errorx.NewWithCode(errorx.ErrCodeAlreadyExists)

	case strings.Contains(errMsg, "foreign key constraint"):
		return errorx.NewWithMsg(errorx.ErrCodeOperationFailed, "外键约束违反")

	case strings.Contains(errMsg, "data too long"):
		return errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "数据长度超出限制")

	case strings.Contains(errMsg, "cannot be null"):
		return errorx.NewWithMsg(errorx.ErrCodeParamMissing, "必填字段不能为空")

	case strings.Contains(errMsg, "connection refused"):
		return errorx.NewWithMsg(errorx.ErrCodeDatabase, "数据库连接被拒绝")

	case strings.Contains(errMsg, "timeout"):
		return errorx.NewWithMsg(errorx.ErrCodeDatabase, "数据库操作超时")

	case strings.Contains(errMsg, "deadlock"):
		return errorx.NewWithMsg(errorx.ErrCodeDatabase, "数据库死锁")

	case strings.Contains(errMsg, "lock wait timeout"):
		return errorx.NewWithMsg(errorx.ErrCodeDatabase, "锁等待超时")

	// SQLite错误
	case strings.Contains(errMsg, "unique constraint"):
		return errorx.NewWithCode(errorx.ErrCodeAlreadyExists)

	case strings.Contains(errMsg, "not null constraint"):
		return errorx.NewWithMsg(errorx.ErrCodeParamMissing, "必填字段不能为空")

	case strings.Contains(errMsg, "database is locked"):
		return errorx.NewWithMsg(errorx.ErrCodeDatabase, "数据库被锁定")

	// 通用数据库错误
	case strings.Contains(errMsg, "syntax error"):
		return errorx.NewWithMsg(errorx.ErrCodeSystem, "SQL语法错误")

	case strings.Contains(errMsg, "table doesn't exist") || strings.Contains(errMsg, "no such table"):
		return errorx.NewWithMsg(errorx.ErrCodeSystem, "数据表不存在")

	case strings.Contains(errMsg, "column") && strings.Contains(errMsg, "doesn't exist"):
		return errorx.NewWithMsg(errorx.ErrCodeSystem, "数据列不存在")

	default:
		// 保留原始错误，但包装为数据库错误
		return errorx.NewWithMsg(errorx.ErrCodeDatabase, err.Error())
	}
}

// IsNotFoundError 判断是否为记录不存在错误
func (a *ORMErrorAdapter) IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	// 检查是否为统一的NotFound错误
	if codeErr, ok := err.(*errorx.CodeError); ok {
		return codeErr.GetCode() == errorx.ErrCodeNotFound
	}

	// 检查原始错误
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows)
}

// IsAlreadyExistsError 判断是否为记录已存在错误
func (a *ORMErrorAdapter) IsAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}

	// 检查是否为统一的AlreadyExists错误
	if codeErr, ok := err.(*errorx.CodeError); ok {
		return codeErr.GetCode() == errorx.ErrCodeAlreadyExists
	}

	// 检查原始错误消息
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "duplicate") || strings.Contains(errMsg, "unique constraint")
}

// IsDatabaseError 判断是否为数据库错误
func (a *ORMErrorAdapter) IsDatabaseError(err error) bool {
	if err == nil {
		return false
	}

	// 检查是否为统一的数据库错误
	if codeErr, ok := err.(*errorx.CodeError); ok {
		return codeErr.GetCode() == errorx.ErrCodeDatabase
	}

	return false
}
