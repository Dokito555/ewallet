package models

type WalletHistoryParam struct {
	Page                  int    `query:"page"`
	Limit                 int    `query:"limit"`
	WalletTransactionType string `query:"wallet_transaction_type"`
}

type WalletStructOTP struct {
	OTP string `json:"otp" validate:"required"`
}