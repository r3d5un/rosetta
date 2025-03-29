package data

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrRecordNotFound                = errors.New("record not found")
	ErrForeignKeyConstraintViolation = errors.New("foreign key constraint violation")
	ErrUniqueConstraintViolation     = errors.New("unique constraint violation")
	ErrNotNullConstraintViolation    = errors.New("not null constraint violation")
	ErrCheckConstraintViolation      = errors.New("check constraint violation")
	ErrSynatxErrorViolation          = errors.New("SQL syntax errors")
	ErrUndefinedResource             = errors.New("undefined resource")
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

func handleError(err error, logger *slog.Logger) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			logger.Error("unique constraint violation", slog.String("error", pgErr.Message))
			return ErrUniqueConstraintViolation
		case "23503": // foreign_key_violation
			logger.Error("foreign key constraint violation", slog.String("error", pgErr.Message))
			return ErrForeignKeyConstraintViolation
		case "23514": // check_violation
			logger.Error("check constraint violation", slog.String("error", pgErr.Message))
			return ErrCheckConstraintViolation
		default:
			logger.Error("unhandled constraint violation", slog.String("error", pgErr.Message))
			return fmt.Errorf("Unhandled constraint violation: %s", pgErr.Message)
		}
	} else if errors.Is(err, pgx.ErrNoRows) {
		return ErrRecordNotFound
	} else {
		logger.Error("unable to perform query", slog.String("error", err.Error()))
		return err
	}
}
