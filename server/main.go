package main

import (
	"fmt"
	"net/http"

	"github.com/heartgg/integri-scan/server/pkg/routes"
	"github.com/heartgg/integri-scan/server/pkg/utils"
)

func main() {

	testList := []string{"Angiography", "Arthrography", "Bone Density Scan", "Bone XRAY", "Chest XRAY", "Crystogram", "Fluoroscopy", "Mammography", "Myelography", "Skull Radiography", "Virtual Colonoscopy"};
	var testListStr string = "";

	for i := 0; i < len(testList); i++ {
		testListStr = testListStr + testList[i] + ", ";
	}

	routes.SetupRoutes()
	utils.AskAI("Cervical Cancer", testList, testListStr);
	fmt.Println("IntegriScan websocket server is listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
