package utils

import (
	"fmt"
	"strings"
)

// BuildInClause builds a safe SQL "IN" clause with numbered placeholders
// and returns the clause string (like `($1, $2, $3)`) and the args slice.
func BuildInClause(startIndex int, items []any) (string, []any) {
	if len(items) == 0 {
		return "(NULL)", nil
	}

	placeholders := make([]string, len(items))
	args := make([]any, len(items))

	for i, val := range items {
		placeholders[i] = fmt.Sprintf("$%d", startIndex+i)
		args[i] = val
	}

	return fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")), args
}
