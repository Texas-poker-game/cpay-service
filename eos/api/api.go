package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"queding.com/go/common/config"

	"cpay/eos/dao"
	"cpay/eos/util"
)

// 提现
func WithdrawHandler(c *gin.Context) {
	var encryptedBody util.EncryptedBody
	if err := c.ShouldBind(&encryptedBody); err != nil {
		c.Error(err)
		return
	}
	var req WithdrawRequest
	if err := util.DecryptRestBody(encryptedBody.Encrypted, &req); err != nil {
		c.Error(err)
		return
	}

	if req.Amount == 0 {
		c.Error(errors.New("amount 参数错误"))
		return
	}

	record := &dao.WithdrawPo{
		ConsumerSid: req.Sid,
		From:        config.GetString("eos.vault.contract"),
		To:          req.To,
		Amount:      req.Amount,
	}
	po, err := dao.GetWithdrawBySid(record.ConsumerSid)
	if err != nil {
		c.Error(err)
		return
	}

	if po != nil {
		record = po
	} else {
		if err = dao.CreateWithdraw(record); err != nil {
			c.Error(err)
			return
		}
	}

	success(c, &WithdrawVo{
		Sid:    record.ConsumerSid,
		Amount: record.Amount,
		To:     record.To,
	})
}

// auth token
func tokenGenerator() string {
	b := make([]byte, 32)
	rand.Read(b)
	token := base64.StdEncoding.EncodeToString(b)
	token = strings.NewReplacer("+", "a", "/", "b", "=", "").Replace(token)
	return token
}

func getSha256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func NewAuthToken(c *gin.Context) {
	token := tokenGenerator()
	po := dao.AuthTokenPo{
		Token:       token,
		TokenSha256: getSha256(token),
	}

	if err := dao.CreateAuthToken(&po); err != nil {
		c.Error(err)
		return
	}
	success(c, token)
}
