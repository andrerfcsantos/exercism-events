package mentoring

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/andrerfcsantos/exercism-events/events"
	"github.com/andrerfcsantos/exercism-events/exgo"
)

type MentoringEventSource struct {
	ctx    context.Context
	cancel context.CancelFunc
	tracks []string
	wg     *sync.WaitGroup
}

func NewMentoringEventSource(tracks ...string) *MentoringEventSource {

	ctx, cancel := context.WithCancel(context.Background())

	return &MentoringEventSource{
		tracks: tracks,
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}
}

func (s *MentoringEventSource) Start() (chan any, error) {
	ch := make(chan any, 10)

	s.wg.Add(len(s.tracks))

	for _, track := range s.tracks {
		go func(t string) {
			s.handleTrack(t, ch)
			s.wg.Done()
		}(track)
	}

	return ch, nil
}

func (s *MentoringEventSource) Stop() error {
	s.cancel()
	s.wg.Wait()
	return nil
}

func (s *MentoringEventSource) handleTrack(track_slug string, ch chan any) {
	client, err := exgo.New()
	if err != nil {
		panic(err.Error())
	}

	req, err := client.GetAllMentoringRequests(track_slug)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("[Mentoring Requests] Found %d requests. Listening for requests changes.\n", len(req.Results))
	status := NewMentoringStatus(req.Results...)
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(time.Second * 10):

		}
		req, err := client.GetAllMentoringRequests(track_slug)
		if err != nil {
			fmt.Printf("[Mentoring Requests] error getting mentoring requests: %v", err)
			continue
		}

		diff := status.Update(req.Results...)

		for _, a := range diff.Added {
			ch <- events.MentoringEvent{
				Track:   a.TrackTitle,
				Type:    events.NewMentoringRequest,
				Request: a,
			}
		}

		for _, d := range diff.Removed {
			ch <- events.MentoringEvent{
				Track:   d.TrackTitle,
				Type:    events.MentoringRequestDeleted,
				Request: d,
			}
		}

	}
}
