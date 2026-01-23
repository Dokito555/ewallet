package models

type UpdateBalance struct {
	Reference	string	`json:"reference"`
	Amount		float64	`json:"amount"`
}

type UpdateBalanceResponse struct {
	Message		string	`json:"message"`
	Data		struct {
		Balance		float64		`json:"balance"`
	}	`json:"data"`
}