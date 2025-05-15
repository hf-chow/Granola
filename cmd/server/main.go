package main

import (
    "context"
    "log"
    "os"
    "strings"
    "time"

    amqp "github.com/rabbitmq/amqp091-go"
)

func bodyFrom(args []string) string {
    var s string
    if (len(args) < 3) || os.Args[2] == "" {
        s = "hello"
    } else {
        s = strings.Join(args[2:], " ")
    }
    return s
}

func severityFrom(args []string) string {
    var s string
    if (len(args) < 2) || os.Args[1] == "" {
        s = "anonymous.info"
    } else {
        s = os.Args[1]
    }
    return s
}

func publish(ch *amqp.Channel) {
    err := ch.ExchangeDeclare(
        "logs_topic",              // name
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

    body := bodyFrom(os.Args)
    err = ch.PublishWithContext(
        ctx,
        "logs_topic",              // exchange
        severityFrom(os.Args),      // routing key
        false,                      // mandatory
        false,                      // immediate
        amqp.Publishing {
            ContentType: "text/plain",
            Body:        []byte(body),
        })
    if err != nil {
        log.Fatalf("failed to publish a message: %s", err)
    }
    log.Printf(" [x] Sent %s\n", body)
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
    publish(ch)
}
