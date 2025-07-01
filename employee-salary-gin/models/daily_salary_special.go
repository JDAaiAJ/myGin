package models

type DailySalarySpecial struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Date       string `json:"date"`
	UserID     uint   `json:"u_id" gorm:"column:u_id"`
	Name       string `json:"name"`        //名称
	Price      string `json:"price"`       //单价
	Quantity   int    `json:"quantity"`    // 件数
	InsertTime int64  `json:"insert_time"` // 插入时间
	UpdateTime int64  `json:"update_time"` // 更新时间
}

func (DailySalarySpecial) TableName() string {
	return "daily_salary_special"
}
