package agent

import (
    "fmt"
    "log"

    model "github.com/hf-chow/tofu/internal/model"
)

func InitAgent(name, port string) (Agent, error){
    var a Agent

    switch name {
    case "or":
        log.Println("Initializing OR Agent")
        a.Topic = ""
    case "qa":
        log.Println("Initializing QA Agent")
        a.Topic = "quest_ans"
    case "pq":
        log.Println("Initializing PQ Agent")
        a.Topic = "prod_query"
    case "ps":
        log.Println("Initializing PS Agent")
        a.Topic = "prod_search"
    default:
        log.Fatalf("Valid as [qa, pq, ps, or]")
    }

    fmt.Println("Starting a client...")

    m, err := model.ServeOllamaModel("gemma3:1b", port, false)
    if err != nil {
        log.Fatalf("failed to serve model %s: %s", m.Name, err)
        return a, err
    }

    a.Name = name
    a.Model = m
    a.setContext()

    log.Printf(" [*] Serving model %s on endpoint %s", m.Name, m.Endpoint)

    return a, nil
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
