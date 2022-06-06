package consumer

type Consumer interface {
	Start(<-chan interface{}) error
	Stop() error
}
