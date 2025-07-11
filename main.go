package main

import (
	"fmt"
	orson "orson/src"
	"orson/src/version"
	"os"
	"sort"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/urfave/cli/v2"
	"moul.io/banner"
)

func main() {
	fmt.Println(banner.Inline("orson"))
	fmt.Println("version:", version.Version)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var (
		directory string
	)

	app := &cli.App{
		EnableBashCompletion: true,
		Flags:                []cli.Flag{},
		UsageText:            "Orson is a CLI for investigating Model Context Protocol content",
		Commands: []*cli.Command{
			{
				Name:    "scan",
				Aliases: []string{"m"},
				Usage:   "scan files in folders for package managers and examine and report for MCP",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "directory",
						Aliases:     []string{"d"},
						Usage:       "Directory to scan (defaults to .)",
						Value:       ".",
						Destination: &directory,
					},
				},
				Action: func(*cli.Context) error {
					//arn, err :=
					orson.GetViolations(directory)
					//if err != nil {
					//	return fmt.Errorf("make failed: %w", err)
					//}

					//if arn != nil {
					//	fmt.Print(*arn)
					//}

					return nil
				},
			},
			{
				Name:      "version",
				Aliases:   []string{"v"},
				Usage:     "Outputs the application version",
				UsageText: "orson version",
				Action: func(*cli.Context) error {
					fmt.Println(version.Version)

					return nil
				},
			},
		},
		Name:     "orson",
		Usage:    "Examine codebase for MCP",
		Compiled: time.Time{},
		Authors:  []*cli.Author{{Name: "James Woolfenden", Email: "james.woolfenden@gmail.com"}},
		Version:  version.Version,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err)
	}
}
