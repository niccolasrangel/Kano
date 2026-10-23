// Kano-server is the server build of Kano: kanod and the kano CLI,
// without the app, for a machine whose sessions Kano on another one shows
// (Settings › Hosts), as a Linux box with no display. It installs in
// ~/.local with install-server.sh, which runs kanod as a service, and
// updates itself (kano daemon update, or the app's Update Kano).
//
// Its version comes with the build:
//
//	go build -ldflags "-X main.version=2.1.0 -X github.com/egoist/mygo.production=1" -o kano ./cmd/kano-server
package main

import (
	"fmt"
	"os"

	"github.com/egoist/mygo"

	"kano/internal/build"
	"kano/internal/cli"
	"kano/internal/daemon"
)

// version is the version of Kano this build is of, set with -ldflags.
var version = "dev"

func main() {
	build.Server = true
	mygo.App.SetVersion(version)
	if len(os.Args) > 1 && os.Args[1] == "__kanod" {
		if err := daemon.Run(daemon.Options{}); err != nil {
			fmt.Fprintln(os.Stderr, "kanod:", err)
			os.Exit(1)
		}
		return
	}
	if len(os.Args) < 2 {
		os.Exit(cli.Run([]string{"help"}))
	}
	if !cli.IsCommand(os.Args[1]) {
		fmt.Fprintf(os.Stderr, "kano: unknown command %q (see kano help)\n", os.Args[1])
		os.Exit(2)
	}
	os.Exit(cli.Run(os.Args[1:]))
}
