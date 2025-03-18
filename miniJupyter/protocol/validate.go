package protocol

import (
	"errors"
	"fmt"
)

// Validator interface defines message validation methods
type Validator interface {
    Validate() error
}

// ValidateMessage validates the entire message structure
func ValidateMessage(msg *Message) error {
    // Validate Header
    if err := validateHeader(&msg.Header); err != nil {
        return fmt.Errorf("invalid header: %w", err)
    }

    // Validate Content
    if validator, ok := msg.Content.(Validator); ok {
        if err := validator.Validate(); err != nil {
            return fmt.Errorf("invalid content: %w", err)
        }
    }

    return nil
}

// validateHeader validates the message header
func validateHeader(h *Header) error {
    if h.MsgId == "" {
        return errors.New("msg_id is required")
    }
    if h.SessionId == "" {
        return errors.New("session_id is required")
    }
    if h.UserId == "" {
        return errors.New("user_id is required")
    }
    if !IsValidMessageType(h.MsgType) {
        return fmt.Errorf("invalid message type: %s", h.MsgType)
    }
    if h.Version != ProtocolVersion {
        return fmt.Errorf("unsupported version: %s", h.Version)
    }
    return nil
}

// ExecuteRequestContent validation
func (c *ExecuteRequestContent) Validate() error {
    if c.CommandId == "" {
        return errors.New("command_id is required")
    }
    if c.Service == "" {
        return errors.New("service is required")
    }
    if c.Method == "" {
        return errors.New("method is required")
    }
    if c.Timeout < 0 {
        return errors.New("timeout cannot be negative")
    }
    if c.Retry.MaxAttempts < 0 {
        return errors.New("retry max_attempts cannot be negative")
    }
    return nil
}

// ExecuteReplyContent validation
func (c *ExecuteReplyContent) Validate() error {
    switch c.Status {
    case StatusError, StatusStarting, StatusWaiting:
        return nil
    default:
        return fmt.Errorf("invalid status: %s", c.Status)
    }
}

// CoreInfoContent validation
func (c *CoreInfoContent) Validate() error {
    if c.CoreVersion == "" {
        return errors.New("core_version is required")
    }
    if c.ActiveConnections < 0 {
        return errors.New("active_connections cannot be negative")
    }
    if c.RunningTasks < 0 {
        return errors.New("running_tasks cannot be negative")
    }
    if c.TaskQueueSize < 0 {
        return errors.New("task_queue_size cannot be negative")
    }
    return nil
}

// ExecuteResultContent validation
func (c *ExecuteResultContent) Validate() error {
    switch c.Status {
    case StatusSuccess, StatusError:
        return nil
    default:
        return fmt.Errorf("invalid status: %s", c.Status)
    }
}

// StreamContent validation
func (c *StreamContent) Validate() error {
    switch c.Type {
    case StreamStdout, StreamStderr:
        return nil
    default:
        return fmt.Errorf("invalid stream type: %s", c.Type)
    }
}

// CommOpenContent validation
func (c *CommOpenContent) Validate() error {
    if c.CommId == "" {
        return errors.New("comm_id is required")
    }
    if c.TargetName == "" {
        return errors.New("target_name is required")
    }
    return nil
}

// CommMsgContent validation
func (c *CommMsgContent) Validate() error {
    if c.CommId == "" {
        return errors.New("comm_id is required")
    }
    return nil
}