package main

import (
	"bufio"
	"context"
    "fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hf-chow/tofu/internal/agent"
	amqp "github.com/rabbitmq/amqp091-go"
)

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

    agent, err := agent.InitAgent("or", "11111")
    if err != nil {
        log.Fatalf("failed initialize agent %s: %s", agent.Name, err)
    }

    shutdown := make(chan struct{})
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        scanner := bufio.NewScanner(os.Stdin)
        fmt.Print("> ")
        for scanner.Scan() {
            select {
            case <- shutdown:
                log.Println("Shutting down...")
                return 
            default:
                input := strings.TrimSpace(scanner.Text())
                decision := agent.Orchestrate([]byte(input))
                fmt.Print("> ")
            }
        }
    }()

    topic, prompt := parse()
    publish(topic, prompt, ch)
}
