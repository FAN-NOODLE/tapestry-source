<template>
  <div>
    <el-dialog
      title="标签成员"
      :visible="true"
      @close="closeDialogBA"
      :close-on-click-modal="false"
    >
      <layout-list>
        <template slot="search">
          <search-item text="顾客：">
            <cy-input
              size="small"
              width-size="small"
              placeholder="请输入顾客昵称"
              v-model.trim="customer_name"
              prefix-icon="el-icon-search"
              clearable
            ></cy-input>
          </search-item>
          <el-button size="small" type="primary" @click="searchEvn">
            查询
          </el-button>
          <el-button class="mar-l-20" size="small" @click="reset">
            重置
          </el-button>
        </template>
        <template slot="table">
          <el-table stripe :data="list" class="mar-t-20">
            <el-table-column label="顾客" min-width="160">
              <div slot-scope="scope">
                <img
                  :src="
                    scope.row.avatar
                      ? scope.row.avatar
                      : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                  "
                  alt=""
                  width="30"
                  height="30"
                  class="table-text-img"
                /><span v-text="scope.row.name"></span>
              </div>
            </el-table-column>
            <el-table-column prop="tag_name" label="标签"></el-table-column>
            <el-table-column prop="sex" label="性别"></el-table-column>
            <el-table-column prop="mobile" label="手机号"></el-table-column>
            <div slot="empty" class="table-empty-label">
              <img src="@/assets/images/empty-box.png" alt="" />
              <p>暂无数据</p>
            </div>
          </el-table>
        </template>
        <template slot="footer">
          <el-pagination
            v-if="list && list.length > 0 && page_init"
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
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import Api from "@/api/wisdom/planTask";
@Component
export default class extends Vue {
  @Prop({ default: () => 0 }) selectedTerm?: number;
  @Prop({ default: () => "" }) selectedBA?: string;
  private mounted() {
    this.searchEvn();
  }
  /************data**************/
  private customer_name = ""; //顾客姓名
  private page_init = false;
  private list = [];
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private closeDialogBA() {
    this.$emit("closeDialogBA");
  }
  private reset() {
    this.customer_name = "";
    this.searchEvn();
  }
  private searchEvn() {
    const params = {
      id: this.selectedTerm,
      user_id: this.selectedBA,
      customer_name: this.customer_name,
      page: this.meta.page,
      size: this.meta.size
    };
    Api.userCompletion(params).then((res: any) => {
      if (res.code === 200) {
        this.page_init = true;
        this.list = res.data.list;
        this.meta.total = res.data.total;
      }
    });
  }
  private handleSizeChange(val: number) {
    this.meta.size = val;
    this.searchEvn();
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.searchEvn();
  }
}
</script>
<style lang="scss" scoped>
.table-text-img {
  width: 30px;
  height: 30px;
  border-radius: 100%;
  vertical-align: middle;
  margin-right: 4px;
}
/deep/ .layout-list {
  padding: 0;
}
/deep/ .search-content {
  padding: 10px !important;
  margin-bottom: 0 !important;
  background-color: #fff !important;
}
/deep/ .el-table__body-wrapper {
  max-height: 300px;
  overflow: auto;
}
</style>
