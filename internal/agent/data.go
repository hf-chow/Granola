package agent

import (
	"github.com/hf-chow/tofu/internal/model"
	"github.com/rabbitmq/amqp091-go"
)

type Agent struct{
    Name            string
    Topic           string 
    ContextPrompt   string
    Model           model.Model
    Channel         *amqp091.Channel
}

type OrchestrationOutcome int

const (
    OrchestrationOutcomeQA OrchestrationOutcome = iota
    OrchestrationOutcomePQ
    OrchestrationOutcomePS
)
