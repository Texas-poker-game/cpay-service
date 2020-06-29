package dao

import (
	"time"

	"github.com/sirupsen/logrus"
	"queding.com/go/common/log"
)

var (
	logger *logrus.Logger
)

func init() {
	logger = log.AppLogger()
}

// deposit
func CreateDeposit(deposit *DepositPo) error {
	logger.Infof("memoe: %s", deposit.Memo)
	return db.FirstOrCreate(deposit).Error
}

func GetConsumerUnConfirmedDeposits() (*[]DepositPo, error) {
	result := make([]DepositPo, 0)
	err := db.Limit(10).Find(&result, "status = ?", "0").Error
	return &result, err
}

func SetDepositBlockConfirmed(id uint) error {
	return db.Model(&DepositPo{}).UpdateColumns(DepositPo{
		ID:             id,
		BlockConfirmed: true,
		UpdatedAt:      time.Now(),
	}).Error
}

func SetDepositConsumerConfirmed(id uint) error {
	return db.Model(&DepositPo{}).Where("id = ? and status <> 1", id).UpdateColumns(DepositPo{
		Status:      1,
		NoticeError: "",
		UpdatedAt:   time.Now(),
	}).Error
}

func SetDepositConsumerConfirmFail(id uint, msg string) error {
	return db.Model(&DepositPo{}).Where("id = ? and status <> 1", id).UpdateColumns(DepositPo{
		Status:      2,
		NoticeError: msg,
		UpdatedAt:   time.Now(),
	}).Error
}

// withdraw
func CreateWithdraw(withdraw *WithdrawPo) error {
	return db.Create(withdraw).Error
}

func GetWithdrawBySid(sid string) (*WithdrawPo, error) {
	var records []WithdrawPo
	if err := db.Where("consumer_sid = ?", sid).Find(&records).Error; err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	return &records[0], nil
}

func GetNotSentWithdraws() (*[]WithdrawPo, error) {
	var result []WithdrawPo
	if err := db.Limit(10).Find(&result, "status = ?", 0).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func SetWithdrawSent(id uint, transactionId string) error {
	return db.Model(&WithdrawPo{}).UpdateColumns(WithdrawPo{
		ID:            id,
		TransactionId: transactionId,
		Status:        1, // 表示已发送未确认
		UpdatedAt:     time.Now(),
	}).Error
}

// auth token
func CreateAuthToken(po *AuthTokenPo) error {
	return db.Create(po).Error
}

func GetConsumerUnConfirmedAuth() (*[]AuthTokenPo, error) {
	result := make([]AuthTokenPo, 0)
	err := db.Limit(10).Find(&result, "status = ?", 1).Error
	return &result, err
}

func UpdateAuthAccount(account, tokenSha256 string) error {
	updateValues := map[string]interface{}{"eos_account": account, "status": 1}
	return db.Model(&AuthTokenPo{}).
		Where("token_sha256 = ? and eos_account=''", tokenSha256).
		Updates(updateValues).Error
}

func SetAuthConfirmed(token string) error {
	return db.Model(&AuthTokenPo{}).UpdateColumns(AuthTokenPo{
		Token:     token,
		Status:    2,
		UpdatedAt: time.Now(),
	}).Error
}

func SetAuthConfirmFail(token, msg string) error {
	return db.Model(&AuthTokenPo{}).UpdateColumns(AuthTokenPo{
		Token:       token,
		Status:      3,
		NoticeError: msg,
		UpdatedAt:   time.Now(),
	}).Error
}
