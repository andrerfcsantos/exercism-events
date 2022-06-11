package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/andrerfcsantos/exercism-events/consumer/database/repository"
	"github.com/andrerfcsantos/exercism-events/events"
)

// Database implements the consumer interface
type Database struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
	r      *repository.Repository
}

func NewDatabase() *Database {

	ctx, cancel := context.WithCancel(context.Background())
	return &Database{
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}
}

func (d *Database) Start(ch <-chan interface{}) error {
	var err error

	d.r, err = repository.New()
	if err != nil {
		return fmt.Errorf("starting repository: %w", err)
	}

	d.wg.Add(1)
	go func() {
		d.listen(ch)
		d.wg.Done()
	}()
	return nil
}

func (d *Database) Stop() error {
	d.cancel()
	d.wg.Wait()
	return nil
}

func (d *Database) listen(ch <-chan interface{}) {
	for {
		select {
		case <-d.ctx.Done():
			return
		case msg := <-ch:
			d.handleEvent(msg)
		}
	}
}

func (d *Database) handleEvent(event interface{}) {
	switch ev := event.(type) {
	case events.MentoringEvent:
		request := ToRequest(ev)
		err := d.r.SaveMentoringRequest(request)
		if err != nil {
			fmt.Printf("error saving request: %v\n", err)
		}
	}

}

func ToRequest(event events.MentoringEvent) repository.MentoringRequest {
	action := ""
	switch event.Type {
	case events.NewMentoringRequest:
		action = "created"
	case events.MentoringRequestDeleted:
		action = "deleted"
	}

	return repository.MentoringRequest{
		UUID:             event.Request.UUID,
		TrackTitle:       event.Request.TrackTitle,
		ExerciseIconURL:  event.Request.ExerciseIconURL,
		ExerciseTitle:    event.Request.ExerciseTitle,
		StudentHandle:    event.Request.StudentHandle,
		StudentAvatarUrl: event.Request.StudentAvatarURL,
		UpdatedAt:        event.Request.UpdatedAt,
		AddedAt:          time.Now(),
		Action:           action,
	}
}
