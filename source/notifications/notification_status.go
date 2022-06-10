package notifications

import (
	"github.com/andrerfcsantos/exercism-events/exgo"
)

type NotificationStatus struct {
	notifications map[string]exgo.Notification
}

func NewNotificationStatus(initialState ...exgo.Notification) *NotificationStatus {
	status := &NotificationStatus{notifications: make(map[string]exgo.Notification)}

	for _, req := range initialState {
		status.notifications[req.UUID] = req
	}

	return status
}

type NotificationDiff struct {
	Added   []exgo.Notification
	Removed []exgo.Notification
}

func (status *NotificationStatus) Update(newState ...exgo.Notification) NotificationDiff {
	newStatus := NewNotificationStatus(newState...)

	var diff NotificationDiff

	for _, notification := range newStatus.notifications {
		if _, ok := status.notifications[notification.UUID]; !ok {
			diff.Added = append(diff.Added, notification)
		}
	}

	for _, notification := range status.notifications {
		if _, ok := newStatus.notifications[notification.UUID]; !ok {
			diff.Removed = append(diff.Removed, notification)
		}
	}

	status.notifications = newStatus.notifications

	return diff
}
