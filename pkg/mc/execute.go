package mc

import (
	"strings"
)

func (client *Client) Execute(sql string, vars map[string]interface{}) error {
	if !strings.HasSuffix(strings.TrimSpace(sql), ";") {
		sql = sql + ";"
	}

	sql, err := CompileTemplate(sql, vars)
	if err != nil {
		return err
	}

	_, err = client.DB.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}
