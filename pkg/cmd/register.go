package cmd

// RegisterCommands 注册所有命令
func RegisterCommands() {
	// 添加所有子命令到根命令
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(buildImageCmd)
	rootCmd.AddCommand(pullImageCmd)
	rootCmd.AddCommand(apiCmd)

	// 这里可以添加更多命令
}
