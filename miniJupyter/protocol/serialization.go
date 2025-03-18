package protocol

import (
	"encoding/json"
	"fmt"
)

// GetContentType returns the corresponding Content struct based on message type
func GetContentType(msgType string) interface{} {
	switch msgType {
	// ROUTER/DEALER messages
	case MsgTypeExecuteRequest:
		return &ExecuteRequestContent{}
	case MsgTypeExecuteReply:
		return &ExecuteReplyContent{}
	case MsgTypeCoreInfoRequest:
		return &struct{}{} // Empty struct, because this request has no content
	case MsgTypeCoreInfoReply:
		return &CoreInfoContent{}

	// PUB/SUB messages
	case MsgTypeExecuteResult:
		return &ExecuteResultContent{}
	case MsgTypeStream:
		return &StreamContent{}

	// Comm messages
	case MsgTypeCommOpen:
		return &CommOpenContent{}
	case MsgTypeCommMsg:
		return &CommMsgContent{}
	case MsgTypeCommClose:
		return &CommMsgContent{} // CommClose uses the same structure as CommMsg

	default:
		return nil
	}
}

// ParseMessage smartly parses messages
func ParseMessage(data string) (*Message, error) {
	// 1. Parse the basic message structure first
	var msg Message
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		return nil, fmt.Errorf("failed to parse message: %w", err)
	}
	
	// 2. Get the correct Content type
	contentType := GetContentType(msg.Header.MsgType)
	if contentType == nil {
		return nil, fmt.Errorf("unknown message type: %s", msg.Header.MsgType)
	}
	
	// 3. Re-parse Content to the correct type
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

// SerializeMessage serializes the message to a JSON string
func SerializeMessage(msg *Message) (string, error) {
	if msg == nil {
		return "", fmt.Errorf("cannot serialize nil message")
	}

	// Validate message
	if err := ValidateMessage(msg); err != nil {
		return "", fmt.Errorf("message validation failed: %w", err)
	}

	// Serialize to JSON
	data, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("failed to serialize message: %w", err)
	}

	return string(data), nil
}

// DeserializeMessage deserializes a JSON string to a message
func DeserializeMessage(data string) (*Message, error) {
	return ParseMessage(data)
}