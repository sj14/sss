package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:           "root",
		DefaultCommand: "subcommand",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "dynamic",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			log.Println("root")
			log.Println(c.StringArg("dynamic"))
			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "subcommand",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "dynamic",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					log.Println("subcommand")
					log.Println(c.StringArg("dynamic"))
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
