<template>
  <div class="contaniner">
    <el-dialog
      class="dialog-box"
      width="800px"
      min-height="900px"
      title="导购关联门店"
      :visible.sync="dialog_visible"
      :close-on-click-modal="false"
      @close="closeDialog"
    >
      <div class="dialog-body">
        <div class="staff-list">
          <div class="title">已选择员工列表</div>
          <div class="list">
            <div
              class="item-list"
              v-for="(item, index) in staff_list"
              :key="index"
            >
              <div class="name">
                【{{ item.employee_code }}】{{ item.name }}
              </div>
              <div class="del" @click="delStaff(index)">
                <i class="el-icon-delete del-icon"></i>
              </div>
            </div>
          </div>
        </div>
        <div class="store-tree">
          <div class="title">选择需要关联的门店</div>
          <div class="search">
            <el-input
              clearable
              size="small"
              placeholder="请输入门店名称或编码"
              v-model="keyword"
            ></el-input>
            <div class="btns">
              <el-button class="mar-l-30" size="mini" @click="search()"
                >搜索</el-button
              >
              <el-button size="mini" type="primary" @click="resetList()"
                >重置</el-button
              >
            </div>
          </div>
          <el-table
            tooltip-effect="light"
            :data="store_list"
            style="width:95%;"
            max-height="300"
            :show-header="false"
            @row-click="rowClick"
          >
            <el-table-column label="" width="40">
              <template slot-scope="{ row }">
                <el-radio v-model="select_id" :label="row.id"></el-radio>
              </template>
            </el-table-column>
            <el-table-column prop="store_name" label="门店名称">
              <template slot-scope="scope">
                {{ scope.row.store_name }}【{{ scope.row.third_store_code }}】
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            class="foot-pages"
            @current-change="handleCurrentChange"
            :current-page.sync="meta.page"
            :page-size="meta.size"
            layout="total, prev, pager, next"
            :total="meta.total"
          >
          </el-pagination>
        </div>
      </div>
      <div slot="footer" class="dialog-footer">
        <el-button @click="closeDialog" size="small">
          取消
        </el-button>
        <el-button @click="sureData" size="small" type="primary">
          确 认
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from "vue-property-decorator";
import storeApi from "@/api/enterprise/business";

@Component
export default class extends Vue {
  @Prop({ default: () => [] }) selectList?: Array<any>; // 选择的数据

  /********************** data **********************/
  private dialog_visible = true; // 控制模态框显示/隐藏
  private staff_list = []; // 选择的员工列表
  private store_list = []; // 可选的部门列表
  private current_store = {
    id: 0
  }; // 选中的门店
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private select_id = "";
  private keyword = "";
  /********************** 方法 ***************************/
  private rowClick(row: any) {
    this.select_id = row.id;
  }
  // 关闭模态框的回调
  @Emit("close")
  private closeDialog() {
    this.$emit("close");
  }
  private errorText() {
    if (this.staff_list.length === 0) return "请选择员工";
    if (!this.select_id) return "请选择关联的门店";
    return "";
  }
  private resetList() {
    this.keyword = "";
    this.search();
  }
  private search() {
    this.meta.page = 1;
    this.getStoreList();
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getStoreList();
  }
  private getStoreList() {
    const params = {
      store_type: "",
      code_or_name: this.keyword,
      page: this.meta.page,
      size: this.meta.size
    };
    storeApi.getStoreListApi(params).then(res => {
      const { total, page, size, list } = res.data;
      this.store_list = list;
      this.meta = {
        total,
        page,
        size
      };
    });
  }
  private sureData() {
    const error_text = this.errorText();
    if (error_text) {
      this.$message.warning(error_text);
      return;
    }
    const params = {
      department_id: this.select_id,
      employee_id_arr: this.staff_list.map(item => item.id)
    };
    storeApi.staffMoveStoreApi(params).then(() => {
      this.$message.success("分配成功");
      this.$emit("success");
    });
  }
  private delStaff(index) {
    this.staff_list.splice(index, 1);
  }
  private async mounted() {
    this.staff_list = JSON.parse(JSON.stringify(this.selectList));
    this.getStoreList();
  }
}
</script>

<style lang="scss" scoped>
.dialog-body {
  height: 400px;
  display: flex;
  justify-content: flex-start;
  .staff-list {
    width: 300px;
    height: 400px;
    border-right: 1px solid #ebebeb;
    flex-shrink: 0;
    overflow-y: scroll;
  }
  .store-tree {
    padding-left: 30px;
    padding-bottom: 50px;
    flex-basis: 0;
    flex: 1 1 auto;
    position: relative;
    .foot-pages {
      position: absolute;
      bottom: 0;
      right: 0;
    }
    /deep/ .el-tree {
      padding-top: 10px !important;
    }
    .store-icon {
      padding-right: 10px;
    }
    .custom-tree-node {
      width: 100%;
      padding-right: 40px;
      position: relative;
      .select-store {
        position: absolute;
        top: 0;
        right: 10px;
        font-size: 16px;
        font-weight: bolder;
      }
    }
  }
  .title {
    font-size: 16px;
    color: #606266;
    margin-bottom: 10px;
  }
  /deep/.el-table td {
    padding: 10px 0 !important;
  }
  /deep/.el-table::before {
    display: none;
  }
  .item-list {
    display: flex;
    justify-content: space-between;
    align-items: center;
    line-height: 20px;
    &:hover {
      background: #ebebeb;
      .del {
        color: $color-text-button;
      }
    }
    .name {
      padding: 10px 0;
    }
    .del {
      cursor: pointer;
      padding: 10px 20px;
    }
  }
  .search {
    display: flex;
    justify-content: space-between;
    width: 100%;
    margin-bottom: 5px;
    .btns {
      margin-left: 50px;
      display: flex;
      justify-content: flex-start;
      flex-shrink: 0;
      /deep/.el-button--mini {
        padding: 10px 15px !important;
      }
    }
  }
}
</style>
