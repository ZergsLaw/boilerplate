package cmd

import (
	"fmt"

	"github.com/go-openapi/loads"
	"github.com/urfave/cli/v2"
	"github.com/zergslaw/boilerplate/internal/api/web/generated/restapi"
)

var (
	Version = &cli.Command{
		Name:         "version",
		Aliases:      []string{"v"},
		Usage:        "prints the service version.",
		UsageText:    "Prints the service version.",
		Action:       versionAction,
		BashComplete: cli.DefaultAppComplete,
	}
)

func versionAction(c *cli.Context) error {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return fmt.Errorf("load embedded swagger spec: %w", err)
	}
	ver := swaggerSpec.Spec().Info.Version

	fmt.Println(c.App.Name, "version", ver)
	return nil
}
