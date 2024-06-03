package main

import (
	"log"
	"net/http"
	"pdrive/pipedrive-test-api/handlers"
)

func main() {
	http.HandleFunc("/deals", handlers.RequestMetricsMiddleware(handlers.Handler))
	http.HandleFunc("/deals/", handlers.RequestMetricsMiddleware(handlers.Handler))
	http.HandleFunc("/metrics", handlers.MetricsHandler)
	http.HandleFunc("/", handlers.InvalidPathHandler)
	log.Println("Server listening on localhost:8080..")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
