<template>
  <layout-form>
    <template slot="content">
      <edit-form ref="editForm" :data="formInfo"></edit-form>
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
import EditForm, { IFormInfo, defaultFormInfo } from '../components/edit-form.vue'
import Api from '@/api/customer/online-order';
@Component({
  components: {
    EditForm
  }
})
export default class DeliverOrder extends Vue {
  /** 表单数据 */
  formInfo: IFormInfo = JSON.parse(JSON.stringify(defaultFormInfo))

  created() {
    this.getDetail()
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
      const { id: order_id, tracking_number, sku_list, transfer_file  } = params
      const refusal_info = sku_list.filter((item: any) => item.refusal_quantity > 0).map(item => {
        return {
          sku_info: item.sku_info,
          quantity: item.refusal_quantity
        }
      })

      const quantity = sku_list.reduce((total, item) => total + item.quantity, 0)
      const refusal_quantity = sku_list.reduce((total, item) => total + item.refusal_quantity, 0)

      await this.$confirm(`订单商品数量:${quantity} 拒单商品数量:${refusal_quantity}`, '请确认订单发货情况', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })

      const res: any = await Api.submitOnlineOrderApi({
        order_id,
        tracking_number,
        transfer_file,
        refusal_info
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