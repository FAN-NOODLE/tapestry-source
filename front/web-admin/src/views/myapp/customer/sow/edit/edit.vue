<template>
  <layout-form>
    <template slot="content">
      <edit-form ref="editForm" :data="formInfo"></edit-form>
    </template>
    <template slot="footer">
      <div>
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" @click="handleSave(2)">保存</el-button>
        <el-button type="primary" @click="handleSave(1)">保存并上架</el-button>
      </div>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'
import EditForm, { IFormInfo, defaultFormInfo } from '../components/edit-form.vue'
import Api from '@/api/customer/sow';
@Component({
  components: {
    EditForm
  }
})
export default class SowEdit extends Vue {
  /** 表单数据 */
  formInfo: IFormInfo = JSON.parse(JSON.stringify(defaultFormInfo))

  mounted() {
    this.getDetail()
  }

  /**
   * 获取详情
   */
  getDetail() {
    Api.getSowDetailApi({
      id: Number(this.$route.query.id)
    }).then((res: any) => {
      if(res.code == 200) {
        res.data.sow_template_notice = res.data.notice_list
        this.formInfo = res.data
      }
    })
  }

  handleCancel() {
    this.$router.back()
  }

  handleSave(status: number) {
    (this.$refs.editForm as any).validate((value: IFormInfo) => {
      const params = JSON.parse(JSON.stringify(value));
      params.status = status
      Api.saveSowApi(params).then((res: any) => {
        if(res.code === 200) {
          this.$message.success('保存成功')
          this.$router.back()
        }
      })
    })
  }
}
</script>