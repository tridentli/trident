///usr/bin/env go run -ldflags "-X main.Version=shell.run" $0 "$@"; exit

package main

import (
	"trident.li/pitchfork/cmd/server"
	tf "trident.li/trident/src/lib"
	tu "trident.li/trident/src/ui"
)

var Version = "unconfigured"

func main() {
	pf_cmd_server.Serve("trident", "Trident", Version, tf.Copyright, tf.Website, tf.AppSchemaVersion, tu.NewTriUI, nil)
}
