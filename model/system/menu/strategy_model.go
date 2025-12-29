package menu

import (
	"context"
	"time"

	"idrm/pkg/orm"

	"github.com/zeromicro/go-zero/core/logx"
)

// StrategyModel 策略模型包装器
// 根据配置的策略动态选择GORM或SQLx实现
type StrategyModel struct {
	gormModel Model
	sqlxModel Model
	selector  *orm.ORMSelector
	logger    logx.Logger
}

// Insert 插入新的菜单记录
func (s *StrategyModel) Insert(ctx context.Context, data *Menu) (*Menu, error) {
	return s.executeWithStrategy(ctx, orm.OpTypeInsert, func(model Model) (*Menu, error) {
		return model.Insert(ctx, data)
	})
}

// FindOne 根据ID查询单个菜单
func (s *StrategyModel) FindOne(ctx context.Context, id int64) (*Menu, error) {
	return s.executeWithStrategy(ctx, orm.OpTypeFindOne, func(model Model) (*Menu, error) {
		return model.FindOne(ctx, id)
	})
}

// FindByPermTag 根据权限标识查询菜单
func (s *StrategyModel) FindByPermTag(ctx context.Context, permTag string) (*Menu, error) {
	return s.executeWithStrategy(ctx, orm.OpTypeFindOne, func(model Model) (*Menu, error) {
		return model.FindByPermTag(ctx, permTag)
	})
}

// Update 更新菜单记录
func (s *StrategyModel) Update(ctx context.Context, data *Menu) error {
	_, err := s.executeWithStrategy(ctx, orm.OpTypeUpdate, func(model Model) (*Menu, error) {
		return nil, model.Update(ctx, data)
	})
	return err
}

// Delete 删除菜单记录（软删除）
func (s *StrategyModel) Delete(ctx context.Context, id int64) error {
	_, err := s.executeWithStrategy(ctx, orm.OpTypeDelete, func(model Model) (*Menu, error) {
		return nil, model.Delete(ctx, id)
	})
	return err
}

// FindAll 查询所有菜单（支持状态过滤）
func (s *StrategyModel) FindAll(ctx context.Context, status int) ([]*Menu, error) {
	start := time.Now()
	
	ormType, err := s.selector.SelectORM(orm.OpTypeFindMany)
	if err != nil {
		return nil, err
	}
	
	model := s.getModel(ormType)
	if model == nil {
		return nil, ErrNotFound
	}
	
	result, err := model.FindAll(ctx, status)
	
	// 记录指标
	latency := time.Since(start)
	s.selector.GetMetrics().RecordQuery(ormType, orm.OpTypeFindMany, latency, err)
	
	// 适配错误
	if err != nil {
		err = s.selector.AdaptError(err, ormType)
	}
	
	s.logOperation("FindAll", ormType, latency, err)
	
	return result, err
}

// FindByParentId 根据父级ID查询子菜单列表
func (s *StrategyModel) FindByParentId(ctx context.Context, parentId int64) ([]*Menu, error) {
	start := time.Now()
	
	ormType, err := s.selector.SelectORM(orm.OpTypeFindMany)
	if err != nil {
		return nil, err
	}
	
	model := s.getModel(ormType)
	if model == nil {
		return nil, ErrNotFound
	}
	
	result, err := model.FindByParentId(ctx, parentId)
	
	// 记录指标
	latency := time.Since(start)
	s.selector.GetMetrics().RecordQuery(ormType, orm.OpTypeFindMany, latency, err)
	
	// 适配错误
	if err != nil {
		err = s.selector.AdaptError(err, ormType)
	}
	
	s.logOperation("FindByParentId", ormType, latency, err)
	
	return result, err
}

// CountByParentId 统计指定父级下的子菜单数量
func (s *StrategyModel) CountByParentId(ctx context.Context, parentId int64) (int64, error) {
	start := time.Now()
	
	ormType, err := s.selector.SelectORM(orm.OpTypeFindOne)
	if err != nil {
		return 0, err
	}
	
	model := s.getModel(ormType)
	if model == nil {
		return 0, ErrNotFound
	}
	
	result, err := model.CountByParentId(ctx, parentId)
	
	// 记录指标
	latency := time.Since(start)
	s.selector.GetMetrics().RecordQuery(ormType, orm.OpTypeFindOne, latency, err)
	
	// 适配错误
	if err != nil {
		err = s.selector.AdaptError(err, ormType)
	}
	
	s.logOperation("CountByParentId", ormType, latency, err)
	
	return result, err
}

// HasChildren 检查菜单是否有子菜单
func (s *StrategyModel) HasChildren(ctx context.Context, parentId int64) (bool, error) {
	start := time.Now()
	
	ormType, err := s.selector.SelectORM(orm.OpTypeFindOne)
	if err != nil {
		return false, err
	}
	
	model := s.getModel(ormType)
	if model == nil {
		return false, ErrNotFound
	}
	
	result, err := model.HasChildren(ctx, parentId)
	
	// 记录指标
	latency := time.Since(start)
	s.selector.GetMetrics().RecordQuery(ormType, orm.OpTypeFindOne, latency, err)
	
	// 适配错误
	if err != nil {
		err = s.selector.AdaptError(err, ormType)
	}
	
	s.logOperation("HasChildren", ormType, latency, err)
	
	return result, err
}

// WithTx 创建事务副本
func (s *StrategyModel) WithTx(tx interface{}) Model {
	// 事务操作优先使用GORM
	if s.gormModel != nil {
		return s.gormModel.WithTx(tx)
	}
	
	if s.sqlxModel != nil {
		return s.sqlxModel.WithTx(tx)
	}
	
	return nil
}

// Trans 在事务中执行多个操作
func (s *StrategyModel) Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error {
	start := time.Now()
	
	ormType, err := s.selector.SelectORM(orm.OpTypeTransaction)
	if err != nil {
		return err
	}
	
	model := s.getModel(ormType)
	if model == nil {
		return ErrNotFound
	}
	
	err = model.Trans(ctx, fn)
	
	// 记录指标
	latency := time.Since(start)
	s.selector.GetMetrics().RecordQuery(ormType, orm.OpTypeTransaction, latency, err)
	
	// 适配错误
	if err != nil {
		err = s.selector.AdaptError(err, ormType)
	}
	
	s.logOperation("Trans", ormType, latency, err)
	
	return err
}

// executeWithStrategy 使用策略执行操作
func (s *StrategyModel) executeWithStrategy(ctx context.Context, opType orm.OperationType, fn func(Model) (*Menu, error)) (*Menu, error) {
	start := time.Now()
	
	ormType, err := s.selector.SelectORM(opType)
	if err != nil {
		return nil, err
	}
	
	model := s.getModel(ormType)
	if model == nil {
		return nil, ErrNotFound
	}
	
	result, err := fn(model)
	
	// 记录指标
	latency := time.Since(start)
	s.selector.GetMetrics().RecordQuery(ormType, opType, latency, err)
	
	// 适配错误
	if err != nil {
		err = s.selector.AdaptError(err, ormType)
	}
	
	s.logOperation(string(opType), ormType, latency, err)
	
	return result, err
}

// getModel 根据ORM类型获取模型实例
func (s *StrategyModel) getModel(ormType orm.ORMType) Model {
	switch ormType {
	case orm.ORMTypeGORM:
		return s.gormModel
	case orm.ORMTypeSQLx:
		return s.sqlxModel
	default:
		return nil
	}
}

// logOperation 记录操作日志
func (s *StrategyModel) logOperation(operation string, ormType orm.ORMType, latency time.Duration, err error) {
	if s.logger == nil {
		s.logger = logx.WithContext(context.Background())
	}
	
	if err != nil {
		s.logger.Errorw("Menu operation failed",
			logx.Field("operation", operation),
			logx.Field("ormType", ormType),
			logx.Field("latency", latency),
			logx.Field("error", err.Error()),
		)
	} else {
		s.logger.Infow("Menu operation completed",
			logx.Field("operation", operation),
			logx.Field("ormType", ormType),
			logx.Field("latency", latency),
		)
	}
}

// GetMetrics 获取性能指标
func (s *StrategyModel) GetMetrics() *orm.ORMMetrics {
	return s.selector.GetMetrics()
}

// ResetMetrics 重置性能指标
func (s *StrategyModel) ResetMetrics() {
	s.selector.ResetMetrics()
}