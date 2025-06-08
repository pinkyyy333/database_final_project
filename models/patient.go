package models

//import "time"

type Patient struct {
	PatientID      string    `gorm:"column:Patient_ID;primaryKey"`
	PatientName    string    `gorm:"column:Patient_Name"`
	PatientGender  string    `gorm:"column:Patient_Gender"`
	PatientBirth   string    `gorm:"column:Patient_Birth"` // string 或 time.Time，視需求
	PatientPhone   string    `gorm:"column:Patient_Phone;unique"`
	Password       string    `gorm:"column:Password"`
	DrugAllergy    string    `gorm:"column:drug_allergy"`
	FoodAllergy    string    `gorm:"column:food_allergy"`
	MedicalHistory string    `gorm:"column:medical_history"`
}

func (Patient) TableName() string {
	return "patients"
}
