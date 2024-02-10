package main

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/ws", initiateWebSocket)
	handler := cors.Default().Handler(mux)
	port := ":3000"
	fmt.Printf("server listening at port %s ", port)
	http.ListenAndServe(port, handler)
}
