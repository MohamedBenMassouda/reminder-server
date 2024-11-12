package utils

import (
	"database/sql"
	"errors"
)

const (
	ErrorUserNotFound      = "User not found"
	ErrorInvalidPassword   = "Invalid password"
	ErrorUserAlreadyExists = "User already exists"
	ErrorFailedToHash      = "Failed to hash password"
	ErrorInternalServer    = "Internal server error"
	ErrorCategoryNotFound  = "Category not found"
)

func ErrorSqlNoRows(err error) error {
	if err == sql.ErrNoRows {
		return errors.New("record not found")
	}

	return err
}
