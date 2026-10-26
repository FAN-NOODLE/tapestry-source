<template>
  <layout-list>
    <template slot="buttons">
      <el-button
        size="small"
        type="primary"
        @click="toCreate"
      >
        新增角色
      </el-button>
    </template>
    <template slot="search">
      <search-item text="名称查询：">
        <select-input
          v-model="text"
          :list="acc_tel_list"
          @selectChange="changed"
          @change="changeText"
          :placeholder="search_place"
        ></select-input>
      </search-item>
      <search-item>
        <el-button
          size="small"
          @click="search()"
          type="primary"
        >
          查询
        </el-button>
        <el-button
          size="small"
          @click="reset()"
        >重置</el-button>
      </search-item>
    </template>
    <template slot="table">
      <cy-table :tableConfig="tableConfig"></cy-table>
    </template>
  </layout-list>
</template>
<script lang="ts">
import Api from "@/api/setting/Role";
import { Vue, Component } from "vue-property-decorator";
import BaseData from "../base-data";
import UserModule from "@/store/modules/user";
import { roleAppName } from "../app-options";
@Component({
  components: {},
  mixins: [BaseData],
})
export default class extends Vue {
  get USER() {
    return UserModule.getUserInfo;
  }
  // 角色-创建人列表
  private acc_tel_list = [
    {
      key: 1,
      name: "角色名称",
    },
    {
      key: 2,
      name: "创建人",
    },
  ];
  private tableConfig = {
    data: [],
    maxHeight: "",
    columnSort: (val: any) => {
      this.sortChange(val);
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
      },
    },
    empty: {
      text: "还没有创建角色哦，快去创建吧",
      emptyType: "list",
    },
    columns: [
      {
        key: "name",
        title: "角色名称",
        minWidth: "180",
      },
      {
        key: "role_range",
        title: "适用后台",
        minWidth: "160",
      },
      {
        key: "admin_relation_count",
        title: "关联员工数量",
        sortable: "custom",
        minWidth: "150",
        type: "int",
      },
      {
        key: "admin_name",
        title: "创建人",
        minWidth: "240",
      },
      {
        key: "created_at",
        title: "创建时间",
        minWidth: "180",
        type: "date",
      },
      {
        key: "updated_at",
        title: "更新时间",
        minWidth: "180",
        type: "date",
      },
    ],
    operate: {
      title: "操作",
      width: 172,
      fixed: "right",
      key: [
        {
          key: "detail",
          text: "查看",
          style: "color: #303133",
        },
        {
          key: "addNote",
          text: "编辑",
          style: "color: #6F3112",
        },
        {
          key: "del",
          text: "删除",
          style: "color: #F56C6C",
        },
      ],
      operateKeyClick: (row: any, key: any) => {
        switch (key) {
          case "detail":
            this.goDetail(row);
            break;
          case "addNote":
            this.goEdit(row);
            break;
          case "del":
            this.delRole(row);
            break;
        }
      },
    },
  };
  private page_init = false; // 页面初始状态
  private search_key = ""; // 搜索角色名称
  private searching = false; //正在搜索中
  private sort_type = "";
  private type_list = [
    { name: "角色名称", id: 1 },
    { name: "创建人", id: 2 },
  ];
  private type = 1;
  private meta = {
    total: 40, // 总数
    page: 1, // 当前页
    size: 10, // 每页数量
  };
  private options = [];
  private list = [];
  private text = "";
  private search_place = "请输入角色名称";
  private changed(val: any) {
    if (val === 1) {
      this.search_place = "请输入角色名称";
    } else if (val === 2) {
      this.search_place = "请输入创建人";
    }
    this.text = "";
    this.search_key = "";
    this.type = val;
  }
  private changeText(val: any) {
    this.search_key = val;
    this.tableConfig.pagination.currentPage = 1;
  }
  private changeInput() {
    this.search_key = "";
    this.fetchList();
  }
  // 新建角色
  private toCreate() {
    this.$router.push({
      name: "myapp_setting_accounts_roles_role-edit",
      query: {
        type: "1",
        bread: "新增角色",
      },
    });
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
    this.sort_type = "";
    this.search_key = "";
    this.tableConfig.pagination.currentPage = 1;
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
  private sortChange(res: any) {
    switch (res.order) {
      case "descending":
        this.sort_type = "DESC";
        break;
      case "ascending":
        this.sort_type = "ASC";
        break;
      default:
        this.sort_type = "";
        break;
    }
    this.fetchList();
  }
  // 获取角色列表
  private fetchList(key?: string) {
    const params = {
      page: this.tableConfig.pagination.currentPage,
      size: this.tableConfig.pagination.pageSize,
      sort_type: this.sort_type,
      role_name: "",
      admin_name: "",
    };
    this.searching = this.search_key != "";
    this.type === 1
      ? (params.role_name = this.search_key)
      : (params.admin_name = this.search_key);
    Api.roleList(params).then((res) => {
      const data = res.data.list.map((item: any) => {
        const isRoles = !(this as any).default_roles.includes(item.mark); // 内置的角色,不可删除编辑
        const isConsoleRoles = !(this as any).console_default_roles.includes(
          item.mark
        ); // 管理后台角色,可编辑不可删除
        return {
          ...item,
          detail: true, //允许查看
          role_range: roleAppName(item.app_code),
          addNote: isRoles, //允许编辑
          addNoteShow: isRoles, // 显示编辑按钮
          del: isRoles && isConsoleRoles, //允许删除
          delShow: isRoles && isConsoleRoles, //显示删除按钮
        };
      });
      this.tableConfig.data = data;
      this.tableConfig.pagination = {
        ...this.tableConfig.pagination,
        currentPage: res.data.page,
        pageSize: res.data.size,
        total: res.data.total,
      };
      this.meta.page = this.meta.page - 0;
      this.page_init = true;
      if (key === "search" && data.length === 0) {
        this.tableConfig.empty = {
          text: "没有找到该角色",
          emptyType: "search",
        };
      }
    });
  }
  // 编辑
  private goEdit(val: any) {
    (this as any).$global.enter_path = this.$route.path;
    (this as any).$global.meta = this.meta;
    (this as any).$global.search_key = this.search_key;
    this.$router.push({
      name: "myapp_setting_accounts_roles_role-edit",
      query: {
        id: val.id,
        type: "2",
        bread: "编辑角色",
      },
    });
  }
  // 查看详情
  private goDetail(val: any) {
    (this as any).$global.enter_path = this.$route.path;
    (this as any).$global.meta = this.meta;
    (this as any).$global.search_key = this.search_key;
    this.$router.push({
      name: "myapp_setting_accounts_roles_role-detail",
      query: {
        id: val.id,
        identity_type: val.identity_type,
      },
    });
  }
  // 删除
  private delRole(params: any) {
    const self = this;
    self
      .$confirm("确认要删除该角色吗? 删除后将无法恢复！", "删除角色", {
        confirmButtonText: "确认删除",
        cancelButtonText: "我再想想",
        type: "warning",
      })
      .then(() => {
        const obj = {
          role_id: params.id,
        };
        Api.deleteRole(obj).then(() => {
          self.$message({
            type: "success",
            message: "删除成功",
          });
          if (this.meta.page > 1 && this.list.length === 1) {
            this.meta.page--;
          }
          self.fetchList();
        });
      })
      .catch();
  }

  mounted() {
    if (
      (this as any).$global.enter_path &&
      (this as any).$global.enter_path === this.$route.path
    ) {
      if ((this as any).$global.meta) this.meta = (this as any).$global.meta;
      this.search_key = (this as any).$global.search_key;
    }
    this.fetchList();
    // this.DEBOUNCE = this.debounce(this.search, 1000, false);
  }
}
</script>

<style lang="scss" scoped>
.bg-bt__color {
  color: $color-danger;
}

.el-button {
  + {
    .el-button {
      margin-left: 20px !important;
    }
  }
}

/deep/ {
  .main-content {
    padding: 30px !important;
  }

  .search-content {
    padding-bottom: 0px !important;
  }
}</style>
