package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"mc-helper/pkg/mc"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query data from MaxCompute",
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		dsn := os.Getenv("DSN")
		if dsn == "" {
			log.Fatalln("DSN is missing")
		}

		vars, err := getFromJsonOrFile(dataworksVarsArg, dataworksVarsPath)
		if err != nil {
			log.Fatalln(err)
		}

		sql, err := getFromStringOrFile(querySQL, querySQLPath)
		if err != nil {
			log.Fatalln(err)
		}

		client, err := mc.NewClient(dsn)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		if format == "csv" {
			writer := csv.NewWriter(os.Stdout)
			defer writer.Flush()
			if err := client.Sql2csv(dsn, sql, vars, writer); err != nil {
				log.Fatalln(err.Error())
			}
		} else if format == "table" {
			t := table.NewWriter()
			if err := client.Sql2Table(dsn, sql, vars, t); err != nil {
				log.Fatalln(err.Error())
			}
			t.SetStyle(table.StyleLight)
			t.SetStyle(table.StyleColoredBright)
			fmt.Println(t.Render())
		} else {

			log.Fatalln(errors.Errorf("Invalid format: %v", format))
		}
	},
}

var querySQL string
var querySQLPath string
var format string

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringVarP(&querySQL, "sql", "s", "", "SQL query")
	queryCmd.Flags().StringVarP(&querySQLPath, "sql-file", "f", "", "Path to sql query")
	queryCmd.MarkFlagsOneRequired("sql", "sql-file")

	queryCmd.Flags().StringVarP(&dataworksVarsArg, "dataworks-vars", "v", "", "Variables in json")
	queryCmd.Flags().StringVarP(&dataworksVarsPath, "dataworks-vars-file", "d", "", "Path to variables file in YAML")

	queryCmd.Flags().StringVar(&format, "format", "table", "Format: csv, table")
}
