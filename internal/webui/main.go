package webui

import "encoding/base64"
import "fmt"
import "os"
import "os/user"
import "path"
import "github.com/zserge/webview"
import "github.com/sorribas/localshare/internal/localsharelib"
import "strconv"

type LocalShareWebBindings struct {
	lsi       *localsharelib.LocalshareInstance
	w         webview.WebView
	Peers     []*localsharelib.Peer    `json:"peers"`
	Files     []serializableSharedFile `json:"files"`
	Downloads map[string]float64       `json:"downloads"`
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
	lswb := &LocalShareWebBindings{
		lsi,
		w,
		[]*localsharelib.Peer{},
		[]serializableSharedFile{},
		map[string]float64{},
	}
	lswb.updateFrontend()

	go lswb.listenForPeers()

	w.Run()
}

func (lswb *LocalShareWebBindings) Download(peerName string, fileName string) {
	go func() {
		lswb.Downloads[peerName+"|"+fileName] = 0
		for _, peer := range lswb.lsi.Peers {
			if peer.Name == peerName {

				// find the file size

				var fileSize int64
				for _, file := range peer.FileList {
					if file.Name == fileName {
						fileSize, _ = strconv.ParseInt(file.Size, 10, 64)
					}
				}

				usr, _ := user.Current()
				f, _ := os.Create(path.Join(usr.HomeDir, "Downloads", fileName))
				defer f.Close()
				ch := make(chan int64)
				peer.DownloadFileWithProgress(fileName, f, ch)
				for progress := range ch {
					lswb.Downloads[peerName+"|"+fileName] = float64(progress) / float64(fileSize)
					lswb.updateFrontend()
				}
				lswb.Downloads[peerName+"|"+fileName] = float64(1)
				lswb.updateFrontend()
				break
			}
		}
	}()
}

func (lswb *LocalShareWebBindings) ChooseFile() {
	filePath := lswb.w.Dialog(webview.DialogTypeOpen, 0, "", "")
	lswb.lsi.AddFile(localsharelib.NewFsFile(filePath, path.Base(filePath)))
	lswb.Files = filesToSerializable(lswb.lsi.SharedFiles())
	lswb.updateFrontend()
}

func (lswb *LocalShareWebBindings) listenForPeers() {
	go func() {
		ch := lswb.lsi.PeerListChannel()
		for {
			<-ch
			lswb.Peers = lswb.lsi.Peers
			fmt.Println("new peer list", lswb.Peers)
			lswb.updateFrontend()
		}
	}()
}

func (lswb *LocalShareWebBindings) updateFrontend() {
	lswb.w.Dispatch(func() {
		lswb.w.Bind("localshare", lswb)
		lswb.w.Eval("window.update()")
	})
}

func filesToSerializable(files map[string]localsharelib.File) []serializableSharedFile {
	result := []serializableSharedFile{}
	for _, file := range files {
		result = append(result, serializableSharedFile{Name: file.Name()})
	}

	return result
}
