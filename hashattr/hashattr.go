package hashattr

import (
	"fmt"
	"os"
	"sync"

	"github.com/isometry/bitrat/hasher"
	"github.com/isometry/bitrat/pathwalk"
	"github.com/pkg/xattr"
)

// HashAttr implements xattr specific methods
type HashAttr struct {
	Name string
}

// New returns an initialised instance of HashAttr
func New(name string) HashAttr {
	return HashAttr{
		Name: name,
	}
}

// Get xattr value
func (attr *HashAttr) Get(path string) []byte {
	data, err := xattr.Get(path, attr.Name)
	if err != nil {
		return nil
	}
	return data
}

// Set xattr value
func (attr *HashAttr) Set(path string, data []byte) error {
	return xattr.Set(path, attr.Name, data)
}

// Remove xattr
func (attr *HashAttr) Remove(path string) error {
	return xattr.Remove(path, attr.Name)
}

// Reader goroutine
func (attr *HashAttr) Reader(input <-chan *pathwalk.File, output chan<- *hasher.FileHash, wg *sync.WaitGroup) {
	defer wg.Done()

	for item := range input {
		hash := attr.Get(item.Path)
		output <- &hasher.FileHash{File: item, Hash: hash}
	}
}

// Writer goroutine
func (attr *HashAttr) Writer(input <-chan *hasher.FileHash, wg *sync.WaitGroup) {
	defer wg.Done()

	for item := range input {
		err := attr.Set(item.File.Path, item.Hash)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}

// Remover goroutine
func (attr *HashAttr) Remover(input <-chan *pathwalk.File, wg *sync.WaitGroup) {
	defer wg.Done()

	for item := range input {
		err := attr.Remove(item.Path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}
