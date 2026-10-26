// Component Object
import { addCrowdApi, gerCrowdListApi, getChoiceOptionApi, getCustomerListApi } from '@/api/customer-beta-api'
import { options } from './data'
import dayjs from 'dayjs'
import { debounce } from '@/utils/tools'
import { update_time } from '@/utils/const'

const { globalData } = getApp()
Component({
  properties: {
    value: {
      type: Object,
      value: {},
      observer: function() {}
    },
    show: {
      type: Boolean,
      value: false,
      observer: function(val) {
        if (val) {
          // this.getChoiceOption()
          // this.gerCrowdList()
          const app = getApp()
          this.setData({
            current: 0,
            scrollToView: `group-${0}`,
            innerValue: {
              rel_create_time_type: app.globalData.is_manager === 1 ? 2 : 1,
              rel_interact_time_type: 2,
              rel_active_time_type: 2,
              ...JSON.parse(JSON.stringify(this.data.value))
            }
          })
          this.setInnerValueMap()
        }
      }
    },
    offset: {
      type: String,
      value: globalData.system.statusBarHeight + 44 + 58 + 'px'
    },
    zIndex: {
      type: Number,
      value: 9
    },
    onlySearch: {
      type: Boolean,
      value: false
    },
    searchButtonText: {
      type: String,
      value: '确定'
    },
    exclude: {
      type: Array,
      value: [],
      observer(exclude) {
        if (exclude.length > 0) {
          this.setData({
            options: options.filter(item => !exclude.includes(item.key))
          })
        }
      }
    },
    currentPage: {
      type: Number,
      value: 4
    }
  },
  data: {
    /** 分组 */
    group: [
      {
        value: 1,
        label: '基本信息'
      },
      {
        value: 2,
        label: '邀约任务'
      },
      {
        value: 3,
        label: '购物行为'
      },
      {
        value: 4,
        label: '沟通互动'
      },
      {
        value: 5,
        label: '偏好信息'
      }
    ],
    /** 当前分组下标 */
    current: 0,
    scrollToView: '',
    /** 页面展示的选项 */
    options,
    /** 底部提示标题 */
    innerTitle: '筛选项数据更新时间说明',
    /** 底部提示内容 */
    innerTips: [
      {
        title: '一、数据变化后约1~5分钟后更新',
        content: '【好友状态】【客户印象】【客户标签】'
      },
      {
        title: '二、状态变化后约1小时后更新',
        content: '【顾客沟通意愿】【近30天邀约状态】\n【近30天任务状况】【近30天任务触达】'
      },
      {
        title: '三、数据变化后T+1日更新',
        content: '【优惠券筛选】【近1年消费金额】【最近消费时间】\n【会员卡等】【顾客类型】'
      }
    ],
    /** 提示信息 */
    tips: {
      show: false,
      title: '',
      list: []
    },
    /** 内部值 */
    innerValue: {
      rel_create_time_type: 1,
      rel_interact_time_type: 1,
      rel_active_time_type: 2
    },
    /** 内部值关系映射 */
    innerValueMap: [],
    /** 打包弹框 */
    showSavePopup: false,
    /** 单个包最多2000 */
    maxCount: 2000,
    /** 已查询总人数 */
    total: 0,
    /** 单包人数限制弹框 */
    showCountMessagePopup: false,
    /** 最多可以打多少个包 */
    maxPackCount: 5,
    /** 已经有多少个包了 */
    packTotal: 0,
    /** 包数弹框限制 */
    showPackCountMessagePopup: false,
    /** 创建成功提示 */
    showPackSuccessPopup: false,

    innerDateValue: [],
    /** 时间选择弹框显示控制 */
    showPickerDate: false,
    /** 时间选择配置 最大时间 */
    maxDate: new Date().getTime(),
    /** 时间选择配置 最小时间 */
    minDate: new Date().setFullYear(new Date().getFullYear() - 3),
    /** 获取数据中 */
    loading: false
  },

  lifetimes: {
    attached: function () {
      this.disabledScroll = true

      this.getChoiceOption()
      this.gerCrowdList()
    },
    detached: function () {

    }
  },
  pageLifetimes: {
    // 组件所在页面的生命周期函数
    show() {
      this.gerCrowdList()
    }
  },
  methods: {
    async getOptions() {
      return new Promise(resolve => {
        wx.request({
          url: 'https://cdn-cos.coach.com.cn/static/wxapp/json/data.json?timestamp=' + Date.now(),
          method: 'GET',
          header: {
            'Content-Type': 'application/json'
          },
          success: function (res) {
            resolve(res.data)
          },
          fail: function () {
            resolve(options)
          }
        })
      })
    },
    async getChoiceOption() {
      const optionsJSON = await this.getOptions()
      const res = await getChoiceOptionApi()
      const { member_list, ...value } = res.data
      // 获取group
      const group = Object.entries(optionsJSON.reduce((prev, cur) => {
        prev[cur.group] = cur.group_name
        return prev
      }, {})).map(item => ({
        value: Number(item[0]),
        label: item[1]
      }))
      console.log('分组', group)

      // 合并选项数据
      const options = optionsJSON.map(item => {
        if (item.type === 1 || item.type === 2) {
          item.options = res.data[item.alias || item.key] || []
          item.hide = item.options.length === 0
          item.span = item.span || (item.options.length === 5 ? 4.8 : 6)
        }
        if (item.key === 'card_type') {
          if (new Date() >= new Date(update_time) && new Date() <= new Date('2026-12-31 23:59:59')) {
            item.tips = [
              {
                title: '',
                content: '因会员卡等体系调整，2026年起，银卡/金卡/钻卡会员会籍到期后将更新为【已购会员】，会籍到期时间可在顾客详情中查看。\n\n如需筛选全部已经购买的会员可全选Retention中的选项。'
              }
            ]
          }
        }

        if (item.key === 'task_type' && res.data.task_type_value) {
          item = JSON.parse(JSON.stringify(item).replaceAll('30', res.data.task_type_value))
          this.data.innerTips[1].content = this.data.innerTips[1].content.replace('近30天任务触达', `近${res.data.task_type_value}天任务触达`)
        }

        if (item.key === 'task_have' && res.data.task_have_value) {
          item = JSON.parse(JSON.stringify(item).replaceAll('30', res.data.task_have_value))
          this.data.innerTips[1].content = this.data.innerTips[1].content.replace('近30天任务状况', `近${res.data.task_have_value}天任务状况`)
        }

        return item
      })
      console.log(options)


      this.setData({
        group,
        options,
        innerTips: this.data.innerTips
      })
    },
    async gerCrowdList() {
      const res = await gerCrowdListApi()
      console.log(res)
      this.setData({
        packTotal: res.data.length || 0,
        packList: res.data.length ? res.data : []
      })
    },

    getCustomerList: debounce(async function() {
      const res = await getCustomerListApi({
        only_count: true,
        is_search: this.data.innerValueMap.length > 0 ? 1 : 0,
        ...this.data.innerValue
      })

      if (res.code === 200) {
        this.setData({
          total: res.data.total,
          loading: false
        })
      } else {
        this.setData({
          loading: false
        })
      }
    }, 1000),
    /**
     * 打开筛选弹框
     */
    handleShow() {
      this.setData({
        show: true
      })
    },
    /**
     * 关闭筛选弹框
     */
    handleClose() {
      this.setData({
        show: false
      })
      this.triggerEvent('close')
    },
    /**
     * 侧边栏分组选择
     * @param {*} e
     */
    handleSideChange(e) {
      const { index } = e.currentTarget.dataset

      this.setData({
        current: index,
        scrollToView: `group-${index}`
      })
      // 关闭滚动监听
      this.disabledScroll = true
    },
    /**
     * 选项变更
     * @param {*} e
     */
    handleFilterChange(e) {
      const { value, key, type } = e.detail
      const { innerValue } = this.data
      console.log(value, key, type)
      if (type === 4) {
        if (value.start_number === '') delete value.start_number
        if (value.end_number === '') delete value.end_number
      }

      if (!value || JSON.stringify(value) === '{}') {
        delete innerValue[key]
        this.setData({
          innerValue
        })
      } else {
        this.setData({
          [`innerValue.${key}`]: value
        })
      }

      this.setInnerValueMap()
      this.setData({
        loading: true
      })
      this.getCustomerList()
    },
    setInnerValueMap() {
      const { innerValue, options } = this.data

      const innerValueMap = []
      Object.entries(innerValue).forEach(([key, value]) => {
        const current = options.find(item => item.key === key)
        if (current && current.options) {
          const result = current.options.filter(item => {
            if (current.type === 1) return item.key === value
            return value.includes(item.key)
          })
          if (result.length) {
            innerValueMap.push({
              key: current.key,
              type: current.type,
              title: current.label,
              result
            })
          }
        } else if (current) {
          innerValueMap.push({
            key: current.key,
            type: current.type,
            title: current.label,
            result: value
          })
        }
      })

      this.setData({
        innerValueMap
      })
    },
    /**
     * 个人/企业开关切换
     * @param {*} e
     */
    handleSwitchChange(e) {
      const { value, key } = e.currentTarget.dataset
      console.log(value, key)
      this.setData({
        [`innerValue.${key}`]: Number(value)
      })
      this.setData({
        loading: true
      })
      this.getCustomerList()
    },
    /**
     * 提示弹框
     * @param {*} e
     * @returns
     */
    handleTips(e) {
      if (e.type === 'close') {
        this.setData({
          'tips.show': false
        })
        return
      }
      const { tips, label } = e.detail

      this.setData({
        'tips.show': !this.data.tips.show,
        'tips.title': label || this.data.innerTitle,
        'tips.list': tips || this.data.innerTips
      })
    },
    /**
     * 打包沟通名单
     */
    handleSave() {
      // 神策数据上报 -- 顾客列表页操作
      wx.track.trackCustomerOperation({
        current_page: '筛选',
        operation_type: '打包沟通名单'
      })

      this.track('打包沟通名单')

      if (this.data.innerValueMap.length === 0) {
        wx.showToast('请选择筛选条件')
        return
      }
      if (this.data.total === 0) {
        wx.showToast('当前打包人数为0，请重新选择')
        return
      }

      const { maxCount, total, maxPackCount, packTotal } = this.data

      // 打包人数拦截 单包最大2000人
      if (total > maxCount) {
        this.setData({
          showCountMessagePopup: true
        })
        return
      }
      // 打包数量拦截 最多可以包5个包
      if (packTotal >= maxPackCount) {
        this.setData({
          showPackCountMessagePopup: true
        })
        return
      }

      this.setData({
        showSavePopup: true,
        packName: dayjs().format('YYYYMMDD')
      })
    },
    /**
     * 打包沟通名单 取消
     */
    handleSaveCancel() {
      // 神策数据上报 -- 顾客列表页操作
      wx.track.trackCustomerOperation({
        current_page: '沟通名单命名弹窗',
        operation_type: '取消'
      })

      this.setData({
        showSavePopup: false
      })
    },
    /**
     * 打包沟通名单 确定
     */
    async handleSaveConfirm(e) {
      // 神策数据上报 -- 顾客列表页操作
      wx.track.trackCustomerOperation({
        current_page: '沟通名单命名弹窗',
        operation_type: '确认打包'
      })
      // 包名称
      const { value } = e.detail

      const { innerValue, innerValueMap, packList } = this.data


      if (!value.trim()) {
        wx.showToast('沟通名单名称不能为空')
        return
      }

      if (packList.some(item => item.name === value)) {
        wx.showToast('沟通名单名称已存在')
        return
      }

      const params = {
        name: value,
        sort: 0,
        search: {
          value: JSON.stringify(innerValue),
          map: JSON.stringify(innerValueMap)
        }
      }

      const res = await addCrowdApi(params)
      if (res.code === 200) {
        // 打包成功
        this.setData({
          packName: value,
          showSavePopup: false,
          showPackSuccessPopup: true
        })
        this.gerCrowdList()
      } else {
        wx.showToast(res.message)
      }
    },
    /**
     * 关闭提示弹框
     */
    handleCloseMessage(e) {
      this.setData({
        showPackSuccessPopup: false,
        showCountMessagePopup: false,
        showPackCountMessagePopup: false
      })

      if (e) {
        const { name, content } = e.currentTarget.dataset
        if (name) {
          // 神策数据上报 -- 顾客列表页操作
          wx.track.trackCustomerOperation({
            current_page: name,
            operation_type: content
          })
        }
      }
    },
    /**
     * 去搜索
     */
    handleConfirm(e) {
      this.handleCloseMessage()
      this.handleClose()
      this.triggerEvent('confirm', {
        value: this.data.innerValue,
        count: this.data.innerValueMap.length
      })

      const { name, content } = e.currentTarget.dataset
      if (name) {
        // 神策数据上报 -- 顾客列表页操作
        wx.track.trackCustomerOperation({
          current_page: name,
          operation_type: content
        })
      }
      this.track('确定')
    },
    /**
     * 包管理
     */
    handleViewPack(e) {
      this.handleCloseMessage()
      wx.navigateTo({
        url: '/packages/customer/pages/communication-list/index'
      })

      const { name, content } = e.currentTarget.dataset
      if (name) {
        // 神策数据上报 -- 顾客列表页操作
        wx.track.trackCustomerOperation({
          current_page: name,
          operation_type: content
        })
      }
    },
    /**
     * 重置
     */
    handleReset(e) {
      const app = getApp()
      this.setData({
        innerValue: {
          rel_create_time_type: app.globalData.is_manager === 1 ? 2 : 1,
          rel_interact_time_type: 2,
          rel_active_time_type: 2
        }
      })
      this.setInnerValueMap()
      this.triggerEvent('confirm', {
        value: this.data.innerValue,
        count: this.data.innerValueMap.length
      })

      this.track('重置')

      const { name, content } = e.currentTarget.dataset
      if (name) {
        // 神策数据上报 -- 顾客列表页操作
        wx.track.trackCustomerOperation({
          current_page: name,
          operation_type: content
        })
      }
    },
    /**
     * 删除选中项
     */
    handleDelete(e) {
      const { key, value, type } = e.currentTarget.dataset

      const { innerValue } = this.data

      if (type === 2) {
        const index = innerValue[key].findIndex(item => item.key === value)
        innerValue[key].splice(index, 1)
      } else {
        delete innerValue[key]
      }

      this.setData({
        innerValue
      })
      this.setInnerValueMap()

      this.setData({
        loading: true
      })
      this.getCustomerList()
    },
    handleObserverChange(e) {
      if (this.disabledScroll) {
        return
      }
      const { value } = e.detail
      const { index } = e.currentTarget.dataset
      if (value) {
        this.setData({
          current: index
        })
      }
    },
    /**
     * 滚动到顶部 回到第一个
     * @param {*} e
     */
    handleScrollToUpper(e) {
      this.setData({
        current: 0
      })
    },
    /**
     * 开启滚动监听
     */
    handleTouchStart() {
      this.disabledScroll = false
    },
    handlePicker(e) {
      const { value } = e.detail
      this.setData({
        showPickerDate: true,
        innerDateValue: value
      })
    },
    /**
     * 日期选择
     * @param {*} e
     */
    handlePickerConfirm(e) {
      const { value } = e.detail
      const [ start, end ] = value
      this.setData({
        'innerValue.purchased_date': {
          start_time: start ? parseInt(start / 1000) : 0,
          end_time: end ? parseInt(end / 1000) : 0
        }
      })
      this.setInnerValueMap()
      this.setData({
        loading: true
      })
      this.getCustomerList()
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
    track(button_name) {
      const current_page = this.data.currentPage
      const { rel_create_time_type, rel_interact_time_type, rel_active_time_type } = this.data.innerValue


      const result = this.data.innerValueMap.reduce((prev, cur) => {
        const opt = this.data.options.find(item => item.key === cur.key)

        let text = ''
        if (cur.key === 'rel_create_time') {
          text = ['', '我的', '企业'][rel_create_time_type]
        }
        if (cur.key === 'rel_interact_time') {
          text = ['', '我的', '企业'][rel_interact_time_type]
        }

        if (cur.key === 'rel_active_time') {
          text = ['', '我的', '企业'][rel_active_time_type]
        }


        const { label } = this.data.group.find(item => item.value === opt.group)
        const label2 = cur.title
        if (cur.type === 1 || cur.type === 2) {
          prev = prev.concat(
            cur.result.map(item => {
              return {
                tag_level_1: label,
                tag_level_2: label2,
                tag_value: text ? (text + '-' + item.value) : item.value
              }
            })
          )
        } else if (cur.type === 3) {
          prev.push({
            tag_level_1: label,
            tag_level_2: label2,
            tag_value: `${dayjs(cur.result.start_time * 1000).format('YYYY/MM/DD')}-${dayjs(cur.result.end_time * 1000).format('YYYY/MM/DD')}`
          })
        } else {
          prev.push({
            tag_level_1: label,
            tag_level_2: label2,
            tag_value: `${cur.result.start_number}-${cur.result.end_number}`
          })
        }
        return prev
      }, [])


      const result_option = result.reduce((prev, cur) => {
        if (prev[cur.tag_level_2]) {
          prev[cur.tag_level_2].push(cur.tag_value)
        } else {
          const item = options.find(item => item.label === cur.tag_level_2)
          prev[cur.tag_level_2] = item.type === 2 ? [cur.tag_value] : cur.tag_value
        }
        return prev
      }, {})

      // console.log(result_option)
      // 神策数据上报 -- 筛选
      wx.track.trackFilter({
        current_page: current_page,
        button_name: button_name,
        filter_type_sv: 1,
        filter_option: result_option
      })

      result.forEach(item => {
      // 神策数据上报 -- 筛选详情
        wx.track.trackFilterDetails({
          current_page: current_page,
          button_name: button_name,
          filter_type_sv: 1,
          ...item
        })
      })
    }
  }
})
