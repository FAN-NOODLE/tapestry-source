<template>
  <div>
    <layout-list>
      <!-- 顶部用户查询部分 -->
      <template slot="search">
        <search-item text="客户昵称：">
          <el-input
            class="search-input"
            size="medium"
            v-model.trim="nick_name"
            placeholder="请输入客户昵称"
          >
          </el-input>
        </search-item>
        <search-item text="客户手机号：">
          <el-input
            class="search-input"
            size="medium"
            v-model.trim="mobile"
            placeholder="请输入客户手机号"
          >
          </el-input>
        </search-item>
        <div>
          <el-button
            type="primary"
            @click="userSearch"
            size="small"
          >查询</el-button>
          <el-button
            plain
            @click="userSearchReset"
            size="small"
          >重置</el-button>
        </div>
      </template>

      <!-- 客户列表tabs标签页 -->
      <template slot="tabs">
        <el-tabs
          v-model="activeName"
        >
          <el-tab-pane
            label="全部客户"
            name="first"
            v-if="tagControl.show_shop == true || tagControl.show_official == true
              "
          ></el-tab-pane>
        </el-tabs>
      </template>
      <!-- 客户列表内容 -->
      <template slot="table">
        <div class="user-table">
          <el-table
            :data="lists"
            stripe
            ref="tableList"
            style="width:100%"
            :default-sort="{ prop: 'pay_order_num', order: 'descending' }"
            :cell-style="{ height: '20px' }"
          >
            <el-table-column
              prop="union_id"
              label="union_id"
              min-width="250"
            >
              <template slot-scope="scope">
                <div>{{ scope.row.union_id || "--" }}</div>
              </template>
            </el-table-column>
            <el-table-column
              label="微信"
              min-width="180"
              fixed="left"
              prop="name"
            >
              <template slot-scope="scope">
                <div>
                  <img
                    :src="scope.row.avatar"
                    class="m-r-10 user-idImg"
                    alt=""
                  />

                  <el-popover
                    placement="top-start"
                    trigger="hover"
                    :disabled="scope.row.nick.length <= 10"
                  >
                    <div>{{ scope.row.nick ? scope.row.nick : "--" }}</div>
                    <span
                      slot="reference"
                      v-if="scope.row.nick.length <= 10"
                      class="user-idName"
                      @click="routerlink(scope.row.user_id)"
                    >{{ scope.row.nick || "--" }}</span>
                    <span
                      slot="reference"
                      v-if="scope.row.name.length > 10"
                      class="user-idName"
                      @click="routerlink(scope.row.user_id)"
                    >{{ scope.row.nick.substr(0, 10) + "..." || "--" }}</span>
                  </el-popover>
                </div>
              </template>
            </el-table-column>
            <el-table-column
              prop="member_name"
              label="客户姓名"
              min-width="100"
            >
              <template slot-scope="scope">
                <span v-if="scope.row.name">{{ scope.row.name }}</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
            <el-table-column
              prop="mobile"
              label="手机号"
              min-width="150"
            >
              <template slot-scope="scope">
                <span v-if="scope.row.mobile">{{
                  scope.row.mobile.substr(0, 3) +
                  "****" +
                  scope.row.mobile.substr(7, 11)
                }}</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
            <el-table-column
              prop="birthday"
              label="生日"
              min-width="140"
            >
              <template slot-scope="scope">
                <div>{{ scope.row.birthday || "--" }}</div>
              </template>
            </el-table-column>
            <el-table-column
              prop="user_id"
              label="操作"
              min-width="180"
              fixed="right"
            >
              <template slot-scope="scope">
                <el-button
                  style="color:#303133"
                  @click="routerlink(scope.row.id)"
                  type="text"
                  size="small"
                >查看详情</el-button>
              </template>
            </el-table-column>
            <div
              slot="empty"
              class="table-empty-label"
            >
              <img
                v-if="searchLoading"
                src="@/assets/images/empty-search.png"
                alt=""
              />
              <img
                v-else
                src="@/assets/images/empty-box.png"
                alt=""
              />
              <p>{{ searchLoading ? "没有找到相关客户" : "暂无客户" }}</p>
            </div>
          </el-table>
        </div>
      </template>
      <!-- 底部分页器 -->
      <template slot="footer">
        <el-pagination
          v-if="lists && lists.length > 0"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="meta.page"
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
import API from "@/api/customer";
@Component({
  components: {
    // LayoutList
    // SelectInput,
    // DoubleSelect
  }
})
export default class extends Vue {
  /*********************************************data开始**************************************/
  private activeName = "first"; // tabs标签页默认开始位置
  private customerTagsList: number[] = []; // 客户标签选项
  private impressionTagsList: number[] = []; // 印象标签选项
  private userId: Array<any> = []; // 用户Id

  /**** Tags标签显示 ****/
  private tagControl = {
    show_official: true,
    show_shop: true
  };
  private userCustomerLableList: Array<any> = []; // 客户标签
  userCustomerLableHave = []; // 客户标签分页
  userCustomerGruopList: Array<any> = []; // 客户分组
  private nick_name = ""; // 客户昵称
  private mobile = ""; // 客户手机号
  private userLableSearchType = true; //查询类型;
  private searchLoading = false; // 是否正在搜索
  // 表格数据
  private lists: Array<any> = []; // table数据
  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 0
  };
  // 标签和分组分页
  private metas = {
    size: 500,
    page: 1,
    total: 0
  };
  private keyword = "";
  private tagMeta = {
    size: 10,
    page: 1,
    total: 0
  };
  /*********************************************data结束**************************************/
  /*********************************************methods开始***********************************/
  //客户标签分组API接口
  private getTagList() {
    console.log("11客户标签");
    const params = {
      type: 1,
      page: this.tagMeta.page,
      size: this.tagMeta.size,
      search: this.keyword
    };
    API.customerTagList(params).then((res: any) => {
      if (res.code === 200) {
        if (res.data.list) {
          if (this.tagMeta.page == 1) {
            const obj = { id: 0, name: "全部" };
            res.data.list.unshift(obj);
          }
          const { list, ...meta } = res.data;
          this.userCustomerLableList = [...this.userCustomerLableList, ...list];
          this.tagMeta = meta;
        }
      }
    });
  }
  //内部绑定函数（获取每页数据量）
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.tableDataList();
  }
  //内部绑定函数（获取当前第几页）
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.tableDataList();
  }
  //查询
  private userSearch() {
    this.meta = {
      size: 10,
      page: 1,
      total: 0
    };
    // 是否正在搜索
    this.searchLoading = true;
    this.tableDataList();
  }
  //重置按钮
  private userSearchReset() {
    //内容重置
    this.searchLoading = false;
    this.userLableSearchType = true;
    this.customerTagsList = [];
    this.nick_name = "";
    this.mobile = "";
    this.meta = {
      size: 10,
      page: 1,
      total: 0
    };
    this.tableDataList();
  }
  //跳转
  private routerlink(item: string) {
    //用户详情路由跳转
    this.$router.push({
      name: "myapp_customer_customer-list_detail",
      query: {
        user_id: item
      }
    });
  }
  //表格内容API接口
  private tableDataList() {
    const params = {
      tag_list: this.customerTagsList.filter(item => item),
      impression_tag_list: this.impressionTagsList.filter(item => item),
      page: this.meta.page,
      size: this.meta.size,
      name: this.nick_name,
      mobile: this.mobile
    };
    API.customerListApi(params).then((res: any) => {
      if (res.code === 200) {
        /* this.meta = {
          page: res.data.page,
          size: res.data.size,
          total: res.data.total
        }; */
        this.meta.total = res.data.total;
        this.lists = res.data.list;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  private initParams() {
    console.log(this.$route.query, "router");
    const { query } = this.$route;
    if (query.tag_id && query.tag_type) {
      this.getTagList();
      if (query.tag_type == 'manual'){
        this.customerTagsList = [Number(query.tag_id)];
      } else {
        this.impressionTagsList = [Number(query.tag_id)];
      }
    }
    this.tableDataList();
    !query.tag_id ? this.getTagList() : "";
  }
  /*********************************************methods开始***********************************/
  // 生命周期
  private mounted() {
    this.initParams();
    // 解决表格错位
    this.$nextTick(() => {
      (this.$refs.tableList as any).doLayout();
    });
  }
}
</script>

<style lang="scss" scoped>
.search-input {
  width: 280px !important;
}

.user-idImg {
  width: 20px;
  height: 20px;
  border: none;
}

.user-idName {
  position: relative;
  top: -4px;
  cursor: pointer;
  border: none;
}

/deep/.el-tabs__nav {
  .is-top {
    color: $color-text-normal;
  }
}

.tabs-font {
  font-size: 18px;
  color: #303133;
  font-weight: normal;
}

select {
  border-radius: 3px;
  height: 28px;
}

.el-input__count {
  padding-top: 20px;
}

/deep/.slide-box {
  position: absolute;
}

/deep/.cell {
  .user-LableButtonList {
    display: flex;

    .user-groupbutton {
      font-size: 12px;
      width: 52pt;
      height: 24pt;
      line-height: 24pt;
      background: #fef8ea;
      color: #f6ab3c;
      border-radius: 3pt;
      box-sizing: border-box;
      border: 1px solid #fce5c4;
      overflow: hidden;
      cursor: pointer;
      margin-right: 5px;
      text-align: center;
    }
  }

  .user-LableButtonLists {
    display: flex;
  }

  .user-LableButtonList {
    display: flex;
    margin-left: 5px;

    .user-lablebutton {
      font-size: 12px;
      width: 52pt;
      height: 24pt;
      line-height: 24pt;
      background: #f0f3fe;
      color: #708df5;
      border-radius: 3pt;
      box-sizing: border-box;
      border: 1px solid #d4dcfc;
      overflow: hidden;
      cursor: pointer;
      margin-right: 5px;
      text-align: center;
    }
  }
}

.customerGroups-list {
  margin-top: 5px;
  margin-left: 5px;
  display: flex;
}

.user-group {
  /deep/ .el-input__inner {
    height: 36px;
    margin-right: 0px;
  }
}

.cell {
  display: flex;
}

.user-label {
  /deep/ .el-input__inner {
    height: 36px;
    margin-right: 0px;
  }
}

.el-table_1_column_13 .is-hidden {
  .cell {
    display: none;
  }
}

.table-empty-label {
  text-align: center;
  height: 340px;

  img {
    width: 124px;
    display: block;
    margin: 80px auto 0;
  }
}

/deep/.el-table__fixed-right {
  right: 0 !important;
}
</style>
<style lang="scss" scoped>
.popover-Tip {
  height: 100px;
  min-width: 300px;
  max-width: 400px;
  overflow: auto;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
}

.mark-btn {
  font-size: 14px;
}

/deep/ {
  .select-input .select {
    width: 120px !important;
  }
}

/deep/ .el-tabs__active-bar {
  width: 60px !important;
}

/deep/ .el-range-separator {
  line-height: 26px;
}

/deep/ .el-select__tags {
  max-width: 170px !important;
}
</style>
