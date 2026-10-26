import { createSowTaskApi, getSowTaskListApi, finishedSowTaskApi, deleteSowTaskApi, sowTempListApi } from '@/api/sow-api'
import { toDate } from '@/utils/tools'
import { getRect } from '@/components/common/common/utils'
import { toggleTab } from '../../utils'

const dayjs = require('dayjs')

const defaultFilterList = [
  {
    key: 0,
    label: '全部'
  },
  {
    key: 1,
    label: '未完成'
  },
  {
    key: 2,
    label: '已完成'
  }
]
Component({
  properties: {
    detail: {
      type: Object
    },
    sidebar: {
      type: Boolean,
      value: false
    }
  },
  data: {
    editValue: '', // 播种信息
    filterList: defaultFilterList,
    filterKey: 0, // 查询字段
    sowList: [], // 播种记录列表
    dateTime: null, // 选择时间
    dateTimeText: '', // 选择时间文本
    isCreated: false, // 是否正在创建
    isLoading: false, // 是否正在请求列表
    meta: {
      page: 1,
      size: 10,
      total: 0,
      isMore: true
    },
    isOpenTipsModal: false,
    sowId: null,
    currentSow: null,
    isDelete: false,
    /** 播种模版列表 */
    sowTempList: [],
    currentTempIndex: 0,
    disabled: true,
    showOpened: true,
    opened: false
  },
  observers: {
    'detail': function (v) {
      if (v) {
        console.log('observers get list', v)
        // 重置meta参数
        this.setData({
          meta: {
            page: 1,
            size: 10,
            total: 0,
            isMore: true
          },
          isLoading: false
        })
        this.getList()
        this.getSowTempList()
      }
    }
  },
  /*    系统事件开始   */
  created () {},
  attached () {},
  ready () {
  },
  /*    系统事件结束   */

  methods: {
    /*    普通事件开始  */
    /**
     * 获取播种模版
     */
    getSowTempList() {
      sowTempListApi().then(res => {
        const { code, data, message } = res
        if (code !== 200) {
          wx.showToast(message)
          return
        }

        this.setData({
          sowTempList: [{
            id: 0,
            name: '默认播种',
            remind_time: 0,
            content: ''
          }, ...(data.length ? data : []).map(item => {
            const first = item.notice_list.shift()
            console.log(first)
            return {
              remind_time: 0,
              content: first.content,
              ...item
            }
          })]
        })
      }).catch(err => {
        console.log(err)
      })
    },
    /**
     * 获取播种记录
     */
    async getList (page = false) {
      const { meta, filterKey, isLoading, sowList, detail } = this.data

      console.log('get list', page, detail.id)

      // 正在请求
      if (isLoading || !meta.isMore) return

      this.setData({ isLoading: true })

      // 分页
      if (page) {
        this.setData({ 'meta.page': meta.page + 1 })
      }

      const params = {
        page: meta.page,
        size: meta.size,
        status: filterKey,
        work_external_user_id: detail.id
      }
      try {
        const { code, data, message } = await getSowTaskListApi(params)
        if (code !== 200) {
          wx.showToast(message)
          return
        }

        const { total } = data
        const list = data.list || []

        console.log(list, 'list')

        const newList = list.map(i => ({
          ...i,
          remind_time: toDate(i.remind_time, 'yyyy/MM/dd hh:mm:ss').nowTime
        }))
        const newSowList = page ? sowList.concat(newList) : newList

        this.setData({
          'meta.total': total,
          'meta.isMore': newSowList.length < total,
          sowList: newSowList,
          isLoading: false
        })

        console.log(data)
      } catch (error) {
        console.log(error)
        this.setData({ isLoading: false, 'meta.page': Math.max(1, meta.page - 1) })
      }
    },
    loadMore() {
      this.getList(true)
    },

    /*    普通事件结束  */

    /*    cdp分割事件   */
    cdpSplitEvent () {},
    /*    cdp分割事件   */

    /*    行为事件开始   */
    handleOpen() {
      this.setData({
        opened: !this.data.opened
      })
    },
    handleSowTypeChange(e) {
      const { value } = e.detail
      this.setData({
        currentTempIndex: value
      }, () => {
        this.handleSetDisabled()
      })
      wx.nextTick(() => {
        getRect(this, '#theme-text').then(rect => {
          console.log(rect)
          this.setData({
            showOpened: rect.height > 40
          })
        }).catch(err => {
          console.log(err)
        })
      })
    },
    handleSowTimeChange(e) {
      const { currentTempIndex } = this.data
      const { value } = e.detail
      this.setData({
        [`sowTempList[${currentTempIndex}].remind_time`]: value
      }, () => {
        this.handleSetDisabled()
        this.computedTime()
      })
    },
    handleSowContentChange(e) {
      const { currentTempIndex } = this.data
      const { value } = e.detail
      this.setData({
        [`sowTempList[${currentTempIndex}].content`]: value
      }, () => {
        this.handleSetDisabled()
      })
    },
    handleSetDisabled() {
      const { sowTempList, currentTempIndex } = this.data
      const item = sowTempList[currentTempIndex]
      this.setData({
        disabled: item.remind_time === 0 || item.content === ''
      })
    },
    computedTime() {
      const { sowTempList, currentTempIndex } = this.data
      const { remind_time, notice_list = [] } = sowTempList[currentTempIndex]
      const data = {}
      notice_list.forEach((item, index) => {
        if (item.delay_type === 1) {
          const date = new Date((remind_time + item.delay_value * 24 * 60 * 60) * 1000)
          const new_time = item.delay_time ? parseInt(date.setHours(...item.delay_time.split(':')) / 1000) : parseInt(date.getTime() / 1000)
          data[`sowTempList[${currentTempIndex}].notice_list[${index}].remind_time`] = new_time
        } else if (item.delay_type === 2) {
          const date = new Date((remind_time + 7 * item.delay_value * 24 * 60 * 60) * 1000)
          const new_time = item.delay_time ? parseInt(date.setHours(...item.delay_time.split(':')) / 1000) : parseInt(date.getTime() / 1000)
          data[`sowTempList[${currentTempIndex}].notice_list[${index}].remind_time`] = new_time
        } else {
          const date = new Date(remind_time * 1000)
          const new_date = new Date(date.setMonth(date.getMonth() + item.delay_value))
          const new_time = item.delay_time ? parseInt(new_date.setHours(...item.delay_time.split(':')) / 1000) : parseInt(new_date.getTime() / 1000)
          data[`sowTempList[${currentTempIndex}].notice_list[${index}].remind_time`] = new_time
        }
      })

      this.setData(data)
    },
    /**
     * 编辑器输入
     */
    editorInput (e) {
      const value = e.detail.value.trim().slice(0, 200)
      this.setData({ editValue: value })
      return { value }
    },
    /**
     * 创建播种
     */
    async createSow () {
      const { sowTempList, currentTempIndex, detail } = this.data
      const item = sowTempList[currentTempIndex]
      // 好友播种
      if (!detail.is_relation) {
        wx.showToast('该用户非您好友，暂无法设置播种')
        return
      }


      if (!item.remind_time) {
        wx.showToast('请选择创建时间')
        return
      }

      if (!item.content.trim()) {
        wx.showToast('请填写播种内容')
        return
      }

      if (dayjs().isBefore(item.remind_time)) {
        wx.showToast('请选择当前时间之后的日期')
        return
      }
      if (currentTempIndex !== 0) {
        this.selectComponent('#create-sow-popup').handleShow()
      } else {
        this.submit({
          detail: item
        })
      }
    },
    handleClose() {
      toggleTab()
    },
    async submit(e) {
      const value = e.detail
      if (this.data.isCreated) return
      this.setData({ isCreated: true })

      const { detail: { id, openid } } = this.properties

      const params = {
        work_external_user_id: id,
        task_content: [
          {
            content: value.content.trim(),
            remind_time: value.remind_time
          },
          ...(value.notice_list || [])
        ],
        template_id: value.id
      }

      if (params.task_content.find(item => !item.content.trim())) {
        wx.showToast('播种内容不能为空')
        return
      }

      try {
        const { code, data, message } = await createSowTaskApi(params)

        this.setData({ isCreated: false })

        if (code !== 200) {
          wx.showToast(message)
          return
        }
        wx.showToast('创建成功')
        this.selectComponent('#create-sow-popup').handleClose()
        this.triggerEvent('update')
        this.setData({
          meta: {
            page: 1,
            size: 10,
            total: 0,
            isMore: true
          },
          filterKey: 0,
          currentTempIndex: 0,
          'sowTempList[0].remind_time': 0,
          'sowTempList[0].content': ''
        }, () => {
          this.handleSetDisabled()
        })
        this.getList();
        (data || []).forEach(item => {
          // 神策数据上报 -- 创建播种
          wx.track.trackReminderCreate({
            sow_task_id: item.id,
            sow_task_content: item.content,
            user_external_id: openid,
            entrance: this.data.sidebar ? 1 : 2,
            template_name: value.name
          })
        })
      } catch (error) {
        console.log(error)
        this.setData({
          isCreated: false,
          currentTempIndex: 0,
          'sowTempList[0].remind_time': 0,
          'sowTempList[0].content': ''
        })
      }
    },
    /**
     * 弹窗点击确认
     */
    confirm () {
      if (this.data.isDelete) {
        this.deleteSow()
      } else {
        this.confirmFinished()
      }
    },
    /**
     * 删除播种
     */
    async deleteSow () {
      const { sowId, currentSow } = this.data
      const { detail: { id, openid } } = this.properties

      try {
        const { data, code, message } = await deleteSowTaskApi({ id: sowId, work_external_user_id: id })

        if (code !== 200) {
          wx.showToast(message)
          return
        }

        if (data.status === 'OK') {
          this.closeTipsModel()
          this.triggerEvent('update')
          wx.showToast('删除成功')
          this.setData({
            meta: {
              page: 1,
              size: 10,
              total: 0,
              isMore: true
            },
            filterKey: 0
          })
          this.getList()
          // 神策数据上报 -- 删除播种
          wx.track.trackReminderDelete({
            sow_task_id: sowId,
            sow_task_content: currentSow.content,
            user_external_id: openid
          })
        }
      } catch (error) {
        console.log(error)
      }
    },
    /**
     * 结束播种
     */
    async confirmFinished () {
      const { sowId, currentSow } = this.data
      const { detail: { openid } } = this.properties
      try {
        const { data, code, message } = await finishedSowTaskApi({ id: sowId })

        if (code !== 200) {
          wx.showToast(message)
          return
        }

        if (data.status === 'OK') {
          this.closeTipsModel()
          wx.showToast('播种完成')
          this.setData({
            meta: {
              page: 1,
              size: 10,
              total: 0,
              isMore: true
            },
            filterKey: 0
          })
          this.getList()

          // 神策数据上报 -- 播种完成
          wx.track.trackReminderComplete({
            sow_task_id: sowId,
            sow_task_content: currentSow.content,
            user_external_id: openid
          })
        }
      } catch (error) {
        console.log(error)
      }
    },
    /**
     * 时间选择
     */
    timeConfirm (e) {
      const { dateTimeText, dateTime } = e.detail
      this.setData({ dateTimeText, dateTime })
    },
    /**
     * 点击筛选条件
     */
    tapFilter (e) {
      this.setData({
        filterKey: e.target.dataset.key,
        'meta.page': 1,
        sowList: [],
        meta: {
          page: 1,
          size: 10,
          total: 0,
          isMore: true
        }
      })
      this.getList()
    },
    /**
     * 关闭弹窗
     */
    closeTipsModel () {
      this.setData({ isOpenTipsModal: false })
      toggleTab()
    },
    /**
     * 打开确认弹窗
     */
    openTipsModel (e) {
      const { id, item } = e.currentTarget.dataset
      this.setData({ isOpenTipsModal: true, sowId: id, isDelete: false, currentSow: item })
    },
    /**
     * 打开删除弹窗
     */
    openDeleteModel (e) {
      const { id, item } = e.currentTarget.dataset
      this.setData({ isOpenTipsModal: true, sowId: id, isDelete: true, currentSow: item })
      toggleTab()
    },
    /**
     * 无法创建播种时点击事件
     */
    tapDisableEdit () {
      if (!this.data.detail.is_relation) {
        wx.showToast('该用户非您好友，暂无法设置播种')
      }
    }
    /*    行为事件结束   */
  }
})
