<template>
  <div>
    <el-form ref="form" :model="formInfo" :rules="formRules" label-width="110px">
      <el-form-item label="模版名称" prop="name">
        <cy-input
          size="medium"
          width-size="medium"
          placeholder="请输入播种模版名称"
          v-model="formInfo.name"
          maxlength="20"
          show-word-limit
          :readonly="readonly"
        ></cy-input>
      </el-form-item>
      <el-form-item label="主题说明">
        <el-input
          type="textarea"
          size="medium"
          placeholder="请输入该播种模版的相关主题说明"
          v-model="formInfo.theme"
          show-word-limit
          :rows="5"
          :readonly="readonly"
          :maxlength="100" 
          style="width: 800px;"
        ></el-input>
      </el-form-item>
      <el-form-item label="播种配置" prop="sow_template_notice">
        <div class="tips">注：播种1的提醒时间及内容，由导购主动在手机端创建</div>
        <div class="list">
          <div v-for="(item, index) in formInfo.sow_template_notice" class="item" :key="index">
            <div class="item-box">
              <div class="item-label">【播种{{index+1}}】</div>
              <div class="item-value">
                <el-form-item v-if="index == 0" label="提醒时间">由导购在企业微信手机端的播种页面手动选择</el-form-item>
                <el-form-item v-if="index > 0" label="提醒时间">
                  <div class="flex">
                    <span>导购【播种1】提醒后的第</span>
                    <el-form-item class="m-l-10" :prop="`sow_template_notice[${index}].delay_value`" :rules="formRules.sow_delay_value">
                      <el-input v-model.number="item.delay_value" size="small" :readonly="readonly" style="width: 120px;" placeholder="请输入正整数"></el-input>
                    </el-form-item>
                    <el-form-item class="m-l-10" :prop="`sow_template_notice[${index}].delay_type`" :rules="formRules.sow_delay_type">
                      <el-select v-model="item.delay_type" size="small" :readonly="readonly" style="width: 80px;">
                        <el-option label="天" :value="1"></el-option>
                        <el-option label="周" :value="2"></el-option>
                        <el-option label="月" :value="3"></el-option>
                      </el-select>
                    </el-form-item>
                    <span class="m-l-10">，</span>
                    <el-form-item class="m-l-10" :prop="`sow_template_notice[${index}].delay_time`" :rules="formRules.sow_delay_time">
                      <el-time-select size="small" v-model="item.delay_time" value-format="HH:mm" :readonly="readonly" style="width: 140px !important;" placeholder="选择时间" ></el-time-select>
                    </el-form-item>
                    <span class="m-l-10">提醒。</span>
                  </div>
                </el-form-item>
                <el-form-item class="m-t-10" label="播种文案" :prop="`sow_template_notice[${index}].content`" :rules="formRules.sow_content">
                  <el-input
                    type="textarea"
                    size="medium"
                    placeholder="请输入播种提示语"
                    v-model="item.content"
                    maxlength="200"
                    show-word-limit
                    :rows="5"
                    :readonly="readonly"
                    style="width: 100%;"
                  ></el-input>
                </el-form-item>
              </div>
            </div>
            <div v-if="!readonly" class="item-btn">
              <el-button v-if="index==0 && formInfo.sow_template_notice.length <20" type="text" icon="el-icon-circle-plus" @click="handleAdd()"></el-button>
              <el-button v-if="index>0" type="text" icon="el-icon-remove" @click="handleDelete(index)"></el-button>
            </div>
          </div>
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
interface ISowTemplateItem {
  /** 迟提醒类型：1天 2周 3月 */
  delay_type: 1|2|3;
  /** 延迟提醒的时间数值 */
  delay_value: number;
  /** 延迟提醒当天的提醒时间，时:分 */
  delay_time: string;
  /** 提醒文案 */
  content: string
}
export interface IFormInfo {
  /** id */
  id?: number;
  /** 名称 */
  name: string;
  /** 主题说明 */
  theme: string;
  /** 模版状态 */
  status: number;
  /** 播种模版 */
  sow_template_notice: ISowTemplateItem[]
}

export const defaultFormInfo: IFormInfo = {
  name: '',
  theme: '',
  status: 2,
  sow_template_notice: [
    {
      delay_time: '',
      delay_type: 1,
      delay_value: 0,
      content: ''
    },
    {
      delay_time: '',
      delay_type: 1,
      delay_value: null,
      content: ''
    },
  ]
}

@Component({
  components: {
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
    name: [{ required: true, message: '请输入模版名称' }],
    sow_template_notice: [{ 
      required: true, 
      validator: (rule: any, value: any, callback: any) => {
        if (value.length === 0) {
          callback(new Error("请添加播种配置"));
        } else {
          callback();
        }
      },
    }],
    sow_delay_type: [{ required: true, message: '请选择时间单位' }],
    sow_delay_value: [{ 
      required: true,
      validator: (rule: any, value: any, callback: any) => {
        if(value===null || value==='') {
          callback(new Error('请输入提醒时间'))
        } else if (!/^[1-9]\d*$/.test(value)) {
          callback(new Error("请输入正整数"));
        } else {
          callback();
        }
      },
    }],
    sow_delay_time: [{ required: true, message: '请选择提醒时间' }],
    sow_content: [{ 
      required: true, message: '请输入播种提示语',
      validator: (rule: any, value: any, callback: any) => {
        if (value.trim().length === 0) {
          callback(new Error("请输入播种提示语"));
        } else {
          callback();
        }
      },
    }]
  }
  

  @Watch('data')
  onChange(value: IFormInfo) {
    if(value) {
      this.formInfo = value
    }
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

  handleAdd() {
    if(this.formInfo.sow_template_notice.length>=20) {
      this.$message.error('最多添加20条')
      return
    }
    this.formInfo.sow_template_notice.push({
      delay_time: '',
      delay_type: 1,
      delay_value: null,
      content: ''
    })
  }

  handleDelete(index: number) {
    this.formInfo.sow_template_notice.splice(index, 1)
  }
}
</script>

<style scoped lang="scss">
.tips {
  font-size: 12px;
  color: #999;
}
.item {
  display: flex;
  align-items: center;
  margin-bottom: 15px;

  .item-box {
    display: flex;
    align-items: center;
    flex: 1;
    background: #f6f7f8;
  }

  .item-label {
    height: 100%;
    width: 100px;
    text-align: center;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .item-value {
    background: #fff;
    border: 1px double #eee;
    padding: 15px 15px 15px 0;
    flex: 1;

    .flex {
      display: flex;
      align-items: flex-start;
    }
  }

  .item-btn {
    width: 50px;
    text-align: center;
    .el-button {
      font-size: 20px;
    }
  }
}
</style>