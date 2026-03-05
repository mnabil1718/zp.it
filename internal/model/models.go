package model

import "errors"

type Models struct {
	Lookup ILookup
}

func NewModels(lu ILookup) *Models {
	return &Models{
		Lookup: lu,
	}
}

var (
	ErrNotFound = errors.New("record not found")
)
