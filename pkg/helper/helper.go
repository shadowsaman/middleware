package helper

import "database/sql"

func NullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	} else {
		return sql.NullString{
			String: s,
			Valid:  true,
		}
	}
}
