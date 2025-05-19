package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	agent "github.com/hf-chow/tofu/internal/agent"
	model "github.com/hf-chow/tofu/internal/model"
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

    agent, err := agent.InitAgent("or", "11111", ch)
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
                outcome, err := agent.Orchestrate([]byte(input))
                if err != nil {
                    log.Fatalf("failed to orchestrate task: %s", err)
                }
                fmt.Println(outcome)
                fmt.Print("> ")
            }
        }
    }()
    <-sigChan
    log.Println("Shutting down server...")
    model.StopOllamaService()
}
