/** 生成地图连接 */
const generateUrl = () => {
  const URL = 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.3.2/map.png'
  let getTimestamp = new Date().getTime()
  return URL + '?timestamp=' + getTimestamp
}

Component({
  data: {
    /** 是否展示、关闭预约信息弹窗 */
    show: false,
    imageUrl: generateUrl()
  },

  methods: {
    /** 关闭弹窗 */
    handleClose () {
      this.updatePopup(false)
    },

    /** 打开、关闭弹窗 */
    updatePopup (show = true) {
      this.setData({
        show
      })
    }
  }
})
