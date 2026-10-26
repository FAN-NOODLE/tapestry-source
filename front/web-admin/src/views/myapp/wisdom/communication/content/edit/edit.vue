<template>
  <layout-form>
    <template slot="content">
      <edit-form ref="editForm" :data="formInfo"></edit-form>
    </template>
    <template slot="footer">
      <div>
        <el-button @click="handleCancel">取消</el-button>
        <el-button @click="handleSave(1)">保存草稿</el-button>
        <el-button type="primary" @click="handleSave(0)" v-loading="loading">保存并提交</el-button>
        <!-- <el-button type="primary" @click="handlePreview()">预览PDF文件</el-button> -->
      </div>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'
import EditForm, { IFormInfo, defaultFormInfo } from '../../components/edit-form.vue'
import Api from '@/api/wisdom/communication';
@Component({
  components: {
    EditForm
  }
})
export default class ContentEdit extends Vue {
  /** 表单数据 */
  formInfo: IFormInfo = JSON.parse(JSON.stringify(defaultFormInfo))
  config = null
  loading = false
  get message() {
    return (type:number|string) => {
      // 获取今天是周几
      const week = new Date().getDay()
      const weeks = ['','一', '二', '三', '四', '五', '六', '日']
      const value = this.config.find(item => item.info_type === type)
      if(!value) return ''
      const hours: number[] = value.deadline_time.split(':').map(Number)
      const date = new Date().setHours(hours[0], hours[1], hours[2])

      if(value.deadline_week < week || value.deadline_week === week &&Date.now() > date) {
        // 已过截止日期
        return `周${weeks[value.deadline_week]} ${value.deadline_time}前`
      } else {
        return ''
      }
    }
  }

  created() {
    this.getDetail()
    this.getData()
  }

  /**
   * 获取详情
   */
  getDetail() {
    Api.getContentDetailApi({
      id: Number(this.$route.query.id)
    }).then((res: any) => {
      if (res.code == 200) {
        res.data.page_url = res.data.page_url !== '2'? '1': '2'

        res.data.attachment = res.data.attachment.map((item: any) => {
          return {
            ...item,
            _progress: 100,
            _status: 'success'
          }
        })

        if(this.$route.query.copy) {
          res.data.attachment = []
          res.data.date_time = ''
          res.data.status = 0
          delete res.data.id
        }
        this.formInfo = res.data
      } else {
        this.$message.error(res.message);
      }
    }).catch(() => {
      this.loading = false
    })
  }


  getData() {
    Api.getNoticeConfigApi({}).then((res: any) => {
      this.config = res.data
    })
  }

  handleCancel() {
    this.$router.back()
  }

  handleSave(is_draft: number) {
    (this.$refs.editForm as any).validate(async(value: IFormInfo) => {
      const params: IFormInfo = JSON.parse(JSON.stringify({...value, copy: this.$route.query.copy}));

      const message = this.message(Number(params.type)-1)
      if(is_draft === 0 && message) {
        await this.$confirm(`该类型传讯内容本周收集时间（${message}）已截止，当前提交内容可能于下周发布，且审核期间内容不可修改，确定继续提交？`, "确认提交", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
        });
      }
      this.loading = true
      
      try {
        params.contact_info = params.contact_info.filter((item) => {
          return Object.values(item).some(val => val)
        })
        
        params.pdf_url = (await (this.$refs.editForm as any).handleCreate()).url
        
        const res: any = await Api.saveContentApi({...params, is_draft})

        this.loading = false
        if (res.code == 200) {
          this.$message.success('保存成功')
          this.$router.back()
        } else {
          this.$message.error(res.message);
        }
      } catch (error) {
        this.loading = false
      }
      
    })
  }

  /**
   * 预览
   */
  handlePreview() {
    (this.$refs.editForm as any).handlePreview()
  }
}
</script>