/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	source "meow/source"

	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args:    cobra.ExactArgs(1),
	Example: `  meow exec myscript.meow`,
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		source.Start(filepath)
	},
}

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args:    cobra.ExactArgs(1),
	Example: `  meow debug myscript.meow`,
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		source.DebugTree(filepath)
	},
}

var tokensCmd = &cobra.Command{
	Use:   "tokenize",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args:    cobra.ExactArgs(1),
	Example: `  meow tokenize myscript.meow`,
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		source.DebugTokens(filepath)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(debugCmd)
	rootCmd.AddCommand(tokensCmd)
}
