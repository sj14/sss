package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sj14/sss/cli"
)

var (
	// will be replaced during the build process
	version = "undefined"
	commit  = "undefined"
	date    = "undefined"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ver := fmt.Sprintf("version: %s | commit: %s | date: %s", version, commit, date)

	if err := cli.Exec(
		ctx,
		os.Stdout,
		os.Stderr,
		ver,
	); err != nil {
		log.Fatalln(err)
	}
}
