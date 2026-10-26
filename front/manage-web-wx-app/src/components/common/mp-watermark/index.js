import Base64 from './base64'
import { getMemberDetailApi } from '@/api/store-api'

const app = getApp()
Component({
  properties: {
    title: {
      type: String,
      value: '',
      observer() {
        this.init()
      }
    },
    opacity: {
      type: Number,
      value: 0.1
    },
    width: {
      type: Number,
      value: 200
    },
    height: {
      type: Number,
      value: 200
    },
    fontSize: {
      type: Number,
      value: 14
    },
    // 角度
    angle: {
      type: Number,
      value: -45
    }
  },
  data: {

  },
  lifetimes: {
    attached() {
      this.getMemberDetail()
    }
  },
  methods: {
    init() {
      const { title, opacity, width, height, fontSize, angle } = this.data
      const color = `rgba(0,0,0,${opacity})`
      const svg = `<?xml version="1.0" encoding="UTF-8"?><svg width="${width}" height="${height}" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"><text x="50%" y="50%" text-anchor="middle" font-size="${fontSize}"  font-weight="normal" fill="${color}" transform="rotate(${angle} ${width / 2}, ${height / 2})">${title}</text></svg>`

      this.setData({
        src: this.getSrc(svg)
      })
    },
    getSrc(rawData) {
      const base64 = Base64.encode(rawData)
      const src = 'data:image/svg+xml;base64,' + base64
      return src
    },
    getMemberDetail() {
      if (app.globalData.name) {
        const { name, userid } = app.globalData
        this.data.title = `${name || '导购姓名'}-${userid || '导购工号'}`
        this.init()
        return
      }

      getMemberDetailApi({
        id: 0
      })
        .then((res) => {
          const { name } = res.data
          app.globalData.name = res.data.name
          this.data.title = `${name || '导购姓名'}-${app.globalData.userid || '导购工号'}`
          this.init()
        }).catch((err) => {
          console.log(err, '获取会员信息失败')
        })
    }
  }
})
