package couchbase

import (
	"github.com/bgetsug/go-toolbox/config"
	"github.com/bgetsug/go-toolbox/errors"
	"github.com/bgetsug/go-toolbox/logging"
	"github.com/spf13/cobra"
)

var (
	seederLog = logging.NewModuleLog(Couchbase, "seeder")
)

var SeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed registered database(s) with fixtures",
	Run: func(cmd *cobra.Command, args []string) {
		Cb.Seed(nil)
	},
}

type SeederResults struct {
	Seeds  []interface{}
	Errors []error
}

func NewSeederResults(seeds []interface{}, errs []error) SeederResults {
	return SeederResults{Seeds: seeds, Errors: errs}
}

type SeedFunc func(cb *DB) SeederResults

// Register one or more seeders
func (c *DB) RegisterSeeders(seeders []SeedFunc) {
	c.seeders = seeders
}

// Flush and seed the Couchbase bucket with all registered seeders
func (c *DB) Seed(seederResults chan<- SeederResults) {
	env := c.config.Environment

	if env != config.DEVELOPMENT && env != config.LOCAL && env != config.TESTING {
		seederLog.Panic("Seeding may only be performed in development or local environments")
	}

	if len(c.seeders) == 0 {
		abortSeederWithError(
			seederResults,
			errors.New("NO_REGISTERED_SEEDERS", "No seeders have been registered"),
		)
		return
	}

	if err := c.Bucket.Manager("", "").Flush(); err != nil {
		abortSeederWithError(seederResults, err)
		return
	}

	if err := c.WaitForHealth(); err != nil {
		abortSeederWithError(seederResults, err)
		return
	}

	for _, seeder := range c.seeders {
		results := seeder(c)

		if seederResults != nil {
			seederResults <- results
		}

		seederLog.With(
			"seeds", len(results.Seeds),
			"errors", len(results.Errors),
		).Info("Seeder completed")
	}

	abortSeederWithError(seederResults, nil)
}

func abortSeederWithError(seederResults chan<- SeederResults, err error) {
	if err != nil {
		seederLog.With("error", err).Error()

		if seederResults != nil {
			seederResults <- SeederResults{Errors: []error{err}}
		}
	}

	if seederResults != nil {
		close(seederResults)
		seederLog.Debug("Closed seeder results channel")
	}
}
