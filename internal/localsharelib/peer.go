package localsharelib

import "encoding/json"
import "github.com/grandcat/zeroconf"
import "github.com/parnurzeal/gorequest"
import "github.com/sorribas/localshare/internal/writercounter"
import "io"
import "strconv"
import "time"

type Peer struct {
	Name     string       `json:"name"`
	FileList []RemoteFile `json:"files"`
	entry    zeroconf.ServiceEntry
}

type RemoteFile struct {
	Name string `json:"name"`
	Size string `json:"size"`
}

func NewPeer(entry zeroconf.ServiceEntry) *Peer {
	return &Peer{Name: entry.Instance, entry: entry}
}

// Get the list of files from the peer's http api.
func (peer *Peer) ListFiles() ([]RemoteFile, error) {
	r := []RemoteFile{}
	address := "http://" + peer.entry.AddrIPv4[0].String() + ":" + strconv.Itoa(peer.entry.Port)
	_, body, errs := gorequest.New().Get(address + "/api/files").End()
	if len(errs) > 0 {
		return []RemoteFile{}, errs[0]
	}

	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		return []RemoteFile{}, err
	}

	return r, nil
}

func (peer *Peer) DownloadFile(name string, w io.Writer) error {
	address := "http://" + peer.entry.AddrIPv4[0].String() + ":" + strconv.Itoa(peer.entry.Port)
	agent := gorequest.New()
	req, err := agent.Get(address + "/api/files/" + name).MakeRequest()
	if err != nil {
		return err
	}

	resp, err := agent.Client.Do(req)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (peer *Peer) DownloadFileWithProgress(name string, w io.Writer, progress chan int64) {
	closed := false
	wc := writercounter.NewWriterCounter(w)
	go func() {
		peer.DownloadFile(name, wc)
		closed = true
		close(progress)
	}()

	go func() {
		for {
			time.Sleep(250 * time.Millisecond)
			if closed {
				break
			}
			progress <- wc.Count
		}
	}()
}
