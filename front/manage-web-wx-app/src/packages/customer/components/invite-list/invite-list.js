import { inviteListApi } from '@/api/guide-invite'

Component({
  properties: {
    detail: {
      type: Object,
      observer(detail) {
        if (detail) {
          this.getData()
        }
      }
    },
    role: {
      type: Number,
      value: 3
    }
  },
  data: {
    list: [],
    meta: {
      page: 1,
      size: 10,
      total: 0,
      status: ''
    },
    isLoading: false
  },
  /*    系统事件开始   */
  created () {
  },
  attached () {},
  ready () {
    this.getData()
  },
  /*    系统事件结束   */

  methods: {
    /*    普通事件开始  */
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
        page,
        size,
        external_userid: detail.openid
      }
      try {
        const { data, code, message } = await inviteListApi(params)
        if (code !== 200) {
          this.setData({
            'meta.status': 'loadError'
          })
          wx.showToast(message)
          return
        }

        this.data.list = (page === 1 ? [] : this.data.list)
          .concat((data.list || []).map((item) => {
            item.isOneMonth = this.isOneMonth(item.updated_at * 1000)
            return item
          }))

        this.setData({
          'meta.page': page,
          'meta.total': data.total,
          'meta.status': data.list.length > 0 ? 'loadMore' : 'loaded',
          list: this.data.list
        })
      } catch (error) {
        this.setData({
          'meta.status': 'loadError'
        })
      }
    },
    loadMore() {
      const { page, status } = this.data.meta
      if (status === 'loadMore') {
        this.getData(page + 1)
      }
    },
    async getData1 (page = 1) {
      const { isLoading, meta, inviteList, detail } = this.data


      if (isLoading || !meta.isMore) return

      this.setData({ isLoading: true })

      if (page) {
        this.setData({ 'meta.page': meta.page + 1 })
      }

      const params = {
        page: meta.page,
        size: meta.size,
        store_ids: [],
        member_ids: [],
        external_userid: detail.openid
      }

      try {
        const { code, data, message } = await inviteListApi(params)
        if (code !== 200) {
          wx.showToast(message)
          return
        }

        let { list = [], total } = data
        // todo:判断时间是否在一个月之内
        list = list.map((item) => {
          item.isOneMonth = this.isOneMonth(item.updated_at * 1000)
          return item
        })
        const newInviteList = page ? inviteList.concat(list) : list
        this.setData({
          inviteList: newInviteList,
          'meta.total': total,
          'meta.isMore': list.length === meta.size,
          isLoading: false
        })
      } catch (error) {
        console.log(error)
        this.setData({ isLoading: false, 'meta.page': Math.max(1, meta.page - 1) })
      }
    },
    // todo:判断时间是否在一个月之内
    isOneMonth (time) {
      const now = new Date().getTime()
      const oneMonth = 30 * 24 * 60 * 60 * 1000
      return now - time < oneMonth
    },
    /*    普通事件结束  */

    /*    cdp分割事件   */
    cdpSplitEvent () {}
    /*    cdp分割事件   */

    /*    行为事件开始   */
    /*    行为事件结束   */
  }
})
