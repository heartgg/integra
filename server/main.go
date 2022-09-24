package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/heartgg/integri-scan/server/pkg/routes"
	"github.com/heartgg/integri-scan/server/pkg/utils"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env\n")
		return
	}
	err = utils.InitAI()
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	routes.SetupRoutes()
	port := os.Getenv("PORT")
	fmt.Printf("Listening on port: %v", port)
	http.ListenAndServe(":"+port, nil)
}
