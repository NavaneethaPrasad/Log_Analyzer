package ginhandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

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
	r.GET("/", ShowFilterPage)
	r.POST("/search", RunFilter)
}
