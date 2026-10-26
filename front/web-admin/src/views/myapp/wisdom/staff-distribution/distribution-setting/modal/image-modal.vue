<template>
  <el-dialog
    title="海报标背景上传示例"
    :visible.sync="dialogVisible"
    width="390px"
    :show-close="false"
    :before-close="() => closeImage()"
  >
    <span class="tips">支持jpg、png，请上传尺寸750*980px，图片&lt;500kb</span>
    <img :src="imgs[type]" alt="" class="images" />
    <div slot="footer">
      <el-button type="primary" @click="closeImage">知道了</el-button>
    </div>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from "vue-property-decorator";
@Component
export default class ImageModal extends Vue {
  @Prop(String) type!: string;
  @Prop({ default: false }) dialogVisible?: boolean;
  private imgs: any = {
    example: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/distribution_example.png"
  };

  @Emit("close")
  private closeImage() {
    return false;
  }
}
</script>

<style lang="scss" scoped>
.tips {
  font-size: 12px;
  font-family: PingFangSC-Regular, PingFang SC;
  font-weight: 400;
  color: $color-text-second;
  line-height: 21px;
}
.images {
  width: 335px;
  height: 375px;
}
</style>
