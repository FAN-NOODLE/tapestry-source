<template>
  <!-- 查看分配情况 -->
  <el-dialog
    title="分配情况"
    :visible.sync="dialogVisible"
    width="60%"
    @close="close"
  >
    <el-tabs v-model="activeName">
      <el-tab-pane
        label="待分配好友"
        name="first"
        v-if="pageTag === 'not'"
      >
        <not-customer-detail
          :handoverUserid="handoverUserid"
          @customerDistribution="customerDistribution"
        />
      </el-tab-pane>
      <el-tab-pane
        label="已分配好友"
        name="second"
      >
        <already-customer-detail :handoverUserid="handoverUserid" />
      </el-tab-pane>
      <el-tab-pane
        label="待分配社群"
        name="third"
        v-if="pageTag === 'not'"
      >
        <not-group-detail
          :handoverUserid="handoverUserid"
          @groupDistribution="groupDistribution"
        />
      </el-tab-pane>
      <el-tab-pane
        label="已分配社群"
        name="fourth"
      >
        <already-group-detail :handoverUserid="handoverUserid" />
      </el-tab-pane>
    </el-tabs>
  </el-dialog>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
// 待分配好友
import notCustomerDetail from "./not-customer-detail.vue";
// 待分配社群
import notGroupDetail from "./not-group-detail.vue";
// 已分配好友
import alreadyCustomerDetail from "./already-customer-detail.vue";
// 已分配社群
import alreadyGroupDetail from "./already-group-detail.vue";
@Component({
  components: {
    notCustomerDetail,
    notGroupDetail,
    alreadyCustomerDetail,
    alreadyGroupDetail,
  },
})
export default class extends Vue {
  @Prop({ default: "" }) handoverUserid!: string;
  @Prop({ default: "not" }) pageTag!: string;
  private dialogVisible = true;
  private activeName = "first";
  private close() {
    this.$emit("close");
  }
  // 分配好友
  private customerDistribution() {
    this.$emit("customerDistribution");
  }
  private groupDistribution() {
    this.$emit("groupDistribution");
  }
  // 生命周期
  created() {
    if (this.pageTag === "already") {
      this.activeName = "second";
    }
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-dialog {
  max-width: 1200px !important;
}
</style>
