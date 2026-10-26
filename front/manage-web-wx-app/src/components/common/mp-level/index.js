const { update_time } = require('@/utils/const')

// Component Object
Component({
  properties: {
    level: {
      type: String,
      value: '',
      observer: 'init'
    },
    size: {
      type: String,
      value: 'normal',
      observer: 'init'
    }
  },
  data: {
    imageHost: 'https://cdn-cos.coach.com.cn',
    levelImages: {
      normal: {
        DIAMOND: '/static/wxapp/common/level/level-0.png',
        GOLD: '/static/wxapp/common/level/level-1.png',
        SILVER: '/static/wxapp/common/level/level-2.png',
        C: '/static/wxapp/common/level/level-3.png',
        PROSPECT: '/static/wxapp/common/level/level-4.png',
        O: '/static/wxapp/common/level/level-5.png'
      },
      small: {
        DIAMOND: '/static/wxapp/common/level/level-list-0.png',
        GOLD: '/static/wxapp/common/level/level-list-1.png',
        SILVER: '/static/wxapp/common/level/level-list-2.png',
        C: '/static/wxapp/common/level/level-list-3.png',
        PROSPECT: '/static/wxapp/common/level/level-list-4.png',
        O: '/static/wxapp/common/level/level-list-5.png',
        6: '/static/wxapp/common/level/level-list-0.png',
        4: '/static/wxapp/common/level/level-list-1.png',
        5: '/static/wxapp/common/level/level-list-2.png',
        2: '/static/wxapp/common/level/level-list-3.png',
        1: '/static/wxapp/common/level/level-list-4.png',
        0: '/static/wxapp/common/level/level-list-5.png'
      }
    },
    newLevelImages: {
      normal: {
        DIAMOND: '/static/wxapp/common/level/level-0.png',
        GOLD: '/static/wxapp/common/level/level-1.png',
        SILVER: '/static/wxapp/common/level/level-2.png',
        C: '/static/wxapp/common/level/level-3-new.png',
        PROSPECT: '/static/wxapp/common/level/level-4-new.png',
        O: '/static/wxapp/common/level/level-5.png'
      },
      small: {
        DIAMOND: '/static/wxapp/common/level/level-list-0.png',
        GOLD: '/static/wxapp/common/level/level-list-1.png',
        SILVER: '/static/wxapp/common/level/level-list-2.png',
        C: '/static/wxapp/common/level/level-list-3-new.png',
        PROSPECT: '/static/wxapp/common/level/level-list-4-new.png',
        O: '/static/wxapp/common/level/level-list-5.png',
        6: '/static/wxapp/common/level/level-list-0.png',
        4: '/static/wxapp/common/level/level-list-1.png',
        5: '/static/wxapp/common/level/level-list-2.png',
        2: '/static/wxapp/common/level/level-list-3-new.png',
        1: '/static/wxapp/common/level/level-list-4-new.png',
        0: '/static/wxapp/common/level/level-list-5.png'
      }
    },
    currentLevelImage: ''
  },
  methods: {
    init() {
      const { level, size, levelImages, newLevelImages } = this.data
      const useNew = new Date() >= new Date(update_time)
      console.log(useNew, level, size)

      const path = useNew ? newLevelImages[size][level || 'O'] : levelImages[size][level || 'O']
      console.log(path)
      if (!path) return
      this.setData({
        currentLevelImage: `${this.data.imageHost}${path}`
      })
    }
  }
})
