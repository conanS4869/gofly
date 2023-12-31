package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"gofly/global"
	"gofly/middleware"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type IFnRegistRoute = func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup)

var (
	gfnRoutes []IFnRegistRoute // 存放所有模块注册路由的回调函数
)

func RegistRoute(fn IFnRegistRoute) {
	if fn == nil {
		return
	}
	gfnRoutes = append(gfnRoutes, fn)
}

func InitRouter() {
	// = 创建监听ctrl + c, 应用退出信号的上下文
	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelCtx()

	r := gin.Default()
	r.Use(middleware.Cors())

	rgPublic := r.Group("/api/v1/public")
	rgAuth := r.Group("/api/v1")
	rgAuth.Use(middleware.Auth())

	// 初始基础平台的路由
	initBasePlatformRoutes()

	// 开始注册系统各模块对应的路由信息
	for _, fnRegistRoute := range gfnRoutes {
		fnRegistRoute(rgPublic, rgAuth)
	}

	stPort := viper.GetString("server.port")
	if stPort == "" {
		stPort = "8999"
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", stPort),
		Handler: r,
	}
	// = 启动一个goroutine来开启web服务, 避免主线程的信号监听被阻塞
	go func() {
		global.Logger.Info(fmt.Sprintf("Start Listen: %s", stPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Error(fmt.Sprintf("Start Server Error: %s", err.Error()))
			return
		}
	}()

	// = 等待停止服务的信号被触发
	<-ctx.Done()
	// = 关闭Server， 5秒内未完成清理动作会直接退出应用
	ctx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctx); err != nil {
		global.Logger.Error(fmt.Sprintf("Stop Server Error: %s", err.Error()))
		return
	}
	global.Logger.Info("Stop Server Success")
}

func initBasePlatformRoutes() {
	InitUserRoutes()
	InitHostRoutes()
}

// ! 注册自定义验证器
func registCustValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("first_is_a", func(fl validator.FieldLevel) bool {
			if value, ok := fl.Field().Interface().(string); ok {
				if value != "" && 0 == strings.Index(value, "a") {
					return true
				}
			}
			return false
		})
	}
}
