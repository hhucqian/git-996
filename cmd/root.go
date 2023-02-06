package cmd

import (
	"fmt"
	"git-996/cmd/arg"
	"git-996/execute"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "git-996",
	Short:   "git-996 是一个统计代码提交的工具",
	Example: "git-996 <git path>",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		execute.LoadAndPrintFromPath(args[0])
	},
}

func init() {
	rootCmd.Flags().StringVarP(&arg.RootArg.Format, "format", "f", "table", "输出格式, table | json")
	rootCmd.Flags().StringSliceVar(&arg.RootArg.MergeEMail, "merge-email", nil, "合并人员 例如：user-to@mail.com=user-from@mail.com")
	rootCmd.Flags().StringVarP(&arg.RootArg.Sort, "sort", "s", "l", "排序方式： i | increase | d | decrease | l | left")
	rootCmd.Flags().BoolVarP(&arg.RootArg.Revert, "revert", "r", true, "逆序")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
