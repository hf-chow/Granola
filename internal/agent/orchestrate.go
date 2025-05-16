package agent

import (
	"fmt"
	"log"
	"strings"
)

func (agent *Agent) Orchestrate(p []byte) (OrchestrationOutcome, error) {
    var outcome OrchestrationOutcome

    modelResp, err:= agent.Respond(p)
    if err != nil {
        return OrchestrationOutcomeQA, err
    }

    fmt.Printf("[DEBUG] model response: %s", modelResp.Response)

    switch strings.TrimSpace(modelResp.Response) {
    case "QA":
        outcome =  OrchestrationOutcomeQA
        log.Println("The model interpret the query as QA")
        return outcome, nil
    case "PQ":
        outcome =  OrchestrationOutcomePQ
        log.Println("The model interpret the query as PQ")
        return outcome, nil
    case "PS":
       outcome =  OrchestrationOutcomePS
       log.Println("The model interpret the query as PS")
       return outcome, nil
    default:
        outcome = OrchestrationOutcomeQA
        log.Println("The model defaulted to QA")
        return outcome, nil
    }
}
