package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/isometry/bitrat/hasher"
)

// hashCmd represents the hash (and default) command
var dashCmd = &cobra.Command{
	Use:   "stdin",
	Short: "hash stdin",
	Long:  ``,
	Run:   hashDash,
}

func init() {
	rootCmd.AddCommand(dashCmd)
}

func hashDash(cmd *cobra.Command, args []string) {
	h := hasher.New(viper.GetString("hash"), []byte(viper.GetString("hmac")))
	hv := h.HashIoReader(os.Stdin)
	fmt.Printf("%x  (stdin)\n", hv)
}
