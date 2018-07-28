package couchbase

import (
	"fmt"
	"strings"

	"github.com/bgetsug/go-toolbox/errors"
	"github.com/bgetsug/go-toolbox/logging"
	"github.com/fatih/structs"
	pkgerr "github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/couchbase/gocb.v1"
)

var (
	indexerLog               = logging.NewModuleLog(Couchbase, "indexer")
	errorNoRegisteredIndexes = errors.New("NO_REGISTERED_INDEXES", "No indexes have been registered for creation")
)

var IndexCmd = &cobra.Command{
	Use:   "index",
	Short: "Add indexes to registered database(s)",
	Run: func(cmd *cobra.Command, args []string) {
		Cb.CreateIndexes()
	},
}

type Index struct {
	Name   string
	Fields []string
	Where  string
}

// Register one or more indexes
func (c *DB) RegisterIndexes(indexes []*Index) {
	c.indexes = indexes
}

// Create all registered indexes that do not already exist
func (c *DB) CreateIndexes() ([]gocb.IndexInfo, []error) {
	indexErrors := c.createIndexesOnAllHosts()

	indexes, err := c.Bucket.Manager("", "").GetIndexes()

	if err != nil {
		indexErrors = append(indexErrors, err)
	}

	for _, indexInfo := range indexes {
		indexerLog.With(structs.Map(indexInfo)).Info("Registered index")
	}

	return indexes, indexErrors
}

func (c *DB) createIndexesOnAllHosts() []error {
	hosts := strings.Split(c.config.Hosts, ",")

	var indexErrors []error

	if len(c.indexes) == 0 {
		noIndexes := errorNoRegisteredIndexes
		indexerLog.With(noIndexes)
		indexErrors = append(indexErrors, noIndexes)
	}

	for _, index := range c.indexes {
		for _, host := range hosts {
			host = host + ":8091"

			err := c.createIndex(index, host, true)

			if err != nil {
				indexerLog.With("error", pkgerr.WithStack(err)).Error()
				indexErrors = append(indexErrors, err)
			}
		}
	}

	return indexErrors
}

func (c *DB) createIndex(index *Index, node string, ignoreIfExists bool) error {
	nodeParts := strings.Split(strings.Split(node, ":")[0], ".")
	name := fmt.Sprintf("%s_%s", index.Name, nodeParts[0])

	var qs string

	qs += "CREATE INDEX"

	if name != "" {
		qs += " `" + name + "`"
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

	qs += " USING GSI WITH {\"nodes\":[\"" + node + "\"]}"

	rows, err := c.Bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery(qs), nil)

	if err != nil {
		if strings.Contains(err.Error(), "already exist") {
			if ignoreIfExists {
				indexerLog.Infof("Index '%s' already exists...skipping creation", name)
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
