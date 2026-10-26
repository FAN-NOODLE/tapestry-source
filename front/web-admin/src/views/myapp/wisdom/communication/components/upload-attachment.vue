<template>
  <el-row class="upload-attachment">
    <el-row class="m-b-10">
      <el-upload 
        multiple
        :limit="30" 
        :show-file-list="false"
        :file-list="fileList"
        :before-upload="handleUploadBefore"
        :http-request="handleRequest"
        :on-exceed="handleExceed">
        <el-button type="primary">+添加附件</el-button>
        <span slot="tip" class="el-upload__tip m-l-10">*可支持上传word、excel、ppt、pdf、图片、视频（单个附件内容 ≤ 10M ;总附件数量≤30个）</span>
      </el-upload>
    </el-row>
    <!-- 文件列表 -->
    <el-row class="item" v-for="(item, index) in value" :key="item.name">
      <div class="item-label">{{item.name.replace(/\d+/, '')}}{{index+1}}</div>
      <div class="item-value">
        <div class="item-name">
          <svg-icon class="icon" :icon-class="item.type"></svg-icon>
          <div class="name">
            <div class="name-text">{{item.attachment_name}}</div>
            <el-tag v-if="item._status === 'error'" class="m-l-10" type="danger" size="mini">上传失败</el-tag>
          </div>
          <el-button v-if="item._status !== 'uploading'" type="text" @click="handleDownLoad(item)">下载</el-button>
          <el-button v-if="item._status !== 'uploading'" type="text" @click="handleUploadRemove(index)">删除</el-button>
        </div>
        <el-progress color="#000" :percentage="item._progress" :status="item._status" :show-text="false" />
      </div>
    </el-row>
  </el-row>
</template>
<script lang="ts">
import { Vue, Component, Prop } from 'vue-property-decorator'
import {cosUploader} from '@/utils/upload';
import { downloadFile } from '@/utils/download';

export interface IAttachmentItem {
  /** 附件label名称 */
  name: string;
  /** 附件名称 */
  attachment_name: string;
  /** 附件url */
  url: string;
  /** 附件类型 */
  type: string;
  /** 附件大小 */
  size: number;
  /** 上传进度 */
  _progress: number;
  /** 上传状态 */
  _status: string;
}

function getFileType(fileName: string) {
  const fileType = fileName.split('.').pop();
  const accept = {
    word: ['doc', 'docx'], 
    excel: ['xls', 'xlsx'], 
    ppt: ['ppt', 'pptx'], 
    pdf: ['pdf'], 
    img: ['jpeg','jpg', 'png'], 
    video: ['mp4']
  }
  for (const key in accept) {
    if (accept[key].includes(fileType.toLowerCase())) {
      return key
    }
  }
  return null
}

@Component({
  model: {
    prop: 'value',
    event: 'change'
  }
})
export default class UploadAttachment extends Vue {
  /** 附件列表 */
  @Prop({ type: Array, default: () => [] }) value!: IAttachmentItem[];

  get fileList() {
    return this.value.map((item: IAttachmentItem) => {
      return {
        name: item.attachment_name,
        url: item.url
      }
    })
  }
  
  /**
   * 上传文件限制规范
   */
  handleUploadBefore(file: any) {
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

    // 附件数量限制
    if (this.value.length >= 30) {
      this.$message.error('上传文件数量不能超过30个');
      return false;
    }
    
    return true;
  }

  handleExceed() {
    this.$message.error('上传文件数量不能超过30个');
  }

  /**
   * 上传文件
   */
  handleRequest({ file }: any) {
    console.log(file);
    console.log(this.value);
    const item: IAttachmentItem = {
      name: 'Attach',
      attachment_name: file.name,
      url: '',
      size: file.size,
      type: getFileType(file.name) as string,
      _progress: 0,
      _status: 'uploading'
    }

    this.value.push(item);
    this.$emit('change', this.value);
    cosUploader.sliceUploadFile(file, { 
      filename: file.name,
      folder: 'private/attachment',
      ACL: 'private',
      onProgress: (progress: any) => {
        item._progress = progress.percent*100;
      }
    }).then((res) => {
      item._status = 'success';
      item.url = res.url;
    }).catch(() => {
      item._status = 'error';
    })
  }
  /**
   * 移除上传
   */
  handleUploadRemove(index: number) {
    this.value.splice(index, 1)
    this.$emit('change', this.value);
  }

  /**
   * 下载附件
   */
  handleDownLoad(item: IAttachmentItem) {
    
    cosUploader.getObjectUrl(item.url).then(res => {
      downloadFile(res, {
        filename: item.attachment_name
      })
      // window.open(res)
    })
  }

}
</script>

<style lang="scss" scoped>
.upload-attachment {
  width: 80%;
}
.item {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
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
</style>