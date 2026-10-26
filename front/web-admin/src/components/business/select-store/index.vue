<template>
  <el-dialog
    class="dialog-box"
    width="70vw"
    :title="title"
    :visible.sync="dialogVisible"
    :close-on-click-modal="false"
    @close="handleCancel"
  >
    <div class="content">
      <!-- 左侧组织树 -->
      <div class="content__left">
        <div class="data-header">
          <div class="tree-search clearfix">
            <div class="search-ipt pull-left">
              <cy-input
                class="input-with-select"
                v-model.trim="filter_text"
                size="small"
                width-size="small"
                :placeholder="`请输入${list_types[search_type - 1].name}`"
                clearable="clearable"
              >
                <el-select
                  v-model="search_type"
                  slot="prepend"
                  @change="
                    filter_text = '';
                    is_show_tree = true;
                  "
                >
                  <el-option
                    v-for="item in list_types"
                    :label="item.name"
                    :value="item.key"
                    :key="item.key"
                  ></el-option>
                </el-select>
              </cy-input>
            </div>
            <div class="search-btn pull-left m-l-20">
              <el-button size="small" type="primary" @click="handleSearch">
                查询
              </el-button>
              <el-button size="small" @click="reset">重置</el-button>
            </div>
          </div>
          <!-- 门店分类 -->
          <div class="cate-store" v-show="is_show_tree && showCate">
            <div class="cate-item">门店分类</div>
            <div
              class="cate-item"
              v-for="(item, index) in select_cate"
              :key="index"
            >
              <el-checkbox
                v-if="index !== 0"
                :disabled="item.store_list.length == 0"
                @change="handleSelectByCate(item, index)"
                :indeterminate="item.indeterminate"
                v-model="item.checked"
                >{{ item.name }}</el-checkbox
              >
            </div>
          </div>
        </div>
        <div class="data-content">
          <el-tree
            class="scroll"
            :data="dataList"
            ref="tree"
            empty-text="暂无可执行的部门"
            :props="default_props"
            :node-key="rowKey"
            show-checkbox
            default-expand-all
            :filter-node-method="filterNode"
            @check-change="handleCheckChange"
          
          >
            <div slot-scope="{ node, data }" class="custom-tree-node">
              <store-item :item="data"></store-item>
            </div>
          </el-tree>
        </div>
      </div>

      <!-- 右侧选择的数据 -->
      <div class="content__right">
        <div class="selected-title">
          {{ choice_name }}({{ select_arr.length }})
        </div>
        <div class="selected-content">
          <div class="scroll">
            <store-item v-for="(item, index) in select_arr" :key="index" :item="item" clearable @delete="handleRemove(index)"></store-item>
          </div>
        </div>
      </div>
    </div>
    <div slot="footer" class="dialog-footer">
      <template>
        <el-button @click="handleCancel" size="small">
          取消
        </el-button>
        <el-button @click="handleConfirm" size="small" type="primary">
          确 认
        </el-button>
      </template>
    </div>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from "vue-property-decorator";
import Api from "@/api/enterprise/enterprise";
import storeApi from "@/api/enterprise/material-library";
import StoreItem from './store-item.vue'

interface IStoreItem {
  department_id: number;
  dp_id: number;
  dp_name: string;
  dp_parent_id: number;
  dp_third_id: number;
  dp_third_parent_id: number;
  /** 是否有子节点 1-有；0-无 */
  is_child_department: number;
  store_name: string;
  store_type: number;
  third_id: number;
  children_list: Array<IStoreItem>;
}

interface ISelectCateItem {
  name: string;
  checked: boolean;
  store_list: IStoreItem[];
}

@Component({
  components: {
    StoreItem
  }
})
export default class extends Vue {
  @Prop({ default: "选择归属门店" }) title?: string; // 右侧选择标题
  @Prop({ default: () => [] }) selectArr?: Array<IStoreItem>; // 选择的数据
  @Prop({ default: true }) AllExpandClick?: boolean; // 是否允许点击父级展开下级
  @Prop({ default: "已选部门" }) choice_name?: string; // 右侧选择标题
  @Prop({ default: "" }) worker_status?: string | number;
  @Prop({ default: () => [] }) groupIds?: Array<any>;
  @Prop({ default: true }) showCate?: boolean;
  // 唯一标识 新版属性用的是dp_third_id, 默认用的department_id（接口已弃用）
  @Prop({ default: "department_id" }) rowKey?: string;
  /*** data ***/
  dialogVisible = true; // 控制模态框显示/隐藏
  expand_key: Array<number> = []; // 默认展开的key
  search_type = 1; // 搜索类型
  choice_type = 1;
  default_props = {
    label: "dp_name",
    children: "children_list"
  }; // 配置参数
  dataList: Array<IStoreItem> = [];
  filter_text = ""; // 搜索的字段
  choice_text = "";
  is_show_tree = true; // 是否显示tree
  filter_obj: Array<IStoreItem> = []; // 筛选出的数据
  select_arr: Array<IStoreItem> = []; // 已选择的部门数据
  select_cate: Array<ISelectCateItem> = [
    {
      name: "其他",
      checked: false,
      store_list: []
    },
    {
      name: "奥莱",
      checked: false,
      store_list: []
    },
    {
      name: "零售",
      checked: false,
      store_list: []
    }
  ]; // 门店门类
  get list_types() {
    return [
      { name: "门店名称", key: 1 }
      // { name: "门店编码", key: 2 }
    ];
  }

  /*** methods ***/

  mounted() {
    this.init()
  }

  /** 初始化 */
  async init() {
    await this.getDp();
    // 如果有分组ID
    if (this.groupIds.length === 1) {
      this.getList();
    }
    // 如果有选中数据
    if (this.selectArr.length > 0) {
      this.setSelectArr(this.dataList)
      this.setChoice();
    }
  }

  /** 确定 */
  handleConfirm() {
    if(this.select_arr.length === 0) {
      this.$message.error("请选择门店");
      return;
    }
    this.$emit("submit", this.select_arr);
  }
  /** 取消选择 */
  handleCancel() {
    this.$emit("close");
  }

  /** 通过接口获取选中数据 */
  getList() {
    storeApi
      .getGroupStoreData({
        source_group_id: this.groupIds[0],
        page: 1,
        size: 10000
      })
      .then(res => {
        // this.select_arr = res.data.list.map(item => {
        //   item.dp_name = item.department_name;
        //   return {
        //     ...item,
        //     dp_third_id: item.third_id
        //   };
        // });

        this.setSelectArr(this.dataList, res.data.list.map(item => {
          item.dp_name = item.department_name;
          return {
            ...item,
            dp_third_id: item.third_id
          };
        }))

        this.setChoice();
      });
  }


  setSelectArr(list, selectList = this.selectArr) {
    const selectIds = selectList.map(item => item[this.rowKey]);
    
    list.forEach(item => {
      if(selectIds.includes(item[this.rowKey])) {
        console.log(1231231);
        
        this.select_arr.push(JSON.parse(JSON.stringify(item)))
      }
      if(item.children_list && item.children_list.length > 0) {
        this.setSelectArr(item.children_list, selectList)
      }
    });
  }

  /** 通过分类选择数据 */
  handleSelectByCate(item: ISelectCateItem) {
    const oldSelectIds = this.select_arr.map(val => val[this.rowKey]);
    const selectSet = new Set(this.select_arr);
    const currentIds = item.store_list.map(val => val[this.rowKey])
    const parentIds = item.store_list.map(val => val.dp_third_parent_id)
    
    if (item.checked) {
      item.store_list.map(val => {
        if(oldSelectIds.includes(val[this.rowKey])) {
          const cur = this.select_arr.find(v => v[this.rowKey] === val[this.rowKey])
          selectSet.delete(cur);
        }
        selectSet.add(val)
      })
    } else {
      selectSet.forEach(val => {
        if (currentIds.includes(val[this.rowKey])) {
          selectSet.delete(val);
        }
        // 移除相关的父节点
        if (parentIds.includes(val[this.rowKey])) {
          selectSet.delete(val);
        }
      })
    }
    this.select_arr = Array.from(selectSet);
    const selectedIds = this.select_arr.map(val => val[this.rowKey]);
    // 变更选中状态
    (this.$refs.tree as any).setCheckedKeys(selectedIds);
  }

  /** 树节点选择变更 */
  handleCheckChange(data: any, checked: boolean) {
    console.log(data, checked);
    console.log(this.select_arr);
    
    
    if (checked) {
      // 选中
      const index = this.select_arr.findIndex(
        v => v[this.rowKey] === data[this.rowKey]
      );
      console.log(index);
      
      if (index < 0) {
        this.select_arr.push(data);
      }
    } else {
      // 取消选中
      this.select_arr = this.select_arr.filter(
        item => data[this.rowKey] !== item[this.rowKey]
      );
    }
    // 判断奥莱和零售选择的状态
    this.checkCateStatus()
  }

  /** 判断分类选择状态 */
  checkCateStatus() {
    this.select_cate.forEach((item, index) => {
      const selected = JSON.stringify(this.select_arr
      .filter(cell => {
        return cell.store_type === index;
      })
      .map(cell => {
        return cell[this.rowKey];
      })
      .sort())

      const currentId = JSON.stringify(item.store_list.map(val => val[this.rowKey])
      .sort())
      
      item.checked = currentId === selected;
    });
  }
  /** 查询门店 */
  async handleSearch() {

    (this.$refs.tree as any).filter(this.filter_text)
    // await this.getDp();
    // // 如果有分组ID
    // if (this.groupIds.length === 1) {
    //   this.getList();
    // }
    // // 如果有选中数据
    // if (this.selectArr.length > 0) {
    //   // this.setSelectArr(this.dataList)
    //   this.setChoice();
    // }
  }

  /** 搜索框重置 */
  reset() {
    this.filter_text = "";
    this.handleSearch()
  }

  /** 删除已选数据 */
  handleRemove(index: number) {
    // 数组移除
    const [current] = this.select_arr.splice(index, 1);
    // 修改选中状态
    (this.$refs.tree as any).setChecked(current[this.rowKey], false, true);
  }

  /** 获取部门 */
  async getDp() {
    try {
      const params: any = {}
      // 如果有搜索值
      if (this.filter_text.trim()) {
        params.dp_name = this.filter_text.trim();
      }
      const res = await Api.getTreeList(params);
      
      this.dataList = res.data.list || [];
      this.setCateStore(this.dataList);
      
    } catch (error) {
      console.error(error, "error");
    }
  }


  /** 设置分类状态 */
  setCateStore(list) {
    list.forEach(item => {
      this.select_cate[item.store_type].store_list.push(item);
      if(item.children_list) {
        this.setCateStore(item.children_list);
      }
    });
  }

  // 设置树选中
  setChoice() {
    this.select_arr.forEach(item => {
      (this.$refs.tree as any).setChecked(item[this.rowKey], true);
    });
  }

  filterNode(value, data) {
    console.log(value, data);
    
    if (!value) return true;
    return String(data.dp_store_name).toLowerCase().indexOf(value.toLowerCase()) !== -1;
  }
}
</script>

<style lang="scss" scoped>
/deep/.el-dialog {
  min-width: 800px !important;
  max-height: none !important;
  min-height: 600px !important;
  height: 60vh;
  display: flex;
  flex-direction: column;
  .el-dialog__header,
  .el-dialog__body,
  .el-dialog__footer {
    width: 100%;
    box-sizing: border-box;
  }

  .el-dialog__header {
    height: 52px !important;
  }

  .el-dialog__body {
    flex: 1;
    overflow: hidden;
    max-height: none;
  }
}
/deep/.el-select .el-input {
  width: 110px;
}
/deep/.el-input-group__append,
/deep/.el-input-group__prepend {
  background-color: $color-white;
}

/deep/.el-tree {
  padding: 0 !important;
}
/deep/.el-tree .el-tree-node__content:before {
  display: none !important;
}
.cate-store {
  display: flex;
  justify-content: flex-start;
  margin-top: 20px;
  .cate-item {
    margin-right: 20px;
  }
}
.content {
  height: 100%;
  display: flex;

  .content__left {
    flex: 1;
    height: 100%;
    box-sizing: border-box;
    padding-right: 5px;
    border-right: 2px solid #e4e7ed;
    display: flex;
    flex-direction: column;

    .data-header {
      padding-bottom: 20px;
    }
    
    .data-content {
      flex: 1;
      overflow: hidden;
    }
  }

  .content__right {
    width: 300px;
    height: 100%;
    box-sizing: border-box;
    padding: 10px;
    display: flex;
    flex-direction: column;
    .selected-title {
      margin-bottom: 10px;
      width: 100%;
      color: #787878;
    }
    .selected-content {
      flex: 1;
      overflow: hidden;
    }
  }
}
.default {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}



.scroll {
  height: 100%;
  overflow: auto;

  &::-webkit-scrollbar {
    /*滚动条整体样式*/
    width: 5px;
    /*高宽分别对应横竖滚动条的尺寸*/
    height: 1px;
  }

  &::-webkit-scrollbar-thumb {
    /*滚动条里面小方块*/
    border-radius: 2px;
    // box-shadow: inset 0 0 5px rgba(0, 0, 0, .2)
    background: #e4e7ed;
  }

  &::-webkit-scrollbar-track {
    /*滚动条里面轨道*/
    // box-shadow: inset 0 0 5px rgba(0, 0, 0, .2)
    border-radius: 2px;
    background: white;
  }
}


.custom-tree-node {
  height: 30px;
  line-height: 30px;
  margin-left: 6px;
}
.el-select-dropdown__item.selected {
  font-weight: normal;
}
</style>
