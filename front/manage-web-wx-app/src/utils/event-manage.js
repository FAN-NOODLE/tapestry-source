import Event from 'wx-event'

/* 注册管理主题事件 */
wx.event.$increase(new Event('date-update'))
wx.event.$increase(new Event('date-cancel'))
wx.event.$increase(new Event('refresh-list'))
wx.event.$increase(new Event('modal-slot-close'))
wx.event.$increase(new Event('modal-slot-confirm'))
// wx.event.$increase(new Event('close-area'))
/*
订阅事件主题模式 多线程分发事件
接收
wx.event.$on('date-update', params => {
      this.setData(params)
    })
触发
wx.event.$broadcast('date-update', params)
*/

