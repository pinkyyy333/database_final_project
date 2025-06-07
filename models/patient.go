package models

type Patient struct {
	PatientID     string `gorm:"primaryKey;column:patient_id" json:"patient_id"`
	PatientName   string `json:"patient_name"`
	PatientGender string `json:"patient_gender"`
	PatientBirth  string `json:"patient_birth"`
	PatientPhone  string `json:"patient_phone"`
	Password      string `json:"-"` // 不回傳給前端
}
