#!/bin/bash

# 简单版 Go 程序部署脚本
# 使用方法: ./deploy.sh <程序名称> [start|stop|restart]

PROGRAM_NAME=$1
ACTION=${2:-"start"}
INSTALL_DIR="/usr/local/bin"  # 安装目录
LOG_FILE="/tmp/${PROGRAM_NAME}.log"  # 日志文件

# 检查参数
if [ $# -lt 1 ]; then
    echo "使用方法: $0 <程序名称> [start|stop|restart]"
    exit 1
fi

# 编译并安装程序
install() {
    echo "正在编译 $PROGRAM_NAME..."
    go build -o $PROGRAM_NAME

    if [ $? -ne 0 ]; then
        echo "编译失败"
        exit 1
    fi

    echo "正在安装到 $INSTALL_DIR"
    sudo mv $PROGRAM_NAME $INSTALL_DIR
    sudo chmod +x $INSTALL_DIR/$PROGRAM_NAME
}

# 启动程序
start() {
    echo "正在启动 $PROGRAM_NAME..."
    nohup $INSTALL_DIR/$PROGRAM_NAME > $LOG_FILE 2>&1 &
    echo "程序已启动，PID: $!"
    echo "日志输出: $LOG_FILE"
}

# 停止程序
stop() {
    echo "正在停止 $PROGRAM_NAME..."
    pkill -f $INSTALL_DIR/$PROGRAM_NAME
    echo "程序已停止"
}

# 根据参数执行操作
case $ACTION in
    "start")
        install
        start
        ;;
    "stop")
        stop
        ;;
    "restart")
        stop
        start
        ;;
    *)
        echo "无效操作: $ACTION"
        echo "使用方法: $0 <程序名称> [start|stop|restart]"
        exit 1
        ;;
esac

exit 0