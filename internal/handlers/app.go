package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppConfigGetReq struct {
	Key string `form:"key" binding:"required"`
}

type AppConfigSetReq struct {
	Key   string      `json:"key" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
}

func AppConfigGetHandler(c *gin.Context) {
	var req AppConfigGetReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request parameters",
		})
		return
	}

	svc := GetService(c)
	value := svc.GetConfigValue(req.Key)
	if value == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "config item not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"key":   req.Key,
		"value": value,
	})
}

func AppConfigSetHandler(c *gin.Context) {
	var req AppConfigSetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request parameters",
		})
		return
	}

	svc := GetService(c)
	if err := svc.SetConfigValue(req.Key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save config",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "config updated",
		"key":     req.Key,
		"value":   req.Value,
	})
}
