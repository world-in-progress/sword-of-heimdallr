package protocol

import "fmt"

// ProtocolError defines the protocol error type
type ProtocolError struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *ProtocolError) Error() string {
    if e.Details != nil {
        return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Details)
    }
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NewProtocolError creates a new protocol error
func NewProtocolError(code int, message string, details interface{}) *ProtocolError {
    return &ProtocolError{
        Code:    code,
        Message: message,
        Details: details,
    }
}

// Predefined error codes
const (
    // 1000-1099: Protocol level errors
    ErrCodeInvalidMessage     = 1000  // Message format does not comply with protocol specification
    ErrCodeInvalidMessageType = 1001  // Message type not in predefined type list
    ErrCodeInvalidVersion     = 1002  // Protocol version mismatch or unsupported
    ErrCodeInvalidFormat      = 1003  // Incorrect message structure (e.g., missing required fields)
    ErrCodeValidationFailed   = 1004  // Message content validation failed (e.g., invalid field values)
    ErrCodeSerializeFailed    = 1005  // Message serialization failed (e.g., JSON conversion)
    ErrCodeDeserializeFailed  = 1006  // Message deserialization failed (e.g., JSON parsing)

    // 1100-1199: Authentication/Authorization errors
    ErrCodeUnauthorized      = 1100  // Unauthorized access
    ErrCodeInvalidToken      = 1101  // Invalid authentication token
    ErrCodeInsufficientPerms = 1102  // Insufficient permissions
    ErrCodeSessionExpired    = 1103  // Session expired

    // 1200-1299: Execution errors
    ErrCodeExecutionFailed   = 1200  // Execution failed
    ErrCodeTimeout           = 1201  // Operation timeout
    ErrCodeDependencyFailed  = 1202  // Dependency execution failed
    ErrCodeServiceNotFound   = 1203  // Service not found
    ErrCodeMethodNotFound    = 1204  // Method not found
    ErrCodeInvalidParams     = 1205  // Invalid parameters

    // 1300-1399: Communication errors
    ErrCodeConnectionFailed  = 1300  // Connection failed
    ErrCodeHeartbeatTimeout  = 1301  // Heartbeat timeout
    ErrCodeSubscribeFailed   = 1302  // Subscribe failed
    ErrCodePublishFailed     = 1303  // Publish failed
    ErrCodeCommFailed        = 1304  // Communication operation failed
)

// Predefined error instances
var (
    // Protocol errors
    ErrInvalidMessage     = NewProtocolError(ErrCodeInvalidMessage, "Invalid message format", nil)
    ErrInvalidMessageType = NewProtocolError(ErrCodeInvalidMessageType, "Invalid message type", nil)
    ErrInvalidVersion     = NewProtocolError(ErrCodeInvalidVersion, "Invalid protocol version", nil)
    ErrInvalidFormat      = NewProtocolError(ErrCodeInvalidFormat, "Invalid message format", nil)
    ErrValidationFailed   = NewProtocolError(ErrCodeValidationFailed, "Message validation failed", nil)
    ErrSerializeFailed    = NewProtocolError(ErrCodeSerializeFailed, "Message serialization failed", nil)
    ErrDeserializeFailed  = NewProtocolError(ErrCodeDeserializeFailed, "Message deserialization failed", nil)

    // Auth errors
    ErrUnauthorized       = NewProtocolError(ErrCodeUnauthorized, "Unauthorized access", nil)
    ErrInvalidToken       = NewProtocolError(ErrCodeInvalidToken, "Invalid token", nil)
    ErrInsufficientPerms  = NewProtocolError(ErrCodeInsufficientPerms, "Insufficient permissions", nil)
    ErrSessionExpired     = NewProtocolError(ErrCodeSessionExpired, "Session expired", nil)

    // Execution errors
    ErrExecutionFailed    = NewProtocolError(ErrCodeExecutionFailed, "Execution failed", nil)
    ErrTimeout           = NewProtocolError(ErrCodeTimeout, "Operation timeout", nil)
    ErrDependencyFailed  = NewProtocolError(ErrCodeDependencyFailed, "Dependency execution failed", nil)
    ErrServiceNotFound   = NewProtocolError(ErrCodeServiceNotFound, "Service not found", nil)
    ErrMethodNotFound    = NewProtocolError(ErrCodeMethodNotFound, "Method not found", nil)
    ErrInvalidParams     = NewProtocolError(ErrCodeInvalidParams, "Invalid parameters", nil)

    // Communication errors
    ErrConnectionFailed  = NewProtocolError(ErrCodeConnectionFailed, "Connection failed", nil)
    ErrHeartbeatTimeout  = NewProtocolError(ErrCodeHeartbeatTimeout, "Heartbeat timeout", nil)
    ErrSubscribeFailed   = NewProtocolError(ErrCodeSubscribeFailed, "Subscribe failed", nil)
    ErrPublishFailed     = NewProtocolError(ErrCodePublishFailed, "Publish failed", nil)
    ErrCommFailed        = NewProtocolError(ErrCodeCommFailed, "Comm operation failed", nil)
)

// WithDetails adds error details
func (e *ProtocolError) WithDetails(details interface{}) *ProtocolError {
    return &ProtocolError{
        Code:    e.Code,
        Message: e.Message,
        Details: details,
    }
}

// IsProtocolError checks if the error is a protocol error
func IsProtocolError(err error) bool {
    _, ok := err.(*ProtocolError)
    return ok
}

// GetErrorCode gets the error code, returns 0 if it's not a protocol error
func GetErrorCode(err error) int {
    if pe, ok := err.(*ProtocolError); ok {
        return pe.Code
    }
    return 0
}