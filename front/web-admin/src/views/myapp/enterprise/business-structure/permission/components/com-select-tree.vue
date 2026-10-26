<template>
  <!-- 选择树弹框 -->
  <el-dialog
    :title="title"
    :visible.sync="visible"
    append-to-body
    @opened="handleOpened"
    @closed="handleClose"
  >
    <div class="com-select-tree__content">
      <div class="content-top">
        <el-input
          v-model="keyword"
          placeholder="请输入层级或门店名称"
          style="width:280px"
          class="m-r-20"
          clearable
          @keyup.enter.native="handleSearch"
        ></el-input>
        <el-button size="small" type="primary" @click="handleSearch">查询</el-button>
        <el-button size="small" @click="handleReset">重置</el-button>
      </div>
      <div class="content-main">
        <div class="content-left">
          <el-tree
            ref="tree"
            :data="treeList"
            :props="defaultProps"
            :default-expand-all="true"
            :default-checked-keys="checked"
            :show-checkbox="true"
            node-key="id"
            :check-strictly="false"
            :expand-on-click-node="true"
            :filter-node-method="filterNode"
            @check-change="handleCheckChange"
          >
          </el-tree>
        </div>
        <div class="content-right">
          <div class="content-right__title">已选管理范围</div>
          <div class="content-right__list">
            <div class="content-right__item" v-for="(item, index) in checkedList" :key="index">
              <div class="content-right__item__name">
                <svg-icon v-if="item.type===4" class-name="m-r-5" icon-class="store"></svg-icon>
                <svg-icon v-else class-name="m-r-5" icon-class="organization"></svg-icon>
                <span class="name">{{ item.name }}</span>
              </div>
              <el-button class="content-right__item__del" size="mini" type="text" @click="handleDel(item)">
                <i class="el-icon-delete"></i>
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div slot="footer" class="dialog-footer">
      <el-button @click="handleClose">取 消</el-button>
      <el-button type="primary" @click="handleConfirm">确 定</el-button>
    </div>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Model, Prop, Vue } from "vue-property-decorator";


@Component({
  name: "ComSelectTree"
})
export default class ComSelectTree extends Vue {
  @Model('update:visible') visible!: boolean;
  @Prop({ default: "选择管理范围" }) title!: string;
  @Prop({ default: false }) hideDefaultButton!: boolean;
  @Prop({ default: () => [] }) checked!: any[];
  @Prop({ default: () => [] }) treeList!: any[];
  
  defaultProps = {
    children: "children",
    label: 'name'
  };

  /** 搜索关键字 */
  keyword = ''

  /** 当前选中节点 */
  checkedList = []

  handleOpened() {
    (this.$refs.tree as any).setCheckedKeys(this.checked)
    this.handleCheckChange()
  }
  
  /** 点击事件 */
  handleClick() {
    this.$emit("update:visible", true);
  }
  /** 关闭弹框 */
  handleClose() {
    this.$emit("update:visible", false);
  }
  /** 确定 */
  handleConfirm() {
    const ids = this.checkedList.map((item: any) => item.id);
    this.$emit("confirm", ids);
  }
  /** 选中节点 */
  handleCheckChange() {
    this.checkedList = (this.$refs.tree as any).getCheckedNodes()
  }

  /**
   * 搜索
   */
  // @Watch('keyword')
  handleSearch() {
    (this.$refs.tree as any).filter(this.keyword);
  }
  /** 
   * 重置
   */
  handleReset() {
    this.keyword = '';
    this.handleSearch()
  }

  /**
   * 搜索节点
   * @param value 
   * @param data 
   */
  filterNode(value: string, data: any) {
    if (!value) return true;
    return data.name.indexOf(value) !== -1;
  }

  handleDel(item: any) {
    (this.$refs.tree as any).setChecked(item, false, true)
  }
}

</script>
<style lang="scss" scoped>

.ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}
.flex {
  display: flex;
  align-items: center;
}

.com-select-tree__content {
  height: 50vh;
  display: flex;
  flex-direction: column;

  .content-top {
    width: 100%;
  }

  .content-main {
    width: 100%;
    flex: 1;
    overflow: hidden;
    display: flex;
    margin-top: 20px;
  }

  .content-left {
    flex: 1;
    overflow-x: hidden;
    overflow-y: scroll;
  }

  .content-right {
    width: 300px;
    overflow-y: scroll;
    border-left: 1px solid #e6e6e6;
    padding: 0 15px;

    .content-right__title {
      color: #999;
    }

    .content-right__list {
      margin-top: 10px;
    }

    .content-right__item {
      display: flex;
      align-items: center;
      &:hover {
        background-color: #f5f7fa;
        cursor: pointer;
      }
      .content-right__item__name {
        flex: 1;
        padding-right: 20px;
        overflow: hidden;
        .name {
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }
      }
    }
  }
}
</style>