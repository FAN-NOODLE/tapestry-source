<template>
  <div class="pistribution-modal">
    <div class="mina-box">
      <div class="major">
        <div class="title clearfix">
          <div class="left">添加支付即导购</div>
          <div @click="cancelSave()" class="right el-icon-close"></div>
        </div>
        <el-row class="search">
          <el-input
            :placeholder="`请输入${type_list[gudie_type].name}`"
            class="search-input"
            size="medium"
            v-model="search_key"
            @keyup.enter.native="search"
          >
            <el-select
              slot="prepend"
              v-model="gudie_type"
              class="input-with-select"
            >
              <el-option
                v-for="item in type_list"
                :value="item.id"
                :label="item.name"
                :key="item.id"
              ></el-option>
            </el-select>
          </el-input>

          <el-select
            :value="select_dp.name"
            placeholder="请选择所属部门"
            size="medium"
            clearable="clearable"
            @change="changeSearch"
            class="tree-select"
          >
            <el-option :value="select_dp.name" :label="select_dp.name">
              <el-tree
                class="treeStyle"
                :props="dp_props"
                :load="loadNode"
                node-key="dp_id"
                :expand-on-click-node="true"
                :check-strictly="true"
                :check-on-click-node="true"
                @node-click="handleNodeClick"
                lazy
                ref="tree"
              ></el-tree>
            </el-option>
          </el-select>

          <el-button
            size="small"
            @click="search"
            class="mar-l-20"
            type="primary"
            >查询</el-button
          >
          <el-button size="small" @click="reset">重置</el-button>
        </el-row>

        <el-row class="box">
          <el-table
            ref="selectShopList"
            :data="list"
            size="mini"
            height="350"
            tooltip-effect="dark"
            @selection-change="handleSelectionChange"
            @select-all="changeAll"
            :row-key="getRowKeys"
            :header-cell-style="{
              fontSize: '14px',
              color: '#606266',
              fontWeight: '400'
            }"
          >
            <el-table-column
              type="selection"
              :selectable="checkSelectable"
              width="55"
              :reserve-selection="true"
            ></el-table-column>

            <el-table-column
              prop="worker_no"
              label="导购ID"
              min-width="150"
            ></el-table-column>

            <el-table-column label="导购" min-width="150">
              <div slot-scope="scope">
                <img
                  :src="
                    scope.row.mem_avatar
                      ? scope.row.mem_avatar
                      : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                  "
                  alt=""
                  width="30"
                  height="30"
                  class="table-text-img"
                /><span v-text="scope.row.mem_name"></span>
              </div>
            </el-table-column>

            <el-table-column
              prop="mem_mobile"
              label="手机号"
              min-width="120"
            ></el-table-column>

            <el-table-column
              prop="dp_name"
              :show-overflow-tooltip="true"
              label="所在部门"
              min-width="150"
            ></el-table-column>
            <div slot="empty">
              <layout-empty
                v-if="!search_key"
                :message="'暂无导购信息'"
              ></layout-empty>
              <layout-empty
                v-if="search_key"
                empty_type="tag"
                :message="'没有找到相关导购'"
              ></layout-empty>
            </div>
          </el-table>
        </el-row>

        <el-row class="footer">
          <el-col :span="24">
            <el-pagination
              :page-size="meta.size"
              :total="meta.total"
              layout="prev, pager, next"
              :current-page="meta.page"
              @current-change="handleCurrentChange"
              class="page-right font-weight-400"
            ></el-pagination>
          </el-col>

          <el-col :span="24" class="p-t-20">
            <el-button type="primary" size="small" @click="save" class="fr"
              >确定</el-button
            >
            <el-button size="small" @click="cancelSave" class="fr m-r-20"
              >取消</el-button
            >
          </el-col>
        </el-row>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
import Api from "@/api/wisdom/Guide";
import personApi from "@/api/wisdom/Personnel";
// import { component } from "vue/types/umd";

// import Api from "@/api/organization/Guide";
@Component({
  name: "AddGuide",
  components: {
    LayoutEmpty
  }
})
export default class extends Vue {
  private list = [];
  private meta = {
    total: 40, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private search_key = "";
  private gudie_type = 0;
  private type_list = [
    { name: "导购姓名", id: 0, key: "mem_name" },
    { name: "导购手机号", id: 1, key: "mem_mobile" },
    { name: "导购ID", id: 2, key: "worker_no" }
  ];
  private dp_props = {
    label: "dp_name",
    children: [],
    isLeaf: "leaf"
  };
  private select_dp = { name: "", value: "" };
  private checked_shop: any = [];
  private fetchList() {
    const params = {
      search: {
        dp_id: this.select_dp.value,
        filter_type: 1
      },
      page: this.meta.page,
      size: this.meta.size
    };
    params.search[this.type_list[this.gudie_type].key] = this.search_key;
    Api.wisdomGuideList(params).then(res => {
      const { list, page, size, total } = res.data;
      this.list = list;
      list.forEach((item: any) => {
        if (item.is_guide === 1) {
          (this.$refs.selectShopList as any).toggleRowSelection(item, true);
        }
      });
      this.meta = { page, size, total };
      // this.page_init = true;
    });
  }
  //返回Id
  private getRowKeys(row: any) {
    return row.id;
  }
  //多选选择商品
  private handleSelectionChange(val: any) {
    this.checked_shop = [];
    this.checked_shop = [...val];
  }
  //全选数量限制
  private changeAll(selection: any) {
    console.log(selection);
  }
  //判断是否禁用
  private checkSelectable(row: any) {
    return row.is_guide === 0;
  }
  // 当前页
  private handleCurrentChange(val: any) {
    this.meta.page = val;
    this.fetchList();
  }
  // 保存
  private save() {
    const arr = this.checked_shop.filter((v: any) => v.is_guide !== 1);
    if (arr.length <= 0) {
      this.$message.warning("请选择导购");
      return;
    }
    this.$emit("success", arr);
  }
  // 取消
  private cancelSave() {
    this.$emit("fail", "取消");
  }
  // 选择部门
  private handleNodeClick(node: any) {
    this.select_dp.name = node.dp_name;
    this.select_dp.value = node.dp_id;
    this.search();
  }
  // 搜索
  private search() {
    this.meta.page = 1;
    this.fetchList();
  }
  //动态加载节点
  private async loadNode(node: any, resolve: any) {
    if (node.level === 0) {
      const {
        data: { list }
      } = await personApi.depListApi({});
      list.forEach((item: any) => {
        if (!item.is_child_department) item.leaf = true;
      });
      resolve(list);
    }
    if (node.level >= 1) {
      this.getIndex(node, resolve);
    }
  }
  //判断部门是否有子节点
  private setChecked(arr: any) {
    arr.forEach((item: any) => {
      if (!item.is_child_department) item.leaf = true;
    });
  }
  // 加载叶子节点数据函数
  private getIndex(node: any, resolve: any) {
    const params = { dp_id: node.data.dp_id };
    personApi.depListApi(params).then(res => {
      this.setChecked(res.data.list);
      resolve(res.data.list);
    });
  }
  private changeSearch() {
    this.select_dp.value = "";
    this.select_dp.name = "";
    this.search();
    // this.$refs.tree.setCurrentKey(null);
  }
  // 重置
  private reset() {
    this.search_key = "";
    this.gudie_type = 0;
    this.select_dp.value = "";
    this.select_dp.name = "";
    this.search();
  }
  // 生命周期
  private created() {
    console.log("created");
  }
  private mounted() {
    this.fetchList();
    console.log("mounted");
  }
  private destroyed() {
    console.log("destroyed");
  }
  // private changeGuideType() {}
}
</script>
<style lang="scss" scoped>
.search {
  margin-bottom: 20px;
}

.mar-l-20 {
  margin-left: 20px;
}

.search-input {
  width: 280px !important;
  margin-bottom: 0px !important;
  margin-right: 20px;
}

.input-with-select {
  width: 115px;
}

/deep/ {
  .el-scrollbar {
    .el-scrollbar__view {
      .el-select-dropdown__item {
        height: auto;
      }
    }
  }
}

.el-select-dropdown__item.selected {
  font-weight: normal;
  height: auto !important;
}
.el-icon-refresh-right {
  color: #409eff;
  cursor: pointer;
}

.treeStyle {
  .el-tree-node.is-current.is-focusable.is-checked {
    color: #409eff !important;
  }
}

.main-content .title {
  position: static;
}

.left {
  padding-bottom: 0;
}

.tree-select {
  width: 280px !important;
}

/deep/.el-pagination {
  padding-top: 0 !important;
}

/deep/.el-table__fixed-right-patch {
  background: #f0f2f6;
}

/deep/.el-button + .el-button {
  margin-left: 10px !important;
}

.page-right {
  float: right;
}

.p-t-20 {
  padding-top: 20px;
}

.pistribution-modal {
  width: 100%;
  height: 100%;
  position: fixed;
  top: 0;
  left: 0;
  margin: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 10;
}

.pistribution-modal {
  ::-webkit-scrollbar {
    width: 4px;
  }
}

::-webkit-scrollbar-track {
  background-color: #ffffff;
}

::-webkit-scrollbar-thumb {
  background-color: #dfe1e4;
}

::-webkit-scrollbar-thumb:hover {
  background-color: #dfe1e4;
}

::-webkit-scrollbar-thumb:active {
  background-color: #dfe1e4;
}

.fr {
  float: right;
}

.m-r-20 {
  margin-right: 20px;
}

.mina-box {
  position: relative;
  width: 100%;
  height: 100%;
}
.table-text-img {
  border-radius: 100%;
  vertical-align: middle;
  margin-right: 4px;
}
.major {
  position: absolute;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  margin: auto;
  width: 800px;
  height: 600px;
  background: #ffffff;
  border-radius: 4px;
  overflow: hidden;
  padding: 0 30px;
}
.title {
  padding-top: 28px;
  padding-bottom: 20px;
}

.right {
  float: right;
  cursor: pointer;
  color: #000;
  font-size: 20px;
}

.empty-page {
  padding-top: 0 !important;
  /deep/ .default {
    .empty-list-img {
      width: 124px !important;
      height: 134px !important;
    }
  }
}

.left {
  float: left;
  color: #33363c;
  font-size: 18px;
  font-weight: 600;
}

.footer {
  padding-top: 20px;
}
.empty-box {
  display: flex;
  flex-direction: column;
  align-items: center;
}
</style>
