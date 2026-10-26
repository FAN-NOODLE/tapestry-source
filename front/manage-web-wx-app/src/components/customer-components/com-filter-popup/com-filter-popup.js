
// Component Object
import { getChoiceOptionApi } from '@/api/customer-api'
import { memberPositionApi } from '@/api/compass-api'
import dayjs from 'dayjs'

const defaultSelectList = [
  {
    type: 'select',
    title: '顾客沟通意愿',
    key: 'chat_tag',
    multiple: true
  },
  {
    type: 'select',
    title: '好友状态',
    key: 'is_union',
    multiple: false
  },
  {
    type: 'select',
    title: '顾客类型',
    key: 'age_tag',
    multiple: false
  },
  {
    type: 'select',
    title: 'CLT类型',
    key: 'clt_type',
    multiple: true,
    isExpand: true
  },
  {
    type: 'select',
    title: '近30天邀约状态',
    key: 'reserve_type',
    multiple: false,
    tips: [
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
    type: 'select',
    title: '近5天任务状况',
    key: 'task_have',
    multiple: false,
    tips: [
      {
        title: '一、无任务',
        content: '企业内，近5天无企微沟通任务关联匹配的顾客'
      },
      {
        title: '二、有任务',
        content: '企业内，近5天有企微沟通任务关联匹配的顾客'
      }
    ]
  },
  {
    type: 'select',
    title: '近30天任务触达',
    key: 'task_type',
    multiple: false,
    tips: [
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
    type: 'select',
    title: '优惠券筛选',
    key: 'coupon_type',
    multiple: false,
    tips: [
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
    type: 'select',
    title: '近1年消费金额',
    key: 'consume_type',
    multiple: true,
    tips: false,
    showExpand: true,
    isExpand: true
  },
  {
    type: 'date',
    title: '最近消费时间',
    key: 'recent_consume_time',
    keys: ['recent_consume_start_time', 'recent_consume_end_time'],
    tips: false,
    isExpand: true
  },
  {
    type: 'select',
    title: '顾客筛选',
    key: 'manage_type',
    multiple: false,
    isExpand: true
  },
  {
    type: 'select',
    title: '员工筛选',
    key: 'member_list',
    multiple: true,
    showExpand: true,
    isExpand: true
  },
  {
    type: 'select',
    title: '会员卡等',
    key: 'card_type',
    multiple: true,
    showExpand: true,
    isExpand: true,
    tips: [
      {
        title: '',
        content: '因会员卡等体系调整，2026年起，银卡/金卡/钻卡会员会籍到期后将更新为【已购会员】，会籍到期时间可在顾客详情中查看。\n\n如需筛选全部已经购买的会员可全选会员卡等中的选项。'
      }
    ]
  },
  {
    type: 'select',
    title: '客户印象',
    key: 'impression_tag_list',
    multiple: true,
    showExpand: true,
    isExpand: true
  },
  {
    type: 'select',
    title: '客户标签',
    key: 'tag_list',
    multiple: true,
    showExpand: true,
    isExpand: true
  }
]
const keys = {
  default: ['chat_tag', 'clt_type', 'reserve_type', 'task_have', 'task_type', 'coupon_type', 'consume_type', 'recent_consume_time', 'manage_type', 'member_list', 'card_type', 'age_tag', 'impression_tag_list', 'tag_list'],
  task: ['chat_tag', 'is_union', 'clt_type', 'reserve_type', 'task_have', 'task_type', 'coupon_type', 'consume_type', 'recent_consume_time', 'manage_type', 'card_type', 'age_tag', 'impression_tag_list', 'tag_list']
}
/**
 * 获取选择列表
 * @param {Array<String>} include
 * @returns
 */
function getSelectList(include = []) {
  if (include.length === 0) {
    return defaultSelectList
  } else {
    return include.map(item => {
      return defaultSelectList.find(i => i.key === item)
    })
  }
}

Component({
  properties: {
    value: {
      type: Object,
      value: {},
      observer: function (newVal) {
        if (newVal) {
          this.init()
        }
      }
    },
    count: {
      type: Number,
      value: 0
    },
    show: {
      type: Boolean,
      value: false,
      observer: function (newVal) {
        if (newVal) {
          this.init()
        }
      }
    },
    offset: {
      type: String,
      value: 0
    },
    zIndex: {
      type: Number,
      value: 5
    },
    include: {
      type: String,
      value: 'default'
    },
    currentPage: {
      type: Number,
      value: 4
    }
  },
  data: {
    /** 选择列表 */
    selectList: [],
    /** 提示信息显示控制 */
    showTip: false,
    /** 时间选择弹框显示控制 */
    showPickerDate: false,
    /** 时间选择配置 最大时间 */
    maxDate: new Date().getTime(),
    /** 时间选择配置 最小时间 */
    minDate: new Date().setFullYear(new Date().getFullYear() - 3),
    /** 筛选条件 */
    innerValue: {},
    /** 筛选条件数量 */
    innerCount: 0,
    /** 选择选项 */
    selectOption: {},
    /** 当前提示信息 */
    current: {
      title: '',
      list: []
    },
    /** 时间选择 */
    innerDateValue: [],
    /** 筛选项数据更新时间说明 */
    innerTips: {
      title: '筛选项数据更新时间说明',
      list: [
        {
          title: '一、数据变化后约1~5分钟后更新',
          content: '【好友状态】【客户印象】【客户标签】'
        },
        {
          title: '二、状态变化后约1小时后更新',
          content: '【顾客沟通意愿】【近30天邀约状态】\n【近5天任务状况】【近30天任务触达】'
        },
        {
          title: '三、数据变化后T+1日更新',
          content: '【CLT类型】【优惠券筛选】【近1年消费金额】\n【最近消费时间】【会员卡等】【顾客类型】'
        }
      ]
    }
  },
  lifetimes: {
    async attached() {
      await this.getSelectList()
      await this.getMember()
    },
    ready() {
    }
  },
  methods: {
    /**
     * @description 初始化
     */
    init() {
      const { include, role, value, count, selectOption } = this.data
      console.log(value)
      // 顾虑筛选列表
      const selectList = getSelectList(keys[include]).filter(item => {
        if (role === 5) {
          return true
        } else {
          return item.key !== 'member_list'
        }
      })
      // 设置选中状态
      const result = Object.keys(selectOption).reduce((prev, key) => {
        const { multiple, type } = selectList.find(item => item.key === key) || {}
        if (type === 'select' && selectOption[key]) {
          selectOption[key].forEach((val, index) => {
            if (multiple) {
              prev[`selectOption.${key}[${index}].selected`] = value[key] ? value[key].includes(val.key) : false
            } else {
              prev[`selectOption.${key}[${index}].selected`] = value[key] === val.key
            }
          })
        }
        return prev
      }, {})
      // const result = Object.entries(value).reduce((prev, [key, value]) => {
      //   const { multiple, type } = selectList.find(item => item.key === key) || {}
      //   if (type === 'select' && selectOption[key]) {
      //     selectOption[key].forEach((val, index) => {
      //       if (multiple) {
      //         prev[`selectOption.${key}[${index}].selected`] = value.includes(val.key)
      //       } else {
      //         prev[`selectOption.${key}[${index}].selected`] = val.key === value
      //       }
      //     })
      //   }
      //   return prev
      // }, {})
      console.log(result)
      this.setData({
        ...result,
        selectList,
        innerValue: {...value},
        innerCount: count
      })
      // wx.nextTick(() => {
      //   setTimeout(() => {
      //     this.computedDoubleLine()
      //   }, 2000)
      // })
    },
    /**
     * @description 获取筛选数据
     */
    async getSelectList () {
      try {
        const { code, data, message } = await getChoiceOptionApi()
        if (code !== 200) {
          wx.showToast(message)
          return
        }
        this.setData({ selectOption: data })
      } catch (error) {
        console.log(error)
      }
    },
    /**
     * @description 获取身份信息
     */
    async getMember () {
      await memberPositionApi({ is_new: 0 }).then(res => {
        if (res.code === 200) {
          // role 3 企业导购 5 企业管理（目前只有3和5）
          this.setData({
            role: res.data.role
          })
        } else {
          wx.hideLoading()
        }
      }).catch(err => {
        wx.hideLoading()
        console.log(err)
      })
    },
    /**
     * @description 提示弹框
     */
    handleShowTips(e) {
      const { index } = e.currentTarget.dataset
      const { selectList, innerTips } = this.data
      if (index === -1) {
        this.setData({
          current: innerTips,
          showTip: true
        })
      } else {
        this.setData({
          current: {
            title: selectList[index].title,
            list: selectList[index].tips
          },
          showTip: true
        })
      }
    },
    /**
     * @description 关闭提示弹框
     */
    handleHideTips() {
      this.setData({
        showTip: false
      })
    },
    /**
     * @description 选择事件
     * @param {*} e
     */
    handleSelect(e) {
      const { key, index } = e.currentTarget.dataset
      const { selectList, selectOption } = this.data
      const item = selectList.find(i => i.key === key)

      // 设置选中项
      const result = {
        [`selectOption.${key}[${index}].selected`]: !selectOption[key][index].selected
      }
      // 如果是单选取消其他选中项
      if (!item.multiple) {
        selectOption[key].forEach((item, i) => {
          if (i !== index && item.selected) {
            result[`selectOption.${key}[${i}].selected`] = false
          }
        })
      }
      this.setData(result, () => {
        this.updateInnerValue()
      })
    },
    /**
     * @description 更新筛选条件
     */
    updateInnerValue() {
      const { selectOption, selectList, innerValue, innerCount } = this.data

      const {result, count} = selectList.reduce((prev, item) => {
        // 选择列表
        if (item.type === 'select') {
          const value = selectOption[item.key].filter(item => item.selected).map(item => item.key)
          // 有选中数据
          if (value.length > 0) {
            // 筛选数量变更
            if (!prev.result[item.key]) {
              prev.count++
            }
            // 多选与单选赋值
            prev.result[item.key] = item.multiple ? value : value[0]
          } else {
            // 筛选数量变更
            if (prev.result[item.key]) {
              prev.count--
            }
            delete prev.result[item.key]
          }
        }
        return prev
      }, {result: innerValue, count: innerCount})

      this.setData({
        innerValue: result,
        innerCount: count
      })
    },
    /**
     * @description 选择时间
     */
    handlePickerShow(e) {
      const { key } = e.currentTarget.dataset
      const { selectList, innerValue, innerDateValue } = this.data
      const { keys } = selectList.find(item => item.key === key)
      if (innerValue[keys[0]]) {
        innerDateValue[0] = innerValue[keys[0]] * 1000
      }
      if (innerValue[keys[1]]) {
        innerDateValue[1] = innerValue[keys[1]] * 1000
      }
      this.setData({
        showPickerDate: true,
        innerDateValue: innerDateValue
      })
    },
    /**
     * @description 选择时间确定
     * @param {*} e
     */
    handlePickerConfirm(e) {
      const { value } = e.detail
      const [ start, end ] = value
      let { innerValue, innerCount } = this.data

      if (start || end) {
        if (!innerValue.recent_consume_start_time) {
          innerCount++
        }
      } else {
        if (innerValue.recent_consume_start_time) {
          innerCount--
        }
      }

      innerValue.recent_consume_start_time = start ? parseInt(start / 1000) : 0
      innerValue.recent_consume_end_time = end ? parseInt(end / 1000) : 0
      this.setData({
        innerValue,
        innerCount,
        innerDateValue: value
      })

      // 关闭弹框
      this.handlePickerHide()
    },
    /**
     * @description 选择时间取消
     */
    handlePickerHide() {
      this.setData({
        showPickerDate: false
      })
    },
    /**
     * @description 关闭弹框
     */
    handleClose() {
      // todo: 将展开项重置
      // this.handleReset()
      this.triggerEvent('close')
    },
    /**
     * @description 重置筛选
     */
    handleReset() {
      const { selectOption, selectList } = this.data

      const result = selectList.reduce((prev, item) => {
        // 选择列表
        if (item.type === 'select') {
          selectOption[item.key].forEach((val, index) => {
            prev[`selectOption.${item.key}[${index}].selected`] = false
          })
        }
        // todo: 将展开项重置
        // item.isExpand = true
        // item.showExpand = false
        return prev
      }, {})
      // todo: 将展开项重置
      this.setData({
        ...result,
        innerValue: {},
        innerCount: 0,
        innerDateValue: []
      }, () => {
        const { innerValue, innerCount } = this.data
        this.triggerEvent('confirm', {
          value: innerValue,
          count: innerCount
        })
      })
      this.track('重置')
    },
    /**
     * @description 确认选择
     */
    handleConfirm() {
      const { innerValue, innerCount } = this.data
      this.triggerEvent('confirm', {
        value: innerValue,
        count: innerCount
      })

      this.track('确定')
    },
    computedDoubleLine () {
      var query = wx.createSelectorQuery().in(this)
      query.selectAll('.selection__content').boundingClientRect((rects) => {
        const result = rects.reduce((prev, item) => {
          const { height, dataset } = item
          prev[`selectList[${dataset.index}].showExpand`] = height > 39 * 1.5
          return prev
        }, {})
        this.setData(result)
      }).exec()
    },
    handleExpand(e) {
      const { index } = e.currentTarget.dataset
      this.setData({
        [`selectList[${index}].isExpand`]: !this.data.selectList[index].isExpand
      })
    },

    track(button_name) {
      const { innerValue, currentPage, selectList, selectOption } = this.data

      const result = []

      const { recent_consume_start_time, recent_consume_end_time, ...data } = innerValue

      Object.keys(data).forEach(key => {
        const item = selectList.find(item => item.key === key)
        const value = data[key]
        if (item.type === 'select') {
          // 选择
          if (item.multiple) {
            // 多选
            selectOption[key].forEach(val => {
              if (value.includes(val.key)) {
                result.push({
                  tag_level_1: '',
                  tag_level_2: item.title,
                  tag_value: val.value
                })
              }
            })
          } else {
            // 单选
            selectOption[key].forEach(val => {
              if (value === val.key) {
                result.push({
                  tag_level_1: '',
                  tag_level_2: item.title,
                  tag_value: val.value
                })
              }
            })
          }
        }
      })

      if (recent_consume_start_time && recent_consume_end_time) {
        result.push({
          tag_level_1: '',
          tag_level_2: '最近消费时间',
          tag_value: `${dayjs(recent_consume_start_time * 1000).format('YYYY/MM/DD')}-${dayjs(recent_consume_end_time * 1000).format('YYYY/MM/DD')}`
        })
      }

      const result_option = result.reduce((prev, cur) => {
        if (prev[cur.tag_level_2]) {
          prev[cur.tag_level_2].push(cur.tag_value)
        } else {
          const item = selectList.find(item => item.title === cur.tag_level_2)
          prev[cur.tag_level_2] = item.multiple ? [cur.tag_value] : cur.tag_value
        }
        return prev
      }, {})

      // 神策数据上报 -- 筛选
      wx.track.trackFilter({
        current_page: currentPage,
        button_name: button_name,
        filter_type_sv: 1,
        filter_option: result_option
      })

      result.forEach(item => {
      // 神策数据上报 -- 筛选详情
        wx.track.trackFilterDetails({
          current_page: currentPage,
          button_name: button_name,
          filter_type_sv: 1,
          ...item
        })
      })
    }

  }
})

