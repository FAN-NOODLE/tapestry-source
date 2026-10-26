<template>
  <div class="supernatant">
    <layout-list>
      <template slot="buttons">
        <el-button class="sync-btn" type="primary" @click="goSyncList()"
          >People Hub同步日志</el-button
        >
      </template>
      <template slot="search">
        <div class="top-content">
          <div class="top-item" v-if="false">
            <el-button type="primary" @click="searchStaff()"
              >从People Hub同步</el-button
            >
            <div class="right-text ">上次同步时间：2022年2月11日</div>
          </div>
          <!--  m-l-30 -->
          <div class="top-item">
            <el-button type="primary" v-if="epName != 'coach'" @click="syncChat(1)"
              >同步到企业微信</el-button
            >
            <div class="right-text ">
              上次同步时间：{{
                wechat_time | formatDate("yyyy-mm-dd hh:mm:ss")
              }}
            </div>
          </div>
        </div>
      </template>
      <template slot="slide">
        <div class="slide-list">
          <div class="search">
            <el-input
              clearable
              placeholder="请输入门店名称或编码"
              v-model="search.value"
            >
              <template slot="append">
                <el-select v-model="search.select" placeholder="请选择">
                  <el-option label="全部" :value="0"></el-option>
                  <el-option label="奥莱" :value="1"></el-option>
                  <el-option label="零售" :value="2"></el-option>
                </el-select>
                <el-button
                  type="primary"
                  icon="el-icon-search"
                  class="button-search"
                  @click="goSearch"
                ></el-button>
              </template>
            </el-input>
            <el-button
              type="primary"
              icon="el-icon-refresh"
              class="reset-search"
              @click="resetSearch"
            ></el-button>
          </div>
          <div class="table" v-if="storeList.length > 0">
            <div class="tr header">
              <div class="code">门店编号</div>
              <div class="name">门店名称</div>
              <div class="type">门店属性</div>
            </div>
            <div
              class="tr cursor"
              :class="item.id === staffStore.id ? 'is-select' : ''"
              v-for="(item, index) in storeList"
              :key="index"
              @click="clickStore(item)"
            >
              <div class="code ellipsis">{{ item.third_store_code }}</div>
              <div class="name ellipsis">{{ item.store_name }}</div>
              <div class="type ">
                <span v-if="item.store_type == 1">奥莱 </span>
                <span v-else-if="item.store_type == 2">零售 </span>
                <span v-else>--</span>
              </div>
            </div>
          </div>
          <div class="empty-table" v-if="storeList.length == 0">
            <img class="empty-img" src="@/assets/images/empty-box.png" alt="" />
            <div class="desc">暂无门店</div>
          </div>
        </div>
        <div class="footer" v-if="storeList.length > 0">
          <div class="prefix cur" @click="goPrePage">上一页</div>
          <div class="page">
            <span class="current-page">{{ search.page }}</span
            >/{{ search.total }}
          </div>
          <div class="append cur" @click="goAppendPage">下一页</div>
        </div>
      </template>
      <template slot="slide-table">
        <div class="table-top">
          <div class="store-info">
            <div class="store-name" v-if="staffStore.store_name">
              <div>{{ staffStore.store_name }}</div>
              <div class="tag" v-if="staffStore.store_type">
                {{ staffStore.store_type == 1 ? "奥莱店" : "零售店" }}
                <div class="tag-left"></div>
              </div>
            </div>
            <div class="store-code">{{ staffStore.third_store_code }}</div>
          </div>
          <div class="right-search">
            <div class="search">
              <el-input
                clearable
                placeholder="请输入员工姓名或手机号"
                v-model="searchData.mobile_or_name"
              >
                <template slot="append">
                  <el-select v-model="searchData.status" placeholder="请选择">
                    <el-option label="全部" :value="0"></el-option>
                    <el-option label="在职" :value="1"></el-option>
                    <el-option label="离职" :value="2"></el-option>
                  </el-select>
                  <el-button
                    type="primary"
                    icon="el-icon-search"
                    class="button-search"
                    @click="searchStaff"
                  ></el-button>
                </template>
              </el-input>
            </div>
            <el-button
              class="mar-l-30"
              size="small"
              type="primary"
              @click="resetStaffList()"
              >重置</el-button
            >
            <el-button
              class="mar-l-30"
              size="small"
              type="primary"
              @click="transferAll()"
              >批量转岗</el-button
            >
          </div>
        </div>
        <el-table :data="staffList" stripe>
          <el-table-column label="员工姓名" width="150" show-overflow-tooltip>
            <template slot-scope="scope">
              <img
                v-if="scope.row.avatar"
                :src="
                  scope.row.avatar
                    ? scope.row.avatar
                    : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                "
                alt=""
                class="table-text-img"
              />
              <span>{{ scope.row.name }}</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="employee_code"
            label="员工ID"
            min-width="120"
          ></el-table-column>
          <el-table-column
            prop="mobile"
            label="手机号"
            min-width="120"
            show-overflow-tooltip
          >
            <template slot-scope="scope">
              {{
                scope.row.mobile
                  ? scope.row.mobile.replace(
                      scope.row.mobile.substring(3, 7),
                      "****"
                    )
                  : "--"
              }}
            </template>
          </el-table-column>
          <el-table-column prop="gender" label="性别" min-width="70">
            <template slot-scope="scope">
              {{
                scope.row.gender == 1
                  ? "男"
                  : scope.row.gender == 2
                  ? "女"
                  : "未定义"
              }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="在职状态" min-width="100">
            <template slot-scope="scope">
              {{ scope.row.status == 1 ? "在职" : "离职" }}
            </template>
          </el-table-column>
          <el-table-column label="已同步到企微" width="130">
            <template slot-scope="{ row }">
              {{ row.sync_status === 1 ? "是" : "否" }}
            </template>
          </el-table-column>
          <el-table-column prop="job_location" label="工作地点" min-width="150">
          </el-table-column>
          <div slot="empty">
            <div class="table-empty-label">
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
    <edit-store
      v-if="storeModel"
      :storeInfo="storeInfo"
      @closeDialog="closeModel"
      @success="saveModel"
    ></edit-store>
    <tansfer-staff
      @close="tansferClose"
      @success="tansferSuccess"
      :selectList="selectList"
      v-if="tansfer_show"
    ></tansfer-staff>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import Api from "@/api/enterprise/business";
import ApiStore from "@/api/store-invite/index";
import EditStore from "./../component/edit-store.vue";
import TansferStaff from "./../component/tansfer-staff.vue";
@Component({
  components: {
    LayoutList,
    EditStore,
    TansferStaff
  }
})
export default class extends Vue {
  // ---------------------------------------------------------------此页面所需变量---------------------------------------------------------------
  // 左侧门店搜索
  private search = {
    value: "",
    select: 0,
    page: 1,
    total: 0
  };
  // 分页配置
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private searchData = {
    mobile_or_name: "",
    status: 0
  };
  private epName = "";
  // 控制编辑门店的弹框
  private storeModel = false;
  // 编辑的门店
  private storeInfo = {};
  // 右侧展示员工的门店
  private staffStore = {};
  // 员工列表
  private staffList: Array<any> = [];
  // 门店列表
  private storeList: Array<any> = [];
  private tansfer_show = false;
  private wechat_time = "";
  private selectList = [];
  // ---------------------------------------------------------------方法---------------------------------------------------------------
  private goSyncList() {
    this.$router.push({
      name: "myapp_enterprise_business-structure_sync-list"
    });
  }
  private syncChat(is_sync) {
    Api.syncWechatApi({ is_sync }).then((res: any) => {
      console.log(res, 122);
      if (is_sync === 1) {
        this.$message.success("更新成功");
      }
      this.wechat_time = (res.data.last_sync_time || 0) + "";
    });
  }
  // 关闭编辑门店弹框
  private closeModel() {
    this.storeModel = false;
  }
  // 打开编辑门店弹框
  private showEdit(item: any) {
    this.storeInfo = item;
    this.storeModel = true;
  }
  // 编辑门店成功的回调
  private saveModel() {
    // 更新数据
    this.getStoreList();
    this.closeModel();
  }
  private goPrePage() {
    if (this.search.page == 1) return;
    this.search.page -= 1;
    this.getStoreList();
  }
  private goAppendPage() {
    if (this.search.page == this.search.total) return;
    this.search.page += 1;
    this.getStoreList();
  }
  private goSearch() {
    this.search.page = 1;
    this.getStoreList();
  }
  private resetSearch() {
    this.search = {
      value: "",
      select: 0,
      page: 1,
      total: 0
    };
    this.getStoreList();
    this.staffStore = {};
    this.resetStaffList();
  }
  // 获取企业信息
  private getEpConfig(){
     ApiStore.getEpConfig().then(res => {
      // // 企业名称：coach;ks;sw
      const { name } = res.data;
      this.epName = name;      
    });
  }
  // 获取门店列表
  private getStoreList() {
    const params = {
      store_type: this.search.select || "",
      code_or_name: this.search.value,
      page: this.search.page,
      size: 20
    };
    // Math.ceil()
    this.search.total = 3;
    Api.getStoreListApi(params).then(res => {
      const { total, list } = res.data;
      this.storeList = list;
      this.search.total = Math.ceil(total / 20);
    });
  }
  private clickStore(item: any) {
    this.staffStore = item;
    this.meta.page = 1;
    this.getStaffList();
  }
  // ------------员工操作---------
  // 批量转岗
  private transferAll() {
    this.selectList = this.staffList.filter(item => item.status === 1);
    this.tansfer_show = true;
  }
  private tansferClose() {
    this.tansfer_show = false;
  }
  private tansferSuccess() {
    this.tansferClose();
    this.searchStaff();
  }
  // 搜索
  private searchStaff() {
    this.meta.page = 1;
    this.getStaffList();
  }
  // 分页
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getStaffList();
  }
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getStaffList();
  }
  // 重置员工列表
  private resetStaffList() {
    this.meta.page = 1;
    this.searchData = {
      mobile_or_name: "",
      status: 0
    };
    this.getStaffList();
  }
  // 获取员工列表
  private getStaffList() {
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      third_store_code: (this.staffStore as any).third_store_code || "",
      no_department: 0,
      ...this.searchData
    };
    Api.getStaffListApi(params).then(res => {
      const { list, ...meta } = res.data;
      this.staffList = list;
      this.meta = meta;
    });
  }
  // ---------------------------------------------------------------生命周期---------------------------------------------------------------

  // 实例挂载完
  mounted() {
    this.getStoreList();
    this.getStaffList();
    this.syncChat(0);
    this.getEpConfig();
  }
}
</script>
<style lang="scss" scoped>
/deep/.layout-list .buttons-content {
  margin-bottom: 0;
}
/deep/.layout-list .slide-table {
  width: 70%;
}
/deep/.layout-list .search-content {
  background: white;
  padding-left: 0;
  margin-bottom: 5px;
}
/deep/.layout-list .buttons-content {
  display: flex;
  justify-content: flex-end;
}
/deep/.layout-list .slide-menu {
  width: 30%;
  flex-shrink: 0;
  background: #f7f7f7;
  position: relative;
}
.empty-table {
  height: 300px;
  display: flex;
  justify-content: center;
  flex-direction: column;
  align-items: center;
}
.search {
  /deep/.el-select {
    width: 70px;
    margin: 0;
    position: absolute;
    top: 0;
    left: 0;
    .el-input__inner {
      padding-left: 10px;
    }
  }
  .button-search {
    width: 35px;
    margin: 0;
    position: absolute;
    top: 0;
    right: 0;
    background: $color-main;
    color: white;
    padding: 12px 0;
    cursor: pointer;
  }
}
/deep/.el-input-group__append {
  padding: 0 52px;
  background: white;
}
.m-l-30 {
  margin-left: 30px;
}
.top-content {
  display: flex;
  justify-content: space-between;
  .top-item {
    width: 50%;
    height: 92px;
    display: flex;
    justify-content: flex-start;
    align-items: center;
    background: #f7f7f7;
    flex-shrink: 0;
    box-sizing: border-box;
    padding: 25px 28px;

    .right-text {
      font-size: 14px;
      color: #c0c4cc;
      margin-left: 20px;
    }
  }
}
.slide-list {
  padding-top: 15px;
  position: relative;
  min-height: 450px;
  .search {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 15px;
  }
  .reset-search {
    height: 38px;
    width: 35px;
    padding-left: 10px;
  }
}
.table {
  padding: 20px 0 50px;
  color: #606266;
  font-size: 14px;
  .pad-l-40 {
    padding-left: 40px;
    flex-shrink: 0;
  }
  .header {
    background: #ebebeb;
    font-weight: 500;
  }
  .cursor {
    cursor: pointer;
  }
  .tr {
    display: flex;
    justify-content: space-around;
    flex-wrap: nowrap;
    align-items: center;
    line-height: 30px;
    padding: 0 15px;
    &.is-select {
      background: #ebebeb;
      color: #6f3112;
    }
    &:hover {
      background: #ebebeb;
    }
    .code {
      width: 100px;
      flex-shrink: 0;
    }
    .name {
      width: 100%;
    }
    .ellipsis {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .type {
      width: 60px;
      flex-shrink: 0;
    }
  }
}
.footer {
  display: flex;
  justify-content: center;
  align-items: center;
  color: #1a1a1a;
  width: 100%;
  position: absolute;
  bottom: 15px;
  left: 0;
  .page {
    margin: 0 20px;
    .current-page {
      color: #6f3112;
    }
  }
  .cur {
    cursor: pointer;
  }
}
.table-top {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
  min-height: 55px;
  .right-search {
    width: 530px;
    margin-left: 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    .search {
      width: 320px;
    }
    .mar-l-30 {
      margin-left: 30px;
    }
  }
  .store-info {
    .tag {
      width: 45px;
      height: 21px;
      position: relative;
      background: #000;
      font-size: 12px;
      font-family: PingFangSC-Regular, PingFang SC;
      font-weight: 400;
      color: #ffffff;
      line-height: 21px;
      text-align: center;
      margin-left: 8px;
      flex-shrink: 0;
      &::after {
        position: absolute;
        content: "";
        clear: both;
        width: 2px;
        height: 2px;
        background: #fff;
        border-radius: 50%;
        top: 10px;
        right: -3px;
        z-index: 2;
      }
    }
    .tag-left {
      position: absolute;
      top: 0;
      right: -11px;
      width: 0;
      height: 0;
      border-top: 11px solid transparent;
      border-left: 11px solid #000;
      border-bottom: 10px solid transparent;
    }
    .store-name {
      font-size: 16px;
      font-family: PingFangSC-Medium, PingFang SC;
      font-weight: 500;
      color: #414244;
      line-height: 20px;
      display: flex;
      justify-content: flex-start;
    }
    .store-code {
      font-size: 14px;
      font-weight: 400;
      color: #c0c4cc;
      line-height: 14px;
      margin-top: 10px;
    }
  }
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
