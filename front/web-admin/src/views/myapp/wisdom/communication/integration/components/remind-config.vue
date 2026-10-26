<template>
  <span class="m-l-20">
    <el-button type="primary" @click="handleShow">提交传讯提醒</el-button>
    <el-dialog
      title="传讯提醒配置"
      :visible.sync="visible"
      width="600px"
      :append-to-body="true"
      center>
      <div>
        <el-form ref="form" :model="formInfo" label-width="120px">
          <div v-for="(item, index) in formInfo.notice_config_list" :key="index">
            <div class="title">{{ ['周传讯提醒', '转货传讯提醒'][item.info_type] }}</div>
            <el-form-item label="提交截止时间:"> 
              <div class="flex">
                <span>每周</span>
                <el-form-item class="m-l-10" label-width="0px" :prop="`notice_config_list[${index}].deadline_week`" :rules="formRules.deadline_week">
                  <el-select size="small" v-model="item.deadline_week" placeholder="请选择" style="width: 150px;">
                    <el-option
                      v-for="item in weekOptions"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value"
                    />
                  </el-select>
                </el-form-item>
                <el-form-item class="m-l-10" label-width="0px" :prop="`notice_config_list[${index}].deadline_time`" :rules="formRules.deadline_time">
                  <el-time-picker
                    size="small"
                    v-model="item.deadline_time"
                    value-format="HH:mm:ss"
                    placeholder="选择时间"
                    style="width: 150px;">
                  </el-time-picker>
                </el-form-item>
              </div>
            </el-form-item>
            <el-form-item label="发布时间:" > 
              <div class="flex">
                <span>每周</span>
                <el-form-item class="m-l-10" label-width="0px" :prop="`notice_config_list[${index}].release_week`" :rules="formRules.release_week">
                  <el-select size="small" v-model="item.release_week" placeholder="请选择" style="width: 150px;">
                    <el-option
                      v-for="item in weekOptions"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value"
                    />
                  </el-select>
                </el-form-item>
                <el-form-item class="m-l-10" label-width="0px" :prop="`notice_config_list[${index}].release_time`" :rules="formRules.release_time">
                  <el-time-picker
                    size="small"
                    v-model="item.release_time"
                    value-format="HH:mm:ss"
                    placeholder="选择时间"
                    style="width: 150px;">
                  </el-time-picker>
                </el-form-item>
              </div>
            </el-form-item>
          </div>
        </el-form>
      </div>
      <div slot="footer">
        <el-button @click="visible=false">取 消</el-button>
        <el-button type="primary" @click="handleConfirm">确 定</el-button>
      </div>
    </el-dialog>
  </span>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Api from '@/api/wisdom/communication';

@Component
export default class RemindConfig extends Vue {
  visible = false

  weekOptions = [
    {label: '星期一', value: 1},
    {label: '星期二', value: 2},
    {label: '星期三', value: 3},
    {label: '星期四', value: 4},
    {label: '星期五', value: 5},
    {label: '星期六', value: 6},
    {label: '星期日', value: 7}
  ]

  formInfo = {
    notice_config_list: [
      {
        info_type: 0,
        deadline_week: '',
        deadline_time: '',
        release_week: '',
        release_time: ''
      },
      {
        info_type: 1,
        deadline_week: '',
        deadline_time: '',
        release_week: '',
        release_time: ''
      },
    ]
  }

  formRules = {
    'deadline_week': [
      {required: true, message: '请选择提交截止日期', trigger: 'change'}
    ],
    'deadline_time': [
      {required: true, message: '请选择提交截止时间', trigger: 'change'}
    ],
    'release_week': [
      {required: true, message: '请选择发布日期', trigger: 'change'}
    ],
    'release_time': [
      {required: true, message: '请选择发布时间', trigger: 'change'}
    ]
  }

  getData() {
    Api.getNoticeConfigApi({}).then((res: any) => {
      this.formInfo.notice_config_list = res.data
    })
  }

  handleShow() {
    this.visible = true
    this.getData()
  }

  handleConfirm() {
    // this.visible = false;
    (this.$refs.form as any).validate(valid => {
      if (valid) {
        Api.updateNoticeConfigApi(this.formInfo).then((res: any)=> {
          if (res.code == 200) {
            this.$message.success('保存成功!')
            this.visible = false;
          } else {
            this.$message.error(res.message);
          }
        })
      }
    })
  }
}
</script>

<style lang="scss" scoped>
.flex {
  display: flex;
  align-items: center;
}
.title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 10px;
}
</style>