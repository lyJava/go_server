package controller

import (
	"apiProject/api/expressAPI/rabbitmq"
	"apiProject/api/expressAPI/types"
	"apiProject/api/response"
	"encoding/json"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"io"
	"log"
	"net/http"
)

// OrderController 订单控制器
type OrderController struct {
	rabbitmqConn *amqp.Connection //rabbitmq连接
}

// OrderControllerInit 订单控制器初始化
func OrderControllerInit(conn *amqp.Connection) *OrderController {
	return &OrderController{
		rabbitmqConn: conn,
	}
}

func (mq *OrderController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/order/create", mq.CreateHandler).Methods("POST")
}

var queueName = "create_order_queue"
var exchangeName = "create_order_exchange"
var routingKey = "create_order_routing_key"

// CreateHandler 创建订单
func (mq *OrderController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var order types.Order

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("新增失败"))
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, &order)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	// 创建 RabbitMQ 信道
	channel := rabbitmq.OpenChannel(mq.rabbitmqConn)
	// 声明交换器
	rabbitmq.DeclareExchange(channel, exchangeName, "direct")
	// 声明队列
	queue := rabbitmq.DeclareQueue(channel, queueName)
	// 绑定交换器与队列
	rabbitmq.BindQueue(channel, queue.Name, routingKey, exchangeName)
	// 发布订单消息到交换器
	rabbitmq.PublishMessage(channel, exchangeName, routingKey, order)

	// 用于通知消费者协程已完成消费
	done := make(chan struct{})
	defer close(done) // 确保在函数退出时关闭通道

	// 启动消费者协程
	go func() {
		defer func() {
			// 发送完成通知
			done <- struct{}{}
		}()

		// 用于接收消息的通道
		delivery := make(chan amqp.Delivery)

		// 使用协程接收消息

		msgs, err := channel.Consume(
			queueName, // queue
			"",        // consumer
			false,     // auto ack
			false,     // exclusive
			false,     // no local
			false,     // no wait
			nil,       // args
		)
		if err != nil {
			log.Printf("Failed to register a consumer: %v", err)
			return
		}

		// 将消息发送到 delivery 通道中
		go func() {
			for msg := range msgs {
				delivery <- msg
			}
		}()

		// 开始消费消息
		for msg := range delivery {
			log.Printf(" 接收到到消息=== %s", msg.Body)
			// 手动确认消息
			err := msg.Ack(false)
			if err != nil {
				log.Printf("Failed to acknowledge message: %v", err)
				// 如果确认消息失败，可以选择进行日志记录或其他错误处理
			}
		}
	}()

	response.WriteJson(w, response.OkMessageResp("订单消息发送成功"))
	return
}
