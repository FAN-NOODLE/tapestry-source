import { groupArr } from '../../utils/index'

Component({
  externalClasses: ['custom-class'],
  properties: {
    show: {
      type: Boolean
    },
    list: {
      type: Array,
      observer: function (newVal) {
        this.setData({
          childList: newVal
        })
      }
    },
    storeList: {
      type: Array
    }
  },
  data: {
    show: false,
    childList: [],
    chooseItem: {}
  },
  methods: {
    matchStoreKeyword(item, value) {
      const keyword = String(value || '')
      const storeCode = String(item.third_store_code || item.loc_code || '')
      const city = String(item.city || '')
      return storeCode.indexOf(keyword) > -1 || city.indexOf(keyword) > -1
    },
    /**
     * @description: 筛选门店编码，城市编码,有库存
     */
    changeInput(e) {
      const { value } = e.detail
      let { storeList } = this.properties
      // 从childList中筛选出符合条件的数据
      storeList = storeList.filter(item => {
        return this.matchStoreKeyword(item, value)
      })
      // 创建树形结构,将门店按照城市分组
      const treeData = groupArr(storeList, 'city')
      this.setData({
        storeValue: value,
        childList: treeData
      })
    },
    /**
     * @description: 单个门店选择
     * @return {*}
     */
    handleCheck(e) {
      const { item } = e.currentTarget.dataset
      this.setData({
        chooseItem: item
      })
    },
    /**
     * @description: 折叠城市项
     * @return {*}
     */
    collapseItem(e) {
      const { cityIndex } = e.currentTarget.dataset
      const { childList } = this.data
      childList[cityIndex].isCollapse = !childList[cityIndex].isCollapse
      this.setData({
        childList
      })
    },
    /**
     * @description: 确定
     */
    handleConfirm() {
      // 判断是否选择门店
      if (!this.data.chooseItem.third_store_code) {
        wx.showToast('请选择门店')
        return
      }
      this.triggerEvent('confirm', this.data.chooseItem)
    },
    /**
     * @description: 取消
     */
    handleCancel() {
      this.setData({
        show: false,
        chooseItem: {}
      })
      this.triggerEvent('close')
    },
    /**
     * @description:清空搜索框
     */
    clearSearch() {
      this.setData({
        storeValue: ''
      })
      let { storeList } = this.properties
      // 创建树形结构,将门店按照城市分组
      const treeData = groupArr(storeList, 'city')
      this.setData({
        childList: treeData
      })
    }
  }
})

