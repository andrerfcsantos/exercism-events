package source

type Source interface {
	Start() (chan interface{}, error)
	Stop() error
}
