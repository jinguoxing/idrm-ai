package types

// ========== 菜单管理相关类型定义 ==========

// CreateMenuReq 创建菜单请求
type CreateMenuReq struct {
	ParentId      int64  `json:"parentId"`      // 父级菜单ID，0表示根菜单
	Name          string `json:"name"`          // 菜单名称
	RoutePath     string `json:"routePath"`     // 路由路径
	ComponentPath string `json:"componentPath"` // 组件路径
	Icon          string `json:"icon"`          // 菜单图标
	SortOrder     int    `json:"sortOrder"`     // 排序号
	PermTag       string `json:"permTag"`       // 权限标识
	Status        int    `json:"status"`        // 状态：1-启用，0-禁用
}

// CreateMenuResp 创建菜单响应
type CreateMenuResp struct {
	Id   int64  `json:"id"`   // 菜单ID
	Name string `json:"name"` // 菜单名称
}

// UpdateMenuReq 更新菜单请求
type UpdateMenuReq struct {
	Id            int64  `json:"id"`            // 菜单ID
	ParentId      int64  `json:"parentId"`      // 父级菜单ID
	Name          string `json:"name"`          // 菜单名称
	RoutePath     string `json:"routePath"`     // 路由路径
	ComponentPath string `json:"componentPath"` // 组件路径
	Icon          string `json:"icon"`          // 菜单图标
	SortOrder     int    `json:"sortOrder"`     // 排序号
	PermTag       string `json:"permTag"`       // 权限标识
	Status        int    `json:"status"`        // 状态
}

// UpdateMenuResp 更新菜单响应
type UpdateMenuResp struct {
	Id   int64  `json:"id"`   // 菜单ID
	Name string `json:"name"` // 菜单名称
}

// DeleteMenuReq 删除菜单请求
type DeleteMenuReq struct {
	Id int64 `json:"id"` // 菜单ID
}

// DeleteMenuResp 删除菜单响应
type DeleteMenuResp struct {
	Id int64 `json:"id"` // 菜单ID
}

// GetMenuReq 获取菜单详情请求
type GetMenuReq struct {
	Id int64 `json:"id"` // 菜单ID
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
	Status        int    `json:"status"`         // 状态
	CreatedAt     string `json:"createdAt"`      // 创建时间
	UpdatedAt     string `json:"updatedAt"`      // 更新时间
}

// GetMenuTreeReq 获取菜单树请求
type GetMenuTreeReq struct {
	Status int `json:"status"` // 状态筛选
}

// MenuTreeNode 菜单树节点
type MenuTreeNode struct {
	Id            int64          `json:"id"`             // 菜单ID
	ParentId      int64          `json:"parentId"`       // 父级菜单ID
	Name          string         `json:"name"`           // 菜单名称
	RoutePath     string         `json:"routePath"`      // 路由路径
	ComponentPath string         `json:"componentPath"`  // 组件路径
	Icon          string         `json:"icon"`           // 菜单图标
	SortOrder     int            `json:"sortOrder"`      // 排序号
	PermTag       string         `json:"permTag"`        // 权限标识
	Status        int            `json:"status"`         // 状态
	Children      []MenuTreeNode `json:"children"`       // 子菜单
}

// GetMenuTreeResp 获取菜单树响应
type GetMenuTreeResp struct {
	Tree []MenuTreeNode `json:"tree"` // 菜单树
}

// ListMenusReq 获取菜单列表请求
type ListMenusReq struct {
	ParentId int64 `json:"parentId"` // 父级菜单ID
	Status   int   `json:"status"`   // 状态筛选
	Page     int   `json:"page"`     // 页码
	PageSize int   `json:"pageSize"` // 每页数量
}

// ListMenusResp 获取菜单列表响应
type ListMenusResp struct {
	Items    []GetMenuResp `json:"items"`    // 菜单列表
	Total    int64         `json:"total"`    // 总数
	Page     int           `json:"page"`     // 当前页
	PageSize int           `json:"pageSize"` // 每页数量
}
