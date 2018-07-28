package couchbase

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bgetsug/go-toolbox/config"
	"github.com/bgetsug/go-toolbox/logging"
	"github.com/pkg/errors"
	"gopkg.in/couchbase/gocb.v1"
	"gopkg.in/couchbase/gocbcore.v7"
)

const Couchbase = "couchbase"

var (
	log = logging.NewModuleLog(Couchbase)
	Cb  *DB
)

type Configuration struct {
	Environment       config.Environment
	MaxConnectRetries int
	Hosts             string
	BucketName        string
	BucketPassword    string
}

// Init sets up a new Couchbase cluster and bucket connection, and health check.
// If logger is not nil, the Debug and DebugVerbose fields of the Configuration will be ignored.
func Init(configuration Configuration, logger gocbcore.Logger) {
	Cb = &DB{config: configuration, logger: logger, consistencyMode: gocb.NotBounded}
	Cb.Connect()
	InitChecker(Cb)
}

// DB store information about the DB
type DB struct {
	*gocb.Bucket
	config          Configuration
	cluster         *gocb.Cluster
	logger          gocbcore.Logger
	connectRetries  int
	indexes         []*Index
	seeders         []SeedFunc
	consistencyMode gocb.ConsistencyMode
}

// Connect initializes a bucket connection
func (c *DB) Connect() {
	if c.logger != nil {
		gocbcore.SetLogger(c.logger)
	}

	maxConnectRetries := c.config.MaxConnectRetries
	log.Info("Database maximum connection retries: " + strconv.Itoa(maxConnectRetries))

	connectString := "couchbase://" + c.config.Hosts
	log.Info("Connecting to: " + connectString)

	cluster, clusterError := gocb.Connect(connectString)

	if clusterError != nil {
		log.With("error", errors.WithStack(clusterError)).Fatal()
	}

	c.cluster = cluster

	bucketName := c.config.BucketName
	bucketPassword := c.config.BucketPassword

	bucket, bucketError := cluster.OpenBucket(bucketName, bucketPassword)

	if bucketError != nil {
		log.With("error", errors.WithStack(bucketError)).Error()
		log.Info("Retrying connecting to database...")

		if c.connectRetries < maxConnectRetries {
			c.connectRetries++
			time.Sleep(2 * time.Second)
			c.Connect()
		}

		if c.connectRetries >= maxConnectRetries {
			log.With("error", errors.WithStack(bucketError)).Panic("Max DB connection retries reached")
		}
	}

	c.Bucket = bucket
}

func (c *DB) ConsistencyMode() gocb.ConsistencyMode {
	return c.consistencyMode
}

func (c *DB) SetConsistencyMode(mode gocb.ConsistencyMode) {
	c.consistencyMode = mode
}

func (c *DB) ResetConsistencyMode() {
	c.consistencyMode = gocb.NotBounded
}

func (c *DB) WaitForHealth() error {
	hosts := strings.Split(c.config.Hosts, ",")

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://"+hosts[0]+":8091/pools/default/buckets/"+c.config.BucketName, nil)

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
				seederLog.With("node", node["hostname"]).Warn("Cluster not yet healthy. Waiting...")
				time.Sleep(2 * time.Second)
				continue GetStatus
			}
		}

		break
	}

	return nil
}
