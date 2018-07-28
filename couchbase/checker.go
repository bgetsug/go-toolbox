package couchbase

import (
	"github.com/dimiro1/health"

	"gopkg.in/couchbase/gocb.v1"
)

var CbChecker *Checker

func InitChecker(cb *DB) {
	CbChecker = &Checker{cb}
}

type Checker struct {
	cb *DB
}

func (c *Checker) Check() health.Health {
	status := health.NewHealth()

	_, err := c.cb.Bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery("SELECT 1"), nil)

	if err != nil {
		status.Down()
	} else {
		status.Up()
	}

	return status
}
