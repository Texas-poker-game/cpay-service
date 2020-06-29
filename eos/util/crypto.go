package util

import (
	"encoding/hex"
	"encoding/json"

	"queding.com/go/common/log"
)

var (
	memoKey     = []byte("MXKP~C~$z*jGjwYZVlwbBAoh(Qr&eEy9")
	restBodyKey = []byte("DaDgEFX:yJLwW3.V@V-kbnQId,W!36LH")
)

var (
	logger = log.AppLogger()
)

func encrypt(v interface{}, key []byte) (hexStr string, err error) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	encryptData, err := AesEncrypt(data, key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(encryptData), nil
}

func decrypt(hexStr string, key []byte, v interface{}) (err error) {
	logger.Info("decrypt, hexStr=", hexStr)
	encryptData, err := hex.DecodeString(hexStr)
	if err != nil {
		logger.Info("decrypt err:", err)
		return
	}
	data, err := AesDecrypt(encryptData, key)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, v)
	return
}

func DecryptMemo(buf string) error {
	data := new(map[string]interface{})
	return decrypt(buf, memoKey, data)
}

func EncryptRestBody(v interface{}) (encryptedBody interface{}, err error) {
	if err != nil {
		return nil, err
	}
	encrypted, err := encrypt(v, restBodyKey)
	if err != nil {
		return
	}
	encryptedBody = map[string]interface{}{
		"encrypted": encrypted,
	}
	return
}

func DecryptRestBody(hexStr string, v interface{}) error {
	return decrypt(hexStr, restBodyKey, v)
}

type EncryptedBody struct {
	Encrypted string `json:"encrypted"`
}
