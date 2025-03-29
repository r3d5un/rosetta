package data

import (
	"strings"
)

// CreateOrderByClause creates a SQL statement for ordering based on a given
// list of column names. Column names prefixed in "-" will be in a descending
// order.
//
// This function always adds an "id" column at the end to guarantee the order
// of cursor SQL queries.
func CreateOrderByClause(orderBy []string) string {
	length := len(orderBy)
	if length == 0 {
		return "ORDER BY id"
	}

	orderClauses := make([]string, length+1)
	for i, item := range orderBy {
		if strings.HasPrefix(item, "-") {
			orderClauses[i] = strings.TrimPrefix(item, "-") + " DESC"
		} else {
			orderClauses[i] = item + " ASC"
		}
	}
	orderClauses[len(orderClauses)-1] = "id ASC"

	return "ORDER BY " + strings.Join(orderClauses, ", ")
}
