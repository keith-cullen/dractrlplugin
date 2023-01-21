// implements resource driver controller

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/keith-cullen/dractrlplugin/plugin/pkg/plugin"
)

func main() {
	plugin := plugin.New()
	err := plugin.Run()
	if err != nil {
		return
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigCh
}
