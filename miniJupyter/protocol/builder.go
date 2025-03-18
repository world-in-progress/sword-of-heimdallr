package protocol

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// MessageBuilder creates messages using the builder pattern
type MessageBuilder struct {
    message *Message
}

func NewMessageBuilder() *MessageBuilder {
    return &MessageBuilder{
        message: &Message{
            Header: Header{
                MsgId:       GenerateUUID(),
                SessionId:   "",  // Need to be explicitly set
                UserId:      "",  // Need to be explicitly set
                Timestamp:   time.Now(),
                MsgType:     "",  // Need to be explicitly set
                Compression: CompressNone,
                Encoding:    EncodeJSON,
                Transport:   "",  // Need to be explicitly set
                Version:     ProtocolVersion,
            },
            ParentHeader: Header{},  // Keep empty initialization, can be set with WithParentMessage/WithParentHeader
            Meta:        Metadata{}, // Keep empty initialization, can be set with WithPriority/WithTags etc.
            Trace:       NewMessageTrace(), // Keep initialization, ensure trace functionality is available
            Security:    SecurityConfig{},  // Keep empty initialization, can be set with WithSecurity etc.
        },
    }
}

// Required setting methods
func (b *MessageBuilder) WithType(msgType string) *MessageBuilder {
    b.message.Header.MsgType = msgType
    return b
}

func (b *MessageBuilder) WithSession(sessionId string) *MessageBuilder {
    b.message.Header.SessionId = sessionId
    return b
}

func (b *MessageBuilder) WithUser(userId string) *MessageBuilder {
    b.message.Header.UserId = userId
    return b
}

func (b *MessageBuilder) WithTransport(transport Transport) *MessageBuilder {
    b.message.Header.Transport = transport
    return b
}

// Optional setting methods
func (b *MessageBuilder) WithCompression(compression Compression) *MessageBuilder {
    b.message.Header.Compression = compression
    return b
}

func (b *MessageBuilder) WithEncoding(encoding Encoding) *MessageBuilder {
    b.message.Header.Encoding = encoding
    return b
}

func (b *MessageBuilder) WithContent(content interface{}) *MessageBuilder {
    b.message.Content = content
    return b
}

// Meta related methods
func (b *MessageBuilder) WithPriority(priority Priority) *MessageBuilder {
    b.message.Meta.Priority = priority
    return b
}

func (b *MessageBuilder) WithTags(tags []string) *MessageBuilder {
    b.message.Meta.Tags = tags
    return b
}

func (b *MessageBuilder) AddTag(tag string) *MessageBuilder {
    b.message.Meta.Tags = append(b.message.Meta.Tags, tag)
    return b
}

// Security related methods
func (b *MessageBuilder) WithToken(token string) *MessageBuilder {
    b.message.Security.Token = token
    return b
}

func (b *MessageBuilder) WithEncryption(encryption string) *MessageBuilder {
    b.message.Security.Encryption = encryption
    return b
}

// Convenient method: set both token and encryption
func (b *MessageBuilder) WithSecurity(token, encryption string) *MessageBuilder {
    b.message.Security.Token = token
    b.message.Security.Encryption = encryption
    return b
}

// ParentHeader related methods
func (b *MessageBuilder) WithParentMessage(parent *Message) *MessageBuilder {
    if parent != nil {
        b.message.ParentHeader = parent.Header
    }
    return b
}

// Also can directly set ParentHeader
func (b *MessageBuilder) WithParentHeader(header Header) *MessageBuilder {
    b.message.ParentHeader = header
    return b
}

// Trace related methods
func (b *MessageBuilder) WithNewTrace() *MessageBuilder {
    b.message.Trace = NewMessageTrace()
    return b
}

func (b *MessageBuilder) WithTrace(trace *MessageTrace) *MessageBuilder {
    b.message.Trace = trace
    return b
}

// Convenient method: add trace and record the first hop
func (b *MessageBuilder) WithTraceHop(serviceId, serviceName, hostName string) *MessageBuilder {
    if b.message.Trace == nil {
        b.message.Trace = NewMessageTrace()
    }
    b.message.Trace.AddHop(serviceId, serviceName, hostName)
    return b
}

// GenerateUUID generates a UUID
func GenerateUUID() string {
    bytes := make([]byte, 16)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

// Build method includes necessary validation
func (b *MessageBuilder) Build() (*Message, error) {
    // Validate required fields
    if b.message.Header.MsgType == "" {
        return nil, ErrInvalidMessageType
    }
    if b.message.Header.SessionId == "" {
        return nil, NewProtocolError(ErrCodeInvalidMessage, "session_id is required", nil)
    }
    if b.message.Header.UserId == "" {
        return nil, NewProtocolError(ErrCodeInvalidMessage, "user_id is required", nil)
    }
    if b.message.Header.Transport == "" {
        return nil, NewProtocolError(ErrCodeInvalidMessage, "transport is required", nil)
    }

    return b.message, nil
}