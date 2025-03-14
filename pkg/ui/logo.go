package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	// 定义用于logo的颜色
	cyan    = color.New(color.FgCyan).SprintFunc()
	red     = color.New(color.FgRed, color.Bold).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
	blue    = color.New(color.FgBlue).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
)

// 定义ASCII艺术字体的Logo
var logoText = `
   _____       _     _ _     _   ____      
  / ____|     | |   | (_)   | | |___ \     
 | |  __  ___ | |___| |_ ___| |_  __) |_ __ 
 | | |_ |/ _ \| / __| | / __| __||__ <| '__|
 | |__| | (_) | \__ \ | \__ \ |_ ___) | |   
  \_____|\___/|_|___/_|_|___/\__|____/|_|   
`

// 打印版本行和链接
var versionText = "v1.0.0"
var repoLink = "https://github.com/cyberspacesec/go-Sublist3r"

// DisplayLogo 显示彩色ASCII logo和版本信息
func DisplayLogo() {
	// 随机为logo的每一行分配颜色
	lines := strings.Split(logoText, "\n")
	colorFuncs := []func(a ...interface{}) string{cyan, red, green, yellow, blue, magenta}

	// 打印logo，为每一行随机选择颜色
	for i, line := range lines {
		if line == "" {
			continue
		}
		colorIndex := i % len(colorFuncs)
		fmt.Println(colorFuncs[colorIndex](line))
	}

	// 打印版本和项目链接信息
	fmt.Printf("\n%s - %s\n\n", green(versionText), cyan(repoLink))
}

// DisplayLogoWithText 显示logo并附加自定义文本
func DisplayLogoWithText(text string) {
	DisplayLogo()
	fmt.Println(yellow(text))
	fmt.Println()
}
