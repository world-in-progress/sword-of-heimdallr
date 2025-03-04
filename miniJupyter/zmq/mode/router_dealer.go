package mode

import (
	"fmt"
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
func (d *DealerNode) SendToServer(serverID string, msg string) error {
	return d.Send(serverID, msg)
}

// ReceiveFromServer 接收来自服务端的消息
func (d *DealerNode) ReceiveFromServer() (serverID string, msg string, err error) {
	msgs, err := d.Receive()
	if err != nil {
		return "", "", err
	}
	return msgs[0], msgs[1], nil
}
