package protocol

import (
	"encoding/json"
	"time"
)

// MessageTrace defines message tracing structure
type MessageTrace struct {
    TraceId    string       `json:"trace_id"`     // Trace ID
    StartTime  time.Time    `json:"start_time"`   // Message creation time
    Hops       []MessageHop `json:"hops"`         // Service nodes the message passed through
    TotalTime  Duration     `json:"total_time"`   // Total processing time
}

// MessageHop defines information for each service node the message passes through
type MessageHop struct {
    ServiceId   string   `json:"service_id"`    // Service ID
    ServiceName string   `json:"service_name"`  // Service name
    HostName    string   `json:"host_name"`     // Host name
    EntryTime   time.Time `json:"entry_time"`   // Time entered service
    ExitTime    time.Time `json:"exit_time"`    // Time left service
    Duration    Duration  `json:"duration"`      // Processing duration
    Status      string    `json:"status"`       // Processing status
    Error       string    `json:"error,omitempty"` // Error message (if any)
}

// Duration custom time type, supports more friendly JSON serialization
type Duration time.Duration

// MarshalJSON implements JSON serialization for Duration
func (d Duration) MarshalJSON() ([]byte, error) {
    return json.Marshal(time.Duration(d).String())
}

// UnmarshalJSON implements JSON deserialization for Duration
func (d *Duration) UnmarshalJSON(b []byte) error {
    var v string
    if err := json.Unmarshal(b, &v); err != nil {
        return err
    }
    parsed, err := time.ParseDuration(v)
    if err != nil {
        return err
    }
    *d = Duration(parsed)
    return nil
}

// NewMessageTrace creates a new message trace
func NewMessageTrace() *MessageTrace {
    return &MessageTrace{
        TraceId:   GenerateUUID(),
        StartTime: time.Now(),
        Hops:      make([]MessageHop, 0),
    }
}

// AddHop adds a service node
func (mt *MessageTrace) AddHop(serviceId, serviceName, hostName string) *MessageHop {
    hop := MessageHop{
        ServiceId:   serviceId,
        ServiceName: serviceName,
        HostName:    hostName,
        EntryTime:   time.Now(),
    }
    mt.Hops = append(mt.Hops, hop)
    return &mt.Hops[len(mt.Hops)-1]
}

// CompleteHop completes the processing of the current service node
func (h *MessageHop) Complete(status string, err error) {
    h.ExitTime = time.Now()
    h.Duration = Duration(h.ExitTime.Sub(h.EntryTime))
    h.Status = status
    if err != nil {
        h.Error = err.Error()
    }
}

// CalculateTotalTime calculates the total processing time of the message
func (mt *MessageTrace) CalculateTotalTime() {
    if len(mt.Hops) == 0 {
        mt.TotalTime = 0
        return
    }
    
    firstHop := mt.Hops[0]
    lastHop := mt.Hops[len(mt.Hops)-1]
    mt.TotalTime = Duration(lastHop.ExitTime.Sub(firstHop.EntryTime))
}

// GetHopByService gets the processing information for the specified service
func (mt *MessageTrace) GetHopByService(serviceName string) *MessageHop {
    for i := range mt.Hops {
        if mt.Hops[i].ServiceName == serviceName {
            return &mt.Hops[i]
        }
    }
    return nil
}

// String implements the string representation of the message trace
func (mt *MessageTrace) String() string {
    data, _ := json.MarshalIndent(mt, "", "  ")
    return string(data)
}