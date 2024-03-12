package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首页
// @Accept json
// @Produce json
// @Success 200 {string} Welcome!!!
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome Gin Chat!!!",
	})
}
