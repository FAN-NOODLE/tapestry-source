<template>
  <div class="research-list">
    <el-tabs
      v-model="activeName"
      type="card"
    >
      <el-tab-pane
        label="全部"
        name="first"
      >
        <div class="table-content">
          <el-table
            :data="list"
            stripe
            tooltip-effect="light"
          >
            <el-table-column
              prop="name"
              label="活动标题"
              width="240"
            >
            </el-table-column>
            <el-table-column
              label="活动类型"
              width="120"
            >
              <template slot-scope="scope">
                {{ scope.row.type === 1 ? "报名活动" : "签到活动" }}
              </template>
            </el-table-column>
            <el-table-column
              prop="nums"
              label="参与人数"
              width="120"
            >
            </el-table-column>
            <el-table-column
              label="活动时间"
              min-width="260"
            >
              <template slot-scope="scope">
                <span>
                  {{ scope.row.start_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
                  ~
                  {{ scope.row.end_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
                </span>
              </template>
            </el-table-column>
            <el-table-column
              prop="status_text"
              label="状态"
              width="100"
            >
            </el-table-column>
            <el-table-column
              label="操作"
              width="240"
              fixed="right"
            >
              <template slot-scope="scope">
                <el-button
                  type="text"
                  size="small"
                  class="marginR20"
                  @click="check(scope.row)"
                >查看</el-button>
                <el-button
                  type="text"
                  size="small"
                  class="marginR20 disabled"
                >编辑</el-button>
                <el-button
                  type="text"
                  size="small"
                  class="marginR20 disabled"
                >删除</el-button>
              </template>
            </el-table-column>
            <div
              slot="empty"
              class="table-empty-label-page"
            >
              <img
                class="empty-img"
                :src="require('@/assets/images/empty-box.png')"
                alt=""
              />
              <div>暂无调研活动</div>
            </div>
          </el-table>
          <div class="table_footer">
            <el-pagination
              v-if="list && list.length > 0"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
              :current-page="meta.page"
              :page-sizes="[10, 20, 30, 40]"
              :page-size="meta.size"
              layout="total, sizes, prev, pager, next, jumper"
              :total="meta.total"
            >
            </el-pagination>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/wisdom/researchActivities";
@Component
export default class extends Vue {
  // data
  private activeName = "first";
  private list: Array<any> = [];
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };

  private mounted() {
    this.getList();
  }

  // 查看
  private check(row: any) {
    const { id, type } = row;
    if (type == 1) {
      this.$router.push({
        name: "myapp_wisdom_research-activities_events-list",
        query: { id: id }
      });
    } else if (type == 2) {
      this.$router.push({
        name: "myapp_wisdom_research-activities_signIn-list",
        query: { id: id }
      });
    }
  }

  // 每页
  private handleSizeChange(val: any) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getList();
  }

  // 当前页
  private handleCurrentChange(val: any) {
    this.meta.page = val;
    this.getList();
  }

  private getList() {
    const params = {
      page: this.meta.page,
      size: this.meta.size
    };
    Api.getTypeList(params).then((res: any) => {
      if (res.code === 200) {
        this.list = res.data.list;
        this.meta.total = res.data.total;
      }
    });
  }
}
</script>
<style lang="scss" scoped>.research-list {
  background-color: #fff;
  padding: 20px;
  border-radius: 10px;
}

/deep/ .el-tabs__nav {
  .is-active {
    color: #fff;
    background-color: #708df5;
  }
}

.table-empty-label-page {
  display: flex;
  height: 340px;
  flex-direction: column;

  img {
    width: 124px;
    margin: 80px auto 0;
  }

  div {
    line-height: 20px;
    margin-top: 14px;
  }
}

.table_footer {
  margin: 20px 0;
  text-align: right;
}

.disabled {
  color: #d7d7d7;
  cursor: no-drop;
}</style>
