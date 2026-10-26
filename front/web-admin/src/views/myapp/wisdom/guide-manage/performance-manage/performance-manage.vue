<template>
  <div>
    <layout-list>
      <template slot="buttons">
        <el-button
          type="primary"
          @click="exportData"
        >批量导入数据</el-button>
      </template>
      <template slot="search">
        <search-item text="导购：">
          <el-input
            :placeholder="`请输入${type_list[gudie_type].name}`"
            class="search-input"
            size="medium"
            v-model.trim="search_key"
            clearable
          >
            <el-select
              slot="prepend"
              v-model="gudie_type"
              class="input-with-select"
            >
              <el-option
                v-for="item in type_list"
                :key="item.key"
                :value="item.key"
                :label="item.name"
              ></el-option>
            </el-select>
          </el-input>
        </search-item>
        <search-item text="所在部门：">
          <el-select
            style="width:280px"
            v-model="department.name"
            ref="selectBlur"
          >
            <el-option
              :value="department.value"
              :label="department.name"
              style="height: auto;padding: 0px;"
            >
              <el-tree
                ref="tree"
                class="treeStyle"
                v-if="department.isShowTree"
                :props="department.data"
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
        <search-item text="导购类型：">
          <cy-select
            v-model="search_active_type"
            size="medium"
            placeholder="请选择导购类型"
            clearable
          >
            <cy-option
              v-for="item in active_type_list"
              :value="item.id"
              :label="item.name"
              :key="item.id"
            >
            </cy-option>
          </cy-select>
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
          </div>
          <el-button
            class="mar-l-20"
            size="small"
            @click="exportList"
          >
            数据导出
          </el-button>
        </div>
      </template>
      <template slot="table">
        <el-table
          stripe
          :data="list"
          class="mar-t-20"
        >
          <el-table-column
            prop="counter_code"
            label="柜台code"
          ></el-table-column>
          <el-table-column
            prop="counter_name"
            label="柜台名称"
          ></el-table-column>
          <el-table-column
            prop="ba_name"
            label="导购姓名"
          ></el-table-column>
          <el-table-column
            prop="regional"
            label="区域"
          ></el-table-column>
          <el-table-column
            prop="monthly_sales"
            label="月度销售目标（元）"
          ></el-table-column>
          <el-table-column
            prop="monthly_second_buy_and_three_scan"
            label="新二购且一码三关注月度目标（人）"
          ></el-table-column>
          <el-table-column
            prop="monthly_old_customer_and_three_scan"
            label="老客购且一码三关注月度目标（人）"
          ></el-table-column>
          <el-table-column
            prop="monthly_new_customer_and_three_scan"
            label="新客购且一码三关注月度目标（人）"
          ></el-table-column>
          <div slot="empty">
            <layout-empty :message="'暂无数据'"></layout-empty>
          </div>
        </el-table>
      </template>
      <template slot="footer">
        <el-pagination
          v-if="list && list.length !== 0"
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
    <el-dialog
      title="导购业绩批量导入"
      :visible="isShowDialog"
      @close="closeDialog"
      :close-on-click-modal="false"
    >
      <div class="dialog-wraper">
        <el-row class="mt-20">
          <el-col :span="4">模板下载</el-col>
          <el-col :span="20">
            <p
              class="download"
              @click="downloadTemplate"
            >下载当前导入模版</p>
          </el-col>
        </el-row>
        <el-row class="mt-20">
          <el-col :span="4">模板上传</el-col>
          <el-col :span="20">
            <el-upload
              class="upload-demo"
              drag
              accept=".xls,.xlsx"
              name="excel"
              :headers="headers"
              :action="action"
              :multiple="false"
              :file-list="fileList"
              :before-upload="beforeUpload"
              :on-success="uploadSuccess"
              :on-error="uploadError"
            >
              <i class="el-icon-upload"></i>
              <div class="el-upload__text">
                将文件拖到此处，或<em>点击上传</em>
              </div>
              <div
                class="el-upload__tip"
                slot="tip"
              >只能上传xls、xlsx文件</div>
            </el-upload>
          </el-col>
        </el-row>
        <el-row class="mt-20">
          <el-col :span="4">上传结果</el-col>
          <el-col :span="20">
            <p>{{ uploadMessage }}</p>
          </el-col>
        </el-row>
      </div>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/enterprise/enterprise";
import GuideApi from "@/api/wisdom/Guide";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
import { VUE_APP_API_ROOT } from "@/utils/env-config";
import { dateFormat } from "../../../customer/auto-market/util/index";
import LoginAboat from "@/utils/LoginManager";
@Component({
  components: {
    LayoutEmpty,
  },
})
export default class extends Vue {
  // 数据
  private type_list: Array<any> = [
    { name: "导购ID", key: 0, keyword: "guide_no" },
    { name: "导购名称", key: 1, keyword: "guide_name" },
  ];
  private gudie_type = 1;
  private search_key = "";
  private department = {
    // 部门
    value: 0, // 部门id
    name: "",
    isShowTree: false,
    data: {
      label: "name",
      children: [],
      isLeaf: "leaf",
    },
  };
  // 部门列表
  private departmentList: Array<any> = [];
  private search_active_type = "";
  private active_type_list = [
    { name: "企业管理", id: 4 },
    { name: "企业导购", id: 3 },
    { name: "社群导购", id: 2 },
    { name: "门店导购", id: 1 },
  ];
  private updateTime = 0;
  private list: any = [];
  private meta = {
    page: 1,
    size: 10,
    total: 0,
  };
  private isShowDialog = false;
  private action = `${VUE_APP_API_ROOT}ba-performance/import`;
  private headers = {
    "Auth-Token": LoginAboat.getToken(),
  };
  private uploadMessage = "暂未检测到上传信息";
  private fileList: Array<any> = [];
  // 生命周期
  private mounted() {
    this.fetchStoreList();
    this.search();
  }
  private exportData() {
    this.isShowDialog = true;
  }
  private closeDialog() {
    this.isShowDialog = false;
    this.fileList = [];
    this.uploadMessage = "暂未检测到上传信息";
  }
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
      pid: node.data.third_id,
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
  // 获取部门
  private fetchStoreList() {
    Api.depatListApi({ pid: 0 }).then((res: any) => {
      if (res.code == 200) {
        const obj = { id: 0, name: "全部" };
        this.departmentList = res.data;
        this.departmentList.unshift(obj);
        this.department.isShowTree = true;
      } else {
        this.$message.error(res.message);
        return;
      }
    });
  }
  // 选择部门
  handleNodeClick(node: any) {
    this.department.name = node.name;
    this.department.value = node.id;
    (this.$refs.selectBlur as any).blur();
    this.search();
  }
  private search() {
    const params: any = {
      search: {
        dp_id: this.department.value,
        type: this.search_active_type,
      },
      page: this.meta.page,
      size: this.meta.size,
    };
    if (this.gudie_type === 1) {
      params.search.name = this.search_key;
    } else {
      params.search.worker_no = this.search_key;
    }
    GuideApi.baPerformanceList(params).then((res: any) => {
      if (res.code === 200) {
        const { list, page, size, total } = res.data;
        this.list = list;
        this.meta = { page, size, total };
        if (this.list && this.list.length > 0) {
          this.updateTime = this.list[0].created_at || 0;
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  private reset() {
    this.gudie_type = 1;
    this.search_key = "";
    this.department.value = 0;
    this.department.name = "";
    this.search_active_type = "";
    this.search();
  }
  // 导出
  private async exportList() {
    try {
      const data: any = {
        task_id: String(this.$route.query.val),
      };
      const res: any = await GuideApi.baPerformanceExport(data);
      const blob: Blob = new Blob([res], {
        type: "application/x-xls;charset=utf-8",
      });
      const a: HTMLElement = document.createElement("a");
      const href = window.URL.createObjectURL(blob);
      const name = `导购业绩管理${dateFormat(
        Date.now() / 1000,
        "yyyy_M_d"
      )}.xlsx`;
      a.setAttribute("download", name);
      a.setAttribute("href", href);
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(href);
    } catch (error) {
      console.warn("导出失败！");
    }
  }
  // 每页
  private handleSizeChange(val: any) {
    this.meta.page = 1;
    this.meta.size = val;
    this.search();
  }

  // 当前页
  private handleCurrentChange(val: any) {
    this.meta.page = val;
    this.search();
  }

  // 下载当前导入模板
  private async downloadTemplate() {
    try {
      const res: any = await GuideApi.baPerformanceDownload({});
      const blob: Blob = new Blob([res], {
        type: "application/x-xls;charset=utf-8",
      });
      const a: HTMLElement = document.createElement("a");
      const href = window.URL.createObjectURL(blob);
      const name = `导购业绩目标导入模板.xlsx`;
      a.setAttribute("download", name);
      a.setAttribute("href", href);
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(href);
    } catch (error) {
      console.warn("导出失败！");
    }
  }

  private beforeUpload(file: any) {
    console.log("file", file);
    const fileSuffix = file.name.substring(file.name.lastIndexOf(".") + 1);
    const whiteList = ["xls", "xlsx"];
    if (whiteList.indexOf(fileSuffix) === -1) {
      this.$message.error("上传文件只能是xls和xlsx");
      return false;
    }
  }
  private uploadSuccess() {
    this.$message.success("上传成功");
    setTimeout(() => {
      this.closeDialog();
      this.search();
    }, 1000);
  }

  private uploadError(res: any) {
    const error: any = JSON.parse(res.message).message;
    this.$message.error("上传失败,请重试");
    this.uploadMessage = error;
  }
}
</script>
<style lang="scss" scoped>
.search-input {
  width: 280px !important;
}
.input-with-select {
  width: 110px;
}
.data-button {
  display: inline-block;
  margin: 0 40px;
  color: #666666;
  font-size: 12px;
}
.mar-l-10 {
  margin-left: 10px;
}
.mt-20 {
  margin: 20px 0;
}
.download {
  color: #708df5;
  cursor: pointer;
}
</style>
