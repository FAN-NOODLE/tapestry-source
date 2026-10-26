<template>
  <div class="friend-circle-list">
    <layout-list>
      <template slot="buttons">
        <el-button
          type="primary"
          @click="routerLinkToEdit(0)"
          size="small"
        >发布朋友圈</el-button>
      </template>
      <template slot="top-tabs">
        <el-tabs
          v-model="tabType"
          @tab-click="resetFn"
        >
          <el-tab-pane
            label="按员工"
            name="first"
          ></el-tab-pane>
          <el-tab-pane
            label="按内容"
            name="second"
          ></el-tab-pane>
        </el-tabs>
      </template>
      <template slot="search">
        <div
          class="employee-search"
          v-if="tabType == 'first'"
        >
          <search-item text="所属部门：">
            <el-select
              size="small"
              style="width:280px"
              v-model="department.name"
              placeholder="请选择部门"
              ref="selectBlur"
            >
              <el-option
                :value="department.value"
                :label="department.name"
                style="height: auto;padding: 0px;"
              >
                <el-tree
                  ref="tree"
                  class="treeStyle"
                  v-if="department.isShowTree"
                  :props="department.data"
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
          <search-item text="员工：">
            <cy-select
              v-model="userEmployee"
              placeholder="请选择员工"
              search-input
              :optionsList="searchStaffList"
              :options="{ label: 'name', value: 'id' }"
            >
            </cy-select>
          </search-item>
          <search-item text="最近操作时间：">
            <el-date-picker
              v-model="timeChoose"
              type="datetimerange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              :default-time="['00:00:00', '23:59:59']"
              style="height:36px"
              format="yyyy-MM-dd HH:mm:ss"
              value-format="timestamp"
            >
            </el-date-picker>
          </search-item>
          <div class="user-buttongroup">
            <el-button
              type="primary"
              class="user-remarksbutton"
              @click="userSearch"
              size="small"
            >查询</el-button>
            <el-button
              plain
              class="user-remarksbuttons"
              @click="resetFn"
              size="small"
            >重置</el-button>
          </div>
        </div>
        <div
          class="content-search"
          v-if="tabType == 'second'"
        >
          <search-item text="主题名称：">
            <cy-input
              ref="userSubjectName"
              placeholder="请输入主题名称"
              v-model.trim="userSubjectName"
              size="small"
            ></cy-input>
          </search-item>
          <search-item text="创建时间：">
            <el-date-picker
              v-model="timeChoose"
              type="datetimerange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              :default-time="['00:00:00', '23:59:59']"
              style="height:36px"
              format="yyyy-MM-dd HH:mm:ss"
              value-format="timestamp"
            >
            </el-date-picker>
          </search-item>
          <search-item text="状态：">
            <cy-select
              v-model="userStatus"
              size="medium"
              placeholder="请选择状态"
            >
              <cy-option
                label="全部"
                :value="0"
              ></cy-option>
              <cy-option
                label="已发布"
                :value="1"
              ></cy-option>
              <cy-option
                label="定时未发布"
                :value="2"
              ></cy-option>
            </cy-select>
          </search-item>
          <div>
            <el-button
              type="primary"
              @click="userSearch"
              size="small"
            >查询</el-button>
            <el-button
              plain
              class="user-resetbutton"
              @click="resetFn"
              size="small"
            >重置</el-button>
          </div>
        </div>
      </template>
      <template slot="table">
        <div
          v-if="tabType == 'first'"
          :key="Math.random()"
        >
          <el-table
            style="width:100%"
            :data="employeeDataList"
            stripe
            max-height="560"
          >
            <el-table-column
              prop="user_name"
              label="所属员工"
            >
              <template slot-scope="scope">
                <span>{{ scope.row.user_name || "--" }}</span>
              </template>
            </el-table-column>
            <el-table-column
              prop="department_str"
              label="所属部门"
            >
              <template slot-scope="scope">
                <el-popover
                  placement="top-start"
                  width="300"
                  trigger="hover"
                  :disabled="scope.row.department_str.length <= 10"
                >
                  <div>{{ scope.row.department_str || "--" }}</div>
                  <span
                    slot="reference"
                    v-if="scope.row.department_str.length <= 10"
                  >{{ scope.row.department_str || "--" }}</span>
                  <span
                    slot="reference"
                    v-if="scope.row.department_str.length > 10"
                  >{{ scope.row.department_str.substr(0, 10) + "..." }}</span>
                </el-popover>
              </template>
            </el-table-column>
            <el-table-column
              prop="topic_count"
              label="朋友圈数"
            ></el-table-column>
            <el-table-column
              prop="updated_at"
              label="最近操作时间"
            >
              <template slot-scope="scope">
                <span>{{
                  scope.row.updated_at | formatDate("yyyy-mm-dd hh:mm:ss")
                }}</span>
              </template>
            </el-table-column>
            <el-table-column
              label="操作"
              width="103"
              fixed="right"
            >
              <template slot-scope="scope">
                <el-button
                  type="text"
                  class="work-look"
                  @click.native.prevent="goDetail(scope.row)"
                >查看详情</el-button>
              </template>
            </el-table-column>
            <div
              slot="empty"
              class="table-empty-label"
            >
              <img
                class="empty-img"
                :src="searchLoading
                    ? 'https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/empty-search.png'
                    : require('@/assets/images/empty-box.png')
                  "
                alt=""
              />
              <div>{{ searchLoading ? "没有找到相关员工" : "暂无朋友圈" }}</div>
            </div>
          </el-table>
        </div>
        <div
          v-if="tabType == 'second'"
          :key="Math.random()"
        >
          <el-table
            style="width:100%"
            :data="contentDataList"
            stripe
            max-height="560"
          >
            <el-table-column
              prop="name"
              label="主题名称"
              width="360"
            >
              <template slot-scope="scope">
                <el-popover
                  placement="top-start"
                  width="300"
                  trigger="hover"
                  :disabled="scope.row.name.length <= 20"
                >
                  <div>{{ scope.row.name }}</div>
                  <span
                    slot="reference"
                    v-if="scope.row.name.length <= 20"
                  >{{
                    scope.row.name
                  }}</span>
                  <span
                    slot="reference"
                    v-if="scope.row.name.length > 20"
                  >{{
                    scope.row.name.substr(0, 20) + "..."
                  }}</span>
                </el-popover>
              </template>
            </el-table-column>
            <el-table-column
              prop="content"
              label="朋友圈内容"
              min-width="410"
            >
              <template slot-scope="scope">
                <el-popover
                  placement="top-start"
                  width="300"
                  trigger="hover"
                  :disabled="scope.row.content_json.text.length <= 20"
                >
                  <div>{{ scope.row.content_json.text }}</div>
                  <span
                    slot="reference"
                    v-if="scope.row.content_json.text.length <= 20"
                  >{{ scope.row.contentText }}</span>
                  <span
                    slot="reference"
                    v-if="scope.row.content_json.text.length > 20"
                  >{{
                    scope.row.content_json.text.substr(0, 20) +
                    "..." +
                    scope.row.materialType
                  }}</span>
                </el-popover>
              </template>
            </el-table-column>
            <el-table-column
              prop="member_count"
              label="发布员工数"
              width="120"
            ></el-table-column>
            <el-table-column
              v-if="false"
              prop="created_member_id"
              label="创建人"
              width="100"
            ></el-table-column>
            <el-table-column
              prop="created_at"
              label="创建时间"
              width="190"
            >
              <template slot-scope="scope">
                <span>{{
                  scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss")
                }}</span>
              </template>
            </el-table-column>
            <el-table-column
              prop="status"
              label="状态"
              width="100"
            >
              <template slot-scope="scope">
                <div>
                  {{ scope.row.status == 1 ? "已发布" : "定时未发布" }}
                </div>
              </template>
            </el-table-column>
            <el-table-column
              prop="id"
              label="操作"
              fixed="right"
              width="124"
            >
              <template slot-scope="scope">
                <el-button
                  type="text"
                  class="operating-look button-group"
                  @click="routerLinkToEdit(1, scope.row.id)"
                >编辑</el-button>
                <el-button
                  type="text"
                  class="color-danger button-groups"
                  @click.native.prevent="deleteRow(scope.$index, scope.row.id)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>

            <div
              slot="empty"
              class="table-empty-label"
            >
              <img
                class="empty-img"
                :src="searchLoading
                    ? 'https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/empty-search.png'
                    : require('@/assets/images/empty-box.png')
                  "
                alt=""
              />
              <p>{{ searchLoading ? "没有找到相关朋友圈" : "暂无朋友圈" }}</p>
            </div>
          </el-table>
        </div>
      </template>
      <template slot="footer">
        <div
          v-if="tabType == 'first' && employeeDataList.length > 0"
          :key="Math.random()"
        >
          <el-pagination
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
        <div
          v-if="tabType == 'second' && contentDataList.length > 0"
          :key="Math.random()"
        >
          <el-pagination
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
      </template>
    </layout-list>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import API from "@/api/enterprise/friend-circle";
import Api from "@/api/enterprise/enterprise";

@Component({
  components: {
    LayoutList
  }
})
export default class extends Vue {
  /****************************** data开始 ******************************/
  private tabType = "first"; // top-tabs分页
  private userEmployee = "全部"; // 所选员工
  private searchStaffList: Array<any> = []; // 员工下拉列表数据
  // 分页
  private meta = {
    size: 10,
    page: 1,
    total: 0
  };
  private department = {
    value: 0, // 部门id
    list: [], // 部门数据
    name: "全部",
    isShowTree: false,
    data: {
      label: "name",
      children: [],
      isLeaf: "leaf"
    }
  };
  private userSubjectName = ""; // 主题
  private userStatus = 0; // 状态
  private searchLoading = false; // 是否正在搜索
  private employeeDataList: Array<any> = []; // 员工列表
  private contentDataList: Array<any> = []; // 内容列表
  private timeChoose = 0; // 时间
  private info = {
    list: {
      content: "", // 内容
      created_at: 0, // 创建时间
      created_member_id: 0, // 朋友圈创建员工ID
      id: 0, // id
      member_count: 0, // 员工数量
      name: "", // 主题名称
      publish_time: 0, // 定时发布时间
      publish_type: 0, // 发布类型 1立即发送 2定时发送
      status: 0
    },
    page: 0,
    size: 0,
    total: 0
  };
  /****************************** methods开始 ******************************/
  // 获取部门
  fetchStoreList() {
    Api.depatListApi({}).then(
      res => {
        const obj: any = { id: 0, name: "全部" };
        res.data.unshift(obj);
        this.department.list = res.data;
        this.department.isShowTree = true;
      },
      err => {
        if (err.error === "ERROR") {
          this.$message({
            showClose: true,
            message: err.message,
            type: "error"
          });
        }
      }
    );
  }
  // 获取朋友圈按内容列表展示
  private goDetailList() {
    const params = {
      start_time: (this.timeChoose && this.timeChoose[0] / 1000) || undefined,
      end_time: (this.timeChoose && this.timeChoose[1] / 1000) || undefined,
      name: this.userSubjectName || undefined,
      page: this.meta.page,
      size: this.meta.size,
      status: Number(this.userStatus) || undefined
    };
    API.momentContentListApi(params).then((res: any) => {
      if (res.code === 200) {
        res.data.list.map((item: any) => {
          if (item.content_json.picture.length > 0) {
            item.contentText = item.content_json.text
              ? item.content_json.text + "[图片]"
              : "[图片]";
            item.materialType = "[图片]";
          } else if (
            item.content_json.link &&
            item.content_json.link.link_url
          ) {
            item.contentText = item.content_json.text
              ? item.content_json.text + "[链接]"
              : "[链接]";
            item.materialType = "[链接]";
          } else if (
            item.content_json.applet &&
            item.content_json.applet.app_name
          ) {
            item.contentText = item.content_json.text
              ? item.content_json.text + "[小程序]"
              : "[小程序]";
            item.materialType = "[小程序]";
          } else if (item.content_json.video && item.content_json.video.url) {
            item.contentText = item.content_json.text
              ? item.content_json.text + "[视频]"
              : "[视频]";
            item.materialType = "[视频]";
          } else {
            item.contentText = item.content_json.text;
            item.materialType = "";
          }
        });
        this.meta = {
          page: res.data.page,
          size: res.data.size,
          total: res.data.total
        };
        this.contentDataList = res.data.list;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 获取朋友圈按员工列表展示
  private getEmployeeList() {
    const params = {
      department_id: this.department.value || undefined,
      start_time: (this.timeChoose && this.timeChoose[0] / 1000) || undefined,
      end_time: (this.timeChoose && this.timeChoose[1] / 1000) || undefined,
      member_id:
        (this.userEmployee !== "全部" && this.userEmployee) || undefined,
      page: this.meta.page,
      size: this.meta.size
    };
    API.momentMemberListApi(params).then((res: any) => {
      if (res.code === 200) {
        this.meta = {
          page: res.data.page,
          size: res.data.size,
          total: res.data.total
        };
        this.employeeDataList = res.data.list;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 搜索员工列表
  private searchMemberList(is_reset?: number) {
    const params = {
      department_id: this.department.value,
      page: 1,
      size: 999
    };
    Api.memberListApi(params).then((res: any) => {
      if (res.code == 200) {
        const obj: any = { id: 0, name: "全部" };
        this.searchStaffList = res.data.member_list || [];
        this.searchStaffList.unshift(obj);
        if (is_reset && this.tabType == "first") {
          (this.$refs["tree"] as any).setCurrentKey(0);
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //动态加载节点
  loadNode(node: any, resolve: any) {
    // 一级节点处理
    if (node.level === 0) {
      this.setChecked(this.department.list);
      return resolve(this.department.list);
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
    this.department.name = node.name;
    this.department.value = node.id;
    (this.$refs.selectBlur as any).blur();
    this.searchStaffList = [];
    this.searchMemberList();
  }
  // 当前请求应该请求的列表
  private currentRequestList() {
    if (this.tabType === "first") {
      this.getEmployeeList();
    } else {
      this.goDetailList();
    }
  }
  //查询
  private userSearch() {
    this.meta.page = 1;
    this.searchLoading = true;
    this.currentRequestList();
  }
  //重置
  private resetFn() {
    this.meta.page = 1;
    this.userEmployee = "全部"; // 员工重置
    this.department.name = "全部"; // 部门重置
    this.department.value = 0; // 部门重置
    this.userSubjectName = ""; // 主题重置
    this.userStatus = 0; // 状态重置
    this.searchLoading = false; //是否正在查询重置
    this.timeChoose = 0; //时间重置
    this.fetchStoreList();
    this.searchMemberList(1);
    this.currentRequestList();
  }
  // 按员工列表-点击查看详情按钮
  private goDetail(item: any) {
    this.$router.push({
      name: "myapp_enterprise_friend-circle_friend-circle-detail",
      query: {
        mark: item.mark,
        member_id: item.member_id,
        user_name: item.user_name
      }
    });
    sessionStorage.setItem("LISTTYPE", this.tabType);
  }
  // 去发布、编辑朋友圈页面
  private routerLinkToEdit(is_edit?: any, moment_topic_id?: any) {
    this.$router.push({
      name: "myapp_enterprise_friend-circle_friend-circle-edit",
      query: {
        is_edit,
        moment_topic_id
      }
    });
    sessionStorage.setItem("LISTTYPE", this.tabType);
  }
  // 删除行
  private deleteRow(index: any, id: number) {
    this.$confirm(
      "你确认要删除该朋友圈吗？删除后相关员工的朋友圈都不再展示该内容",
      "删除朋友圈",
      {
        confirmButtonText: "确定删除",
        cancelButtonText: "我再想想",
        type: "warning"
      }
    )
      .then(() => {
        API.momentTopicDeleteApi({
          id
        }).then((res: any) => {
          if (res.code === 200) {
            this.$message({
              type: "success",
              message: "删除成功!"
            });
            this.goDetailList();
          }
        });
      })
      .catch(() => {
        this.$message({
          type: "info",
          message: "已取消删除"
        });
      });
  }
  //内部绑定函数（获取每页数据量）
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.currentRequestList();
  }
  //内部绑定函数（获取当前第几页）
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.currentRequestList();
  }
  /****************************** 生命周期开始 ******************************/
  private created() {
    if (sessionStorage.getItem("LISTTYPE")) {
      this.tabType = sessionStorage.getItem("LISTTYPE") || "first";
      sessionStorage.removeItem("LISTTYPE");
    }
  }
  private mounted() {
    this.currentRequestList();
    // 获取部门下拉列表
    this.fetchStoreList();
    this.searchMemberList();
  }
}
</script>
<style lang="scss" scoped>
/deep/.treeStyle {
  padding: 5px !important;

  .el-tree-node.is-current.is-focusable.is-checked {
    color: #409eff !important;
  }
}

/deep/.el-table__fixed-right {
  right: 0 !important;
}

/deep/.el-input--small .el-input__inner {
  height: 36px;
}

/deep/.gutter {
  width: 0px !important;
  display: none;
}

.friend-circle-list {
  .date-updatatime {
    display: inline-block;
    margin-top: 7px;
    margin-bottom: 10px;
    font-family: SourceHanSansSC;
    font-weight: 400;
    font-size: 12px;
    color: rgba(16, 16, 16, 0.7);
    font-style: normal;
    letter-spacing: 0px;
    line-height: 18px;
    text-decoration: none;
  }
}

.user-search div {
  margin-right: 12px;
  margin-top: 5px;
}

.user-search span {
  margin-right: 5px;
  text-align: center;
}

.user-buttongroup {
  margin-left: 0;
}

.work-look {
  color: #708df5;
}

/deep/.slide-box {
  position: absolute;
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

.button-group {
  margin-left: 0 !important;
  margin-right: 20px;
  padding-left: 0px;
}

.button-groups {
  margin-left: 0 !important;
  padding-left: 0px;
}</style>
