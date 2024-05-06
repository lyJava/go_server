package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %v", msg, err)
	}
}

// ConnectRabbitmq 连接RabbitMQ
func ConnectRabbitmq(connectUrl string) *amqp.Connection {
	conn, err := amqp.Dial(connectUrl)
	failOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

// OpenChannel 创建信道
func OpenChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	return ch
}

// DeclareQueue 申明队列
func DeclareQueue(ch *amqp.Channel, queueName string) *amqp.Queue {
	queue, err := ch.QueueDeclare(queueName, // 队列名称
		false, // 是否持久化
		false, // 是否自动删除
		false, // 是否排他
		false, // 是否阻塞
		nil,   // 额外参数
	)
	failOnError(err, "Failed to declare a queue")
	return &queue
}

func DeclareExchange(ch *amqp.Channel, exchangeName, typeStr string) {
	err := ch.ExchangeDeclare(
		exchangeName, // name
		typeStr,      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")
}

func BindQueue(ch *amqp.Channel, queueName, routeKey, exchangeName string) {
	err := ch.QueueBind(
		queueName,    // queue name
		routeKey,     // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")
}

// PublishMessage 发布消息
func PublishMessage(ch *amqp.Channel, exchange, queueName string, message interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bytes, _ := json.Marshal(message)
	err := ch.PublishWithContext(
		ctx,
		exchange,  // 交换机名称
		queueName, // 路由键
		false,     // 是否强制
		false,     // 是否立即
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytes,
		})

	failOnError(err, "Failed to publish a message: %v")
	fmt.Println("RabbitMQ Message send successfully ===", string(bytes))
}

// Consume 消费
func Consume(ch *amqp.Channel, queueName string, autoAsk bool, done chan struct{}) []string {
	delivery, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		autoAsk,   // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	var messages []string

	//var forever chan struct{}
	go func() {
		/*for d := range msg {
			log.Printf(" [x] %s", d.Body)
			d.Ack(false)
		}*/

		for msg := range delivery {
			log.Printf(" [x] %s", msg.Body)
			messages = append(messages, string(msg.Body))
			// 根据 autoAck 参数决定是否手动确认消息
			if !autoAsk {
				// 手动确认消息
				err := msg.Ack(false)
				if err != nil {
					log.Printf("Failed to acknowledge message: %v", err)
					// 如果确认消息失败，可以选择进行日志记录或其他错误处理
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	// 等待通道关闭，以确保消费者协程完成
	<-done
	return messages
}
