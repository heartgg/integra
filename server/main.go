package main

import (
	"fmt"
	"net/http"

	"github.com/heartgg/integri-scan/server/pkg/routes"
)

func main() {
	routes.SetupRoutes()
	fmt.Println("IntegriScan websocket server is listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
