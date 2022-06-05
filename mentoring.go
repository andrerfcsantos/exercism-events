package main

import (
	"fmt"
	"time"

	"github.com/andrerfcsantos/exercism-events/exgo"
	"github.com/andrerfcsantos/exercism-events/track"
)

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

func handleTrack(track_slug string, ch chan MentoringEvent) {
	client, err := exgo.New()
	if err != nil {
		panic(err.Error())
	}

	req, err := client.GetAllMentoringRequests(track_slug)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Intializing with %d requests\n", len(req.Results))
	status := track.NewMentoringStatus(req.Results...)
	for {
		time.Sleep(time.Second * 5)
		req, err := client.GetAllMentoringRequests(track_slug)
		if err != nil {
			panic(err.Error())
		}

		diff := status.Update(req.Results...)

		for _, a := range diff.Added {
			ch <- MentoringEvent{
				Track:   a.TrackTitle,
				Type:    NewMentoringRequest,
				Request: a,
			}
		}

		for _, d := range diff.Removed {
			ch <- MentoringEvent{
				Track:   d.TrackTitle,
				Type:    MentoringRequestDeleted,
				Request: d,
			}
		}

	}
}

func TrackMentoringRequests(tracks ...string) chan MentoringEvent {

	ch := make(chan MentoringEvent, 10)

	for _, track := range tracks {
		go handleTrack(track, ch)
	}

	return ch
}
