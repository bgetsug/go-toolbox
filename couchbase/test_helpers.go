package couchbase

import (
	"testing"
	"time"

	"github.com/bgetsug/go-toolbox/test_helpers"
	"gopkg.in/couchbase/gocb.v1"
)

type BucketFlush struct {
	watchIndexTimeout time.Duration
}

func WithBucketFlush(watchIndexTimeout time.Duration) *BucketFlush {
	if watchIndexTimeout == 0 {
		watchIndexTimeout = time.Second * 60
	}

	return &BucketFlush{watchIndexTimeout: watchIndexTimeout}
}

func (d *BucketFlush) Before(ctx *test_helpers.Context) {
	if err := Cb.Bucket.Manager(Cb.config.BucketName, Cb.config.BucketPassword).Flush(); err != nil {
		ctx.MustGet("t").(*testing.T).Fatal(err)
	}

	if err := Cb.WaitForHealth(); err != nil {
		ctx.MustGet("t").(*testing.T).Fatal(err)
	}

	d.createIndexes(ctx)
}

func (d *BucketFlush) createIndexes(ctx *test_helpers.Context) {
	bucketMgr := Cb.Bucket.Manager(Cb.config.BucketName, Cb.config.BucketPassword)

	indexes, errs := Cb.CreateIndexes()

	if len(errs) > 0 {
		ctx.MustGet("t").(*testing.T).Fatal(errs)
	}

	var watchList []string

	for _, index := range indexes {
		watchList = append(watchList, index.Name)
	}

	if err := bucketMgr.WatchIndexes(watchList, false, d.watchIndexTimeout); err != nil {
		ctx.MustGet("t").(*testing.T).Fatal(err)
	}
}

func (d *BucketFlush) After(ctx *test_helpers.Context) {}

type QueryConsistency struct{}

func WithQueryConsistency() *QueryConsistency {
	return &QueryConsistency{}
}

func (d *QueryConsistency) Before(ctx *test_helpers.Context) {
	Cb.SetConsistencyMode(gocb.RequestPlus)
}

func (d *QueryConsistency) After(ctx *test_helpers.Context) {
	Cb.ResetConsistencyMode()
}
