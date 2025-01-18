package main

import (
	"github.com/gorilla/mux"
	"github.com/jamieyoung5/pooblet/api"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	r := mux.NewRouter()
	// REST best practices dictate that the endpoint should be a noun i.e /pub
	r.HandleFunc("/getPub", api.Handler).Methods(http.MethodGet)
	r.Use(api.CORSMiddleware)

	log.Println("Server is running on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
