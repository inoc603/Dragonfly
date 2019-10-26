package dfget

import "github.com/dragonflyoss/Dragonfly/dfget/config"

type Dfget struct {
	config.Config

	scheduler Scheduler
}

type Option func(g *Dfget) error

func New(opts ...Option) (*Dfget, error) {
	g := &Dfget{}
	for _, opt := range opts {
		if err := opt(g); err != nil {
			return nil, err
		}
	}
	return g, nil
}

func WithExample() Option {
	return func(g *Dfget) error {
		return nil
	}
}

func (g *Dfget) Download(task *Task) error {
	return nil
}

func (g *Dfget) StartPeerServer() (started bool, err error) {
	return false, nil
}
