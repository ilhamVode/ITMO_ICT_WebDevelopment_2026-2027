package subject

type Subject struct {
	ID   int64  `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}
