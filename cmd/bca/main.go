package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Home endpoint hit")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	})
	log.Println(":INFO: Server listening in port: ", port)
	log.Panic(http.ListenAndServe(":"+port, nil))
}
