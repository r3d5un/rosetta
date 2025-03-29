package data

import (
	"strings"
)

const (
	// PgxIntegrityConstraintViolationCode is the code for general integrity violations
	PgxIntegrityConstraintViolationCode string = "23000"
	// PgxRestrictViolationCode is the error code for deleting or updating records referenced by other resources
	PgxRestrictViolationCode = "23001"
	// PgxNotNullViolationCode is the error code for attempting to insert or update NULL values in a NOT NULL column
	PgxNotNullViolationCode = "23502"
	// PgxForeignKeyViolationCode is for foreign key violation, e.g. creating values not in the referenced table
	PgxForeignKeyViolationCode = "23503"
	// PgxUniqueViolationCode is the error code for inserting duplicate records
	PgxUniqueViolationCode = "23505"
	// PgxCheckViolationCode is for CHECK violations
	PgxCheckViolationCode = "23514"
	// PgxSyntaxErrorCode is for general syntax errors
	PgxSyntaxErrorCode = "42601"
	// PgxUndefinedColumnCode is the error code for referencing columns that does not exists
	PgxUndefinedColumnCode = "42703"
	// PgxUndefinedTableCode is the error code for refercing tables that does not exists
	PgxUndefinedTableCode = "42P01"
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
