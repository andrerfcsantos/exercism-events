package notifications

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/andrerfcsantos/exercism-events/events"
	"github.com/andrerfcsantos/exercism-events/exgo"
)

type NotificationEventSource struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

func NewNotificationEventSource() *NotificationEventSource {

	ctx, cancel := context.WithCancel(context.Background())

	return &NotificationEventSource{
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}
}

func (s *NotificationEventSource) Start() (chan any, error) {
	ch := make(chan any, 10)

	go func() {
		s.handleNotifications(ch)
		s.wg.Done()
	}()

	return ch, nil
}

func (s *NotificationEventSource) Stop() error {
	s.cancel()
	s.wg.Wait()
	return nil
}

func (s *NotificationEventSource) handleNotifications(ch chan any) {
	client, err := exgo.New()
	if err != nil {
		panic(err.Error())
	}

	req, err := client.GetAllUnreadNotifications()
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("[Exercism Notifications] %d unread notifications. Listening for notification changes.\n", len(req.Results))
	status := NewNotificationStatus(req.Results...)
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(time.Second * 5):
		}

		req, err := client.GetAllUnreadNotifications()
		if err != nil {
			fmt.Printf("[Exercism Notifications] error getting unread notifications: %v", err)
			continue
		}

		diff := status.Update(req.Results...)

		for _, a := range diff.Added {
			ch <- events.NotificationEvent{
				Notification: a,
				Type:         events.NotificationAdded,
			}
		}

		for _, d := range diff.Removed {
			ch <- events.NotificationEvent{
				Notification: d,
				Type:         events.NotificationDeleted,
			}
		}
	}
}
