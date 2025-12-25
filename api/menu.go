package main

import (
	"flag"
	"fmt"

	"idrm/api/internal/config"
	menuhandler "idrm/api/internal/handler/sys/menu"
	"idrm/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "api/etc/menu.yaml", "配置文件路径")

func main() {
	flag.Parse()

	// 加载配置
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建服务器
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 初始化ServiceContext（包含数据库连接）
	serverCtx := svc.NewServiceContext(c)

	// 注册路由
	registerRoutes(server, serverCtx)

	// 打印启动信息
	fmt.Printf("╔════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                                                              ║\n")
	fmt.Printf("║  🚀 IDRM 菜单管理服务启动成功                                  ║\n")
	fmt.Printf("║                                                              ║\n")
	fmt.Printf("║  服务地址: http://%s:%d                                   \n", c.Host, c.Port)
	fmt.Printf("║  API文档: http://%s:%d/swagger                          \n", c.Host, c.Port)
	fmt.Printf("║                                                              ║\n")
	fmt.Printf("║  可用端点:                                                    ║\n")
	fmt.Printf("║    POST   /api/v1/sys/menu              - 创建菜单          ║\n")
	fmt.Printf("║    PUT    /api/v1/sys/menu/:id          - 更新菜单          ║\n")
	fmt.Printf("║    DELETE /api/v1/sys/menu/:id          - 删除菜单          ║\n")
	fmt.Printf("║    GET    /api/v1/sys/menu/:id          - 获取菜单详情       ║\n")
	fmt.Printf("║    GET    /api/v1/sys/menu/tree         - 获取菜单树        ║\n")
	fmt.Printf("║    GET    /api/v1/sys/menus             - 获取菜单列表       ║\n")
	fmt.Printf("║                                                              ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════╝\n")

	// 启动服务
	server.Start()
}

// registerRoutes 注册所有路由
func registerRoutes(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 注册菜单管理路由
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  "POST",
				Path:    "/api/v1/sys/menu",
				Handler: menuhandler.NewCreateMenuHandler(serverCtx).CreateMenu,
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  "PUT",
				Path:    "/api/v1/sys/menu/:id",
				Handler: menuhandler.NewUpdateMenuHandler(serverCtx).UpdateMenu,
			},
			{
				Method:  "DELETE",
				Path:    "/api/v1/sys/menu/:id",
				Handler: menuhandler.NewDeleteMenuHandler(serverCtx).DeleteMenu,
			},
			{
				Method:  "GET",
				Path:    "/api/v1/sys/menu/:id",
				Handler: menuhandler.NewGetMenuHandler(serverCtx).GetMenu,
			},
			{
				Method:  "GET",
				Path:    "/api/v1/sys/menu/tree",
				Handler: menuhandler.NewGetMenuTreeHandler(serverCtx).GetMenuTree,
			},
			{
				Method:  "GET",
				Path:    "/api/v1/sys/menus",
				Handler: menuhandler.NewListMenusHandler(serverCtx).ListMenus,
			},
		},
	)
}
