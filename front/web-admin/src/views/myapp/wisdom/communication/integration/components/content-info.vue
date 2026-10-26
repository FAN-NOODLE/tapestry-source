<template>
  <div class="content-info">
    <el-button type="primary" size="medium" @click="handleAdd">+增加部门传讯</el-button>
    <div class="item" v-for="(item, index) in data" :key="index">
      <el-button v-if="index>0" class="del" type="text" @click="handleDelete(index)">删除</el-button>
      <el-form-item label="所属部门" :prop="`content[${index}].department_name`" :rules="rules.department_name">
        <el-select
          v-model="item.department_name"
          size="medium"
          placeholder="请选择所属部门"
          filterable
          style="width: 400px;"
          @change="handleChange(index, [])"
        >
          <el-option v-for="(option, optionIndex) in departmentOptions" :key="optionIndex" :label="option.label" :value="option.value"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item class="m-t-10" label="关联内容" :prop="`content[${index}].info_list`" :rules="rules.info_list">
        <el-button type="primary" @click="handleSelect(index, item.department_name)">关联部门传讯</el-button>
        <el-row class="tags">
          <el-tag v-for="(info, infoIndex) in item.info_list" :key="infoIndex" class="m-r-10 m-t-10" closable @close="handleClose(item, infoIndex)">{{info.title}}</el-tag>
        </el-row>
      </el-form-item>
    </div>
    <select-content ref="select" :value="currentValue" :disabledList="disabledList" @change="handleSelectConfirm"></select-content>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import Api from '@/api/wisdom/communication';
import SelectContent from './select-content.vue'

interface IInfoItem {
  [key:string]:any;
  id: number;
  title: string;
}

export interface IContentItem {
  department_name: string;
  info_list?: IInfoItem[];
  info_ids?: number[];
}

export const defaultItem: IContentItem = {
  department_name: '',
  info_list: []
}


@Component({
  model: {
    prop: 'data',
    event: 'change'
  },
  components: {
    SelectContent
  }
})
export default class ContentInfo extends Vue {
  @Prop({ type: Array, default: () => []}) data: IContentItem[];

  rules = {
    department_name: [{ required: true, message: '请选择所属部门' }],
    info_list: [{
      required: true,
      validator: (rule: any, value: any, callback: any) => {
        if (value.length === 0) {
          callback(new Error("请选择关联内容"));
        } else {
          callback();
        }
      },
    }]
  }
  
  /** 所属部门选项 */
  departmentOptions = []

  /** 选择内容 */
  current = -1
  currentValue: IInfoItem[] = []

  get disabledList() {
    return this.data.filter((item, index) => {
      return index !== this.current
    }).reduce((pre, cur) => {
      return pre.concat(cur.info_list.map(item => item.id))
    }, [])
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
   * 选择部门传讯
   */
  handleSelect(index: number, department:string) {
    this.current = index;
    this.currentValue = this.data[index].info_list;
    (this.$refs.select as any).handleShow(department)
  }

  /** 
   * 选择部门传讯 确认
   */
  handleSelectConfirm(value: IInfoItem[]) {
    this.data[this.current].info_list = value
    this.$emit('change', this.data)
  }

  /**
   * 添加
   */
  handleAdd() {
    this.data.push(JSON.parse(JSON.stringify(defaultItem)))
    this.$emit('change', this.data)
  }

  /**
   * 删除一组数据
   * @param index 
   */
  handleDelete(index: number) {
    this.data.splice(index, 1)
    this.$emit('change', this.data)
  }

  /**
   * 解除关联
   * @param item 
   * @param index 
   */
  handleClose(item: IContentItem, index: number) {
    item.info_list.splice(index, 1)
  }

  /** 
   * 传讯选择
   */
  handleChange(index: number, value: any[]) {
    this.data[index].info_list = value
    this.$emit('change', this.data)
  }
}
</script>

<style lang="scss" scoped>
.item {
  border: 1px dashed #ccc;
  margin-top: 10px;
  padding: 15px;
  border-radius: 5px;
  position: relative;

  .del {
    position: absolute;
    top: 10px;
    right: 10px;
    z-index: 1;
  }
}
.tags {
  line-height: normal;
  
}
#crs-scrm .el-tag {
  display: inline-block;
  max-width: 100% !important;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  position: relative;
  padding: 0 30px 0 10px !important;

  ::v-deep .el-tag__close {
    position: absolute;
    right: 5px;
    top: 50%;
    transform: translateY(-50%) !important;
  }
}

</style>