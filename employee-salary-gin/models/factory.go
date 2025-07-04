package models

type Factory struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `json:"name"`
	InsertTime int64  `json:"insert_time"` // 插入时间
	UpdateTime int64  `json:"update_time"` // 更新时间
}

func (Factory) TableName() string {
	return "factory" // 不加复数
}
