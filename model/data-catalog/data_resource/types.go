package data_resource

import "time"

// ResourceType 数据资源类型枚举
type ResourceType string

const (
	ResourceTypeDatabaseTable ResourceType = "database_table" // 数据表
	ResourceTypeLogicalView   ResourceType = "logical_view"   // 逻辑视图
	ResourceTypeApiInterface  ResourceType = "api_interface"  // API接口
	ResourceTypeFile          ResourceType = "file"           // 文件
)

// DataResource 数据资源挂接数据结构
// 对应表：t_data_resource
type DataResource struct {
	// ========== 主键和关联 ==========
	Id        string `json:"id" gorm:"column:id;primaryKey;size:20"`
	CatalogId string `json:"catalogId" gorm:"column:catalog_id;size:20;not null;index:idx_catalog_id"`

	// ========== 资源信息 ==========
	ResourceType ResourceType `json:"resourceType" gorm:"column:resource_type;size:20;not null"`
	ResourceName string       `json:"resourceName" gorm:"column:resource_name;size:500;not null"`
	ResourceCode string       `json:"resourceCode" gorm:"column:resource_code;size:100"`
	ResourceDesc string       `json:"resourceDesc" gorm:"column:resource_desc;size:1000"`

	// ========== 资源详情（根据类型不同使用不同字段） ==========
	// 数据表/逻辑视图
	DatabaseId   string `json:"databaseId" gorm:"column:database_id;size:50"`
	TblName      string `json:"tableName" gorm:"column:table_name;size:128"`
	TableComment string `json:"tableComment" gorm:"column:table_comment;size:500"`

	// API接口
	ApiId       string `json:"apiId" gorm:"column:api_id;size:50"`
	ApiUrl      string `json:"apiUrl" gorm:"column:api_url;size:500"`
	ApiMethod   string `json:"apiMethod" gorm:"column:api_method;size:10"`
	ApiProtocol string `json:"apiProtocol" gorm:"column:api_protocol;size:20"`

	// 文件
	FileId   string `json:"fileId" gorm:"column:file_id;size:50"`
	FileName string `json:"fileName" gorm:"column:file_name;size:500"`
	FileUrl  string `json:"fileUrl" gorm:"column:file_url;size:500"`
	FileSize int64  `json:"fileSize" gorm:"column:file_size;default:0"`

	// ========== 更多信息 ==========
	Source      string `json:"source" gorm:"column:source;size:20;default:manual"` // 来源: manual, sync
	SortOrder   int    `json:"sortOrder" gorm:"column:sort_order;default:0"`

	// ========== 时间戳 ==========
	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

// TableName 指定表名
func (DataResource) TableName() string {
	return "t_data_resource"
}

// IsLogicalView 是否为逻辑视图
func (r *DataResource) IsLogicalView() bool {
	return r.ResourceType == ResourceTypeLogicalView
}

// IsDatabaseTable 是否为数据表
func (r *DataResource) IsDatabaseTable() bool {
	return r.ResourceType == ResourceTypeDatabaseTable
}

// IsApiInterface 是否为API接口
func (r *DataResource) IsApiInterface() bool {
	return r.ResourceType == ResourceTypeApiInterface
}

// IsFile 是否为文件
func (r *DataResource) IsFile() bool {
	return r.ResourceType == ResourceTypeFile
}
