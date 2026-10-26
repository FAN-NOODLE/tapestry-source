<template>
  <div>
    <el-form ref="form" :model="formInfo" :rules="formRules" label-width="110px">
      <el-row class="title">订单信息</el-row>
      <el-row class="section m-b-20">
        <el-col :span="6">
          <el-form-item label="销售单号：">{{ formInfo.order_sn }}</el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="订单状态：">{{ ['', '待填写', '待发货', '已发货', '已拒单', '已作废'][formInfo.status] }}</el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item v-if="formInfo.receiver_phone" label="收件人信息：">
            <span>{{ formInfo.receiver_name }}，</span>
            <span>{{ formInfo.receiver_phone }}，</span>
            <span>{{ formInfo.receiver_address }}</span>
          </el-form-item>
          <el-form-item v-else label="收件人信息：">
            <span>待填写</span>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="门店代码：">{{ formInfo.store_code }}</el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="门店名称：">{{ formInfo.store_name }}</el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="创建人：">{{ formInfo.member_name }}</el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="创建时间：">{{ formInfo.created_at | formatDate('yyyy/mm/dd hh:mm:ss') }}</el-form-item>
        </el-col>
        <el-col v-if="formInfo.status == 3" :span="6">
          <el-form-item label="发货时间：">{{ formInfo.shipped_at | formatDate('yyyy/mm/dd hh:mm:ss') }}</el-form-item>
        </el-col>
        <el-col v-if="formInfo.status == 4" :span="6">
          <el-form-item label="拒单时间：">{{ formInfo.operate_time | formatDate('yyyy/mm/dd hh:mm:ss') }}</el-form-item>
        </el-col>
      </el-row>
      <el-row class="title">商品信息</el-row>
      <el-row class="m-b-10">
        <el-col :span="2">商品数量：{{quantity}}</el-col>
        <el-col :span="2">拒单数量：{{refusalQuantity}}</el-col>
      </el-row>
      <el-table class="m-b-20" :data="formInfo.sku_list">
        <el-table-column prop="image_url" label="SKU图片">
          <template slot-scope="scope">
            <el-image :src="scope.row.image_url || require('@/assets/images/icon_empty.png')" :preview-src-list="scope.row.image_url ? [scope.row.image_url]:[]" style="width: 50px;height: 50px;">
              <div slot="error" class="image-slot">
                <i class="el-icon-picture-outline"></i>
              </div>
            </el-image>
          </template>
        </el-table-column>
        <el-table-column prop="sku_info" label="SKU"></el-table-column>
        <el-table-column prop="quantity" label="商品数量"></el-table-column>
        <el-table-column prop="refusal_quantity" label="拒单数量">
          <template slot-scope="scope">
            <span v-if="readonly">{{ scope.row.refusal_quantity }}</span>
            <el-input-number v-else v-model="scope.row.refusal_quantity" :step="1" :min="0" :max="scope.row.quantity" @change="() => scope.row.refusal_quantity = scope.row.refusal_quantity||0" />
          </template>
        </el-table-column>
      </el-table>
      <el-row class="title">物流信息<slot></slot></el-row>
      <el-row v-if="readonly && formInfo.status !==3">
        <span style="padding: 30px;">暂无物流信息</span>
      </el-row>
      <div v-else>
        <el-form-item v-if="quantity !== refusalQuantity" label="顺丰单号" prop="tracking_number">
          <el-input v-model="formInfo.tracking_number" :disabled="readonly && !isEdit" placeholder="请输入顺丰单号" />
        </el-form-item>
        <el-form-item v-if="quantity !== refusalQuantity" label="转货单" prop="transfer_file.url" :rules="formRules.transfer_file">
          <el-upload 
            v-if="!(readonly && !isEdit)"
            :multiple="false"
            :limit="1" 
            :show-file-list="false"
            :file-list="fileList"
            :before-upload="handleUploadBefore"
            :http-request="handleRequest"
            :on-exceed="handleExceed">
            <el-button type="primary">上传转货单</el-button>
            <span slot="tip" class="el-upload__tip m-l-10">*可支持上传word、excel、ppt、pdf、图片</span>
          </el-upload>
          <el-row v-if="formInfo.transfer_file && formInfo.transfer_file.name" class="item m-t-20">
            <div class="item-value">
              <div class="item-name">
                <svg-icon class="icon" :icon-class="formInfo.transfer_file.type"></svg-icon>
                <div class="name">
                  <div class="name-text">{{formInfo.transfer_file.name}}</div>
                  <el-tag v-if="formInfo.transfer_file._status === 'error'" class="m-l-10" type="danger" size="mini">上传失败</el-tag>
                </div>
                <el-button v-if="formInfo.transfer_file._status !== 'uploading' && formInfo.transfer_file._status !== 'error'" type="text" @click="handleDownLoad(formInfo.transfer_file)" icon="el-icon-download"></el-button>
                <el-button v-if="!(readonly && !isEdit) && formInfo.transfer_file._status !== 'uploading'" type="text" @click="handleUploadRemove()" icon="el-icon-close"></el-button>
              </div>
              <el-progress v-if="formInfo.transfer_file._status === 'uploading'" color="#000" :percentage="formInfo.transfer_file._progress" :status="formInfo.transfer_file._status" :show-text="false" />
            </div>
          </el-row>
        </el-form-item>
      </div>
    </el-form>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import { cosUploader } from '@/utils/upload';
import { downloadFile } from '@/utils/download';

export interface IFormInfo {
    id?: number;
    all_num: number;
    created_at: number;
    member_id: number;
    member_name: string;
    order_sn: string;
    receiver_address: string;
    receiver_name: string;
    receiver_phone: string;
    refusal_num: number;
    shipped_at: number;
    operate_time: number;
    status: number;
    sku_list: Array<{
      name: string;
      quantity: number;
      refusal_quantity: number;
      sku_info: string;
      status: number;
      upc_code: string;
    }>;
    store_code: string;
    store_name: string;
    tracking_number: string;
    transfer_file: {
      /** 名称 */
      name: string;
      /** 链接 */
      url: string;
      /** 类型 */
      type: string;
      /** 大小 */
      size: number;
      /** 进度 */
      _progress?: number;
      /** 状态 */
      _status?: string
    };
    updated_at: number;
  }

export const defaultFormInfo: IFormInfo = {
  all_num: 0,
  created_at: 0,
  member_id: 0,
  member_name: '',
  order_sn: '',
  receiver_address: '',
  receiver_name: '',
  receiver_phone: '',
  refusal_num: 0,
  shipped_at: 0,
  operate_time: 0,
  sku_list: [],
  store_code: '',
  store_name: '',
  tracking_number: '',
  transfer_file: {
    name: '',
    url: '',
    type: '',
    size: 0
  },
  updated_at: 0,
  status: 0
}

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

@Component({
  components: {
  }
})
export default class EditForm extends Vue {
  /** 是否只读 */
  @Prop({ type: Boolean, default: false }) readonly: boolean;
  /** 编辑 */
  @Prop({ type: Boolean, default: false }) isEdit: boolean;
  /** 表单数据 */
  @Prop({ type: Object, default: () => defaultFormInfo }) data: IFormInfo;
  /** 表单数据 */
  formInfo: IFormInfo = JSON.parse(JSON.stringify(defaultFormInfo))
  /** 表单校验规则 */
  formRules = {
    tracking_number: [{ required: true, message: '请输入顺丰单号' }],
    transfer_file: [{ required: true, message: '请上传转货单' }],
  }
  

  @Watch('data')
  onChange(value: IFormInfo) {
    if(value) {
      this.formInfo = value
    }
  }
  
  get fileList() {
    if(!this.formInfo.transfer_file.url) return []
    return [
      {
        name: this.formInfo.transfer_file.name,
        url: this.formInfo.transfer_file.url
      }
    ]
  }
  /** 商品数量 */
  get quantity() {
    return this.formInfo.sku_list.reduce((total, item) => total + item.quantity, 0)
  }
  /** 拒单数量 */
  get refusalQuantity() {
    return this.formInfo.sku_list.reduce((total, item) => total + item.refusal_quantity, 0)
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
    
    return true;
  }

  handleExceed() {
    this.$message.error('上传文件数量不能超过1个');
  }

  /**
   * 上传文件
   */
  handleRequest({ file }: any) {
    this.formInfo.transfer_file = {
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
        this.formInfo.transfer_file._progress = progress.percent*100;
      }
    }).then((res) => {
      this.formInfo.transfer_file._status = 'success';
      this.formInfo.transfer_file.url = res.url;
    }).catch(() => {
      this.formInfo.transfer_file._status = 'error';
    })
  }
  /**
   * 移除上传
   */
  handleUploadRemove() {
    this.formInfo.transfer_file = {
      name: '',
      url: '',
      type: '',
      size: 0
    };
  }

  /**
   * 下载文件
   */
  handleDownLoad(item: any) {
    cosUploader.getObjectUrl(item.url).then(res => {
      downloadFile(res, {
        filename: item.name
      })
    })
  }

  /**
   * 表单校验
   */
  validate(callback: (data: IFormInfo) => void) {
    (this.$refs.form as any).validate((valid: boolean) => {
      if (valid) {
        callback(this.formInfo)
      }
    })
  }

  handleBlur(key: string) {
    this.formInfo[key] = this.formInfo[key].trim()
  }

}
</script>
<style lang="scss" scoped>

::v-deep .image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: #f5f7fa;
  color: #909399;
  font-size: 30px;
}

.title {
  margin-bottom: 20px;
  font-weight: bold;
  font-size: 16px;
}

.section {
  padding: 10px 10px 0;
  background-color: #f6f7f8;

  .el-form-item {
    margin-bottom: 10px;
  }
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