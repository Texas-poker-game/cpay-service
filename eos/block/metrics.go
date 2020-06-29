package block

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"queding.com/go/common/config"
)

type RpcType string

const (
	GetTable   RpcType = "get_table"
	PushAction RpcType = "push_action"
)

var (
	metricsAddress = config.GetString("server.metrics.address")
	rpcCounterVec  = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cpay_eos_rpc_call_total",
		Help: "The total number of EOS Rpc call",
	}, []string{"type", "success", "table", "action"})
)

type RpcOpts struct {
	Table  string
	Action string
}

func RecordEosRpcCall(rpcType RpcType, err error, opts *RpcOpts) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error(r)
			}
		}()

		labels := prometheus.Labels{
			"type":    string(rpcType),
			"success": fmt.Sprintf("%v", err == nil),
		}
		if opts == nil {
			opts = &RpcOpts{}
		}
		labels["table"] = opts.Table
		labels["action"] = opts.Action

		rpcCounterVec.With(labels).Inc()
	}()
}
