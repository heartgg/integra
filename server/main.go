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
	godotenv.Load()
	err := utils.InitAI()
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	routes.SetupRoutes()
	port := os.Getenv("PORT")
	fmt.Printf("Listening on port: %v\n", port)
	http.ListenAndServe(":"+port, nil)
}
