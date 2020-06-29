package eos

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"queding.com/go/common"
	"queding.com/go/common/config"
	"queding.com/go/common/rest"
	"queding.com/go/common/task"

	"cpay/eos/block"
	"cpay/eos/dao"
	"cpay/eos/util"
)

var (
	httpClient       = http.Client{}
	depositNoticeUrl = config.GetString("depositNoticeUrl")
)

func HandleDeposits() {
	common.SetInterval("deposit", time.Second*5, checkNewDeposit)
	common.SetInterval("deposit", time.Second, handleNewDeposit)
}

func checkNewDeposit() {
	deposits, err := block.GetConsumerUnConfirmedDeposits()
	if err != nil {
		logger.Errorln("checkNewDeposit error:", err)
		return
	}
	for _, deposit := range deposits {
		ms, _ := strconv.ParseUint(deposit.TimeMs, 10, 64)
		depositPo := dao.DepositPo{
			ID:            deposit.ID,
			Amount:        deposit.Amount(),
			From:          deposit.From,
			To:            deposit.To,
			Memo:          deposit.Memo,
			TransactionMs: ms,
		}

		if err := dao.CreateDeposit(&depositPo); err != nil {
			logger.Errorln(err)
		}
	}
	handleNewDeposit()
}

func handleNewDeposit() {
	deposits, err := dao.GetConsumerUnConfirmedDeposits()
	if err != nil {
		logger.Errorln("handleNewDeposit error:", err)
		return
	}
	for _, deposit := range *deposits {
		txId, err := task.ExcuteTask(&task.Task{
			Key: fmt.Sprintf("block.SetDepositConfirm-%v", deposit.ID),
			GetResult: func() (interface{}, error) {
				return block.SetDepositConfirm(uint64(deposit.ID))
			},
		})
		if err != nil {
			if _, ok := err.(task.DuplicationError); !ok {
				logger.Errorln(err)
			}
			continue
		}
		logger.Infof("Call SetDepositConfirm, txId = %v", txId)
		if err := callNotice(deposit); err != nil {
			if _, ok := err.(task.DuplicationError); ok {
				continue
			}
			logger.Errorln(err)
			dao.SetDepositConsumerConfirmFail(deposit.ID, err.Error())
		} else {
			dao.SetDepositConsumerConfirmed(deposit.ID)
		}
	}
}

type depositNoticeOpts struct {
	From   string  `json:"from"`
	Sid    string  `json:"sid"`
	Amount float64 `json:"amount"`
	Memo   string  `json:"memo"`
}

func callNotice(deposit dao.DepositPo) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()

	opts := &depositNoticeOpts{
		From:   deposit.From,
		Sid:    strconv.Itoa(int(deposit.ID)),
		Amount: deposit.Amount,
		Memo:   deposit.Memo,
	}

	logger.Info("memo:", deposit.Memo)
	if err := util.DecryptMemo(deposit.Memo); err != nil {
		logger.Info("memoError:", err)
		// 解析 memo 失败
		return err
	}

	_, err = task.ExcuteTask(&task.Task{
		Key: fmt.Sprintf("callDepositNotice-%s", opts.Sid),
		GetResult: func() (interface{}, error) {
			encrypted, err := util.EncryptRestBody(opts)
			if err != nil {
				return nil, err
			}
			return rest.PostJson(depositNoticeUrl, encrypted)
		},
	})
	if err != nil {
		if _, ok := err.(task.DuplicationError); ok {
			return err
		}
		errStr := fmt.Sprintf("deposit notice fail, depositId=%d. error: %s", deposit.ID, err)
		logger.Errorln(errStr)
		return errors.New(errStr)
	}

	logger.Infof("deposit notice done, from=%s amount=%f memo=%s", deposit.From, deposit.Amount, deposit.Memo)
	return nil
}
