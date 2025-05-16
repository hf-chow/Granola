package agent

import "github.com/hf-chow/tofu/internal/model"

type Agent struct{
    Name            string
    Topic           string 
    ContextPrompt   string
    Model           model.Model
}

type OrchestrationOutcome int

const (
    OrchestrationOutcomeQA OrchestrationOutcome = iota
    OrchestrationOutcomePQ
    OrchestrationOutcomePS
)
