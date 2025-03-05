package main

import (
	"log"
	"miniJupyter/protocol"
	"miniJupyter/zmq/base"
	"miniJupyter/zmq/mode"
	"time"
)

type Client struct {
	dealer *mode.DealerNode
	config *base.Config
}

func NewClient(configPath string) (*Client, error) {
	// 加载配置
	config, err := base.LoadConfig(configPath)
	log.Printf("Load config: %+v", config)
	if err != nil {
		return nil, err
	}

	// 创建Dealer节点
	dealer, err := mode.NewDealer(config.Zmq.DealerAddress)
	log.Println("Create dealer:", dealer)
	if err != nil {
		return nil, err
	}

	return &Client{
		dealer: dealer,
		config: config,
	}, nil
}

func (c *Client) SendExecuteRequest() error {
	// 创建执行请求消息
	msg, err := protocol.NewMessageBuilder().
		WithType(protocol.MsgTypeExecuteRequest).
		WithSession("test-session").
		WithUser("test-user").
		WithTransport(protocol.TransportZMQ).
		WithContent(&protocol.ExecuteRequestContent{
			CommandId: protocol.GenerateUUID(),
			Service:   "test-service",
			Method:    "test-method",
			Params: map[string]interface{}{
				"param1": "value1",
				"param2": 42,
			},
			Timeout: 30,
		}).
		WithTraceHop("client-1", "jupyter-client", "localhost").
		Build()

	if err != nil {
		return err
	}
	log.Println("Create message")

	// 序列化消息
	msgData, err := protocol.SerializeMessage(msg)
	if err != nil {
		return err
	}
	log.Printf("Serialize message: %s", msgData)

	// 发送消息
	err = c.dealer.SendToServer(msgData)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}
	log.Println("Send message to server")

	// 接收响应
	responseData, err := c.dealer.ReceiveFromServer()
	if err != nil {
		return err
	}

	// 解析响应
	response, err := protocol.ParseMessage(responseData)
	if err != nil {
		return err
	}

	// 打印响应
	log.Printf("Received response: %+v\n", response)
	return nil
}

func (c *Client) Close() {
	c.dealer.Close()
}

func main() {
	client, err := NewClient("../../../zmq/config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	log.Println("Client starting...")
	defer client.Close()

	// 等待服务器启动
	time.Sleep(time.Second)
	log.Println("Client started")

	// 发送测试请求
	err = client.SendExecuteRequest()
	log.Println("Send request to server")
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
}