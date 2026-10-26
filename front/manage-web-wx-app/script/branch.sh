#!/bin/bash
# 指定分支
branch=$1

# 获取当前所在分支
current_branch=$(git rev-parse --abbrev-ref HEAD)
if [ $branch != $current_branch ]; then
  current_branch=$branch
  # 检查本地是否存在该分支
  if [ $(git branch | grep $current_branch | wc -l) -eq 0 ]; then
    # 不存在则创建
    git checkout -b $current_branch
  else
    # 存在则切换
    git checkout $current_branch
  fi
fi

echo 当前分支: $current_branch

# 拉取最新代码
git pull origin $current_branch
if [ $? -ne 0 ]; then
  echo 拉取最新代码失败
  exit 1
fi
