package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/cyberspacesec/go-Sublist3r/pkg/docker"
	"github.com/cyberspacesec/go-Sublist3r/pkg/ui"
	"github.com/spf13/cobra"
)

// scanCmd 表示scan子命令
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "扫描子域名",
	Long:  `扫描指定域名的子域名，并可选择性地对发现的子域名执行端口扫描。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查必需的domain参数
		if domain == "" {
			cmd.Help()
			os.Exit(1)
		}

		// 检查Docker可用性
		dockerAvailable := true
		if err := docker.CheckAvailability(); err != nil {
			fmt.Printf("Docker不可用: %v\n", err)
			fmt.Println("将使用模拟模式运行...")
			dockerAvailable = false
		}

		// 只有当Docker可用时才检查和构建镜像
		if dockerAvailable {
			// 检查Docker镜像是否存在
			if !docker.ImageExists() {
				// 镜像不存在，先尝试从Docker Hub拉取
				ui.DisplayLogoWithText("本地未找到镜像，正在尝试从Docker Hub拉取...")

				pullErr := docker.PullImage()
				if pullErr != nil {
					fmt.Printf("从Docker Hub拉取镜像失败: %v\n", pullErr)
					fmt.Println("尝试构建本地镜像...")

					// 拉取失败，尝试构建本地镜像
					buildErr := docker.BuildImage()
					if buildErr != nil {
						fmt.Printf("构建Docker镜像失败: %v\n", buildErr)
						fmt.Println("将使用模拟模式运行...")
						dockerAvailable = false
					}
				} else {
					fmt.Println("成功从Docker Hub拉取镜像")
				}
			}
		}

		// 运行Sublist3r
		ui.DisplayLogoWithText(fmt.Sprintf("正在扫描域名: %s", domain))

		// 准备参数
		var dockerArgs []string
		dockerArgs = append(dockerArgs, "-d", domain)

		if output != "" {
			dockerArgs = append(dockerArgs, "-o", output)
		}

		if bruteforce {
			dockerArgs = append(dockerArgs, "-b")
		}

		if ports != "" {
			dockerArgs = append(dockerArgs, "-p", ports)
		}

		if verbose {
			dockerArgs = append(dockerArgs, "-v")
		}

		if threads != 30 {
			dockerArgs = append(dockerArgs, "-t", fmt.Sprintf("%d", threads))
		}

		if engines != "" {
			dockerArgs = append(dockerArgs, "-e", engines)
		}

		if noColor {
			dockerArgs = append(dockerArgs, "-n")
		}

		var err error
		if dockerAvailable {
			// 使用Docker运行
			err = docker.RunSublist3r(dockerArgs)
		} else {
			// 使用模拟模式运行
			err = docker.SimulateSublist3r(domain, bruteforce, ports, verbose, threads, engines, output, noColor)
		}

		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	},
}

func init() {
	// 添加扫描命令标志
	scanCmd.Flags().StringVarP(&domain, "domain", "d", "", "要枚举子域名的域名 (必需)")
	scanCmd.Flags().BoolVarP(&bruteforce, "bruteforce", "b", false, "启用暴力破解模块")
	scanCmd.Flags().StringVarP(&ports, "ports", "p", "", "对发现的子域名扫描指定的TCP端口")
	scanCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "启用详细输出，实时显示结果")
	scanCmd.Flags().IntVarP(&threads, "threads", "t", 30, "用于暴力破解的线程数")
	scanCmd.Flags().StringVarP(&engines, "engines", "e", "", "指定逗号分隔的搜索引擎列表")
	scanCmd.Flags().StringVarP(&output, "output", "o", "", "将结果保存到文本文件")
	scanCmd.Flags().BoolVarP(&noColor, "no-color", "n", false, "禁用输出中的颜色")

	// 标记domain为必需参数
	scanCmd.MarkFlagRequired("domain")
}
