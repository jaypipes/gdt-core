package run

import (
	"path/filepath"
	"slices"
	"time"

	"github.com/samber/lo"

	"github.com/gdt-dev/core/api"
	"github.com/gdt-dev/core/testunit"
)

// Run stores state of a test run when tests are executed with the `gdt` CLI
// tool.
type Run struct {
	// scenarioResults is a map, keyed by the Scenario path, of slices of
	// testunit.Result structs corresponding to the test specs in the scenario.
	// There is guaranteed to be exactly the same number of testunit.Results in
	// the slice as scenarios in the scenario.
	scenarioResults map[string][]testunit.Result
}

// OK returns true if all Scenarios in the Run had all successful test units.
func (r *Run) OK() bool {
	return lo.EveryBy(
		lo.Values(r.scenarioResults),
		func(results []testunit.Result) bool {
			return lo.EveryBy(results, func(tur testunit.Result) bool {
				return tur.OK()
			})
		},
	)
}

// ScenarioPaths returns a sorted list of Scenario Paths.
func (r *Run) ScenarioPaths() []string {
	paths := lo.Keys(r.scenarioResults)
	slices.Sort(paths)
	return paths
}

// ScenarioResults returns the set of testunit.Results for a Scenario with the
// supplied path.
func (r *Run) ScenarioResults(path string) []testunit.Result {
	return r.scenarioResults[path]
}

// StoreResult stores a test unit result to the Run for the supplied test unit.
func (r *Run) StoreResult(
	index int,
	path string, // the Scenario.Path
	tu *testunit.TestUnit,
	res *api.Result,
) {
	if _, ok := r.scenarioResults[path]; !ok {
		r.scenarioResults[path] = []testunit.Result{}
	}
	r.scenarioResults[path] = append(
		r.scenarioResults[path],
		testunit.NewResult(index, tu, res),
	)
}

// XUnit returns the Run's scenario results as a slice of structs that can be
// serialized to either XML (JUnit/XUnit-style) or JSON.
func (r *Run) XUnit() []api.XUnitTestSuite {
	suites := []api.XUnitTestSuite{}
	paths := r.ScenarioPaths()
	for _, path := range paths {
		shortPath := filepath.Base(path)
		suite := api.XUnitTestSuite{
			Name: shortPath,
			Properties: []api.XUnitProperty{
				{
					Name:  "path",
					Value: path,
				},
			},
			Timestamp: time.Now(),
		}

		var scenElapsed time.Duration

		unitResults := r.ScenarioResults(path)

		testcases := make([]api.XUnitTestCase, len(unitResults))
		tcFails := 0

		for x, res := range unitResults {
			tc := res.XUnit()
			if res.Failed() {
				tcFails++
			}
			scenElapsed += res.Elapsed()
			testcases[x] = tc
		}
		suite.Failures = tcFails
		suite.Tests = len(testcases)
		suite.Time = scenElapsed.String()
		suite.TestCases = testcases
		suites = append(suites, suite)
	}
	return suites
}
