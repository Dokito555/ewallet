package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/Dokito555/ewallet/ewallet-transaction/constants"
	"github.com/Dokito555/ewallet/ewallet-transaction/external"
	"github.com/Dokito555/ewallet/ewallet-transaction/internal/entity"
	"github.com/Dokito555/ewallet/ewallet-transaction/internal/models"
	"github.com/Dokito555/ewallet/ewallet-transaction/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type TransactionService struct {
	Log *logrus.Logger
	Validate *validator.Validate
	Config *viper.Viper
	DB *gorm.DB
	Repository *repositories.TransactionRepository
	WalletExt *external.ExtWallet
}

func NewTransactionService(log *logrus.Logger, va *validator.Validate, config *viper.Viper, db *gorm.DB, repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{
		Log: log,
		Validate: va,
		Config: config,
		DB: db,
		Repository: repo,
	}
}

func GenerateReference() string {
	now := time.Now()
	randomNum := rand.Intn(100)
	reference := fmt.Sprintf("%s%d", now.Format("20060102150405"), randomNum)
	return reference
}

func (s *TransactionService) CreateTransaction(req *entity.Transaction) (models.CreateTransactionResponse, error) {
	resp := models.CreateTransactionResponse{}

	jsonAddInfo := map[string]interface{}{}
	if req.AdditionalInfo != "" {
		err := json.Unmarshal([]byte(req.AdditionalInfo), &jsonAddInfo)
		if err != nil {
			s.Log.Warnf("additional info type is invalid: ", err)
			return resp, fmt.Errorf("additional info type is invalid")
		}
	}

	req.TransactionStatus = constants.TRANSACTION_STATUS_PENDING
	req.Reference = GenerateReference()
	err := s.Repository.CreateTransaction(req)
	if err != nil {
		s.Log.Warnf("failed to create transaction: ", err)
		return resp, fmt.Errorf("failed to create transaction")
	}

	resp.Reference = req.Reference
	resp.TransactionStatus = req.TransactionStatus
	return resp, nil
}

func (s *TransactionService) UpdateStatusTransaction(ctx context.Context, tokenData models.TokenData, req *models.UpdateStatusTransaction) error {
	trx, err := s.Repository.GetTransactionByReference(req.Reference, false)
	if err != nil {
		s.Log.Warnf("failed to get transaction: ", err)
		return fmt.Errorf("failed to get transaction")
	}

	statusValid := false
	mapStatusFlow := constants.MapTransactionStatusFlow[trx.TransactionStatus]
	for i := range mapStatusFlow {
		if mapStatusFlow[i] == req.TransactionStatus {
			statusValid = true
		}
	}

	if !statusValid {
		return fmt.Errorf("status transaction status flow invalid. request status = %s", req.TransactionStatus)
	}

	currentAdditionalInfo := map[string]interface{}{}
	if trx.AdditionalInfo != "" && trx.AdditionalInfo != "null" {
		if err := json.Unmarshal([]byte(trx.AdditionalInfo), &currentAdditionalInfo); err != nil {
			s.Log.Warnf("failed to unmarshal current additional info: ", err)
			return fmt.Errorf("failed to unmarshal current addtional info")
		}
	}

	newAdditionalInfo := map[string]interface{}{}
	if req.AdditionalInfo != "" && req.AdditionalInfo != "null" {
		if err := json.Unmarshal([]byte(req.AdditionalInfo), &newAdditionalInfo); err != nil {
			s.Log.Warnf("failed to unmarshal current additional info: ", err)
			return fmt.Errorf("failed to unmarshal current addtional info")
		}
	}

	for key, val := range newAdditionalInfo {
		currentAdditionalInfo[key] = val
	}

	byteAdditionalInfo, err := json.Marshal(currentAdditionalInfo)
	if err != nil {
		s.Log.Warnf("failed to marshal addtional info: ", err)
		return fmt.Errorf("failed to marshal addtional info")
	}

	reqUpdateBalance := models.UpdateBalance{
		Amount:    trx.Amount,
		Reference: trx.Reference,
	}

	if req.TransactionStatus == constants.TRANSACTION_STATUS_REVERSED {
		reqUpdateBalance.Reference = "REVERSED-" + req.Reference
		now := time.Now()
		expiredReversalTime := trx.CreatedAt.Add(constants.MAXIMUM_REVERSED_DURATION)
		if now.After(expiredReversalTime) {
			s.Log.Warnf("reversal duration is already expired")
			return fmt.Errorf("reversal duration is already expired")
		}
	}

	var (
		errUpdate error
	)

	switch trx.TransactionType {
	case constants.TRANSACTION_TYPE_TOPUP:
		switch req.TransactionStatus {
		case constants.TRANSACTION_STATUS_SUCCESS:
			_, errUpdate = s.WalletExt.CreditBalance(ctx, reqUpdateBalance, tokenData.Token)
		case constants.TRANSACTION_STATUS_REVERSED:
			_, errUpdate = s.WalletExt.DebitBalance(ctx, reqUpdateBalance, tokenData.Token)
		}
	case constants.TRANSACTION_TYPE_PURCHASE:
		switch req.TransactionStatus {
		case constants.TRANSACTION_STATUS_SUCCESS:
			_, errUpdate = s.WalletExt.DebitBalance(ctx, reqUpdateBalance, tokenData.Token)
		case constants.TRANSACTION_STATUS_REVERSED:
			_, errUpdate = s.WalletExt.CreditBalance(ctx, reqUpdateBalance, tokenData.Token)
		}
	}
	if errUpdate != nil {
		s.Log.Warnf("failed to update balance: ", errUpdate)
		return fmt.Errorf("failed to update balance")
	}

	err = s.Repository.UpdateStatusTransaction(req.Reference, req.TransactionStatus, string(byteAdditionalInfo))
	if err != nil {
		s.Log.Warnf("failed to update status transaction")
		return fmt.Errorf("failed to update status transaction")
	}

	trx.TransactionStatus = req.TransactionStatus
	// s.sendNotification(ctx, tokenData, trx)

	return nil
}

func (s *TransactionService) GetTransactions(userID int) ([]entity.Transaction, error) {
	return s.Repository.GetTransactionsByUserID(userID)
}

func (s *TransactionService) GetTransactionDetail( reference string) (entity.Transaction, error) {
	return s.Repository.GetTransactionByReference(reference, true)
}

// func (s *TransactionService) sendNotification(tokenData models.TokenData, trx entity.Transaction) {
// 	if trx.TransactionType == constants.TRANSACTION_TYPE_PURCHASE && trx.TransactionStatus == constants.TRANSACTION_STATUS_SUCCESS {
// 		s.Wall.SendNotification(ctx, tokenData.Email, "PURCHASE_SUCCESS", map[string]string{
// 			"full_name":   tokenData.FullName,
// 			"description": trx.Description,
// 			"reference":   trx.Reference,
// 			"date":        trx.CreatedAt.Format("2006-01-02 15:04:05"),
// 		})
// 	}
// }

func (s *TransactionService) RefundTransaction(ctx context.Context, tokenData models.TokenData, req *models.RefundTransaction) (models.CreateTransactionResponse, error) {
	var (
		resp models.CreateTransactionResponse
	)
	trx, err := s.Repository.GetTransactionByReference(req.Reference, false)
	if err != nil {
		s.Log.Warnf("failed to get transaction: ", err)
		return resp, fmt.Errorf("failed to get transaction")
	}

	if trx.TransactionStatus != constants.TRANSACTION_STATUS_SUCCESS && trx.TransactionType != constants.TRANSACTION_TYPE_PURCHASE {
		s.Log.Warnf("current transaction status is not SUCCESS or transaction type is not PURCHASE")
		return resp, fmt.Errorf("urrent transaction status is not SUCCESS or transaction type is not PURCHASE")
	}
	refundReference := "REFUND-" + req.Reference
	reqCreditBalance := models.UpdateBalance{
		Reference: refundReference,
		Amount:    trx.Amount,
	}
	_, err = s.WalletExt.CreditBalance(ctx, reqCreditBalance, tokenData.Token)
	if err != nil {
		s.Log.Warnf("failed to credit transaction: ", err)
		return resp, fmt.Errorf("failed to credit transaction")
	}

	transaction := entity.Transaction{
		UserID:            int(tokenData.UserID),
		Amount:            trx.Amount,
		TransactionType:   constants.TRANSACTION_TYPE_REFUND,
		TransactionStatus: constants.TRANSACTION_STATUS_SUCCESS,
		Reference:         refundReference,
		Description:       req.Description,
		AdditionalInfo:    req.AdditionalInfo,
	}
	err = s.Repository.CreateTransaction(&transaction)
	if err != nil {
		s.Log.Warnf("failed to insert new transaction: ", err)
		return resp, fmt.Errorf("failed to insert new transaction")
	}

	resp.Reference = refundReference
	resp.TransactionStatus = transaction.TransactionStatus
	return resp, nil
}
