<template>
  <div class="content-block">
    <layout-form>
      <template slot="content">
        <cy-form
          :model="formData"
          ref="form"
          :rules="fromRule"
          label-position="left"
          label-width="90px"
        >
          <cy-form-item label="企业名称">
            <span> {{ formData.enterprise_name }}</span>
          </cy-form-item>
          <cy-form-item label="企业简称" prop="short_name" required>
            <cy-input
              width-size="medium"
              placeholder="请输入企业简称"
              maxlength="8"
              v-model.trim="formData.short_name"
              show-word-limit
            ></cy-input>
          </cy-form-item>
          <cy-form-item label="企业logo" widthSize="100%">
            <div class="flex">
              <upload-image
                :urls="formData.logo"
                @change="bannerPicSuccess"
                :before-upload="(file) => beforeBackgroundUpload(file)"
                :on-error="uploadError"
              ></upload-image>
              <div class="message-block">
                <div class="message-top">新头像不允许涉及政治敏感与色情;</div>
                <div class="message-bottom">
                  图片格式必须为：png,bmp,jpeg,jpg,gif；不可大于2M；建议使用png格式图片，以保持最佳效果；
                </div>
              </div>
            </div>
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="footer">
        <div>
          <el-button type="primary" @click="saveData" size="small">
            保存
          </el-button>
          <!-- <el-button @click="goBack" size="small">取消</el-button> -->
        </div>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import UploadImage from "@/components/business/upload-image/index.vue";
import Api from "@/api/setting/public";
import { Loading } from "element-ui";
interface IDataForm {
  enterprise_name: string; // 企业名称
  short_name: string; // 企业简介
  logo: string; // 企业logo
}
@Component({
  components: {
    UploadImage
  }
})
export default class extends Vue {
  private formData: IDataForm = {
    enterprise_name: "",
    short_name: "",
    logo: ""
  };
  // 表单校验
  private fromRule = {
    short_name: [{ message: "请填写企业简称" }]
  };

  // 编辑商户信息
  private goEdit() {
    this.$router.push({
      name: "myapp_setting_basic-setup_merchant-info_edit"
    });
  }
  // 商户信息接口调用
  private getDetail() {
    Api.merchantDetail().then((res: any) => {
      const data: IDataForm = res.data;
      this.formData = { ...data };
    });
  }
  // 删除图片
  private bannerPicRemove() {
    this.$set(this.formData, "logo", "");
  }
  /**
   * @param {string} res 图片url地址
   * @description url地址改变执行此回调
   */
  private bannerPicSuccess(res: string) {
    this.$set(this.formData, "logo", res);
    Loading.service({ fullscreen: true }).close();
  }
  // 图片上传报错
  uploadError() {
    this.$message.error("上传失败，请重新上传");
    // 上传失败，关闭loading
    Loading.service({ fullscreen: true }).close();
  }
  /**
   * @param {object} res 文件对象
   * @description 判断上传文件是否合法
   */
  private beforeBackgroundUpload(file: any) {
    const types = ["image/jpeg", "image/png", "image/gif"];
    const type = file.type;
    const self = this;
    if (types.indexOf(type) < 0) {
      self.$message.error("上传的图片只能是 JPG/JPEG/BMP/PNG/GIF 格式!");
      return false;
    }
    if (file.size / 1024 > 2048) {
      self.$message.error("上传的图片大小不能超过 2M !");
      return false;
    }
    const isSize = new Promise<void>((resolve, reject) => {
      const width = 144; // 图片宽度
      const height = 144; // 图片高度
      const url = window.URL || window.webkitURL;
      const img = new Image();
      img.onload = function() {
        // if (img.width !== width) {
        //   reject("上传的图片尺寸必须为144px*144px!");
        // }
        // if (img.height !== height) {
        //   reject("上传的图片尺寸必须为144px*144px!");
        // }
        resolve();
      };
      img.src = url.createObjectURL(file);
    }).then(
      () => {
        Loading.service({ fullscreen: true, text: "图片上传中，请稍后" });
        return file;
      },
      (err: any) => {
        this.$message.error(err);
        return Promise.reject();
      }
    );
    return isSize;
  }
  // 取消编辑并返回上一页
  private goBack() {
    this.$router.go(-1);
  }
  // 保存数据
  private saveData() {
    (this.$refs["form"] as HTMLFormElement).validate((valid: boolean) => {
      if (!valid) {
        return false;
      } else {
        const params = {
          logo: this.formData.logo, //分享图片
          short_name: this.formData.short_name //员工分享文案
        };
        Api.merchantSave(params).then(() => {
          this.$message.success("保存成功");
        });
      }
    });
  }
  mounted() {
    this.getDetail();
  }
}
</script>
<style lang="scss" scoped>
.content-block {
  .flex {
    display: flex;
  }
  .message-block {
    line-height: 20px;
    margin-left: 15px;
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
  }
  .message-top,
  .message-bottom {
    font-size: 12px;
    color: $color-text-second;
  }
}
</style>
