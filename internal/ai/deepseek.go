package aiutils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sapopinguino/internal/config"
	"time"

	"github.com/valyala/fastjson"
)

var DeepSeekClient *http.Client

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type ChatRequest struct {
	Messages         []Message      `json:"messages"`
	Model            string         `json:"model"`
	FrequencyPenalty int            `json:"frequency_penalty"`
	MaxTokens        int            `json:"max_tokens"`
	PresencePenalty  int            `json:"presence_penalty"`
	ResponseFormat   ResponseFormat `json:"response_format"`
	Stop             interface{}    `json:"stop"`
	Stream           bool           `json:"stream"`
	StreamOptions    interface{}    `json:"stream_options"`
	Temperature      float64        `json:"temperature"`
	TopP             float64        `json:"top_p"`
	Tools            interface{}    `json:"tools"`
	ToolChoice       string         `json:"tool_choice"`
	Logprobs         bool           `json:"logprobs"`
	TopLogprobs      interface{}    `json:"top_logprobs"`
}

func newChatRequest(systemRole, userInput string) ChatRequest {
	return ChatRequest{
        Messages: []Message{
            {Role: "system", Content: systemRole},
            {Role: "user", Content: userInput},
        },
		Model:            "deepseek-chat",
		FrequencyPenalty: 0,
		MaxTokens:        4096,
		PresencePenalty:  0,
		ResponseFormat:   ResponseFormat{Type: "json_object"},
		Stop:             nil,
		Stream:           false,
		StreamOptions:    nil,
		Temperature:      1,
		TopP:             1,
		Tools:            nil,
		ToolChoice:       "none",
		Logprobs:         false,
		TopLogprobs:      nil,
	}
}

func ChatCompletionDeepSeek(context context.Context, system_role string, input string) (*Output, error) {
    url := "https://api.deepseek.com/chat/completions"
    method := "POST"

    req_body := newChatRequest(system_role, input) 

	payload_bytes, err := json.Marshal(req_body)
	if err != nil {
		log.Println("Error marshaling request body in ChatCompletion():", err)
		return nil, err
	}

	payload := bytes.NewReader(payload_bytes)
    if err != nil {
		log.Println("Error reading bytes in ChatCompletion()", err)
        return nil, err
    }

    DeepSeekClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest(method, url, payload)
    if err != nil {
		log.Println("Error creating request in ChatCompletion():", err)
		return nil, err
	}
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Authorization", "Bearer " + config.C.OpenAI.Key)

    res, err := DeepSeekClient.Do(req)
    defer res.Body.Close()
    if err != nil {
        log.Println("Error making an http request in ChatCompletion(): %s", err)
        return nil, err
    }

    body, err := io.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }

    var p fastjson.Parser
	v, err := p.ParseBytes(body)
	if err != nil {
		log.Println("Parse error while reading http response in ChatCompletion(): %s", err)
		return nil, err
	}

    content := v.GetStringBytes("choices", "0", "message", "content")
	if content == nil {
        fmt.Println("Error getting content from response in ChatCompletion()")
		return nil, err
	}

    content_string := string(content)

    return &content_string, nil
}
