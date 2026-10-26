import { impressionTagsOperateApi, workTagsOperateApi } from '@/api/customer-api'

Component({
  properties: {
    info: {
      type: Object,
      value: {}
    },
    showOtherDetail: {
      type: Boolean,
      value: true
    },
    showModulePage: {
      type: Boolean,
      value: true
    },
    openid: {
      type: String,
      value: ''
    }
  },
  data: {
    isOpenTipsModal: false,
    description: '',
    dec1: '客户标签仅提供推荐建议，导购无法编辑客户标签。',
    dec2: `
    为顾客添加或移除印象标签
    1、点击客户印象标签可以为顾客添加或移除印象标签
    2、已添加的标签以深色底色展示
    3、所有导购都可以为顾客添加或移除印象标签，页面展示最后一次编辑
    `,
    isLoading: false
  },

  /*    系统事件开始   */
  created () {},
  attached () {},
  ready () {},
  /*    系统事件结束   */

  methods: {
    /*    普通事件开始  */

    /*    普通事件结束  */

    /*    cdp分割事件   */
    cdpSplitEvent () {},
    /*    cdp分割事件   */

    /*    行为事件开始   */
    /**
     * 去购物清单页面
     */
    goShoppingList () {
      console.log('去购物清单页面')
      const url = `/pages/order-list/order-list?info=${JSON.stringify(this.properties.info)}`
      wx.navigateTo({ url })
    },
    /**
     * 去卡券详情页面
     */
    goCardDetail () {
      console.log('去卡券详情页面')
      const url = `/pages/coupon-list/coupon-list?info=${JSON.stringify(this.properties.info)}`
      wx.navigateTo({ url })
    },
    /**
     * 打开提示语
     */
    openTipsModel (e) {
      const { dec } = e.currentTarget.dataset

      this.setData({ description: dec, isOpenTipsModal: true })
    },
    /**
     * 关闭提示语
     */
    closeTipsModel () {
      this.setData({ isOpenTipsModal: false })
    },

    async tagsImpressionOperate (e) {
      const { index } = e.currentTarget.dataset
      const { tag, tagarr: tagArr } = e.detail
      if (!this.properties.info.isRelation) {
        return wx.showToast('抱歉，非本人好友暂不可操作')
      }
      if (!this.properties.info.customer_code) {
        wx.showToast('顾客还未注册会员，无法使用标签功能')
        return
      }
      if (this.data.isLoading) return
      let add_tag_ids = []
      const del_tag_ids = []

      tagArr.forEach((item) => {
        if (item.is_have) {
          del_tag_ids.push(item.tag_id)
        }
      })
      if (!tag.is_have) {
        add_tag_ids = [tag.tag_id]
      }

      this.setData({ isLoading: true })
      const params = {
        external_userid: this.properties.openid,
        work_external_user_id: this.properties.info.work_external_user_id,
        add_tag_ids: add_tag_ids,
        del_tag_ids: del_tag_ids
      }

      try {
        const { code, message } = await impressionTagsOperateApi(params)

        if (code !== 200) {
          this.setData({ isLoading: false })
          wx.showToast(message)
          return
        }

        let newTags = tagArr.map(item => {
          return {...item, is_have: item.tag_id === tag.tag_id && !tag.is_have ? 1 : 0}
        })

        this.data.info.user_impression_tags[index] = {
          ...this.data.info.user_impression_tags[index],
          impression_tags_list: newTags
        }

        this.setData({
          info: this.data.info
        }, () => {
          this.setData({ isLoading: false })
        })

        // this.triggerEvent('reGetDetail')
      } catch (error) {
        this.setData({ isLoading: false })
        console.log(error)
      }
    },

    async handleChatTagChange(e) {
      if (!this.data.info.isRelation) {
        return wx.showToast('抱歉，非本人好友暂不可操作')
      }
      const { tag } = e.detail
      const { impression_tags_list } = this.data.info.user_chat_tags
      const add_tag_ids = []
      const del_tag_ids = []

      const isDoubleHave = impression_tags_list.every(item => item.is_have)

      if (isDoubleHave) {
        del_tag_ids.push(tag.tag_id)
      } else {
        // 遍历当前选中的标签
        impression_tags_list.forEach((item) => {
          if (item.is_have) {
            del_tag_ids.push(item.tag_id)
          }
        })
        // 记录选中标签
        if (!tag.is_have) {
          add_tag_ids.push(tag.tag_id)
        }
      }

      this.setData({ isLoading: true })
      const params = {
        external_userid: this.data.openid,
        work_external_user_id: this.data.info.work_external_user_id,
        add_tag_ids,
        del_tag_ids
      }

      try {
        const { code, message } = await workTagsOperateApi(params)

        if (code !== 200) {
          this.setData({ isLoading: false })
          wx.showToast(message)
          return
        }

        const newArr = impression_tags_list.map(item => {
          if (isDoubleHave) {
            if (tag.tag_id === item.tag_id) {
              item.is_have = 0
            }
          } else {
            if (tag.is_have && tag.tag_id === item.tag_id) {
              item.is_have = 0
            } else {
              item.is_have = tag.tag_id === item.tag_id ? 1 : 0
            }
          }
          return item
        })
        this.setData({
          'info.user_chat_tags.impression_tags_list': newArr,
          isLoading: false
        })

        // this.triggerEvent('reGetDetail')
      } catch (error) {
        this.setData({ isLoading: false })
        console.log(error)
      }
    }
    /*    行为事件结束   */
  }
})
