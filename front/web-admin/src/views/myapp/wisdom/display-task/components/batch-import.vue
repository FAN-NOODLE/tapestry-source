
<template>
  <el-dialog
    :title="title"
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
          <el-button type="text" class="link" @click="handleDownload">{{downloadBtnText}}</el-button>
        </el-col>
      </el-row>
      <el-row class="m-t-20">
        <el-col :span="4">模板上传</el-col>
        <el-col :span="20">
          <el-upload
            v-bind="uploadProps"
            action=""
            :http-request="handleRequest"
            :before-upload="beforeUpload"
            :on-success="uploadSuccess"
            :on-error="uploadError"
            :file-list="fileList"
            :limit="1"
            :on-exceed="handleExceed"
            :on-remove="handleRemove"
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
            <span>上传{{uploadResult.upload_num}}，</span>
            <span>解析成功{{uploadResult.success}}，</span>
            <span>重复{{uploadResult.duplicate_num}}，</span>
            <span>失败{{uploadResult.failed_num}}条。</span>
            <el-button v-if="uploadResult.fail_url" type="text" class="link" @click="handleDownloadFail">下载失败数据</el-button>
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
import { Component, Vue, Model, Prop, Watch } from 'vue-property-decorator'
import { downloadFile } from '@/utils/download';
import { cosUploader } from '@/utils/upload';
import Api from '@/api/wisdom/displayTask';
import dayjs from 'dayjs';

@Component
export default class BatchImport extends Vue {
  @Model('update:visible', { type: Boolean, default: false }) visible!: boolean;
  @Prop({ type: String, default: '导入分配文件' }) title: string
  @Prop({ type: String, default: '下载分配文件demo' }) downloadBtnText: string
  @Prop({ type: String, default: '分配文件' }) filename: string
  @Prop({ type: String, default: '' }) templateUrl: string
  @Prop({ type: Object, default: () => ({}) }) params: any
  @Prop({ type: String, default: '' }) defaultValue: string
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
    upload_num: number;
    success: number;
    duplicate_num: number;
    failed_num: number;
    fail_url: string;
  } = null

  @Watch('visible')
  defaultValueChange(val: string) {
    if(val) {
      console.log('defaultValueChange', this.defaultValue);
      
      this.$nextTick(() => {
        this.file_url = this.defaultValue
        if(this.file_url) {
          this.fileList.push({
            name: `导入的分配文件.xlsx`,
            url: this.defaultValue
          })
        }
      })
    }
  }

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
    Api.importListApi({
      file_url: this.file_url
    }).then((res: any) => {
      if(res.code === 200) {
        this.$message.success("上传成功");
        this.uploadResult = res.data
      } else {
        this.fileList = []
        this.file_url = ''
        this.uploadResult = null
        this.$message.error(res.message);
      }
    }).catch(() => {
      this.fileList = []
      this.file_url = ''
      this.uploadResult = null
    })
  }

  /**
   * 上传失败
   */
  uploadError(res: any) {
    const error: any = JSON.parse(res.message).message;
    this.$message.error("上传失败,请重试");
    this.message = error;
    this.fileList = []
    this.file_url = ''
    this.uploadResult = null

  }

  /**
   * 上传超出限制
   */
  handleExceed() {
    this.$message.error("只能上传一个文件");
  }

  /**
   * 删除文件
   */
  handleRemove(file: any, fileList: any[]) {
    this.file_url = ''
    this.uploadResult = null
    this.fileList = []
  }


  /**
   * 下载模版
   */
  handleDownload() {
    if(this.templateUrl) {
      downloadFile(this.templateUrl, {
        filename: `${this.filename}模版.xlsx`,
      })
    } else {
      this.$emit('download')
    }
  }

  /**
   * 下载失败数据
   */
  async handleDownloadFail() {
    const url = await cosUploader.getObjectUrl(this.uploadResult.fail_url)
    downloadFile(url, {
      filename: `${this.filename}_导入失败数据.xlsx`,
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
    if(!this.file_url) {
      this.$message.error("请上传文件");
      return
    }
    this.$emit('update:visible', false)
    this.$emit('confirm', {
      url: this.file_url
    })
  }
}
</script>