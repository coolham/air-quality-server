package web

import (
	"os"
	"path/filepath"
)

// Config Web模块配置
type Config struct {
	TemplatesPath string
	StaticPath    string
	AssetsPath    string
}

// GetWebPaths 获取Web资源路径
func GetWebPaths() (templatesPath, staticPath, assetsPath string) {
	// 优先级1: 环境变量指定Web根目录
	if webRoot := os.Getenv("AIR_QUALITY_WEB_ROOT"); webRoot != "" {
		templatesPath = filepath.Join(webRoot, "templates")
		staticPath = filepath.Join(webRoot, "static")
		assetsPath = filepath.Join(webRoot, "assets")
		return
	}

	// 优先级2: 智能获取项目根目录
	projectRoot := getProjectRoot()
	templatesPath = filepath.Join(projectRoot, "web", "templates")
	staticPath = filepath.Join(projectRoot, "web", "static")
	assetsPath = filepath.Join(projectRoot, "web", "assets")

	// 确保使用绝对路径
	templatesPath, _ = filepath.Abs(templatesPath)
	staticPath, _ = filepath.Abs(staticPath)
	assetsPath, _ = filepath.Abs(assetsPath)

	return
}

// getProjectRoot 智能获取项目根目录
func getProjectRoot() string {
	// 方法1: 尝试从可执行文件路径推断
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		// 检查是否在bin目录下
		if filepath.Base(exeDir) == "bin" {
			return filepath.Dir(exeDir)
		}
		// 检查是否在cmd/air-quality-server目录下
		if filepath.Base(exeDir) == "air-quality-server" {
			parent := filepath.Dir(exeDir)
			if filepath.Base(parent) == "cmd" {
				return filepath.Dir(parent)
			}
		}
	}

	// 方法2: 从当前工作目录开始向上查找
	workDir, _ := os.Getwd()
	currentDir := workDir

	for {
		// 检查当前目录是否包含go.mod文件
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir
		}

		// 检查当前目录是否包含web目录
		if _, err := os.Stat(filepath.Join(currentDir, "web")); err == nil {
			return currentDir
		}

		// 向上查找
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			// 已经到达根目录
			break
		}
		currentDir = parent
	}

	// 方法3: 使用当前工作目录作为后备
	return workDir
}
