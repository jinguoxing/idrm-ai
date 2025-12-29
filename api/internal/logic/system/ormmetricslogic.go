package system

import (
	"context"

	"idrm/api/internal/svc"
	"idrm/pkg/orm"

	"github.com/zeromicro/go-zero/core/logx"
)

// ORMMetricsLogic ORM指标业务逻辑
type ORMMetricsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewORMMetricsLogic 创建ORM指标业务逻辑实例
func NewORMMetricsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ORMMetricsLogic {
	return &ORMMetricsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ORMMetricsResponse ORM指标响应
type ORMMetricsResponse struct {
	*orm.ORMMetrics
	DatabaseConnections DatabaseConnectionInfo `json:"databaseConnections"`
	ORMStrategy         string                  `json:"ormStrategy"`
}

// DatabaseConnectionInfo 数据库连接信息
type DatabaseConnectionInfo struct {
	SQLxAvailable bool `json:"sqlxAvailable"`
	GORMAvailable bool `json:"gormAvailable"`
}

// ORMMetrics 获取ORM性能指标
func (l *ORMMetricsLogic) ORMMetrics() (*ORMMetricsResponse, error) {
	// 获取ORM指标
	var metrics *orm.ORMMetrics
	if l.svcCtx.ORMMetrics != nil {
		metrics = l.svcCtx.ORMMetrics.GetMetrics()
	} else {
		// 如果没有启用指标，尝试从MenuModel获取
		if strategyModel, ok := l.svcCtx.MenuModel.(interface{ GetMetrics() *orm.ORMMetrics }); ok {
			metrics = strategyModel.GetMetrics()
		} else {
			// 返回空指标
			metrics = orm.NewORMMetrics().GetMetrics()
		}
	}

	// 检查数据库连接状态
	dbConnInfo := DatabaseConnectionInfo{
		SQLxAvailable: l.svcCtx.SqlConn != nil,
		GORMAvailable: l.svcCtx.GormDB != nil,
	}

	// 如果连接存在，进行健康检查
	if l.svcCtx.SqlConn != nil {
		if err := l.svcCtx.SqlConn.PingContext(l.ctx); err != nil {
			dbConnInfo.SQLxAvailable = false
			l.Logger.Errorw("SQLx connection health check failed", logx.Field("error", err))
		}
	}

	if l.svcCtx.GormDB != nil {
		if sqlDB, err := l.svcCtx.GormDB.DB(); err == nil {
			if err := sqlDB.PingContext(l.ctx); err != nil {
				dbConnInfo.GORMAvailable = false
				l.Logger.Errorw("GORM connection health check failed", logx.Field("error", err))
			}
		} else {
			dbConnInfo.GORMAvailable = false
			l.Logger.Errorw("Failed to get GORM underlying DB", logx.Field("error", err))
		}
	}

	response := &ORMMetricsResponse{
		ORMMetrics:          metrics,
		DatabaseConnections: dbConnInfo,
		ORMStrategy:         string(l.svcCtx.Config.ORM.Strategy),
	}

	l.Logger.Infow("ORM metrics retrieved",
		logx.Field("totalQueries", metrics.GetTotalQueries()),
		logx.Field("gormQueries", metrics.GORMQueryCount),
		logx.Field("sqlxQueries", metrics.SQLxQueryCount),
		logx.Field("failoverCount", metrics.FailoverCount),
	)

	return response, nil
}

// ResetORMMetrics 重置ORM性能指标
func (l *ORMMetricsLogic) ResetORMMetrics() error {
	// 重置全局指标
	if l.svcCtx.ORMMetrics != nil {
		l.svcCtx.ORMMetrics.Reset()
	}

	// 重置MenuModel指标（如果支持）
	if strategyModel, ok := l.svcCtx.MenuModel.(interface{ ResetMetrics() }); ok {
		strategyModel.ResetMetrics()
	}

	l.Logger.Info("ORM metrics reset successfully")

	return nil
}