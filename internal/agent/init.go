package agent

import (
    "fmt"
    "log"

    model "github.com/hf-chow/tofu/internal/model"
)

func InitAgent(name, port string) (Agent, error){
    var agent Agent

    switch name {
    case "or":
        log.Println("Initializing OR Agent")
        agent.Topic = ""
    case "qa":
        log.Println("Initializing QA Agent")
        agent.Topic = "quest_ans"
    case "pq":
        log.Println("Initializing PQ Agent")
        agent.Topic = "prod_query"
    case "ps":
        log.Println("Initializing PS Agent")
        agent.Topic = "prod_search"
    default:
        log.Fatalf("Valid agents [qa, pq, ps, or]")
    }

    fmt.Println("Starting agent client...")

    m, err := model.ServeOllamaModel("gemma3:1b", port, false)
    if err != nil {
        log.Fatalf("failed to serve model %s: %s", m.Name, err)
        return agent, err
    }

    agent.Name = name
    agent.Model = m

    log.Printf(" [*] Serving model %s on endpoint %s", m.Name, m.Endpoint)

    return agent, nil
}

func (a *Agent) setContext() {
    switch a.Name {
    case "or":
        a.ContextPrompt = ORContextPrompt
    case "qa":
        a.ContextPrompt = QAContextPrompt
    case "ps":
        a.ContextPrompt = PSContextPrompt
    case "pq":
        a.ContextPrompt = PQContextPrompt
    }
}
