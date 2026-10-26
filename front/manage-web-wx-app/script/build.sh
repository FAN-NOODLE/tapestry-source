#!/bin/bash
echo 开始打包 $1

yarn build --env=$1 --output-path=$2
if [ $? -eq 0 ]; then
  echo "打包完成"
else
  echo "打包失败"
  exit 1
fi
