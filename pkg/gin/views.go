package ginhandler

import (
	databasemodel "loggenerator/pkg/database_model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DBRef *gorm.DB

func ShowFilterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func RunFilter(c *gin.Context) {
	rawFilter := c.PostForm("filter")

	if strings.TrimSpace(rawFilter) == "" {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Error": "Filter cannot be empty",
		})
		return
	}

	parts := databasemodel.SplitUserFilter(rawFilter)

	entries, err := databasemodel.Query(DBRef, parts)

	if err != nil {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	// render same page but with results
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Entries": entries,
	})
}
