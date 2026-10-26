<template>
  <div>
    <div class="add-link">
      <cy-form ref="cform" :rules="rules" :model="formData" label-width="90px">
        <cy-form-item label="素材分组" prop="group_id">
          <cy-select
            v-model="formData.group_id"
            size="medium"
            placeholder="请选择素材分组"
            search-input
            :popperAppendToBody="true"
          >
            <cy-option
              :label="item.group_name"
              :value="item.id"
              :key="index"
              v-for="(item, index) in newGroupData"
            ></cy-option>
          </cy-select>
        </cy-form-item>
        <cy-form-item label="链接标题" prop="title" required>
          <cy-input
            ref="inp"
            placeholder="请输入链接标题"
            v-model.trim="formData.title"
            size="medium"
            maxlength="30"
            show-word-limit
          ></cy-input>
        </cy-form-item>
        <cy-form-item label="链接描述" prop="description" required>
          <cy-input
            ref="inp"
            placeholder="请输入链接描述"
            v-model.trim="formData.description"
            size="medium"
            maxlength="30"
            show-word-limit
          ></cy-input>
        </cy-form-item>
        <cy-form-item label="封面图" required prop="cover_url" trigger="change">
          <upload-image
            class="upload-template"
            :urls="formData.cover_url"
            accept="image/jpeg,image/png"
            @change="picSuccess"
            :before-upload="beforeUpload"
          ></upload-image>
          <div class="upload-tip">支持JPG、PNG格式，图片＜2MB</div>
        </cy-form-item>
        <cy-form-item label="链接内容" prop="link_type">
          <el-radio
            v-model="formData.link_type"
            :label="1"
            @change="radioChange"
            >URL
          </el-radio>
          <el-radio
            v-model="formData.link_type"
            :label="2"
            @change="radioChange"
            >自定义</el-radio
          >
        </cy-form-item>
        <cy-form-item
          label="URL"
          required
          v-if="formData.link_type == 1"
          prop="link_url"
        >
          <cy-input
            ref="inp"
            placeholder="请输入URL"
            v-model.trim="formData.link_url"
            size="medium"
          ></cy-input>
        </cy-form-item>
        <cy-form-item
          label="作者"
          required
          v-if="formData.link_type == 2"
          prop="author"
        >
          <cy-input
            ref="inp"
            placeholder="请输入作者名称"
            v-model.trim="formData.author"
            size="medium"
          ></cy-input>
        </cy-form-item>
        <cy-form-item
          label="正文"
          required
          v-if="formData.link_type == 2"
          prop="content"
        >
          <div id="editContent"></div>
        </cy-form-item>
      </cy-form>
    </div>
    <div class="footer">
      <el-button @click="goBack" size="small">返 回</el-button>
      <el-button type="primary" @click="addConfirm" size="small"
        >确 认</el-button
      >
    </div>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import wangEditor from "wangeditor";
import Api from "@/api/enterprise/material-library";
import UploadImage from "@/components/business/upload-image/index.vue";
@Component({
  components: {
    UploadImage
  }
})
export default class extends Vue {
  //------------------------prop-------------------------------
  @Prop() newGroupData?: any;
  //-------------------data-----------------------
  private formData = {
    title: "", //素材标题
    description: "", //素材描述
    group_id: "",
    link_type: 1, //(1 URL 2 自定义)
    cover_url: "", //封面图URL
    link_url: "", //链接URL
    author: "", //作者
    content: "" //正文内容
  };
  private editorSms: any = null; //编辑器
  private rules = {
    link_url: [
      {
        validator: this.validateself,
        trigger: "blur"
      }
    ]
  };
  private validateself(rule: any, value: any, callback: any) {
    const match = /^(https?:\/\/)(([a-zA-Z0-9_-]+\.)+[a-zA-Z]{2,})(:\d+)?(\/[^\s?#]*)?(\?[^\s#]*)?(#[^\s]*)?$/;
    const flag = match.test(value);
    flag ? callback() : callback(new Error("请输入正确的URL"));
  }
  //------------------------method--------------------------
  private radioChange() {
    if (this.formData.link_type == 2) {
      this.editorSms = new wangEditor(`#editContent`);
      this.editorSms.config.zIndex = 500;
      this.editorSms.create();
      if (this.formData.content) {
        this.editorSms.txt.html(this.formData.content);
      }
    }
    this.$nextTick(() => {
      (this.$refs["cform"] as any).clearValidate(["link_url"]);
      (this.$refs["cform"] as any).clearValidate(["content"]);
      (this.$refs["cform"] as any).clearValidate(["author"]);
    });
  }
  private addConfirm() {
    if (this.editorSms) {
      this.$set(this.formData, "content", this.editorSms.txt.html());
    }
    const data = JSON.parse(JSON.stringify(this.formData));
    if (this.formData.link_type == 1) {
      data.author = "";
      data.content = "";
    } else {
      data.link_url = "";
      data.content = this.editorSms.txt.html();
    }
    data.group_id = parseInt(data.group_id) || -1;
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        Api.linkSaveApi(data).then((res: any) => {
          if (res.code == 200) {
            this.$message.success("新增链接成功");
            this.$emit("submit");
          } else {
            this.$message.error(res.message);
          }
        });
      } else {
        return;
      }
    });
  }
  private goBack() {
    this.$emit("submit");
  }
  //
  private picSuccess(res: any) {
    this.$set(this.formData, "cover_url", res);
    (this.$refs["cform"] as any).validateField("cover_url");
  }
  // 上传文件限制规范
  private beforeUpload(file: any) {
    const name = file.name
      .substring(file.name.lastIndexOf(".") + 1, file.name.length)
      .toLocaleLowerCase();
    if (name !== "jpg" && name !== "png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }
    if (file.size > 2 * 1024 * 1024) {
      this.$message.error("请上传2MB大小以内的图片");
      return false;
    }
  }
}
</script>
<style lang="scss" scoped>
.add-link {
  max-height: 500px;
  overflow: auto;
}
.footer {
  text-align: right;
  margin-top: 20px;
}
#editContent {
  width: 700px;
}
</style>
