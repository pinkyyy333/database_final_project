// controllers/doctor_controller.go
package controllers

import (
	"database/sql" // 因為我們還是用 sql.NullString、sql.NullInt64 來處理 NULL 值
	"net/http"

	"clinic-backend/db"

	"github.com/gin-gonic/gin"
)

type Appointment struct {
	AppointmentID   int    `json:"appointment_id"`
	PatientName     string `json:"patient_name"`
	AppointmentTime string `json:"appointment_time"`
	Status          string `json:"status"`
}

type Record struct {
	AppointmentTime string `json:"appointment_time"`
	Status          string `json:"status"`
	Comment         string `json:"comment"`
	Rating          int    `json:"rating"`
}

type FeedbackStats struct {
	AverageRating float64 `json:"average_rating"`
	TotalCount    int     `json:"total_count"`
}

// 排程與預約管理
func GetDoctorSchedule(c *gin.Context) {
	doctorID := c.Param("id")

	rows, err := db.DB.Raw(`
		SELECT a.appointment_id, p.patient_name, a.appointment_time, a.status
		FROM appointments a
		JOIN patients p ON a.patient_id = p.patient_id
		WHERE a.doctor_id = ?
		ORDER BY a.appointment_time ASC`, doctorID).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "DB query error"})
		return
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var appt Appointment
		if scanErr := rows.Scan(
			&appt.AppointmentID,
			&appt.PatientName,
			&appt.AppointmentTime,
			&appt.Status,
		); scanErr == nil {
			appointments = append(appointments, appt)
		}
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "data": appointments})
}

// 病歷與回饋查詢
func GetPatientRecords(c *gin.Context) {
	patientID := c.Param("patient_id")
	doctorID := c.Param("id")

	rows, err := db.DB.Raw(`
		SELECT a.appointment_time, a.status, f.patient_comment, f.feedback_rating
		FROM appointments a
		LEFT JOIN feedback f ON a.appointment_id = f.appointment_id
		WHERE a.patient_id = ? AND a.doctor_id = ?`, patientID, doctorID).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "DB error"})
		return
	}
	defer rows.Close()

	var records []Record
	var totalRating int
	var ratingCount int
	for rows.Next() {
		var r Record
		var comment sql.NullString
		var rating sql.NullInt64
		if scanErr := rows.Scan(&r.AppointmentTime, &r.Status, &comment, &rating); scanErr == nil {
			r.Comment = comment.String
			if rating.Valid {
				r.Rating = int(rating.Int64)
				totalRating += r.Rating
				ratingCount++
			}
			records = append(records, r)
		}
	}

	var stats FeedbackStats
	if ratingCount > 0 {
		stats.AverageRating = float64(totalRating) / float64(ratingCount)
		stats.TotalCount = ratingCount
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  records,
		"stats": stats,
	})
}
