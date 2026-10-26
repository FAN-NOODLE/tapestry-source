import { throttle } from '@/utils/tools'
Component({
  properties: {
    info: {
      type: Object,
      value: {}
    }
  },
  data: {
    expand: true
  },

  /*    系统事件开始   */
  created () {},
  attached () {},
  ready () {
  },
  /*    系统事件结束   */

  methods: {
    /*    普通事件开始  */

    /*    普通事件结束  */

    /*    cdp分割事件   */
    cdpSplitEvent () {},
    /*    cdp分割事件   */

    /*    行为事件开始   */
    /**
     * 展开折叠
     */
    expand () {
      this.setData({ expand: !this.data.expand })
    },
    /**
     * 选择和取消选择标签
     */
    tagsImpressionOperate: throttle(function (e) {
      this.triggerEvent('tagsImpressionOperate', e.currentTarget.dataset)
    }, 300)
    /*    行为事件结束   */
  }
})
