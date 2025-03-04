package protocol

import (
	"encoding/json"
	"fmt"
)

// GetContentType 根据消息类型返回对应的 Content 结构体
func GetContentType(msgType string) interface{} {
	switch msgType {
	// ROUTER/DEALER 消息
	case MsgTypeExecuteRequest:
		return &ExecuteRequestContent{}
	case MsgTypeExecuteReply:
		return &ExecuteReplyContent{}
	case MsgTypeCoreInfoRequest:
		return &struct{}{} // 空结构体，因为该请求没有content
	case MsgTypeCoreInfoReply:
		return &CoreInfoContent{}

	// PUB/SUB 消息
	case MsgTypeExecuteResult:
		return &ExecuteResultContent{}
	case MsgTypeStream:
		return &StreamContent{}

	// Comm 消息
	case MsgTypeCommOpen:
		return &CommOpenContent{}
	case MsgTypeCommMsg:
		return &CommMsgContent{}
	case MsgTypeCommClose:
		return &CommMsgContent{} // CommClose 使用相同的结构

	default:
		return nil
	}
}

// ParseMessage 智能解析消息
func ParseMessage(data string) (*Message, error) {
	// 1. 先解析基础消息结构
	var msg Message
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		return nil, fmt.Errorf("failed to parse message: %w", err)
	}
	
	// 2. 获取正确的 Content 类型
	contentType := GetContentType(msg.Header.MsgType)
	if contentType == nil {
		return nil, fmt.Errorf("unknown message type: %s", msg.Header.MsgType)
	}
	
	// 3. 重新解析 Content 到正确的类型
	contentData, err := json.Marshal(msg.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal content: %w", err)
	}
	
	if err := json.Unmarshal(contentData, contentType); err != nil {
		return nil, fmt.Errorf("failed to parse content: %w", err)
	}
	
	msg.Content = contentType
	return &msg, nil
}

// SerializeMessage 序列化消息为JSON字符串
func SerializeMessage(msg *Message) (string, error) {
	if msg == nil {
		return "", fmt.Errorf("cannot serialize nil message")
	}

	// 验证消息
	if err := ValidateMessage(msg); err != nil {
		return "", fmt.Errorf("message validation failed: %w", err)
	}

	// 序列化为JSON
	data, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("failed to serialize message: %w", err)
	}

	return string(data), nil
}

// DeserializeMessage 反序列化JSON字符串为消息
func DeserializeMessage(data string) (*Message, error) {
	return ParseMessage(data)
}