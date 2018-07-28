package couchbase

import (
	"github.com/wheniwork/n1ql-query-builder"
	"gopkg.in/couchbase/gocb.v1"
)

type Repository struct {
	Cb          *DB
	BucketIdent *nqb.Expression
}

func (r Repository) ExecuteN1qlAndLog(q *gocb.N1qlQuery, params interface{}) (gocb.QueryResults, error) {
	rows, err := r.Cb.ExecuteN1qlQuery(q, params)

	if err != nil {
		log.With("error", err).Errorf("Bad N1QL Query: %v, Params: %+v", q, params)
		return rows, err
	}

	log.Debugf("N1QL Query: %v, Params: %+v", q, params)

	return rows, nil
}

func (r Repository) FetchResults(q *gocb.N1qlQuery, params interface{}, valuePtr interface{}, resultHandler func(valuePtr interface{}) error) (gocb.QueryResults, error) {
	res, err := r.ExecuteN1qlAndLog(q, params)

	if err != nil {
		return res, err
	}

	for res.Next(valuePtr) {
		resultHandler(valuePtr)
	}

	if err := res.Close(); err != nil {
		return res, err
	}

	return res, nil
}
