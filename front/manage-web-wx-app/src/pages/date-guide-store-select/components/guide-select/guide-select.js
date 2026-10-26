import { guideListApi, baStoreOrDepartmentListApi } from '@/api/compass-api'

const { globalData } = getApp()
Component({
  properties: {
    // 已选的导购所属门店id数组
    storeIds: {
      type: Array,
      value: []
    },
    // 已选的导购所属部门id数组
    dpIds: {
      type: Array,
      value: ''
    },
    // 已选的导购id列表
    hasChooseGuideIds: {
      type: Array,
      value: []
    },
    role: {
      type: Number,
      value: 3
    },
    isBa: {
      type: Boolean,
      value: false
    },
    mautiple: {
      // 是否多选
      type: Boolean,
      value: true
    }
  },
  data: {
    ready: false,
    chooseList: [], // 选中的人员数组
    currentList: [], // 当前展示列表
    staffList: [], // 所有导购的列表
    mobileList: [],
    isAll: false,
    list: [],
    selectNum: 0
  },
  /*    系统事件开始   */
  created () { },
  attached () {
  },
  ready () {
    if (globalData.mobileList.length > 0) {
      this.data.mobileList = globalData.mobileList
      globalData.mobileList = []
    }
    if (!this.data.isBa) {
      this.getList()
    } else {
      this.getLeaderList()
    }
  },
  /*    系统事件结束   */

  methods: {
    /*    普通事件开始  */
    // 获取导购列表
    async getList () {
      const params = {
        store_ids: this.properties.storeIds || undefined,
        dp_ids: this.properties.dpIds,
        store_id: globalData.store.store_id || 0
      }
      // if (this.data.mobileList.length > 0) {
      //   params.mobiles = this.data.mobileList
      // }
      try {
        const { code, data, message } = await guideListApi(params)
        if (code === 200) {
          if (data.list.length > 0) {
            data.list.forEach((item, index) => {
              item.is_show = true
              item.typeName = this.typeDeal(item.type)
              item.is_check = this.properties.hasChooseGuideIds.indexOf(item.mem_uid) > -1
            })
          }
          this.setData({
            ready: true,
            dataList: data.list,
            selectNum: data.list.length,
            list: this.getShowList(data.list || []),
            isAll: this.getShowList(data.list || []).every(item => { return item.is_check })
          })
        } else {
          wx.showToast(message)
          this.setData({
            ready: true
          })
        }
      } catch (error) {
        this.setData({
          ready: true
        })
      }
    },
    getShowList (list) {
      const arr = []
      list.forEach(item => {
        const currentIndex = arr.findIndex(cell => { return cell.store_name === item.store_name })
        if (currentIndex === -1) {
          arr.push({
            store_name: item.store_name,
            is_check: false,
            is_show: true,
            show_all: false,
            guide_list: [item]
          })
        } else {
          arr[currentIndex].guide_list.push(item)
        }
      })
      arr.forEach(item => {
        item.is_show = item.guide_list.every(cell => { return cell.is_show })
        item.is_check = item.guide_list.every(cell => { return cell.is_check })
      })
      return arr
    },

    // 部门负责人列表
    async getLeaderList () {
      try {
        let query = {
          enterprise_id: globalData.currentEnterprise.enterprise_id,
          mobile: globalData.currentEnterprise.account_mobile
        }
        const { code, data } = await baStoreOrDepartmentListApi(query)
        if (code === 200) {
          if (data.leader_ba && data.leader_ba.length > 0) {
            data.leader_ba.forEach((item, index) => {
              item.is_check = this.properties.hasChooseGuideIds.indexOf(item.mem_uid) > -1
              item.sort = index
            })
          }
          this.setData({
            staffList: data.leader_ba,
            currentList: data.leader_ba,
            chooseList: data.leader_ba.filter((item) => item.is_check),
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

    // 导购类型处理
    typeDeal (type) {
      switch (type) {
        case 1:
          return '门店导购'
        case 2:
          return '社群导购'
        case 3:
          return '企业导购'
        case 4:
          return '企业管理'
      }
    },
    /*    普通事件结束  */

    /*    cdp分割事件   */
    cdpSplitEvent () { },
    /*    cdp分割事件   */

    /*    行为事件开始   */
    // 点击搜索的叉叉
    clearValue () {
      this.data.list.forEach(item => {
        item.is_show = true
        item.guide_list.forEach(cell => {
          cell.is_show = true
        })
      })
      this.setData({
        list: this.data.list,
        keywords: '',
        selectNum: this.data.dataList.length
      })
    },
    // 搜索人员操作
    searchStaff (e) {
      const { value } = e.detail
      this.data.list.forEach(item => {
        item.guide_list.forEach(cell => {
          cell.is_show = false
          if (cell.name.indexOf(value) > -1 || (cell.mobile + '').indexOf(value) > -1) {
            cell.is_show = true
          }
        })
        item.is_show = item.guide_list.some(cell => { return cell.is_show })
      })
      this.setData({
        list: this.data.list,
        keywords: value
      })
    },
    showGuideList (e) {
      const { index, item } = e.currentTarget.dataset
      this.setData({
        [`list[${index}].show_all`]: !item.show_all
      })
    },
    chooseCell (e) {
      const { ind, index } = e.currentTarget.dataset
      const { list } = this.data
      if (this.properties.mautiple) {
        const item = list[index]
        const cell = item.guide_list[ind]
        this.setData({
          [`list[${index}].guide_list[${ind}].is_check`]: !cell.is_check
        })
        const childCheck = item.guide_list.every(item => { return item.is_check })
        this.setData({
          [`list[${index}].is_check`]: childCheck
        })
        const listCheck = list.every(item => { return item.is_check })
        this.setData({
          isAll: listCheck
        })
        return
      }
      list.forEach((item, indexs) => {
        item.guide_list.forEach((cell, inds) => {
          cell.is_check = false
          if (index === indexs && ind === inds) {
            cell.is_check = true
          }
        })
      })
      this.setData({
        list
      })
    },
    // 点击选择导购操作(单选)
    chooseItem (e) {
      const { item, index } = e.currentTarget.dataset
      item.guide_list.forEach(cell => {
        cell.is_check = !item.is_check
      })
      this.setData({
        [`list[${index}].is_check`]: !item.is_check,
        [`list[${index}].guide_list`]: item.guide_list
      })
      const listCheck = this.data.list.every(cell => { return cell.is_check })
      this.setData({
        isAll: listCheck
      })
    },
    // 点击全部导购
    chooseAll () {
      const resetList = this.data.dataList.map(item => {
        item.is_check = !this.data.isAll
        return item
      })
      this.setData({
        isAll: !this.data.isAll,
        list: this.getShowList(resetList)
      })
    },
    reset () {
      const resetList = this.data.dataList.map(item => {
        item.is_check = false
        return item
      })
      this.setData({
        isAll: false,
        list: this.getShowList(resetList)
      })
    },
    // 确认搜索
    confirmTap () {
      const arr = []
      this.data.list.forEach(item => {
        item.guide_list.forEach(cell => {
          cell.is_check && arr.push(cell)
        })
      })
      globalData.selectGuide = arr.length === 0 ? this.data.dataList : arr
      globalData.selectAllGuide = (this.data.isAll || arr.length === 0) ? 1 : 0
      console.log(globalData.selectAllGuide, 'select=====')
      wx.navigateBack({
        delta: 1
      })
      // 神策数据上报 -- 顾客列表筛选
      let currentPages = getCurrentPages()
      if (currentPages[currentPages.length - 2].data.tabIndex) {
        const currentPage = 2 - currentPages[currentPages.length - 2].data.tabIndex
        wx.track.trackFilter({
          current_page: currentPage,
          filter_type_sv: 6,
          filter_option: {
            '导购': globalData.selectGuide.map(item => item.name).join()
          }
        })
        globalData.selectGuide.map(item => {
          // 神策数据上报 -- 筛选详情
          wx.track.trackFilterDetails({
            current_page: currentPage,
            button_name: '确定',
            filter_type_sv: 6,
            tag_level_1: '',
            tag_level_2: '导购',
            tag_value: item.name
          })
        })
      }
    }
    /*    行为事件结束   */
  }
})
