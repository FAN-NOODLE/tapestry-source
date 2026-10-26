<template>
  <div></div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/setting/setting";
@Component({})
export default class extends Vue {
  private loadThenJump() {
    this.$loading({
      lock: true,
      text: "Loading",
      spinner: "el-icon-loading",
      background: "rgba(0, 0, 0, 0.4)"
    });
    const urls = this.$router.resolve({
      path: "/myapp/setting/authorization-settings/store-set/store-sets"
    }).href;
    const back_url = window.location.origin + urls;
    console.log("-----------------");
    console.log("back_url", back_url);
    console.log("this.$route.query", this.$route.query);
    console.log("this.$route.query.back", this.$route.query.back);
    window.location.href =
      this.$route.query.back +
      "?enterprise_id=" +
      this.$route.query.enterprise_id +
      "&auth_for=" +
      this.$route.query.auth_for +
      "&auth_code=" +
      this.$route.query.auth_code +
      "&back=" +
      back_url;
  }
  mounted() {
    this.loadThenJump();
  }
}
</script>

<style lang="scss" scoped></style>
