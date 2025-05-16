package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	agent "github.com/hf-chow/tofu/internal/agent"
	model "github.com/hf-chow/tofu/internal/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

const connString = "amqp://guest:guest@localhost:5672/"

func parse() (string, string) {
    if len(os.Args) < 3 {
        log.Printf("Usage: %s [agent_name] [port]...", os.Args[0])
        os.Exit(0)
    }

    if (os.Args[1] != "qa" && os.Args[1] != "pq" && os.Args[1] != "ps") {
        log.Println("Valid agent names: [qa, pq, ps]")
        os.Exit(0)
    }

    portStr := os.Args[2]
    port, err := strconv.Atoi(portStr)
    if err != nil {
        log.Println("Port number has to be numeric intergers within [1024 - 49151]")
        os.Exit(0)
    }

    if (port < 1024 || port > 49151) {
        log.Println("Valid port numbers: 1024 - 49151")
        os.Exit(0)
    }
        
    return os.Args[1], os.Args[2]
}

func initAgent(name, port string) (agent.Agent, error){
    var agent agent.Agent
    agent.Name = name

    switch agent.Name {
    case "qa":
        log.Println("Initializing QA Agent")
        agent.Topic = "quest_ans"
    case "pq":
        log.Println("Initializing PQ Agent")
        agent.Topic = "prod_query"
    case "ps":
        log.Println("Initializing PS Agent")
        agent.Topic = "prod_search"
    }

    fmt.Println("Starting agent client...")

    m, err := model.ServeOllamaModel("gemma3:1b", port, false)
    if err != nil {
        log.Fatalf("failed to serve model %s: %s", m.Name, err)
        return agent, err
    }

    agent.Model = m

    log.Printf(" [*] Serving model %s on endpoint %s", m.Name, m.Endpoint)

    return agent, nil
}

func main() {
    agentName, agentPort := parse()
    agent, err := initAgent(agentName, agentPort)
    if err != nil {
        log.Fatalf("failed initialize agent %s: %s", agentName, err)
    }

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
        "",             // name
        false,          // duarable
        false,          // delete when unused
        true,           // exclusive
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
        agent.Topic,   // exchange
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
            case <- shutdown:
                log.Println("Shutting down consumer...")
                return
            }
        }
    }()

    log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
    <- sigChan

    log.Printf("Recieved interrupt, shutting down...")
    close(shutdown)
    model.StopOllamaService()
}

