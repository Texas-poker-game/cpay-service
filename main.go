package main

import (
	"cpay/eos"
	"cpay/server"
)

func main() {
	go eos.HandleDeposits()
	go eos.HandleWithdraws()
	go eos.HandleAuths()

	server.Start()
}
