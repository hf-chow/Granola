package model

type Model struct {
	Name	            string
    Endpoint            string
    Stream              bool
}

type ModelRequest struct {
    Model               string          `json:"model"`
    Prompt              string          `json:"prompt"`
    Stream              bool            `json:"stream"`
}

type ModelResponse struct {
	Model               string          `json:"model"`
	CreatedAt           string          `json:"created_at"`
	Response            string          `json:"response"`
	Done                bool            `json:"done"`
	Context             []int           `json:"context"`
	TotalDuration       int             `json:"total_duration"`
	LoadDuration        int             `json:"load_duration"`
	PromptEvalCount     int             `json:"prompt_eval_count"`
	PromptEvalDuration  int             `json:"prompt_eval_duration"`
	EvalCount           int             `json:"eval_count"`
	EvalDuration        int             `json:"eval_duration"`
}

