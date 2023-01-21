// implements resource driver controller

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/keith-cullen/dractrlplugin/controller/pkg/controller"
)

const numWorkers = 2

var flags *flag.FlagSet

func usage() {
        writer := flags.Output()
        fmt.Fprintf(writer, "Usage: %s [OPTIONS]...\n", os.Args[0])
        flags.PrintDefaults()
}

func main() {
	flags = flag.NewFlagSet("", flag.ExitOnError)
	kubeconfigPath := flags.String("f", os.Getenv("KUBECONFIG"), "filename")
	flags.Usage = usage
	flags.Parse(os.Args[1:])  // ExitOnError so no need to check the return value
	ctrl, err := controller.New(*kubeconfigPath)
	if err != nil {
		return
	}
	ctrl.Run(context.Background(), numWorkers)
}
