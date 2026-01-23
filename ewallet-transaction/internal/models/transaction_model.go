package models

import "github.com/go-playground/validator"


type CreateTransactionResponse struct {
	Reference         string `json:"reference"`
	TransactionStatus string `json:"transaction_status"`
}

type UpdateStatusTransaction struct {
	Reference         string `json:"reference" validation:"required"`
	TransactionStatus string `json:"transaction_status" validation:"required"`
	AdditionalInfo    string `json:"additional_info"`
}

func (l UpdateStatusTransaction) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type RefundTransaction struct {
	Reference      string `json:"reference" validation:"required"`
	Description    string `json:"description" validation:"required"`
	AdditionalInfo string `json:"additional_info"`
}

func (l RefundTransaction) Validate() error {
	v := validator.New()
	return v.Struct(l)
}
