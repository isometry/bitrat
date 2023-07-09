package cmd

import (
	"fmt"
	"sort"

	"github.com/isometry/bitrat/hasher"
	"github.com/spf13/cobra"
)

// listAlgorithmsCmd represents the attr command
var listAlgorithmsCmd = &cobra.Command{
	Use:   "list-algorithms",
	Short: "list supported hasher algorithms",
	Long:  ``,
	Run:   listAlgorithms,
}

func init() {
	rootCmd.AddCommand(listAlgorithmsCmd)
}

func listAlgorithms(cmd *cobra.Command, args []string) {
	fmt.Println("Supported algorithms:")
	supportedAlgorithms := hasher.SupportedAlgorithms
	algoNames := make([]string, 0, len(supportedAlgorithms))
	for name, _ := range supportedAlgorithms {
		algoNames = append(algoNames, name)
	}
	sort.Strings(algoNames)
	for _, name := range algoNames {
		fmt.Printf("- %s\n", name)
	}

}
