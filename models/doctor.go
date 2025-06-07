package models

//"database/sql"
// 或者其他需要的 import
type Doctor struct {
	DoctorID   int    `json:"doctor_id" db:"doctor_id"`
	DeptID     int    `json:"dept_id" db:"dept_id"`
	DoctorName string `json:"doctor_name" db:"doctor_name"`
	DoctorInfo string `json:"doctor_info" db:"doctor_info"`
	Password   string `json:"password" db:"password"` // 儲存加密過後密碼
}
