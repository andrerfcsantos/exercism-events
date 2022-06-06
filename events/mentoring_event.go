package events

import "github.com/andrerfcsantos/exercism-events/exgo"

type EventType int

const (
	NewMentoringRequest EventType = iota
	MentoringRequestDeleted
)

type MentoringEvent struct {
	Track   string
	Type    EventType
	Request exgo.MentoringRequest
}
