package classroom

type Classroom struct {
	ID       int64  `gorm:"primaryKey;column:id" json:"id"`
	Number   string `gorm:"column:number" json:"number"`
	Building string `gorm:"column:building" json:"building"`
}
