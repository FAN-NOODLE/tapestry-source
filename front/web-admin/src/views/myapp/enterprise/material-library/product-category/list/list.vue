<template>
  <div class="welcome_container">
    <div class="container_content">
      <div class="search_content">
        <div class="line">
          <search-item text="标签搜索：">
            <cy-input
              width-size="small"
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入标签名称"
              clearable
              v-model="searchGroup.search_text"
            ></cy-input>
          </search-item>
          <el-button
            type="primary"
            @click="searchFn"
            size="small"
          >查询</el-button>
          <el-button
            @click="resetFn"
            size="small"
          >重置</el-button>
        </div>
      </div>
      <div class="button-group">
        <el-button
          type="primary"
          @click="createFn"
          size="small"
        >
          新增素材标签
        </el-button>
      </div>
      <div class="section">
        <el-table
          :data="categoryTable.table"
          stripe
          tooltip-effect="light"
        >
          <el-table-column
            label="标签名称"
            prop="name"
            show-overflow-tooltip
          >
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="scope">
              <el-button
                @click="categorySourceRelation(scope.row.id)"
                type="text"
                style="color:#303133;"
                size="small"
              >删除</el-button>
              <el-button
                type="text"
                size="small"
                @click="editRowFn(scope.row)"
              >编辑</el-button>
            </template>
          </el-table-column>
          <div
            slot="empty"
            class="table-empty-label"
          >
            <img
              :src="searchLoading
                  ? 'https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/empty-search.png'
                  : require('@/assets/images/empty-box.png')
                "
              alt=""
            />
            <div>{{ searchLoading ? "没有找到该标签" : "暂无标签" }}</div>
          </div>
        </el-table>
        <div class="table_footer">
          <el-pagination
            v-if="categoryTable.total > 0"
            @size-change="handleSizeChangeFn"
            @current-change="handlePageChangeFn"
            :current-page="categoryTable.page"
            :page-sizes="[10, 20, 30, 40]"
            :page-size="categoryTable.size"
            layout="total, sizes, prev, pager, next, jumper"
            :total="categoryTable.total"
          >
          </el-pagination>
        </div>
      </div>
    </div>
    <add-category
      v-if="addVisible"
      :selectedCate="selectedCate"
      :dialogType="dialogType"
      @closeVisible="closeVisible"
      @submit="addsubmit"
    ></add-category>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import addCategory from "../../components/add-category.vue";
import API from "@/api/enterprise/material-library";
@Component({
  components: {
    addCategory
  }
})
export default class extends Vue {
  private categoryTable = {
    table: [],
    size: 10,
    page: 1,
    total: 0
  };
  private selectedCate = null;
  private addVisible = false; // 控制添加素材标签弹窗
  private dialogType = "1";
  private searchGroup = {
    search_text: ""
  };
  private searchLoading = false;
  private handleSizeChangeFn(val: number) {
    this.categoryTable.page = 1;
    this.categoryTable.size = val;
    this.getTableDataFn();
  }
  private handlePageChangeFn(val: number) {
    this.categoryTable.page = val;
    this.getTableDataFn();
  }
  // 获取管理素材ID集合
  private categorySourceRelation(id: any) {
    API.categorySourceRelation({
      id: id
    }).then((res: any) => {
      if (res.code === 200) {
        const title = res.data.length
          ? "该素材标签已被素材关联，确定删除该标签吗?"
          : "确定删除该标签吗?";
        this.delItem(id, title);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //删除
  private delItem(id: any, title: string) {
    this.$confirm(title, "提示", {
      confirmButtonText: "确认",
      cancelButtonText: "取消",
      type: "warning"
    }).then(async () => {
      const res: any = await API.categoryDelete({ id: id });
      if (res.code == 200) {
        this.$message({
          type: "success",
          message: "删除成功!"
        });
        this.getTableDataFn();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 获取数据
  async getTableDataFn() {
    const param: object = {
      page: this.categoryTable.page,
      size: this.categoryTable.size,
      name: this.searchGroup.search_text || ""
    };
    const res: any = await API.categoryList(param);
    if (res.code == 200) {
      this.categoryTable.table = res.data.list || [];
      this.categoryTable.total = res.data.total;
      if (this.categoryTable.table.length == 0) {
        this.searchLoading = true;
      } else {
        this.searchLoading = false;
      }
    } else {
      this.$message.error(res.message);
    }
  }
  // 搜索
  private searchFn() {
    this.categoryTable.page = 1;
    this.getTableDataFn();
  }
  // 重置
  private resetFn() {
    this.searchGroup.search_text = "";
    this.categoryTable.page = 1;
    this.categoryTable.size = 10;
    this.getTableDataFn();
  }
  // 编辑
  private editRowFn(row: any) {
    this.dialogType = "2";
    this.selectedCate = row;
    this.addVisible = true;
  }
  // 关闭取消弹窗
  private closeVisible() {
    this.addVisible = false;
  }
  private addsubmit() {
    if (this.dialogType == "2") {
      this.resetFn();
    } else {
      this.getTableDataFn();
    }
  }
  // 新增
  private createFn() {
    this.dialogType = "1";
    this.addVisible = true;
  }
  mounted() {
    this.getTableDataFn();
  }
}
</script>
<style lang="scss" scoped>
.welcome_container {
  box-shadow: 0px 0px 20px 1px rgba(0, 0, 0, 0.08);

  .el-popover-ref {
    pointer-events: auto;
  }

  .el-tag {
    background: #f7f7f7;
    color: #000000;
    border: 1px solid #dcdfe6;
  }

  .container_content {
    background-color: #fff;
    padding: 30px;
    padding-bottom: 10px;
    border-radius: 10px;

    .button-group {
      margin-bottom: 20px;
    }

    .search_content {
      background-color: $color-background;
      padding: 20px;
      padding-bottom: 0px;
      margin-bottom: 20px;
    }

    .section {
      .marginR10 {
        margin-right: 10px;
      }

      .aligns-center {
        display: flex;
        align-items: center;

        .goods-name {
          flex: 1;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
        }
      }

      .deleteBtn.el-button--text {
        color: #f67170;
      }

      .deleteBtn.el-button--text:hover {
        color: #f56c6c;
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

      .table_footer {
        margin: 20px 0;
        text-align: right;
      }
    }
  }
}</style>
