import { groupArr } from '../../utils/index'

Component({
  externalClasses: ['custom-class'],
  properties: {
    show: {
      type: Boolean
    },
    // 城市门店树列表
    childList: {
      type: Array,
      value: []
    },
    // 原始门店列表
    storeList: {
      type: Array
    },
    // 选中的门店列表
    checkedList: {
      type: Array,
      observer: function() {
        this.handleCheckStock()
        this.getCheckedNum()
      }
    }
  },
  data: {
    isAllChecked: false,
    checkedNum: 0,
    checkActive: true
  },
  methods: {
    matchStoreKeyword(item, value) {
      const keyword = String(value || '')
      const storeCode = String(item.third_store_code || item.loc_code || '')
      const city = String(item.city || '')
      return storeCode.indexOf(keyword) > -1 || city.indexOf(keyword) > -1
    },
    /**
     * @description: 重置数据
     * @return {*}
     */
    resetData(value) {
      let { storeList, checkedList } = this.properties
      // 从childList中筛选出符合条件的数据
      storeList = storeList.filter(item => {
        return this.matchStoreKeyword(item, value)
      })
      this.resetChecked(storeList, checkedList)
    },
    /**
     * @description: 筛选门店编码，城市编码,有库存
     */
    changeInput(e) {
      const { value } = e.detail
      if (this.data.checkActive) {
        this.handleCheckStock(value)
      } else {
        this.resetData(value)
      }
      this.setData({
        storeValue: value
      })
    },
    /**
     * @description: 单个门店选择
     * @return {*}
     */
    handleCheck(e) {
      const { index, storeIndex } = e.currentTarget.dataset
      this.handleCheckOne(index, storeIndex)
    },
    /**
     * @description: 单个选中方法
     */
    handleCheckOne(index, storeIndex) {
      // 更新checked状态
      const updateData = `childList[${index}].children[${storeIndex}].checked`
      this.setData({
        [updateData]: !this.data.childList[index].children[storeIndex].checked
      })
      // 获取的checked数量
      const checkedNum = this.data.childList[index].children.filter(item => item.checked).length
      const updateCheckedNum = `childList[${index}].checkedNum`
      //  更新checkedNum
      this.setData({
        [updateCheckedNum]: checkedNum
      })
      // 判断每个city是否全选
      const allChecked = this.data.childList[index].checkedNum === this.data.childList[index].children.length
      const updateAllChecked = `childList[${index}].allChecked`
      this.setData({
        [updateAllChecked]: allChecked
      })
      // 判断是否全选
      const isAllChecked = this.data.childList.filter(item => item.allChecked).length === this.data.childList.length
      this.setData({
        isAllChecked
      })
      this.getCheckedNum()
    },

    /**
     * @description: 城市全选
     * @return {*}
     */
    cityCheckAll(e) {
      const { index } = e.currentTarget.dataset
      // 找到该城市下的所有门店
      const storeList = this.data.childList[index].children
      // 将每个门店的checked设置为true
      storeList.map(item => {
        item.checked = !this.data.childList[index].allChecked
      })
      const updateChildren = `childList[${index}].children`
      // 更新checkedNum
      const updateCheckedNum = `childList[${index}].checkedNum`
      this.setData({
        [updateCheckedNum]: this.data.childList[index].children.filter(item => item.checked).length
      })
      // 更新checked状态
      const updateData = `childList[${index}].allChecked`
      this.setData({
        [updateData]: !this.data.childList[index].allChecked,
        [updateChildren]: storeList,
        [updateData]: !this.data.childList[index].allChecked
      })
      // 判断是否全选
      const isAllChecked = this.data.childList.filter(item => item.allChecked).length === this.data.childList.length
      this.setData({
        isAllChecked
      })
      this.getCheckedNum()
    },
    /**
     * @description: 所有城市全选
     * @return {*}
     */
    allCheckAll() {
      // 将每个门店的checked设置为true
      this.data.childList.map(item => {
        item.allChecked = !this.data.isAllChecked
        item.children.map(child => {
          child.checked = !this.data.isAllChecked
        })
      })
      // 更新checkedNum
      this.data.childList.map((item, index) => {
        item.checkedNum = item.children.filter(item => item.checked).length
      })
      // 更新checked状态
      this.setData({
        childList: this.data.childList,
        isAllChecked: !this.data.isAllChecked
      })
      this.getCheckedNum()
    },
    /**
     * @description: 获取选中数量
     * @return {*}
     */
    getCheckedNum() {
      let checkedNum = 0
      this.data.childList.map(item => {
        checkedNum += item.children.filter(item => item.checked).length
      })
      // 判断是否全选
      const isAllChecked = this.data.childList.filter(item => item.allChecked).length === this.data.childList.length
      this.setData({
        checkedNum,
        isAllChecked
      })
    },
    /**
     * @description: 确定
     */
    handleConfirm() {
      // 将childList中的checked为true的数据过滤出来
      const storeList = this.data.childList.map(item => {
        return item.children.filter(child => child.checked)
      })
      // 将storeList中的数据合并成一个数组
      const storeListArr = [].concat.apply([], storeList)
      // todo:判断是否有选中的门店
      if (storeListArr.length === 0) {
        wx.showToast('请选择门店')
        return
      }
      this.getCheckedNum()
      this.triggerEvent('confirm', {
        checkedList: storeListArr,
        checkedNum: this.data.checkedNum
      })
    },
    /**
     * @description: 取消
     */
    handleCancel() {
      this.triggerEvent('cancel')
    },
    /**
     * @description: 仅看有库存筛选
     */
    handleCheckStock(value) {
      // 仅看有库存
      let { storeList, checkedList } = this.properties
      // 如果是仅看有库存，筛选出有库存的门店
      if (this.data.checkActive) {
        storeList = storeList.filter(item => {
          return item.on_hand + item.in_transit > 0
        })
      }
      if (value) {
        storeList = storeList.filter(item => {
          return this.matchStoreKeyword(item, value)
        })
      }
      this.resetChecked(storeList, checkedList)
    },
    /**
     * @description: 重置数据设置选中状态
     */
    resetChecked(storeList, checkedList) {
      // 创建树形结构,将门店按照城市分组
      const treeData = groupArr(storeList, 'city')
      // 设置treeData中的对应checkedList的checked为true
      treeData.map(item => {
        item.children.map(child => {
          const index = checkedList.findIndex(item1 => item1.third_store_code === child.third_store_code)
          if (index > -1) {
            child.checked = true
          }
        })
        // 判断是否全选
        const allChecked = item.children.filter(item => item.checked).length === item.children.length
        item.allChecked = allChecked
        // 更新checkedNum
        item.checkedNum = item.children.filter(item => item.checked).length
      })
      // 判断是否全选
      const isAllChecked = treeData.filter(item => item.allChecked).length === treeData.length
      this.setData({
        childList: treeData,
        isAllChecked
      })
    },
    /**
     * @description: 点击仅看有库存
     */
    checkStock() {
      this.setData({
        checkActive: !this.data.checkActive
      })
      let inputVal = ''
      if (this.data.storeValue) {
        inputVal = this.data.storeValue
      }
      if (this.data.checkActive) {
        this.handleCheckStock(inputVal)
      } else {
        this.resetData(inputVal)
      }
    },
    /**
     * @description: 折叠城市项
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
     * @description:清空搜索框
     */
    clearSearch() {
      this.setData({
        storeValue: ''
      })
      if (this.data.checkActive) {
        this.handleCheckStock()
      } else {
        this.resetData('')
      }
    }
  }
})

