package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jamieyoung5/pubroulette-api/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	r := mux.NewRouter()
	r.HandleFunc("/pub", api.Handler).Methods(http.MethodGet)
	r.Use(CORSMiddleware)

	log.Println("Server is running on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
