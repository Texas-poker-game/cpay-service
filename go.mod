module cpay

go 1.12

require (
	github.com/eoscanada/eos-go v0.8.16
	github.com/gin-gonic/gin v1.4.0
	github.com/jinzhu/gorm v1.9.10
	github.com/prometheus/client_golang v0.9.3
	github.com/sirupsen/logrus v1.4.2
	github.com/tidwall/gjson v1.3.2 // indirect
	github.com/tidwall/sjson v1.0.4 // indirect
	github.com/zsais/go-gin-prometheus v0.1.0
	queding.com/go/common v0.0.0
)

replace queding.com/go/common => ../queding-common
