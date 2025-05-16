package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func parse() (string, string) {
    if len(os.Args) < 3 {
        log.Printf("Usage: %s [topic] [prompt]", os.Args[0])
    }
    if (os.Args[1] != "quest_ans" && os.Args[1] != "prod_query" && os.Args[1] != "prod_search") {
        log.Println("Valid topics: [quest_ans, prod_query, prod_search]")
    }

    return os.Args[1], os.Args[2]
}

func makeBody(args []string) string {
    var s string
    if (len(args) < 2) || os.Args[1] == "" {
        s = "hello"
    } else {
        s = strings.Join(args[1:], " ")
    }
    return s
}

func publish(topic, prompt string, ch *amqp.Channel) {
    err := ch.ExchangeDeclare(
        topic,                      // name
        "topic",                    // type
        true,                       // durable
        false,                      // auto-deleted
        false,                      // internal
        false,                      // no-wait
        nil,                        // arugments
    )
    if err != nil {
        log.Fatalf("failed to declare an exchange: %s", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err = ch.PublishWithContext(
        ctx,
        topic,                      // exchange
        "",                         // routing key
        false,                      // mandatory
        false,                      // immediate
        amqp.Publishing {
            ContentType: "text/plain",
            Body:        []byte(prompt),
        })
    if err != nil {
        log.Fatalf("failed to publish a message: %s", err)
    }
    log.Printf(" [x] Prompt sent: %s\n", prompt)
}

func repl() {
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
                modelResp, err := agent.Model.Prompt(d.Body)
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

func main() {
    conn, err :=  amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("failed to connect to RabbitMQ: %s", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("failed to create a channel: %s", err)
    }
    defer ch.Close()

    topic, prompt := parse()
    publish(topic, prompt, ch)
}
