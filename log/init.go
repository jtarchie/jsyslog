package log

import (
	"go.uber.org/zap"
	"log"
)

var Logger *zap.Logger

func init() {
	var err error

	Logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatalf("could not start zap logger: %s", err)
	}
}
