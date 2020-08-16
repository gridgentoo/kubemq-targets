package athena

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type metadata struct {
	method string

	query          string
	catalog        string
	outputLocation string
	DB             string
	executionId    string
}

var methodsMap = map[string]string{
	"list_databases":     "list_databases",
	"list_data_catalogs": "list_data_catalogs",
	"query":              "query",
	"get_query_result":   "get_query_result",
}

func getValidMethodTypes() string {
	s := "invalid method type, method type should be one of the following:"
	for k := range methodsMap {
		s = fmt.Sprintf("%s :%s,", s, k)
	}
	return s
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf(getValidMethodTypes())
	}
	if m.method != "list_data_catalogs" {
		m.catalog, err = meta.MustParseString("catalog")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing catalog, %w", err)
		}
		if m.method == "query" {
			m.query, err = meta.MustParseString("query")
			if err != nil {
				return metadata{}, fmt.Errorf("error parsing query, %w", err)
			}
			m.DB, err = meta.MustParseString("db")
			if err != nil {
				return metadata{}, fmt.Errorf("error parsing db, %w", err)
			}
			m.DB, err = meta.MustParseString("output_location")
			if err != nil {
				return metadata{}, fmt.Errorf("error parsing output_location, %w", err)
			}
		} else if m.method == "get_query_result" {
			m.executionId, err = meta.MustParseString("execution_id")
			if err != nil {
				return metadata{}, fmt.Errorf("error parsing execution_id, %w", err)
			}
		}
	}
	return m, nil
}
