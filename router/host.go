package router

import (
	"github.com/gin-gonic/gin"
	"gofly/api"
)

func InitHostRoutes() {
	RegistRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		hostApi := api.NewHostApi()
		rgAuthUser := rgAuth.Group("host")
		{
			rgAuthUser.POST("/shutdown", hostApi.Shutdown)
		}
	})
}
