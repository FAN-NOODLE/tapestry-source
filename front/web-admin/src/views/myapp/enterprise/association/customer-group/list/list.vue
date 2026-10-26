<template>
  <div>
    <layout-list>
      <template slot="search">
        <div class="m-b-15">
          <span>数据更新时间: {{ refreshTime }}</span>
          <!-- <span class="m-l-10">
            <i class="el-icon-refresh refresh"></i>
            <el-button type="text">刷新</el-button>
          </span> -->
        </div>
        <search-item text="群名称：">
          <cy-input
            size="small"
            prefix-icon="el-icon-search"
            placeholder="请输入群名称"
            clearable
            v-model="search.title"
          ></cy-input>
        </search-item>
        <search-item text="群主：">
          <cy-select
            v-model="search.staff"
            placeholder="全部"
            search-input
            :optionsList="search.groupList"
            :options="{ label: 'name', value: 'id' }"
          >
          </cy-select>
        </search-item>
        <search-item text="所属渠道：">
          <cy-select
            v-model="search.channel"
            placeholder="全部"
            search-input
            :optionsList="search.channelList"
            :options="{ label: 'name', value: 'id' }"
          >
          </cy-select>
        </search-item>
        <search-item text="群ID：">
          <cy-input
            size="small"
            prefix-icon="el-icon-search"
            placeholder="请输入群ID"
            clearable
            v-model="search.groupId"
          ></cy-input>
        </search-item>
        <search-item text="群运营者：">
          <cy-input
            size="small"
            prefix-icon="el-icon-search"
            placeholder="请输入群运营人员名称"
            clearable
            v-model="search.mangerName"
          ></cy-input>
        </search-item>
        <search-item text="所属部门：">
          <el-select
            style="width:280px"
            v-model="search.department.name"
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
        <search-item text="时间：">
          <el-date-picker
            style="height:36px"
            v-model="choose_date"
            type="daterange"
            range-separator="-"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :default-time="['00:00:00', '23:59:59']"
            :clearable="false"
            value-format="timestamp"
            :picker-options="startTimeOptSetting"
          ></el-date-picker>
        </search-item>
        <search-item text="群类型：">
          <cy-select
            v-model="search.groupType"
            placeholder="全部"
            search-input
            :optionsList="search.groupTypeList"
            :options="{ label: 'name', value: 'id' }"
          >
          </cy-select>
        </search-item>
        <search-item>
          <el-button type="primary" @click="getGroupList">查询</el-button>
          <el-button @click="reSearch">重置</el-button>
          <el-button @click="downData">导出数据明细</el-button>
        </search-item>
      </template>
      <template slot="tabs" class="m-t-20 m-b-20">
        <div class="dashboard">
          <div class="item">
            <div class="item-title">
              社群总数（个）
              <el-popover
                placement="top-start"
                width="200"
                trigger="hover"
                content="当前查询条件下，已创建的企微客户社群的总个数（不含已删除社群）"
              >
                <i slot="reference" class="el-icon-question"></i>
              </el-popover>
            </div>
            <div class="item-num">{{ meta.total | thousandBit }}</div>
          </div>
          <div class="item">
            <div class="item-title">
              新入群人数
              <el-popover
                placement="top-start"
                width="200"
                trigger="hover"
                content="当前查询条件下，所有社群新入群人数的总和"
              >
                <i slot="reference" class="el-icon-question"></i>
              </el-popover>
            </div>
            <div class="item-num">{{ newMember | thousandBit }}</div>
          </div>
          <div class="item">
            <div class="item-title">
              累计退群人数
              <el-popover
                placement="top-start"
                width="200"
                trigger="hover"
                content="当前查询条件下，所有社群累计退群人数的总和"
              >
                <i slot="reference" class="el-icon-question"></i>
              </el-popover>
            </div>
            <div class="item-num">{{ leaveMember | thousandBit }}</div>
          </div>
          <div class="item">
            <div class="item-title">
              留存人数
              <el-popover
                placement="top-start"
                width="200"
                trigger="hover"
                content="当前查询条件下，所有社群的群成员人数总和"
              >
                <i slot="reference" class="el-icon-question"></i>
              </el-popover>
            </div>
            <div class="item-num">{{ validMember | thousandBit }}</div>
          </div>
        </div>
      </template>
      <template slot="operate">
        <div>
          <el-button type="primary" @click="updateGroup">同步群信息</el-button>
          <el-button @click="group_dialog = true">管理群类型</el-button>
          <el-button @click="group_set_dialog = true">批量社群设置</el-button>
        </div>
      </template>
      <template slot="table">
        <div>
          <!-- stripe  是否为斑马纹 table -->
          <el-table
            tooltip-effect="light"
            :data="showGroupList"
            stripe
            max-height="600"
            style="width:100%;"
          >
            <el-table-column
              prop="chat_id"
              label="群ID"
              min-width="180"
            ></el-table-column>
            <el-table-column
              prop="name"
              label="群名"
              min-width="150"
            ></el-table-column>
            <el-table-column
              prop="owner"
              label="群主"
              min-width="100"
            ></el-table-column>
            <el-table-column
              prop="channel_name"
              label="所属渠道"
              min-width="100"
            >
            </el-table-column>
            <el-table-column
              prop="chat_type_name"
              label="群类型"
              min-width="100"
            ></el-table-column>
            <el-table-column
              prop="operator"
              label="群运营者"
              min-width="100"
            ></el-table-column>
            <el-table-column
              prop="operator_departments"
              label="运营者所属部门"
              min-width="160"
            >
              <template slot-scope="{ row }">
                {{ row.operator_departments }}
              </template>
            </el-table-column>
            <el-table-column
              prop="valid_uv"
              label="群成员数"
              align="center"
              min-width="100"
            ></el-table-column>
            <el-table-column
              prop="all_loss_uv"
              align="center"
              label="累计退群人数"
              min-width="160"
            ></el-table-column>
            <el-table-column
              prop="add_uv"
              align="center"
              label="新入群人数"
              min-width="160"
            ></el-table-column>
            <el-table-column
              prop="create_time"
              label="创建时间"
              min-width="240"
            >
              <template slot-scope="{ row }">
                {{ row.create_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
              </template>
            </el-table-column>
            <el-table-column
              fixed="right"
              prop="deperment_name"
              label="操作"
              width="280"
            >
              <template slot-scope="{ row }">
                <el-button type="text" @click="setItemGroup(row)"
                  >社群设置</el-button
                >
                <el-button type="text" @click="groupDetail(row)"
                  >查看详情</el-button
                >
                <el-button type="text" @click="groupData(row)"
                  >数据统计</el-button
                >
              </template>
            </el-table-column>

            <div slot="empty" v-if="!emptyFlag && groupList.length == 0">
              <img
                class="empty-img m-t-40"
                src="@/assets/images/empty-box.png"
                alt=""
              />
              <div class="empty-box">
                <div class="desc">暂无客户社群</div>
              </div>
            </div>
            <div slot="empty" v-if="emptyFlag && groupList.length == 0">
              <img
                class="empty-img m-t-40"
                src="@/assets/images/empty-search.png"
                alt=""
              />
              <div class="empty-box">
                <div class="desc">没有找到对应社群</div>
              </div>
            </div>
          </el-table>
          <el-pagination
            class="float-right m-t-20"
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
      </template>
    </layout-list>
    <!-- 群类型 -->
    <addGroupType
      v-if="group_dialog"
      @close="group_dialog = false"
    ></addGroupType>
    <!-- 批量社群设置 -->
    <uploadFile
      v-if="group_set_dialog"
      @close="group_set_dialog = false"
      @submit="
        group_set_dialog = false;
        getGroupList();
      "
    ></uploadFile>
    <!-- 社群设置 -->
    <setGroup
      v-if="set_group_dialog"
      :chatId="chatId"
      @close="
        set_group_dialog = false;
        getGroupList();
      "
    ></setGroup>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import Api from "@/api/enterprise/enterprise";
import addGroupType from "../components/addGroupType.vue";
import uploadFile from "../components/upload-file.vue";
import setGroup from "../components/set-group.vue";
import {
  chatListApi,
  chatTypeManage,
  chatTypeList,
  chatChannelApi,
  getRefreshTime,
  updateGroup
} from "@/api/enterprise/association";
import dateFormat from "@/filter/FormatDate";
@Component({
  components: {
    addGroupType,
    uploadFile,
    setGroup
  },
  filters: {
    thousandBit(num: number | string) {
      num = num + "";
      return (
        num &&
        (num.toString().indexOf(".") !== -1
          ? num.toString().replace(/(\d)(?=(\d{3})+\.)/g, function($0, $1) {
              return $1 + ",";
            })
          : num.toString().replace(/(\d)(?=(\d{3})+$)/g, function($0, $1) {
              return $1 + ",";
            }))
      );
    }
  }
})
export default class extends Vue {
  private emptyFlag = false;
  private group_dialog = false;
  private group_set_dialog = false;
  private set_group_dialog = false;
  private search: any = {
    title: "",
    staff: 0,
    groupId: "",
    mangerName: "",
    groupList: [{ name: "全部", id: 0 }],
    channelList: [
      {
        name: "全部",
        id: 0
      }
    ],
    channel: 0,
    // 部门
    department: {
      value: 0, // 部门id
      name: "",
      isShowTree: false,
      data: {
        label: "name",
        children: [],
        isLeaf: "leaf"
      }
    },
    groupTypeList: [{ name: "全部", id: 0 }],
    groupType: 0
  }; // 搜索数据
  private choose_date: any = []; // 筛选时间
  private newMember = 0;
  private leaveMember = 0;
  private validMember = 0;
  // 部门列表
  private departmentList: Array<any> = [];
  private startTimeOptSetting = {
    disabledDate(time: Date) {
      const now: number =
        +new Date(new Date(new Date().toLocaleDateString()).getTime()) - 1;
      return time.getTime() > now;
    }
  };
  private groupList = [];
  private chatId = "";
  // 分页配置
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  // 刷新时间
  private refreshTime = "";
  // ------------------------------------ methods
  // 下载
  private downData() {
    import("@/utils/Export2Excel").then(excel => {
      const tHeader: any = [
        "群ID",
        "群名",
        "群主",
        "所属渠道",
        "群类型",
        "群运营者",
        "所属部门",
        "群成员人数",
        "退群人数",
        "新入群人数",
        "创建时间"
      ];
      const data: any = [];
      this.groupList.map((v: any) => {
        const create_time_d = dateFormat.func(
          v.create_time,
          "yyyy-mm-dd hh:mm:ss"
        );
        data.push([
          v.chat_id + "",
          v.name,
          v.owner,
          v.channel_name,
          v.chat_type_name,
          v.operator,
          v.operator_departments,
          v.valid_uv,
          v.all_loss_uv,
          v.add_uv,
          create_time_d
        ]);
      });
      excel.export_json_to_excel({
        header: tHeader,
        data,
        filename: `社群列表数据明细${new Date()
          .toLocaleDateString()
          .replace(/\//g, "")}`,
        sheetname: "社群列表数据明细",
        autoWidth: true,
        bookType: "xlsx"
      });
    });
  }
  // 重置
  private reSearch() {
    this.search.title = "";
    this.search.staff = 0;
    this.search.channel = 0;
    this.search.groupType = 0;
    this.search.groupId = "";
    this.search.mangerName = "";
    this.search.department.value = 0;
    this.search.department.name = "全部";
    this.choose_date = [
      +new Date(new Date(new Date().toLocaleDateString()).getTime()) -
        86400000 * 30,
      +new Date(new Date(new Date().toLocaleDateString()).getTime()) - 1
    ];
    this.meta.page = 1;
    this.meta.size = 10;
    this.getGroupList();
  }
  // 获取社群类型列表
  private getChatType() {
    chatTypeList().then(res => {
      console.log(res);
      res.data = res.data || [];
      this.search.groupTypeList = [{ name: "全部", id: 0 }, ...res.data];
    });
  }
  // 获取渠道列表
  private getChannel() {
    chatChannelApi({
      filter_name: "channels"
    }).then(res => {
      console.log(res);
      res.data = res.data || [];
      this.search.channelList = [{ name: "全部", id: 0 }, ...res.data];
    });
  }
  // 群类型管理
  private chatTypeManage(type: number, id?: number) {
    chatTypeManage({
      opt: type, // 1 - 添加，2 - 编辑，3 - 删除
      id: id
    }).then(res => {
      console.log(res);
    });
  }
  // 获取社群列表
  private getGroupList() {
    chatListApi({
      chat_name: this.search.title,
      chat_id: this.search.groupId,
      channel_id: this.search.channel,
      owner: this.search.staff,
      operator: this.search.mangerName,
      department_id: this.search.department.value,
      chat_type_id: this.search.groupType,
      sts: this.choose_date[0] / 1000,
      ets: parseInt(this.choose_date[1] / 1000 + "")
    }).then(res => {
      console.log(res);
      res.data = res.data || [];
      this.newMember = this.leaveMember = this.validMember = 0;
      res.data.forEach((v: any) => {
        this.newMember += v.add_uv || 0;
        this.leaveMember += v.all_loss_uv || 0;
        this.validMember += v.valid_uv || 0;
        // 无 显示 '--'
        v.chat_type_name = v.chat_type_name || "--";
        v.operator = v.operator || "--";
        v.operator_departments = v.operator_departments
          ? v.operator_departments.join("-")
          : "--";
        v.name = v.name || "--";
        v.channel_name = v.channel_name || "--";
      });
      this.groupList = res.data || [];
      this.meta.total = res.data.length;
      this.meta.page = 1;
      this.meta.size = 10;
    });
  }
  // 查看详情
  private groupDetail(row: any) {
    this.$router.push({
      name: "myapp_enterprise_association_customer-group_group-detail",
      query: {
        chatId: row.chat_id,
        ets: parseInt(this.choose_date[1] / 1000 + "") + ""
      }
    });
  }
  // 数据统计
  private groupData(row: any) {
    this.$router.push({
      name: "myapp_enterprise_association_customer-group_group-data",
      query: {
        chatId: row.chat_id
      }
    });
  }
  // 社群设置
  private setItemGroup(row: any) {
    this.chatId = row.chat_id;
    this.set_group_dialog = true;
  }
  // 执行群主下拉列表
  private async memberlistFn() {
    const res: any = await chatChannelApi({
      filter_name: "chat_owner"
    });
    if (res.code !== 200) {
      this.$message.error(res.message);
      return;
    }
    const data = res.data;
    this.search.groupList = [...this.search.groupList, ...data];
  }
  // 换页
  private handleSizeChange(e: number) {
    this.meta.page = 1;
    this.meta.size = e;
  }
  private handleCurrentChange(e: number) {
    this.meta.page = e;
  }
  // 选择部门
  handleNodeClick(node: any) {
    this.search.department.name = node.name;
    this.search.department.value = node.id;
    (this.$refs.selectBlur as any).blur();
  }
  // 获取部门
  fetchStoreList() {
    Api.depatListApi({ pid: 0 }).then((res: any) => {
      if (res.code == 200) {
        const obj = { id: 0, name: "全部" };
        this.departmentList = res.data;
        this.departmentList.unshift(obj);
        this.search.department.isShowTree = true;
      } else {
        this.$message.error(res.message);
        return;
      }
    });
  }
  //动态加载节点
  loadNode(node: any, resolve: any) {
    // 一级节点处理
    if (node.level === 0) {
      this.setChecked(this.departmentList);
      return resolve(this.departmentList);
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
      if (res.code == 200) {
        const data = res.data;
        if (data.length <= 0) {
          return;
        }
        this.setChecked(data);
        resolve(data);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //mounted
  private mounted() {
    this.search.department.name = "全部";
    this.reSearch();
    // 获取部门
    this.fetchStoreList();
    // 获取群主
    this.memberlistFn();
    // 获取社群类型
    this.getChatType();
    // 获取渠道列表
    this.getChannel();
    // 获取社群列表
    this.getGroupList();
    // 获取数据更新时间
    this.getRefreshTime();
  }
  // 获取数据更新时间
  private getRefreshTime() {
    getRefreshTime({}).then((res: any) => {
      if (res.code === 200) this.refreshTime = res.data;
    });
  }
  // 同步群数据
  private updateGroup() {
    const loading = this.$loading({
      lock: true,
      text: "群数据同步中...",
      spinner: "el-icon-loading",
      background: "rgba(0, 0, 0, 0.7)"
    });
    updateGroup({
      sts: this.choose_date[0] / 1000,
      ets: parseInt(this.choose_date[1] / 1000 + "")
    }).then((res: any) => {
      loading.close();
      if (res.code === 200) {
        this.getRefreshTime();
        this.getGroupList();
        this.$message.success("群数据同步成功");
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // --------------------------------计算属性
  // 群列表
  private get showGroupList() {
    return this.groupList.slice(
      (this.meta.page - 1) * this.meta.size,
      this.meta.page * this.meta.size
    );
  }
}
</script>
<style lang="scss" scoped>
.refresh {
  color: #708df5;
  cursor: pointer;
}
.dashboard {
  background-color: #f5f6fa;
  border-radius: 4px;
  padding: 40px 20px;
  margin: 20px 0;
  display: flex;
  .item {
    flex: 1;
    position: relative;
    text-align: center;
    &::after {
      content: "";
      position: absolute;
      right: 0;
      top: 50%;
      height: 70px;
      border-left: 1px solid #d7d7d7;
      transform: translateY(-50%);
    }
    &:nth-child(4):after {
      content: "";
      border-left: none;
    }
    .item-title {
      font-size: 17px;
    }
    .item-num {
      font-size: 34px;
      font-weight: 600;
      line-height: 40px;
      margin-top: 15px;
    }
  }
}
.float-right {
  float: right;
}
.m-t-40 {
  margin-top: 40px;
}
</style>
