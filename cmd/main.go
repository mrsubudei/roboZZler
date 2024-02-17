package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mrsubudei/roboZZler/internal/server"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/solve", server.SolveRobozzle)
	fmt.Println("server started at 8080 port")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
