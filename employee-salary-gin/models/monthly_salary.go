package models

type MonthlySalary struct {
	ID     uint    `gorm:"primaryKey" json:"id"`
	Month  string  `json:"month"` // 年月主键
	UserID uint    `json:"u_id" gorm:"column:u_id"`
	Total  float64 `json:"total"` // 该月所有日工资总和
}

func (MonthlySalary) TableName() string {
	return "monthly_salary"
}
