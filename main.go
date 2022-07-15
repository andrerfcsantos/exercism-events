package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"

	"github.com/andrerfcsantos/exercism-events/consumer"
	"github.com/andrerfcsantos/exercism-events/consumer/database"
	"github.com/andrerfcsantos/exercism-events/consumer/desktopnotifier"
	"github.com/andrerfcsantos/exercism-events/consumer/pushovernotifier"
	"github.com/andrerfcsantos/exercism-events/forward"
	"github.com/andrerfcsantos/exercism-events/source"
	"github.com/andrerfcsantos/exercism-events/source/mentoring"
	"github.com/andrerfcsantos/exercism-events/source/notifications"

	log "github.com/sirupsen/logrus"
)

var (
	tracksFlag         *string
	sourcesFlag        *string
	consumersFlag      *string
	pushoverTracksFlag *string
)

func init() {
	tracksFlag = flag.String("tracks", "", "tracks to be notified")
	sourcesFlag = flag.String("sources", "notifications,mentoring", "sources to be used")
	consumersFlag = flag.String("consumers", "desktopnotifier", "consumers to be used")
	pushoverTracksFlag = flag.String("pushovertracks", "", "pushover tracks to be notified")
}

type FlagInfo struct {
	Tracks         []string
	Sources        map[string]struct{}
	Consumers      map[string]struct{}
	PushoverTracks map[string]struct{}
}

func ParseFlags() FlagInfo {
	flag.Parse()
	tracks := strings.Split(*tracksFlag, ",")
	sources := strings.Split(*sourcesFlag, ",")
	consumers := strings.Split(*consumersFlag, ",")
	pushoverTracks := strings.Split(*pushoverTracksFlag, ",")

	sourcesMap := make(map[string]struct{})
	consumersMap := make(map[string]struct{})
	pushoverTracksMap := make(map[string]struct{})

	for _, source := range sources {
		sourcesMap[source] = struct{}{}
	}

	for _, consumer := range consumers {
		consumersMap[consumer] = struct{}{}
	}

	for _, track := range pushoverTracks {
		pushoverTracksMap[track] = struct{}{}
	}

	return FlagInfo{
		Tracks:         tracks,
		Sources:        sourcesMap,
		Consumers:      consumersMap,
		PushoverTracks: pushoverTracksMap,
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

	if _, ok := flagInfo.Consumers["pushover"]; ok {
		consumers = append(consumers, pushovernotifier.NewPushoverNotifier(flagInfo.PushoverTracks))
	}

	return consumers
}

func main() {
	flagInfo := ParseFlags()

	file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	sources := GetSources(flagInfo)
	consumers := GetConsumers(flagInfo)

	fw := forward.NewForwarder(sources, consumers)

	err = fw.Start()
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
