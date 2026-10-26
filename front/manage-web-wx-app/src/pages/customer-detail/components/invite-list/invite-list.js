import { inviteListApi } from '@/api/guide-invite'

Component({
  properties: {
    external_userid: {
      type: String,
      value: ''
    },
    role: {
      type: Number,
      value: 3
    }
  },
  data: {
    inviteList: [],
    meta: {
      page: 1,
      size: 10,
      total: 0,
      isMore: true
    },
    isLoading: false
  },
  /*    系统事件开始   */
  created () {
  },
  attached () {},
  ready () {
    this.getList()
  },
  /*    系统事件结束   */

  methods: {
    /*    普通事件开始  */
    async getList (page = false) {
      const { isLoading, meta, inviteList } = this.data
      const { external_userid } = this.properties

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
        external_userid
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
