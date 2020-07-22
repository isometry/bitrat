package cmd

import (
	"sync"

	"github.com/isometry/bitrat/hasher"
	"github.com/isometry/bitrat/pathwalk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// hashCmd represents the hash (and default) command
var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "generate file hashes fast",
	Long:  ``,
	Run:   hashWalk,
}

func init() {
	rootCmd.AddCommand(hashCmd)
}

/*
 * Walk the supplied paths finding files that match the supplied criteria,
 * pass these files to a scalable number of hashProcessors, pass the results
 * to an optional sorter, and output the results.
 *                 ┌─┐                       ┌─┐                 ┌─┐
 *                 │p│    ┌─────────────┐    │h│                 │s│
 * ┌──────────┐    │a│ ┌─>│HashProcessor│─┐  │a│                 │o│
 * │ pathWalk │─┐  │t│ │  └─────────────┘ │  │s│                 │r│
 * └──────────┘ │  │h│ │  ┌─────────────┐ │  │h│   ┌─────────┐   │t│   ┌───────────┐
 *              ├─>│C│─┼─>│HashProcessor│─┼─>│C│──>│SortBy...│──>│C│──>│ Output... │
 * ┌──────────┐ │  │h│ │  └─────────────┘ │  │h│   └─────────┘   │h│   └───────────┘
 * │ pathWalk │─┘  │a│ │  ┌─────────────┐ │  │a│                 │a│
 * └──────────┘    │n│ └─>│HashProcessor│─┘  │n│                 │n│
 *                 └─┘    └─────────────┘    └─┘                 └─┘
 */
func hashWalk(cmd *cobra.Command, args []string) {
	var (
		fileWaitGroup   sync.WaitGroup
		hashWaitGroup   sync.WaitGroup
		sortWaitGroup   sync.WaitGroup
		outputWaitGroup sync.WaitGroup
	)

	fileChan := make(chan *pathwalk.File, viper.GetInt("readahead"))
	hashChan := make(chan *hasher.FileHash, defaultWriteahead)
	sortChan := make(chan *hasher.FileHash, 1024)

	outputWaitGroup.Add(1)
	if viper.GetBool("protobuf") {
		go hasher.OutputProtobufFile(sortChan, &outputWaitGroup)
	} else {
		go hasher.OutputTextFile(sortChan, &outputWaitGroup)
	}

	sortWaitGroup.Add(1)
	if viper.GetBool("sort") {
		go hasher.SortByPath(hashChan, sortChan, &sortWaitGroup)
	} else {
		go hasher.SortByFifo(hashChan, sortChan, &sortWaitGroup)
	}

	for i := 0; i < viper.GetInt("parallel"); i++ {
		hashWaitGroup.Add(1)
		h := hasher.New(viper.GetString("hash"), []byte(viper.GetString("hmac")))
		go h.HashProcessor(fileChan, hashChan, &hashWaitGroup)
	}

	for _, path := range pathsToWalk(args) {
		var walker pathwalk.PathWalker
		if viper.GetBool("alt-walker") {
			walker = pathwalk.NewAltWalker(path, pathwalkOptions(), fileChan, &fileWaitGroup)
		} else {
			walker = pathwalk.NewWalker(path, pathwalkOptions(), fileChan, &fileWaitGroup)
		}
		fileWaitGroup.Add(1)
		go walker.Walk()
	}

	fileWaitGroup.Wait()
	close(fileChan)
	hashWaitGroup.Wait()
	close(hashChan)
	sortWaitGroup.Wait()
	close(sortChan)
	outputWaitGroup.Wait()
}
