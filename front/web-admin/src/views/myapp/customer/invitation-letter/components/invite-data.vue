<template>
  <div class="subscribe-list">
    <div class="search">
      <search-item text="模板名称：">
        <cy-input
          width-size="small"
          size="small"
          prefix-icon="el-icon-search"
          placeholder="请输入模板名称"
          clearable
          v-model="search.title_name"
        ></cy-input>
      </search-item>
      <search-item text="关联服务预约：">
        <cy-select
          search-input
          width-size="small"
          size="small"
          v-model="search.reserve_service_id"
          placeholder="请选择服务预约"
          :options-list="associateServiceList"
          :options="{ label: 'reserve_name', value: 'id' }"
        >
        </cy-select>
      </search-item>
      <search-item text="邀约导购：">
        <cy-input
          width-size="small"
          size="small"
          prefix-icon="el-icon-search"
          placeholder="请输入导购姓名"
          clearable
          v-model="search.member_name"
        ></cy-input>
      </search-item>
      <search-item text="顾客姓名：">
        <cy-input
          width-size="small"
          size="small"
          prefix-icon="el-icon-search"
          placeholder="请输入顾客姓名"
          clearable
          v-model="search.user_name"
        ></cy-input>
      </search-item>
      <search-item text="所在部门：">
        <el-select
          ref="selectBlur"
          class="search-select"
          :value="search.storeName"
          placeholder="请选择所属部门"
          size="medium"
          clearable
          @change="changeStore"
        >
          <el-option :value="search.storeName" :label="search.storeName">
            <el-tree
              :props="dp_props"
              :load="loadNode"
              node-key="dp_id"
              :expand-on-click-node="true"
              :check-strictly="true"
              :check-on-click-node="true"
              @node-click="handleNodeClick"
              lazy
              @current-change="blurFn"
            ></el-tree>
          </el-option>
        </el-select>
      </search-item>
      <search-item text="到店时间： ">
        <el-date-picker
          v-model="search.arrive_time"
          type="datetimerange"
          prefix-icon="el-icon-date"
          format="yyyy-MM-dd HH:mm:ss"
          value-format="timestamp"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          :default-time="['00:00:00', '23:59:59']"
          style="height:36px"
        >
        </el-date-picker>
      </search-item>
      <search-item text="创建时间： ">
        <el-date-picker
          v-model="search.created_time"
          type="datetimerange"
          prefix-icon="el-icon-date"
          format="yyyy-MM-dd HH:mm:ss"
          value-format="timestamp"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          :default-time="['00:00:00', '23:59:59']"
          style="height:36px"
        >
        </el-date-picker>
      </search-item>
      <search-item>
        <el-button type="primary" @click="goSearch" size="small"
          >查询</el-button
        >
        <el-button plain @click="reset" size="small">重置</el-button>
        <el-button plain @click="exportList" size="small">导出</el-button>
      </search-item>
    </div>
    <el-table tooltip-effect="light" :data="list" ref="multipleTable">
      <el-table-column label="模板ID" prop="id" min-width="180">
      </el-table-column>
      <el-table-column label="模板名称" prop="title_name" min-width="180">
      </el-table-column>
      <el-table-column
        label="关联预约服务"
        prop="reserve_service_name"
        min-width="170"
      >
      </el-table-column>
      <el-table-column label="邀约导购" prop="member_name" min-width="100">
      </el-table-column>
      <el-table-column label="所在部门" prop="store_name" min-width="100">
      </el-table-column>
      <el-table-column label="顾客姓名" prop="user_name" min-width="100">
      </el-table-column>
      <el-table-column label="到店时间" prop="reserve_time_str" min-width="240">
        <template slot-scope="scope">
          <span>{{ scope.row.reserve_time_str || "--" }}</span>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" prop="create_time_str" min-width="200">
      </el-table-column>
      <div slot="empty" v-if="!emptyFlag" style="margin-top: 50px;">
        <img class="empty-img" src="@/assets/images/empty-box.png" alt="" />
        <div class="empty-box">
          <div class="desc">暂无数据</div>
        </div>
      </div>
      <div slot="empty" v-if="emptyFlag" style="margin-top: 50px;">
        <img class="empty-img" src="@/assets/images/empty-search.png" alt="" />
        <div class="empty-box">
          <div class="desc">没有找到相关数据</div>
        </div>
      </div>
    </el-table>
    <div class="table_footer" v-if="list.length > 0">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page.sync="meta.page"
        :page-sizes="[10, 20, 30, 40]"
        :page-size="meta.size"
        layout="total, sizes, prev, pager, next, jumper"
        :total="meta.total"
      >
      </el-pagination>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import API from "@/api/customer/invitation-letter";
import storeAPI from "@/api/enterprise/enterprise";

import { downFilePost } from "@/utils/tools";
@Component({})
export default class extends Vue {
  @Prop({ default: { id: "", title_name: "" } }) emitData?: any;

  private search = {
    title_name: "",
    reserve_service_id: "",
    member_name: "",
    user_name: "",
    arrive_time: [],
    created_time: [],
    third_id: "", //所在门店third_id
    storeName: "" //所在门店名称
  };

  private meta = {
    page: 1,
    size: 10,
    total: 5
  };

  private list = [];

  private associateServiceList = [];

  private emptyFlag = false;

  /** 生命周期 */
  private created() {
    if (Object.keys(this.emitData).length) {
      this.search.title_name = this.emitData.title_name;
    }
  }

  async mounted() {
    await this.getAssociateServiceList();
    this.getList();
  }

  /** 获取关联预约服务列表 */
  async getAssociateServiceList() {
    const res: any = await API.associateServiceList({
      page: 1,
      size: 999,
      status: 1
    });
    if (res.code == 200) {
      this.associateServiceList = res.data.list || [];
    } else {
      this.$message.error(res.message);
    }
  }

  /** 列表数据 */
  private getList() {
    const params = this.getParams();
    API.invitationDataList(params).then((res: any) => {
      if (res.code === 200) {
        const { total, list } = res.data;
        this.list = list;
        this.meta.total = total;
      } else {
        this.$message.error(res.message);
      }
    });
  }

  /** 整理请求参数 */
  private getParams() {
    const { page, size } = this.meta;

    const {
      title_name,
      reserve_service_id,
      member_name,
      user_name,
      arrive_time,
      created_time,
      third_id
    } = this.search;

    const third_ids = third_id ? [Number(third_id)] : [];

    const params = {
      page,
      size,
      invite_template_id: this.emitData.id || 0, //邀请函模板id
      title_name, //模板名称
      reserve_service_id: reserve_service_id || 0, //服务类型id
      member_name, //导购名称
      user_name, //user_name
      third_ids: third_ids, //门店ids
      reserve_service_start_time: arrive_time[0]
        ? Number(arrive_time[0] / 1000)
        : 0, //服务开始时间
      reserve_service_end_time: arrive_time[1]
        ? Number(arrive_time[1] / 1000)
        : 0,
      order_create_start_time: created_time[0]
        ? Number(created_time[0] / 1000)
        : 0, //订单筛选开始
      order_create_end_time: created_time[1]
        ? Number(created_time[1] / 1000)
        : 0
    };

    if (title_name !== this.emitData.title_name) {
      delete params.invite_template_id;
    }

    return params;
  }

  /** 导出 */
  private exportList() {
    const params = this.getParams();
    downFilePost(API.invitationDataExport, params, "邀请数据");
  }

  /** 查询 */
  private goSearch() {
    this.meta.page = 1;
    this.getList();
  }

  /** 重置 */
  private reset() {
    this.search = {
      title_name: "",
      reserve_service_id: "",
      member_name: "",
      user_name: "",
      arrive_time: [],
      created_time: [],
      third_id: "",
      storeName: ""
    };
    this.meta.page = 1;
    this.getList();
  }

  /** 每页数量改变 */
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getList();
  }

  /** 页码改变 */
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getList();
  }
  /** 部门树 */
  private dp_props = {
    label: "dp_name",
    children: [],
    isLeaf: "leaf"
  };
  /**动态加载节点 */
  private async loadNode(node: any, resolve: any) {
    if (node.level === 0) {
      const {
        data: { list }
      } = await storeAPI.storeDataApi({});
      list.forEach((item: any) => {
        if (!item.is_child_department) item.leaf = true;
      });
      resolve(list);
    }
    if (node.level >= 1) {
      this.getIndex(node, resolve);
    }
  }
  /** 清空回调 */
  private changeStore() {
    this.search.storeName = "";
    this.search.third_id = "";
    this.goSearch();
  }
  /** 选择部门节点 */
  private handleNodeClick(node: { dp_name: string; dp_third_id: string }) {
    this.search.storeName = node.dp_name;
    this.search.third_id = node.dp_third_id;
    this.goSearch();
  }
  /**关闭门店筛选弹窗 */
  private blurFn() {
    (this.$refs.selectBlur as HTMLFormElement).blur();
  }
  /**加载叶子节点数据函数 */
  private getIndex(node: any, resolve: any) {
    const params = { dp_id: node.data.dp_id };
    storeAPI.storeDataApi(params).then(res => {
      this.setChecked(res.data.list);
      resolve(res.data.list);
    });
  }
  /**判断部门是否有子节点 */
  private setChecked(arr: any) {
    arr.forEach((item: any) => {
      if (!item.is_child_department) item.leaf = true;
    });
  }
}
</script>
<style lang="scss" scoped>
.search {
  background-color: $color-background;
  padding: 20px;
  padding-bottom: 0px;
  margin-bottom: 20px;
}
.table_footer {
  margin: 20px 0 0;
  text-align: right;
}
.search-select {
  width: 280px;
}
.el-select-dropdown__list .el-select-dropdown__item.hover {
  background-color: #fff !important;
}
.el-select-dropdown__item.selected {
  font-weight: normal !important;
  height: auto !important;
}
/deep/ .el-range-separator {
  line-height: 26px;
}
</style>
