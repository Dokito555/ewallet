package entity

import (
	"time"

	"github.com/go-playground/validator"
)

type Transaction struct {
	ID                int       `json:"-"`
	UserID            int       `json:"user_id"`
	Amount            float64   `json:"amount" gorm:"column:amount";type:decimal(15,2) validation:"required"`
	TransactionType   string    `json:"transaction_type" gorm:"column:transaction_type";type:enum('TOPUP','PURCHASE','REFUND) validation:"required"`
	TransactionStatus string    `json:"transaction_status" gorm:"column:transaction_status";type:enum('PENDING','SUCCESS','FAILED','REVERSED')`
	Reference         string    `json:"reference" gorm:"column:reference";type:varchar(255)`
	Description       string    `json:"description" gorm:"column:description";type:varchar(255) validation:"required"`
	AdditionalInfo    string    `json:"additional_info" gorm:"column:additional_info";type:text`
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
}

func (l Transaction) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

func (*Transaction) TableName() string {
	return "transactions"
}
