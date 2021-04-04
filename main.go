package main

import (
	"github.com/alecthomas/kong"
	"github.com/jtarchie/jsyslog/commands"
	"github.com/jtarchie/jsyslog/log"
	"go.uber.org/zap"
)

type CLI struct {
	Forwarder commands.ForwardCmd `cmd:"" help:"forward syslog messages"`
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
		log.Logger.Fatal(
			"could not execute",
			zap.String("error", err.Error()),
		)
	}
}
