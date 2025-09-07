@echo off
REM 测试Go代理配置是否有效

echo 🔧 测试Go代理配置...

REM 设置Go代理环境变量
set GOPROXY=https://goproxy.cn,direct
set GOSUMDB=sum.golang.google.cn

echo 当前Go代理配置:
echo   GOPROXY=%GOPROXY%
echo   GOSUMDB=%GOSUMDB%

echo.
echo 📋 测试Go模块下载...

REM 测试下载一个常用的Go模块
go list -m -versions github.com/gin-gonic/gin

if errorlevel 1 (
    echo ❌ Go模块下载测试失败
    echo 💡 请检查网络连接或尝试其他代理服务器
    echo    备用代理: https://goproxy.io,direct
) else (
    echo ✅ Go模块下载测试成功
    echo 💡 代理配置有效，可以正常构建Docker镜像
)

echo.
echo 📝 使用方法:
echo   构建生产环境: scripts\docker\docker-build-china.bat
echo   构建开发环境: scripts\docker\docker-build-china.bat dev
