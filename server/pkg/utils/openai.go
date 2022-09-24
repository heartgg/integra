package utils

import (
	"context"
	"fmt"
	"os"
	"strings"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/joho/godotenv"
)

func GetAiKey() string {
	err := godotenv.Load();
	if err != nil {
		fmt.Printf("Error loading .env\n");
		return "";
	}
	return os.Getenv("OPEN_AI_KEY");
}

// Given a string like "Lung Cancer", returns string array of possible tests.
func AskAI(diagnosis string, testList [] string, testListStr string) []string {

	c := gogpt.NewClient(GetAiKey()) // from .env
	ctx := context.Background()

	prompt := "Given "+testListStr+"what are the best exams for a patient with "+diagnosis+"?";
	fmt.Println("\n\n\nThe prompt is ",prompt,"\n");
	req := gogpt.CompletionRequest{
		Model: "text-babbage-001",
		MaxTokens: 120,
		Prompt:    prompt,
		Temperature: 0.19,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return nil;
	}
	
	fmt.Println("The response from OpenAI is ",resp.Choices[0].Text);
	fmt.Println("\nThe matches Are!!!")

	matchList := make([]string,0);
	noMatchList := make([]string,0);

	lowerResp := strings.ToLower(resp.Choices[0].Text);
	for i := 0; i < len(testList); i++ {
		if (strings.Contains(lowerResp,strings.ToLower(testList[i]))) {
			fmt.Println(testList[i]);
			matchList = append(matchList,testList[i]);
		} else {
			noMatchList = append(noMatchList,testList[i]);
		}
	}

	fmt.Println("\n",matchList,noMatchList);
	
	return nil;
}