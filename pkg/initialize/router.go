package initialize

import (
	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

import (
	_ "github.com/dubbogo/pixiu-admin/docs"
	"github.com/dubbogo/pixiu-admin/pkg/controller/account"
	"github.com/dubbogo/pixiu-admin/pkg/controller/auth"
	"github.com/dubbogo/pixiu-admin/pkg/controller/configInfo"
)

// Routers init router
func Routers() *gin.Engine {
	var router = gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Guest router
	router.POST("/login", account.Login)
	router.POST("/register", account.Register)

	// auth router
	taR := router.Group("/", auth.JWTAuth())

	// The following router needs to check the token
	{
		// user router
		taR.POST("/user/logout", account.Logout)
		taR.POST("/user/password/edit", account.EditPassword)
		taR.POST("/user/getInfo", account.GetUserInfo)
		taR.POST("/user/getUserRole", account.GetUserRole)
		taR.POST("/user/checkIsAdmin", account.CheckUserIsAdmin)

		taR.GET("/config/api/base", configInfo.GetBaseInfo)
		taR.POST("/config/api/base/", configInfo.SetBaseInfo)
		taR.PUT("/config/api/base/", configInfo.SetBaseInfo)

		taR.GET("/config/api/resource/list", configInfo.GetResourceList)
		taR.GET("/config/api/resource/detail", configInfo.GetResourceDetail)
		taR.POST("/config/api/resource", configInfo.CreateResourceInfo)
		taR.PUT("/config/api/resource", configInfo.ModifyResourceInfo)
		taR.DELETE("/config/api/resource", configInfo.DeleteResourceInfo)

		taR.GET("/config/api/cluster/list", configInfo.GetClusterList)
		taR.GET("/config/api/cluster/detail", configInfo.DetailCluster)
		taR.POST("/config/api/cluster", configInfo.UpdateCluster)
		taR.PUT("/config/api/cluster", configInfo.CreateCluster)
		taR.DELETE("/config/api/cluster", configInfo.DeleteCluster)

		taR.GET("/config/api/listener/list", configInfo.GetListenerList)
		taR.GET("/config/api/listener/detail", configInfo.DetailListener)
		taR.POST("/config/api/listener", configInfo.UpdateListener)
		taR.PUT("/config/api/listener", configInfo.CreateListener)
		taR.DELETE("/config/api/listener", configInfo.DeleteListener)

		taR.GET("/config/api/resource/method/list", configInfo.GetMethodList)
		taR.GET("/config/api/resource/method/detail", configInfo.GetMethodDetail)
		taR.POST("/config/api/resource/method", configInfo.CreateMethodInfo)
		taR.PUT("/config/api/resource/method", configInfo.ModifyMethodInfo)
		taR.DELETE("/config/api/resource/method", configInfo.DeleteMethodInfo)

		// taR.GET("/config/api/plugin_group/list", configInfo.GetPluginGroupList)
		// taR.GET("/config/api/plugin_group/detail", configInfo.GetPluginGroupDetail)
		// taR.POST("/config/api/plugin_group", configInfo.CreatePluginGroup)
		// taR.PUT("/config/api/plugin_group", configInfo.ModifyPluginGroup)
		// taR.DELETE("/config/api/plugin_group", configInfo.DeletePluginGroup)

		// taR.GET("/config/api/plugin/ratelimit", configInfo.GetPluginRatelimitDetail)
		// taR.POST("/config/api/plugin/ratelimit", configInfo.CreatePluginRatelimit)
		// taR.PUT("/config/api/plugin/ratelimit", configInfo.ModifyPluginRatelimit)
		// taR.DELETE("/config/api/plugin/ratelimit", configInfo.DeletePluginRatelimit)

		// Which request method to choose, Temporarily choose put method
		taR.PUT("/config/api/resource/publish", configInfo.BatchReleaseResource)
		taR.PUT("/config/api/resource/method/publish", configInfo.BatchReleaseMethod)
		taR.PUT("/config/api/plugin_group/publish", configInfo.BatchReleasePluginGroup)
		// taR.PUT("/config/api/plugin/ratelimit/publish", configInfo.BatchReleasePluginRatelimit)
	}

	return router
}
