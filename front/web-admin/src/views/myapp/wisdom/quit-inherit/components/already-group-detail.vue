<template>
  <!-- 已分配社群情况 -->
  <div class="already-detail">
    <el-table
      :data="qun_data"
      height="400px"
      stripe
    >
      <el-table-column
        prop="takeover_department_name"
        label="所属门店"
      >
      </el-table-column>
      <el-table-column
        prop="takeover_name"
        label="承接导购"
      > </el-table-column>
      <el-table-column
        prop="name"
        label="承接群"
        show-overflow-tooltip
      >
      </el-table-column>
      <el-table-column
        prop="takeover_time"
        label="接替时间"
        width="190"
      >
        <template slot-scope="scope">
          <span>{{
            scope.row.takeover_time | formatDate("yyyy-mm-dd hh:mm:ss")
          }}</span>
        </template>
      </el-table-column>
      <el-table-column
        prop="status"
        label="状态"
        width="120"
      >
        <template slot-scope="scope">
          {{ scope.row.status == 1 ? "已完成" : "未完成" }}
        </template>
      </el-table-column>
      <div
        slot="empty"
        v-if="!qun_data.length"
      >
        <layout-empty :message="'暂无已分配社群'"></layout-empty>
      </div>
    </el-table>
    <div class="page-view">
      <el-pagination
        layout="prev, pager, next"
        :total="meta.total"
        :current-page="meta.page"
        @current-change="handleCurrentChange"
      >
      </el-pagination>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import Api from "@/api/wisdom/quit-inherit";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
@Component({
  components: {
    LayoutEmpty,
  },
})
export default class extends Vue {
  @Prop({ default: "" }) handoverUserid!: string;
  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 0,
  };
  private qun_data: Array<any> = [];
  //内部绑定函数（获取当前第几页）
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.takeoverGroupList();
  }
  private takeoverGroupList() {
    Api.takeoverGroupList({
      page: this.meta.page,
      size: this.meta.size,
      search_name: "",
      handover_userid: this.handoverUserid,
    }).then((res: any) => {
      if (res.code === 200) {
        this.qun_data = res.data.list || [];
        this.meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 生命周期
  mounted() {
    this.takeoverGroupList();
  }
}
</script>
<style lang="scss" scoped>
.already-detail {
  .page-view {
    padding-top: 10px;
    text-align: center;
  }
}
/deep/.el-table--scrollable-y .el-table__body-wrapper {
  /*修改滚动条样式*/
  &::-webkit-scrollbar {
    width: 5px;
    height: 10px;
    /**/
  }
  &::-webkit-scrollbar-track {
    background: rgb(239, 239, 239);
    border-radius: 2px;
  }
  &::-webkit-scrollbar-thumb {
    background: #bfbfbf;
    border-radius: 5px;
  }
  &::-webkit-scrollbar-thumb:hover {
    background: #333;
  }
  &::-webkit-scrollbar-corner {
    background: #179a16;
  }
}
</style>
