package peevee

//PeeVee Representation of the PV
type PeeVee[T any] struct {
	Name string

	// readChan This is the channel the user will read from
	readChan chan T

	// writeChan This is the channel the user will write to
	writeChan chan T
}

// GetReadableChan Returns the readable channel. This is a safe way to expose the channel
// forcing the return to be only readable whilst keeping an internal reference to the "open"
// version of the channel.
func (pv *PeeVee[T]) GetReadableChan() <-chan T {
	return pv.readChan
}

// GetWritableChan Returns the writable channel. This is a safe way to expose the channel
// forcing the return to be only writable whilst keeping an internal reference to the "open"
// version of the channel.
func (pv *PeeVee[T]) GetWritableChan() chan<- T {
	return pv.writeChan
}

//NewPeeVee Configures and returns a new PeeVee
func NewPeeVee[T any](PVname string, ops ...PVOptions[T]) PeeVee[T] {
	pv := PeeVee[T]{
		Name:      PVname,
		readChan:  make(chan T),
		writeChan: make(chan T),
	}

	if len(ops) == 0 {
		// by default, maintain normal behaviour
		ops = append(ops, WithDefault[T]())
	}

	for _, option := range ops {
		option(&pv)
	}

	return pv
}
