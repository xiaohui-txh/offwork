#!/bin/bash
set -e

# 设置输出目录
BIN_NAME=offwork

echo "开始编译 $BIN_NAME ..."

# 编译 main.go
go build -o $BIN_NAME main.go

echo "编译完成: ./$BIN_NAME"
