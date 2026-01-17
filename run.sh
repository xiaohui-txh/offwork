#!/bin/bash
set -e

# 可执行文件名称
BIN_NAME=offwork

# 日志目录
LOG_DIR=logs
LOG_FILE=$LOG_DIR/offwork.log

# 创建日志目录
mkdir -p $LOG_DIR

echo "启动 $BIN_NAME ..."

# 后台运行，并将 stdout/stderr 输出到日志
nohup ./$BIN_NAME > $LOG_FILE 2>&1 &

# 输出进程号
echo "$BIN_NAME 已启动, PID: $!"
echo "日志文件: $LOG_FILE"
