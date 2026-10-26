// Component Object
import { upload } from '@/utils/upload'
import { compressFileApi } from '@/api/material-api'
import { downloadFile2 } from '@/pagesMaterial/utils/util'

Component({
  externalClasses: ['custom-class'],
  properties: {
    item: {
      type: Object
    },
    styleType: {
      type: Number,
      value: 1
    },
    /** 素材类型：image|video */
    fileType: {
      type: String,
      value: 'image'
    },
    compress: {
      type: Boolean,
      value: true
    }
  },
  data: {
    name: '',
    show: false,
    keyboardHeight: 20,
    innerSrc: '',
    percent: 0,
    uploaded: false
  },
  lifetimes: {
    attached() {
      this.upload()
    }
  },
  methods: {
    async upload() {
      const { item, fileType, compress } = this.data

      if (item.loaded) return
      try {
        // 上传素材到cos
        const [uploadRes] = await upload(item.path, {
          onProgress: (res) => {
            this.setData({
              percent: res.percent
            })
          }
        })

        // 图片、朋友圈素材的图片
        // 上传图片：
        // 1、选原图上传不压缩，不选原图超过1M压缩上传
        // 2、选原图超限制也可上传不压缩，但是在聊天框发送素材的时候会压缩后发送

        // 上传朋友圈素材：
        // 1、选原图上传不压缩，不选原图超过1M压缩上传
        // 2、选原图超限制不能上传，并提醒哪张图片上传超限制了
        // if (fileType === 'image' && item.size > 30 * 1024 * 1024) {
        //   try {
        //     const compressRes = await compressFileApi({
        //       file_path: uploadRes.url,
        //       file_type: fileType
        //     })
        //     if (compressRes.code === 200) {
        //       uploadRes.compress_url = compressRes.data.file_url
        //     }
        //   } catch (error) {
        //     console.log(error)
        //   }
        // }

        // 视频素材 获取视频封面图并上传cos ci-process=snapshot&time=1.5
        if (fileType === 'video') {
          // 下载视频封面图
          const coverRes = await downloadFile2(`${uploadRes.url}?ci-process=snapshot&time=1.5`)
          item.cover_path = coverRes.tempFilePath
          // 上传视频封面图到cos
          const [uploadCoverRes] = await upload(coverRes.tempFilePath)
          item.cover_url = uploadCoverRes.url
        }
        this.setData({
          uploaded: true,
          item: {
            ...uploadRes,
            ...item
          }
        })
        this.triggerEvent('change', {
          item: {...this.data.item, loaded: true},
          source: 'upload'
        })
      } catch (error) {
        console.log(error)
      }
    },
    /**
     * @description 编辑名称
     */
    handleEditName() {
      const { item } = this.data
      const name = item.name.split('.')[0]
      this.setData({
        name,
        show: true
      })
    },
    /**
     * @description 编辑取消
     */
    handleEditNameCancel() {
      this.setData({
        show: false
      })
    },
    handleInput(e) {
      this.data.name = e.detail.value
      console.log(e.detail.value)
    },
    /**
     * @description 编辑确认
     */
    handleEditNameConfirm() {
      const { item, name } = this.data
      if (!name) {
        wx.showToast('请输入文件名称')
        return
      }
      // todo: 名称不能全是空格
      if (!name.trim()) {
        wx.showToast('文件名称不能全是空格')
        return
      }
      // todo: 验证图片名称，不能是表情符
      if (/[\ud800-\udbff][\udc00-\udfff]/g.test(name)) {
        wx.showToast('文件名称不能包含表情符')
        return
      }
      item.name = name

      this.triggerEvent('change', {
        item,
        source: 'editName'
      })
      this.setData({
        show: false
      })
    },
    /**
     * @description 使用原图
     */
    handleUseOriginal() {
      const { item } = this.data

      item.is_origin = item.is_origin === 1 ? 0 : 1

      this.setData({
        'item.is_origin': item.is_origin
      })
      this.triggerEvent('change', {
        item,
        source: 'useOriginal'
      })
    },
    /**
     * @description 删除图片
     */
    handleDelete() {
      const { item } = this.data
      this.triggerEvent('delete', item)
    },
    handleKeyboardHeightChange(e) {
      this.setData({
        keyboardHeight: e.detail.height + 20
      })
    },
    /**
     * @description 使用AI识别
     */
    handleUseAI() {
      const { item } = this.data
      if (this.data.percent < 1) return

      item.is_ai = item.is_ai === 1 ? 0 : 1

      this.setData({
        'item.is_ai': item.is_ai
      })
      this.triggerEvent('change', {
        item,
        source: 'useAI'
      })
    }
  }
})

