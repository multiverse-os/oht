package channels

import ()

type Interface struct {
}

func NewInterface() (i *Interface) {
	return &Interface{}
}

// CHANNELS
func (i *Interface) ListChannels() (channels []string) {
	return
}

func (i *Interface) Channel() (successful bool) {
	return
}

func (i *Interface) JoinChannel(channelId string) (successful bool) {
	return
}

func (i *Interface) LeaveChannel(channelId string) (successful bool) {
	return
}

func (i *Interface) ChannelCast(channelId string, message string) (successful bool) {
	return
}
