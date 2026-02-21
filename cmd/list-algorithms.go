package cmd

import (
	"fmt"
	"maps"
	"slices"

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
	for _, name := range slices.Sorted(maps.Keys(hasher.SupportedAlgorithms)) {
		fmt.Printf("- %s\n", name)
	}

}
