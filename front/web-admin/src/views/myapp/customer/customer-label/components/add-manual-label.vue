<template>
  <div class="add-manual-label">
    <el-dialog
      :title="currentEditItem.tag_id ? '编辑标签' : '新增标签'"
      :visible.sync="show_dialog"
      width="650px"
      @close="closeDialog"
    >
      <cy-form
        v-if="!currentEditItem.tag_id"
        class="main-box"
        ref="cform"
        :rules="rules"
        :model="ruleForm"
        label-width="120px"
      >
        <!-- todo:标签组名称 -->
        <cy-form-item
          label="标签组名称"
          prop="group_name"
          required
        >
          <cy-input
            style="width: 70% !important"
            v-model.trim="ruleForm.group_name"
            clearable
            width-size="medium"
            maxlength="30"
            placeholder="请输入标签组名称"
            show-word-limit
            @blur="changeInput"
          >
          </cy-input>
          <div class="choose-tag">
            <el-popover
              placement="bottom"
              trigger="click"
            >
              <el-select
                v-model="ruleForm.group_id"
                placeholder="请选择标签组"
                filterable
                clearable
                @change="changeGroup"
              >
                <el-option
                  :label="item.tag_group_name"
                  :value="item.tag_group_id"
                  :key="item.tag_group_id"
                  v-for="item in groupList"
                ></el-option>
              </el-select>
              <el-button
                type="text"
                slot="reference"
              >选择已有标签组</el-button>
            </el-popover>

          </div>
        </cy-form-item>
        <cy-form-item
          label="标签名称"
          prop="tag_name"
          required
        >
          <cy-input
            v-model.trim="ruleForm.tag_name"
            clearable
            width-size="medium"
            maxlength="30"
            placeholder="请输入标签名称"
            show-word-limit
          ></cy-input>
        </cy-form-item>
        <!-- todo:标签有效期长期有效,自定义删除时间 -->
        <cy-form-item label="标签有效期">
          <el-radio-group
            v-model="ruleForm.tag_validity"
            @change="changeTimeGroup"
          >
            <el-radio :label="1">长期有效</el-radio>
            <el-radio :label="2">自定义删除时间</el-radio>
          </el-radio-group>
        </cy-form-item>
        <cy-form-item
          prop="expired_time"
          v-if="ruleForm.tag_validity == 2"
        >
          <el-date-picker
            v-model="ruleForm.expired_time"
            type="datetime"
            placeholder="选择日期"
            width-size="medium"
            style="width: 100%"
            :picker-options="{
              disabledDate: (time) => {
                return time.getTime() < Date.now() - 8.64e7;
              }
            }"
          ></el-date-picker>
        </cy-form-item>
        <!-- todo: 上传标签人群包 -->
        <!-- dimension 3.4的时候需要下载模板和导入 -->
        <template>
          <cy-form-item label="下载模板">
            <el-button
              type="text"
              @click="temDown()"
            >下载标签人群包模板</el-button>
          </cy-form-item>
          <cy-form-item>
            <large-file-upload
              :on-success="onSuccessFun"
              :beforeUpload="beforeUpload"
              :onRemove="onRemoveFn"
              accept=".xlsx,.xls,.csv"
            >
              <el-button
                type="primary"
                size="small"
              >上传文件</el-button>
            </large-file-upload>
          </cy-form-item>
        </template>
      </cy-form>
      <cy-form
        v-if="currentEditItem.tag_id"
        class="main-box"
        ref="cform"
        :rules="rules"
        :model="ruleForm"
        label-width="120px"
      >
        <cy-form-item
          label="标签名称"
          prop="tag_name"
          required
        >
          <cy-input
            v-model.trim="ruleForm.tag_name"
            clearable
            width-size="medium"
            maxlength="30"
            placeholder="请输入标签名称"
            show-word-limit
          ></cy-input>
        </cy-form-item>
      </cy-form>
      <span
        slot="footer"
        class="dialog-footer"
      >
        <el-button
          size="small"
          @click="closeDialog"
        >取 消</el-button>
        <el-button
          size="small"
          type="primary"
          @click="submit"
        >保 存</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import API from "@/api/customer";
import { VUE_APP_API_ROOT_RE } from "@/utils/env-config";
import LargeFileUpload from "@/components/business/large-file-upload/index_o.vue";
type tagFormTye = {
  group_name: string;
  group_id: string | number;
  tag_name: string;
  tag_validity: number;
  expired_time: string | number;
  file_url: string;
};
@Component({
  components: {
    LargeFileUpload,
  },
})
export default class extends Vue {
  @Prop() currentEditItem!: any;
  checkGroupName = (rule: any, value: any, callback: any) => {
    callback();
  };
  checkName = (rule: any, value: any, callback: any) => {
    if (!value) {
      return callback(new Error("请输入标签名称"));
    }
    callback();
  };
  checkExpiredTime = (rule: any, value: any, callback: any) => {
    if (!value) {
      return callback(new Error("请选择时间"));
    }else if (value < Date.now()) {
      return callback(new Error("请选择大于当前时间的时间"));
    }
    callback();
  };
  /************data**************/
  show_dialog = false;
  rules = {
    group_name: [{ validator: this.checkGroupName }],
    tag_name: [{ validator: this.checkName }],
    expired_time: [{ validator: this.checkExpiredTime }],
  };
  ruleForm: tagFormTye = {
    group_name: "",
    group_id: "",
    tag_name: "",
    tag_validity: 1,
    expired_time: 0,
    file_url: "",
  };
  // 标签组下拉列表
  groupList: any = [];
  /************methods**************/
  // 有效期改变事件
  changeTimeGroup(val: any) {
    if (val == 1) {
      this.ruleForm.expired_time = 0;
    } else {
      this.ruleForm.expired_time = "";
    }
  }
  // 获取标签组下拉列表接口
  getGroupList() {
    API.labelListApi({}).then((res: any) => {
      if (res.code === 200) {
        this.groupList = res.data.manual || [];
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 输入标签组名称，清空标签组id
  changeInput() {
    this.ruleForm.group_id = "";
  }
  // 选择标签组
  changeGroup(val: any) {
    if (!val) {
      this.ruleForm.group_name = "";
      this.ruleForm.group_id = "";
      return;
    }
    // 获取groud_name
    const group_name = this.groupList.filter((item: any) => item.tag_group_id === val)[0].tag_group_name;
    this.ruleForm.group_name = group_name;
    this.ruleForm.group_id = val;
  }
  // 点击完成
  submit() {
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        if (this.currentEditItem.tag_id) {
          this.editTag()
        } else {
          this.createTag();
        }

      } else {
        return false;
      }
    });
  }
  // 编辑标签接口
  editTag() {
    const params = {
      tag_name: this.ruleForm.tag_name,
      tag_id: this.currentEditItem.tag_id
    };
    console.log(params);
    API.tagEditTagApi(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("标签编辑完成");
        this.$emit("submit", this.ruleForm.group_id);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 新增标签接口
  createTag() {
    // 判断是否是长期有效
    if (this.ruleForm.tag_validity == 2) {
      // 判断是否选择了时间
      if (this.ruleForm.expired_time) {
        this.ruleForm.expired_time = new Date(this.ruleForm.expired_time).getTime();
      }
    }
    // 判断是否上传了文件
    if (!this.ruleForm.file_url) {
      this.$message.warning("请上传文件");
      return;
    }
    const params = {
      group_name: this.ruleForm.group_name,
      group_id: Number(this.ruleForm.group_id),
      // 需要数组
      tag_names: [this.ruleForm.tag_name],
      expired_time: Number(this.ruleForm.expired_time) / 1000,
      file_url: this.ruleForm.file_url,
    };
    console.log(params);
    API.tagCreateGroupApi(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("标签添加完成");
        this.$emit("submit", this.ruleForm.group_id);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 关闭弹框
  closeDialog() {
    this.$emit("closeDialog");
  }
  // 模板下载
  temDown() {
    window.location.href = `${VUE_APP_API_ROOT_RE}tag/export-crowd-template`;
  }
  // 上传文件限制规范
  beforeUpload(file: any, type: string) {
    // 允许上传的文件类型
    const type_arr = ["xls", "xlsx", "csv"];
    const nameArr = file.name.split("."); // 通过文件后缀获取文件类型
    const tyle_suffix = nameArr[nameArr.length - 1];
    if (type === "fileUplod" && !type_arr.includes(tyle_suffix)) {
      this.$message.error("请上传xls, xlsx,csv格式文件");
      return false;
    }
    return true;
  }
  // 上传成功
  onSuccessFun(res: any) {
    this.ruleForm.file_url = res.data.url;
  }
  // 移除上传
  onRemoveFn() {
    this.ruleForm.file_url = "";
  }
  mounted() {
    this.getGroupList(); // 获取标签组下拉列表
    this.show_dialog = true;
    if (this.currentEditItem.tag_id) {
      this.ruleForm.tag_name = this.currentEditItem.tag_name
    }
  }
}
</script>
<style lang="scss" scoped>
.add-manual-label {
  .choose-tag {
    display: inline-block;
    margin-left: 10px;
  }

  ::v-deep .el-dialog {
    max-height: inherit !important;
  }

  ::v-deep.cy-select-small {
    width: 400px !important;
  }
  ::v-deep.el-input__inner {
     padding-right: 75px!important;
  }
}
</style>
