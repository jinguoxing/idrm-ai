package main

import (
	"idrm/api/internal/handler/sys/menu"
	"idrm/api/internal/svc"
	"idrm/api/internal/types"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterMenuRoutes 注册菜单管理相关路由
// 这是一个示例文件，展示如何在main.go中注册路由
func RegisterMenuRoutes(server *rest.Server, serverCtx *svc.ServiceContext) {

	// 菜单管理路由组
	server.AddRoutes(
		[]rest.Route{
			// 创建菜单
			{
				Method:  "POST",
				Path:    "/api/v1/sys/menu",
				Handler: menu.NewCreateMenuHandler(serverCtx).CreateMenu,
			},
		},
		// 可以在这里添加中间件，例如认证中间件
		// rest.WithJwt(serverCtx.Config.JwtAuth),
	)

	server.AddRoutes(
		[]rest.Route{
			// 更新菜单
			{
				Method:  "PUT",
				Path:    "/api/v1/sys/menu/:id",
				Handler: menu.NewUpdateMenuHandler(serverCtx).UpdateMenu,
			},
			// 删除菜单
			{
				Method:  "DELETE",
				Path:    "/api/v1/sys/menu/:id",
				Handler: menu.NewDeleteMenuHandler(serverCtx).DeleteMenu,
			},
			// 获取菜单详情
			{
				Method:  "GET",
				Path:    "/api/v1/sys/menu/:id",
				Handler: menu.NewGetMenuHandler(serverCtx).GetMenu,
			},
			// 获取菜单树
			{
				Method:  "GET",
				Path:    "/api/v1/sys/menu/tree",
				Handler: menu.NewGetMenuTreeHandler(serverCtx).GetMenuTree,
			},
			// 获取菜单列表
			{
				Method:  "GET",
				Path:    "/api/v1/sys/menus",
				Handler: menu.NewListMenusHandler(serverCtx).ListMenus,
			},
		},
	)
}

// 使用示例：
// 在 main.go 中调用：
//
// func main() {
//     c := config.NewConfig()
//     server := rest.MustNewServer(c.RestConf)
//     serverCtx := svc.NewServiceContext(c)
//
//     // 注册菜单路由
//     RegisterMenuRoutes(server, serverCtx)
//
//     server.Start()
// }
