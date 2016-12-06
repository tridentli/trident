///usr/bin/env go run -ldflags "-X main.Version=shell.run" $0 "$@"; exit

/*
 * Trident Setup (tsetup)
 *
 * tsetup is only meant for initial setup tasks.
 * It should be run as the 'postgres' user.
 *
 * For general use, use 'tcli' or the webinterface and log in.
 */

package main

import (
	"trident.li/pitchfork/cmd/setup"
	tr "trident.li/trident/src/lib"
)

var Version = "unconfigured"

func main() {
	pf_cmd_setup.Setup("tsetup", "trident", tr.AppName, Version, tr.Copyright, tr.Website, tr.AppSchemaVersion, "TRIDENT_SERVER", "http://127.0.0.1:8333", tr.NewTriCtx)
}
