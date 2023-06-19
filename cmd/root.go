/*
Copyright © 2023 Daichi TAKAHASHI
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/daichitakahashi/modlist/internal/golangci"
	"github.com/spf13/cobra"
	"golang.org/x/exp/rand"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "modlist",
	Short: "List your go modules/packages.",
	Long:  "List your go modules/packages.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := directoryPath != nil && *directoryPath

		list, err := listModules(dir)
		if err != nil {
			return err
		}
		if packages != nil && *packages {
			list, err = expand(list, func(s string) ([]string, error) {
				if dir {
					return listPackagePaths(s)
				}
				return listPackageNames(s)
			})
			if err != nil {
				return err
			}
		}
		if shuffle != nil && *shuffle {
			rand.Seed(uint64(time.Now().UnixNano()))
			rand.Shuffle(len(list), func(i, j int) {
				list[i], list[j] = list[j], list[i]
			})
		}

		f := filters([]filterFunc{})
		if matchPatterns != nil && len(*matchPatterns) > 0 {
			f = append(f, matchFilter(*matchPatterns))
		}
		if excludePatterns != nil && len(*excludePatterns) > 0 {
			f = append(f, excludeFilter(*excludePatterns))
		}
		if golangCILintSkipDirs != nil && *golangCILintSkipDirs {
			cfg, err := golangci.ReadConfig()
			if err != nil {
				if os.IsNotExist(err) {
					log.Println("golangci-lint configuration file not found")
				} else {
					return err
				}
			}
			if err == nil {
				rxs, err := cfg.SkipDirectories()
				if err != nil {
					return err
				}
				f = append(f, excludeRegexpFilter(rxs))
			}
		}

		filtered := make([]string, 0, len(list))
		for _, item := range list {
			if f.excluded(item) {
				continue
			}
			filtered = append(filtered, item)
		}
		lastIdx := len(filtered) - 1
		for idx, item := range filtered {
			fmt.Print(item)
			if idx != lastIdx {
				fmt.Print(*separator)
			}
		}
		fmt.Print("\n")
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

var (
	packages             *bool
	shuffle              *bool
	matchPatterns        *[]string
	excludePatterns      *[]string
	separator            *string
	directoryPath        *bool
	golangCILintSkipDirs *bool
)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.modlist.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	packages = rootCmd.Flags().BoolP("packages", "p", false, "")
	shuffle = rootCmd.Flags().BoolP("shuffle", "s", false, "shuffle module list")
	matchPatterns = rootCmd.Flags().StringArrayP("match", "m", nil, "filter unmatch items")
	excludePatterns = rootCmd.Flags().StringArrayP("exclude", "e", nil, "filter match items")
	separator = rootCmd.Flags().String("separator", "\n", "separator")
	directoryPath = rootCmd.Flags().BoolP("directory", "d", false, "show directory instead of module/package name")
	golangCILintSkipDirs = rootCmd.Flags().Bool("golangci-lint-skip-dirs", false, "if configuration file exists, read run.skip-dirs and run.skip-dirs-use-default option")
}
