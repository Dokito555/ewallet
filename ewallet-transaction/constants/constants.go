package constants

import "time"

const (
	ErrFailedBadRequest = "Data/Request doesn't match"
	ErrServerError      = "There's a problem with the server"
	SUCCESSMessage      = "SUCCESS"
)

var MappingClient = map[string]string {
	"ecommerce": "ecommerce-secret-key",
}

const (
	ERROR_FAILED_BAD_REQUEST = "Data/Request doesn't match"
	ERROR_SERVER_ERROR       = "There's a problem with the server"
	SUCCESS_MESSAGE          = "SUCCESS"
)

const (
	TRANSACTION_STATUS_PENDING  = "PENDING"
	TRANSACTION_STATUS_SUCCESS  = "SUCCESS"
	TRANSACTION_STATUS_FAILED   = "FAILED"
	TRANSACTION_STATUS_REVERSED = "REVERSED"
)

const (
	TRANSACTION_TYPE_TOPUP    = "TOPUP"
	TRANSACTION_TYPE_PURCHASE = "PURCHASE"
	TRANSACTION_TYPE_REFUND   = "REFUND"
)

var MapTransactionType = map[string]bool{
	TRANSACTION_TYPE_TOPUP:    true,
	TRANSACTION_TYPE_REFUND:   true,
	TRANSACTION_TYPE_PURCHASE: true,
}

var MapTransactionStatusFlow = map[string][]string{
	TRANSACTION_STATUS_PENDING: {
		TRANSACTION_STATUS_SUCCESS,
		TRANSACTION_STATUS_FAILED,
	},
	TRANSACTION_STATUS_SUCCESS: {
		TRANSACTION_STATUS_REVERSED,
	},
	TRANSACTION_STATUS_FAILED: {
		TRANSACTION_STATUS_SUCCESS,
	},
}

var (
	MAXIMUM_REVERSED_DURATION = time.Hour * 24
)
