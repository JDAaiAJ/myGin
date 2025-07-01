package models

type DailySalary struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Date       string `json:"date"`
	UserID     uint   `json:"u_id" gorm:"column:u_id"`
	ClothingID uint   `json:"c_id" gorm:"column:c_id"` // 对应服装ID
	Quantity   int    `json:"quantity"`                // 件数
	InsertTime int64  `json:"insert_time"`             // 插入时间
	UpdateTime int64  `json:"update_time"`             // 更新时间
	//Total      float64   `json:"total"`                   // 总计 = Quantity * Price
}

func (DailySalary) TableName() string {
	return "daily_salary"
}
