package cmd

import (
	"fmt"
	"os"

	"github.com/isometry/bitrat/hasher"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hashCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hashCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func hashDash(cmd *cobra.Command, args []string) {
	h := hasher.New(viper.GetString("hash"), []byte(viper.GetString("hmac")))
	hv := h.HashIoReader(os.Stdin)
	fmt.Printf("%x  (stdin)\n", hv)
}
