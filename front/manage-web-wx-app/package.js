
const { execFile } = require('child_process')
const { sendMessageSync, creatMessage, setMdContent } = require('./package-msg')

async function report(data) {
  var _data = []
  _data.push({
    key: '项目名称',
    h: 4,
    text: [
      {
        key: 'color',
        color: 'warning',
        text: '代码发布通知「coach小程序」'
      },
      '\n'
    ]
  })
  _data.push(
    {
      text: [
        '发布人：',
        {
          key: 'color',
          color: 'comment',
          text: data.author
        },
        '\n'
      ]
    },
    {
      text: [
        '发布分支：',
        {
          key: 'color',
          color: 'comment',
          text: data.branch
        },
        '\n'
      ]
    },
    {
      text: [
        '发布信息：',
        {
          key: 'color',
          color: 'comment',
          text: `${data.commit}【${data.message}】`
        },
        '\n'
      ]
    },
    {
      text: [
        '最近修改时间：',
        {
          key: 'color',
          color: 'comment',
          text: data.date
        },
        '\n'
      ]
    },
    {
      text: [
        '前往微信公众平台',
        {
          key: 'link',
          color: 'info',
          text: '点击跳转',
          url: data.wxlink
        },
        '\n'
      ]
    },
  )
  // 发送信息
  return sendMessageSync(
    creatMessage({
      msgtype: 'markdown',
      content: setMdContent(_data)
    }),
  )
}

function runGit(args) {
  return new Promise((resolve, reject) => {
    execFile('git', args, {
      cwd: __dirname,
      encoding: 'utf8'
    }, (error, stdout) => {
      if (error) {
        reject(error)
      } else {
        resolve(stdout)
      }
    })
  })
}

async function getCurrentBranch() {
  const stdout = await runGit(['rev-parse', '--abbrev-ref', 'HEAD'])
  return stdout.trim()
}

async function getGitLog() {
  const stdout = await runGit([
    'log',
    '-1',
    '--date=iso',
    '--pretty=format:%h%x00%aN <%aE>%x00%ad%x00%s'
  ])
  const logData = stdout.split('\0')

  if (logData.length !== 4) {
    throw new Error('无法解析 Git 提交信息')
  }

  return {
    commit: logData[0],
    author: logData[1],
    date: logData[2],
    message: logData[3]
  }
}

async function main() {
  try {
    const result = await Promise.all([getCurrentBranch(), getGitLog()])
    const branch = result[0]
    const logObj = result[1]

    // 如果是uat-tapestry或者prod-tapestry分支，发送消息
    if (branch === 'uat-tapestry' || branch === 'prod-tapestry') {
      await report(Object.assign({}, logObj, {
        branch: branch,
        wxlink: 'https://mp.weixin.qq.com/'
      }))
    } else {
      console.log(`当前分支为 ${branch}，不发送消息`)
    }
  } catch (error) {
    console.error('发布通知执行失败：', error)
    process.exitCode = 1
  }
}

main()
