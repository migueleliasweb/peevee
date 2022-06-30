package peevee

type ChannelPiper interface {
	Pipe()
}

// CallbackChannelPiper This piper will allow a callback in between piping items
type CallbackChannelPiper[T any] struct {
	// readChan This is the channel the user will read from.
	ReadChan chan T

	// writeChan This is the channel the user will write to
	WriteChan chan T

	// callback The callback function is called for each item sent to `WriteChan`.
	//
	// Please note that this callback will impact your app performance greatly. So make it fast
	// or at least non-blocking.
	callback func(T)
}

// Pipe pipe items between channels
func (ccp *CallbackChannelPiper[T]) Pipe() {
	for item := range ccp.WriteChan {
		if ccp.callback != nil {
			ccp.callback(item)
		}

		ccp.ReadChan <- item
	}
}
