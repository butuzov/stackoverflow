package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
)

type Command struct {
	cli *cli.App
}

func New() *Command {
	return &Command{
		cli: bootstrap(os.Stdout),
	}
}

func (cmd *Command) Run(args []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	return cmd.cli.RunContext(ctx, args)
}

func bootstrap(w io.Writer) *cli.App {
	cli.HelpFlag = &cli.BoolFlag{
		Hidden: true,
		Name:   "help",
		Usage:  "help",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "print-version",
		Aliases: []string{"V"},
		Usage:   "print only the version",
	}

	// actual application.
	app := &cli.App{
		Writer:    os.Stdout,
		Name:      "stackoverflow",
		Usage:     "stackoverflow tags monitor",
		UsageText: "stackoverflow -o numpy ",
		// ~~~~~~~~~ Common ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

		// ~~~~~~~~~ Functionality Related ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		Action: func(ctx *cli.Context) error {
			tags, err := ParseTags(ctx.Args().Slice())
			if err != nil {
				return err
			}

			var (
				config = Config{
					Host: Host(ctx.String("host")),
					Tags: tags,
					Open: ctx.Bool("open"),
				}

				client = &http.Client{
					Transport: &http.Transport{
						MaxIdleConns:       3,
						IdleConnTimeout:    2 * time.Second,
						DisableCompression: false,
					},
				}
			)

			monitor := NewMonitor(config, client)

			return monitor.Run(ctx.Context)
		},
		Commands: []*cli.Command{},

		// ~~~~~~~~ Helper ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

		CommandNotFound: func(c *cli.Context, cmd string) {
			_, _ = fmt.Fprintf(c.App.Writer, "[%s] not found\n", cmd)
		},

		// ~~~~~~~~ Flags ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		UseShortOptionHandling: true,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "host",
				Aliases: []string{"h"},
				Value:   "stackoverflow.com",
				Usage:   "Stackoverflow host",
			},
			&cli.BoolFlag{
				Name:    "open",
				Aliases: []string{"o"},
				Usage:   "Open new posts in browser",
				Value:   false,
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}
