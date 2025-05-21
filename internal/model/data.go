package model

type OllamaModelRequest struct {
    Model               string                      `json:"model"`
    Prompt              string                      `json:"prompt"`
    Stream              bool                        `json:"stream"`
}

type OllamaModelResponse struct {
	Model               string                      `json:"model"`
	CreatedAt           string                      `json:"created_at"`
	Response            string                      `json:"response"`
	Done                bool                        `json:"done"`
	Context             []int                       `json:"context"`
	TotalDuration       int                         `json:"total_duration"`
	LoadDuration        int                         `json:"load_duration"`
	PromptEvalCount     int                         `json:"prompt_eval_count"`
	PromptEvalDuration  int                         `json:"prompt_eval_duration"`
	EvalCount           int                         `json:"eval_count"`
	EvalDuration        int                         `json:"eval_duration"`
}

type OpenAIModelRequest struct {
    Model               string                      `json:"model"`
    Messages            []OpenAIModelMessages       `json:"messages"`
}

type OpenAIModelMessages struct {
    Role                string                      `json:"user"`
    Content             string                      `json:"content"`
}

type OpenAIModelResponse struct {
    ID                  string                      `json:"id"`
	Model               string                      `json:"model"`
	CreatedAt           string                      `json:"created_at"`
    Object              string                      `json:"object"`
    Output              OpenAIModelResponseOutput   `json:"output"`          
}

type OpenAIModelResponseOutput struct {
    ID                  string                      `json:"id"`
    Content             OpenAIModelOutputContent    `json:"output"`     
    Role                string                      `json:"assitant"`
    Type                string                      `json:"message"`
}

type OpenAIModelOutputContent struct {
    Text                string                      `json:"text"` 
    Type                string                      `json:"type"` 
}
