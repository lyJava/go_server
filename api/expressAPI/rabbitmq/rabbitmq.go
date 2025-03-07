package rabbitmq

import (
	"apiProject/api/expressAPI/config"
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func failOnError(err error, msg string) {
	if err != nil {
		zap.L().Sugar().Errorf("%s===%+v", msg, err)
	}
}

// ConnectRabbitmq 连接RabbitMQ
func ConnectRabbitmq() *amqp.Connection {
	rabbitmqCfg := config.EnvConfig.Rabbitmq
	if !rabbitmqCfg.Enable {
		return nil
	}

	connectUrl := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		rabbitmqCfg.Username,
		rabbitmqCfg.Password,
		rabbitmqCfg.Host,
		rabbitmqCfg.Port,
		rabbitmqCfg.VirtualHost,
	)

	conn, err := amqp.Dial(connectUrl)

	if err != nil {
		zap.L().Sugar().Errorf("Failed to connect to RabbitMQ===\n%+v", err)
	}
	zap.L().Sugar().Infof("success to connect to RabbitMQ ===%v", rabbitmqCfg)
	// 将json格式化输出
	rabbitmqProperties, _ := json.MarshalIndent(conn.Properties, "", "    ")
	zap.L().Sugar().Infof("Success to connect to RabbitMQ===\r\n%+v", string(rabbitmqProperties))
	return conn
}

// OpenChannel 创建信道
func OpenChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	zap.L().Sugar().Infof("Success to open a channel===%v", ch)
	return ch
}

// DeclareQueue 申明队列
//
// 参数
//
//		ch (*amqp.Channel) 通道
//		queueName (string) 队列名称
//	 durable	(bool) 是否持久化
func DeclareQueue(ch *amqp.Channel, queueName string, durable bool) *amqp.Queue {
	queue, err := ch.QueueDeclare(queueName, // 队列名称
		true,    // 是否持久化
		durable, // 是否自动删除
		false,   // 是否排他
		false,   // 是否阻塞
		nil,     // 额外参数
	)
	failOnError(err, "Failed to declare a queue")
	zap.L().Sugar().Infof("Success to declare a queue===%v", queue)
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
	zap.L().Sugar().Infof("Success to declare an exchange===%s", exchangeName)
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
	zap.L().Sugar().Infof("Success to bind a queue queueName:%s,routeKey:%s,exchangeName:%s", queueName, routeKey, exchangeName)
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
	zap.L().Sugar().Infof("RabbitMQ Message send successfully:%s", string(bytes))

}

// Consume 消费
func Consume(ch *amqp.Channel, queueName string, autoAsk bool) []string {
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
	var forever chan struct{}

	go func() {
		/*for d := range msg {
			log.Printf(" [x] %s", d.Body)
			d.Ack(false)
		}*/

		for msg := range delivery {
			zap.L().Sugar().Infof(" [x] %s", msg.Body)
			messages = append(messages, string(msg.Body))
			// 根据 autoAck 参数决定是否手动确认消息
			if !autoAsk {
				// 手动确认消息
				err := msg.Ack(false)
				if err != nil {
					zap.L().Sugar().Errorf("Failed to acknowledge message: %v", err)
					// 如果确认消息失败，可以选择进行日志记录或其他错误处理
				}
			}
		}
	}()

	zap.L().Sugar().Infof(" [*] Waiting for logs. To exit press CTRL+C")
	// 等待通道关闭，以确保消费者协程完成
	<-forever
	return messages
}
