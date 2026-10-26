import { userBindRelationApi } from '@/api/customer-api'
import { toDate } from '@/utils/tools'

Component({
  properties: {
    openid: {
      type: String,
      value: '',
      // 监听openid变化,获取绑定关系列表
      observer: function (newVal) {
        if (newVal) {
          this.getRelationList()
        }
      }
    }
  },
  data: {
    relationList: []
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
    async getRelationList () {
      const params = { external_userid: this.properties.openid }

      try {
        const { code, data, message } = await userBindRelationApi(params)

        if (code !== 200) {
          wx.showToast(message)
          return
        }

        const newList = data.map(i => ({
          ...i,
          created_at: toDate(i.created_at, 'yyyy/MM/dd hh:mm:ss').nowTime
        }))
        this.setData({
          relationList: newList
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
