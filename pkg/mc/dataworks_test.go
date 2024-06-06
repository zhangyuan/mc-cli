package mc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
	assert.True(t, true, "True is true!")
}

func TestCompileTemplate(t *testing.T) {
	vars := map[string]interface{}{
		"db": "db_dev",
		"ds": "20240101",
	}

	suites := map[string]string{
		"select 1":                "select 1",
		"select * from ${db}.foo": "select * from db_dev.foo",
		"select * from ${db}.foo where ds = '${ds}'": "select * from db_dev.foo where ds = '20240101'",
	}

	for template, expected := range suites {
		text, err := CompileTemplate(template, vars)
		assert.Nil(t, err)
		assert.Equal(t, text, expected)
	}
}
