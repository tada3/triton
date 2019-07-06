package main

import (
	"fmt"
	"os"

	"net/http"

	"github.com/tada3/triton/handler"
	"github.com/tada3/triton/logging"
)

const (
	port = 10700
)

var (
	log *logging.Entry
)

func init() {
	log = logging.NewEntry("main")
}

func main() {
	log.Info("Oni ha soto!")

	log.Info("Listen on port %d", port)

	fileServer := http.FileServer(http.Dir("resources"))
	http.Handle("/resources/", http.StripPrefix("/resources/", fileServer))
	http.HandleFunc("/cek", handler.Dispatch)
	http.HandleFunc("/monitor/l7check", handler.HealthCheck)

	addr := fmt.Sprintf(":%d", port)

	err := http.ListenAndServe(addr, nil)
	log.Error("HTTP server stopped!", err)
	os.Exit(1)
}
