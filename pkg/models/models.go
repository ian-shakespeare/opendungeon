package models

import (
	"database/sql"
	"errors"
	"fmt"

	"modernc.org/sqlite"
	sqllib "modernc.org/sqlite/lib"
)

var (
	ErrUnknown             = errors.New("unknown")
	ErrNotFound            = errors.New("not found")
	ErrUniqueViolation     = errors.New("unique violation")
	ErrCheckViolation      = errors.New("check violation")
	ErrForeignKeyViolation = errors.New("foreign key violation")
)

func convertSqlErr(e error) error {
	if errors.Is(e, sql.ErrNoRows) {
		return ErrNotFound
	}

	sqlErr := new(sqlite.Error)
	if errors.As(e, &sqlErr) {
		switch sqlErr.Code() {
		case sqllib.SQLITE_CONSTRAINT_CHECK:
			return ErrCheckViolation
		case sqllib.SQLITE_CONSTRAINT_UNIQUE:
			return ErrUniqueViolation
		case sqllib.SQLITE_CONSTRAINT_FOREIGNKEY:
			return ErrNotFound
		}
	}

	return fmt.Errorf("%w: %s", ErrUnknown, e.Error())
}
