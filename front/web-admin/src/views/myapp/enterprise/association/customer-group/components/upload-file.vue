<template>
  <div class="counter">
    <el-dialog
      :close-on-click-modal="false"
      :title="upSuccess ? '导入社群设置信息反馈' : '批量社群设置'"
      :width="upSuccess ? '800px' : '600px'"
      :visible.sync="visibleDialog"
      @close="closeEditappDialog()"
    >
      <div v-if="!upSuccess">
        <cy-form ref="cform" :rules="rules" label-width="90px">
          <cy-form-item label="导入模板">
            <el-button type="text" @click="temDown">下载导入模板</el-button>
          </cy-form-item>
          <cy-form-item label="导入文档" required>
            <el-upload
              ref="upload"
              type="file"
              :limit="1"
              :on-exceed="excedFun"
              :data="exceldata"
              :action="uploadUrl"
              :headers="headers"
              :before-upload="beforeUpload"
              :on-success="onSuccessFun"
              :on-remove="onRemoveFn"
              drag
            >
              <i class="el-icon-upload"></i>
              <div class="el-upload__text">
                将文件拖到此处，或<em>点击上传</em>
                <div>只能上传xls、xlsx文件</div>
              </div>
            </el-upload>
          </cy-form-item>
        </cy-form>
      </div>
      <div v-if="upSuccess">
        <div>
          成功{{ upSuccessData }}条，重复{{ repeated }}条，失败<span
            style="color: #F8586C;"
            >{{ errorData.length }}</span
          >条
        </div>
        <div>
          <el-button type="text"
            >注：成功上传的社群设置数据将覆盖列表中原有社群设置数据</el-button
          >
        </div>
        <div class="m-t-20">导入信息错误列表</div>
        <div class="m-t-20">
          <el-table stripe :data="errorData" max-height="400" class="mar-t-20">
            <el-table-column
              min-width="100"
              label="行数"
              prop="line_num"
            ></el-table-column>
            <el-table-column min-width="120" prop="chat_id" label="群ID">
            </el-table-column>
            <el-table-column min-width="120" prop="name" label="群名称">
            </el-table-column>
            <el-table-column
              min-width="120"
              prop="channel_name"
              label="所属渠道"
            >
            </el-table-column>
            <el-table-column min-width="120" prop="chat_type" label="群类型">
            </el-table-column>
            <el-table-column min-width="120" prop="operator" label="群运营者">
            </el-table-column>
            <el-table-column
              min-width="150"
              prop="operator"
              label="群运营者手机号"
            >
              <template slot-scope="{ row }">
                {{
                  row.mobile
                    ? row.mobile.slice(0, 3) + "****" + row.mobile.slice(7, 11)
                    : ""
                }}
              </template>
            </el-table-column>
            <el-table-column min-width="120" prop="chat_target" label="群目标">
            </el-table-column>
            <el-table-column min-width="150" prop="error" label="错误信息">
              <template slot-scope="{ row }">
                <div style="color: #F8586C;">
                  {{ row.errors ? row.errors.join(",") : "" }}
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="cancel()">取 消</el-button>
        <el-button v-if="!upSuccess" size="small" type="primary" @click="save()"
          >确认导入</el-button
        >
        <el-button v-else size="small" type="primary" @click="submit()"
          >确 定</el-button
        >
      </span>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import LoginAboat from "@/utils/LoginManager";
import { batchSetting } from "@/api/enterprise/association";
import { VUE_APP_API_ROOT_RE } from "@/utils/env-config";
@Component({
  components: {}
})
export default class extends Vue {
  private setting_key = "";
  private upSuccess = false;
  private exceldata: any = {
    file_type: "work_group_setting"
  };
  private headers = {
    "Enterprise-Hash": LoginAboat.getRealHash(),
    "Auth-Token": LoginAboat.getToken()
  };
  private errorData = [];
  private upSuccessData = 0;
  private repeated = 0;
  private downloadUrl = `${VUE_APP_API_ROOT_RE}group-chat/download?file_type=chat-setting`;
  // 控制弹框的关闭隐藏
  private visibleDialog = true;
  private updateSelectedTags: Array<any> = [];
  private excelFile = "";
  // 表单验证，此处需要
  private rules = {};
  private url = "";
  // 上传接口路径
  // private uploadUrl = `${this.$VUE_APP_API_ROOT}filefile/upload`;
  private uploadUrl = `${this.$VUE_APP_API_ROOT}chat-group/upload-settings`;
  private selectedTags: Array<any> = [];
  private showTagsList: Array<string> = [];
  private tagDialog = false;
  private uploading = false;
  // 上传文件个数超过范围
  private excedFun() {
    this.$message.error("上传失败，每次仅支持上传一个文件！");
  }
  // 上传文件限制规范
  private beforeUpload(file: any) {
    this.uploading = true;
    const testmsg = file.name.substring(file.name.lastIndexOf(".") + 1);
    const extension = testmsg === "xlsx" || testmsg === "xls";
    if (!extension) {
      this.$message.error("请上传 xls、xlsx 类型的文件!");
      return false;
    }
  }
  // 上传成功
  private onSuccessFun(res: any) {
    if (res.code == 200) {
      if (res.data) {
        if (res.data.success + res.data.errors + res.data.repeated == 0) {
          this.$message.warning("检测上传文件内容为空，请重新上传");
          (this.$refs.upload as any).clearFiles();
          return;
        }
        this.upSuccessData = res.data.success;
        this.errorData = res.data.error_items || [];
        this.repeated = res.data.repeated;
        this.setting_key = res.data.setting_key;
      }
    } else {
      this.$message.error(res.message);
    }
  }
  // 移除上传
  private onRemoveFn() {
    this.upSuccessData = 0;
    this.repeated = 0;
    this.setting_key = "";
    this.errorData = [];
  }
  // 取消，关闭弹框
  cancel() {
    this.$emit("close");
  }
  // 关闭
  closeEditappDialog() {
    this.$emit("close");
  }
  // 保存
  private save() {
    if (this.upSuccessData + this.repeated + this.errorData.length > 0) {
      this.upSuccess = true;
    } else {
      this.$message.warning("请上传文件");
    }
  }
  private submit() {
    batchSetting({
      setting_key: this.setting_key
    }).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("上传成功");
        this.$emit("submit");
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 模板下载
  private temDown() {
    // const a = document.createElement("a");
    // const evt = document.createEvent("MouseEvents");
    // a.download = "导入模板.xlsx";
    // a.href = "/saas/import_chat_setting.xlsx";
    // evt.initEvent("click", true, true);
    // a.dispatchEvent(evt);
    // window.URL.revokeObjectURL(a.href);
    window.location.href = this.downloadUrl;
  }
  private mounted() {
    console.log("mounted");
  }
  public get tagsToTableData(): any {
    return [
      {
        list: this.showTagsList
      }
    ];
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
::v-deep {
  .el-upload {
    width: 300px;
    height: 180px;
  }
  .el-upload-dragger {
    @extend .el-upload;
    border: 1px solid $color-second-line !important;
    border-radius: 4px !important;
  }
  .el-upload__tip {
    color: $color-text-second !important;
  }
  .el-upload-dragger {
    .el-icon-upload {
      font-size: 80px !important;
      margin-top: 30px !important;
      margin-bottom: 20px !important;
    }
  }

  .el-dialog {
    max-height: none !important;
  }
}
.upfile {
  margin-right: 10px;
}
.tips {
  font-size: 13px;
  color: #909399;
  margin-bottom: 10px;
}

.el-tag {
  background: #f0f3fe;
  color: #708df5;
  border: 1px solid #d4dcfc;
}
.tableTag {
  margin: 0 10px 10px 0;
}
</style>
