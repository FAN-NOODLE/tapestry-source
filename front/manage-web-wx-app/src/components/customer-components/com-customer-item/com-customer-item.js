import { getCustomerPhoneApi } from '@/api/customer-api'
import { debounce } from '@/utils/tools'

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
    },
    isBeta: {
      type: Boolean,
      value: false
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
      const { item } = e.currentTarget.dataset
      const params = {
        unionid: item.unionid,
        customer_code: item.customer_id
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
                current_page: ['顾客列表', '顾客分配', '顾客分配选择客户', '创建群发选择客户'][this.data.from - 1],
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
     * @description 查看详情
     */
    handleViewDetail(e) {
      const { openid, sa_user_id, member_id, is_relation, customer_code } = e.currentTarget.dataset.item
      let url = `/pages/customer-detail/customer-detail?openid=${openid}&memberUserid=${sa_user_id}&memberId=${member_id}&isRelation=${is_relation}&customer_id=${customer_code}`
      // beta
      if (this.data.isBeta) {
        url = `/packages/customer/pages/detail/detail?openid=${openid}`
      }
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

