package cmd

import (
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/isometry/bitrat/hasher"
	"github.com/isometry/bitrat/pathwalk"
)

// hashCmd represents the hash (and default) command
func cmdHash() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hash",
		Short: "generate file hashes fast",
		Long:  ``,
		Run:   hashWalk,
	}

	return cmd
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

	if viper.GetBool("protobuf") {
		outputWaitGroup.Go(hasher.OutputProtobufFile(sortChan))
	} else {
		outputWaitGroup.Go(hasher.OutputTextFile(sortChan))
	}

	if viper.GetBool("sort") {
		sortWaitGroup.Go(hasher.SortByPath(hashChan, sortChan))
	} else {
		sortWaitGroup.Go(hasher.SortByFifo(hashChan, sortChan))
	}

	for i := 0; i < viper.GetInt("parallel"); i++ {
		h := hasher.New(viper.GetString("hash"), []byte(viper.GetString("hmac")))
		hashWaitGroup.Go(h.HashProcessor(fileChan, hashChan))
	}

	for _, path := range PathsToWalk(args) {
		var walker pathwalk.PathWalker
		if viper.GetBool("alt-walker") {
			walker = pathwalk.NewAltWalker(path, PathwalkOptions(), fileChan, &fileWaitGroup)
		} else {
			walker = pathwalk.NewWalker(path, PathwalkOptions(), fileChan, &fileWaitGroup)
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
