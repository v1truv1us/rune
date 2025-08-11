package copilot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CopilotRequest is the payload for the Copilot LLM API
// (This is a minimal example; expand as needed)
type CopilotRequest struct {
	Messages []map[string]string `json:"messages"`
}

type CopilotResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// CallCopilotLLM calls the Copilot LLM API with the given prompt and token.
// The model selection is handled by Copilot; there is no model parameter.
func CallCopilotLLM(prompt, token string) (string, error) {
	url := "https://api.githubcopilot.com/chat/completions"
	payload := CopilotRequest{
		Messages: []map[string]string{{
			"role":    "user",
			"content": prompt,
		}},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Copilot API error: %s", string(b))
	}

	var copilotResp CopilotResponse
	if err := json.NewDecoder(resp.Body).Decode(&copilotResp); err != nil {
		return "", err
	}
	if len(copilotResp.Choices) == 0 {
		return "", fmt.Errorf("No choices returned from Copilot")
	}
	return copilotResp.Choices[0].Message.Content, nil
}

// TODO: If Copilot ever adds a model parameter, add it here.
