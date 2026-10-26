import { getCustomerUserUnbindtagApi, bindCustomerUsertagApi } from '@/api/customer-api'

// 获取全局应用程序实例对象
// const app = getApp();

// 创建页面实例对象
Page({
  // 页面的初始数据
  data: {
    inputWidth: '120px',
    relationId: '', // 导购顾客关联id
    currentList: [], // 搜索结果列表
    chooseTagList: [], // 已输入标签列表
    tagsList: [], // 标签列表
    isShow: false, // 是否显示过滤数组
    deleteIndex: -1 // 要删除的标签序号
  },
  /*    普通事件开始   */
  // 获取筛选自定义标签
  getCustomerUsertag () {
    getCustomerUserUnbindtagApi({
      page: 1,
      size: 10000,
      relation_id: this.data.relationId
    }).then(res => {
      const { code, data } = res
      if (code === 200) {
        // 获取之前选中标签
        let idList = wx.getStorageSync('idList') || []
        let chooseTagList = wx.getStorageSync('chooseTagList') || []
        let tagsList = data.length ? data : []
        if (idList && idList.length > 0) {
          for (let j = 0; j < tagsList.length; j++) {
            for (let i = 0; i < idList.length; i++) {
              if (tagsList[j].tag_id === idList[i]) {
                tagsList[j].selected = true
                chooseTagList.push(tagsList[j])
              }
            }
          }
        }
        this.setData({
          chooseTagList: chooseTagList,
          tagsList: tagsList
        })
        // this.setData({
        //   tagsList: data.length ? data : []
        // })
      } else {
        wx.showToast(res.message)
      }
    }).catch(err => {
      wx.showToast(err.message)
    })
  },
  /* 数据整理:给关键词加高亮
  * @param  {二维数组} originArr 筛选出的数组,
  * @param  {string} key 要替换的关键词,需要传入
  * @return {二维数组} 修改后返回
  */
  tidyData (originArr, key) {
    for (var i = 0; i < originArr.length; i++) {
      var tmp = originArr[i].tag_name
      var regExp = new RegExp(key, 'gi')
      originArr[i].light_name = tmp.replace(regExp, '<span class="highlight">' + key + '</span>')
    }
    return originArr
  },

  /*    普通事件结束   */

  /*    系统事件开始   */
  onLoad (options) {
    this.setData({
      relationId: parseInt(options.relation_id)
    })
    wx.removeStorageSync('idList')
    wx.removeStorageSync('chooseTagList')
  },
  onShow () {
    this.getCustomerUsertag()
  },
  onChange (e) {
    const { value } = e.detail
    let filterList = this.data.tagsList.filter((item) => item.tag_name.indexOf(value) > -1)
    this.setData({
      value: value,
      currentList: this.tidyData(filterList, e.detail.value),
      isShow: filterList.length > 0 && value.length > 0
    })
    wx.createSelectorQuery()
      .select('.search-text')
      .boundingClientRect((rect) => {
        if (rect) {
          this.setData({
            inputWidth: value.length > 0 ? rect.width + 'px' : '120px'
          })
        }
      }).exec()
  },
  onBlur () {
    this.setData({
      value: '',
      inputWidth: '120px'
    })
  },
  /*    系统事件结束   */

  /*    行为件开始   */
  // 搜索自定义标签选择
  checkedLabel (e) {
    let string = 'currentList[' + e.currentTarget.dataset.index + '].selected'
    this.setData({
      [string]: !this.data.currentList[e.currentTarget.dataset.index].selected
    })
    let detailValue = this.data.currentList.filter(it => it.selected).map(it => it)
    let newArr = this.data.chooseTagList.concat(detailValue) // 数组拼接
    this.setData({
      chooseTagList: newArr,
      tagsList: this.data.tagsList
    })
    this.clearInput()
  },
  // 自定义标签选择
  checkedGuideLabel (e) {
    let { index, id } = e.currentTarget.dataset
    let string = 'tagsList[' + index + '].selected'
    if (this.data.tagsList[index].selected && this.data.chooseTagList.length > 0) { // 取消选中
      for (let i = 0; i < this.data.chooseTagList.length; i++) {
        if (id === this.data.chooseTagList[i].tag_id) {
          this.data.chooseTagList.splice(i, 1)
        }
      }
    } else { // 选中
      this.data.chooseTagList.push(this.data.tagsList[index])
    }
    this.setData({
      deleteIndex: -1,
      [string]: !this.data.tagsList[index].selected,
      chooseTagList: this.data.chooseTagList
    })
  },
  // 清空输入框
  clearInput () {
    this.setData({
      value: '', // 清空输入框
      inputWidth: '120px',
      currentList: [],
      isShow: false
    })
  },
  // 回车键提交
  confirm (e) {
    if (e.detail.value) {
      // TODO 添加标签
      let chooseTagObj = {
        tag_id: 0,
        tag_name: e.detail.value
      }
      this.data.chooseTagList.push(chooseTagObj)
      this.setData({
        value: '', // 清空输入框
        inputWidth: '120px',
        chooseTagList: this.data.chooseTagList
      })
    }
  },
  // 删除
  deleteTag (e) {
    let { index, id } = e.currentTarget.dataset
    if (this.data.deleteIndex >= 0) {
      if (id !== 0 && this.data.tagsList.length > 0) {
        for (let i = 0; i < this.data.tagsList.length; i++) {
          if (id === this.data.tagsList[i].tag_id) {
            this.data.tagsList[i].selected = false
          }
        }
      }
      this.data.chooseTagList.splice(index, 1)
      this.setData({
        chooseTagList: this.data.chooseTagList,
        tagsList: this.data.tagsList,
        deleteIndex: -1
      })
    } else {
      this.setData({
        deleteIndex: index
      })
    }
  },
  // 保存
  saveData () {
    // 处理数组，删除多余属性
    this.data.chooseTagList.forEach(item => {
      if (item.selected) {
        delete item.selected
      }
      if (parseInt(item.user_num) === 0 || item.user_num) {
        delete item.user_num
      }
      if (item.light_name) {
        delete item.light_name
      }
    })
    let params = {
      relation_id: this.data.relationId,
      tag_names: this.data.chooseTagList
    }
    bindCustomerUsertagApi(params).then(res => {
      if (res.code === 200) {
        wx.showToast('绑定成功')
        // 跳回上一页
        wx.navigateBack({
          delta: 1
        })
      } else {
        wx.showToast(res.message)
      }
    }).catch(err => {
      wx.showToast(err.message)
    })
  },
  goToManageTags () {
    if (this.data.chooseTagList && this.data.chooseTagList.length > 0) {
      let chooseTagList = this.data.chooseTagList.filter(it => it.tag_id === 0).map(it => it)
      // 保存已输入标签列表
      wx.setStorageSync('chooseTagList', chooseTagList)
      // 保存已勾选标签id
      let idList = this.data.chooseTagList.filter(it => it.selected).map(it => it.tag_id)
      wx.setStorageSync('idList', idList)
    }
    // 跳转到管理标签页面
    wx.navigateTo({
      url: '/pages/tag-manage/tag-manage'
    })
  }
  /*    行为件结束   */
})
