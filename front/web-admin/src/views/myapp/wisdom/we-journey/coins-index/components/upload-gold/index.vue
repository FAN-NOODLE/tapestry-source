<!-- 手动补发金币 -->
<template>
  <el-dialog
    :visible.sync="syncVisible"
    width="600px"
    :title="config.title"
    v-if="syncVisible"
  >
    <el-form
      label-width="100px"
      label-position="right"
    >
      <el-form-item label="导入模板">
        <span
          class="blue"
          @click="downLoad"
        >下载导入模板</span>
      </el-form-item>

      <el-form-item
        label="导入文档"
        required
      >
        <large-file-upload
          v-model="formData.file_url"
          :onSuccess="uploadSuccess"
          :onRemove="handleRemove"
          accept=".xlsx"
        >
          <el-button
            type="primary"
            size="small"
          >上传文件</el-button>
        </large-file-upload>
      </el-form-item>
    </el-form>

    <span
      slot="footer"
      class="dialog-footer"
    >
      <el-button @click="syncVisible = false">取 消</el-button>
      <el-button
        type="primary"
        :disabled="!formData.file_url"
        @click="save"
      >确 定</el-button>
    </span>
  </el-dialog>
</template>
<script lang="ts">
import { Component, PropSync, Mixins } from "vue-property-decorator";
import LargeFileUpload from "@/components/business/large-file-upload/index_o.vue";
import DownLoad from "@/mixins/DownLoad";
import { downloadCoinsTemp, importCoinsTemp } from "@/api/we-journey/gold";
import { ApiResult } from "@/types";

@Component({
  components: { LargeFileUpload }
})
export default class AddDialog extends Mixins(DownLoad) {
  /** 判断父属性显示影藏，双向绑定 */
  @PropSync("visible", { type: Boolean }) syncVisible!: boolean;
  /** 默认配置 */
  config = {
    /** 弹窗标题 */
    title: "手动补发金币"
  };

  /** 表单数据 */
  formData = {
    file_url: ""
  };

  /** 上传成功回调 */
  uploadSuccess({ data }: { data: { url: string } }) {
    this.formData.file_url = data.url;
  }

  /** 移除上传的文件 */
  handleRemove() {
    this.formData.file_url = "";
  }

  /** 下载导入模板 */
  async downLoad() {
    try {
      const res = await downloadCoinsTemp();

      const nowTime = new Date().getTime();
      const fileName = `${this.config.title}模板_` + nowTime + ".xlsx";
      this.downLoadFile(res, fileName);
    } catch (e) {
      console.error(e);
    }
  }

  /** 保存 */
  async save() {
    const res: ApiResult = await importCoinsTemp({
      file_url: this.formData.file_url
    });

    if (res.code !== 200) return this.$message.error(res.message as string);

    this.syncVisible = false;
    // 告知父组件导入成功
    this.$emit("success");
  }
}
</script>
<style lang="scss" scoped>.blue {
  color: #409eff;
  cursor: pointer;
}</style>
