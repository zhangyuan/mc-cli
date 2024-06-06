package mc

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Row struct {
	Error   error
	Content []string
}

func (client *Client) Sql2csv(dsn string, query string, dataworkVars map[string]interface{}, writer *csv.Writer) error {
	return client.Query(dsn, query, dataworkVars, func(columnNames []string) error {
		return writer.Write(columnNames)
	}, func(row []any) error {
		values := make([]string, len(row))
		for idx := range row {
			values[idx] = fmt.Sprintf("%v", row[idx])
		}
		return writer.Write(values)
	})
}

func (client *Client) Sql2Table(dsn string, query string, dataworkVars map[string]interface{}, writer table.Writer) error {
	return client.Query(dsn, query, dataworkVars, func(columnNames []string) error {
		row := make([]interface{}, len(columnNames))
		for idx := range columnNames {
			row[idx] = columnNames[idx]
		}
		writer.AppendHeader(row)
		return nil
	}, func(row []any) error {
		writer.AppendRow(row)
		return nil
	})
}

func (client *Client) Query(dsn string,
	query string,
	dataworkVars map[string]interface{},
	onHeaderFunc func(columNames []string) error,
	onRowFunc func(values []any) error,
) error {
	if !strings.HasSuffix(strings.TrimSpace(query), ";") {
		query = query + ";"
	}

	sql, err := CompileTemplate(query, dataworkVars)
	if err != nil {
		return err
	}

	rows, err := client.DB.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	isHeader := true
	columnCount := 0
	for rows.Next() {
		if isHeader {
			columns, err := rows.Columns()
			if err != nil {
				return err
			}

			columnCount = len(columns)
			if err := onHeaderFunc(columns); err != nil {
				return err
			}
			isHeader = false
		}

		var record = make([]interface{}, columnCount)
		var recordPointer = make([]any, columnCount)
		for idx := range record {
			recordPointer[idx] = &record[idx]
		}

		if err := rows.Scan(recordPointer...); err != nil {
			return err
		}

		csvRow := make([]any, columnCount)

		for idx := range record {
			columValue := record[idx]
			if columValue == nil {
				csvRow[idx] = ""
				continue
			}

			dataType := columnTypes[idx]

			if val, ok := columValue.(time.Time); ok {
				if dataType.DatabaseTypeName() == "DATE" {
					csvRow[idx] = val.Format("2006-01-02")
				} else {
					csvRow[idx] = val.Format(time.RFC3339)
				}
				continue
			}
			csvRow[idx] = record[idx]
		}

		if err := onRowFunc(csvRow); err != nil {
			return err
		}
	}
	return nil
}
