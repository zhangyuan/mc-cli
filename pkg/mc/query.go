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

var nilValue = "NULL"

func (client *Client) Sql2csv(dsn string, query string, dataworkVars map[string]interface{}, writer *csv.Writer) error {
	dataChan := make(chan []string, 100)
	errChan := make(chan error, 1)

	go func() {
		if err := client.Query(dsn, query, dataworkVars, func(columnNames []string) error {
			dataChan <- columnNames
			return nil
		}, func(row []any) error {
			values := make([]string, len(row))
			for idx := range row {
				values[idx] = fmt.Sprintf("%v", row[idx])
			}
			dataChan <- values
			return nil
		}); err != nil {
			errChan <- err
		}
		close(errChan)
		close(dataChan)
	}()

Loop:
	for {
		select {
		case data, ok := <-dataChan:
			if !ok {
				break Loop
			}
			if err := writer.Write(data); err != nil {
				return err
			}
		case err, ok := <-errChan:
			if !ok {
				break Loop
			}
			return err
		}
	}

	return nil
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

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	columnCount := len(columns)
	if err := onHeaderFunc(columns); err != nil {
		return err
	}

	for rows.Next() {
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
				csvRow[idx] = nilValue
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
