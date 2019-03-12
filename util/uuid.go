package util

import "github.com/google/uuid"

func UuidStr() (s string, err error) {
	var u uuid.UUID
	if u, err = uuid.NewUUID(); err != nil {
		return
	}
	return u.String(), nil
}
