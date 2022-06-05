package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/gen2brain/beeep"
)

var wordPtr *string

func init() {
	wordPtr = flag.String("tracks", "", "tracks to be notified")
}

func main() {

	flag.Parse()
	track_slugs := strings.Split(*wordPtr, ",")

	ch := TrackMentoringRequests(track_slugs...)

	for event := range ch {

		var title, description string

		request := event.Request

		err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		if err != nil {
			fmt.Print("could not play notification beep")
		}

		switch event.Type {
		case NewMentoringRequest:
			title = fmt.Sprintf("[%s] New Solution", event.Track)
			description = fmt.Sprintf("%s by %s", request.ExerciseTitle, request.StudentHandle)
			err := beeep.Notify(title, description, "assets/exercism.png")
			if err != nil {
				fmt.Printf("could not send notification of solution added: %s\n", err.Error())
			}

		case MentoringRequestDeleted:
			title = fmt.Sprintf("[%s] Solution Mentored", event.Track)
			description = fmt.Sprintf("%s by %s", request.ExerciseTitle, request.StudentHandle)
			err := beeep.Notify(title, description, "assets/exercism.png")
			if err != nil {
				fmt.Printf("could not send notification of solution mentored: %s\n", err.Error())
			}
		}
		fmt.Printf("%s: %s\n", title, description)

	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}
