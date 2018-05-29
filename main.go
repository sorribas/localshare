package main

import "github.com/sorribas/localshare/internal/localsharelib"
import "github.com/sorribas/localshare/internal/webui"

func main() {
	lsi := localsharelib.NewLocalshareInstance()
	lsi.Start()
	webui.Start(&lsi)
}
