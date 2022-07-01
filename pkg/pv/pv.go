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

//New Configures and returns a new PeeVee
func New[T any](PVname string, ops ...PVOptions[T]) PeeVee[T] {
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

//NewReaderWrap Returns a new PeeVee configured to wrap an existing readable channel.
//
// Use this when you have a channel you don't control but still want
// to have PeeVee's benefits for it.
func NewReaderWrap[T any](PVname string, readableChan chan T, ops ...PVOptions[T]) PeeVee[T] {
	pv := PeeVee[T]{
		Name:     PVname,
		readChan: make(chan T),

		// this looks odd, I know but if we think about the implementation
		// we will notice the channel the user reads from is the channel that
		// is receiving the writes
		writeChan: readableChan,
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
