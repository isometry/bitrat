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
		SkipGit:     viper.GetBool("skipGit"),
	}

}

func pathsToWalk(paths []string) []string {
	if len(paths) == 0 {
		paths = append(paths, viper.GetString("path"))
	}
	return paths
}

/*
func pathWalk(root string, walkFn filepath.WalkFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	err := filepath.Walk(root, walkFn)
	if err != nil {
		panic(err)
	}
}

func pathStep(rootPath string, fileChan chan<- *pathwalk.File) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR %v\n", err)
			return nil
		}

		switch {
		case info.IsDir():
			if path != rootPath &&
				((!viper.GetBool("recurse")) ||
					(viper.GetBool("hiddenDirs") && strings.HasPrefix(info.Name(), ".")) ||
					(viper.GetBool("skipGit") && info.Name() == ".git")) {
				return filepath.SkipDir
			}
			return nil
		case !info.Mode().IsRegular(), // must be after info.IsDir() checks
			viper.GetBool("hiddenFiles") && strings.HasPrefix(info.Name(), "."):
			return nil
		}

		if globMatch, _ := filepath.Match(viper.GetString("name"), info.Name()); !globMatch {
			return nil
		}

		fileChan <- &pathwalk.File{Path: path}

		return nil
	}
}
*/

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
