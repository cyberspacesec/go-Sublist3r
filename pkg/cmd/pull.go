package cmd

import (
	"log"

	"github.com/cyberspacesec/go-Sublist3r/pkg/docker"
	"github.com/cyberspacesec/go-Sublist3r/pkg/ui"
	"github.com/spf13/cobra"
)

// pullImageCmd 表示pull-docker-image子命令
var pullImageCmd = &cobra.Command{
	Use:   "pull-docker-image",
	Short: "从Docker Hub拉取Sublist3r镜像",
	Long:  `从Docker Hub拉取Sublist3r镜像并为本地使用添加标签。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查Docker可用性
		if err := docker.CheckAvailability(); err != nil {
			log.Fatalf("需要Docker但不可用: %v", err)
		}

		// 从Docker Hub拉取镜像
		ui.DisplayLogoWithText("正在从Docker Hub拉取Sublist3r镜像...")

		if err := docker.PullImage(); err != nil {
			log.Fatalf("从Docker Hub拉取镜像失败: %v", err)
		} else {
			log.Printf("成功从Docker Hub拉取Sublist3r镜像")
		}
	},
}

func init() {
	// 拉取镜像命令不需要额外的标志
}
