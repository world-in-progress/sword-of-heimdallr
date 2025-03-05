package mode

import (
	"encoding/hex"
	"fmt"
	"log"
	"miniJupyter/zmq/base"

	zmq "github.com/pebbe/zmq4"
)

// RouterNode 扩展基础的 ZmqNode，添加 Router 特定功能
type RouterNode struct {
	*base.ZmqNode
}

// DealerNode 扩展基础的 ZmqNode，添加 Dealer 特定功能
type DealerNode struct {
	*base.ZmqNode
}

// NewRouter 创建并返回 RouterNode
func NewRouter(address string) (*RouterNode, error) {
	node, err := base.NewZmqNode(zmq.ROUTER, address, true)
	if err != nil {
		return nil, err
	}
	return &RouterNode{node}, nil
}

// NewDealer 创建并返回 DealerNode
func NewDealer(address string) (*DealerNode, error) {
	node, err := base.NewZmqNode(zmq.DEALER, address, false)
	if err != nil {
		return nil, err
	}
	return &DealerNode{node}, nil
}

// SendToClient 发送消息给特定客户端
func (r *RouterNode) SendToClient(clientID string, msg string) error {
	// 将十六进制字符串转换回二进制身份标识符
	id, err := hex.DecodeString(clientID)
	if err != nil {
		return fmt.Errorf("invalid client ID format: %w", err)
	}
	return r.Send(string(id), msg)
}

// ReceiveFromClient 接收来自客户端的消息，返回客户端ID和消息
func (r *RouterNode) ReceiveFromClient() (clientID string, msg string, err error) {
	msgs, err := r.Receive()
	if err != nil {
		return "", "", err
	}
	log.Println("ReceiveFromClient:", msgs)
	if len(msgs) < 2 {
		return "", "", fmt.Errorf("invalid message format")
	}
	// 将二进制身份标识符转换为十六进制字符串表示
	clientID = fmt.Sprintf("%x", []byte(msgs[0]))
	return clientID, msgs[1], nil
}

// SendToServer 发送消息给特定服务端
func (d *DealerNode) SendToServer(msg string) error {
	return d.Send(msg)
}

// ReceiveFromServer 接收来自服务端的消息
func (d *DealerNode) ReceiveFromServer() (msg string, err error) {
	msgs, err := d.Receive()
	if err != nil {
		return "", err
	}
	return msgs[0], nil
}
