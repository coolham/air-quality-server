package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// getProjectRoot 获取项目根目录
func getProjectRoot() string {
	// 尝试从可执行文件路径获取
	if execPath, err := os.Executable(); err == nil {
		// 从可执行文件路径向上两级到项目根目录
		return filepath.Join(filepath.Dir(execPath), "..", "..")
	}

	// 尝试从当前工作目录获取
	if workDir, err := os.Getwd(); err == nil {
		// 检查是否在cmd目录下
		if filepath.Base(workDir) == "cmd" || filepath.Base(filepath.Dir(workDir)) == "cmd" {
			// 从cmd目录向上到项目根目录
			return filepath.Join(workDir, "..")
		}
		// 否则假设已经在项目根目录
		return workDir
	}

	// 默认返回当前目录
	return "."
}

func main() {
	fmt.Println("=== 配置路径测试 ===")

	// 获取当前工作目录
	workDir, _ := os.Getwd()
	fmt.Printf("当前工作目录: %s\n", workDir)

	// 获取可执行文件路径
	execPath, _ := os.Executable()
	fmt.Printf("可执行文件路径: %s\n", execPath)

	// 获取项目根目录
	projectRoot := getProjectRoot()
	fmt.Printf("项目根目录: %s\n", projectRoot)

	// 测试配置文件路径
	possiblePaths := []string{
		filepath.Join(projectRoot, "config", "config.yaml"),
		filepath.Join(projectRoot, "config.yaml"),
		"config/config.yaml",
		"config.yaml",
	}

	fmt.Println("\n=== 配置文件路径测试 ===")
	for i, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("✓ 路径 %d: %s (存在)\n", i+1, path)
		} else {
			fmt.Printf("✗ 路径 %d: %s (不存在)\n", i+1, path)
		}
	}

	// 测试环境变量
	fmt.Println("\n=== 环境变量测试 ===")
	configFile := os.Getenv("CONFIG_FILE")
	if configFile != "" {
		fmt.Printf("CONFIG_FILE: %s\n", configFile)
		if !filepath.IsAbs(configFile) {
			fullPath := filepath.Join(projectRoot, configFile)
			fmt.Printf("完整路径: %s\n", fullPath)
			if _, err := os.Stat(fullPath); err == nil {
				fmt.Println("✓ 配置文件存在")
			} else {
				fmt.Println("✗ 配置文件不存在")
			}
		}
	} else {
		fmt.Println("CONFIG_FILE 未设置")
	}
}
