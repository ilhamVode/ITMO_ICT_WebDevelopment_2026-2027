package semester

type Semester struct {
	ID           int64  `gorm:"primaryKey;column:id" json:"id"`
	AcademicYear string `gorm:"column:academic_year" json:"academic_year"`
	Number       int    `gorm:"column:number" json:"number"`
}
