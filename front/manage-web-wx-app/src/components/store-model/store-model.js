import { getMemberDetailApi } from '@/api/store-api'

const { split } = require('@/utils/tools')

const { globalData } = getApp()

Component({
  properties: {
    store: Object
  },
  data: {
    stores: {}
  },

  /*    系统事件开始   */
  created () { },
  attached () { },
  ready () {
    this.properties.store.work_phone = split(
      this.properties.store.work_phone,
      'space'
    )
    this.setData({
      stores: this.properties.store
    })
  },
  /*    系统事件结束   */

  methods: {
    select (e) {
      const { store } = this.data
      wx.showNavigationBarLoading()
      getMemberDetailApi({
        store_id: store.id,
        login_cache_key: globalData.login_cache_key
      }).then(res => {
        if (res.code === 200) {
          const { is_boss, id, store_sun_code, guide_role, guide_permissions } = res.data
          globalData.memberInfo = {
            isBoss: is_boss === 1,
            memberId: id
          }
          // 选择门店时，店员/店长角色从这里重新拿取赋值
          globalData.roleMark = guide_role
          // 选择门店时，店员/店长角色的 展示模块权限 从这里重新拿取
          globalData.currentEnterprise.guide_permissions = guide_permissions
          globalData.store = store
          globalData.store.store_id = store.id
          globalData.store.qr_code = store_sun_code
          this.triggerEvent('select', store)
        } else {
          wx.showToast(res.message)
        }
        wx.hideNavigationBarLoading()
      }).catch(err => {
        wx.hideNavigationBarLoading()
        wx.showToast(err.message)
      })
    }
  }
})
