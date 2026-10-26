<template>
  <div>
    <layout-form type="normal">
      <template slot="content">
        <cy-form ref="cform" :rules="rules" :model="cardForm">
          <cy-form-item label="模板名称" prop="template_name" required>
            <cy-input
              v-model="cardForm.template_name"
              clearable
              width-size="medium"
              placeholder="请输入模板名称"
            ></cy-input>
          </cy-form-item>
          <cy-form-item
            style="max-width:1000px"
            label="背景图"
            prop="template_picture_url"
            required
            trigger="change"
          >
            <upload-image
              class="upload-template"
              :urls="cardForm.template_picture_url"
              accept="image/jpeg,image/png"
              img-type="batch"
              :before-upload="beforeUpload"
              @update="update"
              :limit="limiteNum"
              multiple
            >
            </upload-image>
            <span class="upload-tip"
              >可上传多张图片，建议上传尺寸为750*1624的图片，大小不超过500K，建议上传小于200kb</span
            >
          </cy-form-item>
          <cy-form-item label="文案" prop="template_content" required>
            <cy-input
              v-model="cardForm.template_content"
              clearable
              width-size="medium"
              type="textarea"
              maxlength="50"
              show-word-limit
              :autosize="{ minRows: 3, maxRows: 10 }"
              placeholder="请输入文案"
            ></cy-input>
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="footer">
        <el-button type="primary" @click="saveClickFn">保存</el-button>
        <el-button @click="goCancelFn">取消</el-button>
      </template>
    </layout-form>
    <!-- 创建离开页面弹框 -->
    <el-dialog
      title="离开页面提示"
      :visible.sync="isCancelDialog"
      width="420px"
    >
      <span
        ><i aria-hidden="true" class="el-icon-warning"></i>
        检测到未保存的内容，是否保存到下次新增继续编辑？</span
      >
      <span slot="footer" class="dialog-footer">
        <el-button @click="goBackFn">直接离开</el-button>
        <el-button type="primary" @click="saveGoBackFn">保存并离开</el-button>
      </span>
    </el-dialog>
    <!-- 编辑离开页面弹框 -->
    <el-dialog
      title="离开页面提示"
      :visible.sync="isTipsCancel"
      width="420px"
      @close="goCloseFn"
    >
      <span
        ><i aria-hidden="true" class="el-icon-warning"></i>
        你确定要放弃编辑吗？放弃后编辑的信息将不会被保存</span
      >
      <span slot="footer" class="dialog-footer">
        <el-button @click="goCloseFn">取消</el-button>
        <el-button type="primary" @click="editGoBackFn">确定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import API from "@/api/customer/greeting-card";
import UploadImage from "@/components/business/upload-image/index.vue";
@Component({
  components: { UploadImage }
})
export default class extends Vue {
  @Prop() type?: string;

  /*************************** data start ***************************/
  private isCancelDialog = false; //取消按钮是否弹框
  private isTipsCancel = false; //编辑取消按钮是否弹框
  private cardId = 0; //当前id
  private cardForm = {
    template_name: "",
    template_picture_url: [],
    template_content: ""
  };
  private preObj = {
    template_name: "",
    template_picture_url: [],
    template_content: ""
  };
  private rules = {};
  private isUploading = false;
  private limiteNum = 10;
  /*************************** method start ***************************/
  private created() {
    if (this.type == "create") {
      this.createShow();
    } else if (this.type == "edit") {
      this.cardId = Number(this.$route.query.id);
      this.getDetailFn();
    }
  }

  //保存
  private saveClickFn() {
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        const urls = this.cardForm.template_picture_url.map(
          (item: any) => item.url
        );
        const params = {
          ...this.cardForm,
          template_picture_url: urls,
          id: this.cardId
        };
        if (this.type !== "edit") {
          delete params.id;
        }
        API.templateInfo(params).then((res: any) => {
          if (res.code === 200) {
            this.$message.success("提交成功");
            this.routeBackFn();
          } else {
            this.$message.error(res.message);
          }
        });
      }
    });
  }
  //取消
  private goCancelFn() {
    if (this.type === "create") {
      this.createCancelFn();
    } else if (this.type === "edit") {
      this.editCancelFn();
    }
  }
  //创建取消
  private createCancelFn() {
    if (
      this.cardForm.template_name ||
      this.cardForm.template_picture_url.length != 0 ||
      this.cardForm.template_content
    ) {
      this.isCancelDialog = true;
    } else {
      this.goBackFn();
    }
  }
  //编辑取消
  private editCancelFn() {
    if (
      this.cardForm.template_name != this.preObj.template_name ||
      this.cardForm.template_content != this.preObj.template_content ||
      JSON.stringify(this.cardForm.template_picture_url) !=
        JSON.stringify(this.preObj.template_picture_url)
    ) {
      this.isTipsCancel = true;
    } else {
      this.editRouteBackFn();
    }
  }
  //返回上一页
  private goBackFn() {
    this.isCancelDialog = false;
    this.routeBackFn();
  }
  //保存返回上一页
  private saveGoBackFn() {
    this.isCancelDialog = false;
    // 新建的值保存sessionStorage
    sessionStorage.setItem("create_card_obj", JSON.stringify(this.cardForm));
    this.routeBackFn();
  }
  //路由跳转
  private routeBackFn() {
    this.$router.push({
      name: "myapp_customer_greeting-card_list"
    });
  }
  //路由跳转
  private editRouteBackFn() {
    this.$router.push({
      name: "myapp_customer_greeting-card_list",
      query: { id: this.cardId.toString() }
    });
  }

  //编辑关闭和取消
  private goCloseFn() {
    this.isTipsCancel = false;
  }
  //编辑确定
  private editGoBackFn() {
    this.isTipsCancel = true;
    this.editRouteBackFn();
  }

  //获取详情
  private getDetailFn() {
    API.templateDetail({ id: this.cardId }).then((res: any) => {
      if (res.code === 200) {
        const urls = res.data.template_picture_url.map(
          (item: any, index: number) => {
            return {
              name: "name_" + index,
              url: item
            };
          }
        );
        const params = {
          template_name: res.data.template_name,
          template_content: res.data.template_content,
          template_picture_url: urls
        };
        // 页面提交数据
        this.cardForm = JSON.parse(JSON.stringify(params));
        // 详情备份，用来点击取消按钮进行数据对比
        this.preObj = JSON.parse(JSON.stringify(params));
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //创建页面进来，拿sessionStorage值
  private createShow() {
    this.cardForm = JSON.parse(
      sessionStorage.getItem("create_card_obj") || JSON.stringify(this.cardForm)
    );
  }
  // 上传文件限制规范
  private beforeUpload(file: any) {
    if (file.type != "image/jpeg" && file.type != "image/png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }
    const name = file.name
      .substring(file.name.lastIndexOf(".") + 1, file.name.length)
      .toLocaleLowerCase();
    if (name === "jpeg") {
      this.$message.warning("企业微信不支持JPEG的图片");
      return false;
    }
    // if (file.size > 2 * 1024 * 1024) {
    if (file.size > 1024 * 500) {
      this.$message.error("请上传500kb大小以内的图片");
      return false;
    }
    this.isUploading = true;
  }
  private update(res: any) {
    this.cardForm.template_picture_url = res;
    this.isUploading = false;
  }
  /*************************** 生命周期 start ***************************/
}
</script>
<style lang="scss" scoped>
.upload-tip {
  font-size: 12px;
  color: #666666;
}
</style>
