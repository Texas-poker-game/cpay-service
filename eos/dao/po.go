package dao

import (
	"encoding/json"
	"time"
)

// Deposit
type DepositPo struct {
	ID             uint `gorm:"primary_key"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Amount         float64
	From           string
	To             string
	TransactionId  string
	Memo           string
	BlockConfirmed bool
	Status         int // 0 初始 1 通知成功 2 通知失败
	TransactionMs  uint64
	NoticeError    string
}

func (DepositPo) TableName() string {
	return "eos_deposits"
}

// Withdraw
type WithdrawPo struct {
	ID            uint `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Amount        float64
	From          string
	To            string
	Memo          string
	TransactionId string
	Status        int // 0 初始 1 已发送未确认 2 已确认
	ConsumerSid   string
}

func (WithdrawPo) TableName() string {
	return "eos_withdraws"
}

func (po *WithdrawPo) String() string {
	bytes, _ := json.Marshal(po)
	return string(bytes)
}

// AuthToken
type AuthTokenPo struct {
	Token       string `gorm:"primary_key"`
	TokenSha256 string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	EosAccount  string
	Status      int // 0 初始 1 已绑定账号 2 通知成功 3 通知失败
	NoticeError string
}

func (AuthTokenPo) TableName() string {
	return "eos_auth_tokens"
}
