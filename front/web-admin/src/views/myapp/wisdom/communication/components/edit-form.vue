<template>
  <div>
    <el-form ref="form" :model="formInfo" :rules="formRules" label-width="110px">
      <el-form-item label="传讯类型" prop="type">
        <cy-select
          v-model="formInfo.type"
          size="medium"
          width-size="medium"
          placeholder="请选择传讯类型"
          :optionsList="typeOptions"
        ></cy-select>
      </el-form-item>
      <el-form-item label="所属部门" prop="department_name">
        
        <el-select
          v-model="formInfo.department_name"
          size="medium"
          placeholder="请选择所属部门"
          filterable
          style="width: 400px;"
        >
          <el-option v-for="(item, index) in departmentOptions" :key="index" :label="item.label" :value="item.value"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item v-if="formInfo.type===3" label="主题" prop="theme_title">
        <cy-input
          size="medium"
          width-size="medium"
          placeholder="请输入主题"
          v-model="formInfo.theme_title"
          maxlength="100"
          show-word-limit
          :disabled="readonly"
          @blur="handleBlur('theme_title')"
        ></cy-input>
      </el-form-item>
      <el-form-item label="标题" prop="title">
        <cy-input
          size="medium"
          width-size="medium"
          placeholder="请输入标题"
          v-model="formInfo.title"
          maxlength="100"
          show-word-limit
          :disabled="readonly"
          @blur="handleBlur('title')"
        ></cy-input>
      </el-form-item>
      <el-form-item label="至" prop="until_name">
        <cy-input
          size="medium"
          width-size="medium"
          placeholder="请输入通知人群"
          v-model="formInfo.until_name"
          maxlength="2000"
          :disabled="readonly"
          @blur="handleBlur('until_name')"
        ></cy-input>
      </el-form-item>
      <el-form-item label="内容" prop="content">
        <rich-text-editor editorId="content" v-model="formInfo.content" @change="handleContentChange"></rich-text-editor>
        <div class="tip">*为不影响员工手机端查看效果，插入复杂表格效果建议在Excle表格中调整好格式后截图，通过插入图片放置，或放置在附件内容中</div>
      </el-form-item>
      <el-form-item label="执行事务" prop="perform">
        <rich-text-editor editorId="perform" v-model="formInfo.perform"></rich-text-editor>
        <div class="tip">*为不影响员工手机端查看效果，插入复杂表格效果建议在Excle表格中调整好格式后截图，通过插入图片放置，或放置在附件内容中</div>
      </el-form-item>
      <el-form-item label="日期" prop="date_time">
        <cy-input
          size="medium"
          width-size="medium"
          placeholder="请输入日期"
          v-model="formInfo.date_time"
          maxlength="100"
          show-word-limit
          :disabled="readonly"
          @blur="handleBlur('date_time')"
        ></cy-input>
      </el-form-item>
      <el-form-item label="联系人" prop="contact_info">
        <contact-info v-model="formInfo.contact_info" @change="handleContactChange"></contact-info>
      </el-form-item>
      <el-form-item label="附件">
        <upload-attachment v-model="formInfo.attachment"></upload-attachment>
      </el-form-item>
      <el-form-item v-if="formInfo.type === 1" label="页脚设置">
        <el-radio-group v-model="formInfo.page_url">
          <el-radio label="1">显示</el-radio>
          <el-radio label="2">隐藏</el-radio>
        </el-radio-group>
        <div v-if="formInfo.page_url === '1'" class="page-footer">
          <div class="page-footer__item">
            <span class="page-footer__item__label">负责跟进同事：</span>
            <span class="page-footer__item__value">负责跟进同事</span>
          </div>
          <div class="page-footer__item">
            <span class="page-footer__item__label">执行行动日期：</span>
            <span class="page-footer__item__value">负责跟进同事</span>
          </div>
          <div class="page-footer__item">
            <span class="page-footer__item__label">完成行动日期：</span>
            <span class="page-footer__item__value">负责跟进同事</span>
          </div>
        </div>
      </el-form-item>
    </el-form>
    <temp-one ref="detail" :data="formInfo"></temp-one>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import { typeOptions } from '../const'
import UploadAttachment, { IAttachmentItem } from "./upload-attachment.vue";
import RichTextEditor from "./rich-text-editor.vue";
import ContactInfo, { IContactItem } from "./contact-info.vue";
import Api from '@/api/wisdom/communication';
import TempOne from './temp-one.vue'


export interface IFormInfo {
  id?: number|string;
  /** 传讯类型 */
  type: number|string;
  /** 传讯类型名称 */
  type_name?: string;
  /** 部门id */
  department_name: string;
  /** 标题 */
  title: string;
  /** 主题 */
  theme_title: string;
  /** 至 */
  until_name: string;
  /** 内容 */
  content: string;
  /** 执行事务 */
  perform: string;
  /** 时间 */
  date_time: string;
  /** 联系人 */
  contact_info: IContactItem[];
  /** 附件 */
  attachment: IAttachmentItem[];
  /** 审核人 */
  review_name?: string;
  /** 审核时间 */
  review_time?: number;
  /** 审核留言 */
  comment?: string;
  /** 状态 */
  status?: number;
  /** 状态名称 */
  status_name?: string;
  /** 页脚设置 */
  page_url?: string;
  /** 导购ID */
  member_id?: number;
  /** 导购名称 */
  member_name?: string;
  /** 提交时间 */
  submit_time?: number;
  /** pdf 地址 */
  pdf_url?: string
}

export const defaultFormInfo: IFormInfo = {
  type: '',
  department_name: '',
  title: '',
  theme_title: '',
  until_name: '',
  content: '',
  perform: '',
  date_time: '',
  contact_info: [{
    department_name: '',
    name: '',
    email: '',
    tel_phone: '',
    remark: ''
  }],
  attachment: [],
  page_url: '1'
}

@Component({
  components: {
    UploadAttachment,
    RichTextEditor,
    ContactInfo,
    TempOne
  }
})
export default class EditForm extends Vue {
  /** 是否只读 */
  @Prop({ type: Boolean, default: false }) readonly: boolean;
  /** 表单数据 */
  @Prop({ type: Object, default: () => defaultFormInfo }) data: IFormInfo;
  /** 表单数据 */
  formInfo: IFormInfo = JSON.parse(JSON.stringify(defaultFormInfo))
  /** 表单校验规则 */
  formRules = {
    type: [{ required: true, message: '请选择传讯类型' }],
    department_name: [{ required: true, message: '请选择所属部门' }],
    title: [{ required: true, message: '请输入标题名称' }],
    theme_title: [{ required: true, message: '请输入主题' }],
    until_name: [{ required: true, message: '请输入通知人群' }],
    content: [{ required: true, message: '请输入通知内容' }],
    contact_info: [{
      required: true,
      validator: (rule: any, value: any, callback: any) => {
        const isValid = value.every((item: IContactItem) => {
          return Object.values(item).every(val => !val)
        })
        if (value.length === 0 || isValid) {
          callback(new Error("请添加联系人"));
        } else {
          callback();
        }
      },
    }]
  }
  /** 传讯类型选项 */
  typeOptions = typeOptions
  /** 所属部门选项 */
  departmentOptions = []
  /** 上传配置 */
  uploadProps = {
    accept: ".xlsx,.xls,.csv"
  }
  

  @Watch('data')
  onChange(value: IFormInfo) {
    if(value) {
      this.formInfo = value
    }
  }

  created() {
    this.getDepartmentList()
  }

  async getDepartmentList() {
    const res: any = await Api.getDepartmentListApi({})
    if(res.code === 200) {
      this.departmentOptions = res.data.map((item: any) => {
        return {
          label: item.name,
          value: item.name
        }
      })
    } else {
      this.$message.error(res.message)
    }
  }

  /**
   * 内容变化
   */
  handleContentChange() {
    (this.$refs.form as any).validateField('content')
  }

  /**
   * 联系人变化
   */
  handleContactChange() {
    
    (this.$refs.form as any).validateField('contact_info')
  }
  
  /**
   * 表单校验
   */
  validate(callback: (data: IFormInfo) => void) {
    (this.$refs.form as any).validate((valid: boolean) => {
      if (valid) {
        // 判断是否有正在上传的附件
        if (this.formInfo.attachment.some(item => item._status === 'uploading')) {
          this.$message.error('附件正在上传中，请等待上传完成')
          return
        }
        callback(this.formInfo)
      }
    })
  }

  handleBlur(key: string) {
    this.formInfo[key] = this.formInfo[key].trim()
  }

  handlePreview() {
    (this.$refs.detail as any).handlePreview()
  }
  handleCreate() {
    return (this.$refs.detail as any).handleCreate()
  }
}
</script>
<style lang="scss" scoped>
::v-deep {
  .show-word-limit .el-input__inner {
    padding-right: 60px !important;
  }
}
.el-select {
  width: 400px;
}
.page-footer {
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  border-bottom: 2px dotted #000;

  .page-footer__item__value {
    color: transparent;
    text-decoration: underline #000;
    user-select: none;
  }
}

.tip {
  font-size: 12px;
  color: #606266;
}
</style>