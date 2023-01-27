package cmd

import (
	"git-analyse/repository"

	"github.com/spf13/cobra"
)

func init() {
	loadCmd.Flags().String("db", "db.sqlite", "db path")
	loadCmd.SetArgs([]string{"."})
	rootCmd.AddCommand(loadCmd)
}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "加载git数据",
	Long:  `加载git数据`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := repository.LoadFromPath(args[0], cmd.Flag("db").Value.String()); err != nil {
			panic(err.Error())
		}
	},
}
