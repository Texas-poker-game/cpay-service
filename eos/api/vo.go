package api

type WithdrawRequest struct {
	Sid    string  `form:"sid"`
	Amount float64 `form:"amount"`
	To     string  `form:"to"`
}

type WithdrawVo struct {
	Sid    string  `json:"sid"`
	Amount float64 `json:"amount"`
	To     string  `json:"to"`
}
