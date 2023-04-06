package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type EmbeddingsRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
	User  string   `json:"user"`
}

type EmbeddingsResponse struct {
	Object string `json:"object"`
	Data   []data `json:"data"`
	Model  string `json:"model"`
	Usage  usage  `json:"usage"`
}

type data struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

type usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type EmbeddingsErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    int    `json:"code"`
	} `json:"error"`
}

func Embeddings(req EmbeddingsRequest, token string) (resp EmbeddingsResponse, err error) {

	reqBytes, _ := json.Marshal(req)

	r, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(reqBytes))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", token)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Println("Error:", string(body))
		err = errors.New("Request returned != 200")
		return
	}

	err = json.Unmarshal(body, &resp)
	return
}
