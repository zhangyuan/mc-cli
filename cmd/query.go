package cmd

import (
	"log"
	"os"

	"mc-helper/pkg/mc"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
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

		if err := client.Sql2csv(dsn, sql, vars); err != nil {
			log.Fatalln(err.Error())
		}
	},
}

var querySQL string
var querySQLPath string

func init() {
	rootCmd.AddCommand(dumpCmd)
	dumpCmd.Flags().StringVarP(&querySQL, "sql", "s", "", "SQL query")
	dumpCmd.Flags().StringVarP(&querySQLPath, "sql-path", "p", "", "Path to sql query")
	dumpCmd.MarkFlagsOneRequired("sql", "sql-path")

	dumpCmd.Flags().StringVarP(&dataworksVarsArg, "dataworks-vars", "v", "", "Variables in json")
	dumpCmd.Flags().StringVarP(&dataworksVarsPath, "dataworks-vars-path", "d", "", "Path to variables file in YAML")
}
