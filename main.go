package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"tests-selector/filesystem"
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
	descriptionInput := os.Args[2]

	vector, _ := filesystem.ReadFloatArrayFromCSV("result.csv")

	prDescriptionEmbedding, err := generateEmbedding(descriptionInput)
	if err != nil {
		log.Fatal(err)
	}

	testCaseSimilarity := make([]float64, len(vector))
	for i, testCase := range vector {
		cosineSim := cosineSimilarity(prDescriptionEmbedding, testCase)
		fmt.Printf("Test case %d: cosine similarity: %f\n", i+1, cosineSim)
		testCaseSimilarity[i] = cosineSim
	}

	fmt.Printf("FINAL RESULT %+v", mostSimilar(0.8, testCaseSimilarity))
}

func mostSimilar(threshold float64, testCaseSimilarity []float64) []int {
	result := make([]int, 0)
	for i, v := range testCaseSimilarity {
		if v >= threshold {
			result = append(result, i)
		}
	}
	return result
}

func embeddingsToCSV() {
	files := filesystem.GetFileList(`.\testcases`)

	// iterate over the files in the directory
	for _, file := range files {
		scenarios := filesystem.FeaturesToSingleLine(`.\testcases`, file)
		for _, scenario := range scenarios {
			vector, _ := generateEmbedding(scenario)
			filesystem.AppendToCSV("result.csv", file.Name(), scenario, vector)
		}
	}
}

// Generate an embedding for a text prompt using OpenAI API.
func generateEmbedding(prompt string) ([]float64, error) {
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
func cosineSimilarity(vec1, vec2 []float64) float64 {
	// Calculate the dot product of the two vectors
	dotProduct := float64(0)
	for i := 0; i < len(vec1); i++ {
		dotProduct += vec1[i] * vec2[i]
	}

	// Calculate the magnitude of the first vector
	vec1Magnitude := float64(0)
	for _, v := range vec1 {
		vec1Magnitude += v * v
	}
	vec1Magnitude = float64(math.Sqrt(float64(vec1Magnitude)))

	// Calculate the magnitude of the second vector
	vec2Magnitude := float64(0)
	for _, v := range vec2 {
		vec2Magnitude += v * v
	}
	vec2Magnitude = float64(math.Sqrt(float64(vec2Magnitude)))

	// Calculate the cosine similarity of the two vectors
	if vec1Magnitude == 0 || vec2Magnitude == 0 {
		return 0
	}
	return dotProduct / (vec1Magnitude * vec2Magnitude)
}
