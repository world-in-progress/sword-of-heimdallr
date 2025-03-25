package mode

import (
	"fmt"
	"miniJupyter/zmq/base"

	zmq "github.com/pebbe/zmq4"
)

// RouterNode extends the basic ZmqNode, adding Router-specific functionality
type RouterNode struct {
	*base.ZmqNode
}

// DealerNode extends the basic ZmqNode, adding Dealer-specific functionality
type DealerNode struct {
	*base.ZmqNode
}

// NewRouter creates and returns a RouterNode (no identity required)
func NewRouter(address string, bind bool) (*RouterNode, error) {
	node, err := base.NewZmqNode(zmq.ROUTER, address, bind, "")  // Empty identity
	if err != nil {
		return nil, err
	}
	return &RouterNode{node}, nil
}

// NewDealer creates and returns a DealerNode (requires identity)
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

// SendToClient sends a message to a specific client
func (r *RouterNode) SendToClient(clientID string, msg string) error {
	return r.Send(clientID, msg)
}

// ReceiveFromClient receives a message from a client, returning the client ID and message
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

// SendToServer sends a message to a specific server
func (d *DealerNode) SendToServer(msg string) error {
	return d.Send(msg)
}

// ReceiveFromServer receives a message from a server
func (d *DealerNode) ReceiveFromServer() (msg string, err error) {
	msgs, err := d.Receive()
	if err != nil {
		return "", err
	}
	return msgs[0], nil
}
