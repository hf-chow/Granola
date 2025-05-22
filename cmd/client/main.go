package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	agent "github.com/hf-chow/tofu/internal/agent"
	model "github.com/hf-chow/tofu/internal/model"
	"github.com/hf-chow/tofu/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const connString = "amqp://guest:guest@localhost:5672/"

func parse() (string, string, string) {
    if len(os.Args) < 4 {
        log.Printf("Usage: %s [model_provider] [agent_name] [port]...", os.Args[0])
        os.Exit(0)
    }

    if (os.Args[1] != "vllm-cpu" && os.Args[1] != "ollama") {
        log.Println("Valid model provider: [ollama, vllm-cpu]")
        os.Exit(0)
    }

    if (os.Args[2] != "qa" && os.Args[2] != "pq" && os.Args[2] != "ps") {
        log.Println("Valid agent names: [qa, pq, ps]")
        os.Exit(0)
    }

    portStr := os.Args[3]
    port, err := strconv.Atoi(portStr)
    if err != nil {
        log.Println("Port number has to be numeric intergers within [1024 - 49151]")
        os.Exit(0)
    }

    if (port < 1024 || port > 49151) {
        log.Println("Valid port numbers: 1024 - 49151")
        os.Exit(0)
    }
    return os.Args[0], os.Args[1], os.Args[3]
}

func main() {
    conn, err := amqp.Dial(connString)
    if err != nil {
        log.Fatalf("failed to connect to Rabbitmq: %s", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("failed to create a channel: %s", err)
    }
    defer ch.Close()

    provider, agentName, agentPort := parse()

    models := []model.Model{
        model.NewOllamaModel("gemma3:1b", agentPort, false),
        model.NewVLLMModel("google/gemma-3-1b-it", agentPort, "cpu"),
    }

    var m model.Model
    if provider == "ollama" {
        m := models[0]
        err = m.Start()
        if err != nil {
            log.Fatalf("failed to start model: %s", err)
        }

    } else if provider == "vllm-cpu" {
        m := models[1]
        err = m.Start()
        if err != nil {
            log.Fatalf("failed to start model: %s", err)
        }
    }
   
    agent, err := agent.InitAgent(agentName, m, ch)
    if err != nil {
        log.Fatalf("failed initialize agent %s: %s", agentName, err)
    }

    err = ch.ExchangeDeclare(
        agent.Topic,    // name
        "topic",        // type
        true,           // durable
        false,          // auto-deleted
        false,          // internal
        false,          // no-wait
        nil,            // arguments
    )
    if err != nil {
        log.Fatalf("failed to declare an exchange: %s", err)
    }

    q, err := ch.QueueDeclare(
        agent.Topic,    // name
        false,          // duarable
        false,          // delete when unused
        false,          // exclusive
        false,          // no-wait
        nil,            // arguments
    )
    if err != nil {
        log.Fatalf("failed to declare a queue :%s", err)
    }

    log.Printf("Binding queue %s to exchange %s with routing key '#'", 
                q.Name, 
                agentName)

    err = ch.QueueBind(
        q.Name,         // queue name
        "#",            // routing key
        agent.Topic,    // exchange
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("failed to bind a queue: %s", err)
    }

    msgs, err := ch.Consume(
        q.Name,         // queue
        "",             // consumer
        true,           // auto-ack
        false,          // exclusive
        false,          // no-local
        false,          // no-wait
        nil,            // args
    )
    if err != nil {
        log.Fatalf("failed to register a consumer: %s", err)
    }

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    shutdown := make(chan struct{})

    go func() {
        for {
            select {
            case d, ok := <-msgs: 
                if !ok {
                    log.Fatal("Message channel closed")
                    return 
                }
                log.Printf(" [x] receive prompt %s", d.Body)
                modelResp, err := agent.Respond(d.Body)
                if err != nil {
                    log.Printf(" [x] model failed to generate a response: %s", err)
                    continue
                }
                log.Printf(" [x] model response: %s", modelResp.Response)
                switch agent.Name{
                case "pq":
                    agent.SendDown(
                        []byte(modelResp.Response),
                        routing.AgentQATopic,
                    )
                case "ps":
                    agent.SendDown(
                        []byte(modelResp.Response),
                        routing.AgentPQTopic,
                    )
                }
            case <- shutdown:
                log.Println("Shutting down consumer...")
                return
            }
        }
    }()

    log.Printf(" [*] Waiting for prompts. To exit press CTRL+C")
    <- sigChan

    log.Printf("Recieved interrupt, shutting down...")
    close(shutdown)
    model.StopOllamaService()
}

