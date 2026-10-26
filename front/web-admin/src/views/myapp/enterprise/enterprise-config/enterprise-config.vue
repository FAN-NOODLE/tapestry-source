<template>
  <div class="config">
    <div class="box">
      <!-- 密钥管理 -->
      <div class="module">
        <h4>密钥管理</h4>
        <div class="desc">
          请前往【企业微信后台-客户联系-客户】中点击【api】，复制Secret
        </div>
        <div class="row">
          <span class="napa">外部联系人Secret</span
          ><span>{{ tempData.external_secret }}</span>
          <el-button
            type="text"
            @click="openDialog(1)"
            v-if="!tempData.external_secret"
            >添加</el-button
          >
          <el-button
            v-else
            type="text"
            @click="openDialog(1, tempData.external_secret)"
            >修改</el-button
          >
        </div>
      </div>
      <div class="line"></div>
      <!-- 客户联系配置 -->
      <div class="module">
        <h4>客户联系配置</h4>
        <div class="desc">
          复制至【企业微信后台-客户联系-客户】中点击【api-接收事件服务器】
        </div>
        <div class="module-box">
          <div class="module-left">
            <div class="row">接收推送请求URL</div>
            <div class="row">接收推送请求token</div>
            <div>接收推送请求key</div>
          </div>
          <div class="module-right">
            <div class="row" v-if="tempData.uc_callback_url">
              <span>{{ tempData.uc_callback_url }}</span>
              <el-button type="text" @click="copy(tempData.uc_callback_url)"
                >复制</el-button
              >
            </div>
            <div class="row" v-else>-</div>
            <div class="row" v-if="tempData.token">
              <span>{{ tempData.token }}</span>
              <el-button type="text" @click="copy(tempData.token)"
                >复制</el-button
              >
            </div>
            <div class="row" v-else>-</div>
            <div v-if="tempData.key">
              <span>{{ tempData.key }}</span>
              <el-button type="text" @click="copy(tempData.key)"
                >复制</el-button
              >
            </div>
            <div class="row" v-else>-</div>
          </div>
        </div>
      </div>
    </div>
    <div class="box">
      <!-- 自建应用管理 -->
      <div class="module">
        <h4>
          自建应用管理<el-button
            type="text"
            @click="editapp(tempData.agent_id)"
            class="editapp"
            >{{ tempData.agent_id ? "修改" : "配置" }}自建应用</el-button
          >
        </h4>
        <div class="desc">
          请前往【企业微信后台-应用管理-应用-自建应用】，复制AgentId和Secret
        </div>
        <div class="module-box">
          <div class="module-left">
            <div class="row">
              Agentid
            </div>
            <div>
              Secret
            </div>
          </div>
          <div class="module-right">
            <div class="row">
              <span v-if="tempData.agent_id > 0">{{ tempData.agent_id }}</span>
              <el-button
                type="text"
                @click="editapp()"
                v-if="!tempData.agent_id"
                >添加</el-button
              >
              <el-button v-else type="text" @click="editapp(tempData.agent_id)"
                >修改</el-button
              >
            </div>
            <div>
              <span>{{ tempData.secret }}</span>
              <el-button type="text" @click="editapp()" v-if="!tempData.secret"
                >添加</el-button
              >
              <el-button v-else type="text" @click="editapp(tempData.secret)"
                >修改</el-button
              >
            </div>
          </div>
        </div>
      </div>
      <div class="line"></div>
      <!-- 网页授权及JS-SDK -->
      <div class="module">
        <h4>网页授权及JS-SDK</h4>
        <div class="desc">
          请前往【企业微信后台-应用管理-应用-自建应用-网页授权及JS-SDK】，“粘贴可信域名”与“下载校验文件”
        </div>
        <div class="module-box">
          <div class="module-left">
            <div class="row">
              可信域名
            </div>
            <div>
              域名校验文件
            </div>
          </div>
          <div class="module-right">
            <div class="row" v-if="tempData.domain">
              <span>{{ tempData.domain }}</span>
              <el-button type="text" @click="copy(tempData.domain)"
                >复制</el-button
              >
            </div>
            <div class="row" v-else>-</div>
            <div>
              <span v-if="tempData.file_url" class="listupload">{{
                tempData.file_url
              }}</span>
              <el-button @click="editapp" type="text" class="listupload" v-else
                >上传</el-button
              >
              <el-button @click="editapp" type="text" v-if="tempData.file_url"
                >重新上传</el-button
              >
            </div>
          </div>
        </div>
      </div>
      <div class="line"></div>
      <!-- 接受消息 -->
      <!-- <div class="module">
        <h4>接受消息</h4>
        <div class="desc">
          客户联系配置：复制至【企业微信后台-应用管理-应用-自建应用-接受消息-设置API接收】，粘贴
        </div>
        <div class="module-box">
          <div class="module-left">
            <div class="row">URL</div>
            <div class="row">Token</div>
            <div>EncodingAESKey</div>
          </div>
          <div class="module-right">
            <div class="row">
              <span>{{ tempData.domain }}</span>
              <el-button type="text" @click="copy(tempData.domain)"
                >复制</el-button
              >
            </div>
            <div class="row">
              <span>{{ tempData.domain }}</span>
              <el-button type="text" @click="copy(tempData.domain)"
                >复制</el-button
              >
            </div>
            <div>
              <span>{{ tempData.domain }}</span>
              <el-button type="text" @click="copy(tempData.domain)"
                >复制</el-button
              >
            </div>
          </div>
        </div>
      </div>
      <div class="line"></div> -->
      <!-- 企微右侧插件参数 -->
      <div class="module">
        <h4>企微右侧插件参数</h4>
        <div class="desc">
          复制至【企业微信后台-应用管理-应用-自建应用-配置到聊天工具栏】中配置页面和粘贴URL
        </div>
        <div class="module-box">
          <div class="module-left">
            <div class="row">营销推广</div>
            <div class="row">素材库</div>
            <div class="row">客户画像</div>
            <div>裂变活动</div>
          </div>
          <div class="module-right">
            <div class="row" v-if="tempData.marketing_promotion_url">
              <span>{{ tempData.marketing_promotion_url }}</span>
              <el-button
                type="text"
                @click="copy(tempData.marketing_promotion_url)"
                >复制</el-button
              >
            </div>
            <div class="row" v-else>-</div>
            <div class="row" v-if="tempData.material_library_url">
              <span>{{ tempData.material_library_url }}</span>
              <el-button
                type="text"
                @click="copy(tempData.material_library_url)"
                >复制</el-button
              >
            </div>
            <div class="row" v-else>-</div>
            <div class="row" v-if="tempData.customer_information_url">
              <span>{{ tempData.customer_information_url }}</span>
              <el-button
                type="text"
                @click="copy(tempData.customer_information_url)"
                >复制</el-button
              >
            </div>
            <div class="row" v-else>-</div>
            <!-- 裂变活动 -->
            <div v-if="tempData.fission_activity_url">
              <span>{{ tempData.fission_activity_url }}</span>
              <el-button
                type="text"
                @click="copy(tempData.fission_activity_url)"
                >复制</el-button
              >
            </div>
            <div class="row" v-else>-</div>
          </div>
        </div>
      </div>
    </div>
    <div class="box">
      <!-- 会话存档 -->
      <div class="module">
        <h4>会话存档</h4>
        <div>
          <div class="module-item">
            <div class="module-left w-100">会话存档secret</div>
            <div class="module-right">
              <div class="row m-b-10">
                <span v-if="configData.secret">{{ configData.secret }}</span>
                <el-button
                  v-if="!configData.secret"
                  type="text"
                  @click="openDialog(2)"
                  >添加</el-button
                >
                <el-button
                  v-else
                  type="text"
                  @click="openDialog(2, configData.secret)"
                  >修改</el-button
                >
              </div>
              <div class="desc">
                请前往【企业微信后台--管理工具--会话内容存档】中申请开通，开通后将secret复制至此处
              </div>
            </div>
          </div>
          <div class="module-item">
            <div class="module-left w-100">可信IP地址</div>
            <div class="module-right">
              <div class="row">
                <div v-if="configData.trusted_ips">
                  <div
                    v-for="(item, index) in configData.trusted_ips"
                    :key="index"
                  >
                    <span class="ips-item">{{ item }}</span>
                    <el-button type="text" @click="copy(item)">复制</el-button>
                  </div>
                </div>
                <div class="row" v-else>-</div>
              </div>
              <div class="desc">
                复制至【企业微信后台--管理工具--会话内容存档】中可信IP地址
              </div>
            </div>
          </div>
          <div class="module-item">
            <div class="module-left w-100">消息加密公钥</div>
            <div class="module-right">
              <div class="row flex" v-if="configData.message_public_key">
                <el-input
                  type="textarea"
                  readonly
                  resize="none"
                  v-model="configData.message_public_key"
                >
                </el-input>
                <el-button
                  type="text"
                  @click="copy(configData.message_public_key)"
                  >复制</el-button
                >
              </div>
              <div class="row" v-else>-</div>
              <div class="desc-bot">
                复制至【企业微信后台-管理工具-会话存档】中消息加密公钥
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- 外部联系人Secret编辑弹框 -->
    <el-dialog
      :title="addOrEdit + channel_title"
      width="500px"
      :visible.sync="dialog_visible"
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
    <!-- 配置自建应用弹框 -->
    <app-mage-dialog
      v-if="dialog_editapp"
      :editappTitle="editappTitle"
      :channelEditappAgeid="channelEditappAgeid"
      :channelEditAppSecret="channelEditAppSecret"
      :tocheckFile="tocheckFile"
      @close="close"
      @formLabelAlign="formLabelAlign"
    />
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/enterprise/enterprise";
import appMageDialog from "../components/app-mage-dialog.vue";
@Component({
  components: {
    appMageDialog
  }
})
export default class extends Vue {
  // 自建应用弹框
  private dialog_editapp = false;
  private channel_name = "";
  private flag: any;
  // 传给组件dialog的数据
  private channelEditAppSecret = "";
  private channelEditappAgeid = 0;
  private tocheckFile = "";
  //外部联系人Secret，会话编辑弹框
  private dialog_visible = false;
  private addOrEdit = "";
  private channel_title = "";
  private editappTitle = "";
  private configData = {
    message_public_key: "",
    trusted_ips: "",
    secret: ""
  };
  // 数据
  private tempData: any = {};

  // 企微配置
  private settings() {
    Api.settingsApi().then((res: any) => {
      if (res.code == 200) {
        this.tempData = res.data;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 会话存档
  private getChatHistory() {
    Api.onGetChatSetting().then((res: any) => {
      if (res.code == 200) {
        if (Object.keys(res.data).length > 0) {
          this.configData = {
            secret: res.data.secret,
            message_public_key: res.data.message_public_key,
            trusted_ips: res.data.trusted_ips
          };
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 来自子组件传值
  formLabelAlign() {
    this.settings();
  }
  // 打开自建应用弹框
  private editapp(item: any) {
    if (item) {
      this.editappTitle = "修改自建应用";
      this.channelEditappAgeid = this.tempData.agent_id;
      this.channelEditAppSecret = this.tempData.secret;
      this.tocheckFile = this.tempData.file_url;
    } else {
      this.editappTitle = "配置自建应用";
      this.channelEditappAgeid = 0;
      this.channelEditAppSecret = "";
      this.tocheckFile = "";
    }
    console.log(this.dialog_editapp);
    this.dialog_editapp = true;
  }
  close() {
    this.dialog_editapp = false;
  }
  // 外部联系人secret,会话更新
  private authConacUpdate(id: any) {
    Api.authContacUpdateApi({ external_secret: id }).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("保存成功");
        this.settings();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 会话存档自建应用更新
  private updateChatSetting(secret: string) {
    const params = {
      agent_id: 0,
      type: this.configData.secret == "" ? 1 : 2,
      secret: secret
    };
    Api.onUpdateChatSetting(params).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("保存成功");
        this.getChatHistory();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 复制
  copy(data: string) {
    this.flag ? this.flag.close() : "";
    const url = data;
    const oInput: HTMLTextAreaElement = document.createElement("textarea");
    oInput.value = url;
    document.body.appendChild(oInput);
    oInput.select(); // 选择对象;
    document.execCommand("Copy"); // 执行浏览器复制命令
    this.flag = this.$message({
      message: "复制成功",
      type: "success"
    });
    oInput.remove();
  }
  // 打开秘钥管理，会话存档弹框
  private openDialog(type: number, item: any) {
    if (item) {
      this.addOrEdit = "修改";
      this.channel_name = item;
    } else {
      this.addOrEdit = "添加";
      this.channel_name = "";
    }
    this.dialog_visible = true;
    if (type === 1) {
      this.channel_title = "外部联系人Secret";
    } else {
      this.channel_title = "会话存档secret";
    }
  }
  // 秘钥管理，会话存档弹框点击保存
  private save() {
    if (this.channel_name === "") {
      this.$message.error(`请输入${this.channel_title}`);
      return false;
    } else {
      if (this.channel_title == "会话存档secret") {
        this.updateChatSetting(this.channel_name);
      } else {
        this.authConacUpdate(this.channel_name);
      }
    }
    this.dialog_visible = false;
    this.channel_name = "";
  }
  // 关闭弹框
  private closeDialog(close: boolean) {
    this.dialog_visible = close;
    this.settings();
    this.getChatHistory();
  }
  private mounted() {
    this.settings();
    this.getChatHistory();
  }
}
</script>
<style lang="scss" scoped>
::v-deep .el-dialog__body {
  padding-left: 0;
}
.editapp-box {
  display: flex;
  align-items: center;
  ::v-deep .el-input--small {
    width: 70%;
    margin-left: 20px;
  }
}
/deep/ .el-textarea {
  width: 475px;
}
/deep/ .el-textarea__inner {
  // width: 650px;
  height: 210px;
  border: none;
  margin-right: 0;
  padding: 0;
}
.flex {
  display: flex;
}
.module-item {
  display: flex;
}
.ips-item {
  // width: 200px;
  display: inline-block;
  margin-bottom: 15px;
}
.w-100 {
  width: 100px;
  text-align: right;
  margin-right: 20px;
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
.box {
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

.module-box {
  display: flex;
  // align-items: center;
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
}
.file {
  margin-top: 20px;
  display: flex;
  label {
    margin-right: 20px;
  }
  .relative-upload {
    position: relative;
    input {
      position: absolute;
      right: 0;
      top: 0;
      opacity: 0;
      filter: alpha(opacity=0);
      cursor: pointer;
    }
  }
}
</style>
