package mc

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"
)

type Record struct {
	Error   error
	Content []string
}

func (client *Client) Sql2csv(dsn string, query string, dataworkVars map[string]interface{}, csvWriter *csv.Writer) error {
	rowChan := make(chan Record, 100)
	go client.Query(dsn, query, dataworkVars, rowChan)

	for {
		if row, ok := <-rowChan; ok {
			if row.Error != nil {
				return row.Error
			}

			if err := csvWriter.Write(row.Content); err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}

func (client *Client) Query(dsn string, query string, dataworkVars map[string]interface{}, dataChan chan Record) {
	if !strings.HasSuffix(strings.TrimSpace(query), ";") {
		query = query + ";"
	}

	sql, err := CompileTemplate(query, dataworkVars)
	if err != nil {
		dataChan <- Record{Error: err}
		close(dataChan)
		return
	}

	rows, err := client.DB.Query(sql)
	if err != nil {
		dataChan <- Record{Error: err}
		close(dataChan)
		return
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		dataChan <- Record{Error: err}
		close(dataChan)
		return
	}

	isHeader := true
	columnCount := 0
	for rows.Next() {
		if isHeader {
			columns, err := rows.Columns()
			if err != nil {
				dataChan <- Record{Error: err}
				close(dataChan)
				return
			}

			columnCount = len(columns)

			dataChan <- Record{Content: columns}
			isHeader = false
		}

		var record = make([]interface{}, columnCount)
		var recordPointer = make([]any, columnCount)
		for idx := range record {
			recordPointer[idx] = &record[idx]
		}

		if err := rows.Scan(recordPointer...); err != nil {
			dataChan <- Record{Error: err}
			close(dataChan)
			return
		}

		csvRow := make([]string, columnCount)

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

			csvRow[idx] = fmt.Sprintf("%v", columValue)
		}

		dataChan <- Record{Content: csvRow}
	}

	close(dataChan)
}
