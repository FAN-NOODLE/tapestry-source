<template>
  <div ontouchstart="">
    <!-- ontouchstart="" 解决按钮组件ios点击没有阴影反馈效果 -->
    <keep-alive>
      <router-view v-if="$route.meta.keepAlive"></router-view>
    </keep-alive>
    <router-view v-if="!$route.meta.keepAlive"></router-view>
    <!-- 水印 -->
    <watermark :opacity="0.2"></watermark>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Watermark from '@/components/base/watermark.vue'
@Component({
  components: {
    Watermark
  }
})
export default class extends Vue {
}
</script>
