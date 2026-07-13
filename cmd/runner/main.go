package main

import (
	"log"
	"net/http"
	"os"

	"github.com/qianqiu0404/wallet-reliability-lab/internal/lab"
)

func main() {
	server, err := lab.NewServerFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	address := os.Getenv("LAB_RUNNER_ADDR")
	if address == "" {
		address = "127.0.0.1:8787"
	}
	log.Printf("Wallet Reliability Lab runner listening on http://%s", address)
	if err := http.ListenAndServe(address, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
