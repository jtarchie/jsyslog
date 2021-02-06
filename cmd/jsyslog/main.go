package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/jtarchie/jsyslog/listeners"
	"log"
	"net/url"
	"os"
)

type ForwardCmd struct {
	From []*url.URL `help:"the uri to listen to messages" default:"tcp://0.0.0.0:9000" required:""`
	To   []*url.URL `help:"the uri to forward messages to" required:""`
}

func (l *ForwardCmd) Run() error {
	server, err := listeners.New(l.From[0].String())
	if err != nil {
		return fmt.Errorf("could not start from (%s): %w", l.From[0].String(), err)
	}

	file, err := os.OpenFile(l.To[0].Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("could not start to (%s): %w", l.To[0].String(), err)
	}

	log.Printf("starting %s", l.From[0].String())
	return server.ListenAndServe(func(message string) error {
		_, err := file.WriteString(message)
		if err != nil {
			return fmt.Errorf("could not write to (%s): %w", l.To[0].String(), err)
		}

		return nil
	})
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
