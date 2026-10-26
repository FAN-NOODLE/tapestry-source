<template>
  <el-image v-if="innerSrc" :src="innerSrc" fit="cover" :preview-src-list="previewList" v-bind="$attrs" v-on="$listeners">
    <template #error>
      <slot name="error"></slot>
    </template>
  </el-image>
</template>
<script lang="ts">
import { VUE_APP_API_ROOT_RE } from "@/utils/env-config";
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
@Component
export default class FImage extends Vue {
  @Prop({ type: String, default: "" }) src: string;
  @Prop({ type: Array, default: () => [] }) previewSrcList: string[];

  innerSrc = '';
  previewList:string[] = []

  @Watch('src', { immediate: true })
  handleSrcChange(val:string) {
    if(
      val.indexOf('https://tapestry-uat-1308608494.cos.ap-shanghai.myqcloud.com') != -1 || 
      val.indexOf('https://cdn-cos.coach.com.cn') != -1
    ) {
      const encoder = new TextEncoder();
      const encodedData = encoder.encode(val);
      // 使用 btoa 编码字节数组为 Base64 字符串
      const encodedText = btoa(String.fromCharCode(...encodedData));
      // 获取随机2位字符串
      const randomStr = Math.random().toString(36).substring(2, 4);
      this.innerSrc = `${VUE_APP_API_ROOT_RE}common/get-file-by-url?url=${randomStr}${encodedText}`;

    } else {
      this.innerSrc = val
    }
  }

  @Watch('previewSrcList', { immediate: true, deep: true })
  handlePreviewChange(val) {
    this.previewList = val.map(item => {
      if(
        item.indexOf('https://tapestry-uat-1308608494.cos.ap-shanghai.myqcloud.com') != -1 || 
        item.indexOf('https://cdn-cos.coach.com.cn') != -1
      ) {
        const encoder = new TextEncoder();
        const encodedData = encoder.encode(item);
        // 使用 btoa 编码字节数组为 Base64 字符串
        const encodedText = btoa(String.fromCharCode(...encodedData));
        // 获取随机2位字符串
        const randomStr = Math.random().toString(36).substring(2, 4);
        
        return `${VUE_APP_API_ROOT_RE}common/get-file-by-url?url=${randomStr}${encodedText}`;

      } else {
        return item
      }
    })
  }
}
</script>
<style lang="scss" scoped>
.search-item {
  display: inline-block;
  color: $color-text-normal;
  margin-right: 20px;
  margin-bottom: 15px;
}
</style>