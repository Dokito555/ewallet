package services

import (
	"fmt"
	"math/rand/v2"

	"github.com/Dokito555/ewallet/ewallet-ewallet/internal/entity"
	"github.com/Dokito555/ewallet/ewallet-ewallet/internal/models"
	"github.com/Dokito555/ewallet/ewallet-ewallet/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type WalletService struct {
	Log        *logrus.Logger
	Validate   *validator.Validate
	Config     *viper.Viper
	DB         *gorm.DB
	WalletRepo *repositories.WalletRepository
}

func NewWalletService(log *logrus.Logger, val *validator.Validate, config *viper.Viper, db *gorm.DB, walletRepo *repositories.WalletRepository) *WalletService {
	return &WalletService{
		Log: log,
		Validate: val,
		Config: config,
		DB: db,
		WalletRepo: walletRepo,
	}
}

func (s *WalletService) Create(wallet *entity.Wallet) error {
	return s.WalletRepo.CreateWallet(wallet)
}

func (s *WalletService) CreditBalance(userID int, req models.TransactionRequest) (models.BalanceResponse, error) {
	var (
		resp models.BalanceResponse
	)

	history, err := s.WalletRepo.GetWalletTransactionByReference(req.Reference)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			s.Log.Warnf("failed to check reference")
			return resp, fmt.Errorf("failed to check reference")
		}
	}
	if history.ID > 0 {
		s.Log.Warnf("reference is duplicate")
		return resp, fmt.Errorf("reference is duplicate")
	}

	wallet, err := s.WalletRepo.UpdateBalance(userID, req.Amount)
	if err != nil {
		s.Log.Errorf("failed to update balance: ", err)
		return resp, fmt.Errorf("failed to update balance")
	}

	walletTrx := &entity.WalletTransaction{
		WalletID:              wallet.ID,
		Amount:                req.Amount,
		Reference:             req.Reference,
		WalletTransactionType: "CREDIT",
	}

	err = s.WalletRepo.CreateWalletTrx(walletTrx)
	if err != nil {
		s.Log.Warnf("failed to insert wallet transaction: ", err)
		return resp, fmt.Errorf("failed to insert wallet transaction")
	}

	resp.Balance = wallet.Balance + req.Amount

	return resp, nil
}

func (s *WalletService) DebitBalance(userID int, req models.TransactionRequest) (models.BalanceResponse, error) {
	var (
		resp models.BalanceResponse
	)

	history, err := s.WalletRepo.GetWalletTransactionByReference(req.Reference)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			s.Log.Warnf("failed to check reference: ", err)
			return resp, fmt.Errorf("failed to check reference")
		}
	}
	if history.ID > 0 {
		s.Log.Warnf("reference is duplicate")
		return resp, fmt.Errorf("reference is duplicate")
	}

	wallet, err := s.WalletRepo.UpdateBalance(userID, -req.Amount)
	if err != nil {
		s.Log.Warnf("failed to update balance: ", err)
		return resp, fmt.Errorf("failed to update balance")
	}

	walletTrx := &entity.WalletTransaction{
		WalletID:              wallet.ID,
		Amount:                req.Amount,
		Reference:             req.Reference,
		WalletTransactionType: "DEBIT",
	}

	err = s.WalletRepo.CreateWalletTrx(walletTrx)
	if err != nil {
		s.Log.Warnf("failed to insert wallet transaction: ", err)
		return resp, fmt.Errorf("failed to insert wallet transaction")
	}

	resp.Balance = wallet.Balance - req.Amount

	return resp, nil
}

func (s *WalletService) GetBalance(userID int) (models.BalanceResponse, error) {
	var (
		resp models.BalanceResponse
	)

	wallet, err := s.WalletRepo.GetWalletByUserID(userID)
	if err != nil {
		s.Log.Warnf("failed to get wallet", err)
		return resp, fmt.Errorf("failed to get wallet")
	}

	resp.Balance = wallet.Balance

	return resp, nil
}

func (s *WalletService) ExGetBalance(walletId int) (models.BalanceResponse, error) {
	var (
		resp models.BalanceResponse
	)

	wallet, err := s.WalletRepo.GetWalletByID(walletId)
	if err != nil {
		s.Log.Warnf("failed to get wallet: ", err)
		return resp, fmt.Errorf("failed to get wallet")
	}

	resp.Balance = wallet.Balance

	return resp, nil
}


func (s *WalletService) GetWalletHistory(userID int, param models.WalletHistoryParam) ([]entity.WalletTransaction, error) {
	var (
		resp []entity.WalletTransaction
	)

	wallet, err := s.WalletRepo.GetWalletByUserID(userID)
	if err != nil {
		s.Log.Warnf("failed to get wallet: ", err)
		return resp,fmt.Errorf("failed to get wallet")
	}

	offset := (param.Page-1) * param.Limit 
	resp, err = s.WalletRepo.GetWalletHistory(wallet.ID, offset, param.Limit, param.WalletTransactionType)
	if err != nil {
		s.Log.Warnf("failed to get wallet history: ", err)
		return resp, fmt.Errorf("failed to get wallet history")
	}

	return resp, nil
}


func (s *WalletService) CreateWalletLink(clientSource string, req *entity.WalletLink) (models.WalletStructOTP, error) {
	req.ClientSource = clientSource
	req.Status = "PENDING"
	req.OTP = fmt.Sprintf("%d", rand.IntN(999999))

	resp := models.WalletStructOTP{
		OTP: req.OTP,
	}

 	err := s.WalletRepo.InsertWalletLink(req)	
	if err != nil {
		s.Log.Warnf("failed to insert wallet link: ", err)
		return resp, fmt.Errorf("failed to insert wallet link")
	}
	return resp, nil
}

func (s *WalletService) WalletLinkConfirmation(walletID int, clientSource string, otp string) (error) {
	walletLink, err := s.WalletRepo.GetWalletLink(walletID, clientSource)
	if err != nil {
		s.Log.Warnf("failed to get wallet link")
		return fmt.Errorf("failed to get wallet link")
	}

	if walletLink.Status != "PENDING" {
		s.Log.Warnf("wallet status is not pending")
		return fmt.Errorf("wallet status is not pending")
	}

	if walletLink.OTP != otp {
		s.Log.Warnf("invalid otp, request = %s, stored = %s", otp, walletLink.OTP)
		return fmt.Errorf("invalid otp")
	}

	return s.WalletRepo.UpdateStatusWalletLink(walletID, clientSource, "LINKED")
}

func (s *WalletService) WalletUnlink(walletID int, clientSource string) (error) {
	return s.WalletRepo.UpdateStatusWalletLink(walletID, clientSource, "UNLINKED")
}