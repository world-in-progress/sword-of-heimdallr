package mode

import (
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

// NewRouter 创建并返回 RouterNode (不需要身份标识)
func NewRouter(address string, bind bool) (*RouterNode, error) {
	node, err := base.NewZmqNode(zmq.ROUTER, address, bind, "")  // 空身份标识
	if err != nil {
		return nil, err
	}
	return &RouterNode{node}, nil
}

// NewDealer 创建并返回 DealerNode (需要身份标识)
func NewDealer(address string, bind bool, identity string) (*DealerNode, error) {
	if identity == "" {
		return nil, fmt.Errorf("dealer requires an identity")
	}
	node, err := base.NewZmqNode(zmq.DEALER, address, bind, identity)
	if err != nil {
		return nil, err
	}
	return &DealerNode{node}, nil
}

// SendToClient 发送消息给特定客户端
func (r *RouterNode) SendToClient(clientID string, msg string) error {
	log.Printf("SendToClient to %s", clientID)
	return r.Send(clientID, msg)
}

// ReceiveFromClient 接收来自客户端的消息，返回客户端ID和消息
func (r *RouterNode) ReceiveFromClient() (clientID string, msg string, err error) {
	msgs, err := r.Receive()
	if err != nil {
		return "", "", err
	}
	
	if len(msgs) < 2 {
		return "", "", fmt.Errorf("invalid message format")
	}
	return msgs[0], msgs[1], nil
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
