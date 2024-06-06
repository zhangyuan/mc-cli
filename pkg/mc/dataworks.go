package mc

import (
	"fmt"
	"regexp"
)

var re = regexp.MustCompile(`\${([\w.]+)}`)

func CompileTemplate(template string, vars map[string]interface{}) (string, error) {
	compiled := re.ReplaceAllStringFunc(template, func(s string) string {
		for key, value := range vars {
			if fmt.Sprintf("${%s}", key) == s {
				return fmt.Sprint(value)
			}
		}
		return s
	})
	return compiled, nil
}
