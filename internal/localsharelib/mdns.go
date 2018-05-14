package localsharelib

import "github.com/grandcat/zeroconf"
import "log"
import "strconv"

func (instance *LocalshareInstance) startMdnsService() {
	instance.peerId = "localshare" + getIps() + ":" + strconv.Itoa(instance.port)
	go func() {
		server, err := zeroconf.Register(instance.peerId,
			"_edsfoobar._tcp",
			"local.",
			instance.port,
			[]string{},
			nil)
		if err != nil {
			panic(err)
		}
		defer server.Shutdown()
		<-instance.ctx.Done()
	}()
}

func (instance *LocalshareInstance) startMdnsDiscoverer() {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			if entry.Instance != instance.peerId {
				peer := NewPeer(*entry)
				instance.peers = append(instance.peers, peer)
				instance.peerCh <- peer
			}
		}
	}(entries)

	err = resolver.Browse(instance.ctx, "_edsfoobar._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}
}

func (instance *LocalshareInstance) PeerChannel() chan *Peer {
	return instance.peerCh
}
