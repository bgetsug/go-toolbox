package couchbase

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bgetsug/go-toolbox/errors"
	"github.com/bgetsug/go-toolbox/logging"
	"github.com/fatih/structs"
	pkgerr "github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/couchbase/gocb.v1"
)

var (
	indexerLog               = logging.NewModuleLog(couchbase, "indexer")
	errorNoRegisteredIndexes = errors.New("NO_REGISTERED_INDEXES", "No indexes have been registered for creation")
)

var IndexCmd = &cobra.Command{
	Use:   "index [<num replicas>]",
	Short: "Add indexes to registered database(s)",
	Run: func(cmd *cobra.Command, args []string) {
		numReplicas := 1

		if len(args) == 1 {
			nr, err := strconv.Atoi(args[0])

			if err != nil {
				log.Fatal(err)
			}

			numReplicas = nr
		}

		Cb.CreateIndexes(numReplicas)
	},
}

type Index struct {
	Name   string
	Fields []string
	Where  string
}

// Register one or more indexes
func (c *Couchbase) RegisterIndexes(indexes []*Index) {
	c.indexes = indexes
}

// Create all registered indexes that do not already exist
func (c *Couchbase) CreateIndexes(numReplicas int) ([]gocb.IndexInfo, []error) {
	if numReplicas == 0 {
		numReplicas = 1
	}

	hosts := strings.Split(c.config.Hosts, ",")

	var indexErrors []error

	if len(c.indexes) == 0 {
		noIndexes := errorNoRegisteredIndexes
		indexerLog.With(noIndexes)
		indexErrors = append(indexErrors, noIndexes)
	}

	for _, index := range c.indexes {

		err := c.createIndex(index, hosts[0], numReplicas, true)

		if err != nil {
			indexerLog.With("error", pkgerr.WithStack(err)).Error()
			indexErrors = append(indexErrors, err)
		}
	}

	indexes, err := c.Bucket.Manager(Cb.config.BucketName, Cb.config.BucketPassword).GetIndexes()

	if err != nil {
		indexErrors = append(indexErrors, err)
	}

	for _, indexInfo := range indexes {
		indexerLog.With(structs.Map(indexInfo)).Info("Registered index")
	}

	return indexes, indexErrors
}

func (c *Couchbase) createIndex(index *Index, node string, numReplicas int, ignoreIfExists bool) error {
	var qs string

	qs += "CREATE INDEX"

	if index.Name != "" {
		qs += " `" + index.Name + "`"
	}

	qs += " ON `" + c.config.BucketName + "`"

	if len(index.Fields) > 0 {
		qs += " ("
		for i := 0; i < len(index.Fields); i++ {
			if i > 0 {
				qs += ", "
			}
			qs += "`" + index.Fields[i] + "`"
		}
		qs += ")"
	}

	if len(index.Where) > 0 {
		qs += " WHERE " + index.Where
	}

	qs += " USING GSI WITH {\"num_replica\": " + fmt.Sprintf("%d", numReplicas) + "}"

	rows, err := c.Bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery(qs), nil)

	if err != nil {
		if strings.Contains(err.Error(), "already exist") {
			if ignoreIfExists {
				indexerLog.Infof("Index '%s' already exists...skipping creation", index.Name)
				return nil
			}
			return gocb.ErrIndexAlreadyExists
		}

		return err
	}

	if err := rows.Close(); err != nil {
		return err
	}

	return nil
}
