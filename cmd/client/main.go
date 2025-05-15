package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
    model "github.com/hf-chow/tofu/internal/model"
)

func main() {
    m, err := model.ServeOllamaModel("gemma3:1b", false)
    if err != nil {
        log.Fatalf("failed to serve model %s: %s", m.Name, err)
    }
    log.Printf("serving model %s on endpoint %s", m.Name, m.Endpoint)

    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
        "logs_topic",   // name
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
        true,          // exclusive
        false,          // no-wait
        nil,            // arguments
    )
    if err != nil {
        log.Fatalf("failed to declare a queue :%s", err)
    }

    if len(os.Args) < 2 {
        log.Printf("Usage: %s [binding_key]...", os.Args[0])
        os.Exit(0)
    }

    for _, s := range os.Args[1:] {
        log.Printf("Binding queue %s to exchange %s with routing key %s", 
                    q.Name, "logs_topic", s)
        err = ch.QueueBind(
            q.Name,         // queue name
            s,              // routing key
            "logs_topic",   // exchange
            false,
            nil,
        )
        if err != nil {
            log.Fatalf("failed to bind a queue: %s", err)
        }

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

    var forever chan struct{}

    go func() {
        for d := range msgs {
            log.Printf(" [x] %s", d.Body)
        }
    }()

    log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
    <- forever
}

