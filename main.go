package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/renanmoreirasan/go-rate-limiter/app/middleware"
	"github.com/renanmoreirasan/go-rate-limiter/infra/config"
)

func main() {

	config.LoadEnvConfigs()
	config.LoadCache()

	mux := http.NewServeMux()
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request permitted!")
		w.Write([]byte("Request permitted!"))
	})
	mux.Handle("/", middleware.RateLimiterMiddleware(finalHandler))

	fmt.Println("Server is listening on port 8080.")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)

}
