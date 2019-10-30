package glue

import (
	"errors"
	"fmt"
	"strings"
)

const (
	rds_glue_prefix = "###glue"
)

var (
	charInvalidForName = errors.New("invalid char for a glue name, '_' and ':' can not present in a name")
	glueNameInvalid    = errors.New("invalid glue name")
)

func hasSpecialChar(names ...string) bool {
	for _, n := range names {
		for _, c := range n {
			switch c {
			case '_', ':':
				return true
			}
		}
	}

	return false
}

func formatServerKey(server string) (string, error) {
	if hasSpecialChar(server) {
		return "", charInvalidForName
	}

	return fmt.Sprintf("%s_%s", rds_glue_prefix, server), nil
}

func formatClientKey(server, client string) (string, error) {
	if hasSpecialChar(server, client) {
		return "", charInvalidForName
	}

	return fmt.Sprintf("%s_%s:%s", rds_glue_prefix, server, client), nil
}

func parseKey(key string) (string, error) {
	var i = strings.LastIndexByte(key, ':')
	if i == -1 {
		i = strings.LastIndexByte(key, '_')
	}
	if i == -1 {
		return key, glueNameInvalid
	}

	return key[i:], nil
}
