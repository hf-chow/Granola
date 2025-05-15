package model

type Model struct {
	Name	            string
    Endpoint            string
    Stream              bool
}

type ChatRequest struct {
    Model               string          `json:"model"`
    Messages            []ChatMessage   `json:"messages"`
    Stream              bool            `json:"stream"`
}

type ChatMessage struct {
    Role                string          `json:"role"`
    Content             string          `json:"content"`
}

type ChatResponse struct {
    Model               string          `json:"Model"`
    CreatedAt           string          `json:"created_at"`
    Message             ChatMessage     `json:"message"`
    DoneReason          string          `json:"done_reason"`
    Done                bool            `json:"done"`
    TotalDuration       int             `json:"total_duration"`
    LoadDuration        int             `json:"load_duration"`
    PromptEvalCount     int             `json:"prompt_eval_count"`
    PromptEvalDuration  int             `json:"prompt_eval_duration"`
    EvalCount           int             `json:"eval_count"`
    EvalDuration        int             `json:"eval_duration"`
}
