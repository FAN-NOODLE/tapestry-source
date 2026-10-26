
import {
  getDays
} from '@/utils/invite-date'

Component({
  properties: {
    timeObj: {
      type: Object
    },
    allDay: {
      type: Number,
      value: 30
    },
    timeSlot: {
      type: Array,
      value: []
    },
    showNum: {
      type: Boolean,
      value: true
    }
  },
  data: {
    dataList: [],
    currentIndex: 0,
    dateIndex: -1
  },

  /*    系统事件开始   */
  created () { },
  attached () { },
  ready () {
    console.log(this.properties.timeSlot)
    let startDate = new Date(new Date().getTime() + 1000 * 60 * 60 * 24)
    const dataList = getDays(startDate, this.properties.allDay, this.properties.timeSlot)
    const timeObj = this.properties.timeObj
    const keys = Object.keys(timeObj)
    dataList.forEach(item => {
      item.date_list.forEach(cell => {
        if (keys.includes(cell.startTime + '')) {
          cell.inviteNum = timeObj[cell.startTime]
        }
      })
    })
    this.setData({
      dataList
    })
  },
  methods: {
    /*    普通事件开始  */
    // 预览图片
    clickItem (e) {
      const currentIndex = e.currentTarget.dataset.index
      this.setData({
        currentIndex,
        dateIndex: -1
      })
    },
    clickInvite (e) {
      const dateIndex = e.currentTarget.dataset.index
      this.setData({
        dateIndex
      })
    },
    sumbit () {
      const { dateIndex, currentIndex, dataList } = this.data
      if (dateIndex === -1) {
        wx.showToast('请选择时间段')
        return
      }
      const value = dataList[currentIndex].date_list[dateIndex]
      // todo: 拼接时间，不改变原有数据结构
      const select_name = dataList[currentIndex].select_name
      const startHour = value.name.split(' - ')[0]
      const endHour = value.name.split(' - ')[1]
      const startTime = new Date(select_name + ' ' + startHour).getTime() / 1000
      const endTime = new Date(select_name + ' ' + endHour).getTime() / 1000
      const obj = {
        ...value,
        select_name,
        startTime,
        endTime
      }
      this.triggerEvent('sumbit', obj)
    },
    eventStop () { },
    closeModel () {
      this.triggerEvent('close')
    },
    /*    普通事件结束  */

    /*    cdp分割事件   */
    cdpSplitEvent () { }
    /*    cdp分割事件   */

    /*    行为事件开始   */

    /*    行为事件结束   */
  }
})
