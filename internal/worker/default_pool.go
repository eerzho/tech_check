package worker

import (
	"context"
	"log/slog"
)

var defaultPool *Pool

func Start(lg *slog.Logger, count int) {
	defaultPool = NewPool(lg, count)
	defaultPool.Start(context.Background())
}

func Stop() {
	defaultPool.Stop()
}

func AddTask(task Task) {
	defaultPool.AddTask(task)
}
