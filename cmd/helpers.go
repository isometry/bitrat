package cmd

import (
	"crypto/hmac"
	"fmt"
	"os"
	"sync"

	"github.com/isometry/bitrat/hasher"
	"github.com/isometry/bitrat/pathwalk"
	"github.com/spf13/viper"
)

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func pathwalkOptions() *pathwalk.Options {
	return &pathwalk.Options{
		Pattern:     viper.GetString("name"),
		Recurse:     viper.GetBool("recurse"),
		HiddenDirs:  viper.GetBool("hiddenDirs"),
		HiddenFiles: viper.GetBool("hiddenFiles"),
		IncludeGit:  viper.GetBool("includeGit"),
		Parallel:    viper.GetInt("parallel"),
	}
}

func pathsToWalk(paths []string) []string {
	if len(paths) == 0 {
		paths = append(paths, viper.GetString("path"))
	}
	return paths
}

func hashConsumer(input <-chan *hasher.FileHash, wg *sync.WaitGroup) {
	defer wg.Done()

	for item := range input {
		fmt.Printf("%x  %s\n", item.Hash, item.File.Path)
	}
}

// Diff two hashes
func hashDiff(fileHash []byte, attrHash []byte) string {
	switch {
	case fileHash == nil:
		return "!"
	case attrHash == nil:
		return "+"
	case hmac.Equal(fileHash, attrHash):
		return "="
	case attrHash != nil:
		return "~"
	default:
		return "?"
	}
}
