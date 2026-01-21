package entity

import "time"

type Wallet struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id" gorm:"column:user_id; unique"`
	Balance   float64   `json:"balance" gorm:"column:balance;type:decimal(15, 2)"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (*Wallet) TableName() string {
	return "wallets"
}

type WalletTransaction struct {
	ID                    int       `json:"-"`
	WalletID              int       `json:"wallet_id" gorm:"column:wallet_id"`
	Amount                float64   `json:"amount" gorm:"column:amount;type:decimal(15, 2)"`
	WalletTransactionType string    `json:"wallet_transaction_type" gorm:"column:wallet_transaction_tpe;type:ENUM('CREDIT', 'DEBIT')"`
	Reference             string    `json:"reference" gorm:"column:reference;type:varchar(100);unique"`
	CreatedAt             time.Time `json:"date"`
	UpdatedAt             time.Time `json:"-"`
}

func (*WalletTransaction) TableName() string {
	return "wallet_transactions"
}

type WalletLink struct {
	ID           int       `json:"id"`
	WalletID     int       `json:"wallet_id" gorm:"column:wallet_id" validate:"required"`
	ClientSource string    `json:"client_source" gorm:"column:client_source;type:varchar(100)"`
	OTP          string    `json:"otp" gorm:"column:otp;type:varchar(6)"`
	Status       string    `json:"status" gorm:"column:status;type:varchar(10)"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

func (*WalletLink) TableName() string {
	return "wallet_links"
}
