package agent

import (
    model "github.com/hf-chow/tofu/internal/model"
)

func (a *Agent) Respond(p []byte) (model.ModelResponse, error) {
    context := a.ContextPrompt
    prompt := "Here is the context:\n" + context + "\n" +
    "Here is the query: \n" + string(p) +
    "Please reply with the given context"

    modelResp, err := a.Model.Prompt([]byte(prompt))
    if err != nil {
        return model.ModelResponse{}, err
    }
    return modelResp, err
}


