<#
.SYNOPSIS
    简单的 Golang 程序 Windows 部署脚本
.DESCRIPTION
    此脚本用于编译 Golang 程序并复制到指定目录
#>

# 设置参数
$appName = "broadcaster_2"          # 应用程序名称
$installDir = "C:\Apps"     # 安装目录

# 1. 编译 Go 程序
Write-Host "building... $appName..." -ForegroundColor Cyan
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o "$appName.exe"

if (-not $?) {
    Write-Host "build failure!" -ForegroundColor Red
    exit 1
}

# 2. 创建安装目录
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

# 3. 复制文件
Copy-Item -Path "$appName.exe" -Destination "$installDir\$appName.exe" -Force
Write-Host "install app: $installDir\$appName.exe" -ForegroundColor Green

# 4. 可选: 添加到PATH
$choice = Read-Host "是否要添加到系统PATH? (y/n)"
if ($choice -eq 'y') {
    $path = [Environment]::GetEnvironmentVariable("Path", "User")
    if (-not $path.Contains($installDir)) {
        [Environment]::SetEnvironmentVariable("Path", "$path;$installDir", "User")
        Write-Host "ADD PATH" -ForegroundColor Green
    }
}

Write-Host "well done!" -ForegroundColor Green