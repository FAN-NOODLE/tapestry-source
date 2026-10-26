import { getRect, getAllRect } from '@/pagesMaterial/utils/util'

// Component Object
const app = getApp()
Component({
  externalClasses: ['custom-class'],
  options: {
    multipleSlots: true
  },
  properties: {
    value: {
      type: Array,
      observer: 'init'
    },
    /** 素材类型：image|video */
    fileType: {
      type: String,
      value: 'image'
    },
    /** 数量限制 */
    limit: {
      type: Number,
      value: 1
    },
    /** 展示形式 1｜2  宫格｜平铺 */
    styleType: {
      type: Number,
      value: 1
    },
    disabled: {
      type: Boolean,
      value: false
    }
  },
  data: {
    /** 是否全部原图 */
    isAll: false,
    /** 是否全部AI识别 */
    isAllAI: false,
    /** 拖拽区域 X */
    areaX: 0,
    /** 拖拽区域 Y */
    areaY: 0,
    /** 拖拽元素位置 X */
    x: 0,
    /** 拖拽元素位置 Y */
    y: 0,
    /** 触点位置 X */
    sX: 0,
    /** 触点位置 Y */
    sY: 0,
    /** 所有可以拖拽的元素 */
    rects: [],
    /** 当前拖动元素的下标 */
    current: -1
  },
  methods: {
    init() {
      const { value } = this.data
      if (value.length > 0) {
        this.setData({
          isAll: value.every(item => item.is_origin === 1),
          isAllAI: value.every(item => item.is_ai === 1)
        })
      }
    },
    /**
     * @description 初始化
     */
    initDrag() {
      if (this.data.disabled) return
      return Promise.all([
        getRect(this, '#movable-area'),
        getAllRect(this, '.drag-item')
      ]).then(([area, rects]) => {
        if (!area || !rects) return
        this.data.areaX = area.left
        this.data.areaY = area.top
        this.data.rects = rects.map(item => {
          return {
            x: item.left - area.left,
            y: item.top - area.top,
            width: item.width,
            height: item.height
          }
        })
        // console.log(area, rects)
      }).catch(err => {
        console.log(err)
      })
    },
    /**
     * @description 触屏开始
     * @param {*} e
     */
    async onTouchStart(e) {
      if (this.data.disabled) return
      await this.initDrag()
      const { rects } = this.data
      const { clientX, clientY } = e.changedTouches[0]
      const { index } = e.currentTarget.dataset
      if (index !== -1) {
        const { x, y, width, height } = rects[index]
        this.setData({
          x,
          y,
          width,
          height,
          sX: clientX,
          sY: clientY,
          current: index
        })
        // this.triggerEvent('dragBefore', { current: index })
      }
    },
    /**
     * @description 触屏移动
     * @param {*} e
     */
    onTouchMove(e) {
      if (this.data.disabled || !this.data.dragable) return
      const { sX, sY, x, y } = this.data
      const { clientX, clientY } = e.changedTouches[0]
      this.setData({
        x: x + clientX - sX,
        y: y + clientY - sY,
        sX: clientX,
        sY: clientY
      })
    },
    /**
     * @description 触屏结束
     * @param {*} e
     */
    onTouchEnd(e) {
      if (this.data.disabled) return
      // console.log(e.changedTouches[0])
      const { clientX, clientY } = e.changedTouches[0]
      const { rects, value, current, areaX, areaY } = this.data

      const index = rects.findIndex(item =>
        this.pointInRect({ x: clientX - areaX, y: clientY - areaY }, item),
      )

      // console.log(areaX, areaY, clientX, clientY, index, rects[index])
      if (index !== -1) {
        const temp = value.splice(current, 1)[0]
        value.splice(index, 0, temp)
      }
      this.setData({
        value,
        current: -1,
        sX: 0,
        sY: 0,
        dragable: false
      })
      this.triggerEvent('change', {
        value,
        source: 'drag'
      })
    },
    /**
     * 判断点是否在矩形内
     * @param {Object} point 要判断的点，形如 { x, y }
     * @param {Object} rect 矩形，形如 { x, y, width, height }
     * @returns {boolean} 点在矩形内返回 true，否则返回 false
     */
    pointInRect(point, rect) {
      return (
        point.x >= rect.x &&
        point.x <= rect.x + rect.width &&
        point.y >= rect.y &&
        point.y <= rect.y + rect.height
      )
    },
    handleUpload() {
      const { value, limit, fileType } = this.data
      const count = limit - value.length
      if (fileType === 'image') {
        wx.chooseImage({
          count: count,
          sizeType: ['original', 'compressed'],
          sourceType: ['album', 'camera'],
          success: (result) => {
            console.log(result)
            this.triggerEvent('change', {
              value: value.concat(result.tempFiles),
              source: 'select'
            })
          }
        })
      } else {
        if (app.globalData.isIos) {
          wx.showLoading({
            title: '正在处理中',
            mask: true
          })
          wx.chooseVideo({
            sourceType: ['album', 'camera'],
            compressed: false,
            maxDuration: 30,
            success: (result) => {
              console.log('chooseVideo success', result)
              wx.hideLoading()
              if (result.duration > 30) {
                wx.showToast('视频时长最长不能超过30S')
                return
              }
              if (result.size > 10 * 1024 * 1024) {
                wx.showToast('视频大小最大不能超过10M')
                return
              }
              this.triggerEvent('change', {
                value: value.concat([{
                  path: result.tempFilePath,
                  size: result.size,
                  cover_path: result.thumbTempFilePath || '' // 视频封面图，安卓手机获取不到
                }]),
                source: 'select'
              })
            },
            fail: (err) => {
              console.log('chooseVideo fail', err)
              wx.hideLoading()
            }
          })
        } else {
          wx.chooseMedia({
            count: 1,
            mediaType: ['video'],
            sourceType: ['album', 'camera'],
            maxDuration: 30,
            camera: 'back',
            success: (res) => {
              console.log('chooseMedia success', res)
              const result = res.tempFiles[0]
              // wx.hideLoading()
              if (result.duration > 30) {
                wx.showToast('视频时长最长不能超过30S')
                return
              }
              if (result.size > 10 * 1024 * 1024) {
                wx.showToast('视频大小最大不能超过10M')
                return
              }
              this.triggerEvent('change', {
                value: value.concat([{
                  path: result.tempFilePath,
                  size: result.size,
                  cover_path: result.thumbTempFilePath || '' // 视频封面图，安卓手机获取不到
                }]),
                source: 'select'
              })
            },
            fail: (err) => {
              console.log('chooseMedia fail', err)
              // wx.hideLoading()
            }
          })
        }
      }
    },
    /**
     * @description 全部使用原图
     */
    handleSelectAll() {
      const { value, isAll } = this.data

      value.forEach(item => {
        item.is_origin = isAll ? 0 : 1
      })

      this.triggerEvent('change', {
        value,
        source: 'selectAll'
      })
    },
    /**
     * @description 使用原图
     */
    handleChange(e) {
      const { value } = this.data
      const { index } = e.currentTarget.dataset
      const { item, source } = e.detail

      value[index] = item
      this.setData({
        value
      })
      this.triggerEvent('change', {
        value,
        source
      })
    },
    /**
     * @description 删除
     */
    handleDelete(e) {
      const { value } = this.data
      const { index } = e.currentTarget.dataset
      value.splice(index, 1)
      this.triggerEvent('change', {
        value,
        source: 'delete'
      })
    },
    onLongTap(e) {
      this.setData({
        dragable: true
      })
    },
    /**
     * @description 全部使用AI识别
     */
    handleSelectAllAI() {
      const { value, isAllAI } = this.data

      value.forEach(item => {
        item.is_ai = isAllAI ? 0 : 1
      })
      this.triggerEvent('change', {
        value,
        source: 'selectAllAI'
      })
    }
  }
})

