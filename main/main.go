package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		panic(err)
	}
	model := client.GenerativeModel("gemini-1.5-pro")

	resp, err := model.GenerateContent(ctx, genai.Text("Explain how AI works"))
	fmt.Println(resp.Candidates[0].Content)
}
