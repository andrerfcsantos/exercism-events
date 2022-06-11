package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/andrerfcsantos/exercism-events/consumer"
	"github.com/andrerfcsantos/exercism-events/consumer/database"
	"github.com/andrerfcsantos/exercism-events/consumer/desktopnotifier"
	"github.com/andrerfcsantos/exercism-events/forward"
	"github.com/andrerfcsantos/exercism-events/source"
	"github.com/andrerfcsantos/exercism-events/source/mentoring"
	"github.com/andrerfcsantos/exercism-events/source/notifications"
)

var (
	wordPtrFlag   *string
	sourcesFlag   *string
	consumersFlag *string
)

func init() {
	wordPtrFlag = flag.String("tracks", "", "tracks to be notified")
	sourcesFlag = flag.String("sources", "notifications,mentoring", "sources to be used")
	consumersFlag = flag.String("consumers", "desktopnotifier", "consumers to be used")
}

type FlagInfo struct {
	Tracks    []string
	Sources   map[string]struct{}
	Consumers map[string]struct{}
}

func ParseFlags() FlagInfo {
	flag.Parse()
	tracks := strings.Split(*wordPtrFlag, ",")
	sources := strings.Split(*sourcesFlag, ",")
	consumers := strings.Split(*consumersFlag, ",")

	sourcesMap := make(map[string]struct{})
	consumersMap := make(map[string]struct{})

	for _, source := range sources {
		sourcesMap[source] = struct{}{}
	}

	for _, consumer := range consumers {
		consumersMap[consumer] = struct{}{}
	}

	return FlagInfo{
		Tracks:    tracks,
		Sources:   sourcesMap,
		Consumers: consumersMap,
	}
}

func GetSources(flagInfo FlagInfo) []source.Source {
	sources := []source.Source{}

	if _, ok := flagInfo.Sources["notifications"]; ok {
		sources = append(sources, notifications.NewNotificationEventSource())
	}

	if _, ok := flagInfo.Sources["mentoring"]; ok {
		sources = append(sources, mentoring.NewMentoringEventSource(flagInfo.Tracks...))
	}

	return sources
}

func GetConsumers(flagInfo FlagInfo) []consumer.Consumer {
	consumers := []consumer.Consumer{}

	if _, ok := flagInfo.Consumers["desktopnotifier"]; ok {
		consumers = append(consumers, desktopnotifier.NewDesktopNotifier())
	}

	if _, ok := flagInfo.Consumers["database"]; ok {
		consumers = append(consumers, database.NewDatabase())
	}

	return consumers
}

func main() {
	flagInfo := ParseFlags()

	sources := GetSources(flagInfo)
	consumers := GetConsumers(flagInfo)

	fw := forward.NewForwarder(sources, consumers)

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
