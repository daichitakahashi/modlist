/*
Copyright © 2023 Daichi TAKAHASHI
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/IGLOU-EU/go-wildcard"
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
		list, err := listModules()
		if err != nil {
			return err
		}
		if packages != nil && *packages {
			list, err = expand(list, func(s string) ([]string, error) {
				return listPackages(s)
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
		filtered := make([]string, 0, len(list))
		for _, item := range list {
			_match, tested := match(item, matchPatterns)
			if tested && !_match {
				continue
			}
			_match, tested = match(item, excludePatterns)
			if tested && _match {
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
	packages        *bool
	shuffle         *bool
	matchPatterns   *[]string
	excludePatterns *[]string
	separator       *string
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
}

func match(s string, patterns *[]string) (match, tested bool) {
	if patterns == nil || len(*patterns) == 0 {
		return false, false
	}
	for _, p := range *patterns {
		if wildcard.MatchSimple(p, s) {
			return true, true
		}
	}
	return false, true
}
