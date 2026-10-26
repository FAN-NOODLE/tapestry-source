import { getPhoneApi, getTimes, createReserveApi, editReserveApi } from '@/api/booking-customer-api'
import { getConfigMessageApi, messageSuccessApi } from '@/api/store-invite'

const { globalData } = getApp()
/** 默认小时，用于判断锚点 */
const HOURS = [ 9, 10, 11, 13, 14, 15, 16 ]

/** 生成默认时间戳 */
const generateDefaultTime = (date) => {
  /** 默认时间段 */
  const defaultTimes = [
    {time: '9:00', reserve_num: 0, reserved_num: 0},
    {time: '10:00', reserve_num: 0, reserved_num: 0},
    {time: '11:00', reserve_num: 0, reserved_num: 0},
    {time: '13:00', reserve_num: 0, reserved_num: 0},
    {time: '14:00', reserve_num: 0, reserved_num: 0},
    {time: '15:00', reserve_num: 0, reserved_num: 0},
    {time: '16:00', reserve_num: 0, reserved_num: 0}
  ]

  // 生成上午时间戳
  const times = defaultTimes.map((item) => {
    let time = date + ' ' + item.time
    time = time.replace(/-/g, '/')
    const reserve_time = new Date(time).getTime()
    return {
      ...item,
      reserve_time: reserve_time / 1000
    }
  })

  return times
}

/** 性别 */
const GENDERS = [
  {value: 0, name: '女士', checked: true},
  {value: 1, name: '先生'}
]

Component({
  options: {
    addGlobalClass: true
  },
  properties: {
    /** 预约客户 */
    info: { type: Object }
  },
  data: {
    /** 是否展示、关闭预约信息弹窗 */
    show: false,
    /** 表单信息 */
    form: {
      gender: 0
    },
    /** 用户性别 */
    genders: GENDERS,
    /** 预约时间默认选中下标 */
    current_index: -1,
    /** 默认时间戳 */
    times: [],
    /** 预约时间 */
    today: '',
    /** 是否同意使用协议 */
    agree: false,
    /** 是否编辑 */
    isEdit: false,
    /** 时间段锚点 */
    anchorPointId: '',
    activity_id: globalData.activity_id
  },
  observers: {
    show: function (show) {
      if (!show) return
      this.getTimeList()
      this.getTempId()
    }
  },

  created () {
    // 获取编辑预约组件实例
    this.toastPopup = this.selectComponent('#toast')
  },

  methods: {
    /** 获取订阅模板id */
    getTempId () {
      getConfigMessageApi().then(res => {
        console.log('获取订阅消息配置res=====》', res)
        this.data.tmplIds = [res.data['confirm-subscribe'][0]]
      }).catch(err => {
        console.error(err, '获取订阅模板出错')
      })
    },

    /** 更改是否同意安全协议 */
    changeAgree () {
      this.setData({
        agree: !this.data.agree
      })
    },

    /** 切换选中时间段 */
    changeChecked (e) {
      const { disabled, index } = e.currentTarget.dataset

      if (disabled) return
      this.setData({
        current_index: index
      })
    },

    /** 获取预约时间段 */
    getTimeList () {
      // 获取当前时间戳
      const reserve_time = Date.parse(new Date()) / 1000
      getTimes({reserve_time})
        .then(res => {
          // 错误处理
          if (res.code !== 200) return wx.showToast(res.message)

          const { lists, today } = res.data

          let times = []
          // 判断时间戳list是否为空
          if (!lists.length) {
            // 使用默认时间戳
            times = generateDefaultTime(today)
          } else {
            // 采用接口数据
            times = lists
          }
          let anchorPointId = -1 // 锚点
          const tableData = this.differenceTime(times)

          // 处理修改数据
          const {user_name, mobile, gender, reserve_time, is_reserve} = this.data.info
          let data = {}
          if (is_reserve) {
            const index = tableData.findIndex(item => item.reserve_time === reserve_time)
            data = {
              form: {
                user_name,
                mobile,
                gender
              },
              current_index: index,
              agree: true
            }
            anchorPointId = index
          }

          // 编辑状态，如果当前有选中数据，制空默认状态
          if (is_reserve && data.current_index !== -1) {
            tableData[data.current_index].icon = ''
            tableData[data.current_index].disabled = false
          }

          // 设置时间段锚点，如果没有选中锚点，就判断时间，定位到当前时间段中
          if (anchorPointId === -1) {
            anchorPointId = HOURS.findIndex((item) => item === new Date().getHours())
          }
          this.setData({
            times: tableData,
            today,
            ...data
          })

          // 关闭弹窗之后，再次打开，值没有变，不会再次触发滚动，使用setTimeout触发会再次执行
          setTimeout(() => {
            this.setData({
              anchorPointId: `cont${anchorPointId}`
            })
          }, 100)
        }).catch(err => {
          console.error(err, '获取时间段出错')
        })
    },

    /** 生成时间段 */
    differenceTime (data) {
      const nowHour = new Date().getHours() // 当前时间（小时）

      const newData = []
      data.forEach((item) => {
        // 获取每个时间段的小时
        let hour = new Date(item.reserve_time * 1000).getHours()
        // 是否置灰
        let disabled = false
        // 展示文案、背景图
        let text = ''
        // 展示background-image
        let icon = ''
        const reserve_num = Number(item.reserve_num)
        const reserved_num = Number(item.reserved_num)
        if (nowHour > hour) {
          //  已超过当前时间段内容的，选框置灰，显示【已结束】
          disabled = true
          text = '已结束'
          icon = 'is-over'
        } else if (reserve_num === 0 && reserved_num === 0) {
          // 预约时间段中，若后台数值已预约人数与剩余预约数均为0时，则显示【未开始】
          disabled = true
          text = '未开始'
          icon = 'is-start'
        } else if (reserved_num > 0 && reserve_num === 0) {
          // 若后台数值已预约人数大于0，但剩余可预约数为0时，则显示【已约满】
          disabled = true
          text = '已约满'
          icon = 'is-full'
        } else if (reserve_num > 0) {
          // 预约时间段中，若后台剩余可预约名额有大于0的数值，则展示【可预约】
          disabled = false
          text = '可预约'
        }

        const renderTime = `${hour}:00~${hour + 1}:00`
        newData.push({...item, renderTime, disabled, text, icon})
      })
      // 返回上午、下午时间段，统一setData
      return newData
    },

    /** 保存成功后执行订阅 */
    afterSubmit (id) {
      let self = this
      const { tmplIds } = this.data
      wx.requestSubscribeMessage({
        tmplIds,
        success (res) {
          if (res[tmplIds[0]] === 'accept') {
            self.afterSubscribeInfo(id)
          } else {
            // 关闭弹窗
            self.updatePopup(false)
            // 通知父组件预约成功
            self.triggerEvent('handleSuccess')
          }
        },
        fail (err) {
          console.log(err, '订阅失败')
        }
      })
    },

    /** 订阅消息 */
    afterSubscribeInfo (id) {
      const { tmplIds } = this.data
      messageSuccessApi({
        rel_id: id,
        tmpl_list: tmplIds.map(item => {
          return {
            pri_tmpl_id: item,
            url: 'pages/booking-customer-activities/booking-customer-activities'
          }
        })
      }).then(() => {
        // 关闭弹窗
        this.updatePopup(false)
        // 通知父组件预约成功
        this.triggerEvent('handleSuccess')
      }).catch(e => console.error('订阅出错'))
    },

    /** 提交预约信息 */
    submit () {
      const {current_index, times, form, agree, info} = this.data
      const {user_name, mobile, gender} = form
      // 昵称校验
      if (!user_name) {
        return wx.showToast('请输入昵称')
      }

      let phoneCheck = /^1[3456789]\d{9}$/

      // 手机号校验
      if (!phoneCheck.test(mobile)) {
        const message = mobile ? '请输入正确的手机号' : '请输入预约手机号'
        return wx.showToast(message)
      }

      // 时间段校验
      if (!times[current_index]) {
        return wx.showToast('请选择预约时间段')
      }

      // 安全协议验证
      if (!agree) {
        return wx.showToast('请勾选用户协议')
      }

      // 组合参数
      const params = {
        activity_id: this.data.activity_id,
        reserve_time: times[current_index].reserve_time,
        user_name,
        mobile,
        gender,
        reserve_id: info.id
      }

      const api = info.is_reserve ? editReserveApi : createReserveApi
      api(params)
        .then(res => {
          // 错误处理
          if (res.code !== 200) return wx.showToast(res.message)

          // this.toastPopup.openPopup('success', '预约成功')

          // 新增时需要订阅，单独处理
          if (!info.is_reserve) {
            return this.afterSubmit(res.data.reserve_id)
          }

          // 关闭弹窗
          this.updatePopup(false)
          // 通知父组件预约成功
          this.triggerEvent('handleSuccess')
        }).catch(err => {
          console.error(err, '新建预约信息出错')
        })
    },

    /** 轮播预约时间 */
    changeIndex (e) {
      this.setData({
        current_index: e.detail.current
      })
    },

    /** 关闭弹窗 */
    handleClose () {
      this.updatePopup(false)
    },

    /** 打开、关闭弹窗 */
    updatePopup (show = true, isEdit = false) {
      // 关闭弹窗
      if (!show) {
        // 清除数据
        return this.setData({
          form: { gender: 0 },
          times: [],
          current_index: -1,
          agree: false,
          show: false,
          isEdit: false
        })
      }

      // 打开弹窗
      this.setData({
        show,
        isEdit
      })
    },

    /** input 变化 */
    changeInput (e) {
      // 通过事件参数获取input所对应的全局属性名
      const inputModel = e.currentTarget.dataset.name
      // 通过value获取当前用户输入的内容
      const value = e.detail.value
      // 更新data, 将input所对应的全局属性的属性只更新
      this.setData({
        form: {...this.data.form, [inputModel]: value}
      })
    },

    /** 获取授权手机号 */
    getPhoneNumber (e) {
      const code = e.detail.code
      // todo:授权手机号报错提示
      if (e.detail.errno) {
        wx.report({
          title: 'coach小程序获取手机号失败',
          time: new Date().getTime(),
          errMsg: e.detail
        })
      }
      if (!code) return
      getPhoneApi({code})
        .then(res => {
          // 错误处理
          if (res.code !== 200) return wx.showToast(res.message)

          // 保存手机号
          this.setData({
            'form.mobile': res.data.phoneNumber || ''
          })
        }).catch(err => {
          console.error(err, '授权手机号报错')
        })
    },

    /** 计算滚动条位置 */
    getLeft (e) {
      // 480 滚动条总长度， 160 滚动条内部滑块长度   160/480 = 0.337
      // 最长滑动距离 296, 求滚动比例   480/296 = 1.62
      const slideLeft = e.detail.scrollLeft * 1.62 - 100
      this.setData({
        slideLeft
      })
    }
  }
})
