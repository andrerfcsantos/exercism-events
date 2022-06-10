package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/andrerfcsantos/exercism-events/consumer"
	"github.com/andrerfcsantos/exercism-events/consumer/desktopnotifier"
	"github.com/andrerfcsantos/exercism-events/forward"
	"github.com/andrerfcsantos/exercism-events/source"
	"github.com/andrerfcsantos/exercism-events/source/mentoring"
	"github.com/andrerfcsantos/exercism-events/source/notifications"
)

var wordPtr *string

func init() {
	wordPtr = flag.String("tracks", "", "tracks to be notified")
}

func main() {

	flag.Parse()
	track_slugs := strings.Split(*wordPtr, ",")

	// Sources
	mentoringSource := mentoring.NewMentoringEventSource(track_slugs...)
	notificationSource := notifications.NewNotificationEventSource()

	// Consumers
	notifierConsumer := desktopnotifier.NewDesktopNotifier()

	fw := forward.NewForwarder(
		[]source.Source{
			mentoringSource,
			notificationSource,
		},
		[]consumer.Consumer{
			notifierConsumer,
		},
	)

	err := fw.Start()
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	err = fw.Stop()
	if err != nil {
		log.Fatal(err)
	}
}
