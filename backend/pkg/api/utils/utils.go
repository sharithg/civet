package utils

import "database/sql"

func NullFloat64ToPtr(n sql.NullFloat64) *float64 {
	if n.Valid {
		return &n.Float64
	}
	return nil
}
