package models

import (
	"database/sql"
	"regexp"
	"time"
)

type Account struct {
	ID          int64        `json:"id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
	UniqueName  string       `json:"unique_name"`
	Domain      string       `json:"domain"`
	Civitas     int64        `json:"civitas"`
	DisplayName string       `json:"display_name"`
	Tombstoned  bool         `json:"tombstoned"`
}

var AccountUniqueNameRegex = regexp.MustCompile("[a-z0-9_]+([a-z0-9_.-]+[a-z0-9_]+)?")

func ValidateAccountUniqueName(name string) bool {
	return AccountUniqueNameRegex.MatchString(name)
}
