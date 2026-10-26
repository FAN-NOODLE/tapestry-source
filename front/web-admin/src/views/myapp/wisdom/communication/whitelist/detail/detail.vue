<template>
  <layout-form>
    <template slot="content">
      <edit-form ref="editForm" :data="formInfo" readonly></edit-form>
    </template>
    <template slot="footer">
      <div>
        <el-button @click="handleCancel">返回</el-button>
      </div>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'
import EditForm, { IFormInfo, defaultFormInfo } from '../components/edit-form.vue'
import Api from '@/api/wisdom/communication';
@Component({
  components: {
    EditForm
  }
})
export default class WhitelistDetail extends Vue {
  /** 表单数据 */
  formInfo: IFormInfo = JSON.parse(JSON.stringify(defaultFormInfo))

  mounted() {
    this.getDetail()
  }

  /**
   * 获取详情
   */
  getDetail() {
    Api.getWhitelistDetailApi({
      id: Number(this.$route.query.id)
    }).then((res: any) => {
      if (res.code == 200) {
        this.formInfo = res.data
      } else {
        this.$message.error(res.message);
      }
    })
  }

  handleCancel() {
    this.$router.back()
  }

}
</script>