/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(file)
		fmt.Println("export called")
	},
}

var file string

func init() {
	// 作用：把 export 子命令 注册到根命令
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// 作用：给 export 命令添加一个「字符串类型的 flag」
	exportCmd.Flags().StringVarP(
		&file,
		"file",           // 长名字 --file使用"f",
		"f",              // 短名字 -f使用
		"local",          // 不传的时候的默认值
		"file to output") // --help 时显示的说明

	// yourapp export --file=test.json
	// yourapp export -f=imooc
}
