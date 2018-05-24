package main

import "github.com/sorribas/localshare/internal/localsharelib"
import "github.com/sorribas/localshare/internal/webui"

func main() {
	lsi := localsharelib.NewLocalshareInstance()
	// gui := &ui.UI{Ls: &lsi}
	// gui.Start()

	lsi.Start()
	webui.Start(&lsi)
}
