package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/heartgg/integri-scan/server/pkg/routes"
)

func main() {
	routes.SetupRoutes()
	port := os.Getenv("PORT")
	fmt.Printf("Listening on port: %v", port)
	http.ListenAndServe(":"+port, nil)
}
