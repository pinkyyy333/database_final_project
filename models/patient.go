package models

import (
	"gorm.io/datatypes"
)

// Patient 代表病患基本資料，新增地址、緊急聯絡、過敏及病史欄位
type Patient struct {
	PatientID     string `gorm:"primaryKey;column:patient_id" json:"patient_id"`
	PatientName   string `json:"patient_name"`
	PatientGender string `json:"patient_gender"`
	PatientBirth  string `json:"patient_birth"`
	PatientPhone  string `json:"patient_phone"`
	Password      string `json:"-"`

	// 以下都是前端表單新加
	Address           string `json:"address" gorm:"column:address"`
	EmergencyName     string `json:"emergency_name" gorm:"column:emergency_name"`
	EmergencyPhone    string `json:"emergency_phone" gorm:"column:emergency_phone"`
	EmergencyRelation string `json:"emergency_relation" gorm:"column:emergency_relation"`

	DrugAllergy    datatypes.JSON `json:"drug_allergy" gorm:"column:drug_allergy;type:json"`
	FoodAllergy    datatypes.JSON `json:"food_allergy" gorm:"column:food_allergy;type:json"`
	MedicalHistory datatypes.JSON`json:"medical_history" gorm:"column:medical_history;type:json"`
}
