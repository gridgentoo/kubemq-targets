package bigtable

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/types"
)

type metadata struct {
	tableName      string
	columnFamily   string
	method         string
	rowKeyPrefix   string
	readColumnName string
}

var methodsMap = map[string]string{
	"write":                  "write",
	"write_batch":            "write_batch",
	"get_row":                "get_row",
	"get_all_rows":           "get_all_rows",
	"delete_row":             "delete_row",
	"get_tables":             "get_tables",
	"create_table":           "create_table",
	"delete_table":           "delete_table",
	"create_column_family":   "create_column_family",
	"get_all_rows_by_column": "get_all_rows_by_column",
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}
	if m.method != "get_tables" {
		m.tableName, err = meta.MustParseString("table_name")
		if err != nil {
			return metadata{}, fmt.Errorf("table_name is required for method:%s , error parsing table_name, %w", m.method, err)
		}
	}
	if m.method == "write" || m.method == "write_batch" || m.method == "create_column_family" {
		m.columnFamily, err = meta.MustParseString("column_family")
		if err != nil {
			return metadata{}, fmt.Errorf("column_family is required for method:%s ,error parsing column_family, %w", m.method, err)
		}
	} else if m.method == "delete_row" || m.method == "get_row" {
		m.rowKeyPrefix, err = meta.MustParseString("row_key_prefix")
		if err != nil {
			return metadata{}, fmt.Errorf("row_key_prefix is required for method:%s ,error parsing row_key_prefix, %w", m.method, err)
		}
	}
	if m.method == "get_all_rows_by_column" {
		m.readColumnName, err = meta.MustParseString("column_name")
		if err != nil {
			return metadata{}, fmt.Errorf("column_name is required for method:%s , error parsing column_name, %w", m.method, err)
		}
	}

	return m, nil
}
