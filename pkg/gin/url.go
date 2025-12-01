package ginhandler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DBRef *gorm.DB

func SetupRoutes(r *gin.Engine) {

	// Add custom template functions FIRST
	r.SetFuncMap(map[string]interface{}{
		"contains": func(arr []string, val string) bool {
			for _, v := range arr {
				if v == val {
					return true
				}
			}
			return false
		},
	})

	// Register routes
	r.POST("/api/logs", PaginatedfilterLogs)
}

func PrintHelloWorldBeforeRouting(ctx *gin.Context) {
	fmt.Println("Hello World")
	ctx.Next()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")
		//disable caching
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
