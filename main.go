package main

import (
	"github.com/alecthomas/kong"
	"github.com/jtarchie/jsyslog/commands"
	"go.uber.org/zap"
	"log"
)

type CLI struct {
	Forwarder commands.ForwardCmd `cmd:"" help:"forward syslog messages"`
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("could not start logger: %s", err)
	}

	cli := &CLI{}
	ctx := kong.Parse(cli,
		kong.Name("jsyslog"),
		kong.Description("handle syslog messages"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)

	err = ctx.Run(logger)
	if err != nil {
		logger.Fatal(
			"could not execute",
			zap.String("error", err.Error()),
		)
	}
}
