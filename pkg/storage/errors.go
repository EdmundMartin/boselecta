package storage

import "errors"

var (
	MissingFlag      = errors.New("matching flag does not exist")
	MissingNamespace = errors.New("matching namespace does not exist")
)
