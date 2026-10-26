<template>
  <div class="service-setting">
    <layout-list>
      <template slot="search">
        <search-item text="服务名称：">
          <cy-input
            width-size="small"
            size="small"
            prefix-icon="el-icon-search"
            placeholder="请输入服务名称"
            v-model="searchValue"
          ></cy-input>
        </search-item>

        <search-item text="服务状态：">
          <el-select
            style="width:250px"
            class="select"
            v-model="search.reserve_status"
            size="medium"
          >
            <el-option
              v-for="item in serviceStatusList"
              :key="item.key"
              :label="item.label"
              :value="item.key"
            ></el-option>
          </el-select>
        </search-item>
        <search-item text="有效期： ">
          <el-date-picker
            type="daterange"
            v-model="search.service_time"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            range-separator="至"
            style="height:36px"
          ></el-date-picker>
        </search-item>

        <search-item>
          <el-button type="primary" @click="goSearch" size="small"
            >查询</el-button
          >
          <el-button plain @click="reset" size="small">重置</el-button>
        </search-item>
      </template>
      <template slot="operate">
        <el-button type="primary" @click="goAdd" size="small"
          >新增服务</el-button
        >
      </template>
      <template slot="table">
        <el-table tooltip-effect="light" :data="list">
          <el-table-column label="服务ID" prop="id" min-width="170">
          </el-table-column>
          <el-table-column label="服务名称" prop="reserve_name" min-width="150">
            <template slot-scope="scope">
              {{ scope.row.reserve_name }}
            </template>
          </el-table-column>
          <el-table-column label="时长" prop="reserve_time" min-width="100">
            <template slot-scope="scope">
              {{ scope.row.reserve_time }}小时
            </template>
          </el-table-column>
          <el-table-column label="图片" prop="reserve_img" min-width="155">
            <template slot-scope="scope">
              <el-image fit="cover"
                v-if="scope.row.reserve_img"
                style="width: 150px; height: 50px"
                :src="scope.row.reserve_img"
                :preview-src-list="[scope.row.reserve_img]"
              >
              </el-image>
              <span v-else>--</span>
            </template>
          </el-table-column>

          <el-table-column
            label="服务状态"
            prop="content"
            min-width="140"
            key="reserve_status"
          >
            <template slot-scope="scope">
              {{ serviceStatusMap[scope.row.reserve_status] || "--" }}
            </template>
          </el-table-column>
          <el-table-column
            label="有效期"
            prop="content"
            min-width="120"
            key="service_type"
          >
            <template slot-scope="scope">
              <div v-if="scope.row.service_type === 0">
                {{ scope.row.start_time | formatDate("yyyy-mm-dd") }} ~
                {{ scope.row.end_time | formatDate("yyyy-mm-dd") }}
              </div>
              <div v-else>长期有效</div>
            </template>
          </el-table-column>
          <el-table-column
            label="适用范围"
            prop="content"
            min-width="140"
            key="service_role"
          >
            <template slot-scope="scope">
              <div v-if="scope.row.service_role === 1">
                <span v-if="scope.row.role_rules.length < 3">{{
                  scope.row.role_rules.map(item => item.name).join(",")
                }}</span>
                <el-popover
                  placement="top-start"
                  v-else
                  title=""
                  width="100%"
                  trigger="hover"
                  :content="
                    scope.row.role_rules.map(item => item.name).join(',')
                  "
                >
                  <p slot="reference" class="text-span">
                    {{
                      scope.row.role_rules
                        ? scope.row.role_rules
                            .slice(0, 2)
                            .map(item => item.name)
                            .join(",")
                        : "--"
                    }}...
                  </p>
                </el-popover>
              </div>
              <div v-else>全部</div>
            </template>
          </el-table-column>

          <el-table-column label="操作" width="200" fixed="right">
            <template slot-scope="scope">
              <el-button type="text" @click="goLook(scope.row)" size="small"
                >查看</el-button
              >
              <el-button
                type="text"
                size="small"
                :disabled="scope.row.reserve_status == 1"
                @click="goEdit(scope.row)"
                >编辑</el-button
              >
              <el-button
                type="text"
                size="small"
                class="danger-txt"
                @click="delItem(scope.row)"
                >删除</el-button
              >
            </template>
          </el-table-column>
          <div slot="empty" v-if="!emptyFlag" style="margin-top: 50px;">
            <img class="empty-img" src="@/assets/images/empty-box.png" alt="" />
            <div class="empty-box">
              <div class="desc">暂无服务</div>
            </div>
          </div>
          <div slot="empty" v-if="emptyFlag" style="margin-top: 50px;">
            <img
              class="empty-img"
              src="@/assets/images/empty-search.png"
              alt=""
            />
            <div class="empty-box">
              <div class="desc">没有找到相关服务</div>
            </div>
          </div>
        </el-table>
      </template>
      <template slot="footer" v-if="list.length > 0">
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
    <!-- <ServiceCreate
      :title="title"
      :row="row"
      v-if="showCreate"
      @close="showCreate = false"
      @submit="submitData"
    ></ServiceCreate>
    <service-detail
      :title="title"
      :row="row"
      v-if="showDetail"
      @close="showDetail = false"
    ></service-detail> -->
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
// import ServiceCreate from "../../component/service-create.vue";
import Api from "@/api/store-invite/index";
// import ServiceDetail from "../../component/service-detail.vue";
@Component({
  components: {
    // ServiceCreate,
    // ServiceDetail
  }
})
export default class extends Vue {
  private searchValue = "";
  private meta = {
    page: 1,
    size: 10,
    total: 5
  };
  private title = "新增服务";
  private showCreate = false;
  private showDetail = false;
  private list = [];
  private multipleSelection = [];
  private emptyFlag = false;
  private row = {};
  private epNameDetail = "";
  private search = {
    reserve_status: -1,
    service_time: []
  };
  private serviceStatusList = [
    {
      key: -1,
      label: "全部"
    },
    {
      key: 0,
      label: "有效"
    },
    {
      key: 1,
      label: "失效"
    }
  ];
  private serviceStatusMap = {
    0: "有效",
    1: "失效"
  };
  // 生命周期
  private created() {
    console.log("created");
  }
  private mounted() {
    this.getList();
    this.getEpConfig();
  }
  private goLook(row) {
    this.row = row;
    this.title = "服务详情";
    this.showDetail = true;
    this.$router.push({
      name: "myapp_customer_store-invite_service-setting_detail",
      query: {
        id: row.id,
        title: this.title
      }
    });
  }
  private goEdit(row) {
    // TODO 服务状态为【有效】的服务可再次编辑，点击编辑后可编辑所有的配置项；【失效】的服务编辑置灰，仅可查看，不可编辑
    //1 失效  0有效
    if (row.reserve_status == 1) {
      return;
    }
    this.row = row;
    this.title = "编辑服务";
    this.showCreate = true;
    this.$router.push({
      name: "myapp_customer_store-invite_service-setting_edit",
      query: {
        id: row.id,
        title: this.title
      }
    });
  }
  private goAdd() {
    this.row = {};
    this.title = "新增服务";
    this.showCreate = true;
    this.$router.push({
      name: "myapp_customer_store-invite_service-setting_create",
      query: {
        id: "",
        title: this.title
      }
    });
  }

  // 获取企业信息
  private getEpConfig() {
    Api.getEpConfig().then(res => {
      // // 企业名称：coach;ks;sw
      const { name } = res.data;
      this.epNameDetail = name;
      // 缓存企业名称，给详情页面用
      localStorage.setItem("epNameDetail", name);
    });
  }
  // 删除服务
  private delService(id) {
    this.$confirm("是否需要删除该服务", "删除确认", {
      confirmButtonText: "确认删除",
      confirmButtonClass: "button-danger",
      cancelButtonText: "我再想想",
      type: "warning"
    })
      .then(() => {
        Api.deleteServiceApi({ id }).then(() => {
          this.$message.success("删除成功");
          this.getList();
        });
      })
      .catch(() => {
        this.$message({
          type: "info",
          message: "已取消删除"
        });
      });
  }
  // 删除前判断
  private delItem(row) {
    // TODO 被引用的服务不可删除，点击后弹窗提示“该服务已被其它模块功能引用使用，暂不可删除”

    Api.checkDelete({ id: row.id }).then(res => {
      if (res.data && res.data.status) {
        if (res.data.status === 2) {
          this.delService(row.id);
        } else if (res.data.status === 1) {
          this.$alert(
            "该服务已被其它模块功能引用使用，暂不可删除",
            "删除确认",
            {
              confirmButtonText: "关闭"
            }
          );
        }
      }
    });
  }
  private goSearch() {
    this.meta.page = 1;
    this.getList();
  }
  // 分页
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getList();
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getList();
  }
  // 当选择项发生变化时会触发该事件，数组
  selectChange(val: any) {
    this.multipleSelection = [];
    val.map((item: any) => {
      this.multipleSelection.push(item.id);
    });
  }
  private reset() {
    this.searchValue = "";
    this.search = {
      reserve_status: -1,
      service_time: []
    };
    this.meta.page = 1;
    this.getList();
  }
  private getList() {
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      reserve_name: this.searchValue,
      reserve_status: this.search.reserve_status,
      start_time: this.search.service_time
        ? (this.search.service_time[0] || 0) / 1000
        : 0,
      end_time: this.search.service_time
        ? (this.search.service_time[1] || 0) / 1000
        : 0
    };
    Api.serviceListApi(params).then(res => {
      const { total, list } = res.data;
      this.list =
        list && list.length
          ? list.map(item => ({ ...item, role_rules: item.role_rules || [] }))
          : [];
      this.meta.total = total;
    });
  }
  private submitData() {
    this.showCreate = false;
    this.getList();
  }
  private destroyed() {
    console.log("destroyed");
  }
}
</script>
<style lang="scss" scoped>
.danger-txt {
  color: rgba($color: red, $alpha: 0.8) !important;
}
</style>
