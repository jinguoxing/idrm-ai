package orm

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"idrm/pkg/config"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// ModelFactory 模型工厂接口
type ModelFactory interface {
	// CreateModel 创建模型实例
	CreateModel(sqlConn *sql.DB, gormDB *gorm.DB) interface{}
	
	// GetModelName 获取模型名称
	GetModelName() string
}

// ORMSelector ORM选择器
// 根据配置策略和运行时条件选择最优的ORM实现
type ORMSelector struct {
	strategy     config.ORMStrategy
	sqlConn      *sql.DB
	gormDB       *gorm.DB
	metrics      *ORMMetrics
	errorAdapter *ORMErrorAdapter
	logger       logx.Logger
}

// NewORMSelector 创建ORM选择器
func NewORMSelector(strategy config.ORMStrategy, sqlConn *sql.DB, gormDB *gorm.DB, metrics *ORMMetrics) *ORMSelector {
	return &ORMSelector{
		strategy:     strategy,
		sqlConn:      sqlConn,
		gormDB:       gormDB,
		metrics:      metrics,
		errorAdapter: NewORMErrorAdapter(),
		logger:       logx.WithContext(context.Background()),
	}
}

// SelectORM 选择ORM实现
func (s *ORMSelector) SelectORM(opType OperationType) (ORMType, error) {
	switch s.strategy {
	case config.ORMStrategyGORMFirst:
		return s.selectGORMFirst()
	
	case config.ORMStrategySQLxFirst:
		return s.selectSQLxFirst()
	
	case config.ORMStrategyGORMOnly:
		return s.selectGORMOnly()
	
	case config.ORMStrategySQLxOnly:
		return s.selectSQLxOnly()
	
	case config.ORMStrategyAuto:
		return s.selectAuto(opType)
	
	default:
		s.logger.Errorf("Unknown ORM strategy: %s, fallback to gorm_first", s.strategy)
		return s.selectGORMFirst()
	}
}

// selectGORMFirst GORM优先策略
func (s *ORMSelector) selectGORMFirst() (ORMType, error) {
	if s.isGORMAvailable() {
		return ORMTypeGORM, nil
	}
	
	if s.isSQLxAvailable() {
		s.logger.Infow("GORM not available, fallback to SQLx")
		s.metrics.RecordFailover()
		return ORMTypeSQLx, nil
	}
	
	return "", fmt.Errorf("no ORM available")
}

// selectSQLxFirst SQLx优先策略
func (s *ORMSelector) selectSQLxFirst() (ORMType, error) {
	if s.isSQLxAvailable() {
		return ORMTypeSQLx, nil
	}
	
	if s.isGORMAvailable() {
		s.logger.Infow("SQLx not available, fallback to GORM")
		s.metrics.RecordFailover()
		return ORMTypeGORM, nil
	}
	
	return "", fmt.Errorf("no ORM available")
}

// selectGORMOnly 仅GORM策略
func (s *ORMSelector) selectGORMOnly() (ORMType, error) {
	if s.isGORMAvailable() {
		return ORMTypeGORM, nil
	}
	
	return "", fmt.Errorf("GORM not available")
}

// selectSQLxOnly 仅SQLx策略
func (s *ORMSelector) selectSQLxOnly() (ORMType, error) {
	if s.isSQLxAvailable() {
		return ORMTypeSQLx, nil
	}
	
	return "", fmt.Errorf("SQLx not available")
}

// selectAuto 智能选择策略
func (s *ORMSelector) selectAuto(opType OperationType) (ORMType, error) {
	// 根据操作类型选择最优ORM
	switch opType {
	case OpTypeInsert, OpTypeUpdate, OpTypeDelete:
		// 写操作：优先使用GORM（更好的事务支持）
		if s.isGORMAvailable() {
			return ORMTypeGORM, nil
		}
		if s.isSQLxAvailable() {
			return ORMTypeSQLx, nil
		}
	
	case OpTypeFindOne, OpTypeFindMany:
		// 读操作：根据性能指标选择
		return s.selectByPerformance()
	
	case OpTypeTransaction:
		// 事务操作：优先使用GORM
		if s.isGORMAvailable() {
			return ORMTypeGORM, nil
		}
		if s.isSQLxAvailable() {
			return ORMTypeSQLx, nil
		}
	
	default:
		// 默认策略：GORM优先
		return s.selectGORMFirst()
	}
	
	return "", fmt.Errorf("no ORM available")
}

// selectByPerformance 根据性能指标选择ORM
func (s *ORMSelector) selectByPerformance() (ORMType, error) {
	metrics := s.metrics.GetMetrics()
	
	// 如果没有历史数据，使用默认策略
	if metrics.GetTotalQueries() < 100 {
		return s.selectGORMFirst()
	}
	
	// 比较平均延迟
	gormAvgLatency := metrics.GORMAvgLatency
	sqlxAvgLatency := metrics.SQLxAvgLatency
	
	// 比较成功率
	gormSuccessRate := s.metrics.GetSuccessRate(ORMTypeGORM)
	sqlxSuccessRate := s.metrics.GetSuccessRate(ORMTypeSQLx)
	
	// 综合评分：延迟权重70%，成功率权重30%
	gormScore := s.calculateScore(gormAvgLatency, gormSuccessRate)
	sqlxScore := s.calculateScore(sqlxAvgLatency, sqlxSuccessRate)
	
	s.logger.Infow("Performance comparison",
		logx.Field("gormScore", gormScore),
		logx.Field("sqlxScore", sqlxScore),
		logx.Field("gormAvgLatency", gormAvgLatency),
		logx.Field("sqlxAvgLatency", sqlxAvgLatency),
		logx.Field("gormSuccessRate", gormSuccessRate),
		logx.Field("sqlxSuccessRate", sqlxSuccessRate),
	)
	
	// 选择评分更高的ORM
	if gormScore > sqlxScore && s.isGORMAvailable() {
		return ORMTypeGORM, nil
	}
	
	if s.isSQLxAvailable() {
		return ORMTypeSQLx, nil
	}
	
	if s.isGORMAvailable() {
		return ORMTypeGORM, nil
	}
	
	return "", fmt.Errorf("no ORM available")
}

// calculateScore 计算ORM评分
func (s *ORMSelector) calculateScore(avgLatency time.Duration, successRate float64) float64 {
	// 延迟评分：延迟越低评分越高（最大100分）
	latencyScore := 100.0
	if avgLatency > 0 {
		// 假设1ms为基准，延迟每增加1ms扣10分
		latencyPenalty := float64(avgLatency.Milliseconds()) * 10
		latencyScore = 100.0 - latencyPenalty
		if latencyScore < 0 {
			latencyScore = 0
		}
	}
	
	// 综合评分：延迟权重70%，成功率权重30%
	return latencyScore*0.7 + successRate*0.3
}

// isGORMAvailable 检查GORM是否可用
func (s *ORMSelector) isGORMAvailable() bool {
	if s.gormDB == nil {
		return false
	}
	
	// 简单的健康检查
	sqlDB, err := s.gormDB.DB()
	if err != nil {
		return false
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	return sqlDB.PingContext(ctx) == nil
}

// isSQLxAvailable 检查SQLx是否可用
func (s *ORMSelector) isSQLxAvailable() bool {
	if s.sqlConn == nil {
		return false
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	return s.sqlConn.PingContext(ctx) == nil
}

// GetMetrics 获取指标
func (s *ORMSelector) GetMetrics() *ORMMetrics {
	return s.metrics.GetMetrics()
}

// ResetMetrics 重置指标
func (s *ORMSelector) ResetMetrics() {
	s.metrics.Reset()
}

// AdaptError 适配错误
func (s *ORMSelector) AdaptError(err error, ormType ORMType) error {
	return s.errorAdapter.AdaptError(err, ormType)
}

// Select 根据策略选择 GORM 或 SQLx Model
// gormModel 和 sqlxModel 分别是两个ORM实现的Model接口
// 返回被选中的Model（非nil），如果都不可用则返回nil
func (s *ORMSelector) Select(gormModel interface{}, sqlxModel interface{}) interface{} {
	ormType, err := s.SelectORM(OpTypeFindOne)
	if err != nil {
		s.logger.Errorf("Failed to select ORM: %v, fallback to available", err)
		// 降级处理：优先使用可用的
		if gormModel != nil && s.isGORMAvailable() {
			return gormModel
		}
		if sqlxModel != nil && s.isSQLxAvailable() {
			return sqlxModel
		}
		return nil
	}

	switch ormType {
	case ORMTypeGORM:
		if gormModel != nil {
			return gormModel
		}
		// 降级到SQLx
		if sqlxModel != nil && s.isSQLxAvailable() {
			s.metrics.RecordFailover()
			return sqlxModel
		}
	case ORMTypeSQLx:
		if sqlxModel != nil {
			return sqlxModel
		}
		// 降级到GORM
		if gormModel != nil && s.isGORMAvailable() {
			s.metrics.RecordFailover()
			return gormModel
		}
	}

	return nil
}