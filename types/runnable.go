package gdtcore

import (
	"context"
	"testing"
)

// Runnable represents things that have a Run() method that accepts a Context
// and a pointer to a testing.T. Example things that implement this interface
// are `gdtcore.scenario.Scenario` and `gdtcore.suite.Suite`.
type Runnable interface {
	Run(context.Context, *testing.T)
}
