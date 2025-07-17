package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// For storing []string like text[] in PostgreSQL
type StringArray []string

// Value implements the driver.Valuer interface
func (sa StringArray) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return "{}", nil
	}

	quoted := make([]string, len(sa))
	for i, v := range sa {
		quoted[i] = `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
	}

	return "{" + strings.Join(quoted, ",") + "}", nil
}

// Scan implements the sql.Scanner interface
func (sa *StringArray) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return fmt.Errorf("cannot convert %T to StringArray", src)
	}

	str = strings.Trim(str, "{}")
	if str == "" {
		*sa = []string{}
		return nil
	}

	elems := strings.Split(str, ",")
	for i, s := range elems {
		elems[i] = strings.Trim(s, `"`)
	}
	*sa = elems

	return nil
}
