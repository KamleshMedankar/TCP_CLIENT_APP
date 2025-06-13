package controllers

import (
	"clientgo/db"
	"clientgo/models"
	//"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// POST /generate
func GenerateRecordsHandler(c *gin.Context) {
	var req models.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Count <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	batchSize := 10000

	for i := 1; i <= req.Count; i += batchSize {
		end := i + batchSize - 1
		if end > req.Count {
			end = req.Count
		}

		pipe := db.Rdb.Pipeline()

		for j := i; j <= end; j++ {
			key := "record:" + strconv.Itoa(j)
			record := models.TenantData{
				ID:     j,
				Name:   fmt.Sprintf("User%d", j),
				Phone:  fmt.Sprintf("90000000%04d", j),
				Status: "pending",
			}

			jsonData, err := json.Marshal(record)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("JSON marshal failed at record %d", j)})
				return
			}

			pipe.Set(key, jsonData, 60*time.Minute)
		}

		if _, err := pipe.Exec(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Pipeline error at batch starting %d: %v", i, err),
			})
			return
		}
	}

	

	c.JSON(http.StatusOK, gin.H{
		"message":  fmt.Sprintf("%d records stored successfully", req.Count),
	})
}

// GET /count
func CountRecordsHandler(c *gin.Context) {
	
	var (
		cursor uint64
		count  int
	)

	for {
		keys, newCursor, err := db.Rdb.Scan(cursor, "record:*", 1000).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		count += len(keys)
		cursor = newCursor
		if cursor == 0 {
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"total_records": count})
}
