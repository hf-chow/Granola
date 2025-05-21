package agent

import (
	//"fmt"

	model "github.com/hf-chow/tofu/internal/model"
	"github.com/hf-chow/tofu/internal/pubsub"
)

func (a *Agent) Respond(p []byte) (string, error) {
    context := a.ContextPrompt
    prompt := "Here is the context:\n" + context + "\n" +
    "Here is the user query: \n" + string(p) + "\n" +
    "Please reply with the given context\n"

    //fmt.Printf(" [DEBUG] Prompt: %s", prompt)
    modelResp, err := model.Prompt(prompt, a.Model)
    if err != nil {
        return "", err
    }
    return modelResp, err
}

func (a *Agent) SendDown(prompt, topic string){
     pubsub.PublishText(
        topic,
        []byte(prompt),
        a.Channel,
    )
}
