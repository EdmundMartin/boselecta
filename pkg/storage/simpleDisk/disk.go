package simpleDisk

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/EdmundMartin/boselecta/pkg/flag"
	"github.com/EdmundMartin/boselecta/pkg/storage"
	"os"
	"sync"
	"time"
)

type namespace struct {
	Flags map[string]*flag.FeatureFlag
}

func newNamespace() *namespace {
	return &namespace{
		Flags: make(map[string]*flag.FeatureFlag),
	}
}

type DiskStorage struct {
	Namespaces map[string]*namespace
	RWLock     *sync.RWMutex
	Save bool
}

func NewDiskStore(save bool) *DiskStorage {
	return &DiskStorage{Namespaces: make(map[string]*namespace),
		RWLock: &sync.RWMutex{}, Save: save}
}

func (ms *DiskStorage) All() []*flag.FeatureFlag {
	ms.RWLock.RLock()
	defer ms.RWLock.RUnlock()
	namespaces := []*namespace{}
	for _, v := range ms.Namespaces {
		namespaces = append(namespaces, v)
	}
	allFlags := []*flag.FeatureFlag{}
	for _, space := range namespaces {
		for _, v := range space.Flags {
			allFlags = append(allFlags, v)
		}
	}
	return allFlags
}

func (ms *DiskStorage) Create(name string, flag *flag.FeatureFlag) error {
	ms.RWLock.Lock()
	defer ms.RWLock.Unlock()
	space, ok := ms.Namespaces[name]
	if !ok {
		space = newNamespace()
		ms.Namespaces[name] = space
	}
	space.Flags[flag.FlagName] = flag
	ms.toDisk()
	return nil
}

func (ms *DiskStorage) GetFlag(name string, flagName string) (*flag.FeatureFlag, error) {
	ms.RWLock.RLock()
	defer ms.RWLock.RUnlock()
	res, ok := ms.Namespaces[name]
	if !ok {
		return nil, storage.MissingNamespace
	}
	fl, ok := res.Flags[flagName]
	if !ok {
		return nil, storage.MissingFlag
	}
	return fl, nil
}

func (ms *DiskStorage) Update(name string, flag *flag.FeatureFlag) error {
	ms.RWLock.Lock()
	defer ms.RWLock.Unlock()
	res, ok := ms.Namespaces[name]
	if !ok {
		return storage.MissingNamespace
	}
	_, ok = res.Flags[flag.FlagName]
	if !ok {
		return storage.MissingFlag
	}
	res.Flags[flag.FlagName] = flag
	ms.toDisk()
	return nil
}

func (ms *DiskStorage) Delete(name string, flagName string) error {
	ms.RWLock.Lock()
	defer ms.RWLock.Unlock()
	res, ok := ms.Namespaces[name]
	if !ok {
		return storage.MissingNamespace
	}
	_, ok = res.Flags[flagName]
	if !ok {
		return storage.MissingFlag
	}
	delete(res.Flags, flagName)
	if len(res.Flags) == 0 {
		delete(ms.Namespaces, name)
	}
	ms.toDisk()
	return nil
}

func (ms *DiskStorage) toDisk() error {
	if !ms.Save {
		return nil
	}
	buf := &bytes.Buffer{}
	fo, _ := os.Create(fmt.Sprintf("flagStore-%d.db", time.Now().Unix()))
	defer fo.Close()
	err := gob.NewEncoder(buf).Encode(ms.Namespaces)
	if err != nil {
		return err
	}
	_, err = fo.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}
