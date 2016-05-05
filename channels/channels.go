package channels

import ()

type Channels struct {
	Interface *Interface
}

func InitializeChannels() *Channels {
	return &Channels{
		Interface: &Interface{},
	}
}
