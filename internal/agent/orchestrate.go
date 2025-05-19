package agent

import (
	"fmt"
	"strings"

	pubsub "github.com/hf-chow/tofu/internal/pubsub"
	"github.com/hf-chow/tofu/internal/routing"
)

func (agent *Agent) Orchestrate(prompt []byte) (OrchestrationOutcome, error) {
    var outcome OrchestrationOutcome

    modelResp, err:= agent.Respond(prompt)
    if err != nil {
        return OrchestrationOutcomeQA, err
    }

    fmt.Printf("[DEBUG] model decision is: %s", modelResp.Response)

    switch strings.TrimSpace(modelResp.Response) {
    case "QA":
        outcome =  OrchestrationOutcomeQA
        //fmt.Println("[DEBUG] model interpret the query as QA")
        pubsub.PublishText(routing.AgentQATopic, prompt, agent.Channel)
        return outcome, nil
    case "PQ":
        outcome =  OrchestrationOutcomePQ
        //fmt.Println("[DEBUG] model interpret the query as PQ")
        pubsub.PublishText(routing.AgentPQTopic, prompt, agent.Channel)
        return outcome, nil
    case "PS":
       outcome =  OrchestrationOutcomePS
       //fmt.Println("[DEBUG] model interpret the query as PS")
       pubsub.PublishText(routing.AgentPSTopic, prompt, agent.Channel)
       return outcome, nil
    default:
        outcome = OrchestrationOutcomeQA
        //fmt.Println("[DEBUG] model defaulted to QA")
        pubsub.PublishText(routing.AgentQATopic, prompt, agent.Channel)
        return outcome, nil
    }
}
