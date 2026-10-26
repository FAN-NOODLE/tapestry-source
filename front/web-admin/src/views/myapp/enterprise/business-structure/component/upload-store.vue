<template>
  <div class="counter">
    <el-dialog
      :close-on-click-modal="false"
      :title="'上传门店信息'"
      :width="'600px'"
      :visible.sync="visibleDialog"
      @close="closeEditappDialog()"
    >
      <div>
        <cy-form ref="cform" :rules="rules" label-width="90px">
          <cy-form-item label="导入模板">
            <el-button type="text" @click="temDown" class="btn-tex"
              >下载导入模板</el-button
            >
          </cy-form-item>
          <cy-form-item label="导入文档" required>
            <div class="upload-box">
              <large-file-upload
                v-model="fileUrl"
                :beforeUpload="beforeUpload"
                :onSuccess="onSuccessFun"
                :onRemove="onRemoveFn"
                :onError="onErrorFun"
                accept="application/vnd.ms-excel,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,application/vnd.ms-powerpoint,application/vnd.openxmlformats-officedocument.presentationml.presentation,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document,application/pdf,text/plain,.csv,.doc,.docx"
              >
                <el-button type="primary" size="small">上传文件</el-button>
              </large-file-upload>
            </div>
            <!-- <div class="el-upload__text">
              <div>只能上传xls、xlsx文件</div>
            </div> -->
          </cy-form-item>
        </cy-form>
      </div>

      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="closeEditappDialog">取 消</el-button>
        <el-button size="small" type="primary" @click="submit()"
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
import Api from "@/api/enterprise/business";
import LargeFileUpload from "@/components/business/large-file-upload/index_o.vue";
@Component({
  components: {
    LargeFileUpload
  }
})
export default class extends Vue {
  private exceldata: any = {
    file_type: "file"
  };
  private headers = {
    "Enterprise-Hash": LoginAboat.getRealHash(),
    "Auth-Token": LoginAboat.getToken()
  };
  private enterprise_id = localStorage.getItem("ENTERPRISE_ID");
  // 控制弹框的关闭隐藏
  private visibleDialog = true;
  private excelFile = "";
  // 表单验证，此处需要
  private rules = {};
  private fileUrl = "";
  // 上传文件个数超过范围
  // private excedFun() {
  //   this.$message.error("上传失败，每次仅支持上传一个文件！");
  // }
  // 上传文件限制规范
  private beforeUpload(file: any) {
    console.log("TTType: ", file.type);

    // const testmsg = file.name.substring(file.name.lastIndexOf(".") + 1);
    // const extension = testmsg === "xlsx" || testmsg === "xls";
    // if (!extension) {
    //   this.$message.error("请上传 xls、xlsx 类型的文件!");
    //   return false;
    // }
    return true;
  }
  // 上传错误
  private onErrorFun(res: any) {
    console.log(" this.isFail = true; //上传失败");
    // if (res.code !== 200) {
    //   this.$message.error(res.message);
    // }
  }
  //上传成功
  private onSuccessFun(res: any) {
    const { url } = res.data;
    this.fileUrl = url;
  }
  //移除上传的文件
  private onRemoveFn() {
    this.fileUrl = "";
  }
  // 关闭
  closeEditappDialog() {
    this.$emit("close");
  }

  private submit() {
    Api.storeImportInfo({
      file_url: this.fileUrl
    }).then((res: any) => {
      console.log("RRREESSS: ", res);
      if (res.code === 200) {
        this.closeEditappDialog();
        // 刷新列表数据
        this.$emit("updateStoreList");
        this.$message.success("门店上传成功");
        // } else {
        //   this.$message.error(res.message);
      }
      if (res.code === 1002) {
        this.$message.error(res.message);
      }
    });
  }
  // 模板下载
  private temDown() {
    const params = {
      enterprise_id: this.enterprise_id
    };
    Api.downloadStoreList(params).then((res: any) => {
      // if (res.code !== 200) {
      //   this.$message.error("导出文件失败");
      //   return;
      // }
      const blob = new Blob([res], {
        type: "application/x-xls;charset=utf-8"
      });
      let fileName = "";
      const nowTime = new Date();
      fileName = "门店信息更新模版_" + nowTime + ".xlsx";

      //下载的文件名称及其后缀，后缀要和后台保持的一致

      if ("download" in document.createElement("a")) {
        // 非IE下载
        const elink = document.createElement("a");
        elink.download = fileName;
        elink.style.display = "none";
        elink.href = URL.createObjectURL(blob);
        document.body.appendChild(elink);
        elink.click();
        URL.revokeObjectURL(elink.href); // 释放URL 对象
        document.body.removeChild(elink);
      } else {
        // IE10+下载
        (navigator as any).msSaveBlob(blob, fileName);
      }
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
  .upload-box,
  .upload-box > div,
  .upload-box .el-upload {
    // width: 100px !important;
    height: 50px !important;
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
.btn-tex {
  color: #708df5 !important;
}
</style>
