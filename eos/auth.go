package eos

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"queding.com/go/common"
	"queding.com/go/common/config"
	"queding.com/go/common/log"
	"queding.com/go/common/rest"
	"queding.com/go/common/task"

	"cpay/eos/block"
	"cpay/eos/dao"
	"cpay/eos/util"
)

var (
	authNoticeUrl = config.GetString("authNoticeUrl")
	logger        *logrus.Logger
)

func init() {
	logger = log.AppLogger()
}

func HandleAuths() {
	common.SetInterval("auth", time.Second*5, checkNewAuth)
	common.SetInterval("auth", time.Second, handleNewAuth)
}

func checkNewAuth() {
	authTokens, err := block.GetUnconfirmedAuthTokens()
	if err != nil {
		logger.Errorln("GetUnconfirmedAuthTokens error:", err)
		return
	}
	for _, t := range authTokens {
		go func(t block.AuthToken) {
			if err := dao.UpdateAuthAccount(t.Account, t.TokenSha256); err != nil {
				logger.Errorln(err)
				return
			}
			go handleNewAuth()

			// 去掉合约记录
			trID, err := rmAuthToken(t.TokenSha256)
			if err != nil {
				if _, ok := err.(task.DuplicationError); !ok {
					logger.Errorln(err)
				}
				return
			}
			logger.Infof("Removed auth token=%s account=%s trID=%s", t.TokenSha256, t.Account, trID)
		}(t)
	}
}

func rmAuthToken(tokenSha256Str string) (txId string, err error) {
	result, err := task.ExcuteTask(&task.Task{
		Key: fmt.Sprintf("auth.rmAuthToken-%v", tokenSha256Str),
		GetResult: func() (interface{}, error) {
			return block.RmAuthToken(tokenSha256Str)
		},
	})
	if err != nil {
		return
	}
	return result.(string), nil
}

func handleNewAuth() {
	tokens, err := dao.GetConsumerUnConfirmedAuth()
	if err != nil {
		logger.Errorln("handleNewAuth error:", err)
		return
	}
	for _, token := range *tokens {
		// 通知 consumer
		reqOpts := map[string]string{
			"eos":   token.EosAccount,
			"token": token.Token,
		}
		logger.Infof("contact consumer, account: %s token: %s", token.EosAccount, token.Token)
		encrypted, err := util.EncryptRestBody(reqOpts)
		if err != nil {
			logger.Errorln("handleNewAuth error:", err)
			return
		}
		//reqBody, _ := json.Marshal(encrypted)
		logger.Infof("=======%v", encrypted)
		if _, err := rest.PostJson(authNoticeUrl, encrypted); err != nil {
			errStr := fmt.Sprintf("auth notice fail. error: %s", err)
			logger.Errorln(errStr)
			dao.SetAuthConfirmFail(token.Token, err.Error())
		} else {
			dao.SetAuthConfirmed(token.Token)
			logger.Infof("auth notice done, eosAccount=%s", token.EosAccount)
		}
	}
}
