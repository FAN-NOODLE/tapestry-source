<!-- 新增积分规则弹窗 -->
<template>
  <el-dialog
    :visible.sync="syncVisible"
    width="600px"
    :title="config.title"
    v-if="syncVisible"
  >
    <el-form label-width="100px" label-position="right">
      <el-form-item label="导入模板">
        <span class="blue" @click="downLoad">下载导入模板</span>
      </el-form-item>

      <el-form-item label="导入文档" required>
        <large-file-upload
          v-model="formData.file_url"
          :onSuccess="uploadSuccess"
          :onRemove="handleRemove"
          accept=".xls,.xlsx"
        >
          <el-button type="primary" size="small">上传文件</el-button>
        </large-file-upload>
      </el-form-item>
    </el-form>

    <span slot="footer" class="dialog-footer">
      <el-button @click="syncVisible = false">取 消</el-button>
      <el-button type="primary" :disabled="!formData.file_url" @click="save"
        >确 定</el-button
      >
    </span>
  </el-dialog>
</template>
<script lang="ts">
import { Component, PropSync, Mixins, Prop } from "vue-property-decorator";
import LargeFileUpload from "@/components/business/large-file-upload/index_o.vue";
import DownLoad from "@/mixins/DownLoad";
import { downloadTemplateApi, importPointApi } from "@/api/we-journey/Integral";
import {
  downloadReissueApi,
  importMedalApi,
  recordReissueApi
} from "@/api/we-journey/medal";
import { ApiResult } from "@/types";

@Component({
  components: { LargeFileUpload }
})
export default class AddDialog extends Mixins(DownLoad) {
  /** 判断父属性显示影藏，双向绑定 */
  @PropSync("visible", { type: Boolean }) syncVisible!: boolean;

  /** 类型 1-积分 2-奖章*/
  @Prop() type!: number;

  /** 默认配置 */
  private config = {
    /** 弹窗标题 */
    title: "手动补发积分",
    downloadApi: downloadTemplateApi,
    importApi: importPointApi
  };

  /** 表单数据 */
  private formData = {
    file_url: ""
  };

  // private beforeUpload(file: File) {
  //   const fileSuffix = file.name.substring(file.name.lastIndexOf(".") + 1);
  //   const whiteList = ["xls", "xlsx"];
  //   if (whiteList.indexOf(fileSuffix) === -1) {
  //     this.$message.error("上传文件只能是xls和xlsx");
  //     return false;
  //   }
  //   return true;
  // }

  /** 上传成功回调 */
  private uploadSuccess({ data }: { data: { url: string } }) {
    this.formData.file_url = data.url;
  }

  /** 移除上传的文件 */
  private handleRemove() {
    this.formData.file_url = "";
  }

  /** 下载导入模板 */
  private async downLoad() {
    try {
      const res = await this.config.downloadApi();

      const nowTime = new Date().getTime();
      const fileName = `${this.config.title}模板_` + nowTime + ".xlsx";
      this.downLoadFile(res, fileName);
    } catch (e) {
      console.error(e);
    }
  }

  /** 保存 */
  private async save() {
    const res: ApiResult = await this.config.importApi({
      file_url: this.formData.file_url
    });

    if (res.code !== 200) return this.$message.error(res.message as string);

    // 由于奖章导入变更，所以需要单独调用接口
    if (this.type !== 1) return this.saveMedal(res.data);

    this.syncVisible = false;
    // 告知父组件导入成功
    this.$emit("success");
  }

  /** 保存奖章数据人群包 */
  private async saveMedal(reissue_data: any) {
    const res: ApiResult = await recordReissueApi({ reissue_data });

    if (res.code !== 200) return this.$message.error(res.message as string);

    this.syncVisible = false;
    // 告知父组件导入成功
    this.$emit("success");
  }

  /** 生命周期 */
  private created() {
    if (this.type === 1) {
      // 积分
      this.config = {
        title: "手动补发积分",
        downloadApi: downloadTemplateApi,
        importApi: importPointApi
      };
    } else {
      // 奖章
      this.config = {
        title: "手动补发奖章",
        downloadApi: downloadReissueApi,
        importApi: importMedalApi
      };
    }
  }
}
</script>
<style lang="scss" scoped>
.blue {
  color: #409eff;
  cursor: pointer;
}
</style>
