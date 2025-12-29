package data_catalog_category

import "time"

// DataCatalogCategory 数据资源目录分类关联数据结构
// 对应表：t_data_catalog_category
type DataCatalogCategory struct {
	// ========== 主键和关联 ==========
	Id        string `json:"id" gorm:"column:id;primaryKey;size:20"`
	CatalogId string `json:"catalogId" gorm:"column:catalog_id;size:20;not null;index:idx_catalog_id"`
	// 资源属性分类节点ID（多选）
	CategoryId string `json:"categoryId" gorm:"column:category_id;size:20;not null"`

	// ========== 时间戳 ==========
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

// TableName 指定表名
func (DataCatalogCategory) TableName() string {
	return "t_data_catalog_category"
}
