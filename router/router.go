package router

import (
	"apiserver/handler/sd"
	"apiserver/handler/user"
	"apiserver/router/middleware"
	_ "github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Load 加载中间件、路由、处理程序
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery()) //注册全局中间件
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)

	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "The incorrect API route.",
		})
		//c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// pprof router
	//pprof.Register(g)

	// 用于身份验证功能的api
	g.POST("/login", user.Login)

	//用户路由设置
	u := g.Group("/v1/user")
	u.Use(middleware.AuthMiddleware())    //路由中间件jwt验证
	{
		u.GET("",user.List)			 	 //获取用户列表
		u.GET("/:id",user.Get) 		 	 //获取用户信息
		u.PUT("/:id",user.Update)	 	 //编辑用户信息
		u.DELETE("/:id",user.Delete)       //删除用户
		u.POST("", user.Create)  //创建用户
	}


	// 运行状况检查处理程序 路由
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	return g
}
