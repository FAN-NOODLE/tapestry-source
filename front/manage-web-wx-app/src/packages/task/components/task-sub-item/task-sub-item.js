// Component Object
import { marketingPlanDelete } from '@/api/plan-task-api'
import { throttle } from '@/utils/tools'

Component({
  properties: {
    item: {
      type: Object
    }

  },
  data: {

  },
  methods: {
    handleShowTooltip(e) {
      const { selector } = e.currentTarget.dataset
      let content = '温馨提示：该任务转发后，任务状态变更会有5~10分钟的延迟哦~'
      if ([2].includes(this.data.item.group_send_type)) {
        content = '1、任务触达率=目标客户中联系状态为【已联系】顾客数/目标客户总数*100%<br/>2、任务状态更新会存在5～10分钟延迟'
      }
      this.selectComponent('#mp-tooltip').show(this, selector, content)
    },
    /**
     * 去做任务
     * dimension
     *  1.群/社群数据
     *  2.用户标签/用户画像
     *  3.导入顾客名单
     *  4.顾客与导购&导购匹配关系
     */
    handleViewTaskDetail(e) {
      const { id, group_number, customer_num, dimension } = this.data.item
      const num = dimension === 1 ? group_number : customer_num
      wx.navigateTo({ url: `/packages/task/pages/task-detail/task-detail?id=${id}&num=${num}&from=1` })
    },
    /** 隐藏任务提示 */
    handleHideRemind: throttle(function(e) {
      const { id } = this.data.item
      marketingPlanDelete({ id: id }).then(res => {
        if (res.code === 200) {
          this.setData({
            'item.close_tips': 1
          })
        } else {
          wx.showToast(res.message)
        }
      }).catch(err => {
        console.error('hideRemind 异常', err)
      })
    }, 2000)
  }
})
