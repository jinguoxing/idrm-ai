package data_catalog_column

import "time"

// DataCatalogColumn 数据资源目录信息项数据结构
// 对应表：t_data_catalog_column
type DataCatalogColumn struct {
	// ========== 主键和关联 ==========
	Id         string `json:"id" gorm:"column:id;primaryKey;size:20"`
	CatalogId  string `json:"catalogId" gorm:"column:catalog_id;size:20;not null;index:idx_catalog_id"`

	// ========== 基础信息 ==========
	ColumnName   string `json:"columnName" gorm:"column:column_name;size:128;not null"`
	ColumnType   string `json:"columnType" gorm:"column:column_type;size:20;not null"`
	ColumnLength int    `json:"columnLength" gorm:"column:column_length;default:0"`
	ColumnScale  int    `json:"columnScale" gorm:"column:column_scale;default:0"`

	// ========== 描述信息 ==========
	ColumnCode     string `json:"columnCode" gorm:"column:column_code;size:128"`
	ColumnAlias    string `json:"columnAlias" gorm:"column:column_alias;size:500"`
	ColumnDescribe string `json:"columnDescribe" gorm:"column:column_describe;size:500"`

	// ========== 数据属性 ==========
	DataType      string `json:"dataType" gorm:"column:data_type;size:100"`
	DataFormat    string `json:"dataFormat" gorm:"column:data_format;size:100"`
	DataUnit      string `json:"dataUnit" gorm:"column:data_unit;size:50"`
	ValueRange    string `json:"valueRange" gorm:"column:value_range;size:500"`
	IsPrimaryKey  bool   `json:"isPrimaryKey" gorm:"column:is_primary_key;default:0"`
	IsSensitive   *bool  `json:"isSensitive" gorm:"column:is_sensitive;default:0"`
	IsRequired    bool   `json:"isRequired" gorm:"column:is_required;default:0"`
	IsUnique      bool   `json:"isUnique" gorm:"column:is_unique;default:0"`

	// ========== 扩展信息 ==========
	DictionaryId   string `json:"dictionaryId" gorm:"column:dictionary_id;size:50"`
	StandardCodeId string `json:"standardCodeId" gorm:"column:standard_code_id;size:50"`
	SortOrder      int    `json:"sortOrder" gorm:"column:sort_order;default:0"`

	// ========== 时间戳 ==========
	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

// TableName 指定表名
func (DataCatalogColumn) TableName() string {
	return "t_data_catalog_column"
}
