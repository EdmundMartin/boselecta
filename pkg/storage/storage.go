package storage

import (
	"github.com/EdmundMartin/boselecta/pkg/flag"
)

type FlagStorage interface {
	GetFlag(string, string) (*flag.FeatureFlag, error)
	Create(string, *flag.FeatureFlag) error
	Update(string, *flag.FeatureFlag) error
	Delete(namespace string, flag string) error
	All() []*flag.FeatureFlag
}
