<template>
  <div class="welcome_container">
    <div class="container_content">
      <div class="search_content">
        <div class="line">
          <search-item text="模板搜索：">
            <cy-input
              width-size="small"
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入模板名称"
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
          新建贺卡
        </el-button>
      </div>
      <div class="section">
        <el-table
          :data="cardTable.table"
          stripe
          tooltip-effect="light"
        >
          <el-table-column
            label="模板名称"
            prop="template_name"
            show-overflow-tooltip
          >
          </el-table-column>
          <el-table-column
            label="模板ID"
            prop="id"
            show-overflow-tooltip
          >
          </el-table-column>
          <el-table-column
            label="更新时间"
            prop="updated_at"
            show-overflow-tooltip
          >
            <template slot-scope="scope">
              {{ scope.row.updated_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column
            label="操作"
            fixed="right"
          >
            <template slot-scope="scope">
              <el-button
                @click="templateCheck(scope.row.id)"
                type="text"
                style="color:#303133;"
                size="small"
              >删除</el-button>
              <el-button
                type="text"
                size="small"
                @click="editRowFn(scope.row.id)"
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
            <div>{{ searchLoading ? "没有找到该贺卡" : "暂无贺卡" }}</div>
          </div>
        </el-table>
        <div class="table_footer">
          <el-pagination
            v-if="cardTable.total > 0"
            @size-change="handleSizeChangeFn"
            @current-change="handlePageChangeFn"
            :current-page="cardTable.page"
            :page-sizes="[10, 20, 30, 40]"
            :page-size="cardTable.size"
            layout="total, sizes, prev, pager, next, jumper"
            :total="cardTable.total"
          >
          </el-pagination>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import API from "@/api/customer/greeting-card";
@Component({})
export default class extends Vue {
  private cardTable = {
    table: [],
    size: 10,
    page: 1,
    total: 1
  };
  private selectedCate = null;
  private addVisible = false; // 控制添加素材标签弹窗
  private dialogType = "1";
  private searchGroup = {
    search_text: ""
  };
  private searchLoading = false;
  private handleSizeChangeFn(val: number) {
    this.cardTable.page = 1;
    this.cardTable.size = val;
    this.getTableDataFn();
  }
  private handlePageChangeFn(val: number) {
    this.cardTable.page = val;
    this.getTableDataFn();
  }

  //检查贺卡模板是否在使用
  private templateCheck(id: any) {
    API.templateCheck({ id: id }).then((res: any) => {
      if (res.code === 200) {
        this.delItem(id, res.data.check);
      } else {
        this.$message.error(res.message);
      }
    });
  }

  //删除
  private delItem(id: any, check: any) {
    let title = "确定删除该贺卡吗？";
    if (!check) {
      title = "该贺卡已被子任务引用，无法删除！";
    }
    this.$confirm(title, "删除贺卡", {
      confirmButtonText: "确认",
      cancelButtonText: "取消",
      type: "warning"
    }).then(async () => {
      if (check) {
        API.templateDelete({ id: id }).then((res: any) => {
          if (res.code === 200) {
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
    });
  }
  // 获取数据
  async getTableDataFn() {
    const param: object = {
      page: this.cardTable.page,
      size: this.cardTable.size,
      template_name: this.searchGroup.search_text || ""
    };
    const res: any = await API.templateCard(param);
    if (res.code == 200) {
      this.cardTable.table = res.data.List || [];
      this.cardTable.total = res.data.total;
      if (this.cardTable.table.length == 0) {
        this.searchLoading = true;
      } else {
        this.searchLoading = false;
      }
    } else {
      this.$message.error(res.message);
    }
  }
  // 编辑
  private editRowFn(id: any) {
    this.$router.push({
      name: "myapp_customer_greeting-card_edit",
      query: { id: id }
    });
  }
  // 搜索
  private searchFn() {
    this.cardTable.page = 1;
    this.getTableDataFn();
  }
  // 重置
  private resetFn() {
    this.searchGroup.search_text = "";
    this.cardTable.page = 1;
    this.cardTable.size = 10;
    this.getTableDataFn();
  }
  // 新增
  private createFn() {
    this.$router.push({
      name: "myapp_customer_greeting-card_create"
    });
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
