package mc

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

type CsvRow struct {
	Error   error
	Content []string
}

func Sql2csv(dsn string, query string) error {
	db, err := NewDB(dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if !strings.HasSuffix(strings.TrimSpace(query), ";") {
		query = query + ";"
	}

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	data := make(chan CsvRow, 1000)

	csvWriter := csv.NewWriter(os.Stdout)
	defer csvWriter.Flush()

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	go func() {
		isHeader := true
		columnCount := 0
		for rows.Next() {
			if isHeader {
				columns, err := rows.Columns()
				if err != nil {
					data <- CsvRow{Error: err}
					close(data)
					return
				}

				columnCount = len(columns)

				data <- CsvRow{Content: columns}
				isHeader = false
			}

			var record = make([]interface{}, columnCount)
			var recordPointer = make([]any, columnCount)
			for idx := range record {
				recordPointer[idx] = &record[idx]
			}

			if err := rows.Scan(recordPointer...); err != nil {
				data <- CsvRow{Error: err}
				close(data)
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

			data <- CsvRow{Content: csvRow}
		}

		close(data)
	}()

	for {
		if row, ok := <-data; ok {
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
