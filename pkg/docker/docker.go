package docker

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	// RemoteDockerImage 是Docker Hub上的Sublist3r镜像
	RemoteDockerImage = "trickest/sublist3r"
	// DockerImage 是使用的Docker镜像
	DockerImage = RemoteDockerImage
	// GitHubRepo 是Sublist3r的GitHub仓库URL
	GitHubRepo = "https://github.com/aboul3la/Sublist3r"
	// LocalDockerImage 是本地构建的Docker镜像名称
	LocalDockerImage = "sublist3r-image"
)

// IsDockerAvailable 检查Docker是否可用
func IsDockerAvailable() bool {
	cmd := exec.Command("docker", "info")
	err := cmd.Run()
	return err == nil
}

// CheckAvailability 检查Docker是否可用
func CheckAvailability() error {
	cmd := exec.Command("docker", "info")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("docker not available: %w", err)
	}
	return nil
}

// ImageExists 检查本地是否存在Sublist3r镜像
func ImageExists() bool {
	cmd := exec.Command("docker", "image", "ls", "-q", RemoteDockerImage)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) != ""
}

// RemoteImageAvailable 检查远程镜像是否可用(不下载)
func RemoteImageAvailable() bool {
	cmd := exec.Command("docker", "manifest", "inspect", RemoteDockerImage)
	err := cmd.Run()
	return err == nil
}

// IsRemoteImage 检查当前使用的是否为远程镜像
func IsRemoteImage() bool {
	// 检查本地镜像的标记信息，看是否来自远程仓库
	cmd := exec.Command("docker", "image", "inspect", "--format", "{{.RepoTags}}", DockerImage)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// 如果标记信息包含远程仓库的名称，则认为是远程镜像
	return strings.Contains(string(output), RemoteDockerImage)
}

// PullImage 从Docker Hub拉取Sublist3r镜像
func PullImage() error {
	log.Println("Pulling Sublist3r Docker image...")
	cmd := exec.Command("docker", "pull", RemoteDockerImage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// BuildImage 构建Sublist3r的Docker镜像
func BuildImage() error {
	log.Println("Building Sublist3r Docker image...")

	// 创建一个临时Dockerfile
	dockerfileContent := []byte(`FROM python:3.9-slim
RUN apt-get update && apt-get install -y git
RUN git clone https://github.com/aboul3la/Sublist3r.git /app
WORKDIR /app
RUN pip install -r requirements.txt
ENTRYPOINT ["python", "sublist3r.py"]
`)

	// 将Dockerfile写入临时文件
	err := os.WriteFile("Dockerfile", dockerfileContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to create Dockerfile: %w", err)
	}
	defer os.Remove("Dockerfile")

	// 构建Docker镜像
	cmd := exec.Command("docker", "build", "-t", LocalDockerImage, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build Docker image: %w", err)
	}

	return nil
}

// RunSublist3r 使用Sublist3r镜像进行子域名扫描
// 新版本支持传递任意命令行参数
func RunSublist3r(args []string) error {
	log.Println("Running Sublist3r Docker container...")

	// 创建临时目录以保存输出
	outputDir, err := os.MkdirTemp("", "sublist3r-output")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(outputDir)

	// 构建Docker命令
	dockerArgs := []string{
		"run",
		"--rm",
		"-v", fmt.Sprintf("%s:/output", outputDir),
	}

	// 添加Sublist3r参数
	dockerArgs = append(dockerArgs, RemoteDockerImage)
	dockerArgs = append(dockerArgs, args...)

	// 执行Docker命令
	cmd := exec.Command("docker", dockerArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker run failed: %w", err)
	}

	// 如果指定了输出文件，将其复制到指定位置
	var outputFile string
	for i, arg := range args {
		if arg == "-o" && i+1 < len(args) {
			outputFile = args[i+1]
			break
		}
	}

	if outputFile != "" {
		// 在输出目录中查找结果文件
		files, err := os.ReadDir(outputDir)
		if err != nil {
			return fmt.Errorf("failed to read output directory: %w", err)
		}

		resultFile := ""
		for _, file := range files {
			if !file.IsDir() {
				resultFile = file.Name()
				break
			}
		}

		if resultFile != "" {
			// 复制结果文件到指定的输出文件
			srcPath := fmt.Sprintf("%s/%s", outputDir, resultFile)
			src, err := os.Open(srcPath)
			if err != nil {
				return fmt.Errorf("failed to open result file: %w", err)
			}
			defer src.Close()

			dst, err := os.Create(outputFile)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer dst.Close()

			if _, err := io.Copy(dst, src); err != nil {
				return fmt.Errorf("failed to copy result to output file: %w", err)
			}
		} else {
			log.Println("Warning: No output file found in the container")
		}
	}

	return nil
}

// SimulateSublist3r 模拟Sublist3r的行为（当Docker不可用时使用）
func SimulateSublist3r(domain string, bruteforce bool, ports string, verbose bool, threads int, engines string, outputFile string, noColor bool) error {
	log.Println("Simulating Sublist3r without Docker...")

	// 等待模拟扫描时间
	time.Sleep(2 * time.Second)

	// 生成一些模拟子域名
	subdomains := []string{
		"www." + domain,
		"mail." + domain,
		"blog." + domain,
		"api." + domain,
		"dev." + domain,
		"admin." + domain,
		"m." + domain,
		"mobile." + domain,
		"ftp." + domain,
		"test." + domain,
	}

	// 添加一些随机子域名以增加多样性
	prefixes := []string{"cdn", "app", "shop", "store", "docs", "wiki", "support", "help", "secure", "vpn"}
	for _, prefix := range prefixes {
		subdomains = append(subdomains, prefix+"."+domain)
	}

	// 如果指定了输出文件，将子域名写入文件
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()

		for _, subdomain := range subdomains {
			fmt.Fprintln(file, subdomain)
		}
	}

	// 如果是详细模式，打印子域名
	if verbose {
		log.Println("Found subdomains (simulated):")
		for _, subdomain := range subdomains {
			log.Println(subdomain)
		}
	}

	return nil
}
