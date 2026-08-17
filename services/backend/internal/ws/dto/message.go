package dto

type RoleType string
type EventType string

const (
	UserRole      RoleType = "user"
	AssistantRole RoleType = "assistant"
	ToolRole      RoleType = "tool"
)

const (
	TextDeltaEvent EventType = "text-delta"
	DoneEvent      EventType = "done"
	ToolEvent      EventType = "tool-call"
)

type Message struct {
	Role    RoleType   `json:"role"`
	Content string     `json:"content"`
	Event   *EventType `json:"event_type,omitempty"`
}
