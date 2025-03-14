package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cyberspacesec/go-Sublist3r/pkg/api"
	"github.com/cyberspacesec/go-Sublist3r/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	// API服务器相关配置
	apiPort     int
	apiWorkers  int
	apiCapacity int
	apiKey      string // API密钥
)

// apiCmd 表示API子命令
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "启动HTTP API服务器",
	Long:  `启动用于Sublist3r的HTTP API服务器，允许通过HTTP调用进行远程访问。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 显示启动信息
		authMsg := ""
		if apiKey != "" {
			authMsg = " (已启用API密钥认证)"
		}
		ui.DisplayLogoWithText(fmt.Sprintf("正在启动API服务器，端口: %d%s", apiPort, authMsg))

		// 创建API服务器
		server := api.NewAPIServer(apiPort, apiWorkers, apiCapacity, apiKey)

		// 捕获系统信号以优雅关闭
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		// 启动一个goroutine监听系统信号
		go func() {
			<-c
			log.Println("正在关闭API服务器...")

			// 创建具有超时的上下文
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// 优雅关闭服务器
			if err := server.Stop(ctx); err != nil {
				log.Fatalf("服务器被强制关闭: %v", err)
			}

			log.Println("服务器已优雅关闭")
			os.Exit(0)
		}()

		// 开始服务器
		fmt.Printf("API服务器正在运行: http://localhost:%d\n", apiPort)
		fmt.Println("可用的端点:")
		fmt.Printf("API文档: http://localhost:%d/docs\n", apiPort)
		fmt.Println("POST /api/v1/scan - 开始新的子域名扫描")
		fmt.Println("POST /api/v1/scan/sync - 同步执行扫描并返回结果(等待最多5分钟)")
		fmt.Println("GET  /api/v1/scan/{id} - 获取扫描状态和结果")
		fmt.Println("GET  /api/v1/scans - 列出所有扫描")

		if apiKey != "" {
			fmt.Println("\n认证信息:")
			fmt.Println("所有API请求需要在请求头中包含 'X-API-Key: " + apiKey + "'")
			fmt.Println("或者在URL中添加 '?api_key=" + apiKey + "'")
		}

		fmt.Println("\n按Ctrl+C停止服务器")

		if err := server.Start(); err != nil {
			log.Fatalf("启动服务器时出错: %v", err)
		}
	},
}

func init() {
	// 添加API命令相关的标志
	apiCmd.Flags().IntVarP(&apiPort, "port", "p", 8080, "API服务器运行的端口")
	apiCmd.Flags().IntVarP(&apiWorkers, "workers", "w", 5, "处理请求的工作线程数")
	apiCmd.Flags().IntVarP(&apiCapacity, "capacity", "c", 100, "最大并发扫描请求数")
	apiCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "API密钥，用于请求认证（留空表示不启用认证）")
}
