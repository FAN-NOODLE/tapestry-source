const { getCustomerPhoneApi } = require('@/api/customer-api')
const { debounce } = require('@/utils/tools')
const { openEnterpriseChat } = require('@/utils/wx.qy')

// Component Object
Component({
  properties: {
    item: {
      type: Object,
      observer: function(item) {
        const { selection } = this.data
        if (!item || !selection) return
        this.setData({
          selected: selection.isSelectedAll || selection.value.includes(item.openid)
        })
      }
    },
    selection: {
      type: Object,
      observer: function(selection) {
        if (!selection) return
        const { openid } = this.data.item
        this.setData({
          selected: selection.isSelectedAll || selection.value.includes(openid)
        })
      }
    }

  },
  data: {
    selected: false
  },
  methods: {
    handleSelect() {
      this.triggerEvent('select', {
        item: this.data.item,
        selected: this.data.selected
      })
    },
    /**
     * @description 电话联系
     */
    handleCall: debounce(async function (e) {
      wx.showLoading({ title: '', mask: true })

      try {
        const { item } = this.data
        if (!item.unionid) {
          wx.showToast('该用户暂无手机号')
          return
        }
        const params = {
          unionid: item.unionid
        }
        const { code, data, message } = await getCustomerPhoneApi(params)
        wx.hideLoading()
        if (code !== 200) {
          wx.showToast(message)
          return
        }
        if (data.mobile) {
          // 对某个人点击了打电话，则允许对该用户标记
          wx.makePhoneCall({
            phoneNumber: data.mobile.replace(/\s*/g, ''),
            success: (res) => {
              console.log(res, '成功')
              this.updateCrowdContact(item)
              // 神策数据上报 -- 顾客列表页操作
              wx.track.trackCommunicateOperation({
                current_page: '沟通名单页',
                operation_type: '电话按钮'
              })
            }
          })
        } else {
          wx.showToast('该用户暂无手机号')
        }
      } catch (error) {
        wx.hideLoading()
        console.log(error)
      }
    }),

    /**
     * @description 发送消息
     */
    handleSendMessage: debounce(async function (e) {
      const { item } = this.data
      if (item.self_relation_status !== 0) {
        wx.showToast('该用户非您好友，无法打开对话框')
        return
      }
      // 打开用户详情页
      openEnterpriseChat({
        externalUserIds: item.openid
      }).then(() => {
        this.updateCrowdContact(item)

        // 神策数据上报 -- 顾客列表页操作
        wx.track.trackCommunicateOperation({
          current_page: '沟通名单页',
          operation_type: '发消息按钮'
        })
      }).catch(err => {
        console.log(err)
      })
    }),

    /**
   * @description 标记已联系
   */
    handleContact(e) {
      const { item } = e.currentTarget.dataset
      this.updateCrowdContact(item)

      // 神策数据上报 -- 顾客列表页操作
      wx.track.trackCommunicateOperation({
        current_page: '沟通名单页',
        operation_type: '标记已联系'
      })
    },
    updateCrowdContact(item) {
      this.triggerEvent('contact', item)
    },

    /**
     * @description 查看详情
     */
    handleViewDetail() {
      const { openid, unionid, id } = this.data.item
      // 新链接
      let url = `/packages/customer/pages/detail/detail?id=${id}&openid=${openid}&unionid=${unionid}`
      wx.navigateTo({
        url: url
      })
    }
  }
})
