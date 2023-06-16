/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// packagesCmd represents the packages command
var packagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		list, err := listModules()
		if err != nil {
			return err
		}
		list, err = expand(list, func(s string) ([]string, error) {
			return listPackages(s)
		})
		if err != nil {
			return err
		}
		for _, pkg := range list {
			fmt.Println(pkg)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(packagesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// packagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
