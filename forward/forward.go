package forward

import (
	"fmt"

	"github.com/andrerfcsantos/exercism-events/consumer"
	"github.com/andrerfcsantos/exercism-events/source"
)

type Forwarder struct {
	sources   []source.Source
	consumers []consumer.Consumer
}

func NewForwarder(sources []source.Source, consumers []consumer.Consumer) *Forwarder {
	return &Forwarder{
		sources:   sources,
		consumers: consumers,
	}
}

func (f *Forwarder) Start() error {
	consumerChans := make([]chan interface{}, 0, len(f.consumers))
	for _, c := range f.consumers {
		consumerChan := make(chan interface{}, 10)
		consumerChans = append(consumerChans, consumerChan)
		go c.Start(consumerChan)
	}

	fmt.Printf("[Forwarder] Listening from %d sources\n", len(f.sources))
	for _, s := range f.sources {
		go func(source source.Source) {
			ch, err := source.Start()
			if err != nil {
				fmt.Printf("could not start source: %s", err)
				return
			}

			for ev := range ch {
				ev := ev
				for _, c := range consumerChans {
					c <- ev
				}
			}

		}(s)
	}

	return nil
}

func (f *Forwarder) Stop() error {

	var errors []error

	for _, s := range f.sources {
		err := s.Stop()
		if err != nil {
			errors = append(errors, err)
		}
	}

	for _, c := range f.consumers {
		err := c.Stop()
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("could not stop all sources and consumers: %s", errors)
	}

	return nil
}
