package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	clients "github.com/jtarchie/jsyslog/clients"
	"github.com/jtarchie/jsyslog/listeners"
	"golang.org/x/sync/errgroup"
	"log"
	"net/url"
)

type ForwardCmd struct {
	From []*url.URL `help:"the uri to listen to messages" default:"tcp://0.0.0.0:9000" required:""`
	To   []*url.URL `help:"the uri to forward messages to" required:""`
}

func (l *ForwardCmd) Run() error {
	outputs := []clients.Client{}
	for _, uri := range l.To {
		output, err := clients.New(uri.String())
		if err != nil {
			return fmt.Errorf(
				"could not create client (%s): %w",
				uri.String(),
				err,
			)
		}

		outputs = append(outputs, output)
	}

	errGroup := &errgroup.Group{}

	for _, uri := range l.From {
		uri := uri
		errGroup.Go(func() error {
			server, err := listeners.New(uri.String())
			if err != nil {
				return fmt.Errorf(
					"could not start from (%s): %w",
					uri.String(),
					err,
				)
			}

			return server.ListenAndServe(func(message string) error {
				for _, output := range outputs {
					err := output.WriteString(fmt.Sprintf("%s\n", message))
					if err != nil {
						return fmt.Errorf(
							"could not write to (%s): %w",
							uri.String(),
							err,
						)
					}
				}

				return nil
			})
		})
	}

	return errGroup.Wait()
}

type CLI struct {
	Forwarder ForwardCmd `cmd:"" help:"forward syslog messages"`
}

func main() {
	cli := &CLI{}
	ctx := kong.Parse(cli,
		kong.Name("jsyslog"),
		kong.Description("handle syslog messages"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)

	err := ctx.Run()
	if err != nil {
		log.Fatalf("could not execute: %s", err)
	}
}
