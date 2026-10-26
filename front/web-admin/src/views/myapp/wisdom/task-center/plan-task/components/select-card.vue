<template>
  <div class="select-materials-dialog">
    <el-dialog
      :visible.sync="selectDialogVisible"
      width="800px"
      @close="cancel"
    >
      <div class="header">
        <div>
          <span>选择贺卡</span>
        </div>
        <span @click="cancelDialog">
          <i
            aria-hidden="true"
            class="el-icon-close"
          ></i>
        </span>
      </div>
      <!-- 默认展示 -->
      <div>
        <div class="search">
          <search-item text="搜索：">
            <cy-input
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入模板名称"
              clearable
              v-model="search.keyWord"
            ></cy-input>
          </search-item>
          <el-button
            type="primary"
            size="small"
            @click="touchSearch()"
          >查询</el-button>
          <el-button
            size="small"
            @click="reset()"
          >重置</el-button>
          <el-button
            size="small"
            type="primary"
            @click="addBtn()"
            style="margin-left: 130px;"
          >新增贺卡</el-button>
        </div>
        <div class="content">
          <el-table
            :data="list"
            stripe
            ref="multipleTable"
            highlight-current-row
            @row-click="rowClick"
            :empty-text="'暂无数据'"
            max-height="370"
          >
            <el-table-column width="35">
              <template slot-scope="scope">
                <el-radio
                  :label="scope.row.id"
                  v-model="selectTemplateId"
                ></el-radio>
              </template>
            </el-table-column>
            <el-table-column
              label="模板名称"
              prop="template_name"
            >
              <template slot-scope="scope">
                {{ scope.row.template_name || "--" }}
              </template>
            </el-table-column>
            <el-table-column
              label="模板ID"
              prop="id"
            >
              <template slot-scope="scope">
                {{ scope.row.id || "--" }}
              </template>
            </el-table-column>
            <el-table-column
              label="更新时间"
              prop="updated_time"
            >
              <template slot-scope="scope">
                {{ scope.row.updated_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
              </template>
            </el-table-column>
          </el-table>
        </div>
        <div class="footer-pages">
          <el-pagination
            :page-size="meta.size"
            :total="meta.total"
            layout="prev, pager, next"
            :current-page="meta.page"
            @current-change="handleCurrentChangeFn"
            class="page-right font-weight-400"
          ></el-pagination>
        </div>
      </div>
      <div
        slot="footer"
        class="dialog-footer"
      >
        <el-button
          size="small"
          @click="cancel"
        >取 消</el-button>
        <el-button
          size="small"
          type="primary"
          @click="submitPar"
        >确 认</el-button>
      </div>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import Api from "@/api/enterprise/material-library";
@Component
export default class extends Vue {
  @Prop() selectItem!: object;
  //------------------------------------------------------data-----------------------------------------------------------/
  private selectDialogVisible = true;
  // 搜索区域
  private search = {
    keyWord: "",
  };
  private list = [];
  private selectTemplateId = "";
  private selectRow = {};
  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 10,
  };
  // -------------新增视频data--------------------

  // ------------------生命周期-----------
  private mounted() {
    this.getlist();
  }
  //----------------------------------------------------methods-----------------------------------------------------------/
  // 列表数据
  private getlist() {
    Api.cardListApi({
      page: this.meta.page,
      size: this.meta.size,
      template_name: this.search.keyWord,
    }).then((res: any) => {
      const { List, total } = res.data;
      this.list = List;
      this.meta.total = total;
    });
  }
  private rowClick(row: any) {
    this.selectTemplateId = row.id;
    this.selectRow = row;
  }
  // 取消
  private cancel() {
    this.$emit("submit");
  }
  // 点击×号关闭弹框
  private cancelDialog() {
    this.cancel();
  }
  // 确认
  private submitPar() {
    this.$emit("submit", this.selectRow);
  }
  // 点击查询
  private touchSearch() {
    this.getlist();
  }
  // 点击重置
  private reset() {
    this.search.keyWord = "";
    this.meta.page = 1;
    this.getlist();
  }
  // 点击新增
  private addBtn() {
    // const newLink = this.$router.resolve({
    //   name: "myapp_customer_greeting-card_create",
    // });
    const newLink = `${origin}/web-admin/myapp/customer/greeting-card/create`;
    window.open(newLink, "_blank");
  }

  // 页码选择回调
  private handleCurrentChangeFn(val: number) {
    this.meta.page = val;
    this.getlist();
  }
}
</script>
<style lang="scss" scoped>
.select-materials-dialog {
  .header {
    display: flex;
    justify-content: space-between;
    text-align: center;
    margin-bottom: 20px;
    font-size: 18px;
    font-family: PingFangSC-Semibold, PingFang SC;
    cursor: pointer;

    .chosoe-clas {
      color: #33363c;
      font-weight: 600;
    }

    .add-clas {
      color: #33363c;
      font-weight: 600;
    }
  }

  /deep/.el-dialog__header {
    display: none;
  }

  .search {
    margin-bottom: 5px;
  }

  .footer-pages {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}</style>
