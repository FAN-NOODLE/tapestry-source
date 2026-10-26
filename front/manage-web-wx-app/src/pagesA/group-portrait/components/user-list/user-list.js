
const { globalData } = getApp()

Component({
  properties: {
    /** 用户信息 */
    list: { type: Array },
    /** */
    type: { type: Number },
    /** 用户id */
    ownerId: { type: String }
  },
  data: {
    /** 用户id */
    userId: ''
  },

  ready () {
    this.setData({
      userId: globalData.userid
    })
  },
  methods: {
    /** 打开立即沟通聊天 */
    communicate (id, type) {
      const params = {
        userIds: '',
        externalUserIds: ''
      }
      const attr = type === 1 ? 'userIds' : 'externalUserIds'

      params[attr] = id
      wx.qy.openEnterpriseChat({
        ...params,
        groupName: '',
        success: function (res) {
          console.log(res, '成功')
        },
        fail: function (err) {
          console.log(err, '错误')
        }
      })
    },

    /** 添加好友 */
    addPerson (id, type) {
      wx.qy.openUserProfile({
        type,
        userid: id,
        success: function (res) {
          console.log(res, '成功')
        },
        fail: function (err) {
          console.log(err, '失败')
        },
        complete: function (complete) {
          console.log(complete, '完成')
        }
      })
    },

    /** 操作，分发添加好友或者联系客户 */
    operation (e) {
      const {friend, user_id, type} = e.currentTarget.dataset.item
      if (friend === 1) {
        // 打开会话
        this.communicate(user_id, type)
      } else {
        this.addPerson(user_id, type)
      }
    }
  }
})

