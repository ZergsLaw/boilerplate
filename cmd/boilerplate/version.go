package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	version = &cli.Command{
		Name:         "version",
		Aliases:      []string{"v"},
		Usage:        "prints the service version.",
		UsageText:    "Prints the service version.",
		Action:       versionAction,
		BashComplete: cli.DefaultAppComplete,
	}
)

func versionAction(c *cli.Context) error {
	fmt.Println(c.App.Name, "version", ver)
	return nil
}
