package menu

import "time"

// Menu 菜单数据结构
type Menu struct {
	Id            int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ParentId      int64     `json:"parentId" gorm:"column:parent_id;not null;default:0"`
	Name          string    `json:"name" gorm:"column:name;type:varchar(50);not null"`
	RoutePath     string    `json:"routePath" gorm:"column:route_path;type:varchar(200)"`
	ComponentPath string    `json:"componentPath" gorm:"column:component_path;type:varchar(200)"`
	Icon          string    `json:"icon" gorm:"column:icon;type:varchar(100)"`
	SortOrder     int       `json:"sortOrder" gorm:"column:sort_order;not null;default:0"`
	PermTag       string    `json:"permTag" gorm:"column:perm_tag;type:varchar(100)"`
	Status        int       `json:"status" gorm:"column:status;not null;default:1"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty" gorm:"column:deleted_at"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "sys_menu"
}
