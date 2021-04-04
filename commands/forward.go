package commands

import (
	"fmt"
	"github.com/jtarchie/jsyslog/clients"
	"github.com/jtarchie/jsyslog/listeners"
	"golang.org/x/sync/errgroup"
)

type ForwardCmd struct {
	From []string `help:"the uri to listen to messages" default:"tcp://0.0.0.0:9000" required:""`
	To   []string `help:"the uri to forward messages to" required:""`
}

func (l *ForwardCmd) Run() error {
	outputs := []clients.Client{}
	for _, uri := range l.To {
		output, err := clients.New(uri)
		if err != nil {
			return fmt.Errorf(
				"could not create client (%s): %w",
				uri,
				err,
			)
		}

		outputs = append(outputs, output)
	}

	errGroup := &errgroup.Group{}

	for _, uri := range l.From {
		uri := uri
		errGroup.Go(func() error {
			server, err := listeners.New(uri)
			if err != nil {
				return fmt.Errorf(
					"could not start from (%s): %w",
					uri,
					err,
				)
			}

			return server.ListenAndServe(func(message string) error {
				for _, output := range outputs {
					err := output.WriteString(fmt.Sprintf("%s\n", message))
					if err != nil {
						return fmt.Errorf(
							"could not write to (%s): %w",
							uri,
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
