package desktopnotifier

import (
	"context"
	"fmt"
	"sync"

	"github.com/andrerfcsantos/exercism-events/events"
	"github.com/gen2brain/beeep"
)

// DesktopNotifier implements the consumer interface
type DesktopNotifier struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

func NewDesktopNotifier() *DesktopNotifier {
	ctx, cancel := context.WithCancel(context.Background())
	return &DesktopNotifier{
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}
}

func (d *DesktopNotifier) Start(ch <-chan interface{}) error {

	d.wg.Add(1)
	go func() {
		d.listen(ch)
		d.wg.Done()
	}()
	return nil
}

func (d *DesktopNotifier) Stop() error {
	d.cancel()
	d.wg.Wait()
	return nil
}

func (d *DesktopNotifier) listen(ch <-chan interface{}) {
	for {
		select {
		case <-d.ctx.Done():
			return
		case msg := <-ch:
			d.handleEvent(msg)
		}
	}
}

func (d *DesktopNotifier) handleEvent(event interface{}) {
	switch ev := event.(type) {
	case events.MentoringEvent:
		request := ev.Request

		err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		if err != nil {
			fmt.Print("[Desktop Notifier] could not play notification beep")
		}

		var title, description string
		switch ev.Type {
		case events.NewMentoringRequest:
			title = fmt.Sprintf("[%s] New Solution", ev.Track)
			description = fmt.Sprintf("%s by %s", request.ExerciseTitle, request.StudentHandle)
			err := beeep.Notify(title, description, "assets/exercism.png")
			if err != nil {
				fmt.Printf("could not send notification of solution added: %s\n", err.Error())
			}

		case events.MentoringRequestDeleted:
			title = fmt.Sprintf("[%s] Solution Mentored", ev.Track)
			description = fmt.Sprintf("%s by %s", request.ExerciseTitle, request.StudentHandle)
			err := beeep.Notify(title, description, "assets/exercism.png")
			if err != nil {
				fmt.Printf("could not send notification of solution mentored: %s\n", err.Error())
			}
		}

		fmt.Printf("[Desktop Notifier] %s: %s\n", title, description)
	}

}
