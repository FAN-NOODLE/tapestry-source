<template>
  <div class="supernatant">
    <layout-list>
      <template slot="search">
        <search-item text="客户姓名：">
          <el-input
            style="width:280px;"
            size="small"
            prefix-icon="el-icon-search"
            placeholder="请输入客户姓名"
            clearable
            class="custom-name"
            v-model="search.customeName"
          ></el-input>
        </search-item>
        <search-item text="所属部门：">
          <el-select
            style="width:280px"
            v-model="search.department.name"
            placeholder="全部"
            ref="selectBlur"
          >
            <el-option
              :value="search.department.value"
              :label="search.department.name"
              style="height: auto;padding: 0px;"
            >
              <el-tree
                ref="tree"
                class="treeStyle"
                v-if="search.department.isShowTree"
                :props="search.department.data"
                :load="loadNode"
                node-key="id"
                :expand-on-click-node="true"
                :check-strictly="true"
                :check-on-click-node="true"
                @node-click="handleNodeClick"
                lazy
              ></el-tree>
            </el-option>
          </el-select>
        </search-item>
        <search-item text="所属员工：">
          <cy-select
            v-model="search.staff.value"
            placeholder="全部"
            search-input
            :optionsList="search.staff.list"
            :options="{ label: 'name', value: 'id' }"
          >
          </cy-select>
        </search-item>
        <search-item text="渠道来源：">
          <cy-select
            v-model="search.channel.value"
            placeholder="全部"
            search-input
            :optionsList="search.channel.list"
            :options="{ label: 'name', value: 'id' }"
          >
          </cy-select>
        </search-item>
        <search-item text="时间：" class="search-time">
          <el-select
            v-model="search.time.value"
            placeholder="全部"
            style="width:110px;"
            class="time-select"
          >
            <el-option
              v-for="item in search.time.list"
              :label="item.label"
              :value="item.value"
              :key="item.value"
            ></el-option>
          </el-select>
          <el-date-picker
            v-model="search.time.datePicker"
            :default-time="['00:00:00', '23:59:59']"
            value-format="yyyy-MM-dd HH:mm:ss"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            clearable
          >
          </el-date-picker>
        </search-item>
        <div>
          <el-button type="primary" @click="searchStaff()" size="small"
            >查询</el-button
          >
          <el-button @click="resetStaff()" size="small">重置</el-button>
        </div>
      </template>
      <template slot="tabs">
        <el-tabs v-model="activTab" @tab-click="handleTabClick(activTab)">
          <el-tab-pane label="员工删除客户" name="1">
            <!-- <el-button class="m-b-20" @click="openDialog(2)" size="small"
              >删除提醒通知配置</el-button
            > -->
          </el-tab-pane>
          <el-tab-pane label="客户删除员工" name="2">
            <!-- <el-button class="m-b-20" @click="openDialog(1)" size="small"
              >被删除提醒通知配置</el-button
            > -->
          </el-tab-pane>
        </el-tabs>
      </template>
      <template slot="table">
        <!-- stripe  是否为斑马纹 table -->
        <el-table
          tooltip-effect="light"
          :data="staffList"
          stripe
          style="width: 100%"
        >
          <el-table-column label="客户" width="200" show-overflow-tooltip>
            <template slot-scope="scope">
              <img
                :src="
                  scope.row.customer_avatar
                    ? scope.row.customer_avatar
                    : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                "
                alt=""
                class="table-text-img"
              />
              <span class="hover-name"
                >{{ scope.row.customer_name
                }}<span v-if="scope.row.remark"
                  >({{ scope.row.remark }})</span
                ></span
              >
            </template>
          </el-table-column>
          <!-- <el-table-column label="渠道来源" width="150" show-overflow-tooltip>
            <template>
              <span>-</span>
            </template>
          </el-table-column> -->
          <el-table-column
            prop="mobile"
            label="所属员工"
            min-width="200"
            show-overflow-tooltip
          >
            <template slot-scope="scope">
              <img
                :src="
                  scope.row.member_avatar
                    ? scope.row.member_avatar
                    : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                "
                alt=""
                class="table-text-img"
              />
              <span>{{ scope.row.member_name }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="所属部门"
            min-width="200"
            show-overflow-tooltip
          >
            <template slot-scope="scope">
              <span v-if="scope.row.member_departs">
                {{ scope.row.dep_name }}
              </span>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column
            label="添加时间"
            min-width="200"
            show-overflow-tooltip
          >
            <template slot-scope="scope">
              <span v-if="scope.row.create_time > 0">{{
                scope.row.create_time | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column
            label="删除时间"
            min-width="200"
            prop="relieve_at"
            show-overflow-tooltip
          >
            <template slot-scope="scope">
              <span v-if="scope.row.relieve_at > 0">{{
                scope.row.relieve_at | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <!-- <el-table-column label="最后联系时间" min-width="150">
            <template>
              <span>-</span>
            </template>
          </el-table-column>
          <el-table-column label="累计对话数" min-width="150">
            <template>
              <span>-</span>
            </template>
          </el-table-column> -->
          <el-table-column label="操作" min-width="72" fixed="right">
            <template slot-scope="scope">
              <el-button
                type="text"
                size="small"
                @click.native.prevent="goDetail(scope.row.user_id)"
              >
                查看
              </el-button>
            </template>
          </el-table-column>
          <div slot="empty" v-if="!emptyFlag" class="table-empty-label">
            <img class="empty-img" src="@/assets/images/empty-box.png" alt="" />
            <div>暂无流失客户</div>
          </div>
          <div slot="empty" v-else class="table-empty-label">
            <img
              class="empty-img"
              src="@/assets/images/empty-search.png"
              alt=""
            />
            <div class="desc">没有找到相关客户</div>
          </div>
        </el-table>
      </template>
      <template slot="footer" v-if="meta.total > 0">
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
      </template>
    </layout-list>
    <delete-warn-set
      v-if="dialog_editapp"
      :titleType="titleType"
      @close="close"
    />
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import deleteWarnSet from "../../../components/delete-warn-set.vue";
import Api from "@/api/enterprise/enterprise";
import API from "@/api/enterprise/drainage-coach";
import { publicMethod, isNull } from "@/utils/tools";
@Component({
  components: {
    deleteWarnSet
  }
})
export default class extends Vue {
  // ---------------------------------------------------------------此页面所需变量---------------------------------------------------------------
  // 查询数据
  private search = {
    // 客户姓名
    customeName: "",
    // 部门
    department: {
      value: 0, // 部门id
      list: [], // 部门数据
      name: "",
      isShowTree: false,
      data: {
        label: "name",
        children: [],
        isLeaf: "leaf"
      }
    },
    // // 渠道
    channel: {
      list: [],
      value: 0
    },
    // 员工
    staff: {
      list: [],
      value: "全部"
    },
    time: {
      list: [
        { label: "添加时间", value: 1 },
        { label: "删除时间", value: 2 }
      ],
      value: 1,
      datePicker: ""
    },
    input: ""
  };
  // tab选项中id
  private activTab = "1";
  private handleTab = 1;
  //左侧重复客户列表
  private leftMenu: any = [];
  private loading = false;
  //默认显示条数
  private cou = 20;
  private repeatTotal = 10;
  // 表格数组
  private staffList: Array<any> = [];
  private emptyFlag = false;
  // 分页配置
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  // 弹框组件传值
  private titleType = 1;
  private dialog_editapp = false;
  private member_departsLame = "";
  // ---------------------------------------------------get计算属性--------------------------------------------------------------
  get noMore() {
    // 判断加载条数是否大于列表数据长度
    if (this.cou > this.leftMenu.length) {
      this.loading = false;
      return true;
    } else {
      return false;
    }
  }
  get countData() {
    // 计算属性使用切片生成新数组
    let data = [];
    // 大于三条，使用切片，返回新数组
    if (this.leftMenu.length > 20) {
      data = this.leftMenu.slice(0, this.cou);
      this.loading = true;
      return data;
    } else {
      // 否则使用原来数组，不进行切片处理
      data = this.leftMenu;
      return data;
    }
  }
  // ---------------------------------------------------------------方法---------------------------------------------------------------
  //渠道数据接口API
  private channeData() {
    const params = {
      page: 1,
      size: 9999
    };
    API.channelsListApi(params).then((res: any) => {
      if (res.code == 200) {
        const obj = { id: 0, name: "全部" };
        if (res.data.list) {
          res.data.list.unshift(obj);
          this.search.channel.list = res.data.list;
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 渠道选择
  private selectChange(itemId: number) {
    this.search.channel.value = itemId;
  }
  //动态加载节点
  loadNode(node: any, resolve: any) {
    // 一级节点处理
    if (node.level === 0) {
      this.setChecked(this.search.department.list);
      return resolve(this.search.department.list);
    }
    // 其余节点处理
    if (node.level >= 1) {
      this.getIndex(node, resolve);
    }
  }
  //判断部门是否有子节点
  setChecked(arr: any) {
    arr.forEach((item: any) => {
      if (!item.children_count) item.leaf = true;
    });
  }
  // 加载叶子节点数据函数
  getIndex(node: any, resolve: any) {
    Api.depatListApi({
      pid: node.data.third_id
    }).then((res: any) => {
      this.setChecked(res.data);
      resolve(res.data);
    });
  }
  // 选择部门
  handleNodeClick(node: any) {
    this.search.department.name = node.name;
    this.search.department.value = node.id;
    (this.$refs.selectBlur as any).blur();
    this.meta.page = 1;
    this.search.staff.value = "全部";
    this.search.staff.list = [];
    this.getStaffList();
  }
  // 获取部门
  fetchStoreList() {
    const _this = this;
    Api.depatListApi({}).then(
      res => {
        const obj: any = { id: 0, name: "全部" };
        res.data.unshift(obj);
        _this.search.department.list = res.data;
        _this.search.department.isShowTree = true;
      },
      err => {
        if (err.error === "ERROR") {
          _this.$message({
            showClose: true,
            message: err.message,
            type: "error"
          });
        }
      }
    );
  }
  // 获取部门下的员工
  getStaffList() {
    const params = {
      department_id: this.search.department.value,
      page: 1,
      size: 999
    };
    Api.memberListApi(params).then((res: any) => {
      if (res.code == 200) {
        const obj = { id: 0, name: "全部" };
        if (res.data.member_list) {
          res.data.member_list.unshift(obj);
          this.search.staff.list = res.data.member_list;
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }

  // 流失客户列表
  losscustomersList(type: number) {
    const params = {
      // 1=员工删除客户 2=客户删除员工
      type: this.handleTab,
      // 客户名称
      customer_name: this.search.customeName,
      // 部门id
      department_id: this.search.department.value,
      // 员工id
      member_id: Number(this.search.staff.value) || 0,
      // 时间搜索类型：1=添加时间 2=删除时间 3=最后联系时间
      time_slot_type: Number(this.search.time.value),
      channel_id: Number(this.search.channel.value),
      // 开始时间
      time_start: publicMethod(this.search.time.datePicker[0]),
      // 结束时间
      time_end: publicMethod(this.search.time.datePicker[1]),
      page: this.meta.page,
      size: this.meta.size
    };
    Api.losscustomersListApi(params).then((res: any) => {
      if (res.code == 200) {
        this.staffList = res.data.list || [];
        if (this.staffList.length == 0 && type == 1) {
          this.emptyFlag = true;
        } else {
          this.emptyFlag = false;
        }
        this.staffList
          .map(v => v.member_departs)
          .forEach((item, index) => {
            this.staffList[index].dep_name = item
              .map((v1: any) => v1.department_name)
              .join("；");
          });
        this.meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 搜索查询
  searchStaff() {
    this.meta.page = 1;
    this.losscustomersList(1);
  }
  // 重置
  resetStaff() {
    this.search.time.value = 1;
    this.search.customeName = "";
    this.search.department.value = 0;
    this.search.department.name = "全部";
    this.search.staff.value = "全部";
    this.search.channel.value = 0;
    this.search.time.datePicker = "";
    this.meta.page = 1;
    // 部门
    this.fetchStoreList();
    // 员工
    this.getStaffList();
    (this.$refs.tree as any).setCurrentKey(null);
    this.losscustomersList(0);
  }
  // 切换tabs
  handleTabClick(type: any) {
    this.meta.page = 1;
    this.handleTab = Number(type);
    this.losscustomersList(0);
  }
  // 点击加载更多
  loadMore() {
    // 每次点击加三条
    this.cou += 20;
  }

  // 分页
  handleSizeChange(val: string) {
    this.$set(this.meta, "size", val);
    this.meta.page = 1;
    this.losscustomersList(0);
  }
  handleCurrentChange(val: string) {
    this.$set(this.meta, "page", val);
    this.losscustomersList(0);
  }
  // 发送消息配置弹框
  openDialog(type: number) {
    if (type == 1) {
      this.titleType = 1;
    } else {
      this.titleType = 2;
    }
    this.dialog_editapp = true;
  }
  // 关闭删除提醒弹框
  close() {
    this.dialog_editapp = false;
  }
  // 查看活动详情数据
  private goDetail(id: any) {
    const obj: any = {
      type: this.handleTab,
      // 客户名称
      customer_name: this.search.customeName,
      // 部门id
      department_name: this.search.department.name,
      department_id: this.search.department.value,
      // 员工id
      member_id: Number(this.search.staff.value) || 0,
      // 时间搜索类型：1=添加时间 2=删除时间 3=最后联系时间
      time_type: this.search.time.value,
      // 时间
      time_start: this.search.time.datePicker,
      page: this.meta.page,
      size: this.meta.size
    };
    sessionStorage.setItem("leaveData", JSON.stringify(obj));
    this.$router.push({
      name: "myapp_enterprise_statistics_leave-custome_detail",
      query: {
        user_id: id
      }
    });
  }
  // ---------------------------------------------------------------生命周期---------------------------------------------------------------

  // 实例挂载完
  mounted() {
    this.search.department.name = "全部";
    this.channeData();
    // 部门
    this.fetchStoreList();
    // 员工
    this.getStaffList();
    // 员工删除客户列表
    this.losscustomersList(0);
  }
  created() {
    let seaData: any = sessionStorage.getItem("leaveData");
    seaData = JSON.parse(seaData);
    if (seaData) {
      this.activTab = String(seaData.type);
      this.handleTab = seaData.type;
      this.search.customeName = seaData.customer_name;
      this.search.department.name = seaData.department_name;
      this.search.department.value = seaData.department_id;
      this.search.staff.value = seaData.member_id;
      this.search.time.value = seaData.time_type;
      if (!isNull(this.search.time.datePicker)) {
        this.search.time.datePicker = seaData.time_start;
      }
      this.meta.page = seaData.page;
      this.meta.size = seaData.size;
      sessionStorage.removeItem("leaveData");
    }
  }
}
</script>
<style lang="scss" scoped>
.hover-name {
  cursor: pointer;
}
.hover-name:hover {
  color: #708df5;
}
/deep/.custom-name .el-input__inner {
  height: 36px;
}
/deep/.time-select .el-input__inner {
  border-right: none;
  border-radius: 0px;
  border-top-left-radius: 5px;
  border-bottom-left-radius: 5px;
}
/deep/ .el-range-editor.el-input__inner {
  border-radius: 0px;
  border-top-right-radius: 5px;
  border-bottom-right-radius: 5px;
}

/deep/ {
  .el-scrollbar {
    .el-scrollbar__view {
      .el-select-dropdown__item {
        height: auto;
        padding: 0px;
        background: white;
      }
    }
  }
  .el-select-dropdown__item.selected {
    font-weight: normal;
    height: auto !important;
  }

  .layout-list .buttons-content {
    margin-bottom: 0;
  }
  .treeStyle {
    .el-tree-node.is-current.is-focusable.is-checked {
      color: #409eff !important;
    }
  }
}
.table-empty-label {
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
.table-text-img {
  width: 25px;
  height: 25px;
  border-radius: 50%;
  vertical-align: middle;
  margin-right: 10px;
}
.search-time {
  .search-item {
    margin-left: 30px;
  }
}
</style>
