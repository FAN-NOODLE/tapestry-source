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
import { saveActivityApi } from '@/api/enterprise/activity';
@Component({
  components: {
    EditForm
  }
})
export default class LaunchAdd extends Vue {

  handleCancel() {
    this.$router.back()
  }

  handleSave() {
    (this.$refs.editForm as any).validate((value: IFormInfo) => {

      // 处理参数
      saveActivityApi(value)
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