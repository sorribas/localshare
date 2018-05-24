package localsharelib

import "encoding/json"
import "github.com/grandcat/zeroconf"
import "github.com/parnurzeal/gorequest"
import "io"
import "strconv"

type Peer struct {
	Name  string `json:"name"`
	entry zeroconf.ServiceEntry
}

type RemoteFile struct {
	Name string `json:"name"`
}

func NewPeer(entry zeroconf.ServiceEntry) *Peer {
	return &Peer{Name: entry.Instance, entry: entry}
}

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
