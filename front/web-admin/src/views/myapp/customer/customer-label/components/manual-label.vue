<template>
  <div>
    <el-button
      class="m-b-10"
      plain
      size="small"
      v-if="activeTab == 'manual'"
      @click="exportList"
    >导出Excel</el-button>
    <!-- todo：新增标签 -->
    <el-button
      class="add-btn"
      type="primary"
      size="small"
      @click="addLabelFn"
      v-if="activeTab == 'manual'"
    >新增标签</el-button>
    <div class="tips">温馨提示：接口打标签需要一定时间，新增标签人群上传后，【标签人数】更新根据数据量会一定的时间延迟，请多多理解耐心等待哦！（参考速度：10分钟约打2W客户）</div>
    <div class="layout-label">
      <dl class="menu-left">
        <div class="dt-label">
          <div>标签组名称</div>
        </div>
        <!-- 标签组列表模块-->
        <div
          class="sub-wrap"
          v-if="groupList.length > 0"
        >
          <transition-group name="flip-list">
            <dd
              v-for="(item, index) in groupList"
              :key="'_' + index"
              @click="setMenuFn(item, index)"
              :class="{ menuAactive: menuIndex == index + 1 }"
            >
              <span>{{ item.tag_group_name }}</span>

              <!-- <img
                @click="setMenuFn(item)"
                src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/customer/drag_icon.png"
                class="group-drag-icon"
                alt=""
              /> -->
              <div
                @click.stop="deleteGroup(item.tag_group_id)"
                class="group-delete-icon el-icon-delete"
              ></div>
            </dd>
          </transition-group>
        </div>
        <div
          class="grouplist-null"
          v-if="groupList.length == 0"
        >
          暂未添加标签组
        </div>
      </dl>
      <!-- 标签列表模块-->
      <div class="list-right">
        <div class="font-middle">
          {{ curMenuItem && curMenuItem.tag_group_name }}
        </div>
        <!--table-->
        <el-table
          v-if="showTable"
          :data="labelList"
          style="width: 100%"
          max-height="500"
          stripe
          @sort-change="sortChange"
        >
          <el-table-column
            prop="tag_name"
            label="标签名称"
          > </el-table-column>
          <el-table-column
            prop="value"
            label="标签人数"
            sortable="custom"
          >
            <template slot-scope="{ row }">
              <el-button
                type="text"
                @click.native.prevent="goCustomerList(row,activeTab)"
              >
                {{ row.value }}
              </el-button>
            </template>
          </el-table-column>
          <!-- todo: 标签有效期 -->
          <el-table-column
            prop="tag_validity"
            label="标签有效期"
            v-if="activeTab == 'manual'"
          >
            <template slot-scope="{ row }">
              <span v-if="row.expired_time == 0">长期有效</span>
              <span v-else>{{ row.create_time | formatDate("yyyy-mm-dd hh:mm:ss") }}~{{ row.expired_time |
                formatDate("yyyy-mm-dd hh:mm:ss") }}</span>
            </template>
          </el-table-column>
          <el-table-column
            fixed="right"
            label="操作"
            width="200"
          >
            <template slot-scope="scope">
              <el-button
                @click.native.prevent="goDetail(scope.row)"
                type="text"
              >
                查看
              </el-button>
              <el-button
                @click="editRow(scope.row)"
                type="text"
              >
                编辑
              </el-button>
              <el-button
                type="text"
                slot="reference"
                class="color-danger m-l-20"
                @click="deleteRow(scope.row.tag_id)"
              >删除</el-button>
            </template>
          </el-table-column>
          <div
            slot="empty"
            class="table-empty-label"
          >
            <img
              :src="emptyLabel.img"
              alt=""
            />
            <div>{{ emptyLabel.text }}</div>
          </div>
        </el-table>
        <!-- todo:增加分页 -->
        <el-pagination
          style="text-align:right;margin-top: 20px;"
          v-if="labelList && labelList.length > 0 && showTable"
          @current-change="handleCurrentChange"
          :current-page="meta.page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="meta.size"
          class="page-right font-weight-400"
          layout="total, prev, pager, next, jumper"
          :total="meta.total"
        ></el-pagination>
      </div>
    </div>
    <!-- 新增标签组弹框-->
    <add-manual-group
      v-if="showAddGroup"
      @closeDialog="showAddGroup = false"
      @submit="addGroupSubmit"
    ></add-manual-group>
    <!-- 新增、编辑标签弹框-->
    <add-manual-label
      v-if="showAddLable"
      :currentEditItem="currentEditItem"
      @closeDialog="showAddLable = false"
      @submit="addLabelSubmit"
    ></add-manual-label>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop, Watch } from "vue-property-decorator";
import addManualLabel from "./add-manual-label.vue";
import addManualGroup from "./add-manual-group.vue";
import API from "@/api/customer";
// todo: 获取企业配置
import Api from "@/api/store-invite/index";
@Component({
  components: {
    addManualLabel,
    addManualGroup
  }
})
export default class extends Vue {
  @Prop() private group: any;
  @Prop({ default: "" }) private keywords: string; // 搜关键词
  @Prop({ default: 1 }) private labelId: any; // 搜关键词
  @Prop({ default: false }) private authStatus: boolean;
  @Prop({ default: "" }) private activeTab: string;
  @Watch("activeTab", { immediate: true })
  onActiveTabChange() {
    // 清空分页参数
    this.meta = {
      total: 0, // 总数
      page: 1, // 当前页
      size: 6 // 每页数量
    } 
  }
  private showAddGroup = false; // 控制新增标签组弹框显隐
  private showAddLable = false; // 控制新增标签弹框显隐
  private groupList: any = []; // 左侧menu-标签组列表
  private labelList: any = []; // 右侧标签列表
  private menuIndex = 1; // 当前选中的标签组下标值
  private curMenuItem: any = {}; // 当前选中的菜单按钮-默认选中第menuIndex==1个
  private currentEditItem: object = {}; // 当前要编辑的标签内容
  private oldNum: any; // 标签组排序数组 老
  private newNum: any; // 标签组排序数组 新
  private sortIdArray: any = []; // 标签组排序数组 新
  private order = 0; // 标签列表顺序-1:降序，2:升序， 默认降序
  private emptyLabel = {
    text: "暂未添加标签",
    img:"https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/default/default_noList.png"
  };
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 6 // 每页数量
  };
  private startY = 0;
  private endY = 0;
  private showTable = false;
  // 企业名称
  private epName = "";
  // todo:获取企业配置
  getCompanyConfig() {
    Api.getEpConfig().then(res => {
      // // 企业名称：coach;ks;sw
      const { name } = res.data;
      this.epName = name;
    });
  }
  /************methods**************/
  // todo:新增标签
  addLabelFn() {
    this.showAddLable = true;
    this.currentEditItem = {}
  }
  // todo:新增标签提交
  async addLabelSubmit() {
    this.showAddLable = false;
    await this.getGroupList();
  }
  private goDetail(row) {
    this.$router.push({
      name: "myapp_customer_customer-label_detail",
      query: {
        tag_id: row.tag_id,
        tag_name: row.tag_name,
        group_name: this.curMenuItem.tag_group_name,
        tag_type: this.activeTab === 'manual' ? '0' : '1',
      }
    });
  }
  // 点击导出excel按钮
  private exportList() {
    API.tagExportApi({
      tag_type: 2
    }).then((res: any) => {
      const blob = new Blob([res], {
        type: "application/x-xls;charset=utf-8"
      });
      let fileName = "";
      const nowTime = new Date();
      fileName = "客户标签_" + this.formateDate(nowTime) + ".xlsx";

      //下载的文件名称及其后缀，后缀要和后台保持的一致

      if ("download" in document.createElement("a")) {
        // 非IE下载
        const elink = document.createElement("a");
        elink.download = fileName;
        elink.style.display = "none";
        elink.href = URL.createObjectURL(blob);
        document.body.appendChild(elink);
        elink.click();
        URL.revokeObjectURL(elink.href); // 释放URL 对象
        document.body.removeChild(elink);
      } else {
        // IE10+下载
        (navigator as any).msSaveBlob(blob, fileName);
      }
    });
  }
  // 时间格式转换
  private formateDate(time: any) {
    const year = time.getFullYear();
    let month = time.getMonth() + 1;
    let day = time.getDate();

    month = month < 10 ? "0" + month : month;
    day = day < 10 ? "0" + day : day;

    return year + "-" + month + "-" + day;
  }
  // 获取标签组列表接口
  private getGroupList(last?: boolean) {
    API.labelListApi({}).then((res: any) => {
      if (res.code === 200) {
        if (res.data.manual && res.data.manual.length > 0) {
          this.groupList = res.data.manual;
          // 新增完标签组后光标定位到最后一项
          if (last) {
            this.menuIndex = res.data.manual.length;
          }
          // 修正删除标签组后的默认光标点
          if (this.menuIndex > this.groupList.length) {
            this.menuIndex = this.groupList.length;
          }
          this.curMenuItem = this.groupList[this.menuIndex - 1];
          this.getSubListFn();
        } else {
          this.groupList = [];
          this.curMenuItem = {};
          this.labelList = [];
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getSubListFn();
  }
  // 获取标签列表接口
  private getSubListFn() {
    setTimeout(() => {
      this.showTable = true;
    }, 50);
    const params: any = {
      tag_group_id: this.curMenuItem && this.curMenuItem.tag_group_id,
      page: this.meta.page,
      size: this.meta.size
    };
    let api = "impressionTagListApi";
    if (this.activeTab === "manual") {
      api = "tagInfoApi";
      params.tag_keywords = this.labelId === 2 ? this.keywords.trim() : "";
      params.order = this.order;
    } else {
      params.keyword = this.labelId === 2 ? this.keywords.trim() : "";
    }
    API[api](params).then((res: any) => {
      if (res.code === 200) {
        res.data.list.forEach((item: any) => {
          if (!item.value) {
            item.value = item.use_num || 0;
          }
          if (this.activeTab !== "manual") {
            item.tag_name = item.name;
            item.tag_id = item.id;
          }
        });
        //获取子标签列表
        this.labelList = res.data.list;
        this.meta.total = res.data.total;
        this.$emit("tagSonCase", {
          _this: this,
          data: res.data.value_update_at
        });
      } else {
        this.$message.error(res.message);
      }
    });
  }
  private setMenuFn(item: object, index: number) {
    this.curMenuItem = item;
    this.menuIndex = index + 1;
    this.order = 0;
    this.showTable = false;
    // todo:清空分页参数
    this.meta = {
      total: 0, // 总数
      page: 1, // 当前页
      size: 6 // 每页数量
    };
    this.getSubListFn();
  }
  // 记录初始信息
  private dragStartFn(value: object, $event: any) {
    this.oldNum = value;
    this.startY = $event.clientY;
  }
  private dropFn($event: any) {
    //end
    $event.preventDefault();
    this.endY = $event.clientY;
  }
  private allowDropFn($event: any) {
    //moving
    $event.preventDefault();
    this.endY = $event.clientY;
  }
  // 做最终操作
  private dragEndFn(item: any, $event: any) {
    // 保存当前选中的type
    const tempId = this.groupList[this.menuIndex - 1].tag_group_id;
    $event.currentTarget.className = "";
    if (this.oldNum != this.newNum) {
      const oldIndex = this.groupList.indexOf(this.oldNum);
      const newIndex = this.groupList.indexOf(this.newNum);
      const newItems = [...this.groupList];
      // 删除老的节点
      newItems.splice(oldIndex, 1);
      // 在列表中目标位置增加新的节点
      newItems.splice(newIndex, 0, this.oldNum);
      // this.items一改变，transition-group就起了作用
      this.groupList = [...newItems];
    }
    //判断移动之后当前的选中
    this.groupList.forEach((item: any, index: any) => {
      if (item.tag_group_id == tempId) {
        this.menuIndex = index + 1;
      }
    });
    let curIndex = 0;
    // if (this.menuIndex == index + 1) {
    //   $event.currentTarget.className = "menuAactive"
    // }
    //组装排序接口的数组数据
    this.sortIdArray = [];
    this.groupList.forEach((menuItem: any, index: number) => {
      this.sortIdArray.push(menuItem.tag_group_id);
      if (menuItem.tag_group_name == item.tag_group_name) {
        curIndex = index;
      }
    });
    //排序接口
    this.sortMenuFn();
    // 重新获取标签信息
    this.setMenuFn(item, curIndex);
  }
  // 标签组排序接口
  private sortMenuFn() {
    const params = {
      tag_group_first_type: 2,
      tag_group_ids: this.sortIdArray
    };
    API.tagRankApi(params).then((res: any) => {
      if (res.code !== 200) {
        this.$message.error(res.message);
      }
    });
  }
  //记录移动过程中信息
  private dragEnterFn(value: object, $event: any) {
    $event.preventDefault();
    this.newNum = value;
  }
  //鼠标按下事件
  private mousedownFn($event: any) {
    $event.currentTarget.className = "menuEnter";
  }
  //鼠标抬起事件
  private mouseupFn($event: any, index: number) {
    if (index == this.menuIndex - 1) {
      $event.currentTarget.className = "menuAactive";
    } else {
      $event.currentTarget.className = "";
    }
  }
  // 标签组列表单行-删除按钮
  private deleteGroup(tag_group_id: number) {
    this.$confirm("是否需要删除该标签组", "删除确认", {
      confirmButtonText: "确认删除",
      confirmButtonClass: "button-danger",
      cancelButtonText: "我再想想",
      type: "warning"
    })
      .then(() => {
        const params = {
          tag_group_id
        };
        API.tagGroupDeleteApi(params).then((res: any) => {
          if (res.code === 200) {
            this.$message({
              message: "删除成功",
              type: "success",
              duration: 2000
            });
            this.getGroupList();
          } else {
            this.$message.error(res.message);
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
  // 标签列表单行-编辑按钮
  private editRow(item: any) {
    this.showAddLable = true;
    this.currentEditItem = item;
  }
  // 标签列表单行-删除按钮
  private deleteRow(id: number) {
    this.$confirm("删除后，客户身上的标签信息将被删除，且企业微信后台该标签也将同步删除？", "确认删除", {
      confirmButtonText: "确定",
      confirmButtonClass: "button-danger",
      cancelButtonText: "取消",
      type: "warning"
    }).then(() => {
      const params = {
        tag_id: id
      };
      //删除子标签接口
      API.tagDeleteRowApi(params).then((res: any) => {
        if (res.code === 200) {
          this.$message({
            message: "删除成功",
            type: "success",
            duration: 2000
          });
          this.getSubListFn();
          this.getGroupList();
        } else {
          this.$message.error(res.message);
        }
      });
    });
  }
  // 标签列表排序
  private sortChange(obj: any) {
    // 点击列排序时
    //1:降序，2:升序， 默认降序
    this.order = 0;
    if (obj.order == "descending") {
      //降序
      this.order = 1;
    } else if (obj.order == "ascending") {
      //升序
      this.order = 2;
    }
    if (this.groupList.length > 0) {
      this.getSubListFn();
    }
  }
  // 打开新增标签组弹框
  private openAddManualGroup() {
    this.showAddGroup = true;
  }
  // 新增标签组弹框成功回调
  private addGroupSubmit() {
    this.getGroupList(true);
    this.showAddGroup = false;
  }
  private goCustomerList(row: any,activeTab: string) {
    this.$router.replace({
      name: "myapp_customer_customer-list",
      query: {
        tag_id: row.tag_id,
        tag_type:activeTab
      }
    });
  }
  // 打开新增标签弹框
  private openAddManualLabel() {
    if (this.groupList.length == 0) {
      this.$message.warning("请先新建标签组");
      return;
    }
    if (this.labelList.length >= 20) {
      this.$message.warning("该标签组已添加满20条标签上限");
      return;
    }
    this.currentEditItem = {};
    this.showAddLable = true;
  }
  @Watch("group")
  onChangeLabel() {
    this.groupList = this.group;
    //默认请求第一个详情
    if (this.groupList.length > 0) {
      this.curMenuItem = this.groupList[0];
      this.setMenuFn(this.groupList[0], 0);
    } else {
      this.showTable = true;
    }
  }
  private created() {
    if (this.group.length) {
      this.groupList = this.group;
      this.curMenuItem = this.groupList[this.menuIndex - 1];
      if (!this.showTable) {
        this.getSubListFn();
      }
    }
  }
  private mounted() {
    this.getCompanyConfig();
  }
}
</script>
<style lang="scss" scoped>
.tips{
  font-size: 14px;
  color: red;
  margin-bottom: 10px;
}
.layout-label {
  display: flex;

  // 左边的标签组列表
  .menu-left {
    flex: 0 0 400px;
    margin: 0;
    margin-right: 30px;
    background: $color-background;

    .group-delete-icon {
      width: 14px;
      height: 14px;
      margin-left: 20px;
      float: right;
      position: relative;
      top: -16px;
      display: none;
    }

    .group-drag-icon {
      width: 14px;
      height: 14px;
      margin-left: 20px;
      margin-right: 20px;
      float: right;
      position: relative;
      top: -16px;
      display: none;
    }

    .dt-label {
      background: #ebebf1;
      font-size: 14px;
      color: $color-text-main;
      height: 36px;
      line-height: 36px;
      padding: 0 20px;
      font-weight: bold;
      display: flex;
      justify-content: space-between;

      // 新增标签组按钮
      .add-group-btn {
        cursor: pointer;
        display: flex;
        justify-content: space-between;
        font-size: 14px;
        color: $color-main;

        ::v-deep.el-icon-circle-plus {
          color: $color-main;
          cursor: pointer;
        }

        .add-label-tip {
          color: $color-text-second;
        }
      }

      .add-group-btn:hover {
        color: #6f3112;

        ::v-deep.el-icon-circle-plus {
          color: #6f3112;
          cursor: pointer;
        }
      }

      .add-group-btn:active {
        color: #6f3112;

        ::v-deep.el-icon-circle-plus {
          color: #6f3112;
          cursor: pointer;
        }
      }
    }

    .sub-wrap {
      height: 600px;
      overflow-y: auto;
    }

    dd {
      margin-inline-start: 0px;
      height: 20px;
      line-height: 20px;
      color: $color-text-main;
      font-size: 14px;
      padding: 10px 0;
      cursor: pointer;

      span {
        display: block;
        padding-left: 20px;
        border-left: 4px solid $color-background;
      }
    }

    dd:hover {
      background: #ebebf1;

      span {
        border-left: 4px solid #ebebf1;
        margin-right: 20px;
        // background: url("https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/customer/drag_icon.png") no-repeat
        //   right center;
        background-size: 14px;
      }

      .group-delete-icon,
      .group-drag-icon {
        display: block;
      }
    }

    dd.menuAactive {
      background: $color-background;
      position: relative;

      &::before {
        content: "";
        display: block;
        width: 5px;
        height: 14px;
        background: #000000;
        border-radius: 0px 2px 2px 0px;
        position: absolute;
        left: 0;
        top: 50%;
        transform: translateY(-50%);
      }

      span {
        // border-left: 4px solid $color-main;
        color: $color-main;
        margin-right: 20px;
      }
    }

    dd.menuEnter {
      background: $color-main;

      span {
        border-left: 4px solid $color-main;
        color: $color-white;
        margin-right: 20px;
        background: none;
        //   center;
        // background-size: 20px;
      }
    }
  }

  .list-right {
    flex: 1;
    overflow-x: hidden;

    ::v-deep.el-table__fixed-right {
      right: 0 !important;
    }

    ::v-deep.gutter {
      width: 0px !important;
      display: none;
    }

    ::v-deep.el-table__empty-block {
      min-height: 300px;
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

    .font-middle {
      color: $color-text-main;
      font-size: 14px;
      line-height: 20px;
      padding-bottom: 20px;
      font-weight: 600;

      span {
        color: $color-text-second;
        font-size: 12px;
      }
    }

    .operation-btn {
      padding-bottom: 20px;

      ::v-deep.el-button--small {
        font-size: 14px;
        margin-right: 20px;
      }
    }
  }

  .grouplist-null {
    color: #909399;
    text-align: center;
    height: 399px;
    line-height: 399px;
  }
}

.el-icon-warning {
  color: #faad14;
}
</style>
