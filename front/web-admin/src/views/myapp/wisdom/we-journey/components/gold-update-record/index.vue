<!-- 金币修改记录弹窗 -->
<template>
  <el-dialog
    title="金币修改记录"
    :visible.sync="goldUpdateRecordVisible"
    width="900px"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    @close="closeGoldUpdateRecord"
  >
    <el-table
      :data="coins_settings_edit_record"
      stripe
      style="margin: 0 auto;"
      max-height="400"
    >
      <el-table-column
        prop="id"
        label="记录ID"
        width="180"
      ></el-table-column>
      <el-table-column
        prop="ratio_remark"
        label="修改内容"
        width="180"
      ></el-table-column>
      <el-table-column
        prop="admin_name"
        label="修改人"
        width="210"
      ></el-table-column>
      <el-table-column
        prop="updated_at"
        label="修改时间"
      >
        <template slot-scope="scope">
          {{ scope.row.updated_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
        </template>
      </el-table-column>
    </el-table>
  </el-dialog>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator"
@Component({})
export default class GoldUpdateRecord extends Vue {
  [x: string]: any;
  @Prop({ default: false }) readonly goldUpdateRecordVisible!: boolean
  @Prop({ default: [] }) readonly coins_settings_edit_record!: any
  // 关闭弹窗
  closeGoldUpdateRecord() {
    this.$emit("closeGoldUpdateRecord")
  }
}

</script>
<style lang="scss" scoped>
.pagination {
  text-align: right;
  margin-top: 20px;
}
</style>
