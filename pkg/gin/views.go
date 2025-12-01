package ginhandler

import (
	"fmt"
	databasemodel "loggenerator/pkg/database_model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ShowFilterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Level":     []string{},
		"Component": []string{},
		"Host":      []string{},
		"RequestID": "",
		"Timestamp": "",
	})
}

func ShowAllLogs(c *gin.Context) {
	entries, err := databasemodel.GetAllLogs(DBRef) //empty filter to get all logs
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"entries": entries[1:10000],
	})
}
func PaginatedfilterLogs(c *gin.Context) {
	// 1. Read pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))
	offset := page * pageSize

	// 2. Read JSON body EXACTLY matching frontend keys
	var body struct {
		Level     []string `json:"level"`
		Component []string `json:"component"`
		Host      []string `json:"host"`
		RequestId string   `json:"requestId"`
		StartTime string   `json:"startTime"`
		EndTime   string   `json:"endTime"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println("JSON BIND ERROR:", err)
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	entries, err := databasemodel.FilterLogs(
		DBRef,
		body.Level,
		body.Component,
		body.Host,
		body.RequestId,
		body.StartTime,
		body.EndTime,
	)

	if err != nil {
		fmt.Println("DB FILTER ERROR:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	total := len(entries)

	// 4. Manual pagination
	start := offset
	end := offset + pageSize

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pageEntries := entries[start:end]

	// 5. Return paginated response
	c.JSON(200, gin.H{
		"entries": pageEntries,
		"count":   total,
	})
}
