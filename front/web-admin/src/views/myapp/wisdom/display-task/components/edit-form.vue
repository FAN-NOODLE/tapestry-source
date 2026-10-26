<template>
  <div>
    <el-form
      ref="formRef"
      label-width="110px"
      :rules="formRules"
      :model="formInfo"
    >
      <el-form-item label="标题" prop="display_name">
        <cy-input
          size="medium"
          width-size="medium"
          placeholder="请输入标题名称"
          v-model="formInfo.display_name"
          maxlength="20"
          show-word-limit
          :readonly="readonly"
        ></cy-input>
      </el-form-item>
      <el-form-item label="执行时间" prop="start_time" :rules="formRules.start_time">
        <el-date-picker
          :readonly="readonly"
          class="cy-el-date-picker"
          v-model="executeTime"
          type="daterange"
          store_type-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          value-format="timestamp"
          :picker-options="startTimeOptSetting"
          :default-time="['00:00:00', '23:59:59']"
        >
        </el-date-picker>
      </el-form-item>
      <el-form-item label="通知门店" prop="store_type">
        <el-radio-group 
          v-model="formInfo.store_type"
          :disabled="readonly">
          <el-radio :label="0">全部门店</el-radio>
          <el-radio :label="1">指定门店</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item v-if="formInfo.store_type === 1" label="" prop="store_list">
        <div v-if="!readonly">
          <el-button v-if="formInfo.store_list.length === 0" type="primary" @click="handleSelectStore">选择门店</el-button>
          <el-table v-if="formInfo.store_list.length>0" :data="[{data: formInfo.store_list}]">
            <el-table-column label="可见门店/员工" min-width="220">
              <template slot-scope="{ row }">
                <div id="list" class="list" :class="{'expand': isExpand}">
                  <el-tag class="tag-item m-r-10 m-b-10" v-for="(item, index) in row.data" :key="index" type="info">
                    <store-item :item="item"></store-item>
                  </el-tag>
                </div>
              </template>
            </el-table-column>
            <el-table-column v-if="!readonly" label="操作" width="140">
              <template>
                <el-button type="text" @click="handleSelectStore">修改</el-button>
                <el-button v-if="showExpandBtn" type="text" @click="handleExpand">{{isExpand?'收起':'展开'}}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <div v-else>
          <div class="box">
            <div id="list" class="list" :class="{'expand': isExpand}">
              <el-tag class="tag-item m-r-10 m-b-10" v-for="(item, index) in formInfo.store_list" :key="index" type="info">
                <store-item :item="item"></store-item>
              </el-tag>
            </div>
            <div class="flex">
              <el-button class="more" v-if="showExpandBtn" type="text"  @click="handleExpand">{{isExpand?'收起':'展开'}}
                <i
                  aria-hidden="true"
                  class="el-icon--right"
                  :class="
                    isExpand ? 'el-icon-arrow-up' : 'el-icon-arrow-down'
                  "
                ></i>
              </el-button>
            </div>
          </div>
        </div>
      </el-form-item>
      <el-form-item label="任务内容" prop="content">
        <el-input
          size="medium"
          v-model="formInfo.content"
          :rows="7"
          type="textarea"
          maxlength="10000"
          placeholder="请输入任务内容"
          :readonly="readonly"
        >
        </el-input>
      </el-form-item>
      <el-form-item label="关联课程链接" prop="class_url">
        <el-row class="item" v-for="(item, index) in formInfo.class_url" :key="index">
          <el-col :span="20">
            <el-input
              v-model="item.class_name"
              size="medium"
              maxlength="30"
              show-word-limit
              placeholder="请输入标题及提示语"
              :readonly="readonly"
            >
            </el-input>
            <el-input
              :readonly="readonly"
              v-model="item.url"
              size="medium"
              placeholder="请输入跳转welearning课程链接"
            >
            </el-input>
          </el-col>
          <el-col v-if="!readonly" :span="4">
            <el-button v-if="index===0" class="btn-add" type="text" icon="el-icon-circle-plus" @click="handleAddUrl"></el-button>
            <el-button v-else class="btn-add" type="text" icon="el-icon-remove" @click="handleDelUrl(index)"></el-button>
          </el-col>
        </el-row>
      </el-form-item>
      <el-form-item label="上传模版字段" prop="template_info">
        <el-row class="item-temp" v-for="(item, index) in formInfo.template_info" :key="index">
          <el-input
            v-model="item.template_name"
            size="medium"
            class="input-with-select"
            maxlength="20"
            show-word-limit
            placeholder="请输入字段名称"
            :readonly="readonly"
          >
          </el-input>
          <el-button v-if="!readonly && index===0" class="btn-add" type="text" icon="el-icon-circle-plus" @click="handleAddTemp"></el-button>
          <el-button v-if="!readonly && index!==0" class="btn-add" type="text" icon="el-icon-remove" @click="handleDelTemp(index)"></el-button>
          <el-checkbox v-model="item.is_required" :disabled="readonly">必传模块</el-checkbox>
        </el-row>
      </el-form-item>
      <el-form-item label="审核权分配" prop="review_type">
        <el-radio-group 
          v-model="formInfo.review_type"
          :disabled="readonly">
          <el-radio :label="0">不分配</el-radio>
          <el-radio :label="1">
            <span>随机平均分配</span>
            <el-button v-if="formInfo.review_type==1 && (formInfo.status>0 || readonly)" class="m-l-10" size="mini" type="text" @click="reviewExport">下载分配详情</el-button>
          </el-radio>
          <el-radio :label="2">
            <span>定向分配</span>
            <div class="absolute-right">
              <el-form-item prop="review_file" label-width="0" v-if="formInfo.review_type==2 && !formInfo.status && !readonly">
                <el-button class="m-l-10" size="mini" type="primary" @click="handleBatchImport">导入分配文件</el-button>
              </el-form-item>
              <el-form-item v-if="formInfo.review_type==2 && (formInfo.status>0 || readonly)" label-width="0">
                <el-button class="m-l-10" size="mini" type="text" @click="reviewExport">下载分配详情</el-button>
              </el-form-item>
            </div>
          </el-radio>
        </el-radio-group>
        <div v-if="formInfo.review_type==0" class="tips">不增加VM Specialist 人员的预审核流程，直接由VM Corporate人员审核。</div>
        <div v-if="formInfo.review_type==1" class="tips">增加VM Specialist 人员的预审核流程，并给VM Specialist人员，随机平均分配门店提交的内容。</div>
        <div v-if="formInfo.review_type==2" class="tips">增加VM Specialist 人员的预审核流程，可根据上传的分配文件，分配VM Specialist需要审核的门店。</div>
      </el-form-item>
    </el-form>
    <!-- 选择门店 -->
    <!-- <select-person
      title="添加门店（支持多选）"
      :level="1"
      :selectedData="formInfo.store_list"
      :visible.sync="selectPersonVisible"
      @submit="handlePersonSubmit"
      :appendToBody="true"
    ></select-person> -->
    <select-store
      title="添加门店（支持多选）"
      v-if="selectPersonVisible"
      :selectArr="formInfo.store_list"
      @submit="handlePersonSubmit"
      @close="selectPersonVisible = false"
      row-key="dp_third_id"
    ></select-store>
    <batch-import
      v-model="showBatchImport"
      templateUrl="https://cdn-cos.coach.com.cn/2025/11/26/%E5%AE%9A%E5%90%91%E5%88%86%E9%85%8D%E6%A8%A1%E6%9D%BF.xlsx"
      :defaultValue="formInfo.review_file"
      @confirm="handleBatchImportConfirm"
    />
  </div>
</template>
<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator';
import SelectPerson from "@/components/business/select-person/index.vue";
import Person from "@/components/business/select-person/person.vue";
import { dateFormat } from '@/utils/tools';
import SelectStore from "@/components/business/select-store/index.vue";
import StoreItem from "@/components/business/select-store/store-item.vue";
import BatchImport from './batch-import.vue';
import { cosUploader } from '@/utils/upload';
import { downloadFile } from '@/utils/download';
import Api from '@/api/wisdom/displayTask';
import dayjs from 'dayjs';
interface IClassUrlItem {
  class_name: string;
  url: string;
}
interface ITemplateItem {
  template_name: string;
  is_required: boolean;
}
export interface IFormInfo {
  id?: number;
  /** 任务名称 */
  display_name: string;
  /** 通知门店类型，0全部门店,1选择门店 */
  store_type: number;
  /** 任务内容及描述 */
  content: string;
  /** 任务的课程学习 */
  class_url: IClassUrlItem[];
  /** 任务的模板内容设置 */
  template_info: ITemplateItem[];
  /** 上下架，0为保存，1为保存并上架 */
  show_type: number;
  /** 活动开始时间 */
  start_time: number;
  /** 活动结束时间 */
  end_time: number;
  /** 选择门店的门店数据 */
  store_list: string[];
  /** 审核权分配类型  */
  review_type: number;
  /** 审核权分配文件 */
  review_file: string
}
export const defaultClassUrlItem: IClassUrlItem = {
  class_name: '',
  url: ''
}
const defaultTemplateItem: ITemplateItem = {
  template_name: '',
  is_required: false
}

export const defaultFormInfo: IFormInfo = {
  display_name: '',
  store_type: 0,
  content: '',
  class_url: [
    {...defaultClassUrlItem}
  ],
  template_info: [
    {...defaultTemplateItem}
  ],
  show_type: 0,
  start_time: null,
  end_time: null,
  store_list: [],
  review_file: '',
  review_type: 0
}
@Component({
  components: {
    SelectPerson,
    SelectStore,
    Person,
    StoreItem,
    BatchImport
  }
})
export default class DisplayTaskEditForm extends Vue {
  @Prop({default: false, type: Boolean}) readonly: boolean;
  @Prop({default: ()=>defaultFormInfo, type: Object}) data: IFormInfo
  /** 显示展开收起按钮 */
  showExpandBtn = false
  /** 批量导入 */
  showBatchImport = false
  /** 展开收起按钮状态 */
  isExpand = false
  /** 选择门店弹框 */
  selectPersonVisible = false
  /** 表单数据 */
  formInfo: IFormInfo = JSON.parse(JSON.stringify(defaultFormInfo))
  /** 表单校验规则 */
  formRules = {
    display_name: [
      { required: true, message: '请输入任务名称', trigger: 'blur' },
      {
        validator: (rule: any, value: string, callback: Function) => {
          if(value.trim() === '') {
            callback(new Error('请输入任务名称'))
          } else {
            callback()
          }
        } 
      }
    ],
    store_type: [
      { required: true, message: '请选择通知门店类型', trigger: 'blur' }
    ],
    content: [
      { required: true, message: '请输入任务内容及描述', trigger: 'blur' },
      {
        validator: (rule: any, value: string, callback: Function) => {
          if(value.trim() === '') {
            callback(new Error('请输入任务内容及描述'))
          } else {
            callback()
          }
        },
      }
    ],
    template_info: [
      { 
        required: true, 
        message: '请完善上传模版字段信息', 
        validator: (rule: any, value: ITemplateItem[], callback: Function) => {
          const len = value.filter(item => item.template_name.trim()).length
          if (len === 0) {
            callback(new Error('请完善上传模版字段信息'))
          } else {
            callback()
          }
        } 
      }
    ],
    show_type: [
      { required: true, message: '请选择上下架', trigger: 'blur' }
    ],
    start_time: [
      { 
        required: true, 
        message: '请选择活动开始时间', 
        trigger: 'change' 
      },
      {
        validator: (rule: any, value: number, callback: Function) => {
          if(Date.now() > value * 1000) {
            callback(new Error('开始时间不能小于当前时间'))
          } else {
            callback()
          }
        },
        trigger: 'blur'
      }
    ],
    store_list: [
      { required: true, message: '请选择通知门店', trigger: 'blur' }
    ],
    review_type: [
      { required: true, message: '请选择审核权分配类型', trigger: 'blur' }
    ],
    review_file: [
      { required: true, message: '请上传审核权分配文件', trigger: 'change' }
    ]
  }
  startTimeOptSetting = {
    disabledDate(time: Date) {
      const str: string = dateFormat((Date.now()) / 1000, "yyyy-MM-dd");
      const now: number = new Date(str + " 23:59:59").getTime();
      return time.getTime() < now;
    }
  };
  /** 执行时间 读取 */
  get executeTime() {
    const { start_time, end_time } = this.formInfo
    return [start_time*1000, end_time*1000].filter(Boolean)
  }
  /** 执行时间 写入 */
  set executeTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.formInfo.start_time = parseInt(start_time/1000+'');
    this.formInfo.end_time = parseInt(end_time/1000+'');
  }
  @Watch('data')
  onChange(value: IFormInfo) {
    this.formInfo = value
  }
  @Watch('formInfo.store_list')
  storeListChange() {
    setTimeout(() => {
      const len = this.formInfo.store_list.length
    this.$nextTick(() => {
      const container = document.getElementById("list");
      const content = document.getElementsByClassName("tag-item");
      const last: HTMLElement = content[len - 1] as HTMLElement;
      if(last) {
        this.showExpandBtn = last.offsetTop > container.clientHeight
      }
      console.log(this.showExpandBtn);
      
    })
    }, 10);
  }

  /**
   * @description 选择门店
   */
  handleSelectStore() {
    this.selectPersonVisible = true
  }

  /**
   * @description 选择门店
   * @param e 
   */
  handlePersonSubmit(e) {
    this.formInfo.store_list = e.map(item => {
      item.mem_subscribe_status = item.mem_subscribe_status ? 1 : 0
      return item
    })
    this.selectPersonVisible = false
  }

  /**
   * @description 展开收起
   */
  handleExpand() {
    this.isExpand = !this.isExpand
  }

  /**
   * @description 删除url
   * @param index 
   */
  handleDelUrl(index: number) {
    this.formInfo.class_url.splice(index, 1)
  }
  
  /**
   * @description 添加url
   */
  handleAddUrl() {
    if(this.formInfo.class_url.length >= 5) {
      this.$message({
        type: 'warning',
        message: '最多可添加5条'
      })
      return
    }
    this.formInfo.class_url.push({...defaultClassUrlItem})
  }

  /**
   * @description 删除url
   * @param index 
   */
  handleDelTemp(index: number) {
    this.formInfo.template_info.splice(index, 1)
  }

  /**
   * @description 添加url
   */
  handleAddTemp() {
    if(this.formInfo.template_info.length >= 10) {
      this.$message({
        type: 'warning',
        message: '最多可添加10条'
      })
      return
    }
    this.formInfo.template_info.push({...defaultTemplateItem})
  }

  /** 
   * @description 校验
   */
  validate(callback: (data: IFormInfo) => void) {
    (this.$refs.formRef as any).validate((valid, values) => {
      console.log(valid, values);
      
      if (valid) {
        callback(this.formInfo)
      } else {
        return false
      }
    })
  }

  handleBatchImport(e) {
    e.stopPropagation()
    this.showBatchImport = true
  }

  /**
   * 模版导入确认
   */
  handleBatchImportConfirm(e) {
    this.showBatchImport = false
    this.formInfo.review_file = e.url
  }

  /**
   * @description 下载分配数据
   */
  async reviewExport(e) {
    e.stopPropagation()
    const params = {
      id: this.formInfo.id
    }
    try {
      const res: any = await Api.reviewExportApi(params)
      if (res.code === 200) {
        const url = await cosUploader.getObjectUrl(res.data.url)
        downloadFile(url, {
          filename: res.data.url.split('/').pop(),
        })
      } else {
        this.$message.error(res.message);
      }
    } catch (error) {
      this.$message.error(error.message);
    }
  }
}
</script>
<style lang="scss" scoped>

.el-radio {
  position: relative;

  .absolute-right {
    position: absolute;
    right: 0;
    top: -12px;
    transform: translate(100%, 0);

    /deep/.el-form-item__error {
      margin-top: -5px !important;
      padding-left: 10px;
    }
  }
}
.tips {
  color: #999;
  font-size: 12px;
}
.input-with-select {
  width: 400px;
}
.btn-add {
  font-size: 30px;
  margin-top: 25px;
  margin-left: 20px;
}
.item {
  padding: 20px;
  background-color: #f7f7f7;
  &:not(:last-child) {
    margin-bottom: 20px;
  }
  .el-input:nth-child(1) {
    margin-bottom: 10px;
  }
}

.item-temp {
  display: flex;
  align-items: center;
  &:not(:last-child) {
    margin-bottom: 10px;
  }
  .btn-add {
    margin-top: 0;
  }

  .el-checkbox {
    margin-left: 20px;
  }
}
.box {
  background-color: #f7f7f7;
  padding: 10px 20px;
  position: relative;

  .flex {
    height: 30px;
    line-height: 30px;
    text-align: center;
    padding-top: 10px;
    margin-top: 10px;
    border-top: 1px solid #e8e8e8;
    cursor: pointer;
  }
}
.list {
  max-height: 120px;
  overflow: hidden;

  &.expand {
    max-height: none;
    height: auto;
  }
}
</style>