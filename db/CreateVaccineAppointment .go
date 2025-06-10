// db/vaccine.go
package db

import (
	"clinic-backend/models"
	"log"
)

// GetVaccineCountByDate 回傳某日已預約人數
func GetVaccineCountByDate(date string) int {
	var count int64
	if err := DB.
		Model(&models.VaccineAppointment{}).
		Where("vaccine_date = ?", date).
		Count(&count).Error; err != nil {
		log.Printf("Error retrieving vaccine appointment count: %v", err)
		return 0
	}
	return int(count)
}

// CreateVaccineAppointment 創建疫苗預約
func CreateVaccineAppointment(appointment models.VaccineAppointment) error {
	if err := DB.Create(&appointment).Error; err != nil {
		log.Printf("Error creating vaccine appointment: %v", err)
		return err
	}
	return nil
}
