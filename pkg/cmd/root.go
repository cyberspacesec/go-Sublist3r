package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// 命令行选项变量
	domain     string
	bruteforce bool
	ports      string
	verbose    bool
	threads    int
	engines    string
	output     string
	noColor    bool
)

// rootCmd 表示没有子命令时的基本命令
var rootCmd = &cobra.Command{
	Use:   "go-sublist3r",
	Short: "go-Sublist3r - 子域名枚举工具",
	Long: `go-Sublist3r 是一个用Go语言编写的子域名枚举工具，
它使用Docker来运行Sublist3r，为您提供简单的命令行界面。

示例用法:
  go-sublist3r scan -d example.com
  go-sublist3r scan -d example.com -o results.txt
  go-sublist3r scan -d example.com -v -p 80,443
  go-sublist3r scan -d example.com -b -t 50
  go-sublist3r scan -d example.com -e "baidu,yahoo"
  go-sublist3r api --port 8080 --workers 10`,
	Run: func(cmd *cobra.Command, args []string) {
		// 如果没有子命令或标志，显示帮助
		cmd.Help()
	},
}

// Execute 将所有子命令添加到根命令并设置标志
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// 注册所有命令
	RegisterCommands()

	// 添加全局标志
	rootCmd.PersistentFlags().BoolP("help", "h", false, "显示帮助信息")

	// 输出ASCII艺术标志
	cobra.AddTemplateFunc("ASCII", func() string {
		return `
 _____      _____ _     _    _ _____  _____      
|  __ \    / ____| |   | |  | |_   _|/ ____|     
| |  \/ __| (___ | |_  | |  | | | | | (___   ___ 
| | __ / _ \___ \| __| | |  | | | |  \___ \ / __|
| |_\ \  __/___) | |_  | |__| |_| |_ ____) | (__ 
 \____/\___|____/ \__|  \____/|_____|_____/ \___|
                                                 
                   version 1.0                   
        https://github.com/cyberspacesec/go-Sublist3r      
`
	})

	rootCmd.SetUsageTemplate(`{{ASCII}}
用法:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

别名:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

示例:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

可用命令:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

标志:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

全局标志:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

其他帮助主题:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

使用 "{{.CommandPath}} [command] --help" 获取有关命令的更多信息。{{end}}
`)
}
