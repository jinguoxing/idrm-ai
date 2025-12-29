package orm

import (
	"sync"
	"sync/atomic"
	"time"
)

// ORMType ORM类型
type ORMType string

const (
	ORMTypeGORM ORMType = "gorm"
	ORMTypeSQLx ORMType = "sqlx"
)

// OperationType 操作类型
type OperationType string

const (
	OpTypeInsert      OperationType = "insert"
	OpTypeFindOne     OperationType = "find_one"
	OpTypeFindMany    OperationType = "find_many"
	OpTypeUpdate      OperationType = "update"
	OpTypeDelete      OperationType = "delete"
	OpTypeTransaction OperationType = "transaction"
)

// ORMMetrics ORM性能指标
type ORMMetrics struct {
	// GORM指标
	GORMQueryCount    int64         `json:"gormQueryCount"`    // GORM查询次数
	GORMTotalLatency  int64         `json:"gormTotalLatency"`  // GORM总延迟（纳秒）
	GORMAvgLatency    time.Duration `json:"gormAvgLatency"`    // GORM平均延迟
	GORMErrorCount    int64         `json:"gormErrorCount"`    // GORM错误次数

	// SQLx指标
	SQLxQueryCount    int64         `json:"sqlxQueryCount"`    // SQLx查询次数
	SQLxTotalLatency  int64         `json:"sqlxTotalLatency"`  // SQLx总延迟（纳秒）
	SQLxAvgLatency    time.Duration `json:"sqlxAvgLatency"`    // SQLx平均延迟
	SQLxErrorCount    int64         `json:"sqlxErrorCount"`    // SQLx错误次数

	// 切换指标
	FailoverCount     int64         `json:"failoverCount"`     // 故障切换次数
	LastFailoverTime  time.Time     `json:"lastFailoverTime"`  // 最后切换时间

	// 操作类型统计
	OperationStats    map[OperationType]*OperationMetrics `json:"operationStats"`

	mu sync.RWMutex // 保护OperationStats的并发访问
}

// OperationMetrics 操作指标
type OperationMetrics struct {
	Count        int64         `json:"count"`        // 操作次数
	TotalLatency int64         `json:"totalLatency"` // 总延迟（纳秒）
	AvgLatency   time.Duration `json:"avgLatency"`   // 平均延迟
	ErrorCount   int64         `json:"errorCount"`   // 错误次数
}

// NewORMMetrics 创建新的ORM指标实例
func NewORMMetrics() *ORMMetrics {
	return &ORMMetrics{
		OperationStats: make(map[OperationType]*OperationMetrics),
	}
}

// RecordQuery 记录查询指标
func (m *ORMMetrics) RecordQuery(ormType ORMType, opType OperationType, latency time.Duration, err error) {
	latencyNs := latency.Nanoseconds()

	// 记录ORM类型指标
	switch ormType {
	case ORMTypeGORM:
		atomic.AddInt64(&m.GORMQueryCount, 1)
		atomic.AddInt64(&m.GORMTotalLatency, latencyNs)
		if err != nil {
			atomic.AddInt64(&m.GORMErrorCount, 1)
		}
	case ORMTypeSQLx:
		atomic.AddInt64(&m.SQLxQueryCount, 1)
		atomic.AddInt64(&m.SQLxTotalLatency, latencyNs)
		if err != nil {
			atomic.AddInt64(&m.SQLxErrorCount, 1)
		}
	}

	// 记录操作类型指标
	m.mu.Lock()
	if m.OperationStats[opType] == nil {
		m.OperationStats[opType] = &OperationMetrics{}
	}
	opMetrics := m.OperationStats[opType]
	m.mu.Unlock()

	atomic.AddInt64(&opMetrics.Count, 1)
	atomic.AddInt64(&opMetrics.TotalLatency, latencyNs)
	if err != nil {
		atomic.AddInt64(&opMetrics.ErrorCount, 1)
	}
}

// RecordFailover 记录故障切换
func (m *ORMMetrics) RecordFailover() {
	atomic.AddInt64(&m.FailoverCount, 1)
	m.mu.Lock()
	m.LastFailoverTime = time.Now()
	m.mu.Unlock()
}

// GetMetrics 获取当前指标（计算平均值）
func (m *ORMMetrics) GetMetrics() *ORMMetrics {
	// 创建副本以避免并发问题
	result := &ORMMetrics{
		GORMQueryCount:   atomic.LoadInt64(&m.GORMQueryCount),
		GORMTotalLatency: atomic.LoadInt64(&m.GORMTotalLatency),
		GORMErrorCount:   atomic.LoadInt64(&m.GORMErrorCount),
		SQLxQueryCount:   atomic.LoadInt64(&m.SQLxQueryCount),
		SQLxTotalLatency: atomic.LoadInt64(&m.SQLxTotalLatency),
		SQLxErrorCount:   atomic.LoadInt64(&m.SQLxErrorCount),
		FailoverCount:    atomic.LoadInt64(&m.FailoverCount),
		OperationStats:   make(map[OperationType]*OperationMetrics),
	}

	// 计算平均延迟
	if result.GORMQueryCount > 0 {
		result.GORMAvgLatency = time.Duration(result.GORMTotalLatency / result.GORMQueryCount)
	}
	if result.SQLxQueryCount > 0 {
		result.SQLxAvgLatency = time.Duration(result.SQLxTotalLatency / result.SQLxQueryCount)
	}

	// 复制操作统计
	m.mu.RLock()
	result.LastFailoverTime = m.LastFailoverTime
	for opType, opMetrics := range m.OperationStats {
		count := atomic.LoadInt64(&opMetrics.Count)
		totalLatency := atomic.LoadInt64(&opMetrics.TotalLatency)
		errorCount := atomic.LoadInt64(&opMetrics.ErrorCount)
		
		result.OperationStats[opType] = &OperationMetrics{
			Count:        count,
			TotalLatency: totalLatency,
			ErrorCount:   errorCount,
		}
		
		if count > 0 {
			result.OperationStats[opType].AvgLatency = time.Duration(totalLatency / count)
		}
	}
	m.mu.RUnlock()

	return result
}

// Reset 重置所有指标
func (m *ORMMetrics) Reset() {
	atomic.StoreInt64(&m.GORMQueryCount, 0)
	atomic.StoreInt64(&m.GORMTotalLatency, 0)
	atomic.StoreInt64(&m.GORMErrorCount, 0)
	atomic.StoreInt64(&m.SQLxQueryCount, 0)
	atomic.StoreInt64(&m.SQLxTotalLatency, 0)
	atomic.StoreInt64(&m.SQLxErrorCount, 0)
	atomic.StoreInt64(&m.FailoverCount, 0)

	m.mu.Lock()
	m.LastFailoverTime = time.Time{}
	for opType := range m.OperationStats {
		delete(m.OperationStats, opType)
	}
	m.mu.Unlock()
}

// GetSuccessRate 获取成功率
func (m *ORMMetrics) GetSuccessRate(ormType ORMType) float64 {
	var totalCount, errorCount int64
	
	switch ormType {
	case ORMTypeGORM:
		totalCount = atomic.LoadInt64(&m.GORMQueryCount)
		errorCount = atomic.LoadInt64(&m.GORMErrorCount)
	case ORMTypeSQLx:
		totalCount = atomic.LoadInt64(&m.SQLxQueryCount)
		errorCount = atomic.LoadInt64(&m.SQLxErrorCount)
	}

	if totalCount == 0 {
		return 0.0
	}

	return float64(totalCount-errorCount) / float64(totalCount) * 100.0
}

// GetTotalQueries 获取总查询次数
func (m *ORMMetrics) GetTotalQueries() int64 {
	return atomic.LoadInt64(&m.GORMQueryCount) + atomic.LoadInt64(&m.SQLxQueryCount)
}