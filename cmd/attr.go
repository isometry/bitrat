package cmd

import (
	"fmt"
	"sync"

	"github.com/isometry/bitrat/hashattr"
	"github.com/isometry/bitrat/hasher"
	"github.com/isometry/bitrat/pathwalk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// attrCmd represents the attr command
var attrCmd = &cobra.Command{
	Use:   "attr",
	Short: "extended attribute management",
	Long:  ``,
	Run:   attrWalk,
}

func init() {
	rootCmd.AddCommand(attrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// attrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// attrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func attrWalk(cmd *cobra.Command, args []string) {
	var (
		pathWalkWaitGroup    sync.WaitGroup
		attrReaderWaitGroup  sync.WaitGroup
		pathSortWaitGroup    sync.WaitGroup
		hashPrinterWaitGroup sync.WaitGroup
	)

	fileChan := make(chan *pathwalk.File, viper.GetInt("readahead"))
	readAttrChan := make(chan *hasher.FileHash, defaultReadahead)
	sortChan := make(chan *hasher.FileHash, 1024)

	h := hasher.New(viper.GetString("hash"), []byte(viper.GetString("hmac")))

	a := hashattr.New(fmt.Sprintf("%s.%s", viper.GetString("attrPrefix"), h.Type))

	for i := 0; i < viper.GetInt("parallel"); i++ {
		attrReaderWaitGroup.Add(1)
		go a.Reader(fileChan, readAttrChan, &attrReaderWaitGroup)
	}

	hashPrinterWaitGroup.Add(1)
	go hasher.OutputTextFile(sortChan, &hashPrinterWaitGroup)

	pathSortWaitGroup.Add(1)
	if viper.GetBool("sort") {
		go hasher.SortByPath(readAttrChan, sortChan, &pathSortWaitGroup)
	} else {
		go hasher.SortByFifo(readAttrChan, sortChan, &pathSortWaitGroup)
	}

	for _, path := range pathsToWalk(args) {
		walker := pathwalk.NewWalker(path, pathwalkOptions(), fileChan, &pathWalkWaitGroup)
		pathWalkWaitGroup.Add(1)
		go walker.Walk()
	}

	pathWalkWaitGroup.Wait()
	close(fileChan)
	attrReaderWaitGroup.Wait()
	close(readAttrChan)
	pathSortWaitGroup.Wait()
	close(sortChan)
	hashPrinterWaitGroup.Wait()
}
