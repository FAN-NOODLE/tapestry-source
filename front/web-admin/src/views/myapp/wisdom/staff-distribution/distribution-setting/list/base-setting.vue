<template>
  <div class="base-setting">
    <el-form ref="form" :rules="rules" :model="form" label-width="86px">
      <cy-form-item label="公告标题" prop="distribution_title" required>
        <cy-input
          type="text"
          placeholder="请输入公告标题"
          v-model.trim="form.distribution_title"
          clearable
          width-size="medium"
          maxlength="20"
          :show-word-limit="true"
        ></cy-input>
      </cy-form-item>
      <cy-form-item label="公告内容" prop="distribution_text" required>
        <cy-input
          type="textarea"
          placeholder="请输入公告内容"
          v-model.trim="form.distribution_text"
          :rows="6"
          width-size="medium"
          maxlength="200"
          :show-word-limit="true"
        ></cy-input>
      </cy-form-item>
      <cy-form-item label="执行方式">
        <div class="settlement">已完成订单</div>
        <div class="remind-info">
          分销订单达到已完成状态时（即订单状态再不会发生变化），结算分销佣金，结算单变更为已结算。
        </div>
      </cy-form-item>
      <cy-form-item label="海报背景" prop="poster_img" required>
        <div class="upload-inner">
          <upload-image
            class="upload-template"
            :urls="form.poster_img"
            :action="action"
            @change="onChange"
            :on-success="(file) => bgPicSuccess(file)"
            :on-error="uploadError"
            :before-upload="
              (file) =>
                _uploadValidate(
                  file,
                  [
                    { type: 'png', limit: 500 },
                    { type: 'jpg', limit: 500 },
                  ],
                  { width: 750, height: 980 }
                )
            "
            :on-remove="() => deleteImgs()"
          ></upload-image>
          <div class="upload-tips remind-info">
            支持jpg、png格式 请上传尺寸为750*980，图片&lt;500kb
            <el-button
              class="m-l-5"
              type="text"
              size="mini"
              @click="switchImageDialog(true)"
            >
              查看示例图
            </el-button>
          </div>
        </div>
      </cy-form-item>
    </el-form>
    <div class="footer">
      <el-button type="primary" size="small" @click="save">保存</el-button>
    </div>

    <image-modal
      :type="image_type"
      :dialogVisible="image_show"
      @close="switchImageDialog"
    ></image-modal>
  </div>
</template>

<script lang="ts">
import { Component, Mixins } from "vue-property-decorator";
import UploadImage from "@/components/business/upload-image/index.vue";
import ImageModal from "../modal/image-modal.vue";
import UploadValidate from "@/mixins/UploadValidate";
import { Loading } from "element-ui";
import DistributionApi from "@/api/wisdom/Distribution";
@Component({
  components: {
    UploadImage,
    ImageModal,
  },
})
export default class BaseSetting extends Mixins(UploadValidate) {
  private action = "";
  private form = {
    distribution_title: "", // 分销公告
    distribution_text: "", // 分销描述
    settlement_standard: 100, // 结算状态，默认为已完成订单
    poster_img: "", // 海报背景图
  };
  private rules = {
    distribution_title: [{ required: true, message: "请输入公告标题" }],
    distribution_text: [{ required: true, message: "请输入公告内容" }],
    poster_img: [{ required: true, message: "请上传海报背景图" }],
  };
  private img_modal = false; // 展示查看海报放大图
  private image_show = false;
  private distribution_id: any = ""; // 分销商品id
  private image_type = "example";

  private bgPicSuccess() {
    console.log("res");
  }
  // 海报上传成功
  private onChange(url: string) {
    console.log(url, "res");
    this.$set(this.form, "poster_img", url);
    Loading.service({ fullscreen: true }).close();
  }

  // 海报上传失败
  private uploadError() {
    this.$message.error("上传失败，请重新上传");
    Loading.service({ fullscreen: true }).close();
  }

  // 删除海报
  private deleteImgs() {
    this.$set(this.form, "poster_img", "");
  }

  // 查看示例图
  private switchImageDialog(isShow: boolean) {
    this.image_show = isShow;
  }

  // 保存
  save() {
    (this.$refs["form"] as any).validate((valid: boolean) => {
      if (valid) {
        this.saveApi();
      }
    });
  }

  // 保存接口
  saveApi() {
    const params = {
      distribution_title: this.form.distribution_title,
      distribution_text: this.form.distribution_text,
      settlement_standard: this.form.settlement_standard,
      poster_img: this.form.poster_img,
      id: this.distribution_id,
    };
    DistributionApi.updateDistributionListApi(params).then(() => {
      this.$message.success("保存成功");
      this.getDistributionDetail();
    });
  }

  // 获取分销查看(分销描述/海报背景)
  getDistributionDetail() {
    // 分销查看数据
    DistributionApi.getDistributionDetail({
      distribution_id: this.distribution_id,
    }).then((res) => {
      this.distribution_id = res.data.id;
      this.form.distribution_title = res.data.distribution_title;
      this.form.distribution_text = res.data.distribution_text;
      this.form.poster_img = res.data.poster_img;
    });
  }

  // 生命周期
  private mounted() {
    this.getDistributionDetail();
    this.action = `${this.$VUE_APP_API_VERSION}file/upload`;
  }
}
</script>

<style lang="scss" scoped>
.base-setting {
  color: #303133;
}
/deep/.el-input__inner,
.el-textarea__inner {
  font-family: PingFangSC-Regular, PingFang SC;
}
.remind-info {
  white-space: nowrap;
  line-height: 18px;
  font-size: 12px;
  color: $color-text-second;
}
.upload-inner {
  display: flex;
  align-items: flex-end;
  .upload-tips {
    margin-left: 15px;
  }
}
.footer {
  padding-top: 20px;
  border-top: 1px solid $color-third-line;
}
</style>
