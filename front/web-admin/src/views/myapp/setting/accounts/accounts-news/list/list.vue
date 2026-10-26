<template>
  <layout-list title_name="账号管理">
    <!-- <template slot="buttons">
      <el-button
        size="small"
        type="primary"
        @click="toCreate"
      >
        新增账号
      </el-button>
    </template> -->
    <template slot="search">
      <search-item text="账号查询：">
        <select-input
          ref="selectInput"
          v-model="text"
          :list="acc_tel_list"
          @selectChange="changed"
          @change="changeText"
          :placeholder="search_place"
        ></select-input>
      </search-item>
      <search-item text="选择角色：">
        <cy-select
          placeholder="请选择角色"
          v-model="role_key"
          size="small"
          clearable="clearable"
          :optionsList="role_list"
        >
          <cy-option
            v-for="item in role_list"
            :key="item.key"
            :label="item.name"
            :value="item.key"
          ></cy-option>
        </cy-select>
      </search-item>
      <div>
        <el-button
          type="primary"
          size="small"
          @click="search()"
        >
          查询
        </el-button>
        <el-button
          size="small"
          @click="reset()"
        >
          重置
        </el-button>
      </div>
    </template>
    <template slot="table">
      <cy-table :tableConfig="tableConfig">
        <div
          slot="username"
          slot-scope="scope"
        >
          <div class="privew-box">
            <el-popover
              placement="bottom-start"
              trigger="hover"
              :disabled="scope.row.account.length <= 10"
            >
              <div class="privew-box">
                <el-scrollbar wrap-class="scrollbar-wrapper">
                  <div class="tag-box">
                    <span>{{ scope.row.account }}</span>
                  </div>
                </el-scrollbar>
              </div>
              <div
                class="tag-box"
                slot="reference"
              >
                <span size="small">{{ scope.row.account }}</span>
              </div>
            </el-popover>
          </div>
        </div>
        <div
          slot="admin_role_relation"
          slot-scope="scope"
        >
          <div class="privew-box">
            <el-popover
              placement="bottom-start"
              trigger="hover"
              :disabled="!scope.row.admin_role_relation ||
                scope.row.admin_role_relation.length < 2
                "
            >
              <div class="privew-box">
                <el-scrollbar wrap-class="scrollbar-wrapper">
                  <div class="tag-box">
                    <el-tag
                      class="m-tag"
                      size="small"
                      v-for="item in scope.row.admin_role_relation"
                      :key="item.id"
                    >
                      <div>{{ item.name }}</div>
                    </el-tag>
                  </div>
                </el-scrollbar>
              </div>
              <div
                class="tag-box"
                slot="reference"
                v-if="scope.row.admin_role_relation &&
                  scope.row.admin_role_relation.length
                  "
              >
                <div class="character-hiding-o">
                  <el-tag
                    class="m-tag"
                    size="small"
                    v-for="item in scope.row.admin_role_relation"
                    :key="item.id"
                  >{{ item.name }}
                  </el-tag>
                </div>
              </div>
              <div
                class="tag-box"
                slot="reference"
                v-else
              >-</div>
            </el-popover>
          </div>
        </div>
        <div
          slot="range"
          slot-scope="scope"
        >
          <div class="privew-box">
            <el-popover
              placement="bottom-start"
              trigger="click"
              :disabled="!scope.row.range || scope.row.range.length <= 10"
            >
              <div class="privew-box">
                <el-scrollbar wrap-class="scrollbar-wrapper">
                  <div class="tag-box">
                    <div
                      class="zhfw"
                      size="small"
                      v-for="item in scope.row.range"
                      :key="item.id"
                    >
                      {{ item.name }}
                    </div>
                  </div>
                </el-scrollbar>
              </div>
              <div
                class="tag-box"
                slot="reference"
                v-if="scope.row.range && scope.row.range.length"
              >
                <div class="character-hiding-o">
                  <span
                    class="mr"
                    size="small"
                    v-for="item in scope.row.range"
                    :key="item.id"
                  >{{ item.name }}
                  </span>
                </div>
              </div>
              <div
                class="tag-box"
                slot="reference"
                v-else
              >-</div>
            </el-popover>
          </div>
        </div>
      </cy-table>
    </template>
  </layout-list>
</template>
<script lang="ts">
import Api from "@/api/setting/Account";
// import MaxLen from "@/filters/MaxLen";
import { disOrOpenAccount, deleteAccount } from "./../utils";
// import debounce from "@/mixins/utilities/debounce";
import BaseData from "../base-data";
import { Vue, Component } from "vue-property-decorator";
import UserModule from "@/store/modules/user";
@Component({
  components: {},
  mixins: [BaseData]
})
export default class extends Vue {
  get USER() {
    return UserModule.getUserInfo;
  }
  private visible = true;
  private text = "";
  private search_place = "请输入账号";
  private tableConfig = {
    data: [],
    maxHeight: "",
    columnSort: (val: any) => {
      console.log("排序", val);
    },
    pagination: {
      pageSizes: [10, 20, 30, 40],
      pageSize: 10,
      total: 0,
      currentPage: 1,
      handleSizeChange: (val: any) => {
        this.handleSizeChange(val);
      },
      handleCurrentChange: (val: any) => {
        this.handleCurrentChange(val);
      }
    },
    empty: {
      text: "还没有创建账号哦，快去创建吧",
      emptyType: "list"
    },
    columns: [
      {
        key: "username",
        title: "账号名称",
        minWidth: "100",
        type: "slot"
      },
      {
        key: "mobile",
        title: "账号手机号",
        minWidth: "130"
        // type: "mobile"
      },
      {
        key: "admin_role_relation",
        title: "角色",
        minWidth: "160",
        type: "slot"
      },
      {
        key: "status",
        title: "账号状态",
        minWidth: "100",
        filter: [
          { value: 1, label: "启用" },
          { value: 2, label: "禁用" }
        ]
      },
      {
        key: "created_at",
        title: "创建时间",
        minWidth: "180",
        type: "date"
      },
      {
        key: "updated_at",
        title: "更新时间",
        minWidth: "180",
        type: "date"
      }
    ],
    operate: {
      title: "操作",
      width: 250,
      fixed: "right",
      key: [
        {
          key: "detail",
          text: "查看"
        },
        {
          key: "addNote",
          text: "编辑",
          style: "color: #6F3112"
        },
        {
          key: "isDisable",
          enumKey: "status",
          enum: { 1: "禁用", 2: "启用" },
          style: "color: #F56C6C"
        },
        {
          key: "del",
          text: "删除",
          style: "color: #F56C6C"
        }
      ],
      operateKeyClick: (row: any, key: any) => {
        switch (key) {
          case "detail":
            this.goDetail(row);
            break;
          case "addNote":
            this.goAddOrEdit(row);
            break;
          case "isDisable":
            this.handleOpenOrClose(row);
            break;
          case "del":
            //del
            this.delSpike(row);
            break;
          case "reset":
            //重置密码
            this.resetPassword(row);
            break;
        }
      }
    }
  };
  private DEBOUNCE = null;
  private searching = false;
  //账号-手机号列表配置
  private acc_tel_set = {
    search_value: "", //搜索内容
    search_class: "text", //搜索类型
    search_key: "username", //搜索key
    search_place: "请输入账号" //搜索placeholder
  };
  // 账号-手机号列表
  private acc_tel_list = [
    {
      key: "username",
      type: "text",
      name: "账号名称",
      place: "请输入账号"
    },
    {
      key: "mobile",
      type: "number",
      name: "账号手机号",
      place: "请输入账号手机号"
    }
  ];
  // 账号范围列表
  private range_list = [
    {
      key: "0",
      name: "全部"
    }
  ];
  // 角色列表
  private role_list = [
    {
      key: "0",
      name: "全部"
    }
  ];
  private range_key = ""; //账号范围
  private role_key = ""; //选择角色
  private page_init = false; // 页面初始状态
  private search_key = ""; // 搜索活动名称
  private meta = {
    total: 40, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private list = [];
  private changed(val: any) {
    if (val == "mobile") {
      this.search_place = "请输入账号手机号";
    } else if (val == "username") {
      this.search_place = "请输入账号";
    }
    this.text = "";
    this.acc_tel_set.search_value = "";
    this.acc_tel_set.search_key = val;
  }
  private changeText(val: any) {
    this.acc_tel_set.search_value = val;
  }
  // 账号，手机号切换
  private querysetChange(e: any) {
    // 同步修改placeHoder
    this.acc_tel_set.search_value = "";
    this.acc_tel_list.forEach(item => {
      if (item.key == e) {
        //修改相关搜索类型和搜索placeholder
        this.acc_tel_set.search_class = item.type;
        this.acc_tel_set.search_place = item.place;
      }
    });
    this.search();
  }
  // 返回上一页
  private goBack() {
    this.$router.go(-1);
  }
  // 搜索
  private search() {
    this.tableConfig.pagination.currentPage = 1;
    this.fetchList("search");
  }
  // 重置
  private reset() {
    this.text = "";
    (this.$refs["selectInput"] as any).reset();
    this.changed("username");
    this.acc_tel_list = [];
    this.acc_tel_list = [
      {
        key: "username",
        type: "text",
        name: "账号名称",
        place: "请输入账号"
      },
      {
        key: "mobile",
        type: "number",
        name: "账号手机号",
        place: "请输入账号手机号"
      }
    ];
    this.search_key = "";
    this.acc_tel_set = {
      search_value: "", //搜索内容
      search_class: "text", //搜索类型
      search_key: "username", //搜索key
      search_place: "请输入账号" //搜索placeholder
    };
    this.range_key = "0"; //账号范围
    this.role_key = "0"; //选择角色
    this.page_init = false; // 页面初始状态
    this.search_key = ""; // 搜索活动名称
    // this.meta.page = 1;
    this.fetchList();
  }

  // 每页
  private handleSizeChange(val: any) {
    this.tableConfig.pagination.currentPage = 1;
    this.tableConfig.pagination.pageSize = val;
    this.fetchList();
  }

  // 当前页
  private handleCurrentChange(val: any) {
    this.tableConfig.pagination.currentPage = val;
    this.fetchList();
  }

  // 获取账号列表
  private fetchList(key?: string) {
    const obj = {
      search_type:
        this.acc_tel_set.search_value && this.acc_tel_set.search_value != ""
          ? this.acc_tel_set.search_key
          : "all",
      keyword: this.acc_tel_set.search_value,
      object_id: this.range_key,
      role_id: this.role_key,
      page: this.tableConfig.pagination.currentPage,
      size: this.tableConfig.pagination.pageSize
    };
    this.searching =
      this.role_key != "0" ||
      this.range_key != "0" ||
      this.acc_tel_set.search_value != "";
    Api.accountsList(obj).then(res => {
      const data = res.data.list.map((item: any) => {
        const index = item.admin_role_relation.findIndex((v: any) =>
          this["no_choice"].includes(v.mark)
        );
        if (index > -1) item.disable_edit = true;
        return {
          ...item,
          detail: true, //允许查看
          // addNote: !item.disable_edit, //允许编辑
          addNote: true,
          isDisable: !item.disable_edit,
          del: !item.disable_edit && item.status !== 1, //允许删除
          delShow: !item.disable_edit && item.status !== 1, //允许删除
          reset: this.USER.is_super === 1
        };
      });
      this.tableConfig.data = data;
      this.tableConfig.pagination = {
        ...this.tableConfig.pagination,
        currentPage: res.data.page,
        pageSize: res.data.size,
        total: res.data.total
      };
      this.page_init = true;
      if (key === "search" && data.length === 0) {
        this.tableConfig.empty = {
          text: "没有找到该账号",
          emptyType: "search"
        };
      }
    });
  }
  // todo:重置密码
  private resetPassword(val: any) {
    console.log(val);
    this.$confirm("确定重置密码吗？", "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    })
      .then(() => {
        Api.resetPassword({
          admin_id: val.id
        }).then(() => {
          this.$message({
            type: "success",
            message: "重置密码成功"
          });
        });
      })
      .catch(() => {
        this.$message({
          type: "info",
          message: "已取消重置密码"
        });
      });
  }
  // 新建账号
  private toCreate() {
    Api.checkCount().then(() => {
      this.$router.push({
        name: "myapp_setting_accounts_accounts-news_accounts-edit",
        query: {
          bread: "新增账号"
        }
      });
    });
  }
  private goDetail(val: any) {
    this.$router.push({
      name: "myapp_setting_accounts_accounts-news_accounts-detail",
      query: {
        id: val.id
      }
    });
  }
  private getRole() {
    Api.roleList({
      search_text: "",
      size: 100
    }).then(res => {
      res.data.list.forEach((item: any) => {
        this.role_list.push({
          key: item.id,
          name: item.name
        });
      });
      this.role_list.forEach(item => {
        if (item.name === this.$route.query.role) {
          this.role_key = item.key;
          this.fetchList();
        }
      });
    });
    // GoodsApi.getStroeList({ size: 1000 }).then((res: any) => {
    //   res.data.list.forEach((item: any) => {
    //     const obj = {
    //       key: item.id.toString(),
    //       name: item.name.toString(),
    //     };
    //     this.range_list.push(obj);
    //     console.log(this.range_list);
    //   });
    // });
  }
  // 编辑角色
  private goAddOrEdit(val: any) {
    this.$router.push({
      name: "myapp_setting_accounts_accounts-news_accounts-edit",
      query: {
        id: val.id,
        is_edit: "1",
        bread: "编辑账号"
      }
    });
  }
  // 停用/启用
  private handleOpenOrClose(params: any) {
    disOrOpenAccount(params.id, params.status == 1 ? false : true).then(() => {
      this.fetchList();
    });
  }
  // 删除
  private delSpike(params: any) {
    deleteAccount(params.id).then(() => {
      if (this.meta.page == 1) {
        this.meta.page = 1;
      } else {
        this.meta.page -= 1;
      }
      this.fetchList();
    });
  }
  mounted() {
    if (
      this["$global"].enter_path &&
      this["$global"].enter_path === this.$route.path
    ) {
      if (this["$global"].meta) this.meta = this["$global"].meta;
      this.search_key = this["$global"].search_key;
    }
    this.fetchList();
    this.getRole();
    // this.DEBOUNCE = this.debounce(this.search, 1000, false);
    // todo:如果是超管账号，添加重置密码操作
    if(this.USER.is_super === 1){
      this.tableConfig.operate.key.push({
         // todo:重置密码
          key: "reset",
          text: "重置密码",
          style: "color: #6F3112"
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.account-confirm-content {
  padding-left: 30px !important;
  text-align: left !important;
  margin-bottom: 10px;
}

/deep/ {
  .el-select.search-select {
    width: 200px !important;
  }

  .search-input {
    width: 200px !important;
  }

  .scrollbar-wrapper {
    overflow-x: hidden !important;
    min-height: 50px;

    .el-tag {
      margin-bottom: 10px;
    }
  }

  .main-content {
    padding: 30px;
  }

  .page-right.font-weight-400.el-pagination {
    padding: 20px 5px 2px;
  }

  .main-search {
    margin-bottom: -5px;
  }

  .el-tag--small.el-tag--light {
    background: #f0f3fe !important;
  }
}

.character-hiding-o {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;

  .m-tag {
    margin-right: 10px;
  }
}

.tag-box {
  display: flex;
  flex-wrap: wrap;
  max-width: 150px;

  .el-tag {
    margin-right: 10px;
  }
}

.mr {
  margin-right: 10px;
}

.zhfw {
  font-size: 14px;
  color: #606266;
  width: 100%;
  margin-bottom: 10px;
}

.account.reset {
  margin-left: 20px !important;
  margin-bottom: 24px !important;
}

.list-del {
  color: red !important;
}

.mar-t-15 {
  margin-top: 0px !important;
}

/deep/ {
  .select-input .select {
    width: 120px !important;
  }
}
</style>
