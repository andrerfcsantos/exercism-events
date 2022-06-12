package pushovernotifier

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/andrerfcsantos/exercism-events/events"
	"github.com/gregdel/pushover"
)

// PushoverNotifier implements the consumer interface
type PushoverNotifier struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wg        *sync.WaitGroup
	push      *pushover.Pushover
	recipient *pushover.Recipient
	tracks    map[string]struct{}
}

func NewPushoverNotifier(tracks map[string]struct{}) *PushoverNotifier {
	ctx, cancel := context.WithCancel(context.Background())
	return &PushoverNotifier{
		ctx:       ctx,
		cancel:    cancel,
		wg:        &sync.WaitGroup{},
		push:      pushover.New(os.Getenv("PUSHOVER_TOKEN")),
		recipient: pushover.NewRecipient(os.Getenv("PUSHOVER_USER")),
		tracks:    tracks,
	}

}

func (d *PushoverNotifier) Start(ch <-chan interface{}) error {

	d.wg.Add(1)
	go func() {
		d.listen(ch)
		d.wg.Done()
	}()
	return nil
}

func (d *PushoverNotifier) Stop() error {
	d.cancel()
	d.wg.Wait()
	return nil
}

func (d *PushoverNotifier) listen(ch <-chan interface{}) {
	for {
		select {
		case <-d.ctx.Done():
			return
		case msg := <-ch:
			d.handleEvent(msg)
		}
	}
}

func (d *PushoverNotifier) handleEvent(event interface{}) {
	switch ev := event.(type) {
	case events.MentoringEvent:

		if _, ok := d.tracks[ev.Track]; !ok && len(d.tracks) > 0 {
			return
		}

		request := ev.Request

		var title, description string
		switch ev.Type {
		case events.NewMentoringRequest:
			title = fmt.Sprintf("[%s] New Solution", ev.Track)
			description = fmt.Sprintf("%s by %s", request.ExerciseTitle, request.StudentHandle)

			message := pushover.NewMessageWithTitle(description, title)
			_, err := d.push.SendMessage(message, d.recipient)
			if err != nil {
				fmt.Printf("Error sending message: %s\n", err)
			}

		case events.MentoringRequestDeleted:
			title = fmt.Sprintf("[%s] Solution Mentored", ev.Track)
			description = fmt.Sprintf("%s by %s", request.ExerciseTitle, request.StudentHandle)
			message := pushover.NewMessageWithTitle(description, title)
			_, err := d.push.SendMessage(message, d.recipient)
			if err != nil {
				fmt.Printf("Error sending message: %s\n", err)
			}
		}

	case events.NotificationEvent:

		switch ev.Type {
		case events.NotificationAdded:
			title := "Exercism Notification"
			message := pushover.NewMessageWithTitle(ev.Notification.Text, title)
			_, err := d.push.SendMessage(message, d.recipient)
			if err != nil {
				fmt.Printf("Error sending message: %s\n", err)
			}
		}
	}

}
