package worker

import (
	"context"
	"log/slog"
)

var defaultPool *Pool

func SetupDefaultPool(lg *slog.Logger, count int) {
	defaultPool = NewPool(lg, count)
	defaultPool.Start(context.Background())
}

func StopDefaultPool() {
	if defaultPool != nil {
		defaultPool.Stop()
	}
}

func AddTask(task Task) {
	defaultPool.AddTask(task)
}
