package model

import "time"

// EventType represents the type of SSE event.
type EventType string

const (
	EventTypeIssueCreated EventType = "issue_created"
	EventTypeIssueUpdated EventType = "issue_updated"
	EventTypeIssueDeleted EventType = "issue_deleted"
	EventTypeHeartbeat    EventType = "heartbeat"
)

// Event is the base type for all SSE events.
type Event struct {
	Type      EventType   `json:"event"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// IssueCreatedEvent is sent when a new issue is created.
type IssueCreatedEvent struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// IssueUpdatedEvent is sent when an issue is modified.
type IssueUpdatedEvent struct {
	ID             string    `json:"id"`
	Status         Status    `json:"status"`
	PreviousStatus Status    `json:"previous_status,omitempty"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// IssueDeletedEvent is sent when an issue is removed.
type IssueDeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}

// HeartbeatEvent is sent periodically to keep the connection alive.
type HeartbeatEvent struct {
	Timestamp time.Time `json:"timestamp"`
}

// NewHeartbeat creates a new heartbeat event with current timestamp.
func NewHeartbeat() Event {
	now := time.Now()
	return Event{
		Type:      EventTypeHeartbeat,
		Data:      HeartbeatEvent{Timestamp: now},
		Timestamp: now,
	}
}

// NewIssueCreatedEvent creates an issue_created event.
func NewIssueCreatedEvent(id, title string, status Status) Event {
	now := time.Now()
	return Event{
		Type: EventTypeIssueCreated,
		Data: IssueCreatedEvent{
			ID:        id,
			Title:     title,
			Status:    status,
			CreatedAt: now,
		},
		Timestamp: now,
	}
}

// NewIssueUpdatedEvent creates an issue_updated event.
func NewIssueUpdatedEvent(id string, status, previousStatus Status) Event {
	now := time.Now()
	return Event{
		Type: EventTypeIssueUpdated,
		Data: IssueUpdatedEvent{
			ID:             id,
			Status:         status,
			PreviousStatus: previousStatus,
			UpdatedAt:      now,
		},
		Timestamp: now,
	}
}
