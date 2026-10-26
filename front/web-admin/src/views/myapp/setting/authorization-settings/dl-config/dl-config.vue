<template>
  <div class="config">
    <!-- DL配置 -->
    <div class="box">
      <!-- Access Token管理 -->
      <div class="module">
        <h4>Access Token管理</h4>
        <div class="desc">
          请前往Datalake对接人员获取
        </div>
        <div class="row" v-for="(value, key) in tokenList" :key="key">
          <span class="napa">{{ key }}</span>
          <span>{{ value || "--" }}</span>
          <el-button type="text" @click="openDialogDL(key, value)"
            >修改</el-button
          >
        </div>
      </div>
      <div class="line"></div>
      <!-- 品牌信息配置 -->
      <div class="module">
        <h4>品牌信息配置</h4>
        <div class="desc">
          请前往Datalake对接人员获取
        </div>
        <div class="row" v-for="(value, key) in configList" :key="key">
          <span class="napa">{{ key }}</span>
          <span>{{ value || "--" }}</span>
          <el-button type="text" @click="openDialogDL(key, value)"
            >修改</el-button
          >
        </div>
      </div>
    </div>
    <!-- 新增编辑弹框 -->
    <el-dialog
      :title="'修改' + channel_title"
      width="500px"
      :visible.sync="dialog_visible"
      :close-on-click-modal="false"
      @click="closeDialog()"
    >
      <div class="content">
        <label><span style="color:red;">* </span>{{ channel_title }}</label>
        <el-input
          :placeholder="`请输入${channel_title}`"
          v-model.trim="channel_name"
          class="search-input"
          size="small"
        ></el-input>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialog_visible = false">取 消</el-button>
        <el-button type="primary" @click="save()">保 存</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/setting/enterprise";
@Component
export default class extends Vue {
  private dialog_visible = false;
  private channel_title = ""; // 标题
  private channel_name = ""; // 输入框的值
  private tokenList = []; //Access Token 管理
  private configList = []; // 品牌信息配置

  // 获取DL 配置的详情
  getDLConfigList() {
    this.channel_name = "";
    Api.getDLConfig({}).then((res: any) => {
      if (res.code == 200) {
        this.tokenList = res.data.access_token;
        this.configList = res.data.brand_config;
      } else {
        this.$message.error(res.message);
      }
    });
  }

  // 修改DL 配置
  private openDialogDL(key: string, value: string) {
    this.dialog_visible = true;
    this.channel_title = key;
    this.channel_name = value;
  }

  // 保存
  private save() {
    if (this.channel_name === "") {
      this.$message.error(`请输入${this.channel_title}`);
      return false;
    } else {
      const params = {
        key: this.channel_title,
        value: this.channel_name
      };
      Api.editDLConfig(params).then((res: any) => {
        if (res.code == 200) {
          this.$message.success("修改成功");
          this.dialog_visible = false;
          this.getDLConfigList();
        } else {
          this.$message.error(res.message);
        }
      });
    }
  }

  private mounted() {
    this.getDLConfigList();
  }
}
</script>
<style lang="scss" scoped>
.box {
  min-height: calc(100vh - 180px);
  background: white;
  margin-bottom: 10px;
  padding: 30px;
  border-radius: 10px;
  span {
    margin-right: 20px;
  }
  .editapp {
    margin-left: 20px;
  }
  h4 {
    color: rgba(16, 16, 16, 1);
    font-size: 18px;
    font-weight: bold;
    margin-bottom: 20px;
  }
  .desc {
    color: #909399;
    font-size: 14px;
    margin-bottom: 20px;
  }
  .desc-bot {
    color: #909399;
    font-size: 14px;
  }
  .line {
    width: 100%;
    height: 1px;
    margin: 30px 0;
    background: #ebeef5;
  }
}

.row {
  font-size: 14px;
  margin-bottom: 20px;
  .napa {
    color: #292a30;
  }
}
.module-left {
  text-align: right;
  margin-right: 20px;
  .row-finally {
    margin-bottom: 60px;
  }
}
.module-right {
  flex: 1;
  .listupload {
    margin-right: 20px;
  }
}
.content {
  display: flex;
  justify-content: space-around;
  align-items: center;
  lable {
    font-size: 13px;
  }
  ::v-deep .el-input--small {
    width: 70%;
  }
}
.el-button--text {
  margin-left: 10px;
}
</style>
