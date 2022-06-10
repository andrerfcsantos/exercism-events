package events

import "github.com/andrerfcsantos/exercism-events/exgo"

type NotificationEventType int

const (
	NotificationAdded NotificationEventType = iota
	NotificationDeleted
)

type NotificationEvent struct {
	Type         NotificationEventType
	Notification exgo.Notification
}
