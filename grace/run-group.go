// Package grace provides a way to run a group of functions with interrupt and shutdown handling.
package grace

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
)

type (
	// onExec is a function type that represents a task to be executed.
	onExec func() error
	// onInterrupt is a function type that represents a task to be executed when the RunGroup is interrupted.
	onInterrupt func(error)
)

// NewRunGroup creates and returns a new RunGroup.
func NewRunGroup() *RunGroup {
	return &RunGroup{g: &run.Group{}}
}

// RunGroup is a struct that wraps the oklog/run.Group struct.
type RunGroup struct {
	g *run.Group
}

// Add adds a new task to the RunGroup's task list.
func (o *RunGroup) Add(exec onExec, inter onInterrupt) {
	o.g.Add(exec, inter)
}

// Run starts the RunGroup's tasks, and waits for interrupt signals to be triggered.
// OnShutdown is a callback that will be called when the RunGroup is interrupted.
func (o *RunGroup) Run(onShutdown func(error)) error {
	sigs := make(chan os.Signal, 1)
	allDone := make(chan struct{})
	o.g.Add(func() error {
		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		for {
			select {
			case <-sigs:
				return nil
			case <-allDone:
				return nil
			}
		}
	}, func(err error) {
		onShutdown(err)
		select {
		case allDone <- struct{}{}:
		default:
		}
	})

	//nolint:wrapcheck
	return o.g.Run()
}
