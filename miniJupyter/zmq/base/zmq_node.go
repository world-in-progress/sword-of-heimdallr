package base

import (
	"fmt"
	"os"

	zmq "github.com/pebbe/zmq4"
	"gopkg.in/yaml.v2"
)

// Read configuration file
type Config struct {
    Zmq struct {
        RouterAddress     string `yaml:"RouterAddress"`
        DealerAddress     string `yaml:"DealerAddress"`
        PubAddress        string `yaml:"PubAddress"`
        SubAddress        string `yaml:"SubAddress"`
        HeartbeatInterval int    `yaml:"HeartbeatInterval"`
        HeartbeatTimeout  int    `yaml:"HeartbeatTimeout"`
    } `yaml:"zmq"`
}

// Parse YAML configuration
func LoadConfig(filename string) (*Config, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }
    return &config, nil
}

// ZmqNode structure, encapsulates ZMQ logic
type ZmqNode struct {
    socket *zmq.Socket
}

// NewZmqNode creates a ZMQ endpoint
func NewZmqNode(socketType zmq.Type, address string, bind bool, identity string) (*ZmqNode, error) {
    socket, err := zmq.NewSocket(socketType)
    if err != nil {
        return nil, err
    }

    // If an identity is provided, set it
    if identity != "" {
        if err := socket.SetIdentity(identity); err != nil {
            socket.Close()
            return nil, fmt.Errorf("failed to set identity: %w", err)
        }
    }

    if bind {
        err = socket.Bind(address)
    } else {
        err = socket.Connect(address)
    }
    
    if err != nil {
        socket.Close()
        return nil, err
    }

    return &ZmqNode{socket: socket}, nil
}

// Send message, supports strings and byte arrays
func (z *ZmqNode) Send(parts ...interface{}) error {
    _, err := z.socket.SendMessage(parts...)
    return err
}

// Receive message
func (z *ZmqNode) Receive() ([]string, error) {
    return z.socket.RecvMessage(0)
}

// Close ZMQ connection
func (z *ZmqNode) Close() {
    z.socket.Close()
}

// SetSubscribe sets the subscription topic
func (z *ZmqNode) SetSubscribe(topic string) error {
    return z.socket.SetSubscribe(topic)
}
