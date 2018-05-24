package webui

import "encoding/base64"
import "fmt"
import "os"
import "os/user"
import "path"
import "github.com/zserge/webview"
import "github.com/sorribas/localshare/internal/localsharelib"

type LocalShareWebBindings struct {
	lsi   *localsharelib.LocalshareInstance
	w     webview.WebView
	Peers []*localsharelib.Peer    `json:"peers"`
	Files []serializableSharedFile `json:"files"`
}

type serializableSharedFile struct {
	Name string `json:"name"`
}

func Start(lsi *localsharelib.LocalshareInstance) {
	html := MustAsset("frontend/index.html")
	b64 := base64.StdEncoding.EncodeToString(html)
	w := webview.New(webview.Settings{
		Title:     "LocalShare",
		URL:       "data:text/html;base64," + b64,
		Width:     800,
		Height:    600,
		Resizable: true,
		Debug:     true,
	})

	w.InjectCSS(string(MustAsset("frontend/style.css")))
	w.Eval(string(MustAsset("frontend/bundle.js")))

	lsi.AddFile(localsharelib.NewInMemoryFile("test", []byte("tst")))
	lswb := &LocalShareWebBindings{lsi, w, []*localsharelib.Peer{}, []serializableSharedFile{}}
	w.Dispatch(func() {
		w.Bind("localshare", lswb)
	})

	go lswb.listenForPeers()

	w.Run()
}

func (lswb *LocalShareWebBindings) Download(peerName string, fileName string) {
	for _, peer := range lswb.lsi.Peers {
		if peer.Name == peerName {
			usr, _ := user.Current()
			f, _ := os.Create(path.Join(usr.HomeDir, "Downloads", fileName))
			defer f.Close()
			peer.DownloadFile(fileName, f)
			break
		}
	}
}

func (lswb *LocalShareWebBindings) ChooseFile() {
	filePath := lswb.w.Dialog(webview.DialogTypeOpen, 0, "", "")
	lswb.lsi.AddFile(localsharelib.NewFsFile(filePath, path.Base(filePath)))
	lswb.Files = filesToSerializable(lswb.lsi.SharedFiles())
	lswb.w.Dispatch(func() {
		lswb.w.Bind("localshare", lswb)
		lswb.w.Eval("window.update()")
	})
}

func (lswb *LocalShareWebBindings) listenForPeers() {
	go func() {
		ch := lswb.lsi.PeerListChannel()
		for {
			<-ch
			lswb.Peers = lswb.lsi.Peers
			fmt.Println("new peer list", lswb.Peers)
			lswb.w.Dispatch(func() {
				lswb.w.Bind("localshare", lswb)
				lswb.w.Eval("window.update()")
			})
		}
	}()
}

func filesToSerializable(files map[string]localsharelib.File) []serializableSharedFile {
	result := []serializableSharedFile{}
	for _, file := range files {
		result = append(result, serializableSharedFile{Name: file.Name()})
	}

	return result
}
