package main

import (
	"flag"
	"fmt"

	"github.com/hysios/journalctl-go"
)

var (
	follow      bool
	serviceUint string
	cursor      string
	verbose     bool
)

func init() {
	flag.BoolVar(&follow, "follow", true, "follow journalctl")
	flag.BoolVar(&follow, "f", true, "follow journalctl")
	flag.BoolVar(&verbose, "v", true, "verbose")
	flag.StringVar(&serviceUint, "service", "nginx.service", "service name")
	flag.StringVar(&serviceUint, "u", "nginx.service", "service name")
	flag.StringVar(&cursor, "cursor", "", "cursor")

}

func main() {
	var (
		client journalctl.Client
		opts   journalctl.Options
	)

	flag.Parse()

	if follow {
		opts.Follow = true
		if cursor != "" {
			opts.Cursor = cursor

		} else {
			opts.Skip = -10
			opts.Limit = 20
		}
	}

	entiriesCh, err := client.Entries(&journalctl.Entry{SYSTEMDUNIT: serviceUint}, &opts)
	if err != nil {
		fmt.Println(err)
		return
	}

	for entry := range entiriesCh {
		if verbose {
			fmt.Printf("cursor: %s\n", entry.CURSOR)
			fmt.Printf("sessionId: %s\n", entry.AUDITSESSION)
		}
		fmt.Println(entry.Message())
	}
}
