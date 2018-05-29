package localsharelib

import "context"
import "github.com/grandcat/zeroconf"
import "log"
import "strconv"
import "time"

type fileList struct {
	hash  string
	files []RemoteFile
}

var prevFileLists map[string]fileList = map[string]fileList{}

func (instance *LocalshareInstance) startMdnsService() {
	instance.peerId = "localshare" + getIps() + ":" + strconv.Itoa(instance.port)
	go func() {
		var err error
		instance.mdnsServer, err = zeroconf.Register(instance.peerId,
			"_edsfoobar._tcp",
			"local.",
			instance.port,
			[]string{},
			nil)
		if err != nil {
			panic(err)
		}
		defer instance.mdnsServer.Shutdown()
		<-instance.ctx.Done()
	}()
}

func (instance *LocalshareInstance) startMdnsDiscoverer() {
	go func() {
		for {
			select {
			case <-instance.ctx.Done():
				break
			default:
				instance.Peers = []*Peer{}
				instance.query()
			}
		}
	}()
}

func (instance *LocalshareInstance) query() {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		newFileLists := map[string]fileList{}
		for entry := range results {
			if entry.Instance != instance.peerId {
				peer := NewPeer(*entry)
				peerHash := firstOrEmpty(entry.Text)
				// if the hash does not exist or has changed, re-fetch the file list
				// otherwise use the previous version of the file list
				_, exists := prevFileLists[peer.Name]
				if (exists && prevFileLists[peer.Name].hash != peerHash) || !exists {
					log.Println("\n\nFETCHING!\n\n")
					flist, err := peer.ListFiles()
					if err != nil {
						log.Fatalln("Error getting peer's file list.", err)
					}
					newFileLists[peer.Name] = fileList{peerHash, flist}
				} else {
					newFileLists[peer.Name] = prevFileLists[peer.Name]
				}
				peer.FileList = newFileLists[peer.Name].files
				instance.Peers = append(instance.Peers, peer)
			}
			instance.peerListCh <- instance.Peers
		}

		prevFileLists = newFileLists
		log.Println("DONE")
	}(entries)

	ctx, cancel := context.WithTimeout(instance.ctx, 5*time.Second)
	defer cancel()
	err = resolver.Browse(ctx, "_edsfoobar._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}
	<-ctx.Done()
}

func (instance *LocalshareInstance) PeerListChannel() chan []*Peer {
	return instance.peerListCh
}
