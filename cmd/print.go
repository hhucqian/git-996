package cmd

import (
	"git-analyse/print"

	"github.com/spf13/cobra"
)

func init() {
	printCmd.Flags().String("db", "db.sqlite", "db path")
	rootCmd.AddCommand(printCmd)
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "打印统计信息",
	Long:  `打印统计信息`,
	Run: func(cmd *cobra.Command, args []string) {
		print.Print(cmd.Flag("db").Value.String())
	},
}
