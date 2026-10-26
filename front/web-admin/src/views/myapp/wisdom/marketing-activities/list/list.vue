<template>
  <div>
    <layout-list>
      <template slot="buttons">
        <el-button class="add-btn" @click="addNewActive">新增活动</el-button>
      </template>
      <template slot="search">
        <div>
          <search-item text="活动时间：">
            <el-date-picker
              v-model="active_time"
              type="datetimerange"
              range-separator="-"
              :default-time="['00:00:00', '23:59:59']"
              start-placeholder="开始时间"
              end-placeholder="结束时间"
            >
            </el-date-picker>
          </search-item>
          <search-item text="活动状态：">
            <cy-select
              v-model="active_status"
              size="medium"
              :optionsList="statusList"
              :options="{ label: 'label', value: 'value' }"
            >
            </cy-select>
          </search-item>
          <search-item text="活动类型：">
            <cy-select
              v-model="active_type"
              size="medium"
              :optionsList="activeList"
              :options="{ label: 'label', value: 'value' }"
            >
            </cy-select>
          </search-item>
          <search-item text="上下架状态：">
            <cy-select
              v-model="status"
              size="medium"
              :optionsList="upStatusList"
              :options="{ label: 'label', value: 'value' }"
            >
            </cy-select>
          </search-item>
          <el-button size="small" type="primary" @click="searchActive">
            查询
          </el-button>
          <el-button class="mar-l-20" size="small" @click="reset">
            重置
          </el-button>
        </div>
      </template>
      <template slot="table">
        <el-table stripe :data="list" class="mar-t-20">
          <el-table-column label="活动标题" width="220">
            <template slot-scope="scope">
              <span class="c-blue plan-id-box" @click="viewActivity(scope.row)">
                {{ scope.row.title }}</span
              >
            </template>
          </el-table-column>
          <el-table-column label="活动类型" width="120">
            <template slot-scope="scope">
              <span>
                {{
                  scope.row.type === 1
                    ? "报名活动"
                    : scope.row.type === 2
                    ? "签到活动"
                    : "问卷活动"
                }}</span
              >
            </template>
          </el-table-column>
          <el-table-column width="130" label="活动范围">
            <template slot-scope="scope">
              <span> {{ scope.row.scope === 1 ? "全部" : "指定范围" }}</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="activity_user_count"
            width="130"
            label="参与人数"
          ></el-table-column>
          <el-table-column prop="activit_time" label="活动时间" width="230">
          </el-table-column>
          <el-table-column
            prop="status"
            label="活动状态"
            width="120"
          ></el-table-column>
          <el-table-column prop="last_update" label="最后更新时间" width="230">
          </el-table-column>
          <el-table-column label="上下架状态" width="130">
            <template slot-scope="scope">
              <el-switch
                v-model="scope.row.isCheck"
                @change="updataStatu(scope.row)"
              >
              </el-switch>
            </template>
          </el-table-column>
          <el-table-column prop="zip" label="操作" width="240" fixed="right">
            <template slot-scope="scope">
              <span class="c-blue" @click="viewActivity(scope.row)"
                >查看详情
              </span>
              <span
                class="c-blue"
                @click="editActivity(scope.row)"
                v-if="scope.row.status == '未开始'"
                >编辑
              </span>
              <span class="c-blue" @click="copyActivity(scope.row)">复制 </span>
              <span
                class="c-blue"
                @click="deleteEvn(scope.row.id)"
                v-if="!scope.row.isCheck"
                >删除
              </span>
            </template>
          </el-table-column>
          <div slot="empty" class="table-empty-label">
            <img
              v-if="searchLoading"
              src="@/assets/images/empty-search.png"
              alt=""
            />
            <img v-else src="@/assets/images/empty-box.png" alt="" />
            <p>
              {{ searchLoading ? "没有找到相关营销活动" : "暂无营销活动" }}
            </p>
          </div>
        </el-table>
      </template>
      <template slot="footer">
        <el-pagination
          v-if="list && list.length > 0"
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
import Api from "@/api/wisdom/marketing-activities";
@Component({})
export default class extends Vue {
  private active_time: Array<any> = []; // 活动时间
  private statusList: Array<any> = [
    {
      label: "全部",
      value: 0,
    },
    {
      label: "未开始",
      value: 2,
    },
    {
      label: "进行中",
      value: 1,
    },
    {
      label: "已结束",
      value: 3,
    },
  ]; // 活动状态
  private activeList: Array<any> = [
    {
      label: "全部",
      value: 0,
    },
    {
      label: "报名活动",
      value: 1,
    },
    {
      label: "问卷活动",
      value: 3,
    },
    {
      label: "签到活动",
      value: 2,
    },
  ]; // 活动状态
  private upStatusList: Array<any> = [
    {
      label: "全部",
      value: 0,
    },
    {
      label: "已上架",
      value: 1,
    },
    {
      label: "已下架",
      value: 2,
    },
  ]; // 活动上下架状态
  private active_status = 0; // 活动状态
  private active_type = 0; // 活动类型
  private status = 0; // 活动上下架状态
  private searchLoading = false; // 是否正在搜索

  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10, // 每页数量
  };

  private list: Array<any> = [];
  // 生命周期
  private mounted() {
    this.getActiveList();
  }
  // 查询
  private searchActive() {
    this.meta.page = 1;
    // 是否正在搜索
    this.searchLoading = true;
    this.getActiveList();
  }
  // 重置
  private reset() {
    // 初始化数据
    this.active_time = [];
    this.active_status = 0; // 活动状态
    this.active_type = 0; // 活动类型
    this.status = 0; // 活动上下架状态
    this.searchLoading = false; // 是否正在搜索
    this.meta.page = 1;
    this.meta.size = 10;
    this.getActiveList();
  }
  private handleSizeChange(val: number) {
    this.meta.size = val;
    this.getActiveList();
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getActiveList();
  }
  private getActiveList() {
    const params: any = {
      page: this.meta.page,
      size: this.meta.size,
    };
    if (this.active_time && this.active_time[0]) {
      params.start_time = this.active_time[0] / 1000;
      params.end_time = this.active_time[1] / 1000;
    }
    if (this.active_status !== 0) {
      params.status = this.active_status;
    }
    if (this.active_type !== 0) {
      params.type = this.active_type;
    }
    if (this.status !== 0) {
      params.online = this.status;
    }
    Api.activityList(params)
      .then((res: any) => {
        if (res.httpCode == 200) {
          this.list = res.data.list;
          this.list.forEach((item: any) => {
            if (item.online === 1) {
              item.isCheck = true;
            } else {
              item.isCheck = false;
            }
          });
          this.meta.total = res.data.meta.total;
        }
      })
      .catch((err: any) => {
        console.log("getActiveList:", err);
      });
  }
  // 更新任务状态
  private updataStatu(row: any) {
    const { is_final } = row;
    if (is_final === 0) {
      this.$message.warning("草稿状态不能更新");
      return;
    }
    const params = {
      id: row.id,
    };
    Api.activityUpdataOnline(params)
      .then((res: any) => {
        if (res.httpCode == 200) {
          this.$message.success(res.status);
          this.getActiveList();
        }
      })
      .catch((err: any) => {
        console.log("activityUpdataOnline:", err);
      });
  }
  private deleteEvn(val: number | string) {
    this.$confirm("你确认删除该活动吗？", "删除确认", {
      confirmButtonText: "确认",
      cancelButtonText: "取消",
    })
      .then(() => {
        Api.activityDelete({ id: val })
          .then((res: any) => {
            if (res.httpCode == 200) {
              this.$message.success(res.status);
              this.getActiveList();
            }
          })
          .catch((err: any) => {
            console.log("taskCenterDelete:", err);
          });
      })
      .catch((err: any) => {
        console.log(err);
      });
  }
  // 新增活动
  private addNewActive() {
    this.$router.push({
      name: "myapp_wisdom_marketing-activities_add-activity",
      query: {
        status_type: "1",
      },
    });
  }
  // 查看活动详情
  private viewActivity(row: any) {
    this.$router.push({
      name: "myapp_wisdom_marketing-activities_activity-detail",
      query: { id: row.id },
    });
  }
  // 编辑活动
  private editActivity(row: any) {
    this.$router.push({
      name: "myapp_wisdom_marketing-activities_edit-activity",
      query: { id: row.id, status_type: "2" },
    });
  }
  // 复制活动
  private copyActivity(row: any) {
    this.$router.push({
      name: "myapp_wisdom_marketing-activities_copy-activity",
      query: { id: row.id, status_type: "3" },
    });
  }
}
</script>
<style lang="scss" scoped>
.search-input {
  width: 280px !important;
}
.input-with-select {
  width: 110px;
}
.add-btn {
  color: #fff;
  background: #708cf6;
}
.c-blue {
  cursor: pointer;
  color: #708cf6;
  margin-right: 8px;
}
.table-empty-label {
  text-align: center;
  height: 340px;
  img {
    width: 124px;
    transform: translateX(-230%) translateY(90px);
  }
  p {
    transform: translateX(-36%) translateY(55px);
    line-height: 20px;
    margin-top: 14px;
  }
}
</style>
