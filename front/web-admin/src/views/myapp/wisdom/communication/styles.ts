const fontFamily = {
  'zh-font': {
    font: 'SongTiSC',
  },
};
const colors = {
  red: {
    color: '#FF0D0D'
  },
  blue: {
    color: '#002FE9'
  },
}
const fontSizes = {
  fs10: {
    fontSize: 10
  },
  fs11: {
    fontSize: 11
  },
  fs12: {
    fontSize: 12
  },
  fs13: {
    fontSize: 13
  },
  fs14: {
    fontSize: 14
  },
  fs16: {
    fontSize: 16
  },
  fs24: {
    fontSize: 24
  },
}

const lines = {
  bold: {
    bold: true
  },
  underline: {
    decoration: 'underline'
  },
  lh1: {
    lineHeight: 1
  },
  lh12: {
    lineHeight: 1.2
  },
  lh15: {
    lineHeight: 1.5
  },
  lh2: {
    lineHeight: 2
  },
}

const layout = {
  center: {
    alignment: 'center',
  },
  right: {
    alignment: 'right',
  }
}

export default {
  // 字体
  ...fontFamily,
  // 颜色
  ...colors,
  // 字体大小
  ...fontSizes,
  // 行高
  ...lines,
  // 布局
  ...layout,
  // 通用
  'html-img': {
    margin: [0,0,0, 5]
  },
}