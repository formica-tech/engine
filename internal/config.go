package internal

import (
	"os"
	"strings"
)

// ReadConfig use with json path example: sql.port
func ReadConfig(name string) string {
	key := strings.ToUpper(name)
	key = strings.Replace(key, ".", "_", -1)
	return os.Getenv("FORMICA_" + key)
}

func ReadConfigWithDefault(name string, val string) string {
	v := ReadConfig(name)
	if v == "" {
		return val
	}
	return v
}
