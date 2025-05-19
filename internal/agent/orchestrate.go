package agent

import (
	"fmt"
	"strings"

    pubsub "github.com/hf-chow/tofu/internal/pubsub"
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
        fmt.Println("[DEBUG] model interpret the query as QA")
        pubsub.PublishText("quest_ans", prompt, agent.Channel)
        return outcome, nil
    case "PQ":
        outcome =  OrchestrationOutcomePQ
        fmt.Println("[DEBUG] model interpret the query as PQ")
        pubsub.PublishText("prod_query", prompt, agent.Channel)
        return outcome, nil
    case "PS":
       outcome =  OrchestrationOutcomePS
       fmt.Println("[DEBUG] model interpret the query as PS")
       pubsub.PublishText("prod_search", prompt, agent.Channel)
       return outcome, nil
    default:
        outcome = OrchestrationOutcomeQA
        fmt.Println("[DEBUG] model defaulted to QA")
        pubsub.PublishText("quest_ans", prompt, agent.Channel)
        return outcome, nil
    }
}
