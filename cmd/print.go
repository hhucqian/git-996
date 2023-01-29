package cmd

import (
	"git-analyse/executor"

	"github.com/spf13/cobra"
)

func init() {
	printCmd.Flags().String("db", "", "db path")
	rootCmd.MarkFlagRequired("db")
	rootCmd.AddCommand(printCmd)
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "打印统计信息",
	Long:  `打印统计信息`,
	Run: func(cmd *cobra.Command, args []string) {
		executor.Print(cmd.Flag("db").Value.String())
	},
}
