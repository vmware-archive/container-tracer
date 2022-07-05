package main

import (
	"log"

	api "gitlab.eng.vmware.com/opensource/tracecruncher-api/api"
)

func main() {
	log.Printf("Server started")

	router := api.NewRouter()

	router.Run("localhost:8080")
}
