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
	r.HandleFunc("/getPub", api.Handler).Methods(http.MethodGet)
	r.Use(api.CORSMiddleware)

	log.Println("Server is running on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
