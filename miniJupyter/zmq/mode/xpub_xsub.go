package mode

import (
	"fmt"
	"miniJupyter/zmq/base"
	"strings"

	zmq "github.com/pebbe/zmq4"
)

// TopicPermission 定义主题的权限控制
type TopicPermission struct {
	Topic        string
	AllowedUsers []string
}

// XPublisherNode XPUB节点
type XPublisherNode struct {
	*base.ZmqNode
	permissions map[string][]string // topic -> allowed users
}

// XSubscriberNode XSUB节点
type XSubscriberNode struct {
	*base.ZmqNode
	topics []string
	userID string // 用户ID字段
}

// NewXPublisher 创建新的 XPUB 节点，增加权限控制
func NewXPublisher(address string) (*XPublisherNode, error) {
	node, err := base.NewZmqNode(zmq.XPUB, address, true, "")
	return &XPublisherNode{
		ZmqNode:     node,
		permissions: make(map[string][]string),
	}, err
}

// NewXSubscriber 创建新的 XSUB 节点，需要提供用户ID
func NewXSubscriber(address string, userID string) (*XSubscriberNode, error) {
	node, err := base.NewZmqNode(zmq.XSUB, address, false, "")
	return &XSubscriberNode{
		ZmqNode: node,
		topics:  make([]string, 0),
		userID:  userID,
	}, err
}

// SetTopicPermission 设置主题的访问权限
func (xp *XPublisherNode) SetTopicPermission(topic string, allowedUsers []string) {
	xp.permissions[topic] = allowedUsers
}

// RemoveTopicPermission 移除主题的访问权限
func (xp *XPublisherNode) RemoveTopicPermission(topic string) {
	delete(xp.permissions, topic)
}

// HasPermission 检查用户是否有权限访问主题
func (xp *XPublisherNode) HasPermission(topic, userID string) bool {
	allowedUsers, exists := xp.permissions[topic]
	if !exists {
		return true // 如果没有设置权限，默认允许所有用户访问
	}
	for _, user := range allowedUsers {
		if user == userID {
			return true
		}
	}
	return false
}

// Publish 发布消息
func (xp *XPublisherNode) Publish(topic, message string) error {
	return xp.Send(topic, message)
}

// HandleSubscription 处理订阅请求
func (xp *XPublisherNode) HandleSubscription(data []byte) error {
	if len(data) < 2 { // 至少包含1字节的订阅标志和一些数据
		return fmt.Errorf("invalid subscription data")
	}

	isSubscribe := data[0] == 1 // 1表示订阅，0表示取消订阅
	payload := string(data[1:])

	// 解析用户ID和主题
	parts := strings.Split(payload, "|")
	if len(parts) != 2 {
		return fmt.Errorf("invalid subscription format")
	}
	userID, topic := parts[0], parts[1]

	// 检查权限
	if isSubscribe && !xp.HasPermission(topic, userID) {
		return fmt.Errorf("user %s does not have permission to subscribe to topic %s", userID, topic)
	}

	// 转发订阅消息
	return xp.Send(string(data))
}

// Subscribe XSUB 订阅特定主题，包含权限验证信息
func (xs *XSubscriberNode) Subscribe(topic string) error {
	// 发送订阅消息，包含用户ID信息
	subscribeMsg := fmt.Sprintf("%s|%s", xs.userID, topic)
	err := xs.Send(string([]byte{1}) + subscribeMsg)
	if err != nil {
		return err
	}
	xs.topics = append(xs.topics, topic)
	return nil
}

// Unsubscribe XSUB 取消订阅特定主题
func (xs *XSubscriberNode) Unsubscribe(topic string) error {
	// 发送取消订阅消息，包含用户ID信息
	unsubscribeMsg := fmt.Sprintf("%s|%s", xs.userID, topic)
	err := xs.Send(string([]byte{0}) + unsubscribeMsg)
	if err != nil {
		return err
	}
	// 从topics列表中移除该主题
	for i, t := range xs.topics {
		if t == topic {
			xs.topics = append(xs.topics[:i], xs.topics[i+1:]...)
			break
		}
	}
	return nil
}

// GetTopics 获取当前订阅的所有主题
func (xs *XSubscriberNode) GetTopics() []string {
	return xs.topics
}

// Run XPUB节点的主运行循环
func (xp *XPublisherNode) Run() error {
	for {
		// 接收订阅消息
		data, err := xp.Receive()
		if err != nil {
			return err
		}

		// 处理订阅请求
		if err := xp.HandleSubscription([]byte(data[0])); err != nil {
			// 记录错误但继续运行
			fmt.Printf("Error handling subscription: %v\n", err)
		}
	}
}