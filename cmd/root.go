package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:   "mc-helper",
	Short: "MC Helper",
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var dataworksVarsArg string
var dataworksVarsPath string

func getFromJsonOrFile(varsArg string, varsPath string) (map[string]interface{}, error) {
	vars := map[string]interface{}{}

	if len(varsArg) > 0 {
		err := json.Unmarshal([]byte(varsArg), &vars)
		if err != nil {
			return nil, errors.Errorf("Invalid vars: %v", err)
		}
	} else if len(varsPath) > 0 {
		varsBytes, err := os.ReadFile(varsPath)
		if err != nil {
			return nil, errors.Errorf("Invalid vars file: %v", err)
		}
		if err := yaml.Unmarshal(varsBytes, &vars); err != nil {
			log.Fatalln(errors.Errorf("Invalid vars file: %v", err))
		}
	}

	return vars, nil
}

func getFromStringOrFile(content string, filePath string) (string, error) {
	if content == "" {
		bytes, err := os.ReadFile(filePath)
		if err != nil {
			return "", errors.Errorf("Invalid SQL file: %v", err)
		}
		content = string(bytes)
	}
	return content, nil
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
