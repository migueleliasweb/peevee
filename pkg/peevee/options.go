package peevee

type PVOptions[T any] func(*PeeVee[T])

func WithDefault[T any]() PVOptions[T] {
	return WithCallback[T](nil)
}

func WithCallback[T any](f func(T)) PVOptions[T] {
	return func(pv *PeeVee[T]) {
		ccp := CallbackChannelPiper[T]{
			ReadChan:  pv.readChan,
			WriteChan: pv.writeChan,
			callback:  f,
		}

		go ccp.Pipe()
	}
}
