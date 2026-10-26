<template>
  <div class="version-prompt">
    <el-dialog
      title="WeSocial 系统版本更新提示"
      :visible.sync="versionDialogVisible"
      width="400px"
      center
      :close-on-click-modal="false"
      :show-close="false"
    >
      <div class="dialog-content">
        <h3>{{ version_info.them }}</h3>
        <h4 class="time">{{ version_info.date }}</h4>
        <div
          class="content-item"
          v-for="(item, index) in version_info.detail"
          :key="index"
        >
          <i class="el-icon-check" aria-hidden="true"></i>
          {{ item.title }}
          <p class="item-desc" v-for="(items, i) in item.detail" :key="i">
            {{ items }}
          </p>
        </div>
      </div>
      <div class="dialog-footer">
        <el-button type="primary" size="small" @click="cancel">关闭</el-button>
      </div>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import UserModule from "@/store/modules/user";
import Api from "@/api/home/home";
@Component
export default class extends Vue {
  // 获取用户信息
  get USER() {
    return UserModule.getUserInfo;
  }
  private versionDialogVisible = false;
  private version_info = {};
  private cancel() {
    this.versionDialogVisible = false;
  }
  private mounted() {
    this.getCurrentVersionInfo();
  }
  // 获取当前版本信息
  private getCurrentVersionInfo() {
    const params = {
      type: 1,
      admin_id: this.USER.admin_id
    };
    Api.currentVersionInfo(params).then((res: any) => {
      if (res.data) {
        this.versionDialogVisible = !res.data.is_has_show;
        this.version_info = res.data.version_info;
      }
    });
  }
}
</script>
<style lang="scss" scoped>
.dialog-content {
  padding: 0 20px;
  h3,
  h4 {
    text-align: center;
  }
  .time {
    margin-top: 10px;
  }
  .content-item {
    margin-top: 20px;
    .item-desc {
      font-size: 12px;
    }
  }
}
.dialog-footer {
  width: 100%;
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
