package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alcb1310/bca-backend-go/internals/routes"
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

	r := routes.NewRouter()

	log.Println(":INFO: Server listening in port: ", port)
	log.Panic(http.ListenAndServe(":"+port, r))
}
