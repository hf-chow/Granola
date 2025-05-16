package pubsub

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Publish(topic string, prompt []byte, ch *amqp.Channel) {
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
            Body:        prompt,
        })
    if err != nil {
        log.Fatalf("failed to publish a message: %s", err)
    }
    log.Printf(" [x] Prompt sent: %s\n", prompt)
}

