package db

import "database/sql"

// NullStringToString converts a sql.NullString to a regular string. Returns an empty string if the NullString is not valid.
func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
