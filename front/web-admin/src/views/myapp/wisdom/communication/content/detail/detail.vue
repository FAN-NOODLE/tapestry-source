<template>
  <layout-form v-if="formInfo">
    <template slot="content">
      <detail-form :data="formInfo"></detail-form>
      <temp-one ref="detail" :data="formInfo"></temp-one>
    </template>
    <template slot="footer">
      <el-button @click="handleCancel">返回</el-button>
      <el-button v-if="formInfo.status===4" type="primary" @click="handleDownLoad">下载文件</el-button>
      <el-button v-if="formInfo.status===3" type="primary" @click="handleEdit">再次编辑</el-button>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'
import Api from '@/api/wisdom/communication';
import { statusOptions, typeOptions } from '../../const';
import { IFormInfo, defaultFormInfo } from '../../components/edit-form.vue'
import DetailForm from '../../components/detail-form.vue'
import TempOne from '../../components/temp-one.vue'

@Component({
  components: {
    DetailForm,
    TempOne
  }
})
export default class ContentDetail extends Vue {
  /** 表单数据 */
  formInfo: IFormInfo = null

  mounted() {
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

        res.data.type_name = typeOptions.find(item => item.value === res.data.type).label
        res.data.status_name = statusOptions.find(item => item.value === res.data.status).label
        res.data.page_url = res.data.page_url !== '2'? '1': '2'

        this.formInfo = res.data
      } else {
        this.$message.error(res.message);
      }
    })
  }

  handleCancel() {
    this.$router.back()
  }
  /**
   * 编辑
   */
  handleEdit() {
    this.$router.push({
      name: 'myapp_wisdom_communication_content_edit',
      query: {
        id: String(this.formInfo.id)
      }
    })
  }

  /**
   * 下载文件
   */
  handleDownLoad(){
    // (this.$refs.detail as any).handleDownloadPdf()
    if(this.formInfo.pdf_url) {
      (this.$refs.detail as any).handleDownloadFile()
    } else {
      (this.$refs.detail as any).handleDownloadPdf()
    }
  }

}
</script>