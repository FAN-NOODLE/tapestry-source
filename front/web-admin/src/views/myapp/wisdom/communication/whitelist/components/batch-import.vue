
<template>
  <el-dialog
    title="导入员工人群包"
    :visible="visible"
    :append-to-body="false"
    :modal-append-to-body="false"
    @open="handleOpen"
    @close="handleCancel"
  >
    <div class="dialog-wrapper">
      <el-row class="m-t-20">
        <el-col :span="4">模板下载</el-col>
        <el-col :span="20">
          <el-button type="text" class="link" @click="handleDownload">下载门店code数据模板</el-button>
          <span class="tips">*输入门店code可自动获取门店下员工ID</span>
        </el-col>
      </el-row>
      <el-row class="m-t-20">
        <el-col :span="4">模板上传</el-col>
        <el-col :span="20">
          <el-upload
            v-bind="uploadProps"
            :http-request="handleRequest"
            :before-upload="beforeUpload"
            :on-success="uploadSuccess"
            :on-error="uploadError"
            :file-list="fileList"
            show-file-list
          >
            <i class="el-icon-upload"></i>
            <div class="el-upload__text">
              将文件拖到此处，或<em>点击上传</em>
            </div>
            <div
              class="el-upload__tip"
            >只能上传xls、xlsx文件</div>
          </el-upload>
        </el-col>
      </el-row>
      <el-row class="m-t-20">
        <el-col :span="4">上传结果</el-col>
        <el-col :span="20">
          <p v-if="uploadResult">
            <span>上传{{uploadResult.result.import}}，</span>
            <span>解析成功{{uploadResult.result.success}}，</span>
            <span>重复{{uploadResult.result.repeat}}，</span>
            <span>失败{{uploadResult.result.fail}}条。</span>
            <el-button v-if="uploadResult.result.fail" type="text" class="link" @click="handleDownloadFail">下载失败数据</el-button>
          </p>
          <p v-else>{{ message }}</p>
        </el-col>
      </el-row>
    </div>
    <div slot="footer" class="dialog-footer">
      <el-button @click="handleCancel">取 消</el-button>
      <el-button type="primary" @click="handleConfirm">确 定</el-button>
    </div>
  </el-dialog>
</template>
<script lang="ts">
import { Component, Vue, Model, Prop } from 'vue-property-decorator'
import LoginAboat from "@/utils/LoginManager";
import { VUE_APP_API_ROOT_RE } from "@/utils/env-config";
import { downloadFile } from '@/utils/download';
import { cosUploader } from '@/utils/upload';
import Api from '@/api/wisdom/communication';

@Component
export default class BatchImport extends Vue {
  @Model('update:visible', { type: Boolean, default: false }) visible!: boolean;
  /** 提示信息 */
  message = '暂未检测到上传信息'
  /** 导入结果 */
  result: any[] = []
  /** 上传属性 */
  uploadProps = {
    multiple: false,
    drag: true,
    accept: ".xls,.xlsx",
    fileList: [],
    showFileList: false
  }
  fileList = []

  /** 文件地址 */
  file_url = ''
  /** 上传结果 */
  uploadResult: {
    result: {
      import: number;
      success: number;
      fail: number;
      repeat: number;
      fail_store_code: Record<string, string>;
    };
    member_list: any[];
  } = null

  /**
   * 上传前校验
   */
  beforeUpload(file: any) {
    const fileSuffix = file.name.substring(file.name.lastIndexOf("."));
    const whiteList = this.uploadProps.accept.split(',');
    if (whiteList.indexOf(fileSuffix) === -1) {
      this.$message.error("上传文件只能是xls和xlsx");
      return false;
    }
    return true
  }

  handleRequest({file, onSuccess, onError}) {
    cosUploader.uploadFile(file).then(res => {
      this.file_url = res.url
      this.fileList = [res.url]
      onSuccess(res)
    }).catch(onError)
  }

  /**
   * 上传成功
   * @param res
   */
  uploadSuccess() {
    Api.importWhitelistApi({
      file_url: this.file_url
    }).then((res: any) => {
      if(res.code === 200) {
        this.$message.success("上传成功");
        this.uploadResult = res.data
      } else {
        this.fileList = []
        this.$message.error(res.message);
      }
    })
  }

  /**
   * 上传失败
   */
  uploadError(res: any) {
    const error: any = JSON.parse(res.message).message;
    this.$message.error("上传失败,请重试");
    this.message = error;
  }


  /**
   * 下载模版
   */
  handleDownload() {
    const url = `${VUE_APP_API_ROOT_RE}notification/download-crowd-template`
    downloadFile(url, {
      filename: '传讯_白名单_导入模版.xlsx',
      headers: {
        "Auth-Token": LoginAboat.getToken(),
      },
    })
  }

  /**
   * 下载失败数据
   */
  handleDownloadFail() {
    const url = `${VUE_APP_API_ROOT_RE}notification/export-fail-import-crowd`
    downloadFile(url, {
      method: "POST",
      params: {
        fail_store_codes: this.uploadResult.result.fail_store_code
      },
      filename: '传讯_白名单_导入失败数据.xlsx',
      headers: {
        "Auth-Token": LoginAboat.getToken(),
      },
      
    })
  }

  handleOpen() {
    this.uploadResult = null
    this.fileList = []
  }

  /**
   * 取消
   */
  handleCancel() {
    this.$emit('update:visible', false)
  }

  /**
   * 确定
   */
  handleConfirm() {
    const { member_list } = this.uploadResult ||{}
    if (!member_list ||member_list.length === 0) {
      this.$message.error('暂未检测到上传信息')
      return
    }
    this.$emit('update:visible', false)
    this.$emit('confirm', {
      file_url: this.file_url,
      member_list: member_list || []
    })
  }
}
</script>