<template>
  <layout-form>
    <template slot="content">
      <edit-form ref="editForm" :data="formInfo"></edit-form>

      <el-dialog
        :title="['审核不通过', '审核通过'][auditInfo.type]"
        :visible.sync="visible"
        width="500px"
        center
      >
        <el-form label-width="60px" style="width: 80%;margin:auto">
          <div class="m-b-10">
            <i class="el-icon-warning m-r-5"></i> 
            <span v-if="auditInfo.type === 0">确定当前内容审核不通过？</span>
            <span v-else>确定当前内容审核通过？确定后将不可再修改编辑该传讯内容？</span>
          </div>
          <el-form-item v-if="auditInfo.type === 0" label="留言:">
            <el-input
              v-model="auditInfo.comment"
              type="textarea"
              :rows="4"
              placeholder="请留下您的意见"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="handleAuditCancel">取 消</el-button>
            <el-button type="primary" @click="handleAuditConfirm" v-loading="loading">确 定</el-button>
          </div>
        </template>
      </el-dialog>
    </template>
    <template slot="footer">
      <div>
        <el-button @click="handleCancel">取消</el-button>
        <el-button @click="handleAudit(0)">审核不通过</el-button>
        <el-button type="primary" @click="handleAudit(1)">审核通过</el-button>
      </div>
      <!-- 审核弹框 -->
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
export default class AuditingEdit extends Vue {
  /** 表单数据 */
  formInfo: IFormInfo = JSON.parse(JSON.stringify(defaultFormInfo))
  /** 审核弹框控制 */
  visible = false
  /** 审核文案 */
  auditInfo = {
    type: 0,
    comment: ''
  }
  loading = false

  created() {
    this.getDetail()
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
        
        this.formInfo = res.data
      } else {
        this.$message.error(res.message);
      }
    })
  }

  handleCancel() {
    this.$router.back()
  }

  handleSave() {
    (this.$refs.editForm as any).validate(async(value: IFormInfo) => {

      const loading = this.$loading({
        lock: true,
        text: '保存中...',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)',
      });

      try {


        const params: IFormInfo = JSON.parse(JSON.stringify({...value}));
        params.status = [3, 4][this.auditInfo.type]
        params.comment = this.auditInfo.comment.trim()
        params.contact_info = params.contact_info.filter((item) => {
          return Object.values(item).some(val => val)
        })
        // 审核通过的时候重新生成pdf
        if(this.auditInfo.type == 1) {
          params.pdf_url = (await (this.$refs.editForm as any).handleCreate()).url
        }

        const res: any = await Api.auditContentApi({...params})

        loading.close()
        if (res.code == 200) {
          this.$message.success('保存成功')
          this.$router.back()
        } else {
          this.$message.error(res.message);
          this.getDetail()
        }
      } catch (error) {
        loading.close()
      }
      
    })
  }

  handleAudit(type: number) {
    this.auditInfo.type = type
    this.auditInfo.comment = ''
    this.visible = true
  }

  handleAuditCancel() {
    this.visible = false
  }

  handleAuditConfirm() {
    this.visible = false
    this.handleSave()
  }
}
</script>

<style lang="scss" scoped>
.tips {
  text-align: center;
}
</style>