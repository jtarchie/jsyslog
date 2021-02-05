package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/jtarchie/jsyslog/listeners"
	"log"
	"net/url"
	"os"
)

type CLI struct {
	Listen *url.URL `help:"the uri to listen on" default:"tcp://0.0.0.0:9000" required:""`
	File   string   `type:"path" help:"where to forward to" required:""`
}

func main() {
	cli := &CLI{}
	kong.Parse(cli,
		kong.Name("forwarder"),
		kong.Description("forward syslog messages to another server"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)

	err := execute(cli)
	if err != nil {
		log.Fatalf("could not execute: %s", err)
	}
}

func execute(cli *CLI) error {
	server, err := listeners.New(cli.Listen.String())
	if err != nil {
		return fmt.Errorf("could not create server: %w", err)
	}

	file, err := os.OpenFile(cli.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("could not open file for write: %w", err)
	}

	log.Printf("starting %s", cli.Listen.String())
	return server.ListenAndServe(func(message string) error {
		_, err := file.WriteString(message)
		return fmt.Errorf("could not write to file: %w", err)
	})
}
