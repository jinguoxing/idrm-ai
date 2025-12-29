package orm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestORMMetrics_RecordQuery(t *testing.T) {
	metrics := NewORMMetrics()

	// 记录GORM查询
	metrics.RecordQuery(ORMTypeGORM, OpTypeInsert, 10*time.Millisecond, nil)
	metrics.RecordQuery(ORMTypeGORM, OpTypeFindOne, 5*time.Millisecond, nil)

	// 记录SQLx查询
	metrics.RecordQuery(ORMTypeSQLx, OpTypeInsert, 8*time.Millisecond, nil)
	metrics.RecordQuery(ORMTypeSQLx, OpTypeFindOne, 3*time.Millisecond, assert.AnError)

	// 获取指标
	result := metrics.GetMetrics()

	// 验证GORM指标
	assert.Equal(t, int64(2), result.GORMQueryCount)
	assert.Equal(t, int64(0), result.GORMErrorCount)
	// 平均延迟应该在合理范围内
	assert.True(t, result.GORMAvgLatency >= 7*time.Millisecond && result.GORMAvgLatency <= 8*time.Millisecond)

	// 验证SQLx指标
	assert.Equal(t, int64(2), result.SQLxQueryCount)
	assert.Equal(t, int64(1), result.SQLxErrorCount)
	// 平均延迟应该在合理范围内
	assert.True(t, result.SQLxAvgLatency >= 5*time.Millisecond && result.SQLxAvgLatency <= 6*time.Millisecond)

	// 验证操作统计
	assert.Contains(t, result.OperationStats, OpTypeInsert)
	assert.Contains(t, result.OperationStats, OpTypeFindOne)

	insertStats := result.OperationStats[OpTypeInsert]
	assert.Equal(t, int64(2), insertStats.Count)
	assert.Equal(t, int64(0), insertStats.ErrorCount)

	findOneStats := result.OperationStats[OpTypeFindOne]
	assert.Equal(t, int64(2), findOneStats.Count)
	assert.Equal(t, int64(1), findOneStats.ErrorCount)
}

func TestORMMetrics_RecordFailover(t *testing.T) {
	metrics := NewORMMetrics()

	// 记录故障切换
	metrics.RecordFailover()
	metrics.RecordFailover()

	result := metrics.GetMetrics()
	assert.Equal(t, int64(2), result.FailoverCount)
	assert.False(t, result.LastFailoverTime.IsZero())
}

func TestORMMetrics_GetSuccessRate(t *testing.T) {
	metrics := NewORMMetrics()

	// 记录一些查询（包含错误）
	metrics.RecordQuery(ORMTypeGORM, OpTypeInsert, 10*time.Millisecond, nil)
	metrics.RecordQuery(ORMTypeGORM, OpTypeInsert, 10*time.Millisecond, nil)
	metrics.RecordQuery(ORMTypeGORM, OpTypeInsert, 10*time.Millisecond, assert.AnError)

	// GORM成功率应该是 2/3 = 66.67%
	successRate := metrics.GetSuccessRate(ORMTypeGORM)
	assert.InDelta(t, 66.67, successRate, 0.01)

	// SQLx没有查询，成功率应该是0
	sqlxSuccessRate := metrics.GetSuccessRate(ORMTypeSQLx)
	assert.Equal(t, 0.0, sqlxSuccessRate)
}

func TestORMMetrics_Reset(t *testing.T) {
	metrics := NewORMMetrics()

	// 记录一些数据
	metrics.RecordQuery(ORMTypeGORM, OpTypeInsert, 10*time.Millisecond, nil)
	metrics.RecordFailover()

	// 验证数据存在
	result := metrics.GetMetrics()
	assert.Equal(t, int64(1), result.GORMQueryCount)
	assert.Equal(t, int64(1), result.FailoverCount)

	// 重置指标
	metrics.Reset()

	// 验证数据已清零
	result = metrics.GetMetrics()
	assert.Equal(t, int64(0), result.GORMQueryCount)
	assert.Equal(t, int64(0), result.SQLxQueryCount)
	assert.Equal(t, int64(0), result.FailoverCount)
	assert.True(t, result.LastFailoverTime.IsZero())
	assert.Empty(t, result.OperationStats)
}

func TestORMMetrics_GetTotalQueries(t *testing.T) {
	metrics := NewORMMetrics()

	// 记录查询
	metrics.RecordQuery(ORMTypeGORM, OpTypeInsert, 10*time.Millisecond, nil)
	metrics.RecordQuery(ORMTypeGORM, OpTypeInsert, 10*time.Millisecond, nil)
	metrics.RecordQuery(ORMTypeSQLx, OpTypeInsert, 10*time.Millisecond, nil)

	// 总查询数应该是3
	totalQueries := metrics.GetTotalQueries()
	assert.Equal(t, int64(3), totalQueries)
}

// 并发安全测试
func TestORMMetrics_ConcurrentAccess(t *testing.T) {
	metrics := NewORMMetrics()

	// 启动多个goroutine并发记录指标
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				metrics.RecordQuery(ORMTypeGORM, OpTypeInsert, time.Millisecond, nil)
				metrics.RecordQuery(ORMTypeSQLx, OpTypeFindOne, time.Millisecond, nil)
				metrics.RecordFailover()
			}
			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证结果
	result := metrics.GetMetrics()
	assert.Equal(t, int64(1000), result.GORMQueryCount)
	assert.Equal(t, int64(1000), result.SQLxQueryCount)
	assert.Equal(t, int64(1000), result.FailoverCount)
}

// 性能基准测试
func BenchmarkORMMetrics_RecordQuery(b *testing.B) {
	metrics := NewORMMetrics()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordQuery(ORMTypeGORM, OpTypeInsert, time.Millisecond, nil)
	}
}

func BenchmarkORMMetrics_GetMetrics(b *testing.B) {
	metrics := NewORMMetrics()

	// 预先记录一些数据
	for i := 0; i < 1000; i++ {
		metrics.RecordQuery(ORMTypeGORM, OpTypeInsert, time.Millisecond, nil)
		metrics.RecordQuery(ORMTypeSQLx, OpTypeFindOne, time.Millisecond, nil)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.GetMetrics()
	}
}