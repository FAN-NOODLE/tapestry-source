<template>
  <div class="supernatant">
    <layout-list>
      <template slot="search">
        <search-item text="客户姓名：">
          <cy-input
            size="small"
            prefix-icon="el-icon-search"
            placeholder="请输入客户姓名"
            clearable
            v-model="search.customer_name"
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
        <search-item text="所属员工：">
          <cy-select
            v-model="search.staff.value"
            search-input
            :optionsList="searchStaffList"
            :options="{ label: 'name', value: 'id' }"
          >
          </cy-select>
        </search-item>
        <div>
          <el-button
            type="primary"
            @click="searchStaff()"
            size="small"
          >查询</el-button>
          <el-button
            @click="resetStaff()"
            size="small"
          >重置</el-button>
        </div>
      </template>
      <template slot="slide">
        <dl class="menu-left">
          <dt>
            重复客户<span v-if="leftMeta.total > 0"></span>({{
              leftMeta.total
            }})
          </dt>
          <div class="sub-wrap">
            <dd
              v-for="(item, index) in countData"
              :key="index"
              @click="setMenuFn(item, index)"
              :class="{
                'bag-color': index == isActive,
                dd: index !== isActive
              }"
            >
              <el-image fit="cover" v-if="item.avatar" class="avatar" :src="item.avatar"></el-image>
              <img v-else class="avatar" src="'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'"/>
              <span
                class="item-label"
                v-if="item.name"
              >{{ item.name }}</span>
              <span v-if="item.nick !== item.name">({{ item.nick }})</span>
            </dd>
            <div class="load-box">
              <div
                class="nomore-box"
                v-if="countData.length && leftMeta.total !== countData.length"
              >
                <el-button
                  type="text"
                  @click="loadMore"
                >加载更多</el-button>
                <img
                  class="load-arrow"
                  src="@/assets/images/loadArrow.png"
                  alt=""
                />
              </div>
              <el-button
                type="text"
                v-if="noMore"
                class="no-more"
              >没有更多了</el-button>
              <div
                class="left-empty-box"
                v-if="leftMenu.length == 0 && emptyFlag"
              >
                没有找到相关重复客户
              </div>
            </div>
          </div>
        </dl>
      </template>
      <template slot="slide-table">
        <!-- stripe  是否为斑马纹 table -->
        <el-table
          tooltip-effect="light"
          :data="staffList"
          stripe
          height="600"
          style="width:100%;"
        >
          <el-table-column
            label="所属员工"
            width="190"
            show-overflow-tooltip
          >
            <template slot-scope="scope">
              <img
                :src="scope.row.avatar
                    ? scope.row.avatar
                    : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                  "
                alt=""
                class="table-text-img"
              />
              <span @click="goDetail(scope.row.id)">{{ scope.row.name }}</span>
              <el-tag
                style="margin-left: 10px"
                size="small"
                v-if="scope.row.is_bind_guide == 1"
              >绑定导购</el-tag>
            </template>
          </el-table-column>
          <el-table-column
            prop="deperment_name"
            label="所属部门"
            min-width="100"
            show-overflow-tooltip
          >
            <template slot-scope="scope">
              <span v-if="!scope.row.department_name">-</span>
              <span v-else>{{ scope.row.department_name }}</span>
            </template>
          </el-table-column>
          <!-- <el-table-column prop="channel" label="渠道来源" min-width="100">
            <template>
              <span>-</span>
            </template>
          </el-table-column> -->
          <el-table-column
            prop="create_time"
            label="添加时间"
            min-width="200"
            show-overflow-tooltip
          >
            <template slot-scope="scope">
              <span v-if="scope.row.create_time != 0">{{
                scope.row.create_time | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <!-- <el-table-column
            prop="last_connect_time"
            label="最后联系时间"
            min-width="150"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.last_connect_time != 0">{{
                scope.row.last_connect_time | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
              <span v-else>-</span>
            </template>
          </el-table-column> -->
          <!-- <el-table-column
            prop="connect_num"
            label="累计对话数"
            min-width="150"
          >
            <template>
              <span>-</span>
            </template>
          </el-table-column> -->
          <div
            slot="empty"
            v-if="!emptyFlag && staffList.length == 0"
          >
            <img
              class="empty-img"
              src="@/assets/images/empty-box.png"
              alt=""
            />
            <div class="empty-box">
              <div class="desc">暂无重复客户</div>
            </div>
          </div>
          <div
            slot="empty"
            v-if="emptyFlag && staffList.length == 0"
          >
            <img
              class="empty-img"
              src="@/assets/images/empty-search.png"
              alt=""
            />
            <div class="empty-box">
              <div class="desc">没有找到相关员工</div>
            </div>
          </div>
        </el-table>
        <el-pagination
          class="float-right m-t-20"
          v-if="meta.total > 0"
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
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import Api from "@/api/enterprise/enterprise";
// import { isNull } from "@/utils/tools";
@Component({
  components: {
    LayoutList
  }
})
export default class extends Vue {
  // ---------------------------------------------------------------此页面所需变量---------------------------------------------------------------
  // 查询数据
  private search = {
    // 客户姓名
    customer_name: "",
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
    // 员工
    staff: {
      value: "全部" // 员工id
    }
  };
  // 部门列表
  private departmentList: Array<any> = [];
  // 员工列表
  private searchStaffList: Array<any> = [];
  //左侧重复客户列表
  private leftMenu: any = [];
  private isActive = 0; //默认第一个有样式
  //默认显示条数
  private leftMeta = {
    page: 1,
    size: 10,
    total: 0
  };
  // 表格数组
  private staffList: Array<any> = [];
  // 重复客户id
  private userWorkId = 0;
  private userId = 0;
  private emptyFlag = false;
  // 分页配置
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  // ---------------------------------------------------------------生命周期---------------------------------------------------------------

  // 实例挂载完
  mounted() {
    this.search.department.name = "全部";
    // 获取部门
    this.fetchStoreList();
    // 部门下的员工
    this.searchMemberList();
    // 重复客户
    this.getStaffList(0);
  }

  // ---------------------------------------------------get计算属性--------------------------------------------------------------

  get noMore() {
    // 判断加载条数是否大于列表数据长度
    return (
      this.leftMeta.size >= this.leftMenu.length && this.leftMenu.length != 0
    );
  }
  get countData() {
    // 计算属性使用切片生成新数组
    let data = [];
    // 大于三条，使用切片，返回新数组
    if (this.leftMenu.length > this.leftMeta.size) {
      data = this.leftMenu.slice(0, this.leftMeta.size);
      return data;
    } else {
      // 否则使用原来数组，不进行切片处理
      data = this.leftMenu;
      return data;
    }
  }
  // ---------------------------------------------------------------方法---------------------------------------------------------------

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
  // 选择部门
  handleNodeClick(node: any) {
    this.search.department.name = node.name;
    this.search.department.value = node.id;
    (this.$refs.selectBlur as any).blur();
    this.meta.page = 1;
    this.searchStaffList = [];
    this.searchMemberList();
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
  // 搜索员工列表
  private searchMemberList() {
    const params = {
      department_id: this.search.department.value,
      page: 1,
      size: 999
    };
    Api.memberListApi(params).then((res: any) => {
      if (res.code == 200) {
        const obj: any = { id: 0, name: "全部" };
        this.searchStaffList = res.data.member_list || [];
        this.searchStaffList.unshift(obj);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 搜索查询
  searchStaff() {
    this.isActive = 0;

    this.meta.page = 1;
    this.getStaffList(1);
  }
  // 重置
  resetStaff() {
    this.search.customer_name = "";
    this.search.department.value = 0;
    this.search.department.name = "全部";
    this.search.staff.value = "全部";
    this.isActive = 0;
    (this.$refs.tree as any).setCurrentKey(null);
    // 获取部门
    this.fetchStoreList();
    // 部门下的员工
    this.searchMemberList();
    this.meta = {
      total: 0, // 总数
      page: 1, // 当前页
      size: 10 // 每页数量
    };
    this.getStaffList(0);
  }
  //点击左边菜单
  private setMenuFn(item: any, index: number) {
    this.isActive = index;
    this.userWorkId = item.user_work_id;
    this.userId = item.user_id;
    this.meta.page = 1;
    this.memberList(0);
  }
  // 点击加载更多
  loadMore() {
    // 每次点击加三条
    this.leftMeta.size += this.leftMeta.size;
  }
  // 重复客户列表
  getStaffList(type: number) {
    const search = {
      customer_name: this.search.customer_name,
      department_id: this.search.department.value,
      member_id: Number(this.search.staff.value) || 0
    };
    const params = {
      search: search,
      page: this.meta.page,
      size: 999
    };
    Api.repeatCustomerApi(params).then((res: any) => {
      if (res.code == 200) {
        this.leftMenu = res.data.list ? res.data.list : [];
        this.leftMeta.total = res.data.total;
        if (this.leftMenu.length == 0) {
          if (type == 1) {
            this.emptyFlag = true;
            this.staffList = [];
          } else {
            this.emptyFlag = false;
          }
          this.meta = {
            page: 1,
            size: 10,
            total: 0
          };
        } else {
          this.userWorkId = this.leftMenu[0].user_work_id;
          this.userId = this.leftMenu[0].user_id;
          this.memberList(type);
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 表格员工列表接口
  private memberList(type: number) {
    const params = {
      department_id: this.search.department.value,
      member_id: Number(this.search.staff.value) || 0,
      user_work_id: this.userWorkId,
      user_id: this.userId,
      page: this.meta.page,
      size: this.meta.size
    };

    Api.repeatFindMemberApi(params).then((res: any) => {
      if (res.code == 200) {
        if (res.data.list.length == 0) {
          if (type == 1 || type == 0) {
            this.emptyFlag = true;
            this.staffList = [];
          } else {
            this.emptyFlag = false;
          }
        } else {
          this.staffList = res.data.list ? res.data.list : [];
        }
        this.meta = {
          page: res.data.page || 1,
          size: res.data.size || 10,
          total: res.data.total || 0
        };
      } else {
        this.$message.error(res.message);
      }
    });
  }

  // 分页
  handleSizeChange(val: number) {
    this.meta.size = val;
    this.meta.page = 1;
    this.memberList(0);
  }
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.memberList(0);
  }
}
</script>
<style lang="scss" scoped>
.no-more {
  color: #909399 !important;
}

.empty-img {
  width: 124px;
}

.el-table__empty-text {
  line-height: 0px !important;
}

.bag-color {
  background: #c7c7d0;
}

.left-empty-box {
  color: #909399;
  line-height: 600px;
}

.load-box {
  margin-top: 10px;
  text-align: center;
}

/deep/ .deparment {
  .el-input__inner {
    height: 36px;
  }
}

// aside
.menu-left {
  margin-right: 6px !important;
  flex: 0 0 400px;
  margin: 0;
  margin-right: 30px;
  background: $color-background;

  dt {
    background: #ebebf1;
    font-size: 14px;
    color: #606266;
    height: 36px;
    line-height: 36px;
    padding-left: 20px;
    font-weight: bold;
  }

  .sub-wrap {
    height: 570px;
    overflow-y: auto;
    padding-bottom: 20px;

    dd {
      display: flex;
      align-items: center;
      padding: 17px 20px;
    }

    .avatar {
      display: inline-block;
      width: 20px;
      height: 20px;
      margin-right: 10px;
      border-radius: 50%;
    }

    .nomore-box {
      display: flex;
      justify-content: center;
      align-items: center;

      .load-arrow {
        display: block;
        width: 12px;
        height: 12px;
      }
    }
  }

  dd {
    margin-inline-start: 0px;
    height: 20px;
    line-height: 20px;
    // background-color: rgba(255, 255, 255, 1);
    color: $color-text-main;
    font-size: 14px;
    padding: 10px 0;
    cursor: pointer;
  }

  .dd:hover {
    background: #ebebf1;
  }
}

/deep/ .slide-table {
  margin-bottom: 0 !important;
}

/deep/ .table-content {
  margin-bottom: 0 !important;
}

.float-right {
  float: right;
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
    padding: 5px !important;

    .el-tree-node.is-current.is-focusable.is-checked {
      color: #409eff !important;
    }
  }
}

.empty-box {
  text-align: center;
  line-height: 20px;
  position: relative;
  top: -10px;
}

.custom-tree-node {
  img {
    display: inline-block;
    width: 14px;
    height: 14px;
    margin-right: 5px;
  }
}

.uptatedTime {
  color: gray;
  font-size: 12px;
  margin-left: 10px;
}

.contant {
  display: flex;
}

.box-card {
  background: #f5f6fa;
}

.table-text-img {
  width: 25px;
  height: 25px;
  border-radius: 50%;
  vertical-align: middle;
  margin-right: 10px;
}

/deep/.el-select .el-input .el-select__caret {
  // display: flex;
  // align-items: center;
}

.item-label {
  max-width: 210px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  font-size: 14px;
  font-family: PingFangSC-Regular, PingFang SC;
  font-weight: 400;
  color: #303133;
  line-height: 20px;
}

/deep/ .slide-menu {
  width: 24% !important;
}</style>
