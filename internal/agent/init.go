package agent

import (
	"fmt"
	"log"

	model "github.com/hf-chow/tofu/internal/model"
	"github.com/hf-chow/tofu/internal/routing"
	"github.com/rabbitmq/amqp091-go"
)

func InitAgent(name, port string, ch *amqp091.Channel) (*Agent, error){
    var topic string
    switch name {
    case "or":
        log.Println("Initializing OR Agent")
        topic = ""
    case "qa":
        log.Println("Initializing QA Agent")
        topic = routing.AgentQATopic
    case "pq":
        log.Println("Initializing PQ Agent")
        topic = routing.AgentPQTopic
    case "ps":
        log.Println("Initializing PS Agent")
        topic = routing.AgentPSTopic
    default:
        log.Fatalf("Valid as [qa, pq, ps, or]")
    }

    fmt.Println("Starting a client...")

    m, err := model.ServeOllamaModel("gemma3:1b", port, false)
    if err != nil {
        log.Fatalf("failed to serve model %s: %s", m.Name, err)
        return &Agent{}, err
    }

    a := &Agent{
        Name:           name,
        Topic:          topic,
        ContextPrompt:  getContext(name),
        Channel:        ch,
        Model:          m,
    }

    log.Printf(" [*] Serving model %s on endpoint %s", m.Name, m.Endpoint)

    return a, nil
}

func getContext(name string) string{
    switch name {
    case "or":
        return ORContextPrompt
    case "qa":
        return QAContextPrompt
    case "ps":
        return PSContextPrompt
    case "pq":
        return PQContextPrompt
    }
    return ORContextPrompt 
}
