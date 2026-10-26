<template>
  <el-dialog title="关联门店" :visible.sync="visible" @close="handleCancel" @closed="closed">
    <el-form :model="formData" :rules="formRules" ref="formData" label-width="100px">
      <el-form-item label="选择分区" prop="parent_id">
        <el-cascader
          v-model="formData.parent_id"
          placeholder="请选择分区"
          :props="{
            checkStrictly: false,
            emitPath: false,
            value: 'id',
            label: 'name',
            children: 'children'
          }"
          :show-all-levels="true"
          :options="list"
          style="width: 100%;"></el-cascader>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="handleCancel">取 消</el-button>
      <el-button type="primary" @click="handleConfirm">确 定</el-button>
    </div>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Model, Prop, Vue } from 'vue-property-decorator';
import { forEachTree } from '@/utils/tree-tool';

interface IFormData {
  parent_id: number;
  third_id: number;
}

const defaultFormData: IFormData = {
  parent_id: 0,
  third_id: 0
}

@Component({})
export default class extends Vue {
  @Prop({ default: () => [] }) treeList: any[];
  /** 弹框控制 */
  @Model('update:visible', { default: false } ) visible: boolean;
  /** 表单数据 */
  formData: IFormData = {...defaultFormData};

  formRules = {
    parent_id: [
      { required: true, message: '请选择分区', trigger: 'blur' }
    ]
  }
  get list() {
    const data = JSON.parse(JSON.stringify(this.treeList))
    forEachTree<any>(data, node => {
      if(node.type === 3) {
        node.children = null
      }
      if(node.type < 3 && !node.children) {
        node.disabled = true
      }
    })
    return data
  }

  show(value: any) {
    this.formData = {
      ...value
    }
    this.$emit('update:visible', true)
  }

  handleCancel() {
    this.$emit('update:visible', false)
  }

  handleConfirm() {
    (this.$refs.formData as any).validate(valid => {
      if (valid) {
        this.$emit('confirm', {...this.formData})
      }
    })
  }

  closed() {
    this.formData = {...defaultFormData}
  }
}
</script>