<template>
  <div class="supernatant">
    <layout-list>
      <template slot="buttons">
        <el-button type="primary" size="small" @click="addGuide"
          >添加支付即导购</el-button
        >
      </template>
      <template slot="search">
        <search-item text="导购：">
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
                :value="item.key"
                :label="item.name"
                :key="item.key"
              ></el-option>
            </el-select>
          </el-input>
        </search-item>
        <search-item text="所在部门：">
          <el-select
            class="search-select"
            :value="select_dp.name"
            placeholder="请选择所属部门"
            size="medium"
            clearable
            @change="changeSearch"
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
        </search-item>
        <el-button size="small" type="primary" @click="search"> 查询</el-button>
        <el-button class="mar-l-20 " size="small" @click="reset">
          重置</el-button
        >
      </template>

      <template slot="table">
        <add-guide
          v-if="show_modal"
          @fail="closeModal"
          @success="submitDialog"
        ></add-guide>
        <el-table :data="list" class="mar-t-20" v-if="page_init" stripe>
          <el-table-column label="导购" width="120" fixed="left">
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
            min-width="150"
          ></el-table-column>
          <el-table-column
            prop="dp_name"
            label="所在部门"
            min-width="170"
          ></el-table-column>
          <el-table-column
            prop="already_relation_num"
            label="已添加客户数"
            min-width="150"
            align="right"
          >
            <template slot="header" slot-scope="">
              <div class="popover-box">
                <span>已添加客户数</span>
                <el-popover
                  ref="popover1"
                  placement="right-end"
                  min-width="300"
                  trigger="hover"
                  visible-arrow="true"
                  content="导购添加客户的数量（包含导购删除客户数）"
                >
                  <svg-icon
                    slot="reference"
                    class="help-icon"
                    icon-class="group-help"
                  ></svg-icon>
                </el-popover>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="date" min-width="150">
            <template slot="header" slot-scope="">
              <div class="popover-box">
                <span>支付即导购开关</span>
                <el-popover
                  ref="popover2"
                  placement="right-end"
                  min-width="300"
                  trigger="hover"
                  visible-arrow="true"
                  content="开启后系统自动分配支付即导购，关闭后系统将不再自动分配支付即导购"
                >
                  <svg-icon
                    slot="reference"
                    class="help-icon"
                    icon-class="group-help"
                  ></svg-icon>
                </el-popover>
              </div>
            </template>

            <div slot-scope="scope">
              <el-switch
                v-model="scope.row.trade_guide_status"
                :active-value="1"
                :inactive-value="2"
                active-color="#13ce66"
                inactive-color="#F0F0F2"
                @change="tradeGuideStatusChange(scope.row, scope.$index)"
              ></el-switch>
            </div>
          </el-table-column>

          <el-table-column label="操作" width="70" fixed="right">
            <template slot-scope="scope">
              <el-button
                @click="deleteItem(scope.row)"
                type="text"
                class="operation-error"
                >移除</el-button
              >
            </template>
          </el-table-column>
          <div slot="empty">
            <layout-empty
              empty_type="tag"
              v-if="search_key"
              :message="'没有找到相关导购'"
            ></layout-empty>
            <layout-empty
              v-if="!search_key"
              :message="'暂无导购'"
            ></layout-empty>
          </div>
        </el-table>
      </template>

      <template slot="footer">
        <el-pagination
          v-if="list.length !== 0 && page_init"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="meta.page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="meta.size"
          class="page-right font-weight-400"
          layout="total, sizes, prev, pager, next, jumper"
          :total="meta.total"
        ></el-pagination>
      </template>
    </layout-list>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import Api from "@/api/wisdom/Guide";
import personApi from "@/api/wisdom/Personnel";
import AddGuide from "@/components/sy-components/add-guide.vue";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
// import Api from "@/api/organization/Guide";
@Component({
  components: {
    AddGuide,
    LayoutList,
    LayoutEmpty
  }
})
export default class extends Vue {
  private page_init = false;
  private text: any = "";
  private item: object = {
    date: "2016-05-03",
    name: "王小虎",
    province: "上海",
    city: "普陀区",
    address: "上海市普陀区金沙江路 1518 弄",
    zip: 200333
  };
  private meta = {
    total: 40, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private list: Array<any> = [];
  private type_list: Array<any> = [
    { name: "导购姓名", key: 0, keyword: "mem_user_name" },
    { name: "导购手机号", key: 1, keyword: "mem_mobile" }
  ];
  private gudie_type = 0;
  private guide_text = "";
  private search_active_type = "";
  private active_type_list = [
    { name: "企业管理", id: 4 },
    { name: "企业导购", id: 3 },
    { name: "社群导购", id: 2 },
    { name: "门店导购", id: 1 }
  ];
  private show_modal = false; // 是否显示待添加客户模态框
  private mem_uid = null; // 员工memUid
  private search_key = "";

  private dp_props = {
    label: "dp_name",
    children: [],
    isLeaf: "leaf"
  };
  private select_dp = { name: "", value: "" };
  private isShowTree = false;
  private dialog_type = 1;
  private remind_name = "";
  // private select_dp = { name: "", value: "" };
  private changed(val: any) {
    console.log(val);
  }
  // 生命周期
  private created() {
    this.fetchList();
    console.log("created");
  }
  private mounted() {
    console.log("mounted");
    for (let i = 0; i < 12; i++) {
      this.list.push(this.item);
    }
  }
  private destroyed() {
    console.log("destroyed");
  }
  private changeGuideType() {
    console.log(this.search_active_type);
  }
  private changeSearch() {
    this.select_dp.value = "";
    this.select_dp.name = "";
    this.search();
    // this.$refs.tree.setCurrentKey(null);
  }
  // 搜索
  private search() {
    this.meta.page = 1;
    this.fetchList();
  }
  //打开添加支付即导购弹框
  private addGuide() {
    this.show_modal = true;
  }
  // 选择人员完成
  private submitDialog(res: any) {
    if (!res || res.length <= 0) return;
    const params = { mem_uids: res.map((v: any) => v.mem_uid) };
    Api.guideManageAdd(params).then(() => {
      this.$message.success("添加成功");
      this.show_modal = false;
      this.fetchList();
    });
  }
  // 关闭弹框
  private closeModal() {
    this.show_modal = false;
  }
  // 更改支付即导购状态
  private tradeGuideStatusChange(params: any, index: any) {
    this.list[index].trade_guide_status =
      this.list[index].trade_guide_status === 1 ? 2 : 1;
    const title_txt = params.trade_guide_status == 1 ? "关闭" : "开启";
    const trade_guide_status = params.trade_guide_status === 1 ? 2 : 1;
    const cont_txt = `你确认要为该导购${title_txt}支付即导购吗?`;
    this.$confirm(`${cont_txt}`, `${title_txt}支付即导购确认`, {
      confirmButtonText: "确定",
      cancelButtonText: "取消"
    }).then(() => {
      Api.changeAllocationStatus({
        mem_uid: params.mem_uid,
        trade_guide_status
      }).then(
        () => {
          this.$message.success(
            `${params.trade_guide_status == 1 ? "关闭" : "开启"}成功!`
          );
          this.fetchList();
        },
        err => {
          if (err.data.status === 210 && err.data.code === 1400002) {
            this.$confirm(
              "你未开通支付即导购功能，请联系管理员配置开通",
              "提示",
              {
                confirmButtonText: "我知道了",
                showCancelButton: false
              }
              // eslint-disable-next-line @typescript-eslint/no-empty-function
            ).then(() => {});
          }
        }
      );
    });
  }
  // 移除导购
  private deleteItem(item: any) {
    this.$confirm("你确认要将该导购移除支付即导购吗？", "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消"
    }).then(() => {
      Api.guideManageDelete({ mem_uid: item.mem_uid }).then(() => {
        this.$message.success("移除成功!");
        // 删除时返回上一页判断
        if (this.list.length == 1 && this.meta.page != 1) {
          this.meta.page--;
        }
        this.fetchList();
      });
    });
  }
  // 获取导购列表
  private fetchList() {
    const params = {
      dp_id: this.select_dp.value,
      page: this.meta.page,
      size: this.meta.size
    };
    params[this.type_list[this.gudie_type].keyword] = this.search_key;
    Api.guidePayList(params).then(res => {
      const { list, page, size, total } = res.data;
      this.list = list;
      this.meta = { page, size, total };
      this.page_init = true;
    });
  }

  // 每页
  private handleSizeChange(val: any) {
    this.meta.page = 1;
    this.meta.size = val;
    this.fetchList();
  }
  // 当前页
  private handleCurrentChange(val: any) {
    this.meta.page = val;
    this.fetchList();
  }
  // 选择部门
  private handleNodeClick(node: any) {
    this.select_dp.name = node.dp_name;
    this.select_dp.value = node.dp_id;
    this.search();
  }
  // 搜索
  //  private search() {
  //     this.meta.page = 1;
  //     this.fetchList();
  //   }
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
  // private changeSearch() {
  //   this.select_dp.value = "";
  //   this.select_dp.name = "";
  //   this.search();
  //   // this.$refs.tree.setCurrentKey(null);
  // }
  // 重置
  private reset() {
    this.search_key = "";
    this.gudie_type = 0;
    this.select_dp.value = "";
    this.select_dp.name = "";
    this.search();
  }
}
</script>
<style lang="scss" scoped>
.mar-l-15 {
  margin-left: 15px;
}
.mar-l-30 {
  margin-left: 30px;
}
.mar-l-20 {
  margin-left: 20px;
}
.mar-t-20 {
  margin-top: 20px;
}
.mar-b-20 {
  margin-bottom: 20px;
}
.mar-r-15 {
  margin-right: 15px;
}
.mar-t-15 {
  margin-top: 15px;
}
.select-input-title {
  margin-right: 20px;
  display: inline-block;
  line-height: 14px;
  font-size: 14px;
  // font-weight: 700;
  font-style: Hiragino Sans GB;
}
.select-input {
  display: inline-block;
}

.search-input {
  width: 280px !important;
  margin-bottom: 0px !important;
}
.top-select-bc {
  background: #eeeeee;
  overflow: hidden;
}
.top-table {
  padding: 15px;
  background: #eeeeee;
  display: inline-block;
  width: 246px;
  margin-right: 20px;
  border-radius: 10px;
  .remind-info {
    margin: 15px 0;
    color: #666666;
    font-size: 12px;
  }
}
.popover-box {
  // margin-top: 11px
  // padding-left: 0
  line-height: 44px;
  .el-icon-question {
    margin-left: 5px;
    cursor: pointer;
    &:hover {
      color: #0076ff;
    }
  }
}

.help-icon {
  margin-left: 5px;
}

.search-item {
  margin-bottom: 0;
}

.treeStyle {
  .el-tree-node.is-current.is-focusable.is-checked {
    color: #409eff !important;
  }
}
.el-scrollbar {
  .el-scrollbar__view {
    .el-select-dropdown__item {
      height: auto;
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
.table-text-img {
  border-radius: 100%;
  vertical-align: middle;
  margin-right: 4px;
}
.search-select {
  width: 280px;
}
/deep/.layout-list .slide-box {
  display: none;
}
// /deep/.el-input-group__prepend div.el-select .el-input__inner {
//   width: 120px;
// }
/deep/.el-input__inner {
  padding-right: 0;
}
.input-with-select {
  width: 115px;
}
</style>
