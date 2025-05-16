package agent

import (
    "fmt"

	model "github.com/hf-chow/tofu/internal/model"
)

func (a *Agent) Respond(p []byte) (model.ModelResponse, error) {
    context := a.ContextPrompt
    prompt := "Here is the context:\n" + context + "\n" +
    "Here is the user query: \n" + string(p) + "\n" +
    "Please reply with the given context\n"

    fmt.Printf(" [DEBUG] Prompt: %s", prompt)
    modelResp, err := a.Model.Prompt([]byte(prompt))
    if err != nil {
        return model.ModelResponse{}, err
    }
    return modelResp, err
}


