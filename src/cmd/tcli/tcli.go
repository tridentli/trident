///usr/bin/env go run $0 "$@"; exit

/*
 * Trident CLI - Tickly (tcli)
 *
 * This is effectively a HTTP client for tridentd.
 * All requests are sent over HTTP, there is no access directly to anything.
 *
 * This client also serves as an example on how to talk to the Trident API.
 *
 * tcli stores a token in ~/.trident_token for retaining the logged-in state.
 *
 * Custom environment variables:
 * - Select a custom token file with:
 *     TRIDENT_TOKEN=/other/path/to/tokenfile
 *   This is useful if you want to have multiple identities
 *   or want to keep a token around that has the sysadmin bit set
 *
 * - Enable verbosity with:
 *     TRIDENT_VERBOSE=<anything>
 *
 * - Disable verbosity with
 *     TRIDENT_VERBOSE=off
 *   or unset the environment variable
 *
 * - Select different server with:
 *     TRIDENT_SERVER=https://trident.example.net
 *
 * Acceptable command line options can be requested with -help
 */

package main

import (
	"trident.li/pitchfork/cmd/cli"
)

func main() {
	pf_cmd_cli.CLI(".trident_token", "TRIDENT_TOKEN", "TRIDENT_VERBOSE", "TRIDENT_SERVER", "http://127.0.0.1:8334")
}
