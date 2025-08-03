package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	apiURL       = "http://<NODE-IP>:30081/v1/chat/completions"
	concurrency  = 5  
	totalRequests = 20
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func sendRequest() (int, time.Duration) {
	reqBody := ChatRequest{
		Model: "qwen/qwen3-32B",
		Messages: []ChatMessage{
			{Role: "user", Content: "Write a short poem about the ocean."},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	start := time.Now()

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Request error: %v", err)
		return 0, 0
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var chatResp ChatResponse
	_ = json.Unmarshal(body, &chatResp)

	duration := time.Since(start)
	return chatResp.Usage.CompletionTokens, duration
}

func main() {
	var wg sync.WaitGroup
	reqChan := make(chan struct{}, concurrency)

	var totalTokens int
	var totalTime time.Duration
	var mu sync.Mutex

	startTest := time.Now()

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		reqChan <- struct{}{} // concurrency limit

		go func() {
			defer wg.Done()
			tokens, duration := sendRequest()

			mu.Lock()
			totalTokens += tokens
			totalTime += duration
			mu.Unlock()

			<-reqChan
		}()
	}

	wg.Wait()
	elapsed := time.Since(startTest)

	tps := float64(totalTokens) / elapsed.Seconds()

	fmt.Println("```markdown")
	fmt.Println("## vLLM Load Test Results\n")
	fmt.Println("| Metric | Value |")
		fmt.Println("|--------|-------|")
	fmt.Printf("| Total Requests | %d |\n", totalRequests)
	fmt.Printf("| Total Tokens | %d |\n", totalTokens)
	fmt.Printf("| Total Time (s) | %.2f |\n", elapsed.Seconds())
	fmt.Printf("| Tokens per Second | %.2f |\n", tps)
	fmt.Println("```")
}
