import { getMaterialGroupListApi } from '@/api/material-api'

function traverseTree(tree) {
  if (tree.length === 0) return
  tree.forEach(item => {
    item.children = item.children || []
    item.expand = false
    traverseTree(item.children)
  })
}
// Component Object
Component({
  externalClasses: ['custom-class'],
  properties: {
    // ids
    value: {
      type: Array,
      value: [],
      observer: 'init'
    },
    show: {
      type: Boolean,
      value: false
    },
    useSlot: {
      type: Boolean,
      value: false
    },
    text: {
      type: String,
      value: '未分组'
    }
  },
  data: {
    /** 是否显示 */
    show: false,
    /** 分类列表 */
    list: [],
    /** 选中的分类  */
    selected: [-1, -1, -1, -1],
    valueName: ''
  },
  lifetimes: {
    attached() {
      this.getMaterialCategory()
    }
  },
  methods: {
    getMaterialCategory() {
      getMaterialGroupListApi({
        page: 1,
        size: 999
      }).then(res => {
        res.data.unshift({
          id: this.data.text === '未分组' ? -1 : 0,
          parent_id: 0,
          group_name: this.data.text,
          children: []
        })
        traverseTree(res.data)
        this.setData({
          list: res.data.map(item => {
            item.children.unshift({
              id: 0,
              parent_id: item.parent_id,
              group_name: '全部',
              children: [],
              is_show: item.is_show
            })
            return item
          })
        }, () => {
          this.init()
        })
      }).catch(err => {
        console.log(err)
      })
    },
    /**
     * @description 初始化
     */
    init() {
      const { list, value } = this.data
      if (list.length === 0) return
      // 获取下标 index[]
      const selected = value.reduce((prev, id, index) => {
        if (index === 0) {
          const i = list.findIndex(item => Number(item.id) === Number(id))
          prev.push([i, list[i]])
          return prev
        } else {
          const i = prev[index - 1][1].children.findIndex(item => Number(item.id) === Number(id))
          prev.push([i, prev[index - 1][1].children[i]])
          return prev
        }
      }, []).map(item => item[0])

      // 填充
      while (selected.length < 4) {
        selected.push(-1)
      }

      // 处理空值
      // if (selected[0] === -1) {
      //   selected[0] = 0
      // }

      this.setData({
        selected
      }, () => {
        // TODO 选中项的名称， 可以在上面直接处理
        this.updateValueName()
      })
    },
    updateValueName() {
      const { list, selected } = this.data
      const value = selected.reduce((prev, item, index) => {
        if (item === -1) {
          return prev
        } else if (index === 0) {
          prev.push(list[item])
          return prev
        } else {
          prev.push(prev[index - 1].children[item])
          return prev
        }
      }, [])
      const name = value.filter(item => item && item.id).map(item => item.group_name).join('/')
      this.setData({
        valueName: name
      })
    },
    /**
     * @description 显示
     */
    show() {
      this.setData({
        show: true
      })
    },
    /**
     * @description 取消
     */
    hide() {
      this.setData({
        show: false
      })
    },
    /**
     * @description 切换
     */
    handleChange(e) {
      const { list } = this.data
      const { key, value, disabled } = e.currentTarget.dataset
      if (disabled) return

      // 判断展开收起
      const [first, secend, third] = value

      if (key === 1) {
        const current = `list[${first}].children[${secend}].expand`
        this.setData({
          [current]: !list[first].children[secend].expand
        })
      } else if (key === 2) {
        const current = `list[${first}].children[${secend}].children[${third}].expand`
        this.setData({
          [current]: !list[first].children[secend].children[third].expand
        })
      }

      // 设置选中值
      this.setData({
        selected: value
      })
    },
    /**
     * @description 确定
     */
    handleConfirm() {
      const { list, selected } = this.data
      // 获取选中项
      const value = selected.reduce((prev, item, index) => {
        if (item === -1) {
          return prev
        } else if (index === 0) {
          prev.push(list[item])
          return prev
        } else {
          prev.push(prev[index - 1].children[item])
          return prev
        }
      }, [])
      // 如果选中的是is_show == 0的选项直接过滤
      const isEmpty = value.every(item => item.is_show === 0)
      if (isEmpty) return

      this.triggerEvent('change', {
        value
      })
      this.hide()
    },
    /**
     * @description 取消
     */
    hanldeCancel() {
      this.hide()
    }
  }
})

