
<template>
  <el-dialog
    :title="title"
    :visible="visible"
    :append-to-body="false"
    :modal-append-to-body="false"
    @opened="handleOpen"
    @close="handleCancel"
  >
    <div class="dialog-wrapper">
      <el-form ref="form" :model="formData" :rules="formRules" label-width="110px">
        <el-form-item label="门店" prop="third_id">
          <el-select v-model="formData.third_id" >
            <el-option v-for="item, index in storeList" :label="item.name" :value="item.third_id" :key="index"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="模版上传" prop="file_url">
          <el-upload
            v-bind="uploadProps"
            :limit="1" 
            :http-request="handleRequest"
            :before-upload="beforeUpload"
            :on-success="uploadSuccess"
            :on-error="uploadError"
            :on-exceed="handleExceedTransfer"
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
        </el-form-item>
        <el-form-item label="模板下载">
          <el-button type="text" class="link" @click="handleDownload">下载导入模版示例</el-button>
        </el-form-item>
        <el-form-item label="上传转货单" prop="transfer_file.url" :rules="formRules.transfer_file">
          <el-upload 

            :multiple="false"
            :limit="1" 
            :show-file-list="false"
            :file-list="fileList2"
            :before-upload="handleUploadTransferBefore"
            :http-request="handleRequestTransfer"
            :on-success="handleUploadTransferSuccess"
            :on-exceed="handleExceedTransfer">
            <el-button type="primary">上传转货单</el-button>
            <span slot="tip" class="el-upload__tip m-l-10">*可支持上传word、excel、ppt、pdf、图片</span>
          </el-upload>
          <el-row v-if="formData.transfer_file && formData.transfer_file.name" class="item">
            <div class="item-value">
              <div class="item-name">
                <svg-icon class="icon" :icon-class="formData.transfer_file.type"></svg-icon>
                <div class="name">
                  <div class="name-text">{{formData.transfer_file.name}}</div>
                  <el-tag v-if="formData.transfer_file._status === 'error'" class="m-l-10" type="danger" size="mini">上传失败</el-tag>
                </div>
                <el-button v-if="formData.transfer_file._status !== 'uploading'" type="text" @click="handleUploadRemove()" icon="el-icon-close"></el-button>
              </div>
              <el-progress v-if="formData.transfer_file._status === 'uploading'" color="#000" :percentage="formData.transfer_file._progress" :status="formData.transfer_file._status" :show-text="false" />
            </div>
          </el-row>
        </el-form-item>
        <el-form-item label="上传结果">
          <p v-if="uploadResult">
            <span>上传{{uploadResult.upload_num}}，</span>
            <span>解析成功{{uploadResult.success}}，</span>
            <span>重复{{uploadResult.duplicate_num}}，</span>
            <span>失败{{uploadResult.failed_num}}条。</span>
            <el-button v-if="uploadResult.fail_url" type="text" class="link" @click="handleDownloadFail">下载失败数据</el-button>
          </p>
          <p v-else>{{ message }}</p>
        </el-form-item>
      </el-form>
    </div>
    <div slot="footer" class="dialog-footer">
      <el-button @click="handleCancel">取 消</el-button>
      <el-button type="primary" @click="handleConfirm">确 定</el-button>
    </div>
  </el-dialog>
</template>
<script lang="ts">
import { Component, Vue, Model, Prop } from 'vue-property-decorator'
import { downloadFile } from '@/utils/download';
import { cosUploader } from '@/utils/upload';
import Api from '@/api/customer/online-order';
import { getPrivateFileUrlApi } from '@/api/enterprise/sensitive-word';

function getFileType(fileName: string) {
  const fileType = fileName.split('.').pop();
  const accept = {
    word: ['doc', 'docx'], 
    excel: ['xls', 'xlsx'], 
    ppt: ['ppt', 'pptx'], 
    pdf: ['pdf'], 
    img: ['jpeg','jpg', 'png'], 
  }
  for (const key in accept) {
    if (accept[key].includes(fileType.toLowerCase())) {
      return key
    }
  }
  return null
}

@Component
export default class BatchImport extends Vue {
  @Model('update:visible', { type: Boolean, default: false }) visible!: boolean;
  @Prop({ type: String, default: '批量上传转货单' }) title: string
  @Prop({ type: String, default: '下载待发货订单数据' }) downloadBtnText: string
  @Prop({ type: String, default: '网仓订单表' }) filename: string
  @Prop({ type: String, default: 'online_warehouse/export' }) templateUrl: string
  @Prop({ type: Object, default: () => ({}) }) params: any
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

  storeList:any[] = []
  formData = {
    third_id: '',
    file_url: '',
    transfer_file: {
      name: '',
      type: '',
      url: '',
      size: '',
      _status: '',
      _progress: 0
    }
  }

  formRules = {
    third_id: [{ required: true, message: '请选择门店' }],
    file_url: [{ required: true, message: '请上传模版' }],
    transfer_file: [{ required: true, message: '请上传转货单' }],
  }

  get fileList2() {
    if(!this.formData.transfer_file.url) return []
    return [
      {
        name: this.formData.transfer_file.name,
        url: this.formData.transfer_file.url
      }
    ]
  }

  mounted() {
    this.getDepList()
  }

  async getDepList() {
    const res:any = await Api.getDepListApi()
    this.storeList =res.data
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
    cosUploader.uploadFile(file, {
      folder: 'private/online-order/transfer-import',
      ACL: 'private',
    }).then(res => {
      this.formData.file_url = res.url
      this.fileList = [res.url]
      onSuccess(res)
    }).catch(onError)
  }

  /**
   * 上传成功
   * @param res
   */
  uploadSuccess() {
    
    this.$message.success("上传成功");
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
   * 移除上传
   */
   handleUploadRemove() {
    this.formData.transfer_file = {
      name: '',
      type: '',
      url: '',
      size: '',
      _status: '',
      _progress: 0
    };
  }


  /**
   * 下载模版
   */
  async handleDownload() {
    if(!this.formData.third_id) {
      this.$message.error("请选择门店");
      return
    }
    const res =await Api.getExportListApi({
      third_id: this.formData.third_id
    })
    // const url = await cosUploader.getObjectUrl(res.data.url)
    downloadFile(res.data.url)
  }

  /**
   * 下载失败数据
   */
  async handleDownloadFail() {
    const res: any = await getPrivateFileUrlApi({ url: this.uploadResult.fail_url })
    downloadFile(res.data.url, {
      filename: `${this.filename}_导入失败数据.xlsx`,
    })
  }

  handleOpen() {
    this.uploadResult = null
    this.fileList = []
    this.formData = {
      third_id: '',
      file_url: '',
      transfer_file: {
        name: '',
        type: '',
        url: '',
        size: '',
        _status: '',
        _progress: 0
      }
    };
    (this.$refs.form as any).resetFields()
  }


  /**
   * 上传文件限制规范
   */
   handleUploadTransferBefore(file: any) {
    // 校验类型 word、excel、ppt、pdf、图片、视频
    if (!getFileType(file.name)) {
      this.$message.error('上传文件类型不支持');
      return false;
    }
    // 校验文件大小不能超过10M
    if (file.size > 10 * 1024 * 1024) {
      this.$message.error('上传文件大小不能超过10M');
      return false;
    }
    
    return true;
  }

  handleExceedTransfer() {
    this.$message.error('上传文件数量不能超过1个');
  }

  /**
   * 上传文件
   */
  handleRequestTransfer({ file, onSuccess, onError }: any) {
    this.formData.transfer_file = {
      name: file.name,
      url: '',
      size: file.size,
      type: getFileType(file.name),
      _progress: 0,
      _status: 'uploading'
    }
    cosUploader.sliceUploadFile(file, { 
      filename: file.name,
      folder: 'private/online-order',
      ACL: 'private',
      onProgress: (progress: any) => {
        this.formData.transfer_file._progress = progress.percent*100;
      }
    }).then((res) => {
      this.formData.transfer_file._status = 'success';
      this.formData.transfer_file.url = res.url;
      onSuccess(res)
    }).catch(() => {
      this.formData.transfer_file._status = 'error';
      onError()
    })
  }

  handleUploadTransferSuccess() {
    console.log(1);
    
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
    (this.$refs.form as any).validate((valid: boolean) => {
      if (valid) {
        Api.importDataApi(this.formData).then((res: any) => {
          if(res.code === 200) {
            this.$message.success("上传成功");
            this.uploadResult = res.data
            if(!res.data.fail_url) {
              this.$emit('update:visible', false)
              this.$emit('confirm')
            }
          } else {
            this.fileList = []
            this.$message.error(res.message);
          }
        })
      }
    })
  }
}
</script>
<style lang="scss" scoped>
.el-form-item {
  margin-bottom: 10px;
}
::v-deep .el-upload-dragger {
  line-height: normal !important;
}
.item {
  display: flex;
  align-items: center;
  line-height: normal;
  .item-label {
    width: 80px;
  }
  .item-value {
    flex: 1;
    padding: 10px 15px;
    background-color: #f6f7f8;
    overflow: hidden;
  }

  .item-name {
    display: flex;
    align-items: center;
    overflow: hidden;

    .icon {
      font-size: 20px;
    }

    .name {
      flex: 1;
      overflow: hidden;
      margin: 0 5px;
      display: flex;
      align-items: center;
      .name-text {
        flex: 1;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
}
::v-deep .el-dialog {
  max-height: 650px !important
}
</style>
