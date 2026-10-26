<template>
  <div class="supernatant">
    <layout-list>
      <template slot="search">
        <div class="m-b-10">
          <el-button type="text" size="small" @click="updateStaffInfo()"
            >更新员工信息</el-button
          >
          <span class="uptatedTime"
            >上次更新时间：<span v-if="uptatedTime">{{
              uptatedTime | formatDate("yyyy-mm-dd hh:mm:ss")
            }}</span
            ><span v-else>-</span></span
          >
        </div>
        <search-item text="员工姓名：">
          <cy-input
            width-size="small"
            size="small"
            prefix-icon="el-icon-search"
            placeholder="请输入员工姓名"
            v-model="inputName"
          ></cy-input>
        </search-item>
        <search-item text="加入状态：">
          <cy-select width-size="small" size="small" v-model="select">
            <cy-option
              :label="item.label"
              v-for="(item, index) in addStatus"
              :key="index"
              :value="item.value"
            ></cy-option>
          </cy-select>
        </search-item>
        <div>
          <el-button type="primary" @click="searchStaff()" size="small"
            >查询</el-button
          >
          <el-button @click="resetStaff()" size="small">重置</el-button>
        </div>
      </template>
      <template slot="slide">
        <el-tree
          :props="tree_props"
          :load="loadNode"
          :highlight-current="true"
          node-key="id"
          :default-expanded-keys="expand_key"
          @node-click="handleCheckChange"
          empty-text="未获取到组织结构和成员信息"
          class="box-card"
          lazy
        >
          <span class="custom-tree-node" slot-scope="{ data }">
            <img src="~@/assets/images/group-ic.png" alt="" />
            <span :title="data.name">{{ data.name }}</span>
          </span>
        </el-tree>
      </template>
      <template slot="slide-table">
        <!-- stripe  是否为斑马纹 table -->
        <!-- :disabled="click_id === 0" -->
        <el-button
          :disabled="true"
          class="config-btn"
          size="small"
          @click="notifyConfig()"
          >客户流失提醒通知配置</el-button
        >
        <el-table :data="staffList" stripe>
          <el-table-column label="员工姓名" width="150" show-overflow-tooltip>
            <template slot-scope="scope">
              <img
                :src="
                  scope.row.avatar
                    ? scope.row.avatar
                    : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                "
                alt=""
                class="table-text-img"
              />
              <span @click="goDetail(scope.row.id)" class="hover-name">{{
                scope.row.name
              }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="gender_text" label="性别" min-width="70">
          </el-table-column>
          <el-table-column
            prop="forms"
            label="部门"
            min-width="100"
            show-overflow-tooltip
          >
          </el-table-column>
          <el-table-column
            prop="mobile"
            label="手机号"
            min-width="120"
            show-overflow-tooltip
          >
          </el-table-column>
          <el-table-column prop="status_text" label="加入状态" min-width="100">
          </el-table-column>
          <el-table-column label="操作" width="72" fixed="right">
            <template slot-scope="scope">
              <el-button
                type="text"
                size="small"
                @click.native.prevent="goDetail(scope.row.id)"
              >
                查看
              </el-button>
            </template>
          </el-table-column>
          <div slot="empty">
            <div class="empty-box" v-if="empty1">
              <img
                class="empty-img"
                src="@/assets/images/empty-box.png"
                alt=""
              />
              <div class="desc">未获取到组织结构和用户信息</div>
              <el-button type="text" size="small" @click="openIntroduce">
                如何获取组织结构？
              </el-button>
            </div>
            <div v-if="empty2" class="table-empty-label">
              <img
                class="empty-img"
                src="@/assets/images/empty-search.png"
                alt=""
              />
              <div class="desc">没有找到相关成员</div>
            </div>
            <div v-if="empty3" class="table-empty-label">
              <img
                class="empty-img"
                src="@/assets/images/empty-box.png"
                alt=""
              />
              <div class="desc">暂无成员</div>
            </div>
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
    <!-- 删除客户消息通知配置 -->
    <notify-config
      v-if="openNotify"
      :selectdpId="selectdpId"
      @closeDialog="closeDialog"
    ></notify-config>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import Api from "@/api/enterprise/enterprise";
import { isNull } from "@/utils/tools";
import NotifyConfig from "./components/notify-config.vue";
@Component({
  components: {
    LayoutList,
    NotifyConfig
  }
})
export default class extends Vue {
  // ---------------------------------------------------------------此页面所需变量---------------------------------------------------------------
  // 加入状态数据
  private addStatus = [
    { label: "已激活", value: 1 },
    { label: "已禁用", value: 2 },
    { label: "未激活", value: 4 },
    { label: "退出企业", value: 5 }
  ];
  // tree结构 isLeaf指定节点是否为叶子节点，仅在指定了 lazy 属性的情况下生效
  private tree_props = { label: "name", isLeaf: "leaf" };
  // 加入状态默认为已加入
  private select = 1;
  private empty1 = false;
  private empty2 = false;
  private empty3 = false;
  // 设置默认展开的节点
  private expand_key: any[] = [];
  // 员工姓名
  private inputName = "";
  // 所选部门id
  private click_id = 0;
  // 最新更新时间
  private uptatedTime = "";
  // 分页配置
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  // ---------------------------------------------------------------方法---------------------------------------------------------------
  // 员工列表
  private staffList: Array<any> = [];
  private openNotify = false;
  private selectdpId = 0;
  // 懒加载组织节点
  loadNode(node: any, resolve: any) {
    if (node.level === 0) {
      const self = this;
      Api.depatListApi({ pid: 0 }).then((res: any) => {
        if (res.code == 200) {
          const data = res.data;
          if (isNull(data)) {
            this.empty1 = true;
            this.empty3 = false;
            this.empty2 = false;
            return;
          }
          data.forEach((item: any) => {
            item.leaf = item.children_count > 0 ? false : true;
          });
          if (data[0].children_count > 0) {
            // 通过default-expanded-keys设置默认展开的节点。需要注意的是，
            // 此时必须设置node-key，其值为节点数据中的一个字段名，该字段在整棵树中是唯一的。
            self.expand_key.push(data[0].id);
          }
          self.click_id = 0;
          self.memberList(0);
          resolve(data);
        } else {
          this.$message.error(res.message);
          return;
        }
      });
    }
    if (node.level >= 1) {
      Api.depatListApi({ pid: node.data.third_id }).then((res: any) => {
        if (res.code == 200) {
          const data = res.data;
          if (data.length <= 0) {
            return;
          }
          if (data)
            data.map((item: any) => {
              if (!item.children_count) item.leaf = true;
            });
          resolve(data);
        } else {
          this.$message.error(res.message);
          return;
        }
      });
    }
  }
  // 选中部门 查询部门下的员工
  handleCheckChange(data: any) {
    this.click_id = data.id;
    this.meta.page = 1;
    this.memberList(0);
  }
  // 更新数据
  updateStaffInfo() {
    const _this = this;
    Api.updateRecordApi().then((res: any) => {
      if (res.code == 200) {
        _this.$message.success("数据更新完成");
        setTimeout(() => {
          _this.$router.go(0);
        }, 200);
      } else {
        _this.$message.error(res.message);
      }
    });
  }
  // 首次进来获取时间
  updateRocord() {
    Api.updateRocordApi().then((res: any) => {
      if (res.code == 200) {
        this.uptatedTime = res.data.updated_at;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 员工列表接口
  private memberList(type: any) {
    const params = {
      department_id: this.click_id,
      name: this.inputName,
      status: this.select,
      page: this.meta.page,
      size: this.meta.size
    };
    Api.memberListApi(params).then((res: any) => {
      if (res.code == 200) {
        this.staffList = res.data.member_list;
        // type == 1代表查询所调
        if (type == 1 && isNull(this.staffList)) {
          this.empty2 = true;
          this.empty3 = false;
        } else if (isNull(this.staffList)) {
          this.empty3 = true;
          this.empty2 = false;
        }
        this.meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 搜索查询
  searchStaff() {
    this.meta.page = 1;
    if (this.empty1 == true) {
      return;
    }
    this.memberList(1);
  }
  // 重置
  resetStaff() {
    this.inputName = "";
    this.meta.page = 1;
    this.select = 1;
    if (this.empty1 == true) {
      return;
    }
    this.memberList(0);
  }
  // 分页
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.memberList(0);
  }
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.memberList(0);
  }
  // 打开如何获取组织结构层级文档
  private openIntroduce() {
    window.open(
      "http://xueyuan.vchangyi.com/bin/view/CRS%E5%8A%9F%E8%83%BD%E8%AF%B4%E6%98%8E/%E4%BC%81%E4%B8%9A%E5%BE%AE%E4%BF%A1/%E8%AE%BE%E7%BD%AE/%E4%BC%81%E4%B8%9A%E5%BE%AE%E4%BF%A1%E9%80%9A%E8%AE%AF%E5%BD%95%E6%8E%88%E6%9D%83/",
      "_blank"
    );
  }
  // 查看活动详情数据
  private goDetail(id: any) {
    this.$router.push({
      name: "myapp_enterprise_organize-structure_staff-detail",
      query: {
        id: id
      }
    });
  }
  // 客户流失提醒通知配置
  private notifyConfig() {
    this.selectdpId = this.click_id;
    this.openNotify = true;
  }
  private closeDialog() {
    this.openNotify = false;
  }
  // ---------------------------------------------------------------生命周期---------------------------------------------------------------

  // 实例挂载完
  mounted() {
    this.updateRocord();
  }
}
</script>
<style lang="scss" scoped>
/deep/.layout-list .buttons-content {
  margin-bottom: 0;
}
// /deep/
//   .el-tree--highlight-current
//   .el-tree-node.is-current
//   > .el-tree-node__content {
//   background-color: none !important;
//   color: #708df5;
// }
.hover-name {
  cursor: pointer;
}
.hover-name:hover {
  color: $color-text-button;
}
.empty-box {
  text-align: center;
}
.empty-img {
  margin-top: 50px;
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
.custom-tree-node {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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
  box-shadow: none !important;
  background: $color-background;
}
.table-text-img {
  width: 25px;
  height: 25px;
  border-radius: 50%;
  vertical-align: middle;
  margin-right: 10px;
}
</style>
<style lang="scss">
#crs-scrm .supernatant {
  .el-tree--highlight-current
    .el-tree-node.is-current
    > .el-tree-node__content {
    background-color: #f5f6fa;
    color: #6f3112;
  }
}
.config-btn {
  margin-bottom: 20px;
}
</style>
