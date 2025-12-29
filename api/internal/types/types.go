package types

// ========== Data Catalog Types ==========

type ColumnItem struct {
	TechnicalName string `json:"technicalName"`          // 技术名称
	BusinessName  string `json:"businessName,optional"`  // 业务名称
	DataFormat    int8   `json:"dataFormat,optional"`    // 字段类型
	DataLength    int    `json:"dataLength,optional"`    // 字段长度
	DataPrecision int8   `json:"dataPrecision,optional"` // 数据精度
	Description   string `json:"description,optional"`   // 字段描述
	SharedType    int8   `json:"sharedType,optional"`    // 共享属性
	OpenType      int8   `json:"openType,optional"`      // 开放属性
	PrimaryFlag   bool   `json:"primaryFlag,optional"`   // 是否主键
	NullFlag      bool   `json:"nullFlag,optional"`      // 是否为空
	Index         int    `json:"index"`                  // 顺序
}

type DeleteDraftReq struct {
	Id string `path:"id"` // 草稿ID
}

type DeleteDraftResp struct {
	Id string `json:"id"` // 删除的草稿ID
}

type GetDraftReq struct {
	Id string `path:"id"` // 草稿ID
}

type GetDraftResp struct {
	CatalogId        string          `json:"catalogId"`
	Title            string          `json:"title"`
	Description      string          `json:"description"`
	SourceDeptId     string          `json:"sourceDeptId"`
	SourceDeptName   string          `json:"sourceDeptName"`
	DataDomain       int8            `json:"dataDomain"`
	PublishStatus    string          `json:"publishStatus"` // 发布状态
	DraftId          string          `json:"draftId"`       // 草稿ID
	IsDraftCopy      bool            `json:"isDraftCopy"`   // 是否是草稿副本
	SharedType       int8            `json:"sharedType"`
	SharedCondition  string          `json:"sharedCondition"`
	OpenType         int8            `json:"openType"`
	OpenCondition    string          `json:"openCondition"`
	Columns          []ColumnItem    `json:"columns"`
	MountResources   []MountResource `json:"mountResources"`
	UpdateCycle      int8            `json:"updateCycle"`
	OtherUpdateCycle string          `json:"otherUpdateCycle"`
	SyncMechanism    int8            `json:"syncMechanism"`
	SyncFrequency    string          `json:"syncFrequency"`
	CategoryNodeIds  []string        `json:"categoryNodeIds"`
	CreatedAt        string          `json:"createdAt"`
	UpdatedAt        string          `json:"updatedAt"`
}

type ListDraftsReq struct {
	SourceDeptId string `form:"sourceDeptId,optional"` // 来源部门ID
	Status       string `form:"status,optional"`       // 发布状态
	Page         int    `form:"page,optional,default=1"`
	PageSize     int    `form:"pageSize,optional,default=20"`
}

type ListDraftsResp struct {
	Items    []GetDraftResp `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

type MountResource struct {
	ResourceId   string `json:"resourceId"`   // 资源ID
	ResourceType int8   `json:"resourceType"` // 资源类型：1-逻辑视图 2-接口 3-文件资源
	Name         string `json:"name"`         // 资源名称
}

type SaveDraftReq struct {
	CatalogId        string          `json:"catalogId,optional"`        // 目录ID，空表示创建新草稿
	Title            string          `json:"title"`                     // 目录名称
	Description      string          `json:"description,optional"`      // 资源目录描述
	SourceDeptId     string          `json:"sourceDeptId"`              // 来源部门ID
	DataDomain       int8            `json:"dataDomain,optional"`       // 数据所在领域
	CategoryNodeIds  []string        `json:"categoryNodeIds,optional"`  // 资源属性分类节点ID列表
	SharedType       int8            `json:"sharedType,optional"`       // 共享类型：1-无条件 2-有条件 3-不予共享
	SharedCondition  string          `json:"sharedCondition,optional"`  // 共享条件
	OpenType         int8            `json:"openType,optional"`         // 开放类型：1-向公众开放 2-不向公众开放
	OpenCondition    string          `json:"openCondition,optional"`    // 开放条件
	Columns          []ColumnItem    `json:"columns,optional"`          // 信息项列表
	MountResources   []MountResource `json:"mountResources,optional"`   // 挂接资源列表
	UpdateCycle      int8            `json:"updateCycle,optional"`      // 更新频率
	OtherUpdateCycle string          `json:"otherUpdateCycle,optional"` // 其他更新频率
	SyncMechanism    int8            `json:"syncMechanism,optional"`    // 数据归集机制：1-增量 2-全量
	SyncFrequency    string          `json:"syncFrequency,optional"`    // 数据归集频率
}

type SaveDraftResp struct {
	CatalogId string `json:"catalogId"` // 目录ID
	DraftId   string `json:"draftId"`   // 草稿ID（如果是草稿副本）
	Status    string `json:"status"`    // 当前状态
}

// ========== Menu Types ==========

// CreateMenuReq 创建菜单请求
type CreateMenuReq struct {
	ParentId      int64  `json:"parentId,optional"`      // 父级菜单ID，0表示根菜单
	Name          string `json:"name"`                   // 菜单名称
	RoutePath     string `json:"routePath,optional"`     // 路由路径
	ComponentPath string `json:"componentPath,optional"` // 组件路径
	Icon          string `json:"icon,optional"`          // 菜单图标
	SortOrder     int    `json:"sortOrder,optional"`     // 排序号
	PermTag       string `json:"permTag,optional"`       // 权限标识
	Status        int8   `json:"status,optional"`        // 状态：1-启用，0-禁用
}

// CreateMenuResp 创建菜单响应
type CreateMenuResp struct {
	Id int64 `json:"id"` // 菜单ID
}

// UpdateMenuReq 更新菜单请求
type UpdateMenuReq struct {
	Id            int64  `path:"id"`                    // 菜单ID
	ParentId      int64  `json:"parentId,optional"`      // 父级菜单ID
	Name          string `json:"name"`                   // 菜单名称
	RoutePath     string `json:"routePath,optional"`     // 路由路径
	ComponentPath string `json:"componentPath,optional"` // 组件路径
	Icon          string `json:"icon,optional"`          // 菜单图标
	SortOrder     int    `json:"sortOrder,optional"`     // 排序号
	PermTag       string `json:"permTag,optional"`       // 权限标识
	Status        int8   `json:"status,optional"`        // 状态
}

// UpdateMenuResp 更新菜单响应
type UpdateMenuResp struct {
	Id int64 `json:"id"` // 菜单ID
}

// DeleteMenuReq 删除菜单请求
type DeleteMenuReq struct {
	Id int64 `path:"id"` // 菜单ID
}

// DeleteMenuResp 删除菜单响应
type DeleteMenuResp struct {
	Id int64 `json:"id"` // 菜单ID
}

// GetMenuReq 获取菜单详情请求
type GetMenuReq struct {
	Id int64 `path:"id"` // 菜单ID
}

// GetMenuResp 获取菜单详情响应
type GetMenuResp struct {
	Id            int64  `json:"id"`             // 菜单ID
	ParentId      int64  `json:"parentId"`       // 父级菜单ID
	Name          string `json:"name"`           // 菜单名称
	RoutePath     string `json:"routePath"`      // 路由路径
	ComponentPath string `json:"componentPath"`  // 组件路径
	Icon          string `json:"icon"`           // 菜单图标
	SortOrder     int    `json:"sortOrder"`      // 排序号
	PermTag       string `json:"permTag"`        // 权限标识
	Status        int8   `json:"status"`         // 状态
	CreatedAt     string `json:"createdAt"`      // 创建时间
	UpdatedAt     string `json:"updatedAt"`      // 更新时间
}

// GetMenuTreeReq 获取菜单树请求
type GetMenuTreeReq struct {
	Status int8 `form:"status,optional"` // 状态筛选
}

// MenuTreeNode 菜单树节点
type MenuTreeNode struct {
	Id            int64           `json:"id"`             // 菜单ID
	ParentId      int64           `json:"parentId"`       // 父级菜单ID
	Name          string          `json:"name"`           // 菜单名称
	RoutePath     string          `json:"routePath"`      // 路由路径
	ComponentPath string          `json:"componentPath"`  // 组件路径
	Icon          string          `json:"icon"`           // 菜单图标
	SortOrder     int             `json:"sortOrder"`      // 排序号
	PermTag       string          `json:"permTag"`        // 权限标识
	Status        int8            `json:"status"`         // 状态
	Children      []*MenuTreeNode `json:"children,optional"` // 子菜单
}

// GetMenuTreeResp 获取菜单树响应
type GetMenuTreeResp struct {
	Tree []*MenuTreeNode `json:"tree"` // 菜单树
}

// ListMenusReq 获取菜单列表请求
type ListMenusReq struct {
	Status   int8  `form:"status,optional"`   // 状态筛选
	ParentId int64 `form:"parentId,optional"` // 父级菜单ID
	Page     int   `form:"page,optional,default=1"`
	PageSize int   `form:"pageSize,optional,default=20"`
}

// ListMenusResp 获取菜单列表响应
type ListMenusResp struct {
	Items    []*GetMenuResp `json:"items"`    // 菜单列表
	Total    int64          `json:"total"`    // 总数
	Page     int            `json:"page"`     // 当前页
	PageSize int            `json:"pageSize"` // 每页数量
}
