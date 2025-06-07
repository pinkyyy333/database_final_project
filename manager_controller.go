// controllers/manager_controller.go
package controllers

import (
	"net/http"
	"strconv"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
)

// CreateManager 新增一個管理員
func CreateManager(c *gin.Context) {
	var req models.Manager
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "請求格式錯誤: " + err.Error()})
		return
	}

	// 假設你已經把 GORM.AutoMigrate(&Manager{}) 做過
	if err := db.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "建立失敗: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"error": false, "data": req})
}

// GetManagers 列出所有管理員
func GetManagers(c *gin.Context) {
	var list []models.Manager
	if err := db.DB.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "查詢失敗: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "data": list})
}

// GetManagerByID 依照 ID 取得單一管理員
func GetManagerByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "ID 格式錯誤"})
		return
	}

	var m models.Manager
	// GORM First() 找不到會回傳 ErrRecordNotFound
	if err := db.DB.First(&m, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "message": "找不到該管理員"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "data": m})
}

// UpdateManager 修改管理員資料
func UpdateManager(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "ID 格式錯誤"})
		return
	}

	var req models.Manager
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "請求格式錯誤: " + err.Error()})
		return
	}

	// 先找出這筆紀錄
	var existing models.Manager
	if err := db.DB.First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "message": "找不到該管理員"})
		return
	}

	// 更新欄位（範例只更新 Account, Password；你也可以更新其他欄位）
	existing.Account = req.Account
	existing.Password = req.Password

	if err := db.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "更新失敗: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "data": existing})
}

// DeleteManager 刪除管理員
func DeleteManager(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "ID 格式錯誤"})
		return
	}

	if err := db.DB.Delete(&models.Manager{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "刪除失敗: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "刪除成功"})
}
