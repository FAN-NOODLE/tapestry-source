<template>
  <layout-list class="download-center">
    <!-- 搜索 -->
    <template slot="search">
      <search-item text="文件名称：">
        <cy-input
          v-model="searchParams.name"
          placeholder="请输入文件名称"
        ></cy-input>
      </search-item>
      <search-item text="创建时间：">
        <el-date-picker
          v-model="searchParams.time"
          size="medium"
          value-format="timestamp"
          type="datetimerange"
          range-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          :default-time="['00:00:00', '23:59:59']"
          :picker-options="pickerOptions"
        ></el-date-picker>
      </search-item>
      <search-item text="文件状态：">
        <cy-select
          placeholder="请选择文件状态"
          v-model="searchParams.status"
          size="small"
          :optionsList="statusList"
          :options="{ label: 'name', value: 'id' }"
        >
        </cy-select>
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

    <!-- 表格 -->
    <template slot="table">
      <el-table
        :data="tableList"
        stripe
        style="width: 100%"
        tooltip-effect="light"
      >
        <el-table-column
          label="文件名称"
          min-width="300"
          prop="name"
        >
        </el-table-column>
        <el-table-column
          label="文件状态"
          min-width="200"
        >
          <template slot-scope="scope">
            <el-tag v-if="parseInt(scope.row.status) === 1">已创建</el-tag>
            <el-tag
              type="warning"
              v-if="parseInt(scope.row.status) === 2"
            >生成中</el-tag>
            <el-tag
              type="success"
              v-if="parseInt(scope.row.status) === 3"
            >生成成功
            </el-tag>
            <el-tag
              type="danger"
              v-if="parseInt(scope.row.status) === 4"
            >生成失败</el-tag>
          </template>
        </el-table-column>
        <el-table-column
          label="创建人"
          min-width="250"
          prop="operation_name"
        >
        </el-table-column>
        <el-table-column
          label="创建时间"
          min-width="250"
        >
          <template slot-scope="scope"><span>{{ scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </span></template>
        </el-table-column>
        <el-table-column
          fixed="right"
          label="操作"
          width="160"
        >
          <template slot-scope="scope">
            <el-button
              v-if="scope.row.status === 3"
              @click="download(scope.row)"
              type="text"
              size="small"
              style="color: #6f3112;"
            >下载</el-button>
            <!-- <el-button
              class="operation-error"
              @click="deleted(scope.row)"
              type="text"
              size="small"
            >删除</el-button> -->
            <el-button
              v-if="scope.row.status === 4 || scope.row.status === 5"
              @click="reload(scope.row)"
              type="text"
              size="small"
            >重新生成</el-button>
          </template>
        </el-table-column>
        <div slot="empty">
          <layout-empty :message="'还没有下载记录哦~'"> </layout-empty>
        </div>
      </el-table>
    </template>
    <template slot="footer">
      <el-pagination
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
</template>
<script lang="ts">
import { ApiResult, Meta } from "@/types";
import { Component, Vue } from "vue-property-decorator";
import mediaObj from "./media";
import Api from "@/api/setting/download";
import UserModule from "@/store/modules/user";
import { downLoad } from "@/utils/down-load";
/**
 * 选项状态
 */
const STATUS = [
  { name: "全部", id: 0 },
  { name: "已创建", id: 1 },
  { name: "生成中", id: 2 },
  { name: "生成成功", id: 3 },
  { name: "生成失败", id: 4 },
];

// 搜索参数
type searchParamsType = {
  /** value值 */
  name: string;
  /** 搜索类型 */
  type: number;
  /** 页面状态 */
  status: number;
  /** 时间 */
  time: Array<number>;
};

@Component
export default class DownloadCenter extends Vue {
  /** 获取用户信息 */
  get USER() {
    return UserModule.getUserInfo;
  }
  private media_obj = mediaObj;

  /** 表格参数 */
  private tableList = [];

  /** 选项状态 */
  private statusList = STATUS;

  /** 分页配置 */
  private meta: Meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };

  /** 搜索条件 供双向绑定 */
  private searchParams: searchParamsType = this.initSearchParams();

  /** 搜索条件 供查询使用 */
  private newSearchParams: searchParamsType = this.initSearchParams();

  /** 初始化搜索条件 */
  public initSearchParams(): searchParamsType {
    return {
      name: "", // 搜索名称
      type: 1, // 搜索类型
      status: 0, // 文件状态
      time: [] // 创建时间
    };
  }
  private beforeDate: any = ""
  private pickerOptions = {
    disabledDate(time: any) {
      const nowDate = new Date();
      const beforeDate = new Date();//以当前时间作参考的前时间
      beforeDate.setDate(nowDate.getDate() - 181); //改变值 获取当前时间之前的不同的日期
      // beforeDate之前的时间都不可选
      return time.getTime() < beforeDate.getTime()
    }
  }
  /**
   * 检测字符串是否超长需截取的方法
   */
  checkLength(str: string, len: number) {
    if (str.split("").length < len) {
      return str;
    } else {
      return `${str.substr(0, len)}…`;
    }
  }

  /**
   * 内部绑定函数（获取每页数据量）
   */
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getTableList();
  }

  /**
   * 内部绑定函数（获取当前第几页）
   */
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getTableList();
  }

  /**
   * 搜索
   */
  private search() {
    this.meta.page = 1;
    console.log("点击了搜索", this.searchParams);
    this.newSearchParams = { ...this.searchParams };
    this.getTableList();
  }

  /**
   * 重置
   */
  private reset() {
    this.meta.page = 1;
    this.searchParams = this.initSearchParams();
    this.newSearchParams = this.initSearchParams();
    this.getTableList();
  }

  /**
   * 文件下载
   */
  async download(row: { id: number, name:string }) {
    const res: ApiResult = await Api.DownloadFile({ id: row.id, type: 1 });
    if (res.code !== 200) return this.$message.error(res.message as string);
    if (!res.data || !res.data.path) return this.$message.error("下载地址异常");
    downLoad(res.data.path, row.name);
  }

  /**
   * 删除文件
   */
  deleted(row: { id: number }) {
    this.$confirm("你确认要删除该文件吗?", "提示", {
      confirmButtonText: "确认删除",
      cancelButtonText: "我再想想",
      type: "warning"
    }).then(() => {
      Api.downloadDel({ file_id: row.id }).then((res: ApiResult) => {
        if (res.code !== 200) return this.$message.error(res.message as string);
        this.$message.success("删除文件成功");
        if (this.meta.page > 1 && this.tableList.length === 1) {
          this.meta.page--;
        }
        this.getTableList();
      });
    });
  }

  /**
   * 重新生成
   */
  reload(row: { id: number }) {
    Api.reDownload({ file_id: row.id }).then((res: ApiResult) => {
      if (res.code !== 200) return this.$message.error(res.message as string);
      this.$message.success("正在重新生成");
      this.getTableList();
    });
  }
  /**
   * 获取列表信息
   */
  async getTableList() {
    const params: any = {
      // enterprise_id: Number(this.USER.enterprise_id), // 企业id
      page: this.meta.page,
      page_size: this.meta.size,
      created_at: {
        min: Math.floor((this.newSearchParams.time[0] || this.beforeDate) / 1000),
        max: Math.floor((this.newSearchParams.time[1] || 0) / 1000)
      },
      status: this.searchParams.status // 文件状态
    };

    // 1-文件名称搜索 2-创建人搜索
    if (this.newSearchParams.type === 1) {
      params.filename = this.newSearchParams.name;
    }

    const res: ApiResult = await Api.DownloadList(params);
    if (res.code !== 200) return this.$message.error(res.message as string);

    this.tableList = res.data.list;
    this.meta.total = res.data.count;
  }

  /**
   * 生命周期
   */
  private created() {
    const nowDate = new Date();
    const beforeDate = new Date();//以当前时间作参考的前时间
    beforeDate.setDate(nowDate.getDate() - 181); //改变值 获取当前时间之前的不同的日期
    this.beforeDate = beforeDate.getTime()
    console.log(this.beforeDate)
    this.getTableList();
  }
}
</script>
<style lang="scss" scoped>
.download-center {
  .operation-error {
    color: red;
  }

  .filter {
    opacity: 0;
    width: 80px;
  }

  .file-message {
    cursor: pointer;

    .file-type {
      display: inline-block;
      width: 44px;
      height: 52px;
      vertical-align: middle;
    }

    .file-name {
      display: inline-block;
      vertical-align: middle;
      padding-left: 10px;
    }
  }
}
</style>
