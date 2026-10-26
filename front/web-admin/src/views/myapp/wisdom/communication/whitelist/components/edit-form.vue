<template>
  <div>
    <el-form ref="form" :model="formInfo" :rules="formRules" label-width="110px">
      <el-form-item label="白名单名称" prop="name">
        <cy-input
          size="medium"
          width-size="medium"
          placeholder="请输入白名单名称"
          v-model="formInfo.name"
          maxlength="30"
          show-word-limit
          :disabled="readonly"
        ></cy-input>
      </el-form-item>
      <el-form-item label="白名单内容" prop="member_list">
        <el-radio-group class="m-b-10" v-model="formInfo.type" @change="handleTypeChange" :disabled="readonly">
          <el-radio :label="1">导入人员名单</el-radio>
          <el-radio :label="2">选择企业人员</el-radio>
        </el-radio-group>
        <el-row v-if="!readonly && formInfo.member_list.length==0" class="m-b-10">
          <el-button v-if="formInfo.type == 1" type="primary" @click="showBatchImport=true">导入白名单</el-button>
          <el-button v-if="formInfo.type == 2" type="primary" @click="showSelectMember=true">选择人员</el-button>
        </el-row>
        <!-- 选中的数据 -->
        <el-row v-if="formInfo.member_list.length>0">
          <el-table :data="[{list:formInfo.member_list}]">
            <el-table-column prop="name" label="员工">
              <template slot-scope="scope">
                <el-tag class="m-r-10 m-b-10" v-for="(item, index) in scope.row.list" :key="index"><i class="el-icon-user-solid"></i>{{item.name||item.member_name}}</el-tag>
              </template>
            </el-table-column>
            <el-table-column v-if="!readonly" label="操作" width="120">
              <el-button v-if="formInfo.type===1" type="text" @click="showBatchImport=true">重新导入</el-button>
              <el-button v-if="formInfo.type===2" type="text" @click="showSelectMember=true">修改</el-button>
            </el-table-column>
          </el-table>
        </el-row>

      </el-form-item>
    </el-form>
    <!-- 导入人群包 -->
    <batch-import
      v-model="showBatchImport"
      ref="batchImport"
      @confirm="handleBatchImportConfirm"
    />
    <!-- 选择员工 -->
    <!-- <select-member
      v-if="showSelectMember"
      member-key="member_id"
      :hasSelectArr="formInfo.member_list"
      @close="showSelectMember = false"
      @submit="handleSelectMemberConfirm"
    /> -->

    <select-member      
      :show.sync="showSelectMember"
      :value="selectMember"
      @cancel="showSelectMember = false"
      @confirm="handleSelectMemberConfirm"
    />
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import BatchImport from './batch-import.vue';
// 选择员工
import SelectMember from "../../components/select-member.vue";
// import SelectMember from "@/components/business/select-member-copy/index.vue";

export interface IFormInfo {
  /** 白名单id */
  id?: number;
  /** 白名单名称 */
  name: string;
  /** 模版文件 */
  file_url: string;
  /** 类型 1导入名单 2选择企业人员 */
  type: number;
  /** 白名单内容 */
  member_list: any[];
}

export const defaultFormInfo: IFormInfo = {
  name: '',
  file_url: '',
  type: 1,
  member_list: []
}

@Component({
  components: {
    BatchImport,
    SelectMember
  }
})
export default class EditForm extends Vue {
[x: string]: any;
  /** 是否只读 */
  @Prop({ type: Boolean, default: false }) readonly: boolean;
  /** 表单数据 */
  @Prop({ type: Object, default: () => defaultFormInfo }) data: IFormInfo;
  /** 表单数据 */
  formInfo: IFormInfo = JSON.parse(JSON.stringify(defaultFormInfo))
  /** 表单校验规则 */
  formRules = {
    name: [{ required: true, message: '请输入白名单名称' }],
    member_list: [{ 
      required: true, 
      validator: (rule: any, value: any, callback: any) => {
        if (value.length === 0) {
          callback(new Error("请添加白名单内容"));
        } else {
          callback();
        }
      },
    }]
  }
  /** 显示导入人群包弹框 */
  showBatchImport = false;
  /** 选择员工 */
  showSelectMember = false


  get selectMember() {
    return this.formInfo.member_list.map(item => {
      return {
        label: item.name,
        value: item.member_id,
        dp_parent_id: item.dp_parent_id
      }
    })
  }

  @Watch('data')
  onChange(value: IFormInfo) {
    if(value) {
      this.formInfo = value
    }
  }

  mounted() {
    if(document.getElementsByName('v-modal')[0]) {
      document.getElementsByName('v-modal')[0].style.display='none'
    }
  }
  /** 
   * 类型切换
   */
  handleTypeChange() {
    this.formInfo.member_list = [];
    (this.$refs.form as any).clearValidate('member_list')
  }

  /**
   * 模版导入确认
   */
  handleBatchImportConfirm(value) {
    this.formInfo.member_list = value.member_list
    this.formInfo.file_url = value.file_url
    this.showBatchImport = false
  }

  /**
   * 选择员工确认
   */
  handleSelectMemberConfirm(value) {
    this.formInfo.member_list = value.map(item => {
      return {
        member_id: item.value,
        name: item.label,
      }
    });
    this.showSelectMember = false;
    (this.$refs.form as any).validateField('member_list')
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
}
</script>