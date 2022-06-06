package mentoring

import "github.com/andrerfcsantos/exercism-events/exgo"

type MentoringStatus struct {
	requests map[string]exgo.MentoringRequest
}

func NewMentoringStatus(initialState ...exgo.MentoringRequest) *MentoringStatus {
	status := &MentoringStatus{requests: make(map[string]exgo.MentoringRequest)}

	for _, req := range initialState {
		status.requests[req.UUID] = req
	}

	return status
}

type StatusDiff struct {
	Added   []exgo.MentoringRequest
	Removed []exgo.MentoringRequest
}

func (status *MentoringStatus) Update(newState ...exgo.MentoringRequest) StatusDiff {
	newStatus := NewMentoringStatus(newState...)

	var diff StatusDiff

	for _, req := range newStatus.requests {
		if _, ok := status.requests[req.UUID]; !ok {
			diff.Added = append(diff.Added, req)
		}
	}

	for _, req := range status.requests {
		if _, ok := newStatus.requests[req.UUID]; !ok {
			diff.Removed = append(diff.Removed, req)
		}
	}

	status.requests = newStatus.requests

	return diff
}
