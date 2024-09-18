package worker

import (
	"context"
	"log/slog"
)

type Pool struct {
	lg         *slog.Logger
	count      int
	tasks      chan Task
	cancelFunc context.CancelFunc
}

func NewPool(lg *slog.Logger, count int) *Pool {
	return &Pool{
		lg:    lg,
		count: count,
		tasks: make(chan Task, count*5),
	}
}

func (p *Pool) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	p.cancelFunc = cancel

	for i := 0; i < p.count; i++ {
		go p.worker(ctx)
	}
}

func (p *Pool) Stop() {
	if p.cancelFunc != nil {
		p.cancelFunc()
	}
}

func (p *Pool) AddTask(task Task) {
	p.tasks <- task
}

func (p *Pool) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-p.tasks:
			err := task.Execute(ctx)
			if err != nil {
				p.lg.Error(err.Error())
			}
		}
	}
}
