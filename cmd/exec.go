package cmd

import (
	"log"
	"mc-helper/pkg/mc"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var executeCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute SQL",
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

		sql, err := getFromStringOrFile(execSQL, execSQLPath)
		if err != nil {
			log.Fatalln(err)
		}

		client, err := mc.NewClient(dsn)
		if err != nil {
			log.Fatal(err)
		}

		defer client.Close()

		if err := client.Execute(sql, vars); err != nil {
			log.Fatalln(err.Error())
		}
	},
}

var execSQL string
var execSQLPath string

func init() {
	rootCmd.AddCommand(executeCmd)

	executeCmd.Flags().StringVarP(&execSQL, "sql", "s", "", "SQL statement(s)")
	executeCmd.Flags().StringVarP(&execSQLPath, "sql-file", "f", "", "Path to sql statement(s)")
	executeCmd.MarkFlagsOneRequired("sql", "sql-file")

	executeCmd.Flags().StringVarP(&dataworksVarsArg, "dataworks-vars", "v", "", "Variables in json")
	executeCmd.Flags().StringVarP(&dataworksVarsPath, "dataworks-vars-path", "d", "", "Path to variables file in YAML")
}
