package router

import (
	"dns-check/server/controller"
	"dns-check/server/middleware/adminJwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowAllOrigins:  true,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.POST("/api/login", controller.Login)

	tokenGroup := r.Group("/")
	tokenGroup.Use(adminJwt.JWTAuth())
	tokenGroup.GET("/api/job/list", controller.ListJob)
	tokenGroup.POST("/api/job/add", controller.AddJob)
	tokenGroup.POST("/api/job/start", controller.StartJob)
	tokenGroup.POST("/api/job/paused", controller.PausedJob)
	tokenGroup.POST("/api/job/end", controller.EndJob)
	tokenGroup.POST("/api/job/delete", controller.DeleteJob)
	tokenGroup.POST("/api/job/process", controller.ProcessJob)
	tokenGroup.POST("/api/job/export", controller.ExportJob)
	tokenGroup.POST("/api/domain/list", controller.ListDomain)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "")
	})
	return r
}
