package localsharelib

import "context"

type LocalshareInstance struct {
	ctx    context.Context
	cancel context.CancelFunc
	port   int
	files  map[string]File
	Peers  []*Peer
	peerId string
	peerCh chan *Peer
}

func NewLocalshareInstance() LocalshareInstance {
	return NewLocalshareInstanceWithContext(context.Background())
}

func NewLocalshareInstanceWithContext(ctx context.Context) LocalshareInstance {
	instance := LocalshareInstance{}
	instance.ctx, instance.cancel = context.WithCancel(ctx)
	instance.files = map[string]File{}
	instance.peerCh = make(chan *Peer)
	return instance
}

func (instance *LocalshareInstance) Start() {
	instance.startHttpServer()
	instance.startMdnsService()
	instance.startMdnsDiscoverer()
}

func (instance *LocalshareInstance) Stop() {
	instance.cancel()
}
