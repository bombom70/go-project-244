package main

import (
	"code"
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	if err := (&cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "stylish",
				Usage:   "output format",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			filePath1 := cmd.Args().Get(0)
			filePath2 := cmd.Args().Get(1)
			format := cmd.String("format")

			result, err := code.GenDiff(filePath1, filePath2, format)
			if err != nil {
				return err
			}

			fmt.Println(result)

			return nil
		},
		Name: "gendiff - Compares two configuration files and shows a difference.",
	}).Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
