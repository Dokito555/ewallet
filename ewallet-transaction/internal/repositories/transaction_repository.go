package repositories

import (
	"github.com/Dokito555/ewallet/ewallet-transaction/constants"
	"github.com/Dokito555/ewallet/ewallet-transaction/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	Repository[entity.Transaction]
	Log *logrus.Logger
	DB *gorm.DB
}

func NewTransactionRepository(log *logrus.Logger, db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		Repository: Repository[entity.Transaction]{DB: db},
		Log: log,
		DB: db,
	}
}

func (r *TransactionRepository) CreateTransaction(trx *entity.Transaction) error {
	return r.DB.Create(&trx).Error
}

func (r *TransactionRepository) UpdateStatusTransaction(reference string, status string, additionalInfo string) error {
	return r.DB.Exec("UPDATE transactions SET transaction_status = ?, additional_info = ? WHERE reference = ?", status, additionalInfo, reference).Error
}

func (r *TransactionRepository) GetTransactionByReference(reference string, includeRefund bool) (entity.Transaction, error) {
	var (
		resp entity.Transaction
	)

	sql := r.DB.Where("reference = ?", reference)
	if !includeRefund {
		sql = sql.Where("transaction_type != ?", constants.TRANSACTION_TYPE_REFUND)
	}

	err := sql.Last(&resp).Error
	return resp, err
}

func (r *TransactionRepository) GetTransactionsByUserID(userId int) ([]entity.Transaction, error) {
	var (
		resp []entity.Transaction
	)
	err := r.DB.Where("user_id = ?", userId).Find(&resp).Order("id DESC").Error
	return resp, err
}
