<template>
  <default-page :picture="picture" :text="text"></default-page>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import defaultPage from "@/components/base/empty-page.vue";
@Component({
  name: "no-auth",
  components: {
    defaultPage,
  },
})
export default class extends Vue {
  // 缺省图片
  picture:any = "system";
  // 缺省文字
  text:any = "拒绝外部人员访问";
  mounted() {
    // 获取缺省图片
    if (this.$route.query.picture) {
      this.picture = this.$route.query.picture;
    }
    // 获取缺省文字
    if (this.$route.query.text) {
      this.text = this.$route.query.text;
    }
  }
}
</script>
