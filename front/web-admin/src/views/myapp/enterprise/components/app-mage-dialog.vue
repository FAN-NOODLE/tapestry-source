<template>
  <div class="counter">
    <layout-form>
      <template slot="content">
        <el-dialog
          :close-on-click-modal="false"
          :title="editappTitle"
          width="580px"
          :visible.sync="visibleDialog"
          @close="closeEditappDialog()"
        >
          <cy-form
            ref="cform"
            :rules="rules"
            :model="formLabelAlign"
            label-width="115px"
          >
            <cy-form-item label="Agentid" required number prop="agentid">
              <cy-input
                v-model="formLabelAlign.agentid"
                placeholder="请输入Agentid"
                @keyup.native="trimFn(1)"
                width-size="medium"
              ></cy-input>
            </cy-form-item>
            <cy-form-item label="Secret" required prop="secret">
              <cy-input
                v-model="formLabelAlign.secret"
                placeholder="请输入Secret"
                @keyup.native="trimFn(2)"
                width-size="medium"
              ></cy-input>
            </cy-form-item>
            <cy-form-item label="域名校验文件" required prop="check_file">
              <cy-upload
                :before-upload="beforeUpload"
                type="file"
                freedom
                :action="uploadUrl"
              >
                <span v-if="formLabelAlign.check_file" class="upfile">{{
                  formLabelAlign.check_file
                }}</span>
                <el-button type="text">{{
                  formLabelAlign.check_file ? "重新上传" : "上传"
                }}</el-button>
              </cy-upload>
            </cy-form-item>
          </cy-form>
          <span slot="footer" class="dialog-footer">
            <el-button @click="cancel()">取 消</el-button>
            <el-button type="primary" @click="save()">保 存</el-button>
          </span>
        </el-dialog>
      </template>
    </layout-form>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import Api from "@/api/enterprise/enterprise";
import { isNull } from "@/utils/tools";
import { VUE_APP_API_ROOT } from "@/utils/env-config"
@Component({})
export default class extends Vue {
  @Prop({ default: "", type: String }) editappTitle!: string;
  @Prop({ default: "", type: Number }) channelEditappAgeid!: number;
  @Prop({ default: "", type: String }) channelEditAppSecret!: string;
  @Prop({ default: "", type: String }) tocheckFile!: string;
  // 控制弹框的关闭隐藏
  private visibleDialog = true;
  // 表单验证，此处需要
  private rules = {};
  // 上传接口路径
  private uploadUrl = `${VUE_APP_API_ROOT}gw-scrm/upload-check-file`;
  // 表单默认数据
  private formLabelAlign: any = {
    agentid: Number,
    secret: "",
    check_file: ""
  };
  // 输入框限制输入空格
  private trimFn(type: number) {
    if (type == 1) {
      this.formLabelAlign.agentid = this.formLabelAlign.agentid.replace(
        /^\s+|\s+$/gm,
        ""
      );
    } else {
      this.formLabelAlign.secret = this.formLabelAlign.secret.replace(
        /^\s+|\s+$/gm,
        ""
      );
    }
  }
  // 上传文件限制规范
  private beforeUpload(file: any) {
    if (file.type != "text/plain") {
      this.$message.error("请上传TXT格式的文件");
      return false;
    }
    if (file.size > 1048576) {
      this.$message.error("请上传1M以内的文件");
      return false;
    }
    this.formLabelAlign.check_file = file.name;
  }

  // 取消，关闭弹框
  cancel() {
    this.$emit("close", false);
  }
  // 关闭
  closeEditappDialog() {
    this.$emit("close", false);
  }
  // 自建应用更新接口
  private workPluginUpdate() {
    let type: number;
    if (this.editappTitle == "修改自建应用") {
      type = 2;
    } else {
      type = 1;
    }
    const params = {
      agent_id: parseInt(this.formLabelAlign.agentid),
      secret: this.formLabelAlign.secret,
      file_url: this.formLabelAlign.check_file,
      check_file: 1,
      type: type
    };
    Api.workPluginUpdateApi(params).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("保存成功");
        this.$emit("close", false);
        this.$emit("formLabelAlign", this.formLabelAlign);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 保存
  private save() {
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        if (isNull(this.formLabelAlign.check_file)) {
          this.$message.error("请上传域名校验文件");
          return;
        }
        this.workPluginUpdate();
      } else {
        return;
      }
    });
  }
  private mounted() {
    this.formLabelAlign.agentid = this.channelEditappAgeid;
    this.formLabelAlign.secret = this.channelEditAppSecret;
    this.formLabelAlign.check_file = this.tocheckFile;
  }
}
</script>
<style lang="scss" scoped>
::v-deep {
  .in-item:last-child {
    margin: 0;
  }
  .el-dialog__body {
    padding-bottom: 15px !important;
  }
}
::v-deep.el-form-item__label {
  width: 100px;
  height: 10px;
}
::v-deep.el-dialog__body {
  padding-bottom: 0;
}
.dialog-upload {
  position: relative;
  cursor: pointer;
  input {
    position: absolute;
    width: 100%;
    height: 100%;
    left: 0;
    top: 0;
    opacity: 0;
    filter: alpha(opacity=0);
    cursor: pointer;
  }
}
.upfile {
  margin-right: 10px;
}
</style>
