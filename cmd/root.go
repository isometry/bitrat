package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/klauspost/cpuid/v2"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version string = "snapshot"
	commit  string = "unknown"
	date    string = "unknown"
)

var cfgFile string

// New returns the root command
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bitrat",
		Short:   "Lightning-fast, multi-algorithm file checksums",
		Run:     hashWalk,
		Args:    cobra.ArbitraryArgs,
		Version: fmt.Sprintf("%s-%s (built %s)", version, commit, date),
	}

	cobra.OnInitialize(initConfig)

	pFlags := cmd.PersistentFlags()

	pFlags.Bool("help", false, "help for "+cmd.Name())

	pFlags.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bitrat.yaml)")
	viper.BindPFlag("config", pFlags.Lookup("config"))

	pFlags.Bool("debug", false, "print diagnostics")
	viper.BindPFlag("debug", pFlags.Lookup("debug"))

	pFlags.Bool("stats", false, "print statistics")
	viper.BindPFlag("stats", pFlags.Lookup("stats"))

	pFlags.String("print-format", defaultPrintFormat, "print format")
	viper.BindPFlag("print-format", pFlags.Lookup("print-format"))

	pFlags.StringP("output-file", "o", defaultOutputFile, "output file")
	viper.BindPFlag("output-file", pFlags.Lookup("output-file"))

	pFlags.StringP("hash", "h", defaultHash, "hash algorithm")
	viper.BindPFlag("hash", pFlags.Lookup("hash"))

	pFlags.StringP("hmac", "k", defaultKey, "HMAC key")
	viper.BindPFlag("hmac", pFlags.Lookup("hmac"))

	pFlags.IntP("parallel", "j", cpuid.CPU.PhysicalCores+1, "number of parallel hashers")
	viper.BindPFlag("parallel", pFlags.Lookup("parallel"))

	pFlags.StringP("path", "p", defaultPath, "base path")
	viper.BindPFlag("path", pFlags.Lookup("path"))

	pFlags.StringP("name", "n", defaultNameGlob, "file glob pattern")
	viper.BindPFlag("name", pFlags.Lookup("name"))

	pFlags.Int("readahead", defaultReadahead, "file walk readahead distance")
	viper.BindPFlag("readahead", pFlags.Lookup("readahead"))

	pFlags.BoolP("recursive", "r", defaultRecurse, "recurse into directories")
	viper.BindPFlag("recurse", pFlags.Lookup("recursive"))

	pFlags.BoolP("sort", "s", defaultSort, "sort output by path")
	viper.BindPFlag("sort", pFlags.Lookup("sort"))

	pFlags.Bool("hidden-dirs", defaultHiddenDirs, "process hidden directories")
	viper.BindPFlag("hidden-dirs", pFlags.Lookup("hidden-dirs"))

	pFlags.Bool("hidden-files", defaultHiddenFiles, "process hidden files")
	viper.BindPFlag("hidden-files", pFlags.Lookup("hidden-files"))

	pFlags.Bool("include-git", defaultIncludeGit, "include .git directories")
	viper.BindPFlag("include-git", pFlags.Lookup("include-git"))

	// TODO: alt-walker is likely broken
	pFlags.Bool("alt-walker", defaultAltWalker, "use alternate pathwalker")
	viper.BindPFlag("alt-walker", pFlags.Lookup("alt-walker"))

	pFlags.Bool("protobuf", defaultProtobuf, "output to protobuf")
	viper.BindPFlag("protobuf", pFlags.Lookup("protobuf"))

	pFlags.StringSliceP("exclude", "e", nil, "exclude paths by pattern")
	viper.BindPFlag("exclude", pFlags.Lookup("exclude"))

	cmd.AddCommand(
		cmdDash(),
		cmdHash(),
		cmdListAlgorithms(),
	)

	return cmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".bitrat" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bitrat")
	}

	viper.SetEnvPrefix("bitrat")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("debug") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
