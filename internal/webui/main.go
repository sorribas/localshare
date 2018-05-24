package webui

import "fmt"
import "os"
import "os/user"
import "path"
import "github.com/zserge/webview"
import "github.com/sorribas/localshare/internal/localsharelib"

type LocalShareWebBindings struct {
	lsi   *localsharelib.LocalshareInstance
	w     webview.WebView
	Peers []serializablePeer       `json:"peers"`
	Files []serializableSharedFile `json:"files"`
}

type serializablePeer struct {
	Name  string                   `json:"name"`
	Files []serializableRemoteFile `json:"files"`
}

type serializableRemoteFile struct {
	Name string `json:"name"`
}

type serializableSharedFile struct {
	Name string `json:"name"`
}

func Start(lsi *localsharelib.LocalshareInstance) {
	w := webview.New(webview.Settings{
		Title:     "LocalShare",
		URL:       "file:///home/ed/prog/localshareui/index.html",
		Width:     800,
		Height:    600,
		Resizable: true,
		Debug:     true,
	})

	lsi.AddFile(localsharelib.NewInMemoryFile("test", []byte("tst")))
	lswb := &LocalShareWebBindings{lsi, w, []serializablePeer{}, []serializableSharedFile{}}
	w.Dispatch(func() {
		w.Bind("localshare", lswb)
	})

	go lswb.listenForPeers()

	w.Run()
}

func (lswb *LocalShareWebBindings) Download(peerName string, fileName string) {
	fmt.Println("download")
	for _, peer := range lswb.lsi.Peers {
		if peer.Name == peerName {
			usr, _ := user.Current()
			f, _ := os.Create(path.Join(usr.HomeDir, "Downloads", fileName))
			defer f.Close()
			fmt.Println("downloading " + fileName)
			peer.DownloadFile(fileName, f)
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
		ch := lswb.lsi.PeerChannel()
		for {
			<-ch
			lswb.Peers = peersToSerializable(lswb.lsi.Peers)
			fmt.Println("new peer", lswb.Peers)
			lswb.w.Dispatch(func() {
				lswb.w.Bind("localshare", lswb)
				lswb.w.Eval("window.update()")
			})
		}
	}()
}

func peersToSerializable(peers []*localsharelib.Peer) []serializablePeer {
	result := []serializablePeer{}
	for _, peer := range peers {
		files, _ := peer.ListFiles()
		sfiles := []serializableRemoteFile{}
		for _, file := range files {
			sfiles = append(sfiles, serializableRemoteFile{file.Name})
		}

		speer := serializablePeer{peer.Name, sfiles}
		result = append(result, speer)
	}

	return result
}

func filesToSerializable(files map[string]localsharelib.File) []serializableSharedFile {
	result := []serializableSharedFile{}
	for _, file := range files {
		result = append(result, serializableSharedFile{Name: file.Name()})
	}

	return result
}
