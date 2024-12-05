package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Lesson struct {
	Topic        string
	Submodules   []string
	Proposition  string
	PastProposal string
}

func generateShortLesson(lesson *Lesson) string {
	rand.Seed(time.Now().UnixNano())

	// Check for overly broad topics
	broadTopics := []string{"science", "history", "math", "art", "life"}
	if contains(broadTopics, strings.ToLower(lesson.Topic)) {
		return fmt.Sprintf("Error: '%s' is too broad. Please be more specific.  Try 'photosynthesis' or 'the water cycle'.", lesson.Topic)
	}

	// Capitalize the topic
	lesson.Topic = strings.Title(lesson.Topic)

	var tailoredContent string
	if lesson.PastProposal != "" {
		tailoredContent = fmt.Sprintf("\nBuilding upon the previous suggestion of exploring '%s', let's delve into %s:\n", lesson.PastProposal, lesson.Topic)
	} else {
		tailoredContent = fmt.Sprintf("\nLet's explore %s!\n\n", lesson.Topic)
	}

	lessonContent := fmt.Sprintf("**%s Lesson: %s**\n%s", lesson.Topic, lesson.Topic, tailoredContent)

	// Gemini Client & Context for ALL Gemini Calls within the Function
	ctx := context.Background()
	c, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv(("GEMINI_API_KEY"))))
	if err != nil {
		panic(fmt.Sprintf("Could not initialize Gemini Client: %v", err)) // Improved error message
	}
	model := c.GenerativeModel("gemini-1.5-flash")

	// Generate submodules (minimum 2)
	numSubmodules := rand.Intn(3) + 2 // Generates 2, 3, or 4
	for i := 0; i < numSubmodules; i++ {
		submodule := fmt.Sprintf("%d. Submodule: [Submodule %d on %s -  Content Needed Here]\n", i+1, i+1, lesson.Topic)

		//Gemini Request to populate submodule
		prompt := genai.Text(fmt.Sprintf("Create a concise and informative submodule on the topic of %s. This is submodule number %v.", lesson.Topic, i+1))
		resp, err := model.GenerateContent(ctx, prompt)

		if err != nil {
			panic(fmt.Sprintf("Gemini call failed: %v", err)) // Improved error handling
		}
		if len(resp.Candidates) == 0 {
			panic(fmt.Sprintf("Gemini returned no candidates: %v", err))
		}

		submodule = fmt.Sprintf("%d. %s\n", i+1, resp.Candidates[0].Content)
		lessonContent += submodule
	}

	// Propose a new lesson topic -  NOW using model and ctx
	prompt := genai.Text(fmt.Sprintf("Propose a related lesson topic that builds upon the current lesson about %s. Be creative and specific.", lesson.Topic)) // More specific prompt
	resp, err := model.GenerateContent(ctx, prompt)
	if err != nil {
		panic(fmt.Sprintf("Gemini call failed: %v", err)) // Improved error handling
	}
	if len(resp.Candidates) == 0 {
		panic(fmt.Sprintf("Gemini returned no candidates: %v", err))
	}
	lesson.Proposition = fmt.Sprintf("%s", resp.Candidates[0].Content)

	lessonContent += fmt.Sprintf("\n**Next Lesson Suggestion:** %s", lesson.Proposition)

	return lessonContent
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func main() {
	// First Lesson
	firstLesson := &Lesson{Topic: "Go Programming", PastProposal: ""}
	firstLessonContent := generateShortLesson(firstLesson)
	fmt.Println(firstLessonContent)
	fmt.Println("=====================================================================")
	// Second Lesson, building on the first
	secondLesson := &Lesson{Topic: "Error Handling", PastProposal: firstLesson.Proposition}
	secondLessonContent := generateShortLesson(secondLesson)
	fmt.Println(secondLessonContent)

}
