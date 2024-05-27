package cmd

import (
	"log"
	"os"

	"mc-helper/pkg/mc"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var query string

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "dump data from MC",
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		dsn := os.Getenv("DSN")
		if dsn == "" {
			log.Fatalln("DSN is missing")
		}
		if err := mc.Sql2csv(dsn, query); err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
	dumpCmd.Flags().StringVarP(&query, "query", "q", "", "sql query")
	_ = dumpCmd.MarkFlagRequired("query")
}
