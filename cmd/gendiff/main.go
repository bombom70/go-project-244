package main

import (
	"code"
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	fmt.Println(code.GenDiff("name", "name", "name"))

	if err := (&cli.Command{
		// Flags: []cli.Flag{},
		// Action: func(ctx context.Context, cmd *cli.Command) error {

		// },
		Name: "gendiff - Compares two configuration files and shows a difference.",
		// Usage: "",
	}).Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
