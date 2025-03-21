package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/isometry/bitrat/hasher"
)

// hashCmd represents the hash (and default) command
func cmdDash() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stdin",
		Short: "hash stdin",
		Long:  ``,
		Run:   hashDash,
	}

	return cmd
}

func hashDash(cmd *cobra.Command, args []string) {
	h := hasher.New(viper.GetString("hash"), []byte(viper.GetString("hmac")))
	hv := h.HashIoReader(os.Stdin)
	fmt.Printf("%x  (stdin)\n", hv)
}
