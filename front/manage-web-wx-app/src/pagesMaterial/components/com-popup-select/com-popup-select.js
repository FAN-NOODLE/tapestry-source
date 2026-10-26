import { customWhiteListApi } from '@/api/order'

// Component Object
Component({
  externalClasses: ['custom-class'],
  properties: {
    show: {
      type: Boolean,
      value: false
    }
  },
  data: {
    list: [
      {
        type: 1,
        show: true,
        url: '/pagesMaterial/pages/create/create?source_type=1',
        title: '文本素材',
        icon: 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_type_1.png'
      },
      {
        type: 2,
        show: true,
        url: '/pagesMaterial/pages/create/create?source_type=2',
        title: '图片素材',
        icon: 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_type_1.png'
      },
      {
        type: 5,
        show: true,
        url: '/pagesMaterial/pages/create/create?source_type=5&file_type=2',
        title: '视频素材',
        icon: 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_type_1.png'
      },
      {
        type: 7,
        show: false,
        url: '/pagesMaterial/pages/create/create?source_type=7',
        title: '朋友圈素材',
        icon: 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_type_1.png'
      }
    ]
  },
  lifetimes: {
    attached() {
      this.getCustomWhiteList()
    }
  },
  methods: {
    getCustomWhiteList() {
      customWhiteListApi({
        type: 2
      }).then(res => {
        if (res.code === 200) {
          this.setData({
            'list[3].show': res.data.show === 1
          })
        }
      }).catch(err => {
        console.log(err)
      })
    },
    /**
     * @description 显示弹窗
     */
    show() {
      this.setData({
        show: true
      })
    },
    /**
     * @description 隐藏弹窗
     */
    hide() {
      this.setData({
        show: false
      })
    },
    navigateTo(e) {
      const { url } = e.currentTarget.dataset
      this.hide()
      wx.navigateTo({
        url: url,
        events: {
          update: (data) => {
            this.triggerEvent('update', data)
          }
        }
      })
    }
  }
})

