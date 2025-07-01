package models

type Clothing struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Code       string `json:"code"`                       // 服装编号
	Price      string `json:"price"`                      // 单价
	Source     string `json:"source,omitempty"`           // 来源信息
	UserID     uint   `json:"user_id" gorm:"column:u_id"` //添加人id
	Image      string `json:"image"`                      // 服饰图片
	InsertTime int64  `json:"insert_time"`                // 插入时间
	UpdateTime int64  `json:"update_time"`                // 更新时间
}

func (Clothing) TableName() string {
	return "clothing"
}
