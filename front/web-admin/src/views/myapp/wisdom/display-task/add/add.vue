<template>
  <layout-form>
    <template slot="content">
      <edit-form ref="editForm"></edit-form>
    </template>
    <template slot="footer">
      <div>
        <el-button type="primary" @click="handleSave(1)">保存并上架</el-button>
        <el-button @click="handleSave(0)">仅保存</el-button>
        <el-button @click="handleCancel">取消</el-button>
      </div>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import EditForm, {IFormInfo} from '../components/edit-form.vue'
import Api from '@/api/wisdom/displayTask'
@Component({
  components: {
    EditForm,
  }
})
export default class DisplayTaskAdd extends Vue {
  /** 
   * 保存
   */
  private handleSave(show_type: number) {
    (this.$refs.editForm as any).validate((value: IFormInfo) => {
      const params: IFormInfo = JSON.parse(JSON.stringify(value));
      // 过滤空值
      params.class_url = params.class_url.filter(item => item.url||item.class_name)
      // 过滤空值
      params.template_info = params.template_info.filter(item => item.template_name)
      // 保存状态
      params.show_type = show_type
      if(params.store_type === 0) {
        params.store_list = []
      }
      Api.addDisplayTaskApi(params)
        .then((res: any) => {
          if (res.code == 200) {
            this.$message.success("保存成功")
            this.$router.go(-1)
          } else {
            this.$message.error(res.message);
          }
        })
    })
  }
  /**
   * @description 取消
   */
  private handleCancel() {
    this.$router.go(-1)
  }
}
</script>