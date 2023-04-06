package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"tests-selector/openai"

	"github.com/joho/godotenv"
)

var Envs = GetEnvs()

func GetEnvs() map[string]string {
	envFile := os.Args[1]
	envs, _ := godotenv.Read(envFile)
	return envs
}

func main() {
	prDescription := os.Args[2]

	// Load test cases from file.
	testCases := loadTestCasesFromFile("test_cases.txt")

	// Generate embeddings for the test cases.
	testCaseEmbeddings := make([][][]float32, len(testCases))
	for i, testCase := range testCases {
		// Generate embedding for the test case using OpenAI API.
		resp, err := openai.Embeddings(openai.EmbeddingsRequest{
			Model: "text-embedding-ada-002",
			Input: testCase,
		}, fmt.Sprintf("Bearer %s", Envs["TOKEN"]))
		if err != nil {
			log.Fatal(err)
		}

		testCaseEmbeddings[i] = make([][]float32, len(resp.Data))
		if len(resp.Data) != len(testCase) {
			log.Fatal("Error response size doesnt match")
		}
		for _, dt := range resp.Data {
			testCaseEmbeddings[i][dt.Index] = make([]float32, len(dt.Embedding))
			testCaseEmbeddings[i][dt.Index] = dt.Embedding
		}
	}

	// Compute cosine similarity between the PR description and each test case.
	prDescriptionEmbedding, err := generateEmbedding(prDescription)
	if err != nil {
		log.Fatal(err)
	}

	testCaseSimilarity := make([]float32, len(testCases))
	mostSimilar := 0
	for i, testCase := range testCases {
		for j, _ := range testCase {
			cosineSim := cosineSimilarity(prDescriptionEmbedding, testCaseEmbeddings[i][j])
			testCaseSimilarity[i] = cosineSim
			if testCaseSimilarity[i] > testCaseSimilarity[mostSimilar] {
				mostSimilar = i
			}
			// fmt.Printf("Test case %d: %s, cosine similarity: %f\n", i+1, t, cosineSim)
		}
	}

	testsCasesIndex := []string{"account", "login", "logout", "pos-sale", "customer"}

	fmt.Println(testsCasesIndex[mostSimilar])
}

// Load test cases from a file.
func loadTestCasesFromFile(filename string) [][]string {
	// TODO: Implement file loading.
	return [][]string{
		{"Testing purchase flow"},
		{"User log in for the first time"},
		{"Logout redirect to home page"},
		{"Creating a comment"},
		{"Reviewing the product"},
	}
}

// Generate an embedding for a text prompt using OpenAI API.
func generateEmbedding(prompt string) ([]float32, error) {
	resp, err := openai.Embeddings(openai.EmbeddingsRequest{
		Model: "text-embedding-ada-002", //text-embedding-ada-002
		Input: []string{prompt},
	}, fmt.Sprintf("Bearer %s", Envs["TOKEN"]))
	if err != nil {
		return nil, err
	}

	if len(resp.Data) != 1 {
		return nil, errors.New("Error response size doesn't match")
	}

	embeddingFloats := resp.Data[0].Embedding
	return embeddingFloats, nil
}

// Compute cosine similarity between two vectors.
// Range  ranges from -1 to 1,
// with -1 indicating that two vectors are completely dissimilar,
// 0 indicating that they are orthogonal (i.e., have no relationship),
// and 1 indicating that they are identical or very similar.
func cosineSimilarity(vec1, vec2 []float32) float32 {
	// Calculate the dot product of the two vectors
	dotProduct := float32(0)
	for i := 0; i < len(vec1); i++ {
		dotProduct += vec1[i] * vec2[i]
	}

	// Calculate the magnitude of the first vector
	vec1Magnitude := float32(0)
	for _, v := range vec1 {
		vec1Magnitude += v * v
	}
	vec1Magnitude = float32(math.Sqrt(float64(vec1Magnitude)))

	// Calculate the magnitude of the second vector
	vec2Magnitude := float32(0)
	for _, v := range vec2 {
		vec2Magnitude += v * v
	}
	vec2Magnitude = float32(math.Sqrt(float64(vec2Magnitude)))

	// Calculate the cosine similarity of the two vectors
	if vec1Magnitude == 0 || vec2Magnitude == 0 {
		return 0
	}
	return dotProduct / (vec1Magnitude * vec2Magnitude)
}
