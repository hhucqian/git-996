package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Debug bool

var rootCmd = &cobra.Command{
	Use:   "git-analyse",
	Short: "git-analyse 是一个分析git提交的工具",
	Long:  `git-analyse 是一个分析git提交的工具`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	rootCmd.Flags().BoolVar(&Debug, "debug", false, "debug mode")
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
