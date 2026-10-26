#!/bin/bash
# 运行目录
current_directory=$(pwd)
# 打包目录
output=script/dist
# 打包环境
env=$1
# 打包分支
branch=$2
# 企业
company=$3
# 上传版本
version=$4

# 版本号是否存在
if [ -z "$version" ]; then
  echo "请输入版本号"
  exit 1
fi

# 执行 env.sh
source $current_directory/script/env.sh $env

# 执行 branch.sh
source $current_directory/script/branch.sh $branch

# 执行 build.sh
source $current_directory/script/build.sh $env $output

# 执行 upload.sh
source $current_directory/script/upload.sh $env $output $version

# 清空
rm -rf $current_directory/$output/

# 如果是生产环境，打标签
if [ $env = "prod" ]; then
  # 打标签
  git tag "$version-$company"
  # 推送标签
  git push origin "$version-$company"
fi
