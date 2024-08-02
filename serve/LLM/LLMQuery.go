package LLM

import (
	"MusicBot/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Object  string `json:"object"`
}

func LLMQuery(prompt string) (string, error) {

	url := "https://api.siliconflow.cn/v1/chat/completions"

	jsonReq := fmt.Sprintf("{\"model\":\"%s\",\"messages\":[{\"role\":\"user\",\"content\":\"%s\"}],\"stream\":false,\"max_tokens\":512,\"temperature\":%f,\"top_p\":0.7,\"top_k\":50,\"frequency_penalty\":0.5,\"n\":1}", config.Config.LLM.Model, prompt, config.Config.LLM.Temperature)

	payload := strings.NewReader(jsonReq)

	req, _ := http.NewRequest("POST", url, payload)

	var bearer = "Bearer " + config.Config.LLM.SiliconFlowToken
	req.Header.Add("Authorization", bearer)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var parsed Response

	err := json.Unmarshal(body, &parsed)
	if err != nil {
		return string(body), err
	}

	return parsed.Choices[len(parsed.Choices)-1].Message.Content, nil
}
