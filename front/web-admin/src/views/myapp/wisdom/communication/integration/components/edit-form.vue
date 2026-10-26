<template>
  <div>
    <el-form ref="form" :model="formInfo" :rules="formRules" label-width="110px">
      <el-form-item label="传讯类型">
        <el-radio-group v-model="formInfo.type">
          <el-radio :label="1">周传讯</el-radio>
          <el-radio :label="2">转货传讯</el-radio>
          <el-radio :label="3">紧急传讯</el-radio>
        </el-radio-group>
        <el-form-item v-if="formInfo.type==1" label-width="0" prop="week_num">
          <span class="m-r-10">第</span>
          <el-input v-model.number="formInfo.week_num" size="medium" style="width: 120px;" placeholder="请输入数字" maxlength="2"></el-input>
          <span class="m-l-10">周</span>
        </el-form-item>
      </el-form-item>
      <el-form-item label="标题" prop="title">
        <cy-input
          size="medium"
          width-size="medium"
          placeholder="请输入标题"
          v-model="formInfo.title"
          maxlength="50"
          show-word-limit
          @blur="handleBlur('title')"
        ></cy-input>
      </el-form-item>
      <el-form-item label="封面图">
        <upload-image
          class="upload-template"
          :urls="formInfo.cover"
          accept="image/jpeg,image/png"
          :limit="1"
          :before-upload="handleUploadBefore"
          @change="handleUploadChange"
        >
        </upload-image>
        <div class="tips">注：建议上传图片尺寸为324x324；或上传1:1的封面图；支持jpg、png格式。</div>
      </el-form-item>
      <el-form-item label="内容" prop="content">
        <content-info v-model="formInfo.content" @change="handleContentChange"></content-info>
      </el-form-item>
      <el-form-item label="通知人群" prop="crowd_list">
        <el-radio-group v-model="formInfo.crowd_type" @change="handleCrowdTypeChange">
          <el-radio :label="1">选择白名单人群</el-radio>
          <el-radio :label="2">选择企微人员</el-radio>
        </el-radio-group>
        <el-row v-if="formInfo.crowd_list.length===0" class="m-b-10">
          <el-button v-if="formInfo.crowd_type===1" type="primary" size="small" @click="handleSelectCrowd">选择白名单人群</el-button>
          <el-button v-if="formInfo.crowd_type===2" type="primary" size="small" @click="handleSelectMember">选择人员</el-button>
        </el-row>
        <el-row v-if="formInfo.crowd_list.length">
          <el-table :data="[{data: formInfo.crowd_list}]">
            <el-table-column v-if="formInfo.crowd_type===1" label="白名单人员群组">
              <template slot-scope="scope">
                <el-tag class="m-r-10 m-b-10" v-for="item in scope.row.data" :key="item.id" size="mini"><i class="el-icon-user-solid"></i>{{item.name}}</el-tag>
              </template>
            </el-table-column>
            <el-table-column v-if="formInfo.crowd_type===2" label="员工">
              <template slot-scope="scope">
                <el-tag class="m-r-10 m-b-10" v-for="item in scope.row.data" :key="item.id" size="mini"><i class="el-icon-user-solid"></i>{{item.name}}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100px">
              <el-button v-if="formInfo.crowd_type===1" type="text" @click="handleSelectCrowd">修改</el-button>
              <el-button v-if="formInfo.crowd_type===2" type="text" @click="handleSelectMember">修改</el-button>
            </el-table-column>
          </el-table>
        </el-row>
      </el-form-item>
      <el-form-item label="发布时间" prop="send_type">
        <el-radio-group v-model="formInfo.send_type" @change="handleSendTypeChange">
          <el-radio :label="1">即刻发布</el-radio>
          <el-radio :label="0">定时发布</el-radio>
          <el-radio :label="2">保存为草稿</el-radio>
        </el-radio-group>
        <el-form-item v-if="formInfo.send_type === 0" label="选择时间" prop="send_time">
          <el-date-picker
            size="medium"
            v-model="formInfo.send_time"
            type="datetime"
            format="yyyy-MM-dd HH:mm"
            placeholder="请选择发布时间"
            value-format="timestamp"
            :default-time="defaultTime"
            :picker-options="pickerOptions"
          >
          </el-date-picker>
        </el-form-item>
      </el-form-item>

      <el-form-item v-if="formInfo.type === 1" label="邮件链接" prop="email_loop_url">
        <cy-input
          type="textarea"
          size="medium"
          width-size="medium"
          placeholder="请输入邮件链接"
          v-model="formInfo.email_loop_url"
        ></cy-input>
      </el-form-item>
    </el-form>
    <!-- 选择白名单 -->
    <select-whitelist 
      ref="selectWhitelist" 
      v-model="formInfo.crowd_list" 
      @change="handleSelectCrowdChange"
    ></select-whitelist>
    <!-- 选择员工 -->
    <!-- <select-member      
      v-if="showSelectMember"
      member-key="member_id"
      :hasSelectArr="formInfo.crowd_list"
      @close="showSelectMember = false"
      @submit="handleSelectMemberConfirm"
    /> -->
    <select-member      
      :show.sync="showSelectMember"
      :value="selectMember"
      @cancel="showSelectMember = false"
      @confirm="handleSelectMemberConfirm"
    />

    <temp-two ref="detail" :value="formInfo"></temp-two>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import UploadImage from "@/components/business/upload-image/index.vue";
import ContentInfo, { IContentItem, defaultItem } from './content-info.vue';
import SelectWhitelist from './select-whitelist.vue';
// 选择员工
import SelectMember from "../../components/select-member.vue";
import TempTwo from './temp-two.vue'
// import SelectMember from "@/components/business/select-member-copy/index.vue";
export interface IFormInfo {
  /** RO传讯ID，存在则更新 */
  id?: number|string;
  /** 传讯类型 */
  type: number;
  /** 第几周 */
  week_num: string|number;
  /** 标题 */
  title: string;
  /** 标题 */
  theme_title?: string;
  /** 图url */
  cover: string;
  /** 发送人群类型 1白名单 2直接选人 */
  crowd_type: number;
  /** 白名单id或者member_id集合 */
  crowd_list: any[];
  /** 发送类型 0定时发布，1保存后发布，2暂不发布 */
  send_type: number;
  /** 定时发布时间 */
  send_time: number;
  /** pdf链接 */
  pdf_url?: string;
  /** 整合传讯内容 */
  content: IContentItem[];
  content_list?: IContentItem[];
  member_name?: string;
  status?: number;
  /** 周传讯下载地址 */
  email_loop_url?: string;
}

export const defaultFormInfo: IFormInfo = {
  type: 1,
  week_num: '',
  title: '',
  cover: '',
  crowd_type: 1,
  crowd_list: [],
  send_type: 1,
  send_time: null,
  content: [
    JSON.parse(JSON.stringify(defaultItem))
  ],
  email_loop_url: ''
}

@Component({
  components: {
    UploadImage,
    ContentInfo,
    SelectWhitelist,
    SelectMember,
    TempTwo
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
    email_loop_url: [{ required: true, message: '请输入周传讯下载地址' }],
    type: [{ required: true, message: '请选择传讯类型' }],
    week_num: [{ required: true, message: '请输入第几周' }],
    title: [{ required: true, message: '请输入标题名称' }],
    send_type: [{ required: true, message: '请选择发布时间' }],
    send_time: [{ 
      required: true, 
      validator: (rule: any, value: any, callback: any) => {
        if (!value) {
          callback(new Error("请选择发布时间"));
        } else if (value<=new Date().getTime()) {
          callback(new Error("发布时间不能小于当前时间"));
        } else {
          callback();
        }
      },
    }],  
    crowd_list: [{
      required: true,
      validator: (rule: any, value: any, callback: any) => {
        if (value.length === 0) {
          callback(new Error("请选择通知人群"));
        } else {
          callback();
        }
      },
    }],
    content: [{
      required: true,
      validator: (rule: any, value: any, callback: any) => {
        if (value.length === 0) {
          callback(new Error("请添加内容"));
        } else {
          callback();
        }
      },
    }]
  }
  /** 选择通知人群弹框 */
  showSelectMember = false

  isAddWhitelist = false

  t = ''

  pickerOptions = {
    disabledDate(time) {
      return time.getTime() < new Date().setHours(0,0,0,0);
    }
  }

  get selectMember() {
    return this.formInfo.crowd_list.map(item => {
      return {
        label: item.name,
        value: item.member_id,
        dp_parent_id: item.dp_parent_id
      }
    })
  }

  get defaultTime() {
    const date = new Date(Date.now()+5*60*1000)
    const hour = date.getHours()
    const minute = date.getMinutes()
    const second = date.getSeconds()

    return `${hour<10?'0'+hour:hour}:${minute<10?'0'+minute:minute}:00`
  }
  
  @Watch('data')
  onChange(value: IFormInfo) {
    if(value) {
      this.formInfo = value
    }
  }

  @Watch('formInfo', { deep: true })
  onValueChange(value: IFormInfo) {
    this.$emit('change', value)
  }

  activated() {
    if(this.t !== this.$route.query.t) {
      this.t = this.$route.query.t as string
      this.formInfo = JSON.parse(JSON.stringify(defaultFormInfo));
      (this.$refs.form as any).resetFields()
    }
  }


  /**
   * 上传前校验
   * @param file 
   */
  async handleUploadBefore(file: File) {
    if (file.size > 2 * 1024 * 1024) {
      this.$message.error("请上传大小不超过2MB的图片");
      throw new Error("请上传大小不超过2MB的图片");
    }
    return true
  }

  /**
   * 上传成功
   * @param e 
   */
  handleUploadChange(e: string) {
    this.formInfo.cover = e
  }

  /**
   * 内容变化
   */
  handleContentChange() {
    (this.$refs.form as any).validateField('content')
  }

  /**
   * 通知人群类型变化
   */
  handleCrowdTypeChange() {
    this.formInfo.crowd_list = [];
    (this.$refs.form as any).clearValidate('crowd_list')
  }
  /**
   * 选择通知人群
   */
  handleSelectCrowd() {
    (this.$refs.selectWhitelist as any).handleShow()
  }

  /**
   * 选择通知人群 确定
   */
  handleSelectCrowdChange() {
    (this.$refs.form as any).validateField('crowd_list')
  }

  /**
   * 选择人员
   */
  handleSelectMember() {
    this.showSelectMember = true
  }

  /**
   * 选择人员 确定
   */
  handleSelectMemberConfirm(value) {
    this.showSelectMember = false
    this.formInfo.crowd_list = value.map(item => {
      return {
        member_id: item.value,
        name: item.label,
      }
    });
    (this.$refs.form as any).validateField('crowd_list')
  }

  /**
   * 发送类型改变
   */
  handleSendTypeChange() {
    this.formInfo.send_time = null
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

  handlePreview() {
    (this.$refs.detail as any).handlePreview()
  }
}
</script>
<style lang="scss" scoped>
.el-select {
  width: 400px;
}
.tips {
  color: #999;
}
</style>