<template>
  <div>
    <!-- 待分配导购列表 -->
    <el-table
      stripe
      :data="guideList"
      class="mar-t-20"
    >
      <el-table-column
        prop="handover_userid"
        label="导购工号"
        width="180"
      ></el-table-column>
      <el-table-column
        prop="handover_name"
        label="导购姓名"
        min-width="140"
      ></el-table-column>
      <el-table-column
        prop="handover_mobile"
        label="手机号"
        width="160"
      ></el-table-column>
      <el-table-column
        prop="handover_department_name"
        label="所在门店"
        min-width="160"
      ></el-table-column>
      <el-table-column
        prop="unassigned_user_num"
        label="待分配好友数"
        min-width="140"
      ></el-table-column>
      <el-table-column
        prop="assigned_user_num"
        label="已分配好友数"
        min-width="140"
      ></el-table-column>
      <el-table-column
        prop="unassigned_group_num"
        label="待分配社群数"
        min-width="140"
      ></el-table-column>
      <el-table-column
        prop="assigned_group_num"
        label="已分配社群数"
        min-width="140"
      ></el-table-column>
      <el-table-column
        label="操作"
        min-width="290"
        fixed="right"
      >
        <template slot-scope="scope">
          <el-button
            type="text"
            size="small"
            :disabled="scope.row.unassigned_user_num <= 0"
            @click="openDialog('isDc', scope.row.handover_userid)"
          >分配好友</el-button>
          <el-button
            type="text"
            size="small"
            :disabled="scope.row.unassigned_group_num <= 0"
            @click="openDialog('isDg', scope.row.handover_userid)"
          >分配社群</el-button>
          <el-button
            type="text"
            size="small"
            @click="openDialog('isCheckDetail', scope.row.handover_userid)"
          >查看分配情况</el-button>
        </template>
      </el-table-column>
      <div
        slot="empty"
        v-if="!guideList.length"
      >
        <layout-empty :message="'暂无导购'"></layout-empty>
      </div>
    </el-table>
    <distributionCustomer
      v-if="isDc"
      :handoverUserid="handover_userid"
      @submit="submit"
      @cancel="cancelDialog('isDc')"
    ></distributionCustomer>
    <distributionGroup
      v-if="isDg"
      :handoverUserid="handover_userid"
      @cancel="cancelDialog('isDg')"
      @submit="submit"
    ></distributionGroup>
    <distributionDetail
      v-if="isCheckDetail"
      :handoverUserid="handover_userid"
      @close="cancelDialog('isCheckDetail')"
      @customerDistribution="customerDistribution"
      @groupDistribution="groupDistribution"
    ></distributionDetail>
  </div>
</template>
<script lang="ts">
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
import { Component, Vue, Prop } from "vue-property-decorator";
// 分配好友
import distributionCustomer from "./distribution-customer.vue";
// 分配社群
import distributionGroup from "./distribution-group.vue";
// 查看分配情况
import distributionDetail from "./distribution-detail.vue";
import Api from "@/api/wisdom/quit-inherit";
@Component({
  components: {
    LayoutEmpty,
    distributionCustomer,
    distributionDetail,
    distributionGroup,
  },
})
export default class extends Vue {
  @Prop({ default: "" }) searchKey!: string;
  @Prop({ default: [] }) guideList!: Array<any>;
  private isDc = false; // 分配好友弹窗
  private isDg = false; // 分配社群弹窗
  private isCheckDetail = false; // 查看分配情况
  private handover_userid = 0;
  private openDialog(params: any, userId?: any) {
    if (userId) {
      this.handover_userid = userId;
    }
    this[params] = true;
  }
  private cancelDialog(params: any) {
    this[params] = false;
  }
  private customerDistribution() {
    this.openDialog("isDc");
  }
  private groupDistribution() {
    this.openDialog("isDg");
  }
  // 分配好友与分配社群提交成功
  private submit() {
    this.$emit("submit");
    this.cancelDialog("isCheckDetail");
  }
  // 生命周期
  mounted() {
    // this.allotExternalList();
  }
}
</script>
<style lang="scss" scoped></style>
