<template>
  <layout-form v-if="formInfo">
    <template slot="content">
      <audit-form v-model="formInfo" readonly></audit-form>
    </template>
    <template slot="footer">
      <div>
        <el-button @click="handleCancel">返回</el-button>
        <el-button type="primary" @click="handleExportImage">一键保存所有图片</el-button>
      </div>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Api from '@/api/wisdom/displayTask';
import AuditForm, { IStoreInfo } from '../components/audit-form.vue'
import { saveImageToZip } from '@/utils/zip'
import { IFormInfo } from '../components/edit-form.vue';

@Component({
  components: {
    AuditForm
  }
})
export default class DisplayTaskAuditDetail extends Vue {
  /** 表单数据 */
  formInfo: IStoreInfo = null
  taskDetail: IFormInfo = null

  async mounted() {
    await this.getTaskDetail()
    this.getDetail()
  }

  getDetail() {
    Api.getAllotInfoApi({
      id: Number(this.$route.query.id),
      store_id: Number(this.$route.query.store_id)
    }).then((res: any) => {
      if(res.code === 200) {
        console.log(res.data);
        
        res.data.template_info = this.taskDetail.template_info.map((item) => {
          const imageList = res.data.image_list
            .filter(image => image.identifier === item.template_name)
            .map(value => {
              value._status = value.status
              value._image_list = []
              value.comment = value.comment.map(comment => {
                const images = comment.content.image_url?.map(url => ({
                  url: url as unknown as string,
                }))||[]
                return {
                  ...comment,
                  _reply_content: '',
                  _reply_name: '',
                  _reply_show: false,
                  _prefix_width: 0,
                  content: images.length ? {image_url: images} : {text: comment.content.text}
                }
              })
              return value
            })
          return {
            ...item,
            image_list: imageList
          }
        })
        res.data.grade = res.data.grade || null
        this.formInfo = res.data
      }
    })
  }

  async getTaskDetail() {
    const res: any = await  Api.getDisplayTaskDetailApi({
      id: Number(this.$route.query.id)
    })
    if (res.code == 200) {
      this.taskDetail = res.data
    } else {
      this.$message.error(res.message);
    }
  }
  /**
   * @description 取消
   */
  handleCancel() {
    this.$router.go(-1)
  }

  /** 
   * @description 一键保存所有图片
   */
  handleExportImage() {
    const { template_info, display_name, store_code, store_name } = this.formInfo
    const loading = this.$loading({
      lock: true,
      text: '正在打包中...',
      spinner: 'el-icon-loading',
      background: 'rgba(0, 0, 0, 0.7)'
    })
    const list = template_info.reduce((prev, item) => {
      prev.push(...item.image_list.filter(item => item.status === 0).map((item, index) => {
        // 获取后缀
        const suffix = item.url.split('.').pop()
        return {
          name: `${item.identifier}_${index}.${suffix}`,
          content: item.url
        }
      }))
      return prev
    }, [])

    saveImageToZip(list, `${display_name}_${store_code}_${store_name}.zip`).then(() => {
      loading.close()
      this.$message.success('下载成功')
    })

  }
}
</script>
<style lang="scss" scoped>
.section {
  &:not(:last-child) {
    margin-bottom: 20px;
  }

  &-title {
    position: relative;
    padding-left: 10px;
    font-size: 16px;
    margin-bottom: 10px;
    &::before {
      content: "";
      width: 5px;
      height: 20px;
      position: absolute;
      top: 50%;
      left: 0;
      transform: translateY(-50%);
      background-color: #000;
    }
  }

  &-row {
    &:not(:last-child) {
      margin-bottom: 10px;
    }
  }

  .btn {

    border: none;
    background-color: #f5f6fa;
    /deep/.el-icon-picture,
    /deep/.el-icon-warning-outline {
      font-size: 16px;
      vertical-align: middle;
    }
  }
}

.upload {
  position: relative;
  .btn-tag {
    position: absolute;
    top: 0;
    left: 105px;
  }
}

.item {
  &:not(:last-child) {
    margin-bottom: 20px;
  }
  .item-user {
    display: flex;
    align-items: center;
    .el-avatar {
      margin-right: 10px;
    }

    .item-user-info {
      line-height: normal;
      flex: 1;
    }
  
    .item-user-name {
      font-weight: bold;
    }
    .item-user-store {
      font-size: 12px;
      color: #999;
      margin-top: 5px;
    }

    .item-user-time {
      color: #999;
      line-height: normal;
    }

    
  }
  .item-content {
    margin-top: 10px;
    .item-text {
      padding: 10px;
      background-color: #f7f7f7;
      line-height: 1.5;
      color: #666;
      position: relative;
      text-align: justify;
      border-radius: 4px;
      .reply {
        height: 20px !important;
        line-height: 20px !important;
        padding: 0 5px!important;  
      }
    }
  }
}
</style>