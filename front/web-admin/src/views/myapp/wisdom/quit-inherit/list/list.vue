<template>
  <div class="quit-inherit">
    <layout-list>
      <template slot="buttons">
        <el-tabs
          v-model="activeName"
          type="card"
          @tab-click="tabClick"
        >
          <el-tab-pane
            label="待分配导购列表"
            name="first"
          > </el-tab-pane>
          <el-tab-pane
            label="已分配导购列表"
            name="second"
          > </el-tab-pane>
        </el-tabs>
      </template>

      <template slot="search">
        <div>
          <search-item text="导购：">
            <el-input
              placeholder="请输入内容"
              v-model="keyword"
              clearable
            >
              <el-select
                v-model="search_key"
                style="width: 105px;"
                slot="prepend"
                placeholder="请选择"
              >
                <el-option
                  label="导购姓名"
                  value="1"
                ></el-option>
                <el-option
                  label="导购工号"
                  value="2"
                ></el-option>
              </el-select>
            </el-input>
          </search-item>
          <search-item text="所在部门：">
            <el-select
              class="search-select"
              :value="select_dp.name"
              placeholder="请选择所属部门"
              ref="selectBlur"
              size="medium"
              clearable
              @change="changeSearch"
            >
              <el-option
                :value="select_dp.name"
                :label="select_dp.name"
              >
                <!-- <el-tree
                  class="treeStyle"
                  :props="dp_props"
                  node-key="dp_id"
                  :expand-on-click-node="true"
                  :check-strictly="true"
                  :check-on-click-node="true"
                  @node-click="handleNodeClick"
                  ref="tree"
                  @current-change="blurFn"
                ></el-tree> -->
                <el-tree
                  class="treeStyle"
                  @current-change="blurFn"
                  :data="dp_data"
                  :props="dp_props"
                  :expand-on-click-node="true"
                  :check-strictly="true"
                  :check-on-click-node="true"
                  @node-click="handleNodeClick"
                ></el-tree>
              </el-option>
            </el-select>
          </search-item>
          <div>
            <el-button
              class=""
              size="small"
              type="primary"
              @click="search"
            >
              查询
            </el-button>
            <el-button
              class="mar-l-20"
              size="small"
              @click="reset"
            >
              重置
            </el-button>
            <div class="data-button">
              <span>数据更新时间：</span>
              <span>{{ updateTime | formatDate("yyyy-mm-dd hh:mm:ss") }}</span>
              <el-button
                class="m-l-10"
                type="text"
                @click="synchroData(true)"
                size="medium"
              >更新</el-button>
            </div>
          </div>
        </div>
      </template>
      <template slot="table">
        <div v-if="activeName === 'first'">
          <listNotDistributionVue
            @submit="search"
            :guideList="guide_list"
          ></listNotDistributionVue>
        </div>
        <div v-else>
          <listAlreadyDistributionVue :guideList="guide_list"></listAlreadyDistributionVue>
        </div>
      </template>
      <template slot="footer">
        <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="meta.page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="meta.size"
          class="page-right font-weight-400"
          layout="total, sizes, prev, pager, next, jumper"
          :total="meta.total"
        ></el-pagination>
      </template>
    </layout-list>
  </div>
</template>
<script lang="ts">
import personApi from "@/api/wisdom/Personnel";
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/wisdom/quit-inherit";
import listAlreadyDistributionVue from "../components/list-already-distribution.vue";
import listNotDistributionVue from "../components/list-not-distribution.vue";
@Component({
  components: {
    listNotDistributionVue,
    listAlreadyDistributionVue,
  },
})
export default class extends Vue {
  // data数据
  private dp_props = {
    label: "dp_name",
    children: "children_list",
    isLeaf: "leaf",
  };
  private dp_data = [];
  private activeName = "first";
  private keyword = "";
  private search_key = "1";
  private select_dp = { name: "", value: 0 };
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10, // 每页数量
  };
  private updateTime = 0;
  private updateLock = false; // 更新锁
  private guide_list = [];
  private changeSearch() {
    this.select_dp.value = 0;
    this.select_dp.name = "";
    this.search();
  }
  // 选择部门
  private handleNodeClick(node: any) {
    this.select_dp.name = node.dp_name;
    this.select_dp.value = node.department_id;
    this.search();
  }
  // 关闭筛选
  private blurFn() {
    (this.$refs.selectBlur as any).blur();
  }
  // 重置
  private reset() {
    this.keyword = "";
    this.select_dp.value = 0;
    this.select_dp.name = "";
    this.search();
  }
  // 搜索
  private search() {
    this.meta.page = 1;
    this.inheritList();
  }

  // 手动点击同步数据
  private synchroData(exeSync?: boolean) {
    // const nowTime = new Date().getTime() / 1000;
    // 5分钟内不让再次更新
    // console.log(nowTime - this.updateTime);
    // this.updateLock = nowTime - this.updateTime < 300;
    // if (exeSync) {
    //   if (this.updateLock) {
    //     this.$message("努力同步中，请耐心等待");
    //     return;
    //   }
    // }
    let pramas = {};
    if (exeSync) {
      pramas = {
        exe_sync: exeSync,
      };
    }
    Api.syncCustomerData(pramas).then(
      (res: any) => {
        if (res.code === 200) {
          this.updateTime = parseInt(res.data);
          if (exeSync) {
            this.$message.success("操作成功，数据更新中");
            this.inheritList(exeSync);
          }
        } else {
          this.$message.error(res.message);
        }
      },
      (err) => {
        this.updateLock = false;
      }
    );
  }
  // 每页
  private handleSizeChange(val: any) {
    this.meta.page = 1;
    this.meta.size = val;
    this.inheritList();
  }

  // 当前页
  private handleCurrentChange(val: any) {
    this.meta.page = val;
    this.inheritList();
  }

  // 获取列表
  private inheritList(exeSync?: any) {
    const pramas = {
      page: this.meta.page,
      size: this.meta.size,
      status: this.activeName === "first" ? 0 : 1,
      handover_department_id: this.select_dp.value,
      handover_name: this.search_key === "1" ? this.keyword : "",
      handover_userid: this.search_key === "2" ? this.keyword : "",
    };
    Api.inheritList(pramas).then((res: any) => {
      if (res.code === 200) {
        this.guide_list = res.data.list;
        this.meta.total = res.data.total;
        if (exeSync) {
          this.$message.success("数据更新成功");
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // tab切换
  tabClick(tab: any) {
    this.guide_list = [];
    this.search();
  }
  // 获取部门
  departmentTree() {
    Api.departmentTree().then((res: any) => {
      if (res.code === 200) {
        this.$nextTick(() => {
          this.dp_data = res.data.list;
        });
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 生命周期
  private mounted() {
    this.updateTime = Date.now();
    this.inheritList();
    this.synchroData();
    this.departmentTree();
  }
}
</script>
<style lang="scss" scoped>
/deep/ .layout-list .buttons-content {
  margin-bottom: 0;
}
/deep/ .el-input-group {
  width: calc(100% - 43px);
}
/deep/ .el-tabs--card .el-tabs__nav {
  .is-active {
    color: #fff;
    background-color: #000000;
  }
}
.search-select {
  width: 280px;
}
.el-select-dropdown__item.selected {
  font-weight: normal !important;
  height: auto !important;
}
.data-button {
  display: inline-block;
  margin-left: 30px;
  color: #666666;
  font-size: 14px;
}
</style>
