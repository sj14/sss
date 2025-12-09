package main

import (
	"context"
	"log"
	"os"

	"github.com/sj14/sss/cli"
	"github.com/sj14/sss/util"
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

	ver := util.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	}

	if err := cli.Exec(
		ctx,
		os.Stdout,
		os.Stderr,
		ver,
	); err != nil {
		log.Fatalln(err)
	}
}
