<template>
  <div class="add-edit-group a">
    <el-dialog
      :title="sourceId ? '编辑文件' : '上传文件'"
      :visible.sync="show_dialog"
      width="700px"
      @close="closeDialog"
    >
      <cy-form
        class="main-box"
        ref="cform"
        :model="ruleForm"
        label-width="90px"
        :rules="rules"
        @submit.native.prevent
      >
        <!-- <cy-form-item label="素材分组" prop="group_id">
          <cy-select
            v-model="ruleForm.group_id"
            size="medium"
            placeholder="请选择素材分组"
            search-input
            :optionsList="groupList"
            :options="{ label: 'group_name', value: 'id' }"
          >
          </cy-select>
        </cy-form-item> -->
        <cy-form-item label="素材分组">
          <el-cascader
            v-model="ruleForm.group_id"
            :options="groupList"
            :props="{
              label: 'group_name',
              value: 'id',
              checkStrictly: true,
              emitPath: false
            }"
            :style="{
              width: '400px'
            }"
          ></el-cascader>
          
        </cy-form-item>
        <cy-form-item
          label="有效期"
          prop="datePicker"
          required
          trigger="change"
        >
          <el-date-picker
            v-model="ruleForm.datePicker"
            :default-time="['00:00:00', '23:59:59']"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            clearable
            value-format="timestamp"
            @change="datePickerChange"
          >
          </el-date-picker>
        </cy-form-item>
        <cy-form-item label="文件名" prop="name" required>
          <cy-input
            ref="inp"
            placeholder="请输入文件名"
            v-model.trim="ruleForm.name"
            size="medium"
            maxlength="50"
            show-word-limit
            width-size="medium"
            :disabled="!sourceId && !fileSuccessObj.url"
          ></cy-input>
        </cy-form-item>
        <cy-form-item label="文件" required>
          <div v-if="!sourceId">
            <large-file-upload
              :on-success="onSuccessFn"
              :beforeUpload="beforeUploadFn"
              :onRemove="removeFn"
              :onError="errorFn"
              accept="application/vnd.ms-excel,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,application/vnd.ms-powerpoint,application/vnd.openxmlformats-officedocument.presentationml.presentation,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document,application/pdf,text/plain,.csv,.doc,.docx"
            >
              <el-button type="primary" v-if="isShowBtn" size="small"
                >上传文件</el-button
              >
            </large-file-upload>
          </div>
          <div class="file-show" v-else>
            {{ ruleForm.name }}
          </div>
          <div class="tip-view">
            <p class="tip-text">
              1. 文件格式仅支持*.txt, *.pdf, *.ppt, *.pptx, *.doc, *.docx,
              *.xls, *.xlsx
            </p>
            <p class="tip-text">2. 文件大小不得超过20M</p>
          </div>
        </cy-form-item>
      </cy-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="closeDialog" size="small">取 消</el-button>
        <el-button type="primary" @click="submitFn" size="small"
          >确 认</el-button
        >
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import { Loading } from "element-ui";
import API from "@/api/enterprise/material-library";
import LargeFileUpload from "@/components/business/large-file-upload/index_o.vue";
import { chatDataApi } from "@/api/enterprise/association";
import Times from "./datePicker";
@Component({
  components: {
    LargeFileUpload
  }
})
export default class extends Vue {
  @Prop() groupId?: any;
  @Prop() groupName!: string;
  @Prop() sourceId?: any;

  //------------------------------------------------------data-----------------------------------------------------------/
  private show_dialog = true;
  private groupList: any = []; //下拉列表
  private ruleForm: any = {
    group_name: "",
    name: "",
    group_id: 0,
    datePicker: Times.times,
    start_time: Times.start, //有效期开始时间
    end_time: Times.end //有效期结束时间
  };
  private rules = {};
  private fileObj: any = {}; //上传文件信息
  private fileSuccessObj: any = {}; //上传成功信息
  private isShowBtn = true; //是否显示按钮
  private isClosePop = true; //是否关闭弹框
  private isFail = false; //是否上传失败
  //----------------------------------------------------methods-----------------------------------------------------------/
  // 获取素材详情
  private sourceDetail() {
    API.sourceDetail({
      source_id: this.sourceId,
      type: 6
    }).then((res: any) => {
      if (res.code == 200) {
        const data = res.data[0];
        if (data.start_time) {
          this.ruleForm.start_time = data.start_time;
          this.ruleForm.end_time = data.end_time;
          this.ruleForm.datePicker = [
            data.start_time * 1000,
            data.end_time * 1000
          ];
        }
        this.ruleForm.name = data.name;
        this.ruleForm.group_id = data.group_id || -1;
        this.fileSuccessObj.url = data.url;
        if (!data.group_id) {
          this.ruleForm.group_name = "未分组";
        } else if (data.group_id === -1) {
          this.ruleForm.group_name = "未分组";
        } else {
          this.ruleForm.group_name = data.group_name;
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  selectNode(node: any) {
    this.ruleForm.group_name = node.group_name;
    this.ruleForm.group_id = node.id;
    (this.$refs.selectBlur as any).blur();
  }
  private submitFn() {
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        if (this.sourceId) {
          this.sourceUpdate();
        } else {
          this.addConfirm();
        }
      }
    });
  }
  // 编辑接口
  private sourceUpdate() {
    API.sourceUpdate({
      source_type: 6,
      id: this.sourceId,
      group_id: this.ruleForm.group_id || 0,
      start_time: this.ruleForm.start_time,
      end_time: this.ruleForm.end_time,
      category_ids: [],
      name: this.ruleForm.name
    }).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("编辑成功");
        this.closeDialog();
        this.$emit("save");
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 点击完成
  private addConfirm() {
    if (!this.isClosePop) {
      this.$message.warning("文件正在上传中");
      return false;
    }
    if (this.isFail) {
      this.$message.error("文件上传失败");
      return false;
    }
    if (this.isShowBtn) {
      this.$message.warning("请上传文件");
      return false;
    }
    const suffix = this.fileObj.name.split(".")[
      this.fileObj.name.split(".").length - 1
    ];
    const params = {
      category_ids: [],
      group_id: this.ruleForm.group_id || 0, //分组id
      name: this.ruleForm.name, //文件名称
      url: this.fileSuccessObj.url, //文件存储url
      ext: suffix, //文件后缀
      size: this.fileObj.size, //文件大小
      start_time: this.ruleForm.start_time,
      end_time: this.ruleForm.end_time
    };

    API.fileAddApi(params).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("添加成功");
        this.closeDialog();
        if (suffix == "mp4") {
          this.$emit("isupdate", "video");
        } else if (suffix == "jpg" || suffix == "png") {
          this.$emit("isupdate", "img");
        } else {
          this.$emit("save");
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 上传
  private update(res: Array<string>) {
    this.$set(this.ruleForm, "urls", res);
    Loading.service({ fullscreen: true }).close();
  }
  // 关闭弹框
  private closeDialog() {
    this.$emit("closeDialog");
  }

  //素材分组列表
  private groupListFn() {
    API.sourceGroupListApi({}).then((res: any) => {
      if (res.code == 200) {
        const arr = [{ group_name: "未分组", id: -1 }];
        this.groupList = [...arr, ...res.data];
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //上传成功
  private onSuccessFn(res: any) {
    this.isShowBtn = false;
    this.isClosePop = true;
    this.isFail = false;
    this.ruleForm.name = res.data.fileName;
    this.fileSuccessObj = res.data;
  }
  private beforeUploadFn(file: any, type: string) {
    if (
      (file.type == "image/png" || file.type == "image/jpeg") &&
      file.size > 2 * 1024 * 1024
    ) {
      this.$message.error("请上传2MB大小以内的图片");
      return false;
    }

    if (file.type == "video/mp4" && file.size > 1024 * 1024 * 10) {
      this.$message.error("请上传10M大小以内的视频");
      return false;
    }
    if (
      file.type !== "video/mp4" &&
      file.type !== "image/png" &&
      file.type !== "image/jpeg" &&
      file.size > 1024 * 1024 * 20
    ) {
      this.$message.error("请上传20M大小以内的文件");
      return false;
    }
    // 允许上传的文件类型
    const type_arr = [
      "txt",
      "pdf",
      "ppt",
      "pptx",
      "doc",
      "docx",
      "xls",
      "xlsx"
    ];
    const nameArr = file.name.split("."); // 通过文件后缀获取文件类型
    const tyle_suffix = nameArr[nameArr.length - 1];
    if (type === "fileUplod" && !type_arr.includes(tyle_suffix)) {
      this.$message.error(
        "请上传txt, pdf, ppt, pptx, doc, docx, xls, xlsx格式文件"
      );
      return false;
    }

    this.isShowBtn = false;
    this.isClosePop = false;
    this.isFail = false; //上传失败
    this.fileObj = file;
    return true;
  }
  //移除
  private removeFn() {
    this.isShowBtn = true;
    this.isClosePop = true; //是否关闭弹框
    this.isFail = false; //是否上传失败
  }
  // 有效期选择
  private datePickerChange(val: any) {
    if (val) {
      // 因为 parseInt 的参数类型为 string ,所以处理后的时间戳只能先转成 string 然后再取整
      this.ruleForm.start_time = parseInt((val[0] / 1000).toString());
      this.ruleForm.end_time = parseInt((val[1] / 1000).toString());
    } else {
      this.ruleForm.start_time = 0;
      this.ruleForm.end_time = 0;
    }
  }
  private mounted() {
    //素材分组列表
    this.groupListFn();
    if (this.sourceId) {
      this.sourceDetail();
    } else {
      this.ruleForm.group_id = this.groupId || -1;
      this.ruleForm.group_name = this.groupId ? this.groupName : "未分组";
    }
  }
  private errorFn() {
    this.isClosePop = true;
    this.isFail = true; //上传失败
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-upload-list__item:first-child {
  margin-top: 0px;
}
/deep/.el-upload-list__item-name {
  width: 240px;
}
.tip-view {
  margin-top: 10px;
  .tip-text {
    font-size: 12px;
    color: $color-danger;
    line-height: 18px;
  }
}
.file-show {
  width: 80%;
  height: 36px;
  background: #000;
  cursor: default;
  color: #ffffff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 36px;
  padding: 0 15px;
  border-radius: 4px;
}
</style>
