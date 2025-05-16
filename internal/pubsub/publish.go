package pubsub

import (
    "bytes"
	"context"
	"log"
	"time"
    "encoding/gob"
    "encoding/json"

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


func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	dat, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        dat,
	})
}

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(val)
	if err != nil {
		return err
	}
	return ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/gob",
		Body:        buffer.Bytes(),
	})
}
