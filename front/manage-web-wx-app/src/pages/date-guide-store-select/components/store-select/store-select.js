import { memberPositionApi, storeOrDepartmentListApi, baStoreOrDepartmentListApi } from '@/api/compass-api'

const { globalData } = getApp()
Component({
  properties: {
    // 上次已经选中的门店或部门ids
    hasIds: {
      type: Array,
      value: []
    },
    mautiple: {
      // 是否多选
      type: Boolean,
      value: true
    }
  },
  data: {
    ready: false,
    isIpx: globalData.isIpx,
    chooseList: [], // 选中的人员数组
    currentList: [], // 当前展示列表
    staffList: [] // 全部人的列表
  },
  /*    系统事件开始   */
  created () { },
  attached () { },
  ready () {
    this.memberPosition()
  },
  /*    系统事件结束   */

  methods: {
    // 获取导购身份和权限类型
    memberPosition () {
      memberPositionApi().then(res => {
        if (res.code === 200) {
          this.setData({
            range: res.data.range
          })
          // 部门权限
          if (res.data.range === 2) {
            wx.setNavigationBarTitle({
              title: '选择指定门店'
            })
            this.setData({
              placeholder: '搜索门店名称',
              defultText: '暂无门店'
            })
          } else if (res.data.range === 3) {
            // 部门权限
            wx.setNavigationBarTitle({
              title: '选择指定部门'
            })
            this.setData({
              placeholder: '搜索部门名称',
              defultText: '暂无部门'
            })
          }
          if (this.data.mautiple) {
            this.getList()
          } else {
            this.getLeaderList()
          }
        } else {
          wx.showToast(res.message)
        }
      }).catch(err => {
        wx.showToast(err.message)
      })
    },
    // 获取门店、部门列表接口
    async getList () {
      try {
        const { code, data, message } = await storeOrDepartmentListApi()
        if (code === 200) {
          if (data.list && data.list.length > 0) {
            data.list.forEach((item, index) => {
              item.is_check = this.properties.hasIds.indexOf(item.id) > -1
              item.sort = index
            })
          }
          this.setData({
            staffList: data.list,
            currentList: data.list,
            chooseList: data.list.filter((item) => item.is_check),
            ready: true
          })
        } else {
          this.setData({
            ready: true
          })
          wx.showToast(message)
        }
      } catch (err) {
        this.setData({
          ready: true
        })
        wx.showToast(err.message)
      }
    },

    // 部门负责人
    async getLeaderList () {
      try {
        let query = {
          enterprise_id: globalData.currentEnterprise.enterprise_id,
          mobile: globalData.currentEnterprise.account_mobile
        }
        const { code, data } = await baStoreOrDepartmentListApi(query)
        if (code === 200) {
          this.setData({
            staffList: [{ counter_name: data.counter_name }],
            currentList: [{ counter_name: data.counter_name }],
            chooseList: [{ counter_name: data.counter_name }],
            ready: true
          })
        } else {
          this.setData({
            ready: true
          })
        }
      } catch (err) {
        this.setData({
          ready: true
        })
        wx.showToast(err.message)
      }
    },

    // 点击重置按钮
    reset () {
      const { staffList } = this.data
      this.setData({
        currentList: staffList.map(item => {
          item.is_check = 0
          return item
        }),
        chooseList: [],
        keywords: ''
      })
    },
    // 点击确认操作
    confirmTap () {
      let arr = this.data.chooseList
      arr.sort((a, b) => {
        return (a.id + '') > (b.id + '') ? 1 : -1
      })
      if (this.data.range === 2) {
        globalData.selectStore = this.data.chooseList
      }
      if (this.data.range === 3) {
        globalData.selectDepartment = this.data.chooseList
      }
      globalData.selectGuide = ''
      wx.navigateBack({
        delta: 1
      })
    },
    // 点击搜索的叉叉
    clearValue () {
      this.setData({
        currentList: this.data.staffList,
        keywords: ''
      })
    },
    // 搜索人员操作
    searchStaff (e) {
      const { value } = e.detail
      let filterList = this.data.staffList.filter((item) => (item.short_name && item.short_name.indexOf(value) > -1) || item.name.indexOf(value) > -1 || (item.mobile + '').indexOf(value) > -1)
      this.setData({
        currentList: filterList,
        keywords: value
      })
    },
    // 点击选择门店或部门操作(多选、单选)
    chooseItem (e) {
      const { item, index } = e.currentTarget.dataset
      let { chooseList, staffList } = this.data
      let deleteIndex = 0
      if (!item.is_check) {
        item.is_check = !item.is_check
        if (this.data.mautiple) {
          // 多选
          chooseList.push(item)
        } else {
          // 单选
          chooseList[0] = item
        }
        this.setData({
          [`currentList[${index}].is_check`]: true
        })
      } else {
        item.is_check = !item.is_check
        if (this.data.mautiple) {
          // 多选
          for (let i = 0; i < chooseList.length; i++) {
            if (chooseList[i].id === item.id) {
              deleteIndex = i
              break
            }
          }
        } else {
          // 单选
          deleteIndex = 0
        }
        chooseList.splice(deleteIndex, 1)
        this.setData({
          [`currentList[${index}].is_check`]: false
        })
      }
      console.log('选择的人', index, deleteIndex, chooseList)
      this.setData({
        chooseList,
        currentList: staffList.map(item => {
          item.is_check = chooseList.some(val => val.id === item.id)
          return item
        })
      })
    }
  }
})

