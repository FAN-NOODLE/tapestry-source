<template>
  <layout-form>
    <template slot="content">
      <edit-form ref="editForm" readonly :isEdit="isEdit" :data="formInfo">
        <el-button v-if="formInfo.status === 3" class="m-l-10" type="primary" size="small" @click="handleEdit">编辑</el-button>
      </edit-form>
    </template>
    <template slot="footer">
      <div>
        <el-button @click="handleCancel">取消</el-button>
        <el-button v-if="isEdit" type="primary" @click="handleSave">保存</el-button>
      </div>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'
import EditForm, { IFormInfo, defaultFormInfo } from '../components/edit-form.vue'
import Api from '@/api/customer/online-order';
@Component({
  components: {
    EditForm
  }
})
export default class OrderDetail extends Vue {
  /** 表单数据 */
  formInfo: IFormInfo = JSON.parse(JSON.stringify(defaultFormInfo))
  isEdit = false

  created() {
    this.getDetail()
  }

  handleEdit() {
    this.isEdit = true
  }

  /**
   * 获取详情
   */
  getDetail() {
    Api.getOnlineOrderDetailApi({
      order_id: Number(this.$route.query.id)
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

  handleSave() {
    (this.$refs.editForm as any).validate(async(value: IFormInfo) => {
      const params: IFormInfo = JSON.parse(JSON.stringify({...value}));
      const { id: order_id, tracking_number, transfer_file  } = params
      
      const res: any = await Api.editOnlineOrderApi({
        order_id,
        tracking_number,
        transfer_file,
      })

      if (res.code == 200) {
        this.$message.success('保存成功')
        this.$router.back()
      } else {
        this.$message.error(res.message);
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