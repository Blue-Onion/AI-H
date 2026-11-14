package types

type StopReasons string
type Roles string
type StreamType string

const (
	Stop    StopReasons = "stop"
	ToolUse StopReasons = "toolUse"
	Length  StopReasons = "length"
)

const (
	Assist Roles = "Assistant"
	User   Roles = "User"
)

const (
	Delta StreamType = "text_delta"
	Done  StreamType = "done"
)

type Usage struct {
	Output int64
	Input  int64
}

// Message can be either *UserMessage or *AssistantMessage.
type Message interface {
	isMessage()
}

type UserMessage struct {
	Role    Roles
	content string
}

type AssistantMessage struct {
	Role       Roles
	content    string
	Usage      Usage
	StopReason StopReasons
}

func (u *UserMessage) isMessage()      {}
func (u *AssistantMessage) isMessage() {}

// StreamEvent can be either *StreamEventDelta or *StreamEventDone.
type StreamEvent interface {
	isStreamEvent()
}

type StreamEventDelta struct {
	Type  StreamType
	Delta string
}

type StreamEventDone struct {
	Type    StreamType
	Message AssistantMessage
}

func (s *StreamEventDelta) isStreamEvent() {}
func (s *StreamEventDone) isStreamEvent()  {}

type StreamOpts struct {
	Message []Message
	Model   string
	System  string
}

type Provider struct {
	Name         string
	DefaultModel string
	Stream       func(opts StreamOpts) <-chan StreamEvent
}
