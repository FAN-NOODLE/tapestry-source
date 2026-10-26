import { getCustomerPhoneApi } from '@/api/customer-api'
import { debounce } from '@/utils/tools'
import { openEnterpriseChat } from '@/utils/wx.qy'

// Component Object
Component({
  properties: {
    item: {
      type: Object
    },
    /** 来源: 1-顾客列表, 2-顾客分配, 3-顾客分配选择客户, 4-创建群发选择客户 */
    from: {
      type: Number,
      value: 1
    }
  },
  data: {
    imageHost: 'https://cdn-cos.coach.com.cn/uploads',
    isExpand: false
  },
  methods: {
    /**
     * @description 拨打电话
     */
    handleMakePhoneCall: debounce(async function(e) {
      const { item } = this.data
      console.log(item)

      const params = {
        unionid: item.unionid
      }
      if (!item.unionid) {
        wx.showToast('该用户暂无手机号')
        return
      }
      wx.showLoading({
        title: '',
        mask: true
      })
      console.log(params)
      try {
        const { code, data, message } = await getCustomerPhoneApi(params)
        wx.hideLoading()

        if (code !== 200) {
          wx.showToast(message)
          return
        }
        const { mobile } = data

        if (mobile) {
          wx.makePhoneCall({
            phoneNumber: mobile.replace(/\s*/g, ''),
            success: () => {
              // 神策数据上报 -- 拨打顾客电话
              wx.track.trackCallCustomer({
                user_external_id: item.openid
              })
              // 神策数据上报 -- 顾客列表页操作
              wx.track.trackCustomerOperation({
                current_page: '顾客列表页',
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
        // this.updateCrowdContact(item)
        // 神策数据上报 -- 顾客列表页操作
        wx.track.trackCustomerOperation({
          current_page: '顾客列表页',
          operation_type: '发消息按钮'
        })
      }).catch(err => {
        console.log(err)
      })
    }),
    /**
     * @description 查看详情
     */
    handleViewDetail(e) {
      const { openid, unionid, id } = e.currentTarget.dataset.item
      // 新链接
      let url = `/packages/customer/pages/detail/detail?id=${id}&openid=${openid}&unionid=${unionid}`
      wx.navigateTo({
        url: url
      })
    },
    /**
     * @description 显示tooltip
     */
    handleShowTooltip(e) {
      const { selector, content } = e.currentTarget.dataset
      if (content) {
        this.selectComponent('#mp-tooltip').show(this, selector, content)
      }
    },
    /**
     * @description 选择
     */
    handleCheck() {
      this.setData({
        'item.checked': !this.data.item.checked
      })
      this.triggerEvent('check', this.data.item)
    },
    /**
     * 展开切换
     * @param {*} e
     */
    handleToggleExpand () {
      this.setData({
        isExpand: !this.data.isExpand
      })
    }
  }
})

