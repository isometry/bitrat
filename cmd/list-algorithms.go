package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/isometry/bitrat/hasher"
)

// listAlgorithmsCmd represents the attr command
func cmdListAlgorithms() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-algorithms",
		Short: "list supported hasher algorithms",
		Long:  ``,
		Run:   listAlgorithms,
	}

	return cmd
}

func listAlgorithms(cmd *cobra.Command, args []string) {
	fmt.Println("Supported algorithms:")
	supportedAlgorithms := hasher.SupportedAlgorithms
	algoNames := make([]string, 0, len(supportedAlgorithms))
	for name := range supportedAlgorithms {
		algoNames = append(algoNames, name)
	}
	sort.Strings(algoNames)
	for _, name := range algoNames {
		fmt.Printf("- %s\n", name)
	}

}
