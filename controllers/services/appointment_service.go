package services

import (
	"clinic-backend/db"
	"clinic-backend/models"
)

// AppointmentFilter 用於篩選預約列表
type AppointmentFilter struct {
	Date     string
	DeptID   string
	DoctorID string
	Status   string
}

// AppointmentService 提供預約相關的商業邏輯
type AppointmentService struct{}

// NewAppointmentService 建立 AppointmentService 實例
func NewAppointmentService() *AppointmentService {
	return &AppointmentService{}
}

// ListAppointments 列出符合條件的預約清單
func (s *AppointmentService) ListAppointments(f AppointmentFilter) ([]models.Appointment, error) {
	var apps []models.Appointment
	query := db.DB
	if f.Date != "" {
		query = query.Where("date = ?", f.Date)
	}
	if f.DeptID != "" {
		query = query.Where("dept_id = ?", f.DeptID)
	}
	if f.DoctorID != "" {
		query = query.Where("doctor_id = ?", f.DoctorID)
	}
	if f.Status != "" {
		query = query.Where("status = ?", f.Status)
	}
	err := query.Find(&apps).Error
	return apps, err
}

// UpdateStatus 更新單筆預約的狀態
func (s *AppointmentService) UpdateStatus(id int, status string) error {
	err := db.DB.Model(&models.Appointment{}).
		Where("id = ?", id).
		Update("status", status).Error
	return err
}

// DeleteAppointment 刪除預約
func (s *AppointmentService) DeleteAppointment(id int) error {
	err := db.DB.Delete(&models.Appointment{}, id).Error
	return err
}

// AssignSubstitute 為請假醫師的當日已預約病人指派替代醫師
func (s *AppointmentService) AssignSubstitute(absentID, substituteID int, date string) error {
	err := db.DB.Model(&models.Appointment{}).
		Where("doctor_id = ? AND date = ? AND status = ?", absentID, date, "booked").
		Update("doctor_id", substituteID).Error
	return err
}

// GenerateReport 產生管理員報表 (示範回傳 interface，可依需求擴充)
func (s *AppointmentService) GenerateReport(reportType, month string) (interface{}, error) {
	// TODO: 根據 reportType 與 month 組合查詢並回傳結構化報表資料
	return nil, nil
}
