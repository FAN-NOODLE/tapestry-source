import { userBindRelationApi } from '@/api/customer-api'
import { toDate } from '@/utils/tools'

Component({
  properties: {
    detail: {
      type: Object,
      // 监听openid变化,获取绑定关系列表
      observer: function (detail) {
        if (detail) {
          this.getData()
        }
      }
    }
  },
  data: {
    list: [],
    meta: {
      page: 1,
      size: 10,
      total: 0,
      status: ''
    }
  },
  /*    系统事件开始   */
  created () {
  },
  attached () {},
  ready () {

  },
  /*    系统事件结束   */

  methods: {
    /*    普通事件开始  */
    async getData () {
      const {
        meta: { status },
        detail
      } = this.data

      if (status === 'loading') return
      this.setData({
        'meta.status': 'loading'
      })

      const params = {
        external_userid: detail.openid
      }

      try {
        const { code, data, message } = await userBindRelationApi(params)

        if (code !== 200) {
          this.setData({
            'meta.status': 'loadError'
          })
          wx.showToast(message)
          return
        }

        const newList = data.map(i => ({
          ...i,
          created_at: toDate(i.created_at, 'yyyy/MM/dd hh:mm:ss').nowTime
        }))
        this.setData({
          'meta.status': 'loaded',
          list: newList
        })
      } catch (error) {
        console.log(error)
      }
    },

    /*    普通事件结束  */

    /*    cdp分割事件   */
    cdpSplitEvent () {}
    /*    cdp分割事件   */

    /*    行为事件开始   */
    /*    行为事件结束   */
  }
})
