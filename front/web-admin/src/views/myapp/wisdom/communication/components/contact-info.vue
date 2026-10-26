<template>
  <el-row class="upload-attachment">
    <el-row class="m-b-10">
      <el-button type="primary" @click="handleAdd">+添加一行</el-button>
      <span class="tips m-l-10">*每行添加1位联系人信息，多个联系人可添加多行；同部门联系人仅在部门第一个联系人前填写部门即可。</span>
    </el-row>
    <!-- 文件列表 -->
    <el-table v-if="value.length" :data="value" stripe>
      <el-table-column label="部门">
        <template slot-scope="scope">
          <el-input v-model="scope.row.department_name" @blur="handleBlur('department_name', scope.$index)" placeholder="请输入部门" maxlength="100"></el-input>
        </template>
      </el-table-column>
      <el-table-column label="姓名">
        <template slot="header">
          <span style="color: #F67170">*</span>
          <span>姓名</span>
          <span style="color: #F67170">（必填）</span>
        </template>
        <template slot-scope="scope">
          <el-form-item :prop="`contact_info[${scope.$index}].name`" :rules="{ required: true, message: '请输入姓名' }">
            <el-input v-model="scope.row.name" @blur="handleBlur('name', scope.$index)" placeholder="请输入姓名" maxlength="100"></el-input>
          </el-form-item>
        </template>
      </el-table-column>
      <el-table-column label="邮箱" width="320">
        <template slot-scope="scope">

          <!-- <el-input v-model="scope.row.email" @change="handleChange($event, 'email_suffix', scope.$index)" @blur="handleBlur('email', scope.$index)" placeholder="请输入邮箱" maxlength="100">
            <el-select v-model="scope.row.email_suffix" slot="append" placeholder="请选择" style="width:140px">
              <el-option label="@coach.com" value="@coach.com"></el-option>
              <el-option label="@tapestry.com" value="@tapestry.com"></el-option>
            </el-select>
          </el-input> -->
          <el-select v-model="scope.row.email"
            remote
            filterable
            allow-create
            default-first-option
            clearable
            :remote-method="handleRemoteMethod"
            placeholder="请输入邮箱" style="width:100%"
            @change="handleEmailChange">
            <el-option v-for="(item, index) in emailList" :label="item.label" :value="item.value" :key="index"></el-option>
          </el-select>
        </template>
      </el-table-column>
      <el-table-column label="电话">
        <template slot-scope="scope">
          <el-input v-model="scope.row.tel_phone" @blur="handleBlur('tel_phone', scope.$index)" placeholder="请输入电话" maxlength="100"></el-input>
        </template>
      </el-table-column>
      <el-table-column label="备注">
        <template slot-scope="scope">
          <el-input v-model="scope.row.remark" @blur="handleBlur('remark', scope.$index)" placeholder="请输入备注" maxlength="100"></el-input>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100">
        <template slot-scope="scope">
          <el-button type="text" @click="handleRemove(scope.$index)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-row>
</template>
<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator'

/** 联系人 */
export interface IContactItem {
  /** 部门名称 */
  department_name: string;
  /** 姓名 */
  name: string;
  /** 邮箱 */
  email: string;
  /** 手机号 */
  tel_phone: string;
  /** 备注 */
  remark: string;
}

const defaultItem: IContactItem = {
  department_name: '',
  name: '',
  email: '',
  tel_phone: '',
  remark: ''
}

@Component({
  model: {
    prop: 'value',
    event: 'change'
  }
})
export default class ContactInfo extends Vue {
  /** 联系人列表 */
  @Prop({ default: () => [] }) value!: IContactItem[]

  options = [
    { label: '@coach.com', value: '@coach.com' },
    { label: '@tapestry.com', value: '@tapestry.com' },
  ]

  emailList = []

  handleRemoteMethod(query:string) {
    const name = query.split('@')[0]
    this.emailList = this.options.map(item => {
      return { label: `${name}${item.label}`, value: `${name}${item.label}` }
    })
  }

  handleEmailChange() {
    this.emailList = []
  }

  /**
   * 添加一行
   */
  handleAdd() {
    this.value.push({ ...defaultItem })
    this.$emit('change', this.value);
  }

  /**
   * 删除一行
   */
  handleRemove(index: number) {
    this.value.splice(index, 1)
    this.$emit('change', this.value);
  }

  handleBlur(key: string, index: number) {
    this.value[index][key] = this.value[index][key].trim()
    this.$emit('change', this.value);
  }

  handleChange(value: string, key: string, index: number) {
    const regex = /^[a-zA-Z0-9_.-]+(@coach.com|@tapestry.com)$/
    
    const match = value.match(regex)
    
    if(match) {
      this.value[index][key] = match[1]
      this.value[index].email = value.replace(match[1], '')
    }

    this.$emit('change', this.value);
  }
}
</script>
<style lang="scss" scoped>
::v-deep .el-form-item__content {
  line-height: normal;
}
#crs-scrm .el-table ::v-deep .cell {
  padding: 0 10px;
}
.upload-attachment {
  width: 100%;
}
.tips {
  font-size: 12px;
  color: #606266;
}
</style>