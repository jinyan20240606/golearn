/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ch03",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ch03.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// 给【根命令 rootCmd】添加一个 布尔类型（bool）的 flag，只要加了，就是true
	// yourapp --toggle
	// yourapp -t
	rootCmd.Flags().BoolP(
		"toggle",                  // 1. 长名称：--toggle
		"t",                       // 2. 短名称：-t
		false,                     // 3. 默认值：false（不写就是 false）
		"Help message for toggle") // 4. 帮助说明文字
}
