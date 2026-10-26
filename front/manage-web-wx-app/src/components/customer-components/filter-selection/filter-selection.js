import { getChoiceOptionApi } from '@/api/customer-api'
import { memberPositionApi } from '@/api/compass-api'

const defaultSelectList = [
  {
    title: '近30天邀约状态',
    key: 'reserve_type',
    multiple: false,
    list: [],
    tips: true,
    tipList: [
      {
        title: '一、邀约中',
        content: '企业内，近30天存在进行中的邀约'
      },
      {
        title: '二、邀约成功',
        content: '企业内，近30天邀约成功的顾客'
      },
      {
        title: '三、无邀约',
        content: '企业内，近30天无邀约记录的顾客'
      }
    ]
  },
  {
    title: '近30天任务状况',
    key: 'task_type',
    multiple: false,
    list: [],
    tips: true,
    tipList: [
      {
        title: '一、已触达',
        content: '企业内，近30天通过任务沟通被触达的顾客'
      },
      {
        title: '二、未触达',
        content: '企业内，近30天未出现在沟通任务名单中的顾客，或出现在任务名单中但未被触达的顾客'
      }
    ]
  },
  {
    title: '优惠券筛选',
    key: 'coupon_type',
    multiple: false,
    list: [],
    tips: true,
    tipList: [
      {
        title: '一、30天内临期券',
        content: '持有有效期少于30天的优惠券的顾客'
      },
      {
        title: '二、待使用券',
        content: '持有待使用优惠券的顾客'
      }
    ]
  },
  {
    title: '近1年消费金额',
    key: 'consume_type',
    multiple: false,
    list: [],
    tips: false
  },
  {
    title: '最近消费时间',
    key: 'recent_consume_time',
    multiple: true,
    list: [],
    tips: false,
    type: 'date',
    show: true
  },
  {
    title: '顾客筛选',
    key: 'manage_type',
    multiple: false,
    list: []
  },
  {
    title: '员工筛选',
    key: 'member_list',
    multiple: true,
    list: []
  },
  {
    title: '会员卡等',
    key: 'card_type',
    multiple: true,
    list: []
  },
  {
    title: '印象标签',
    key: 'impression_tag_list',
    multiple: true,
    list: []
  },
  {
    title: '客户标签',
    key: 'tag_list',
    multiple: true,
    list: []
  },
  // {
  //   title: 'CLT类型',
  //   key: 'clt_type',
  //   multiple: true,
  //   list: []
  // }
]

Component({
  // 使用全局 & 页面样式
  options: {
    addGlobalClass: true
  },
  properties: {
    result: {
      type: Object,
      value: {}
    }
  },
  data: {
    selectList: defaultSelectList.filter(i => i.key !== 'member_list'),
    count: 0
  },
  /*    系统事件开始   */
  created () {
  },
  attached () {},
  ready () {
    this.getMember()
  },
  /*    系统事件结束   */

  methods: {
    /**
     * 获取筛选数据
     */
    async getSelectList () {
      try {
        const { code, data, message } = await getChoiceOptionApi()
        wx.hideLoading()
        if (code !== 200) {
          wx.showToast(message)
          return
        }
        const { selectList } = this.data
        const newSelectList = selectList.map(i => ({
          ...i,
          list: data[i.key] || [],
          show: (data[i.key] || []).length > 0 || i.show
        }))

        this.setData({ selectList: newSelectList })
      } catch (error) {
        wx.hideLoading()
        console.log(error)
      }
    },
    // 重置数据
    resetData () {
      const { selectList } = this.data
      const result = {}

      selectList.forEach(v => {
        this.selectComponent(`#${v.key}`).setData({ result: [] })
      })
      console.log('reset data', result)
      this.triggerEvent('confirm', result)
    },
    // 确定提交筛选数据
    confirm () {
      const { selectList } = this.data
      const result = {}

      selectList.forEach(v => {
        const itemResult = this.selectComponent(`#${v.key}`).data.result

        if (itemResult.length > 0) {
          result[v.key] = v.multiple ? itemResult : itemResult[0]
        }
      })

      console.log('confirm', result)
      this.triggerEvent('confirm', result)
    },
    /**
     * 获取身份信息
     */
    async getMember () {
      wx.showLoading({
        title: '加载中'
      })
      await memberPositionApi({ is_new: 0 }).then(res => {
        if (res.code === 200) {
          // role 3 企业导购 5 企业管理（目前只有3和5）
          this.setData({
            role: res.data.role
          })

          if (res.data.role === 5) {
            this.setData({ selectList: defaultSelectList })
          } else {
            this.setData({ selectList: defaultSelectList.filter(i => i.key !== 'member_list') })
          }

          this.getSelectList()
        } else {
          wx.hideLoading()
        }
      }).catch(err => {
        wx.hideLoading()
        console.log(err)
      })
    },
    /**
     * 显示提示
     * @param {*} filterItem
     * @returns
     */
    showTips (e) {
      this.triggerEvent('showTips', e.detail)
    },
    handleSelectDate(e) {
      console.log(e.detail)
      this.triggerEvent('select-date', e.detail)
    },
    /**
     * @description 选项变化
     * @param {*} e
     */
    handleChange(e) {
      const { value, item } = e.detail
      let cacheSet = new Set(this.cacheList || [])
      if (value.length > 0) {
        cacheSet.add(item.key)
      } else {
        cacheSet.delete(item.key)
      }
      this.cacheList = Array.from(cacheSet)
      this.setData({ count: this.cacheList.length })
      this.triggerEvent('change')
    }
  }
})
