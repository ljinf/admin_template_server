package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"poem_server_admin/internal/cache"
	handler "poem_server_admin/internal/handler/sys"
	"poem_server_admin/internal/middleware"
	"poem_server_admin/pkg/jwt"
	"poem_server_admin/pkg/log"
	"poem_server_admin/pkg/server/http"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	accountHandler handler.AccountHandler,
	accountCache cache.AccountCache,
	userHandler handler.UserHandler,
	menuHandler handler.MenuHandler,
	dictHandler handler.DictHandler,
	roleHandler handler.RoleHandler,

) *http.Server {
	gin.SetMode(gin.DebugMode)
	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	s.Use(
		middleware.CORSMiddleware(),
		//middleware.ResponseLogMiddleware(logger),
		//middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)

	/*
		系统相关
	*/
	r := s.Group("/sys")
	accountGroup := r.Group("/account")
	{
		accountGroup.POST("/login", accountHandler.Login)
		accountGroup.POST("/logout", accountHandler.Logout)
		accountGroup.GET("/getInfo", middleware.AuthMiddleware(jwt), accountHandler.GetInfo)
	}

	userGroup := r.Group("/user", middleware.AuthMiddleware(jwt))
	{
		userGroup.POST("/create", middleware.CheckPerm("system:user:add", logger, accountCache), userHandler.CreateUser)
		//用户列表
		userGroup.GET("/list", userHandler.UserList)
		//获取用户信息
		userGroup.GET("/info/:userId", userHandler.GetUserById)
		userGroup.GET("/info", userHandler.GetUserById)
		userGroup.GET("/profile", userHandler.Profile)

		userGroup.PUT("/updateInfo", middleware.CheckPerm("system:user:edit", logger, accountCache), userHandler.UpdateUser)
		userGroup.DELETE("/del/:userId", middleware.CheckPerm("system:user:remove", logger, accountCache), userHandler.DelUser)
		//根据用户编号获取授权角色
		userGroup.GET("/authRole/:userId", userHandler.AuthRoleById)
		//用户授权角色
		userGroup.POST("/authRole", userHandler.AuthRole)
		//获取部门树列表
		userGroup.GET("/deptTree", userHandler.DeptTree)

		userGroup.PUT("/profile/update", userHandler.UpdateProfile)
	}

	menuGroup := r.Group("/menu", middleware.AuthMiddleware(jwt))
	{
		menuGroup.GET("/getRouters", menuHandler.GetRouters)
		menuGroup.GET("/info/:menuId", menuHandler.GetInfo)
		menuGroup.GET("/list", menuHandler.MenuList)
		menuGroup.POST("/add", middleware.CheckPerm("system:menu:add", logger, accountCache), menuHandler.AddMenu)
		menuGroup.PUT("/update", middleware.CheckPerm("system:menu:edit", logger, accountCache), menuHandler.UpdateMenu)
		menuGroup.DELETE("/del/:menuId", middleware.CheckPerm("system:menu:remove", logger, accountCache), menuHandler.DelMenu)
		menuGroup.GET("/menuTreeSelect", menuHandler.MenuTreeSelect)                 //所有菜单
		menuGroup.GET("/roleMenuTreeSelect/:roleId", menuHandler.RoleMenuTreeSelect) //角色拥有的菜单
	}

	roleGroup := r.Group("/role", middleware.AuthMiddleware(jwt))
	{
		roleGroup.GET("/list", roleHandler.RoleList)
		roleGroup.GET("/info/:roleId", roleHandler.RoleInfo)
		roleGroup.POST("/add", middleware.CheckPerm("system:role:add", logger, accountCache), roleHandler.AddRole)
		roleGroup.PUT("/update", middleware.CheckPerm("system:role:edit", logger, accountCache), roleHandler.UpdateRole)
		roleGroup.DELETE("/del/:roleId", middleware.CheckPerm("system:role:remove", logger, accountCache), roleHandler.DelRole)

	}

	dictGroup := r.Group("/dict", middleware.AuthMiddleware(jwt))
	{
		dictGroup.GET("/data/type/:dictType", dictHandler.GetDictDataType)
	}

	return s
}
