package block

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/eoscanada/eos-go"
)

type Deposit struct {
	ID          uint   `json:"id"`
	Timestamp   int64  `json:"timestamp"` // second
	From        string `json:"from"`
	To          string `json:"to"`
	QuantityRaw string `json:"quantity"`
	Memo        string `json:"memo"`
	Status      int    `json:"status"`
	TimeMs      string `json:"time_ms"`
}

func (d *Deposit) Time() time.Time {
	return time.Unix(d.Timestamp, 0)
}

func (d *Deposit) Amount() float64 {
	amount, _ := strconv.ParseFloat(strings.Split(d.QuantityRaw, " ")[0], 10)
	return amount
}

// deposit
func GetConsumerUnConfirmedDeposits() (deposits []Deposit, err error) {
	opts := GetTableOpts{
		Index:      "2",
		KeyType:    "i64",
		UpperBound: "0",
	}
	err = GetTableRows(eosVaultContract, eosVaultContract, "deposits", &opts, &deposits)
	return
}

func isDepositConfirmed(depositID uint64) (confirmed bool, err error) {
	opts := GetTableOpts{
		UpperBound: strconv.Itoa(int(depositID)),
		LowerBound: strconv.Itoa(int(depositID)),
	}
	var deposits []Deposit
	err = GetTableRows(eosVaultContract, eosVaultContract, "deposits", &opts, &deposits)
	if err != nil {
		return
	}

	confirmed = len(deposits) > 0 && deposits[0].Status > 0
	return
}

type confirmDepOpts struct {
	ID uint64 `json:"id"`
}

func SetDepositConfirm(depositID uint64) (txId string, err error) {
	confirmed, err := isDepositConfirmed(depositID)
	if err != nil {
		return
	}
	if confirmed {
		//err = errors.New(fmt.Sprintf("depositID %s has confirmed", txId))
		txId = "重复提交"
		return
	}

	return PushActionByAdmin("confirmdep", &confirmDepOpts{ID: depositID})
}

// withdraw
type addTransferOpts struct {
	ID       uint64
	From     eos.AccountName
	To       eos.AccountName
	Quantity eos.Asset
	Memo     string
}

func AddWithdraw(id uint64, fromStr, toStr string, amount float64) (txId string, err error) {
	return PushActionByAdmin("addwithdraw", &addTransferOpts{
		ID:       id,
		From:     eos.AccountName(fromStr),
		To:       eos.AccountName(toStr),
		Quantity: eos.NewEOSAsset(int64(amount * 10000)),
		Memo:     "",
	})
}

// auth
type AuthToken struct {
	TokenSha256 string `json:"token_sha256"`
	Account     string `json:"account"`
}

type rmAuthTokenOpts struct {
	Checksum256 eos.Checksum256
}

func GetUnconfirmedAuthTokens() (token []AuthToken, err error) {
	err = GetTableRows(eosVaultContract, eosVaultContract, "authtokens", nil, &token)
	return
}

func RmAuthToken(tokenSha256Str string) (string, error) {
	txOpts := &eos.TxOptions{}
	if err := txOpts.FillFromChain(eosApi); err != nil {
		return "", errors.New(fmt.Sprintf("filling tx opts: %s", err))
	}

	tokenSha256, err := hex.DecodeString(tokenSha256Str)
	if err != nil {
		return "", err
	}

	return PushActionByAdmin("rmauthtoken", &rmAuthTokenOpts{
		Checksum256: tokenSha256,
	})
}
