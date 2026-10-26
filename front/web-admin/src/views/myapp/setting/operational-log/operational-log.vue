<template>
  <div class="supernatant">
    <layout-list>
      <template slot="search">
        <search-item>
          <span class="m-r-10">账号查询:</span>
          <cy-input
            width-size="small"
            size="small"
            placeholder="请输入操作人"
            @change="search"
            v-model="name"
          ></cy-input>
        </search-item>
        <search-item>
          <span class="m-r-10">操作模块:</span>
          <cy-select
            v-model="model"
            placeholder="请选择操作模块"
            size="small"
            clearable
            :optionsList="model_list"
            :options="{ label: 'label', value: 'value' }"
          >
          </cy-select>
        </search-item>
        <search-item>
          <span class="m-r-10">选择时间:</span>
          <el-date-picker
            size="medium"
            v-model="time_vlaue"
            type="daterange"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :default-time="['00:00:00', '23:59:59']"
            @change="search"
            value-format="timestamp"
            @keyup.enter.native="search"
          >
          </el-date-picker>
        </search-item>
        <div>
          <el-button class="font14" size="small" @click="search" type="primary">
            查询
          </el-button>
          <el-button class="font14" size="small" @click="reset()">
            重置
          </el-button>
        </div>
      </template>
      <template slot="table">
        <cy-table :tableConfig="tableConfig"></cy-table>
      </template>
      <template slot="footer">
        <el-pagination
          v-if="list.length !== 0 && page_init"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="meta.page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="meta.size"
          layout="total, sizes, prev, pager, next, jumper"
          :total="meta.total"
        ></el-pagination>
      </template>
    </layout-list>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import UserModule from "@/store/modules/user";
import Api from "@/api/setting/bccount";
import Apitest from "@/api/home/home";
@Component({
  components: {
    LayoutList
  }
})
export default class extends Vue {
  //获取用户信息
  get USER() {
    return UserModule.getUserInfo;
  }
  private list = [];
  private tableColumns = [
    {
      key: "admin_name",
      title: "操作人"
    },
    {
      key: "module",
      title: "操作模块"
    },
    {
      key: "page",
      title: "操作页面"
    },
    {
      key: "content",
      title: "操作内容"
    },
    {
      key: "ip",
      title: "IP"
    },
    {
      key: "created_at",
      title: "操作时间",
      type: "date"
    }
  ];
  private empty = {
    text: "还没有操作日志列表呦~"
  };
  private tableOperate = {
    width: "0",
    key: [],
    more: [],
    operateKeyClick: (row: any, key: string) => {
      console.log("operateKey:", key, "-------", row);
    }
  };
  private tableConfig = {
    data: this.list,
    columns: this.tableColumns,
    empty: this.empty
    // operate: this.tableOperate,
  };
  private page_init = false; // 页面初始状态
  private meta = {
    total: 40, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private type_list = [
    {
      label: "停用",
      value: "1"
    },
    {
      label: "启用",
      value: "2"
    }
  ];
  private time_vlaue: any = "";
  private name = "";
  private model = "";
  private model_list: any = [];

  // 搜索
  private search() {
    if (this.time_vlaue === null) this.time_vlaue = "";
    this.meta.page = 1;
    this.fetchList();
    console.log(this.list);
  }
  // 重置
  private reset() {
    this.time_vlaue = "";
    this.meta.page = 1;
    this.model = "";
    this.name = "";
    this.fetchList();
  }

  // 每页
  private handleSizeChange(val: any) {
    this.meta.page = 1;
    this.meta.size = val;
    this.fetchList();
  }

  // 当前页
  private handleCurrentChange(val: any) {
    this.meta.page = val;
    this.fetchList();
  }

  // 获取查询列表
  private fetchList() {
    const obj = {
      search: JSON.stringify({
        admin_name: this.name,
        module: this.model,
        starttime: this.time_vlaue[0] / 1000,
        endtime: this.time_vlaue[1] / 1000
      }),
      page: this.meta.page,
      size: this.meta.size
    };
    Api.operationLog(obj).then(res => {
      if (res && res.data) {
        const { data } = res;
        this.list = data.list;
        this.tableConfig.data = data.list;
        this.meta = {
          page: data.page,
          size: data.size,
          total: data.total
        };
        this.meta.page = this.meta.page - 0;
        this.page_init = true;
      }
    });
  }
  //获取模块列表
  private modelList() {
    Api.operationModule({}).then(res => {
      if (res && res.data) {
        res.data.list.forEach((el: any, index: any) => {
          res.data.list[index] = {
            label: el,
            value: el
          };
        });
        this.model_list = [...res.data.list];
        console.log(this.model_list);
        // 加入并置顶全部
        const obj = {
          label: "全部",
          name: "全部"
        };
        this.model_list.unshift(obj);
        console.log(this.model_list);
      }
    });
  }

  mounted() {
    this.fetchList();
    this.modelList();
  }
}
</script>

<style lang="scss" scoped>
.supernatant {
  .font14 {
    font-size: 14px;
  }

  .table-empty-firstlabel {
    margin-top: 30px;
    img {
      width: 132px;
    }
  }

  // ::v-deep .operationalTable {
  //   .el-table__body-wrapper {
  //     td {
  //       height: 54px !important;
  //       padding: 0 !important;
  //     }
  //   }
  // }
}
.m-r-10 {
  margin-right: 10px;
}
</style>
