package main

import (
	"fmt"
	"log"
	"miniJupyter/protocol"
	"miniJupyter/zmq/base"
	"miniJupyter/zmq/mode"
	"time"
)

type Core struct {
	router *mode.RouterNode
	config *base.Config
}

func NewCore(configPath string) (*Core, error) {
	// 加载配置
	config, err := base.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 创建Router节点
	router, err := mode.NewRouter(config.Zmq.RouterAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create router: %w", err)
	}

	return &Core{
		router: router,
		config: config,
	}, nil
}

func (c *Core) handleExecuteRequest(msg *protocol.Message) *protocol.Message {
	// 解析执行请求内容
	content, ok := msg.Content.(*protocol.ExecuteRequestContent)
	if !ok {
		return createErrorResponse(msg, "Invalid content type")
	}

	// 模拟执行任务
	log.Printf("Executing command: %s on service: %s\n", content.CommandId, content.Service)
	time.Sleep(time.Second) // 模拟处理时间

	// 创建响应消息
	response := protocol.NewMessageBuilder().
		WithType(protocol.MsgTypeExecuteReply).
		WithSession(msg.Header.SessionId).
		WithUser(msg.Header.UserId).
		WithTransport(protocol.TransportZMQ).
		WithParentMessage(msg).
		WithContent(&protocol.ExecuteReplyContent{
			Status: protocol.StatusSuccess,
		}).
		WithTraceHop("core-1", "jupyter-core", "localhost")

	result, err := response.Build()
	if err != nil {
		return createErrorResponse(msg, err.Error())
	}

	return result
}

func createErrorResponse(parentMsg *protocol.Message, errMsg string) *protocol.Message {
	response, _ := protocol.NewMessageBuilder().
		WithType(protocol.MsgTypeExecuteReply).
		WithSession(parentMsg.Header.SessionId).
		WithUser(parentMsg.Header.UserId).
		WithTransport(protocol.TransportZMQ).
		WithParentMessage(parentMsg).
		WithContent(&protocol.ExecuteReplyContent{
			Status: protocol.StatusError,
			Error:  errMsg,
		}).
		Build()
	return response
}

func (c *Core) Run() error {
	log.Println("Core server starting...")
	defer c.router.Close()

	for {
		// 接收消息
		clientID, msgData, err := c.router.ReceiveFromClient()
		if err != nil {
			log.Printf("Error receiving message: %v\n", err)
			continue
		}

		// 解析消息
		msg, err := protocol.ParseMessage([]byte(msgData))
		if err != nil {
			log.Printf("Error parsing message: %v\n", err)
			continue
		}

		// 处理消息
		var response *protocol.Message
		switch msg.Header.MsgType {
		case protocol.MsgTypeExecuteRequest:
			response = c.handleExecuteRequest(msg)
		default:
			response = createErrorResponse(msg, "Unsupported message type")
		}

		// 序列化响应
		responseData, err := protocol.SerializeMessage(response)
		if err != nil {
			log.Printf("Error serializing response: %v\n", err)
			continue
		}

		// 发送响应
		err = c.router.SendToClient(clientID, string(responseData))
		if err != nil {
			log.Printf("Error sending response: %v\n", err)
		}
	}
}

func main() {
	core, err := NewCore("../../zmq/config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to create core: %v", err)
	}

	if err := core.Run(); err != nil {
		log.Fatalf("Core error: %v", err)
	}
}