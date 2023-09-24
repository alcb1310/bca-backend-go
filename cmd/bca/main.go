package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	port, portRead := os.LookupEnv("PORT")
	if !portRead {
		godotenv.Load()
		port, portRead = os.LookupEnv("PORT")
		if !portRead {
			log.Panic("Unable to load environment variables")
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Home endpoint hit")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	})
	log.Println(":INFO: Server listening in port: ", port)
	log.Panic(http.ListenAndServe(":"+port, nil))
}
