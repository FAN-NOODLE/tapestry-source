import { getCurrentPage } from '@/components/common/common/utils'
import { debounce } from '../../../../utils/tools'
import { memberPositionApi, storeOrDepartmentListApi } from '@/api/compass-api'

const { globalData } = getApp()
Component({
  properties: {
    // 上次已经选中的门店或部门id
    hasIds: {
      type: Array,
      value: []
    }
  },
  data: {
    ready: false,
    isIpx: globalData.isIpx,
    currentList: [], // 当前门店列表
    chooseList: [],
    isloading: false,
    allCheck: false,
    checkedStore: []
  },
  /*    系统事件开始   */
  created () { },
  attached () { },
  ready () {
    this.memberPosition()
    console.log('this.properties.hasIds', this.properties.hasIds)
    this.data.checkedStore = this.properties.hasIds
  },
  methods: {
    // 搜索人员操作
    searchStaff: debounce(async function (e) {
      const { value } = e.detail
      this.setData({
        keywords: value
      })
      this.getList()
    }, 500),
    // 点击搜索的叉叉
    clearValue () {
      this.setData({
        currentList: [],
        keywords: ''
      })
      this.getList()
    },
    // 获取导购身份和权限类型
    memberPosition () {
      let params = {}
      if (globalData.brand_code === 'HR') {
        params = {}
      } else {
        params = {is_new: 0}
      }
      memberPositionApi(params).then(res => {
        if (res.code === 200) {
          this.setData({
            range: res.data.range
          })
          this.getList()
        } else {
          wx.showToast(res.message)
        }
      }).catch(err => {
        wx.showToast(err.message)
      })
    },
    // 获取门店、部门列表接口
    async getList () {
      this.setData({
        isloading: true
      })
      try {
        const params = {
          is_new: 0,
          name: this.data.keywords
        }
        const { code, data, message } = await storeOrDepartmentListApi(params)
        this.setData({
          isloading: false
        })
        if (code === 200) {
          data.list.forEach(item => {
            item.is_check = this.data.checkedStore.indexOf(item.third_id) > -1
          })
          this.setData({
            currentList: data.list,
            ready: true,
            allCheck: !data.list.some(item => {
              return !item.is_check
            })
          })
        } else {
          this.setData({
            ready: true
          })
          wx.showToast(message)
        }
      } catch (err) {
        this.setData({
          ready: true,
          isloading: false
        })
        wx.showToast(err.message)
      }
    },
    // 选择门店
    checkedChange (e) {
      const { item, index } = e.currentTarget.dataset
      this.setData({
        [`currentList[${index}].is_check`]: !item.is_check
      })
      this.setData({
        allCheck: this.data.currentList.every(item => { return item.is_check })
      })
    },
    // 全选
    checkAllCallBack (e) {
      this.data.currentList.forEach(item => {
        item.is_check = !this.data.allCheck
        return item
      })
      this.setData({
        currentList: this.data.currentList,
        allCheck: !this.data.allCheck
      })
    },
    // 重置
    reset () {
      this.setData({
        keywords: ''
      })
      this.data.checkedStore = []
      this.getList()
    },
    // 确认门店
    confirmTap () {
      this.data.currentList.forEach(item => {
        if (item.is_check) {
          this.data.chooseList.push(item)
        }
      })
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
      // 是否全选全部部门=统计那边有用到
      globalData.selectAllDepartment = this.data.allCheck ? 1 : 0
      globalData.selectGuide = ''
      wx.navigateBack({
        delta: 1
      })


      // 神策数据上报 -- 顾客列表筛选
      let currentPages = getCurrentPages()
      if (currentPages[currentPages.length - 2].data.tabIndex) {
        const currentPage = 2 - currentPages[currentPages.length - 2].data.tabIndex
        wx.track.trackFilter({
          current_page: currentPage,
          filter_type_sv: 5,
          filter_option: {
            '门店': this.data.chooseList.map(item => item.name).join()
          }
        })

        this.data.chooseList.map(item => {
          // 神策数据上报 -- 筛选详情
          wx.track.trackFilterDetails({
            current_page: currentPage,
            button_name: '确定',
            filter_type_sv: 5,
            tag_level_1: '',
            tag_level_2: '门店',
            tag_value: item.name
          })
        })
      }
    }
  }
})
