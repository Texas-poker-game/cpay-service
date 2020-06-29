package block

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/eoscanada/eos-go"
	"github.com/sirupsen/logrus"
	"queding.com/go/common/config"
	"queding.com/go/common/log"
)

var (
	eosApi           *eos.API
	eosVaultContract = config.GetString("eos.vault.contract")
	eosAdminAccount  = config.GetString("eos.admin.account")
	logger           *logrus.Logger
)

func init() {
	// loogger
	logger = log.AppLogger()

	endpoints := config.GetStringSlice("eos.endpoints")
	eosApi = eos.New(endpoints[0])

	// import admin private key
	keyBag := &eos.KeyBag{}

	if err := keyBag.ImportPrivateKey(config.GetString("eos.admin.private")); err != nil {
		panic(fmt.Sprintf("import private key error: %s", err))
	}

	eosApi.SetSigner(keyBag)
}

type GetTableOpts struct {
	LowerBound string `json:"lower_bound,omitempty"`
	UpperBound string `json:"upper_bound,omitempty"`
	Limit      uint32 `json:"limit,omitempty"`          // defaults to 10 => chain_plugin.hpp:struct get_table_rows_params
	KeyType    string `json:"key_type,omitempty"`       // The key type of --index, primary only supports (i64), all others support (i64, i128, i256, float64, float128, ripemd160, sha256). Special type 'name' indicates an account name.
	Index      string `json:"index_position,omitempty"` // Index number, 1 - primary (first), 2 - secondary index (in order defined by multi_index), 3 - third index, etc. Number or name of index can be specified, e.g. 'secondary' or '2'.
	EncodeType string `json:"encode_type,omitempty"`    // The encoding type of key_type (i64 , i128 , float64, float128) only support decimal encoding e.g. 'dec'" "i256 - supports both 'dec' and 'hex', ripemd160 and sha256 is 'hex' only
}

func GetTableRows(code, scope, table string, opts *GetTableOpts, rows interface{}) (err error) {
	rv := reflect.ValueOf(rows)
	if rv.Kind() != reflect.Ptr || rv.IsNil() || rv.Elem().Kind() != reflect.Slice {
		return errors.New(fmt.Sprintf("rows must be pointer of slice, but: %v", rv.Type()))
	}

	if opts == nil {
		opts = &GetTableOpts{}
	}
	req := eos.GetTableRowsRequest{
		Code:       code,
		Scope:      scope,
		Table:      table,
		JSON:       true,
		LowerBound: opts.LowerBound,
		UpperBound: opts.UpperBound,
		Limit:      opts.Limit,
		KeyType:    opts.KeyType,
		Index:      opts.Index,
		EncodeType: opts.EncodeType,
	}
	resp, err := eosApi.GetTableRows(req)
	RecordEosRpcCall(GetTable, err, &RpcOpts{Table: table})
	if err != nil {
		return
	}

	err = resp.JSONToStructs(rows)
	return
}

func PushActionByAdmin(actionName string, actionData interface{}) (txId string, err error) {
	rv := reflect.ValueOf(actionData)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		err = errors.New(fmt.Sprintf("actionData must be pointer and not nil, but: %v", rv.Type()))
		return
	}

	txOpts := &eos.TxOptions{}
	if err = txOpts.FillFromChain(eosApi); err != nil {
		return
	}

	action := &eos.Action{
		Account: eos.AN(eosVaultContract),
		Name:    eos.ActN(actionName),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AN(eosAdminAccount), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(actionData),
	}

	tx := eos.NewTransaction([]*eos.Action{action}, txOpts)
	_, packedTx, err := eosApi.SignTransaction(tx, txOpts.ChainID, eos.CompressionNone)
	if err != nil {
		return
	}

	resp, err := eosApi.PushTransaction(packedTx)
	RecordEosRpcCall(PushAction, err, &RpcOpts{Action: actionName})
	if err != nil {
		err = errors.New(fmt.Sprintf("%s push transaction: %s", actionName, err))
		return
	}
	return resp.TransactionID, nil
}
