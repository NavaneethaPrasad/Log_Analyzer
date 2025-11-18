package ginhandler

import (
	databasemodel "loggenerator/pkg/database_model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DBRef *gorm.DB

func ShowFilterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Level":     []string{},
		"Component": []string{},
		"Host":      []string{},
		"RequestID": "",
		"Timestamp": "",
	})
}

// func RunFilter(c *gin.Context) {
// 	rawFilter := c.PostForm("filter")

// 	if strings.TrimSpace(rawFilter) == "" {
// 		c.HTML(http.StatusOK, "index.html", gin.H{
// 			"Error": "Filter cannot be empty",
// 		})
// 		return
// 	}

// 	parts := databasemodel.SplitUserFilter(rawFilter)

// 	entries, err := databasemodel.Query(DBRef, parts)

// 	if err != nil {
// 		c.HTML(http.StatusOK, "index.html", gin.H{
// 			"Error": err.Error(),
// 		})
// 		return
// 	}

// 	// render same page but with results
// 	c.HTML(http.StatusOK, "index.html", gin.H{
// 		"Entries": entries,
// 	})
// }

func RunFilter(c *gin.Context) {

	// Multi-select checkboxes
	levels := c.PostFormArray("level")
	components := c.PostFormArray("component")
	hosts := c.PostFormArray("host")

	// Textboxes
	requestID := c.PostForm("request_id")
	timestamp := c.PostForm("timestamp") // e.g., "> 2025-11-17 10:00:00"

	// Call database function
	entries, err := databasemodel.FilterLogs(DBRef, levels, components, hosts, requestID, timestamp)
	if err != nil {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Error":     err.Error(),
			"Level":     levels,
			"Component": components,
			"Host":      hosts,
			"RequestID": requestID,
			"Timestamp": timestamp,
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Entries":   entries,
		"Count":     len(entries),
		"Level":     levels,
		"Component": components,
		"Host":      hosts,
		"RequestID": requestID,
		"Timestamp": timestamp,
	})
}
