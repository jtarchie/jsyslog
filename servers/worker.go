package servers

import (
	"go.uber.org/zap"
	"sync"
)

type Worker struct {
	errorFunc func(error)
	size      uint
	started   bool
	wait      *sync.WaitGroup
	workers   chan func(int) error
	logger    *zap.Logger
}

type WorkerOption func(*Worker)

func WithErrorFunc(errFunc func(error)) func(*Worker) {
	return func(w *Worker) {
		w.errorFunc = errFunc
	}
}

func NewWorker(
	size uint,
	logger *zap.Logger,
	options ...WorkerOption,
) *Worker {
	worker := &Worker{
		size:      size,
		started:   false,
		errorFunc: func(error) {},
		logger:    logger,
	}

	for _, option := range options {
		option(worker)
	}

	return worker
}

func (w *Worker) Start() {
	if w.started {
		return
	}

	workers := make(chan func(int) error, w.size)
	wait := &sync.WaitGroup{}

	for i := 0; i < int(w.size); i++ {
		wait.Add(1)
		go perform(i, workers, w.errorFunc, wait, w.logger)
	}

	w.started = true
	w.wait = wait
	w.workers = workers
}

func (w *Worker) Stop() {
	if w.started {
		w.logger.Info("stopping workers")
		close(w.workers)
		w.wait.Wait()
	}
	w.started = false
}

func (w *Worker) Run(fn func(int) error) {
	w.workers <- fn
}

func perform(
	id int,
	workers <-chan func(int) error,
	errorFunc func(error),
	wait *sync.WaitGroup,
	logger *zap.Logger,
) {
	defer wait.Done()

	logger.Info(
		"starting worker",
		zap.Int("workerID", id),
	)

	for work := range workers {
		if err := work(id); err != nil {
			errorFunc(err)
		}
	}

	logger.Info(
		"stopping worker",
		zap.Int("workerID", id),
	)
}
