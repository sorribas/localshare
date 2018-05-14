package main

import "bytes"
import "fmt"
import "github.com/sorribas/localshare/internal/localsharelib"

func main() {
	var buf bytes.Buffer
	lsi := localsharelib.NewLocalshareInstance()
	lsi.Start()
	lsi.AddFile(localsharelib.NewInMemoryFile("test", []byte("tst")))

	ch := lsi.PeerChannel()
	for {
		peer := <-ch
		files, err := peer.ListFiles()
		fmt.Println(files, err)
		peer.DownloadFile(files[0].Name, &buf)
		fmt.Println(buf.String())
	}
}
