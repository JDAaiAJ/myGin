package models

type User struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username" gorm:"unique"`
	Password   string `json:"password"`
	FID        uint   `json:"f_id"`
	Type       int    `json:"type"`
	Status     int    `json:"status"`
	InsertTime int64  `json:"insert_time"` // 插入时间
	UpdateTime int64  `json:"update_time"` // 更新时间
}

func (User) TableName() string {
	return "user" // 不加复数
}
