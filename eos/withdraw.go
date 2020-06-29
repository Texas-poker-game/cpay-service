package eos

import (
	"strings"
	"time"

	"cpay/eos/block"
	"cpay/eos/dao"
)

func HandleWithdraws() {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorln(r)
		}
	}()

	t := time.NewTicker(3 * time.Second)

	for range t.C {
		go handleNewWithdraw()
	}

}

func handleNewWithdraw() {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorln(r)
		}
	}()

	withdraws, err := dao.GetNotSentWithdraws()
	if err != nil {
		logger.Errorf("handleNewWithdraw error: %s", err)
		return
	}

	for _, withdraw := range *withdraws {
		go func(withdraw dao.WithdrawPo) {
			txId, err := block.AddWithdraw(uint64(withdraw.ID), withdraw.From, withdraw.To, withdraw.Amount)

			if err != nil && !strings.Contains(err.Error(), "the withdraw id exists") {
				logger.Errorf("push transaction: %s", err)
				return
			}

			logger.Errorf("Transaction [%s] submitted to the network succesfully.", txId)
			dao.SetWithdrawSent(withdraw.ID, txId)
		}(withdraw)
	}
}
