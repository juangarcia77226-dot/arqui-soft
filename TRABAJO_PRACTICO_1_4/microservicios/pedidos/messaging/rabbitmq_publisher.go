package messaging

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const colaPedidosConfirmados = "pedidos-confirmados"

// RabbitMQPublisher publica en RabbitMQ los eventos de pedidos confirmados,
// en la cola "pedidos-confirmados", para que logística (conceptual) los
// pueda tomar más adelante.
type RabbitMQPublisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQPublisher(amqpURI string) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	if _, err := ch.QueueDeclare(colaPedidosConfirmados, true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQPublisher{conn: conn, ch: ch}, nil
}

func (p *RabbitMQPublisher) PublicarPedidoConfirmado(evento PedidoConfirmado) error {
	body, err := json.Marshal(evento)
	if err != nil {
		return err
	}

	err = p.ch.PublishWithContext(
		context.Background(),
		"",                     // exchange por defecto
		colaPedidosConfirmados, // routing key = nombre de la cola
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("[rabbitmq] evento publicado en %q: %s", colaPedidosConfirmados, body)
	return nil
}

func (p *RabbitMQPublisher) Close() {
	if p.ch != nil {
		p.ch.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
