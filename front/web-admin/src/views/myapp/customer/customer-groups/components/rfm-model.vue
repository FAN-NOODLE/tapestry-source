<template>
  <div class="rfm-layout">
    <el-dialog
      title="什么是RFM模型"
      width="636px"
      :before-close="handleClose"
      :visible.sync="isShow"
    >
      <dl>
        <img src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/customer/rfm.png" alt="" />
      </dl>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";

@Component
export default class extends Vue {
  @Prop({ default: true }) isShow!: boolean;
  private handleClose() {
    this.$emit("isShow", false);
  }
}
</script>
<style lang="scss" scoped>
.rfm-layout {
  ::v-deep.el-dialog__body {
    padding-top: 0;
  }
  dl {
    overflow-y: auto;
    height: 400px;
    img {
      width: 100%;
    }
  }
}
</style>
