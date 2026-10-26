const dayjs = require('dayjs')

// Component Object
Component({
  properties: {
    from: {
      type: String,
      value: ''
    },
    detail: {
      type: Object,
      value: null
    },
    item: {
      type: Object,
      value: () => ({
        'id': 192,
        'task_id': 3395850918359744,
        'status': 1,
        'task_category': 0,
        'name': '经典真挚',
        'remark': '经典真挚',
        'task_start_time': 1658131200,
        'task_end_time': 1658851199
      }),
      observer(value) {
        if (value) {
          this.setData({
            innerValue: {
              ...value,
              task_start_time_text: dayjs(value.task_start_time * 1000).format('YYYY-MM-DD HH:mm:ss'),
              task_end_time_text: dayjs(value.task_end_time * 1000).format('YYYY-MM-DD HH:mm:ss')
            }
          })
        }
      }
    }

  },
  data: {
    innerValue: null
  },
  methods: {
    /** 查看任务 */
    handleViewTask(e) {
      const { item } = this.data
      wx.navigateTo({
        url: `/packages/task/pages/task-detail/task-detail?id=${item.id}`,
        events: {
          // 为指定事件添加一个监听器，获取被打开页面传送到当前页面的数据
          update: res => {
            // 状态与当前不一致
            if (item.status !== res.status_code) {
              this.triggerEvent('update', res)
            }
          }
        }
      })
      if (this.data.from === '1') {
        // 神策数据上报 -- 顾客详情页操作
        const { openid } = this.data.detail
        wx.track.trackClickDetails({
          user_external_id: openid,
          operation_type: 10
        })
      } else {
        // 神策数据上报 -- 任务列表页操作
        wx.track.trackTaskListOperation({
          current_page: '任务列表页',
          operation_type: e.currentTarget.dataset.name
        })
      }
    }
  }
})
