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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "bitrat",
	Short:   "Lightning-fast, multi-algorithm file checksums",
	Run:     hashWalk,
	Args:    cobra.ArbitraryArgs,
	Version: fmt.Sprintf("%s-%s (built %s)", version, commit, date),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().Bool("help", false, "help for "+rootCmd.Name())

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bitrat.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().Bool("debug", false, "print diagnostics")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().Bool("stats", false, "print statistics")
	viper.BindPFlag("stats", rootCmd.PersistentFlags().Lookup("stats"))

	rootCmd.PersistentFlags().String("print-format", defaultPrintFormat, "print format")
	viper.BindPFlag("print-format", rootCmd.PersistentFlags().Lookup("print-format"))

	rootCmd.PersistentFlags().StringP("output-file", "o", defaultOutputFile, "output file")
	viper.BindPFlag("output-file", rootCmd.PersistentFlags().Lookup("output-file"))

	rootCmd.PersistentFlags().StringP("hash", "h", defaultHash, "hash algorithm")
	viper.BindPFlag("hash", rootCmd.PersistentFlags().Lookup("hash"))

	rootCmd.PersistentFlags().StringP("hmac", "k", defaultKey, "HMAC key")
	viper.BindPFlag("hmac", rootCmd.PersistentFlags().Lookup("hmac"))

	rootCmd.PersistentFlags().IntP("parallel", "j", cpuid.CPU.PhysicalCores+1, "number of parallel hashers")
	viper.BindPFlag("parallel", rootCmd.PersistentFlags().Lookup("parallel"))

	rootCmd.PersistentFlags().StringP("path", "p", defaultPath, "base path")
	viper.BindPFlag("path", rootCmd.PersistentFlags().Lookup("path"))

	rootCmd.PersistentFlags().StringP("name", "n", defaultNameGlob, "file glob pattern")
	viper.BindPFlag("name", rootCmd.PersistentFlags().Lookup("name"))

	rootCmd.PersistentFlags().Int("readahead", defaultReadahead, "file walk readahead distance")
	viper.BindPFlag("readahead", rootCmd.PersistentFlags().Lookup("readahead"))

	rootCmd.PersistentFlags().BoolP("recursive", "r", defaultRecurse, "recurse into directories")
	viper.BindPFlag("recurse", rootCmd.PersistentFlags().Lookup("recursive"))

	rootCmd.PersistentFlags().BoolP("sort", "s", defaultSort, "sort output by path")
	viper.BindPFlag("sort", rootCmd.PersistentFlags().Lookup("sort"))

	rootCmd.PersistentFlags().Bool("hidden-dirs", defaultHiddenDirs, "process hidden directories")
	viper.BindPFlag("hidden-dirs", rootCmd.PersistentFlags().Lookup("hidden-dirs"))

	rootCmd.PersistentFlags().Bool("hidden-files", defaultHiddenFiles, "process hidden files")
	viper.BindPFlag("hidden-files", rootCmd.PersistentFlags().Lookup("hidden-files"))

	rootCmd.PersistentFlags().Bool("include-git", defaultIncludeGit, "include .git directories")
	viper.BindPFlag("include-git", rootCmd.PersistentFlags().Lookup("include-git"))

	// TODO: alt-walker is likely broken
	rootCmd.PersistentFlags().Bool("alt-walker", defaultAltWalker, "use alternate pathwalker")
	viper.BindPFlag("alt-walker", rootCmd.PersistentFlags().Lookup("alt-walker"))

	rootCmd.PersistentFlags().Bool("protobuf", defaultProtobuf, "output to protobuf")
	viper.BindPFlag("protobuf", rootCmd.PersistentFlags().Lookup("protobuf"))

	rootCmd.PersistentFlags().StringSliceP("exclude", "e", nil, "exclude paths by pattern")
	viper.BindPFlag("exclude", rootCmd.PersistentFlags().Lookup("exclude"))
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
