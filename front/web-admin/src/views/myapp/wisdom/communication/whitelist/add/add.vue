<template>
  <layout-form>
    <template slot="content">
      <edit-form ref="editForm"></edit-form>
    </template>
    <template slot="footer">
      <div>
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </div>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'
import EditForm, { IFormInfo } from '../components/edit-form.vue'
import Api from '@/api/wisdom/communication';
@Component({
  components: {
    EditForm
  }
})
export default class WhitelistAdd extends Vue {

  handleCancel() {
    this.$router.back()
  }

  handleSave() {
    (this.$refs.editForm as any).validate((value: IFormInfo) => {
      const { member_list, ...params } = JSON.parse(JSON.stringify(value));
      // 处理参数
      params.member_ids = member_list.map((item: any) => {
        return item.member_id
      })
      Api.saveWhitelistApi(params)
        .then((res: any) => {
          if (res.code == 200) {
            this.$message.success('新建成功')
            this.$router.back()
          } else {
            this.$message.error(res.message);
          }
        })
    })
  }
}
</script>