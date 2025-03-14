package cmd

import (
	"fmt"
	"log"

	"github.com/cyberspacesec/go-Sublist3r/pkg/docker"
	"github.com/cyberspacesec/go-Sublist3r/pkg/ui"
	"github.com/spf13/cobra"
)

// buildImageCmd 表示build-docker-image子命令
var buildImageCmd = &cobra.Command{
	Use:   "build-docker-image",
	Short: "构建Sublist3r Docker镜像",
	Long:  `构建包含Sublist3r的Docker镜像，用于子域名扫描。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查Docker可用性
		if err := docker.CheckAvailability(); err != nil {
			log.Fatalf("需要Docker但不可用: %v", err)
		}

		// 首先尝试从Docker Hub拉取镜像
		ui.DisplayLogoWithText("正在尝试从Docker Hub拉取镜像...")

		pullErr := docker.PullImage()
		if pullErr != nil {
			fmt.Printf("从Docker Hub拉取镜像失败: %v\n", pullErr)
			fmt.Println("尝试构建本地镜像...")

			// 拉取失败，尝试构建本地镜像
			buildErr := docker.BuildImage()
			if buildErr != nil {
				log.Fatalf("构建Docker镜像失败: %v", buildErr)
			} else {
				fmt.Println("成功构建本地Docker镜像")
			}
		} else {
			fmt.Println("成功从Docker Hub拉取镜像")
		}
	},
}

func init() {
	// 构建镜像命令不需要额外的标志
}
