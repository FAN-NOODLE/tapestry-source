import { getCustomerTaskListApi } from '@/api/customer-beta-api'

// Component Object
Component({
  properties: {
    detail: {
      type: Object,
      observer(detail) {
        if (detail) {
          this.getData()
        }
      }
    }
  },
  data: {
    list: [],
    /** 分页数据 */
    meta: {
      page: 1,
      size: 100,
      total: 0,
      status: ''
    }
  },

  lifetimes: {
    attached () {

    }
  },
  methods: {
    async getData(page = 1) {
      const {
        meta: { size, status },
        detail
      } = this.data

      if (status === 'loading') return
      this.setData({
        'meta.status': 'loading'
      })

      const params = {
        union_id: detail.unionid,
        page: page,
        size: size,
        status: 0,
        task_category: -1
      }

      try {
        const { data, code, message } = await getCustomerTaskListApi(params)
        wx.hideLoading()
        if (code !== 200) {
          this.setData({ isLoading: false })
          wx.showToast(message)
          return
        }
        this.data.list = (page === 1 ? [] : this.data.list).concat(data.list || [])

        this.setData({
          'meta.page': page,
          'meta.total': data.total,
          'meta.status': this.data.list.length < data.total ? 'loadMore' : 'loaded',
          list: this.data.list,
          isLoading: false
        })
      } catch (error) {
        wx.hideLoading()
        console.log(error)
      }
    },
    /**
     * @description 加载更多
     */
    loadMore() {
      const { page, status } = this.data.meta
      if (status === 'loadMore') {
        this.getData(page + 1)
      }
    },
    handleViewMore() {
      const { unionid } = this.data.detail
      wx.navigateTo({
        url: `/packages/task/pages/list/list?unionid=${unionid}`,
        events: {
          update: () => {
            this.getData()
          }
        }
      })

      // 神策数据上报 -- 顾客详情页操作
      const { openid } = this.data.detail
      wx.track.trackClickDetails({
        user_external_id: openid,
        operation_type: 11
      })
    },
    handleUpdate() {
      this.getData()
      this.triggerEvent('update')
    }
  }
})
