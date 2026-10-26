// Component Object
import { toggleTab } from '@/packages/customer/utils'
import { StatusList } from '../../constants'
import { changeStatus, commentShareApi, verifyShareApi } from '@/api/guide-invite'

Component({
  properties: {
    item: {
      type: Object
    },
    index: {
      type: Number
    },
    role: {
      type: Number
    },
    from: {
      type: Number,
      value: 1
    }
  },
  data: {
    opened: false,
    statusList: StatusList,
    message: {
      show: false,
      title: '',
      content: ''
    },
    show: false,
    list: []
  },
  methods: {
    /**
     * 显示弹框
     * @param {*} e
     */
    handleShowMessage(e) {
      const { title, content } = e.currentTarget.dataset
      this.setData({
        message: {
          show: true,
          title,
          content
        }
      })
      toggleTab()
    },
    /**
     * 隐藏弹框
     * @param {*} e
     */
    handleHideMessage() {
      this.setData({
        'message.show': false
      })
      toggleTab()
    },
    /**
     * 复制
     */
    handleCopy() {
      wx.setClipboardData({
        data: this.data.item.reserve_number
      })
    },
    /**
     * 点击展开
     */
    handleToggleOpened() {
      this.setData({
        opened: !this.data.opened
      })
    },
    /**
     * 显示订单详情
     */
    handleShowOrder() {
      this.setData({
        show: true
      })
      toggleTab()
      // 神策数据上报 -- 查看订单记录
      wx.track.trackViewOrderHistory({
        entry_point: this.data.from
      })
    },
    /**
     * 查看详情
     */
    handleViewDetail() {
      const { item, index } = this.data
      wx.navigateTo({
        url: `/pages/store-invite/invite-detail/invite-detail?id=${item.id}&reserveNumber=${item.reserve_number}&role=${this.data.role}&index=${index}`,
        events: {
          // 为指定事件添加一个监听器，获取被打开页面传送到当前页面的数据
          acceptDataFromOpenedPage: res => {
            console.log('acceptDataFromOpenedPage', res)
            // 状态与当前不一致
            if (item.status !== res.status) {
              this.triggerEvent('statusChange', {index, item: res})
            }
          }
        }
      })
    },
    /**
     * 查看提示
     * @param {*} e
     */
    handleTooltip(e) {
      this.selectComponent('#mp-tooltip').show(this, `#${e.target.id}`)
    },
    handleClose() {
      this.setData({
        show: false
      })
      toggleTab()
    },
    /**
     * 改变状态
     * @param {*} e
     */
    onChangeStatus(e) {
      const { value } = e.currentTarget.dataset
      const { item } = this.data
      if (item.is_store) return
      const params = {
        id: item.id,
        reserve_number: item.reserve_number,
        is_store: Number(value)
      }
      changeStatus(params).then(res => {
        if (res.code === 200) {
          this.setData({
            'item.is_store': value
          })
        } else {
          // todo: 提示错误信息
          wx.showToast(res.message)
        }
      }).catch(err => {
        console.log(err)
      })
    },
    /**
     * 图片加载失败
     */
    handleError(e) {
      const { index, goodsIndex } = e.currentTarget.dataset
      this.setData({
        [`item.reserve_order_sales[${index}].reserve_order_goods[${goodsIndex}].goods_url`]: 'https://cdn-cos.coach.com.cn/uploads/2024/v4.2.0/icon_empty.png'
      })
    }
  }
})
