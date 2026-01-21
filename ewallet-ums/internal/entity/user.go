package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username" gorm:"column:username;type:varchar(100)" validate:"required"`
	Email       string    `json:"email" gorm:"column:email;type:varchar(100)" validate:"required"`
	FullName    string    `json:"full_name" gorm:"column:full_name;type:varchar(100)" validate:"required"`
	PhoneNumber string    `json:"phone_number" gorm:"column:phone_number;type:varchar(15)" validate:"required"`
	Address     string    `json:"address" gorm:"column:address;type:text"`
	Dob         string    `json:"dob" gorm:"column:dob;type:date"`
	Password    string    `json:"password,omitempty" gorm:"column:password;type:varchar(255)" validate:"required"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

func (*User) TableName() string {
	return "users"
}

func (l User) Validate() error {
	v := validator.New()
	return v.Struct(l)
}


type UserSession struct {
	ID                  uint      `gorm:"primarykey"`
	UserID              int       `json:"user_id" gorm:"int" validate:"required"`
	Token               string    `json:"token" gorm:"text" validate:"required"`
	RefreshToken        string    `json:"refresh_token" gorm:"text" validate:"required"`
	TokenExpired        time.Time `json:"-" validate:"required"`
	RefreshTokenExpired time.Time `json:"-" validate:"required"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (*UserSession) TableName() string {
	return "user_sessions"
}

func (l UserSession) Validate() error {
	v := validator.New()
	return v.Struct(l)
}