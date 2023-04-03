package main

import (
	"context"
	"fmt"
	"github.com/openai/api"
	"log"
	"strings"
)

func main() {
	// Set up OpenAI API credentials.
	apiKey := "YOUR_API_KEY"
	client, err := api.NewClient(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	// Load test cases from file.
	testCases := loadTestCasesFromFile("test_cases.txt")

	// Generate embeddings for the test cases.
	testCaseEmbeddings := make([][]float32, len(testCases))
	for i, testCase := range testCases {
		// Generate embedding for the test case using OpenAI API.
		resp, err := client.Completions(context.Background(), &api.CompletionRequest{
			Model: "text-davinci-002",
			Prompt: testCase,
			MaxTokens: 64,
			N: 1,
			Temperature: 0.0,
		})
		if err != nil {
			log.Fatal(err)
		}
		embedding := resp.Choices[0].Text
		embedding = strings.TrimSpace(embedding)
		embeddingValues := strings.Split(embedding, ",")
		testCaseEmbeddings[i] = make([]float32, len(embeddingValues))
		for j, val := range embeddingValues {
			fmt.Sscanf(val, "%f", &testCaseEmbeddings[i][j])
		}
	}

	// Compute cosine similarity between the PR description and each test case.
	prDescription := "Fix bug in login flow"
	prDescriptionEmbedding, err := generateEmbedding(client, prDescription)
	if err != nil {
		log.Fatal(err)
	}
	for i, testCase := range testCases {
		cosineSim := cosineSimilarity(prDescriptionEmbedding, testCaseEmbeddings[i])
		fmt.Printf("Test case %d: %s, cosine similarity: %f\n", i+1, testCase, cosineSim)
	}
}

// Load test cases from a file.
func loadTestCasesFromFile(filename string) []string {
	// TODO: Implement file loading.
	return []string{
		"Test case 1",
		"Test case 2",
		"Test case 3",
	}
}

// Generate an embedding for a text prompt using OpenAI API.
func generateEmbedding(client *api.Client, prompt string) ([]float32, error) {
	resp, err := client.Completions(context.Background(), &api.CompletionRequest{
		Model: "text-davinci-002",
		Prompt: prompt,
		MaxTokens: 64,
		N: 1,
		Temperature: 0.0,
	})
	if err != nil {
		return nil, err
	}
	embedding := resp.Choices[0].Text
	embedding = strings.TrimSpace(embedding)
	embeddingValues := strings.Split(embedding, ",")
	embeddingFloats := make([]float32, len(embeddingValues))
	for i, val := range embeddingValues {
		fmt.Sscanf(val, "%f", &embeddingFloats[i])
	}
	return embeddingFloats, nil
}

// Compute cosine similarity between two vectors.
func cosineSimilarity(vec1, vec2 []float32) float32 {
	var dotProduct float32
	var mag1, mag2 float32
	for i := range vec1 {
		dotProduct += vec1[i] * vec2[i]
	}