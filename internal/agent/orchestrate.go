package agent

func (agent *Agent) Orchestrate(p []byte) (OrchestrationOutcome, error) {
    var outcome OrchestrationOutcome

    modelResp, err:= agent.Model.Prompt(p)
    if err != nil {
        return OrchestrationOutcomeQA, err
    }

    switch modelResp.Response {
    case "QA":
        outcome =  OrchestrationOutcomeQA
        return outcome, nil
    case "PQ":
        outcome =  OrchestrationOutcomePQ
        return outcome, nil
    case "PS":
        outcome =  OrchestrationOutcomePS
        return outcome, nil
    default:
        outcome = OrchestrationOutcomeQA
        return outcome, nil
    }
}
