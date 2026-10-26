# 获取当前git分支
branch=$(git branch | grep '*' | awk '{print $2}')


# 检查本地是否有代码变动
if [ $(git diff --name-only | wc -l) -eq 0 ]; then
    echo "没有变动，无需提交"
    exit 1
fi

# 获取当前git状态
git status

# 添加文件到暂存区
git add .
if [ $? -ne 0 ]; then
    echo "添加失败，请检查"
    exit 1
fi

# 提交代码
if [ $# -eq 0 ]; then
    read -p "请输入提交信息：" commit_msg
fi
git commit -am "$commit_msg"
if [ $? -ne 0 ]; then
    echo "提交失败"
    exit 1
fi

# 检查是否需要同步远程代码
if [ $(git diff --name-only --diff-filter=U | wc -l) -ne 0 ]; then
    git pull origin $branch --rebase
fi

# 检查代码是否存在冲突
if [ $(git status | grep 'unmerged' | wc -l) -ne 0 ]; then
    echo "代码存在冲突，请手动解决冲突"
    exit 1
fi

# 推送
git push origin $branch
if [ $? -ne 0 ]; then
    echo "推送失败"
    exit 1
fi
