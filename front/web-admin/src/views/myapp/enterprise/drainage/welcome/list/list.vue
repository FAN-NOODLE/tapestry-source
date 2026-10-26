<template>
  <div class="welcome_container">
    <el-tabs v-model="activeName" type="card">
      <el-tab-pane label="默认欢迎语" name="1">
        <welcome-list v-if="activeName === '1'"></welcome-list>
      </el-tab-pane>
      <el-tab-pane v-if="ep.isCoach" label="投放渠道欢迎语" name="3">
        <channel-list v-if="activeName === '3'"></channel-list>
      </el-tab-pane>
      <el-tab-pane v-if="ep.isCoach || ep.isKS" label="个性化欢迎语" name="2">
        <personalise-list v-if="activeName === '2'"></personalise-list>
      </el-tab-pane>
      <el-tab-pane label="联系我欢迎语" name="4">
        <contact-list v-if="activeName === '4'"></contact-list>
      </el-tab-pane>
    </el-tabs>
    <div v-if="ep.isCoach" class="tip">
      <el-tooltip effect="light">
        <i class="el-icon-question"></i>
         <div slot="content">
          <p style="font-size: 14px; font-weight: bold;">欢迎语优先级说明</p>
          <p style="font-size: 14px; margin-top: 5px;">特定场景欢迎语（如QR code活动、联系我、获客助手）各自独立互不影响＞个性化欢迎语＞默认欢迎语</p>
          <p style="font-size: 12px; color: #ccc; margin-top: 5px;">*活动欢迎语可在创建活动时配置</p>
         </div>
      </el-tooltip>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import welcomeList from "../components/welcome-list.vue";
import contactList from "../components/contact-list.vue";
import channelList from "../components/channel-list.vue";
import personaliseList from "../components/personalise-list.vue";
import PerModule from "@/store/modules/permission";
@Component({
  components: {
    welcomeList,
    personaliseList,
    contactList,
    channelList
  }
})
export default class extends Vue {
  private activeName = "1";

  get ep() {
    return PerModule.enterprise
  }
  private beforeRouteEnter(to: any, form: any, next: (vm: any) => void): void {
    next((vm: any) => {
      if (
        form.name === "myapp_enterprise_drainage_welcome_newcreate" ||
        form.name === "myapp_enterprise_drainage_welcome_newedit"
      ) {
        vm.activeName = "2";
      }
      if(form.name === "myapp_enterprise_drainage_welcome_channel-create" || form.name === "myapp_enterprise_drainage_welcome_channel-edit") {
        vm.activeName = "3";
      }
      if(form.name === "myapp_enterprise_drainage_welcome_create" || form.name === "myapp_enterprise_drainage_welcome_edit") {
        vm.activeName = "4";
      }
    });
  }
}
</script>
<style lang="scss" scoped>
.welcome_container {
  background-color: #fff;
  padding: 20px;
  border-radius: 10px;
  position: relative;
  .tip {
    position: absolute;
    top: 20px;
    left: 600px;
    color: #999;
    font-size: 14px;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    i {
      font-size: 16px;
      cursor: pointer;
    }
  }
}

/deep/ .el-tabs__nav {
  .is-active {
    color: #fff;
    background-color: #000000;
  }
}
</style>
