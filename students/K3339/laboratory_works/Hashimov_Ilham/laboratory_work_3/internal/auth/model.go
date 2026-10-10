package auth

import "time"

type User struct {
	ID           int64  `gorm:"primaryKey;column:id" json:"id"`
	Username     string `gorm:"column:username" json:"username"`
	PasswordHash string `gorm:"column:password_hash" json:"-"`
}

type Token struct {
	Key       string    `gorm:"primaryKey;column:key" json:"-"`
	UserID    int64     `gorm:"column:user_id" json:"-"`
	ExpiresAt time.Time `gorm:"column:expires_at" json:"-"`
}

func (User) TableName() string  { return "api_users" }
func (Token) TableName() string { return "api_tokens" }
