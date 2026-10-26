
<template>
  <el-dialog
    title="订单明细"
    width="1000px"
    :visible="visible"
    :append-to-body="false"
    :modal-append-to-body="false"
    @open="fetchData"
    @close="handleCancel"
  >
    <div class="dialog-wrapper">
      <el-row class="m-b-10">
        <el-col :span="4">门店代码：{{ storeInfo.store_code }}</el-col>
        <el-col :span="4">门店名称：{{ storeInfo.store_name }}</el-col>
      </el-row>
      <el-table :data="list" stripe>
        <el-table-column label="销售单号" prop="order_sn"></el-table-column>
        <el-table-column label="商品数量" prop="total_quantity"></el-table-column>
        <el-table-column label="创建人" prop="member_name"></el-table-column>
        <el-table-column label="发货时间" prop="shipped_at" width="180">
          <template slot-scope="scope">
            {{ scope.row.shipped_at | formatDate('yyyy/mm/dd hh:mm:ss') }}
          </template>
        </el-table-column>
        <el-table-column label="创建时间" prop="created_at" width="180">
          <template slot-scope="scope">
            {{ scope.row.created_at | formatDate('yyyy/mm/dd hh:mm:ss') }}
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="meta.page"
        :page-sizes="[10, 20, 30, 40]"
        :page-size="meta.size"
        class="page-right font-weight-400"
        layout="total, sizes, prev, pager, next, jumper"
        :total="meta.total"
      ></el-pagination>
    </div>
    <div slot="footer" class="dialog-footer">
      <el-button @click="handleCancel">关 闭</el-button>
    </div>
  </el-dialog>
</template>
<script lang="ts">
import { Component, Vue, Model, Prop, Watch } from 'vue-property-decorator'
import Api from '@/api/customer/online-order';

@Component
export default class BatchUpdate extends Vue {
  @Model('update:visible', { type: Boolean, default: false }) visible!: boolean;
  @Prop({ type: Object, default: () => ({}) }) storeInfo: any;

  /** 分页 */
  meta = {
    page: 1,
    size: 10,
    total: 0
  }
  /** 列表数据 */
  list: any[] = []

  @Watch('visible')
  visibleChange(val: boolean) {
    if (val) {
      this.fetchData();
    }
  }

  async fetchData() {

    const params = {
      page: this.meta.page,
      size: this.meta.size,
      third_id: this.storeInfo.third_id,
      shipped_time: this.storeInfo.shipped_time
    }
    const res:any = await Api.getStoreOrderListApi(params)
    if (res.code == 200) {
      this.list = res.data.list;
      this.meta.total = res.data.total;
    } else {
      this.$message.error(res.message);
    }
  }

  /**
   * 每页数量变化
   * @param val 
   */
   handleSizeChange(val: number) {
    this.meta.size = val;
    this.meta.page = 1;
    this.fetchData();
  }

  /**
   * 当前页
   * @param val 
   */
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.fetchData();
  }
  /**
   * 取消
   */
  handleCancel() {
    this.$emit('update:visible', false)
  }

  /**
   * 确定
   */
  handleConfirm() {
    this.$emit('update:visible', false)
  }
}
</script>
<style lang="scss" scoped>
.el-pagination {
  margin-top: 20px;
  text-align: right;
}
</style>