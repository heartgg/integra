package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	gogpt "github.com/sashabaranov/go-gpt3"
)

var c *gogpt.Client

func InitAI() error {
	env := os.Getenv("OPEN_AI_KEY")
	if env == "" {
		return errors.New("no OPEN_AI_KEY in .env")
	}
	c = gogpt.NewClient(env)

	return nil
}

// Given a string like "Lung Cancer", returns a JSON-Encoded Map.
// Example possible return:
// {"Angiography":0,"Arthrography":0,"Bone Density Scan":1,"Bone XRAY":0,"Chest XRAY":1,"Crystogram":0,"Fluoroscopy":1,"Mammography":0,"Myelography":0,"Skull Radiography":0,"Virtual Colonoscopy":1}
func AskAI(diagnosis string, testList []string, testListStr string) (string, error) {

	ctx := context.Background()

	prompt := "Given " + testListStr + "what are the best exams for a patient with " + diagnosis + "?"
	fmt.Println("\n\n\nThe prompt is ", prompt, "\n")
	req := gogpt.CompletionRequest{
		Model:       "text-babbage-001",
		MaxTokens:   120,
		Prompt:      prompt,
		Temperature: 0.19,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	fmt.Println("The response from OpenAI is ", resp.Choices[0].Text)
	fmt.Println("\nThe matches Are!!!")

	matchMap := make(map[string]int)

	lowerResp := strings.ToLower(resp.Choices[0].Text)
	for i := 0; i < len(testList); i++ {
		if strings.Contains(lowerResp, strings.ToLower(testList[i])) {
			fmt.Println(testList[i])
			matchMap[testList[i]] = 1
		} else {
			matchMap[testList[i]] = 0
		}
	}

	fmt.Println("\n", matchMap)
	mapJson, err := json.Marshal(matchMap)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("\n", string(mapJson))
	return string(mapJson), nil
}
