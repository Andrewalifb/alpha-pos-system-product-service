package utils

import "github.com/google/uuid"

func ParseUUID(s string) *uuid.UUID {
	u := uuid.MustParse(s)
	return &u
}
