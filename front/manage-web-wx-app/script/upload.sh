#!/bin/bash

echo 开始上传 $1
echo 输出目录 $2
echo 版本号 $3
yarn upload --mode=$1 --output=$2 --version=$3
if [ $? -eq 0 ]; then
  echo "上传完成"
else
  echo "上传失败"
  exit 1
fi
