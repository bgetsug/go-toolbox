package couchbase

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
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
	if err := Cb.Bucket.Manager("", "").Flush(); err != nil {
		ctx.MustGet("t").(*testing.T).Fatal(err)
	}

	if err := d.WaitForNodeHealth(); err != nil {
		ctx.MustGet("t").(*testing.T).Fatal(err)
	}

	if err := Cb.WaitForHealth(); err != nil {
		ctx.MustGet("t").(*testing.T).Fatal(err)
	}

	d.createIndexes(ctx)
}

func (d *BucketFlush) WaitForNodeHealth() error {
	hosts := strings.Split(Cb.config.Hosts, ",")

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://"+hosts[0]+":8091/pools/default", nil)
	req.SetBasicAuth(Cb.config.BucketName, Cb.config.BucketPassword)

	if err != nil {
		return err
	}

	type Bucket struct {
		Nodes []map[string]interface{} `json:"nodes"`
	}

GetStatus:
	for {
		resp, err := client.Do(req)

		if err != nil {
			return err
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		var bucket Bucket

		if err := json.Unmarshal(body, &bucket); err != nil {
			return err
		}

		for _, node := range bucket.Nodes {
			if node["status"] != "healthy" {
				log.With("node", node["hostname"]).Warn("Cluster not yet healthy. Waiting...")
				time.Sleep(2 * time.Second)
				continue GetStatus
			}
		}

		break
	}

	return nil
}

func (d *BucketFlush) createIndexes(ctx *test_helpers.Context) {
	bucketMgr := Cb.Bucket.Manager("", "")

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
