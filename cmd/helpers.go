package cmd

import (
	"crypto/hmac"
	"fmt"
	"sync"

	"github.com/spf13/viper"

	"github.com/isometry/bitrat/hasher"
	"github.com/isometry/bitrat/pathwalk"
)

func PathwalkOptions() *pathwalk.Options {
	return &pathwalk.Options{
		Pattern:     viper.GetString("name"),
		Recurse:     viper.GetBool("recurse"),
		HiddenDirs:  viper.GetBool("hidden-dirs"),
		HiddenFiles: viper.GetBool("hidden-files"),
		IncludeGit:  viper.GetBool("include-git"),
		Parallel:    viper.GetInt("parallel"),
	}
}

func PathsToWalk(paths []string) []string {
	if len(paths) == 0 {
		paths = append(paths, viper.GetString("path"))
	}
	return paths
}

func HashConsumer(input <-chan *hasher.FileHash, wg *sync.WaitGroup) {
	defer wg.Done()

	for item := range input {
		fmt.Printf("%x  %s\n", item.Hash, item.File.Path)
	}
}

// Diff two hashes
func HashDiff(fileHash []byte, attrHash []byte) string {
	switch {
	case fileHash == nil:
		return "!"
	case attrHash == nil:
		return "+"
	case hmac.Equal(fileHash, attrHash):
		return "="
	case len(attrHash) > 0:
		return "~"
	default:
		return "?"
	}
}
