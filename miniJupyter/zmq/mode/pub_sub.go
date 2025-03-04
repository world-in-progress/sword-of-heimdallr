package mode

import (
	"miniJupyter/zmq/base"

	zmq "github.com/pebbe/zmq4"
)

type PublisherNode struct {
	*base.ZmqNode
}

type SubscriberNode struct {
	*base.ZmqNode
	topics []string
}

func NewPublisher(address string) (*PublisherNode, error) {
	node, err := base.NewZmqNode(zmq.PUB, address, true)
	if err != nil {
		return nil, err
	}
	return &PublisherNode{node}, nil
}

func NewSubscriber(address string) (*SubscriberNode, error) {
	node, err := base.NewZmqNode(zmq.SUB, address, false)
	if err != nil {
		return nil, err
	}
	return &SubscriberNode{node, make([]string, 0)}, nil
}

// PublishWithTopic 发布带主题的消息
func (p *PublisherNode) PublishWithTopic(topic, message string) error {
	return p.Send(topic, message)
}

// Subscribe 订阅特定主题
func (s *SubscriberNode) Subscribe(topic string) error {
	err := s.SetSubscribe(topic)
	if err != nil {
		return err
	}
	s.topics = append(s.topics, topic)
	return nil
}

// GetTopics 获取当前订阅的所有主题
func (s *SubscriberNode) GetTopics() []string {
	return s.topics
}