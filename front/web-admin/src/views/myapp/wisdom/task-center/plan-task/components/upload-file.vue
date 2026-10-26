<template>
  <div class="counter">
    <div class="label">
      <el-dialog
        :close-on-click-modal="false"
        :title="is_file ? '导入员工人群包信息反馈' : '导入员工人群包'"
        width="800px"
        :visible.sync="visibleDialog"
        @close="closeEditappDialog()"
      >
        <div v-if="!is_file">
          <cy-form ref="cform" :rules="rules" label-width="90px">
            <cy-form-item style="max-width: none; width: 100%;">
              <div class="type-radio">
                <el-radio class="m-t-10" v-model="upType" label="1"
                  >根据UserID导入</el-radio
                >
                <div>
                  <el-radio v-model="upType" label="2"
                    >根据Marscode导入</el-radio
                  >
                  <div class="tips">
                    需在企微员工通讯录管理预先配置Marscode才能进行匹配
                  </div>
                </div>
              </div>
            </cy-form-item>
            <cy-form-item label="模板下载" required>
              <el-button v-if="upType == 1" type="text" @click="temDown(1)"
                >下载员工UserID数据模版</el-button
              >
              <el-button v-if="upType == 2" type="text" @click="temDown(2)"
                >下载员工Marscode数据模版</el-button
              >
            </cy-form-item>
            <cy-form-item label="模板上传" required>
              <el-upload
                ref="upload"
                type="file"
                :limit="1"
                :on-exceed="excedFun"
                :data="exceldata"
                :action="uploadUrl"
                :before-upload="beforeUpload"
                :headers="headers"
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
            <cy-form-item label="上传结果">
              <div>
                {{
                  upSuccessData.length + errorData.length == 0
                    ? "暂未检测到上传信息!"
                    : "上传成功!"
                }}
              </div>
            </cy-form-item>
          </cy-form>
        </div>
        <div v-else>
          <div>
            成功{{ upSuccessData.length }}条，重复{{ repeat }}条，失败{{
              errorData.length
            }}条
          </div>
          <div class="m-t-20">导入信息错误列表</div>
          <div class="m-t-20">
            <el-table
              stripe
              :data="errorData"
              max-height="400"
              class="mar-t-20"
            >
              <el-table-column
                min-width="100"
                label="行数"
                prop="line_num"
              ></el-table-column>
              <el-table-column
                min-width="100"
                prop="item"
                :label="upType == 1 ? 'User ID' : 'Marscode'"
              >
              </el-table-column>
              <el-table-column min-width="150" prop="error" label="错误信息">
                <template slot-scope="{ row }">
                  <div style="color: #F8586C;">{{ row.error }}</div>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
        <span slot="footer" class="dialog-footer">
          <el-button size="small" @click="cancel()">取 消</el-button>
          <el-button size="small" type="primary" @click="save()"
            >确 定</el-button
          >
        </span>
      </el-dialog>
    </div>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import LoginAboat from "@/utils/LoginManager";
import { VUE_APP_API_VERSION } from "@/utils/env-config";
@Component({
  components: {},
})
export default class extends Vue {
  private headers = {
    "Enterprise-Hash": LoginAboat.getRealHash(),
    "Auth-Token": LoginAboat.getToken(),
  };
  private exceldata: any = {
    file_type: "userid",
  };
  private upSuccessData = []; // 上传成功数据
  private errorData = []; // 错误的数据
  private repeat = 0; // 重复的数据
  private is_file = false; // 上传成功
  private upType = "1";
  // 控制弹框的关闭隐藏
  private visibleDialog = true;
  private updateSelectedTags: Array<any> = [];
  private excelFile = "";
  // 表单验证，此处需要
  private rules = {};
  private url = "";
  // 上传接口路径
  // private uploadUrl = `http://localhost:20001/api/micro-scrm/v1/plan/subtask/upload-users`;
  private uploadUrl = `${VUE_APP_API_VERSION}plan/subtask/upload-users`;
  private selectedTags: Array<any> = [];
  private showTagsList: Array<string> = [];
  private tagDialog = false;
  // 上传文件个数超过范围
  private excedFun() {
    this.$message.error("上传失败，每次仅支持上传一个文件！");
  }
  // 上传文件限制规范
  private beforeUpload(file: any) {
    this.exceldata.file_type = this.upType == "1" ? "userid" : "mars_code";
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
      res.data = res.data || [];
      if (res.data.length == 0) {
        this.$message.warning("检测上传文件内容为空，请重新上传");
        (this.$refs.upload as any).clearFiles();
        return;
      }
      this.upSuccessData = res.data.filter((v: any) => !v.error);
      this.errorData = res.data.filter((v: any) => v.error);
      this.repeat = res.data.filter((v: any) => v.err_type == 2).length;
    } else {
      this.$message.error(res.message);
    }
  }
  // 移除上传
  private onRemoveFn() {
    this.upSuccessData = [];
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
    if (!this.is_file) {
      if (this.upSuccessData.length + this.errorData.length == 0) {
        this.$message.warning("请上传文件");
        return;
      }
      this.is_file = true;
    } else {
      this.upSuccessData.forEach((v: any) => {
        v.name = v.user_name;
      });
      this.$emit("submit", this.upSuccessData);
    }
  }
  // 模板下载
  private temDown(type: number) {
    import("@/utils/Export2Excel").then((excel) => {
      const tHeader = [type == 1 ? "员工Userid" : "员工Marscode"];
      const data: any = [];
      excel.export_json_to_excel({
        header: tHeader,
        data,
        filename: type == 1 ? "UserID数据模版" : "Marscode数据模版",
        autoWidth: true,
        bookType: "xlsx",
      });
    });
  }
  private mounted() {
    console.log("mounted");
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
  line-height: 0;
}
.type-radio {
  display: flex;
}
</style>
<style lang="scss">
.counter {
  .label {
    .el-form-item__label {
      width: 90px !important;
    }
    .el-form-item__content {
      margin-left: 90px !important;
    }
  }
}
</style>
