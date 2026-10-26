
import { cancelReserveApi } from '@/api/booking-customer-api'
import { toDate } from '@/utils/tools'

const { globalData } = getApp()

Component({
  properties: {
    /** 预约客户 */
    info: { type: Object, observer: 'calcToday' }
  },
  data: {
    /** 是否展示、关闭预约信息弹窗 */
    show: false,
    /** 判断是否需要适配 */
    isIPhoneX: globalData.isIPhoneX,
    /** 今天日期 年月日周 */
    today: ''
  },
  created () {
    // 获取编辑预约组件实例
    this.toastPopup = this.selectComponent('#toast')
    this.appointmentMapPopup = this.selectComponent('#appointmentMap')
  },
  methods: {
    /** 打开展位地图 */
    jumpMap () {
      this.appointmentMapPopup.updatePopup()
    },

    /** 计算今天日期 */
    calcToday () {
      if (!this.data.info.reserve_time) return
      const weeks = ['日', '一', '二', '三', '四', '五', '六']
      const {nowYear, nowMonth, nowDate, nowDayOfWeek} = toDate(this.data.info.reserve_time)
      const today = `${nowYear}-${nowMonth}-${nowDate} 星期${weeks[nowDayOfWeek]}`

      this.setData({
        today
      })
    },

    /** 打开、关闭弹窗, 供外部调用 */
    updatePopup (show = true) {
      this.setData({
        show
      })
    },

    /** 关闭弹窗 */
    close () {
      this.updatePopup(false)
    },

    /** 修改预约信息 */
    updateInfo () {
      this.close()
      // 关闭弹窗
      // 携带参数判断是否打开修改弹窗
      this.triggerEvent('operationSuccess', true)
    },

    /** 取消预约 */
    cancelReserve () {
      let self = this
      wx.showModal({
        title: '提示',
        content: '确定取消当前活动预约？',
        success: function (res) {
          if (res.confirm) {
            const params = {
              activity_id: globalData.activity_id,
              reserve_id: self.data.info.id
            }
            cancelReserveApi(params)
              .then((res) => {
                // 错误判断
                if (res.code !== 200) return wx.showToast(res.message)

                self.toastPopup.openPopup('success', '取消成功')

                // 关闭弹窗
                setTimeout(() => {
                  self.close()
                  // 携带参数判断是否打开修改弹窗
                  self.triggerEvent('operationSuccess', false)
                }, 1000)
              }).catch((err) => {
                console.log(err, '获取活动信息失败')
              })
          } else {
            console.log('取消')
          }
        }
      })
    }
  }
})

