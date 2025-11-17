package ginhandler

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine) {

	r.GET("/", ShowFilterPage)
	r.POST("/search", RunFilter)
}
