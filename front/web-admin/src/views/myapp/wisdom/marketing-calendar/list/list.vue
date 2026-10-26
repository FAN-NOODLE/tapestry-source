<template>
  <div class="list">
    <div
      v-if="showDeployTip"
      class="content-wrap m-b-30"
    >
      <div class="align-item-center">
        <img
          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/customer/ico_tip.png"
          alt=""
        />
        营销日历可在企业微信聊天工具栏(侧边栏)使用，
        <span
          class="c-blue"
          @click="checkDeploy"
        >点击查看如何配置</span>
      </div>
      <div @click="closeDeploy"><i class="el-icon-close"></i></div>
    </div>

    <layout-list>
      <template slot="buttons">
        <el-button
          class="add-btn"
          @click="addMarketingPlan"
        >新增营销计划</el-button>
      </template>
      <template slot="search">
        <div>
          <search-item text="计划搜索：">
            <el-input
              :placeholder="`请输入${type_list[gudie_type].name}`"
              class="search-input"
              size="medium"
              v-model.trim="search_key"
              clearable
            >
              <el-select
                slot="prepend"
                v-model="gudie_type"
                class="input-with-select"
              >
                <el-option
                  v-for="item in type_list"
                  :key="item.key"
                  :value="item.key"
                  :label="item.name"
                ></el-option>
              </el-select>
            </el-input>
          </search-item>

          <search-item text="创建时间：">
            <el-date-picker
              v-model="createTime"
              type="daterange"
              range-separator="-"
              start-placeholder="开始时间"
              end-placeholder="结束时间"
            >
            </el-date-picker>
          </search-item>

          <search-item text="编辑时间：">
            <el-date-picker
              v-model="editTime"
              type="daterange"
              range-separator="-"
              start-placeholder="开始时间"
              end-placeholder="结束时间"
            >
            </el-date-picker>
          </search-item>
          <search-item text="创建人：">
            <cy-input
              size="small"
              width-size="small"
              placeholder="请输入创建人名称"
              v-model.trim="creatName"
              prefix-icon="el-icon-search"
              clearable
            ></cy-input>
          </search-item>
          <div>
            <el-button
              class=""
              size="small"
              type="primary"
              @click="searchEvn"
            >
              查询
            </el-button>
            <el-button
              class="mar-l-20"
              size="small"
              @click="reset"
            >
              重置
            </el-button>
          </div>
        </div>
      </template>
      <template slot="table">
        <el-table
          stripe
          :data="list"
          class="mar-t-20"
          v-if="page_init"
        >
          <el-table-column
            label="计划ID"
            width="150"
            fixed="left"
          >
            <template slot-scope="scope">
              <span
                class="c-blue plan-id-box"
                @click="checkPlanDetailEvn(scope.row.plan_id)"
              >{{ scope.row.plan_sn }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="计划名称"
            min-width="220"
          >
            <template slot-scope="scope">
              <el-popover
                placement="top-start"
                width="220"
                trigger="hover"
                :disabled="scope.row.plan_title.length <= 20"
                :content="scope.row.plan_title"
              >
                <div
                  :class="
                    scope.row.plan_title.length > 20 ? 'work-name-wrap' : ''
                  "
                  slot="reference"
                >
                  {{ scope.row.plan_title }}
                </div>
              </el-popover>
            </template>
          </el-table-column>
          <!-- 可见范围 -->
          <el-table-column
            label="可见范围"
            width="250"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.object_sum <= 1">
                {{
                  scope.row.object_array.length > 0
                    ? scope.row.object_array[0].name
                    : ""
                }}
              </span>
              <el-popover
                v-if="scope.row.object_sum > 1"
                popper-class="marketing-calendar-popover-wrap"
                placement="top-start"
                width="220"
                trigger="hover"
              >
                <div class="object-array-content">
                  {{ scope.row.object_array_item }}
                </div>
                <div slot="reference">
                  <span>
                    {{
                      scope.row.object_array.length > 0
                        ? scope.row.object_array[0].name
                        : ""
                    }}等{{ scope.row.object_sum }}个导购
                  </span>
                </div>
              </el-popover>
            </template>
          </el-table-column>
          <el-table-column
            label="营销任务数"
            min-width="120"
          >
            <template slot-scope="scope"><span>{{ scope.row.task_sum }}个</span></template>
          </el-table-column>
          <el-table-column
            prop="plan_resource"
            label="计划包含素材"
            min-width="200"
          >
          </el-table-column>
          <el-table-column
            label="创建时间/编辑时间"
            width="220"
          >
            <template slot-scope="scope"><span>{{ scope.row.created_at | formatDate("yyyy-mm-dd") }}/{{
                  scope.row.updated_at | formatDate("yyyy-mm-dd")
                }}</span></template>
          </el-table-column>
          <el-table-column
            prop="operator_name"
            label="创建人"
            min-width="250"
          ></el-table-column>
          <el-table-column
            prop="zip"
            label="操作"
            width="200"
            fixed="right"
          >
            <template slot-scope="scope">
              <span
                class="c-blue operate-btn"
                @click="editEvn(scope.row.plan_id)"
              >
                编辑
              </span>
              <span
                class="c-blue operate-btn"
                @click="copyEvn(scope.row.plan_id)"
              >复制
              </span>
              <span
                class="c-red operate-btn"
                @click="deleteEvn(scope.row.plan_id)"
              >删除
              </span>
            </template>
          </el-table-column>

          <div slot="empty">
            <layout-empty
              empty_type="tag"
              v-if="search_key || createTime || editTime || creatName"
              :message="'没有找到相关营销计划'"
            ></layout-empty>
            <layout-empty
              v-if="!search_key && !createTime && !editTime && !creatName"
              :message="'暂无营销计划'"
            ></layout-empty>
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
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
import {
  IMarketingPlanListParams,
  marketingPlanList,
  IMarketingPlanDeleteParams,
  marketingPlanDelete,
} from "@/api/wisdom/marketingPlan";
@Component({
  components: {
    LayoutEmpty,
  },
})
export default class extends Vue {
  private showDeployTip = true; // 营销日历配置提示展示flag
  // 计划搜索选项 TODO新组件
  private type_list: Array<any> = [
    { name: "计划ID", key: 0, keyword: "worker_no" },
    { name: "计划名称", key: 1, keyword: "mem_name" },
  ];
  private gudie_type = 1; // 计划搜索选项 导购ID/导购名称
  private search_key = ""; // 搜索val

  private editTime = ""; // 编辑时间
  private createTime = ""; // 创建时间
  private creatName = ""; // 创建人名称
  private list: Array<any> = [
    // {
    //   plan_sn: "JH20201030001", // 计划ID
    //   plan_title: "这里显示计划名称最多显示20个字", // 计划名称
    //   mem_mobile: "西部国际广场店等4个门店",
    //   task_sum: "4个", // 任务数量
    //   plan_resource: "图片", // 计划包含素材
    //   created_at: "", // 创建时间
    //   updated_at: "", //编辑时间
    //   operator_name: "创建人", // 创建人
    //   object_sum: 1, // 可见范围数量
    //   object_type: 1, // 可见范围类型 1门店 2部门
    //   object_array: [{ worker_id: "11111", name: "部门或门店名称" }], // 部门或门店ID 部门或门店名称
    //   object_array_item: "部门或门店名称"
    // }
  ]; // 表格数据
  // 页面初始化设置
  private page_init = false;
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10, // 每页数量
  };
  /* 生命周期 */
  private mounted() {
    // 是否隐藏wiki提示
    this.showDeployTip = localStorage.showCalendarDeployTip
      ? JSON.parse(localStorage.showCalendarDeployTip)
      : true;
    this.getMarketingPlanList();
  }
  // 新增营销计划
  private addMarketingPlan() {
    this.$router.push({
      name: "myapp_wisdom_marketing-calendar_add-marketing-plan",
      query: {
        status: "0",
      },
    });
  }
  // 查询
  private searchEvn() {
    this.meta.page = 1;
    this.getMarketingPlanList();
  }
  // 重置
  private reset() {
    // 初始化数据
    this.gudie_type = 1;
    this.search_key = "";
    this.editTime = "";
    this.createTime = "";
    this.creatName = "";
    this.meta.page = 1;
    this.meta.size = 10;
    this.getMarketingPlanList();
  }
  private handleSizeChange(val: number) {
    this.meta.size = val;
    this.getMarketingPlanList();
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getMarketingPlanList();
  }
  // 编辑
  private editEvn(val: any) {
    this.$router.push({
      name: "myapp_wisdom_marketing-calendar_edit-marketing-plan",
      query: {
        status: "1",
        val,
      },
    });
  }
  // 复制
  private copyEvn(val: any) {
    this.$router.push({
      name: "myapp_wisdom_marketing-calendar_copy-marketing-plan",
      query: {
        status: "2",
        val,
      },
    });
  }
  // 删除
  private deleteEvn(val: number | string) {
    this.$confirm("你确认删除该营销计划吗？", "删除确认", {
      confirmButtonText: "确认",
      cancelButtonText: "取消",
    })
      .then(() => {
        this.deletMarketingPlan(val);
      })
      .catch((err: any) => {
        console.log(err);
      });
  }
  // 查看计划
  private checkPlanDetailEvn(val: any) {
    this.$router.push({
      name: "myapp_wisdom_marketing-calendar_marketing-plan-detail",
      query: {
        status: "3",
        val,
      },
    });
  }
  // 查看配置营销日历
  private checkDeploy() {
    const protocol = "http";
    window.open(
      `${protocol}://xueyuan.vchangyi.com/bin/view/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD%E7%AE%A1%E7%90%86%E5%90%8E%E5%8F%B0/%E8%90%A5%E9%94%80%E6%97%A5%E5%8E%86/%E5%B0%86%E8%90%A5%E9%94%80%E6%97%A5%E5%8E%86%E9%85%8D%E7%BD%AE%E5%88%B0%E4%BC%81%E5%BE%AE%E4%BE%A7%E8%BE%B9%E6%A0%8F/`
    );
  }
  // 配置营销日历提示关闭
  private closeDeploy() {
    this.showDeployTip = false;
    localStorage.showCalendarDeployTip = this.showDeployTip;
  }
  // 获取营销日历列表
  private getMarketingPlanList() {
    const param: IMarketingPlanListParams = {
      page: this.meta.page,
      size: this.meta.size, //条数：默认10
      plan_title: this.gudie_type == 1 ? this.search_key : "", // 计划名称
      plan_sn: this.gudie_type == 0 ? this.search_key : "", //计划ID
      create_time: {
        create_start: this.createTime
          ? Date.parse(this.createTime[0]) / 1000
          : "",
        create_end: this.createTime ? Date.parse(this.createTime[1]) / 1e3 : "",
      },
      //  TODO 时间选择组件可以自定义返回时间格式
      edit_time: {
        edit_start: this.editTime ? Date.parse(this.editTime[0]) / 1000 : "",
        edit_end: this.editTime ? Date.parse(this.editTime[1]) / 1e3 : "",
      },
      operator_name: this.creatName, // 创建人
    };
    marketingPlanList(param).then((res: any) => {
      if (res.code === 200) {
        this.page_init = true;
        // 处理门店/部门名称
        res.data.list.forEach((item: any) => {
          // 处理导购名称
          if (item.object_array && item.object_array.length > 0) {
            item.object_array_item = item.object_array
              .map((item: any) => {
                return item.name;
              })
              .join("、");
          }
          // 素材
          let resource_text = "";
          if (item.plan_resource && item.plan_resource.length <= 0) {
            resource_text = "无素材";
          } else {
            resource_text = "";
            item.plan_resource.forEach(
              (planItem: number | string, planIndex: number) => {
                const resource =
                  planItem == 1 ? "图片" : planItem == 2 ? "小程序" : "视频";
                if (planIndex < item.plan_resource.length - 1) {
                  resource_text += resource + "、";
                } else {
                  resource_text += resource;
                }
                return resource_text;
              }
            );
          }
          item.plan_resource = resource_text;
        });
        this.list = res.data.list;
        this.meta.total = res.data.total;
      }
    });
  }
  // 删除营销日历
  private deletMarketingPlan(id: number | string) {
    const params: IMarketingPlanDeleteParams = {
      plan_id: id, //计划id
    };
    marketingPlanDelete(params)
      .then((res: any) => {
        if (res.code == 200) {
          this.$message.success("删除成功");
          this.getMarketingPlanList();
        }
      })
      .catch((err: any) => {
        console.log("marketingPlanDelete:", err);
      });
  }
}
</script>
<style lang="scss">
.marketing-calendar-popover-wrap {
  & .object-array-content {
    max-height: 150px;
    overflow: auto;
  }
}
</style>
<style lang="scss" scoped>
#crs-isg .el-date-editor {
  width: 280px !important;
}
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
.work-name-wrap {
  display: block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
.content-wrap {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background: #fff;
  border-radius: 4px;
  i {
    cursor: pointer;
    font-weight: bold;
    font-size: 18px;
    color: #c0c4cc;
  }
  span {
    cursor: pointer;
  }
  img {
    width: 14px;
    height: 14px;
    margin-right: 10px;
  }
}
.align-item-center {
  display: flex;
  align-items: center;
}
.operate-btn {
  cursor: pointer;
  margin-right: 20px;
}
.plan-id-box {
  cursor: pointer;
}

.c-blue {
  color: #708cf6;
}
.c-red {
  color: #ff4343;
}
.m-t-20 {
  margin-top: 20px;
}
.m-b-30 {
  margin-bottom: 30px;
}
</style>
