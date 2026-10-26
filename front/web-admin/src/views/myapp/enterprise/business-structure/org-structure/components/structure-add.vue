<template>
  <el-dialog :title="title" :visible.sync="visible" @close="handleCancel" @closed="closed">
    <el-form v-if="visible" :model="formData" :rules="formRules" ref="formData" label-width="100px">
      <el-form-item label="层级名称" prop="name">
        <el-input v-model="formData.name" placeholder="请输入层级名称" maxlength="50" show-word-limit></el-input>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="handleCancel">取 消</el-button>
      <el-button type="primary" @click="handleConfirm">确 定</el-button>
    </div>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Model, Vue } from 'vue-property-decorator';

interface IFormData {
  key: string;
  /** 节点ID */
  id: number;
  /** 名称（改名用） */
  name: string;
  /** 父级id（新增、移动用） */
  parent_id: number;
  /** 门店企微id，非门店不传或传0（新增门店用） */
  third_id: number;
  /** 兄弟节点ID */
  brother_id: number;
  /** true兄弟节点上方 false兄弟节点下方 */
  is_above: boolean;
}

const defaultFormData: IFormData = {
  key: 'n',
  name: '',
  parent_id: 0,
  third_id: 0,
  id: 0,
  brother_id: 0,
  is_above: false
}

@Component({
  name: 'StructureAdd'
})
export default class extends Vue {
  /** 弹框控制 */
  @Model('update:visible', { default: false } ) visible: boolean;
  /** 表单数据 */
  formData: IFormData = {...defaultFormData};

  formRules = {
    name: [
      { required: true, validator:(rule, value, callback) =>{
        if(!value.trim()) {
          callback(new Error('请输入层级名称'))
        } else {
          callback()
        }
      }, message: '请输入层级名称', trigger: 'blur' }
    ]
  }
  /** 弹框标题 */
  get title() {
    return {
      n: '新增一级节点',
      a: '重命名',
      b: '添加子节点',
      c: '向下添加节点',
      d: '向上添加节点'
    }[this.formData.key]
  }

  show(key?: string, value?: IFormData) {
    if(value) {
      this.formData = {key, ...value}
    }
    this.$emit('update:visible', true)
  }

  handleCancel() {
    this.$emit('update:visible', false)
  }

  handleConfirm() {
    (this.$refs.formData as any).validate(valid => {
      if (valid) {
        this.$emit('confirm', {
          ...this.formData,
          name: this.formData.name.trim()
        })
      }
    })
  }

  closed() {
    this.formData = {...defaultFormData}
  }
}
</script>