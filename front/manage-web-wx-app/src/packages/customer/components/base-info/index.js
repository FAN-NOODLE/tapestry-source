// 暂时注释掉导入以避免路径解析错误
import { workTagsOperateApi } from '@/api/customer-api'
import { getCustomerDetailByCrmApi, getCustomerInterestedProductsApi } from '@/api/customer-beta-api'
import { getBrandMaterialListApi } from '@/api/material-api'
import { customWhiteListApi } from '@/api/order'
import { toggleTab } from '../../utils'
import { update_time } from '@/utils/const'

function getBaseInfoItems() {
  // 简化初始化逻辑，移除条件判断以避免复杂性
  if (new Date() >= new Date(update_time)) {
    return [
      { label: '会员卡号', key: 'union_code', source: 'detail', isCopy: true, span: 24, defaultValue: '' },
      // { label: '会员姓名', key: 'name', source: 'detail', isCopy: true },
      { label: '会员生日', key: 'birthday', source: 'detail', span: 13, defaultValue: '' },
      { label: '居住地', key: 'city', source: 'crmDetail', span: 11, defaultValue: '' },
      { label: '首购店铺', key: 'vipStatus.issueShop', source: 'crmDetail', span: 13, defaultValue: '' },
      { label: '首购日期', key: 'vipStatus.issueDate', source: 'crmDetail', span: 11, defaultValue: '' },
      { label: '最近购买店铺', key: 'last_purchase_store_code', source: 'detail', span: 13, defaultValue: '' },
      { label: '会籍到期日', key: 'expiry_date', source: 'crmDetail', span: 11, defaultValue: '' },
      // { label: '注册渠道', key: 'first_channel_code', source: 'crmDetail', span: 13, defaultValue: '' },
      { label: '注册日期', key: 'registration_date', source: 'crmDetail', span: 13, defaultValue: '' }
    ]
  }
  return [
    { label: '会员卡号', key: 'union_code', source: 'detail', isCopy: true, span: 13, defaultValue: '' },
    { label: '会员姓名', key: 'name', source: 'detail', isCopy: true, span: 11 },
    { label: '会员生日', key: 'birthday', source: 'detail', span: 13 },
    { label: '当前可用积分', key: 'total_valid_points', source: 'crmDetail', span: 11 },
    { label: '发行店铺', key: 'vipStatus.issueShop', source: 'crmDetail', span: 13 },
    { label: '升级所需积分', key: 'upgrade_need_points', source: 'crmDetail', span: 11 },
    { label: '居住地', key: 'city', source: 'crmDetail', span: 13 },
    { label: '最近购买店铺', key: 'last_purchase_store_code', source: 'detail', span: 11 },
    { label: '会籍发行日', key: 'vipStatus.issueDate', source: 'crmDetail', span: 13 },
    { label: '最近升级时间', key: 'upgrade_date_time', source: 'crmDetail', span: 11 },
    { label: '会籍到期日', key: 'expiry_date', source: 'crmDetail', span: 13 }
  ]
}

// Component Object
Component({
  properties: {
    detail: {
      type: Object,
      observer(detail) {
        if (detail) {
          if (detail.chat_tag) {
            const index = detail.chat_tag.findIndex(item => item.is_select)
            console.log(index)

            this.setData({
              pickerOptions: detail.chat_tag.map(item => {
                return {
                  label: item.value,
                  value: item.key
                }
              }),
              value: index
            })
          }
          this.drawPie()
          this.getCustomerDetailByCrm()
          this.getInterestedProductsStatus()
        }
      }
    }
  },
  data: {
    crmDetail: null,
    isExpand: false,
    pickerOptions: [],
    value: -1,
    showPicker: false,
    products: [],
    group1: [],
    group2: [],
    group3: [],
    hideInterestedProducts: false,
    interestedProductsEmptyImage: 'https://cdn-cos.coach.com.cn/static/wxapp/V6.3.8/empty.png',
    show: true,
    // 基础信息展示项数组
    baseInfoItems: getBaseInfoItems(),
    update: new Date() >= new Date(update_time)
  },
  lifetimes: {
    attached: function () {
    },
    detached: function () {
      // 清理逻辑
    }
  },
  methods: {
    /**
     * 获取crm详情
     */
    async getCustomerDetailByCrm() {
      const { detail } = this.data
      const params = {
        unionid: detail.unionid
      }
      const res = await getCustomerDetailByCrmApi(params)
      if (res.code === 200) {
        this.setData({
          crmDetail: res.data || {}
        })
        this.triggerEvent('detail', res.data)
      }
    },
    async getInterestedProductsStatus() {
      try {
        const res = await customWhiteListApi({
          type: 6
        })
        const hideInterestedProducts = res.code === 200 && res.data && res.data.show === 1
        this.setData({
          hideInterestedProducts
        })
        if (!hideInterestedProducts) {
          this.getProducts()
        }
      } catch (error) {
        console.log(error)
        this.setData({
          hideInterestedProducts: false
        })
        this.getProducts()
      }
    },
    async getProducts() {
      const { detail } = this.data
      const params = {
        unionid: detail.unionid
      }
      const res = await getCustomerInterestedProductsApi(params)

      if (res.code === 200) {
        const mRes = await getBrandMaterialListApi({
          size: 15,
          upc_codes: res.data.map(item => item.upccode)
        })
        let group1 = []
        let group2 = []
        let group3 = []
        if (mRes.data.length) {
          const result = mRes.data.reduce((prev, cur) => {
            prev[cur.upc_code] = cur
            return prev
          }, {})
          console.log(result)

          group1 = res.data.filter(item => item.group === 1).map(item => {
            if (!result[item.upccode]) {
              return null
            }
            return {
              ...item,
              code: result[item.upccode].product_code,
              name: result[item.upccode].product_name,
              image: result[item.upccode].brand_source_list[0].url || ''
            }
          }).filter(Boolean)
          group2 = res.data.filter(item => item.group === 2).map(item => {
            if (!result[item.upccode]) {
              return null
            }
            return {
              ...item,
              code: result[item.upccode].product_code,
              name: result[item.upccode].product_name,
              image: result[item.upccode].brand_source_list[0].url || ''
            }
          }).filter(Boolean)
          group3 = res.data.filter(item => item.group === 3).map(item => {
            if (!result[item.upccode]) {
              return null
            }
            return {
              ...item,
              code: result[item.upccode].product_code,
              name: result[item.upccode].product_name,
              image: result[item.upccode].brand_source_list[0].url || ''
            }
          }).filter(Boolean)
        }


        this.setData({
          group1,
          group2,
          group3,
          products: res.data,
          showTips: res.data.length && res.data.every(item => item.type === 'hotsales')
        })
      }
    },
    /** 加载图片 */
    handleLoadImage(e) {
      const image_url = `https://cdn-cos.coach.com.cn/static/wxapp/common/icon_no_image.png?time=${new Date().getTime()}`

      const { callback } = e.detail
      const { item, index } = e.currentTarget.dataset
      if (item._loaded) {
        callback()
        return
      }
      console.log(item)
      getBrandMaterialListApi({
        upc_code: item.upccode
      }).then(res => {
        if (res.data && res.data.length) {
          const list = res.data[0].brand_source_list
          const current = list.find(image => image.is_main === 1) || list[0]

          this.setData({
            [`products[${index}].image`]: current ? current.url : image_url,
            [`products[${index}].name`]: res.data[0].product_name || '--'
          })
          setTimeout(() => {
            callback()
          }, 10)
        } else {
          this.setData({
            [`products[${index}].image`]: image_url,
            [`products[${index}].name`]: '',
            [`products[${index}].show`]: false
          })
          callback()
        }
      }).catch(err => {
        console.log(err)
        callback()
      })
    },
    /**
     * 绘制饼图
     */
    drawPie() {
      const { cdp_tag_buy_category_preference, cdp_tag_purchased_month } = this.data.detail

      wx.nextTick(() => {
        if (this.selectComponent('#pie1')) {
          this.selectComponent('#pie1').init(cdp_tag_buy_category_preference.map(item => {
            return {
              name: item.name,
              value: item.percent
            }
          }))
        }
        if (this.selectComponent('#pie2')) {
          this.selectComponent('#pie2').init(cdp_tag_purchased_month.map(item => {
            return {
              name: item.name + '月',
              value: item.percent
            }
          }))
        }
      })
    },

    hide() {
      if (this.selectComponent('#pie1')) {
        this.selectComponent('#pie1').init([])
      }
      if (this.selectComponent('#pie2')) {
        this.selectComponent('#pie2').init([])
      }
    },
    reload() {
      this.drawPie()
    },
    /**
     * 客户信息 - 展开收起
     */
    handleExpand() {
      this.setData({
        isExpand: !this.data.isExpand
      })
    },
    /**
     * 客户信息 - 复制会员卡号
     */
    handleCopy(e) {
      const {
        value
      } = e.currentTarget.dataset
      wx.setClipboardData({
        data: value
      })
    },
    /**
     * 客户信息 - 显示沟通意愿弹框
     */
    handleShowPicker() {
      if (!this.data.detail.is_relation) {
        wx.showToast('该用户非您好友，无法编辑')
        return
      }
      this.setData({
        showPicker: true
      })
      toggleTab()
    },
    /**
     * 客户信息 - 显示沟通意愿弹框 确认
     */
    async handlePickerConfirm(e) {
      const {
        index
      } = e.detail

      const { id, openid } = this.data.detail
      const { pickerOptions } = this.data

      const params = {
        external_userid: openid,
        work_external_user_id: id,
        add_tag_ids: [],
        del_tag_ids: []
      }

      pickerOptions.forEach((item, i) => {
        if (i === index) {
          params.add_tag_ids.push(item.value)
        } else {
          params.del_tag_ids.push(item.value)
        }
      })

      try {
        const { code, message } = await workTagsOperateApi(params)

        if (code !== 200) {
          wx.showToast(message)
        }

        this.setData({
          value: index,
          showPicker: false
        })
        toggleTab()
      } catch (error) {
        console.log(error)
      }
    },
    /**
     * 公用 - 关闭弹框
     */
    handleCancel() {
      this.setData({
        showPicker: false,
        showRemarkPopup: false,
        showDescPopup: false
      })
      toggleTab()
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

    handleViewDetail(e) {
      const { item, key } = e.currentTarget.dataset
      wx.navigateTo({
        url: `/pagesMaterial/pages/select/select?upc_code=${item.upccode}`
      })
      wx.track.trackCommodityClick({
        customer_id: this.data.detail.openid,
        commodity_name: item.name,
        product_UPC: item.upccode,
        product_code: item.code,
        mkt_id: item.image,
        mkt_rank: key
      })
    }
  }
})
