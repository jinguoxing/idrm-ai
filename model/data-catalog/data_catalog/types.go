package data_catalog

import "time"

// PublishStatus 发布状态枚举
type PublishStatus string

const (
	StatusUnpublished    PublishStatus = "unpublished"     // 未发布（草稿）
	StatusPubAuditing    PublishStatus = "pub-auditing"    // 发布审核中
	StatusPublished      PublishStatus = "published"       // 已发布
	StatusPubReject      PublishStatus = "pub-reject"      // 发布审核未通过
	StatusChangeAuditing PublishStatus = "change-auditing" // 变更审核中
	StatusChangeReject   PublishStatus = "change-reject"   // 变更审核未通过
)

// OnlineStatus 上线状态枚举
type OnlineStatus string

const (
	OnlineStatusNotline      OnlineStatus = "notline"       // 未上线
	OnlineStatusOnline       OnlineStatus = "online"        // 已上线
	OnlineStatusOffline      OnlineStatus = "offline"       // 已下线
	OnlineStatusUpAuditing   OnlineStatus = "up-auditing"   // 上线审核中
	OnlineStatusDownAuditing OnlineStatus = "down-auditing" // 下线审核中
	OnlineStatusUpReject     OnlineStatus = "up-reject"     // 上线审核未通过
	OnlineStatusDownReject   OnlineStatus = "down-reject"   // 下线审核未通过
)

// DataCatalog 数据资源目录数据结构
// 对应表：t_data_catalog
type DataCatalog struct {
	// ========== 主键和基础信息 ==========
	Id              string       `json:"id" gorm:"column:id;primaryKey;size:20"`
	Code            string       `json:"code" gorm:"column:code;size:50;not null"`
	Title           string       `json:"title" gorm:"column:title;size:500;not null"`

	// ========== 分类信息 ==========
	GroupId         uint64       `json:"groupId" gorm:"column:group_id;default:0"`
	GroupName       string       `json:"groupName" gorm:"column:group_name;size:128"`
	ThemeId         uint64       `json:"themeId" gorm:"column:theme_id"`
	ThemeName       string       `json:"themeName" gorm:"column:theme_name;size:100"`

	// ========== 描述和范围 ==========
	Description     string       `json:"description" gorm:"column:description;size:1000"`
	DataRange       int8         `json:"dataRange" gorm:"column:data_range"`
	DataDomain      int8         `json:"dataDomain" gorm:"column:data_domain"`
	DataLevel       int8         `json:"dataLevel" gorm:"column:data_level"`
	TimeRange       string       `json:"timeRange" gorm:"column:time_range;size:100"`

	// ========== 更新和同步 ==========
	UpdateCycle     int8         `json:"updateCycle" gorm:"column:update_cycle"`
	OtherUpdateCycle string       `json:"otherUpdateCycle" gorm:"column:other_update_cycle;size:100"`
	SyncMechanism   int8         `json:"syncMechanism" gorm:"column:sync_mechanism"`
	SyncFrequency   string       `json:"syncFrequency" gorm:"column:sync_frequency;size:128"`

	// ========== 基础信息分类 ==========
	DataKind        int32        `json:"dataKind" gorm:"column:data_kind;not null"`

	// ========== 共享开放 ==========
	SharedType      int8         `json:"sharedType" gorm:"column:shared_type;not null"`
	SharedCondition string       `json:"sharedCondition" gorm:"column:shared_condition;size:255"`
	ColumnUnshared  bool         `json:"columnUnshared" gorm:"column:column_unshared;not null"`
	OpenType        int8         `json:"openType" gorm:"column:open_type;not null"`
	OpenCondition   string       `json:"openCondition" gorm:"column:open_condition;size:255"`
	SharedMode      int8         `json:"sharedMode" gorm:"column:shared_mode;not null"`

	// ========== 挂接资源统计 ==========
	ViewCount       int          `json:"viewCount" gorm:"column:view_count;default:0"`
	ApiCount        int          `json:"apiCount" gorm:"column:api_count;default:0"`
	FileCount       int          `json:"fileCount" gorm:"column:file_count;default:0"`
	PhysicalDeletion *int8       `json:"physicalDeletion" gorm:"column:physical_deletion"`

	// ========== 流程相关 ==========
	FlowNodeId      string       `json:"flowNodeId" gorm:"column:flow_node_id;size:50"`
	FlowNodeName    string       `json:"flowNodeName" gorm:"column:flow_node_name;size:200"`
	FlowId          string       `json:"flowId" gorm:"column:flow_id;size:50"`
	FlowName        string       `json:"flowName" gorm:"column:flow_name;size:200"`
	FlowVersion     string       `json:"flowVersion" gorm:"column:flow_version;size:10"`

	// ========== 部门和用户 ==========
	DepartmentId    string       `json:"departmentId" gorm:"column:department_id;size:36;not null"`
	SourceDeptId    string       `json:"sourceDeptId" gorm:"column:source_department_id;size:36;not null"`
	OwnerId         string       `json:"ownerId" gorm:"column:owner_id;size:50;not null"`
	OwnerName       string       `json:"ownerName" gorm:"column:owner_name;size:128;not null"`

	// ========== 状态字段 ==========
	PublishStatus   PublishStatus `json:"publishStatus" gorm:"column:publish_status;size:20;default:unpublished"`
	OnlineStatus    OnlineStatus  `json:"onlineStatus" gorm:"column:online_status;size:20;default:notline"`
	AuditType       string       `json:"auditType" gorm:"column:audit_type;size:50;default:unpublished"`
	AuditState      *int8        `json:"auditState" gorm:"column:audit_state"`
	AuditApplySn    uint64       `json:"auditApplySn" gorm:"column:audit_apply_sn;default:0"`
	AuditAdvice     string       `json:"auditAdvice" gorm:"column:audit_advice;type:text"`

	// ========== 草稿相关 ==========
	DraftId         uint64       `json:"draftId" gorm:"column:draft_id;default:0"` // 草稿ID，用于关联草稿副本

	// ========== 标记字段 ==========
	Source          int8         `json:"source" gorm:"column:source;default:1"`
	CurrentVersion  bool         `json:"currentVersion" gorm:"column:current_version;default:1"`
	PublishFlag     bool         `json:"publishFlag" gorm:"column:publish_flag;default:0"`
	IsIndexed       *bool        `json:"isIndexed" gorm:"column:is_indexed"`

	// ========== 更多信息 ==========
	DataClassify    string       `json:"dataClassify" gorm:"column:data_classify;size:50"`
	DataRelatedMatters string    `json:"dataRelatedMatters" gorm:"column:data_related_matters;size:255"`
	BusinessMatters  string       `json:"businessMatters" gorm:"column:business_matters;type:text"`

	// ========== 时间戳 ==========
	CreatedAt       time.Time    `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt       *time.Time   `json:"updatedAt" gorm:"column:updated_at"`
	PublishedAt     *time.Time   `json:"publishedAt" gorm:"column:published_at"`
	OnlineTime      *time.Time   `json:"onlineTime" gorm:"column:online_time"`

	// ========== 更多业务字段 ==========
	ProcDefKey      string       `json:"procDefKey" gorm:"column:proc_def_key;size:128"`
	FlowApplyId     string       `json:"flowApplyId" gorm:"column:flow_apply_id;size:50"`
	ApplyNum        int          `json:"applyNum" gorm:"column:apply_num;default:0"`
	ExploreJobId    string       `json:"exploreJobId" gorm:"column:explore_job_id;size:64"`
	ExploreJobVersion int        `json:"exploreJobVersion" gorm:"column:explore_job_version"`
	IsImport        bool         `json:"isImport" gorm:"column:is_import;default:0"`
}

// TableName 指定表名
func (DataCatalog) TableName() string {
	return "t_data_catalog"
}
