// controllers/slot_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/slots?doctor_id=123&month=2025-06
func GetAppointmentSlots(c *gin.Context) {
	doctorID := c.Query("doctor_id")
	month := c.Query("month")
	if doctorID == "" || month == "" {
		RespondError(c, http.StatusBadRequest, "doctor_id 與 month 為必填參數")
		return
	}

	// TODO: 從 DB 撈真實資料，以下為範例
	slots := []gin.H{
		{"date": month + "-01", "sessions": []string{"morning", "afternoon", "evening"}},
		{"date": month + "-02", "sessions": []string{"morning", "evening"}},
	}

	RespondOK(c, gin.H{
		"doctor_id": doctorID,
		"month":     month,
		"slots":     slots,
	})
}

// PUT /api/v1/slots
func UpdateAppointmentSlots(c *gin.Context) {
	var payload struct {
		DoctorID string                `json:"doctor_id"`
		Slots    []map[string][]string `json:"slots"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		RespondError(c, http.StatusBadRequest, "參數錯誤")
		return
	}
	if payload.DoctorID == "" || len(payload.Slots) == 0 {
		RespondError(c, http.StatusBadRequest, "doctor_id 與 slots 均為必填")
		return
	}

	// TODO: 將 payload.Slots 儲存至 DB，若操作失敗請回報錯誤
	// if err := db.DB.Save(...); err != nil { RespondError(...); return }

	RespondOK(c, gin.H{"message": "可預約時段已更新"})
}
