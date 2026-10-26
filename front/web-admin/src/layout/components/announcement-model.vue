<template>
  <el-dialog
    class="content-dialog"
    :visible.sync="modelShow"
    :close-on-press-escape="false"
    :close-on-click-modal="false"
    width="560px"
  >
    <div class="template">
      <div class="template-title">{{ announcement.title }}</div>
      <div class="template-time">
        {{ announcement.released_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
      </div>
      <p class="template-content" v-html="announcement.content"></p>
    </div>
    <el-button
      class="close"
      type="primary"
      size="small"
      @click="modelShow = false"
      >知道了</el-button
    >
  </el-dialog>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import Bus from "@/components/bus";
import { AnnounceItem } from "./interface";
import { trimEnter } from "@/utils/tools";
@Component({})
export default class AnnouncementModel extends Vue {
  private modelShow = false;
  private announcement = {};
  beforeCreate() {
    Bus.$on("showAnnounceMent", (data: AnnounceItem) => {
      const { content } = data;
      Object.assign(data, { content: trimEnter(content) });
      this.announcement = data;
      this.modelShow = true;
    });
  }
}
</script>
<style lang="scss" scoped>
.template {
  .template-title {
    font-weight: 600;
    color: #292a30;
    font-size: 18px;
    line-height: 26px;
    text-align: center;
    margin-bottom: 10px;
  }
  .template-time {
    color: #909399;
    font-size: 14px;
    line-height: 20px;
    text-align: center;
    margin-bottom: 30px;
  }
  .template-content {
    height: 200px;
    text-align: left;
    overflow-y: auto;
    color: #606266;
    font-size: 14px;
    line-height: 20px;
  }
}
.close {
  position: absolute;
  bottom: 30px;
  right: 30px;
}
/deep/ {
  .el-dialog {
    height: 397px !important;
  }
  .el-dialog__body {
    padding: 0 30px 30px 30px;
  }
}
</style>
