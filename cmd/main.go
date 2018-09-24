package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tada3/triton/handler"
)

const (
	PORT = 10900
)

func main() {
	log.Println("You won't get lost if you take the road to the left and bear to the left at every crossroad...")
	log.Printf("Listen on port %d\n", PORT)

	fileServer := http.FileServer(http.Dir("resources"))
	http.Handle("/resources/", http.StripPrefix("/resources/", fileServer))
	http.HandleFunc("/meiroya", handler.Dispatch)
	http.HandleFunc("/monitor/l7check", handler.HealthCheck)

	addr := fmt.Sprintf(":%d", PORT)
	log.Fatalln(http.ListenAndServe(addr, nil))
}
