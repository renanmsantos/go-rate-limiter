package main

import (
	"fmt"
	"net/http"

	"github.com/renanmsantos/go-rate-limiter/pkg"
)

func requestHandle(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {

	http.HandleFunc("/", pkg.RateLimiterMiddleware(requestHandle))

	fmt.Println("Server is listening on port 8080.")
	http.ListenAndServe(":8080", nil)
}
