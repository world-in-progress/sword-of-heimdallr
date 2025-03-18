package protocol

// Protocol version
const (
    ProtocolVersion = "0.4"
)

// Enum type definitions
type (
    Compression string
    Encoding    string
    Transport   string
    Priority    string
    Status      string
    StreamType  string
    RetryStrategy string
)

const (
    // Compression
    CompressNone   Compression = "none"
    CompressGzip   Compression = "gzip"
    CompressSnappy Compression = "snappy"

    // Encoding
    EncodeJSON     Encoding = "json"
    EncodeProtobuf Encoding = "protobuf"
    EncodeCustom   Encoding = "custom"

    // Transport
    TransportZMQ   Transport = "zmq"
    TransportGRPC  Transport = "grpc"

    // Priority
    PriorityHigh   Priority = "HIGH"
    PriorityNormal Priority = "NORMAL"
    PriorityLow    Priority = "LOW"

    // Status
    StatusOK       Status = "ok"
    StatusError    Status = "error"
    StatusStarting Status = "starting"
    StatusWaiting  Status = "waiting"
    StatusSuccess  Status = "success"

    // StreamType
    StreamStdout StreamType = "stdout"
    StreamStderr StreamType = "stderr"

    // RetryStrategy
    RetryExponentialBackoff RetryStrategy = "exponential_backoff"

    // Define message type constants
    MsgTypeExecuteRequest  = "execute_request"
    MsgTypeExecuteReply    = "execute_reply"
    MsgTypeExecuteResult   = "execute_result"
    MsgTypeCoreInfoRequest = "core_info_request"
    MsgTypeCoreInfoReply   = "core_info_reply"
    MsgTypeStream         = "stream"
    MsgTypeCommOpen       = "comm_open"
    MsgTypeCommMsg        = "comm_msg"
    MsgTypeCommClose      = "comm_close"
)

// Add message type check
func IsValidMessageType(msgType string) bool {
    switch msgType {
    case MsgTypeExecuteRequest, MsgTypeExecuteReply, MsgTypeExecuteResult,
         MsgTypeCoreInfoRequest, MsgTypeCoreInfoReply, MsgTypeStream,
         MsgTypeCommOpen, MsgTypeCommMsg, MsgTypeCommClose:
        return true
    }
    return false
}