package repositories

import (
	"fmt"

	"github.com/Dokito555/ewallet/ewallet-ewallet/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WalletRepository struct {
	Repository[entity.Wallet]
	Log *logrus.Logger
	DB *gorm.DB
}

func NewWalletRepository(log *logrus.Logger, db *gorm.DB) *WalletRepository {
	return &WalletRepository{
		Repository: Repository[entity.Wallet]{DB: db},
		Log: log,
		DB: db,
	}
}

func (r *WalletRepository) CreateWallet(wallet *entity.Wallet) error {
	return r.DB.Create(wallet).Error 
}

func (r *WalletRepository) UpdateBalance(userID int, amount float64) (entity.Wallet, error) {
	var wallet entity.Wallet
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Raw("SELECT id, user_id, balance FROM wallets WHERE user_id = ? FOR UPDATE", userID).Scan(&wallet).Error
		if err != nil {
			return err
		}

		if (wallet.Balance+amount) < 0 {
			r.Log.Warnf("user id %d current balance is not enough to perform transaction: %f - %f", userID, wallet.Balance, amount)
			return fmt.Errorf("user id %d current balance is not enough to perform transaction: %f - %f", userID, wallet.Balance, amount)
		}

		err = tx.Exec("UPDATE wallets SET balance = balance + ? WHERE user_id = ?", amount, userID ).Error
		if err != nil {
			return err
		}
		return nil
	})
	return wallet, err
}

func (r *WalletRepository) CreateWalletTrx(walletTrx *entity.WalletTransaction) error {
	return r.DB.Create(walletTrx).Error 
}

func (r *WalletRepository) GetWalletTransactionByReference(ref string) (entity.WalletTransaction, error) {
	var (
		resp entity.WalletTransaction
	)
	err := r.DB.Where("reference = ?", ref).Last(&resp).Error

	return resp, err
}

func (r *WalletRepository) GetWalletByUserID( userID int) (entity.Wallet, error) {
	var (
		resp entity.Wallet
	)

	err := r.DB.Where("user_id = ?", userID).First(&resp).Error
	return resp, err
}

func (r *WalletRepository) GetWalletByID(walletId int) (entity.Wallet, error) {
	var (
		resp entity.Wallet
	)

	err := r.DB.Where("id = ?", walletId).First(&resp).Error
	return resp, err
}

func (r *WalletRepository) GetWalletHistory(walletID int, offset int, limit int, transactionType string) ([]entity.WalletTransaction, error) {
	var (
		resp []entity.WalletTransaction
	)

	sql := r.DB
	if transactionType != "" {
		sql = sql.Where("wallet_transaction_type = ?", transactionType)
	}

	err := sql.Limit(limit).Offset(offset).Order("id DESC").Find(&resp).Error

	return resp, err
}

func (r *WalletRepository) InsertWalletLink(req *entity.WalletLink) error {
	return r.DB.Create(req).Error
}

func (r *WalletRepository) GetWalletLink( walletID int, clientSource string) (entity.WalletLink, error) {
	var (
		resp entity.WalletLink
		err error
	)

	err = r.DB.Where("wallet_id = ?", walletID).Where("client_source = ?", clientSource).First(&resp).Error
	return resp, err
}

func (r *WalletRepository) UpdateStatusWalletLink(walletID int, clientSource string, status string) (error) {
	var (
		err error
	)

	err = r.DB.Exec("UPDATE wallet_links SET status = ? WHERE wallet_id = ? AND client_source = ?", status, walletID, clientSource).Error
	return err
}