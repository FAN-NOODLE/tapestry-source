#!/usr/bin/env bash
# 发布环境配置

PROJECT_DIR="$WORKSPACE/"                                 # 发布工程目录
PACKAGE_NAME=`basename $0 sh`"tar.gz"                     # 打包名称

HOST=$1                                   # 服务器IP
DEPLOY_DIR=$2                             # 部署生产目录
TMP_DIR=${DEPLOY_DIR}"_tmp_dir"           # 部署临时目录


echo ">>>>> 开始发布"

echo "===> 0.构建代码 "
node -v
npm config set registry https://registry.npm.taobao.org/
npm install && \
  npm run build

#  mv dist/static/js/app.*.js dist/static/js/app.js -f && \
#  mv dist/static/js/chunk-vendors.*.js dist/static/js/chunk-vendors.js -f && \
#  mv dist/static/css/app.*.css dist/static/css/app.css -f

echo "<=== 0.成功"

echo "===> 1.打包项目dist目录下代码到: ${PACKAGE_NAME} "
tar -czf ${PACKAGE_NAME} dist
echo "<=== 1.成功"

echo "===> 2.在远端创建临时目录：${TMP_DIR} "
ssh -p 22 root@${HOST} "mkdir -p ${TMP_DIR} && cd ${TMP_DIR} && pwd"
echo "<=== 2.成功"

echo "3.从发布目录中将打包文件：${PACKAGE_NAME} 拷贝到 ${HOST} 下的临时目录：${TMP_DIR} 中"
scp -P 22 -r ${PACKAGE_NAME} root@${HOST}:${TMP_DIR}
echo "<=== 3.成功"

echo "4.解压项目源码： ${PACKAGE_NAME}"
ssh -p 22 root@${HOST} "cd ${TMP_DIR} && tar -zxf ${PACKAGE_NAME}"
echo "<=== 4.成功"

echo "5.复制项目源码： ${TMP_DIR}/dist/ "
ssh -p 22 root@${HOST} "cd ${TMP_DIR} && cp -rf dist/. ${DEPLOY_DIR}"
echo "<=== 5.成功"

echo "6.清除临时目录：${TMP_DIR}"
ssh -p 22 root@${HOST} "rm -rf ${TMP_DIR}"
echo "<=== 6.成功"

echo "7.递归修改项目根目录权限"
ssh -p 22 root@${HOST} "chown nginx:nginx ${DEPLOY_DIR} -R"
echo "<=== 7.成功"

echo "<<<<< 发布结束"
