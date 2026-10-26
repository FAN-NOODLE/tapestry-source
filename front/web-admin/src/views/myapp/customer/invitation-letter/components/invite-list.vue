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
          <search-item text="关联预约服务：">
            <cy-select
              width-size="small"
              search-input
              size="small"
              v-model="searchGroup.reserve_service_id"
              placeholder="请选择预约服务"
              :options-list="associateServiceList"
              :options="{ label: 'reserve_name', value: 'id' }"
            >
            </cy-select>
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
          新建邀请函
        </el-button>
      </div>
      <div class="section">
        <el-table
          :data="cardTable.table"
          stripe
          tooltip-effect="light"
        >
          <el-table-column
            label="模板ID"
            prop="id"
            show-overflow-tooltip
            min-width="180"
          >
          </el-table-column>
          <el-table-column
            label="模板名称"
            prop="title_name"
            show-overflow-tooltip
            min-width="180"
          >
          </el-table-column>
          <el-table-column
            label="关联预约服务"
            prop="reserve_service_name"
            show-overflow-tooltip
            min-width="140"
          >
          </el-table-column>
          <el-table-column
            label="邀请函状态"
            prop="status_str"
            show-overflow-tooltip
            min-width="140"
          >
          </el-table-column>
          <el-table-column
            label="适用范围"
            prop="third_names"
            show-overflow-tooltip
            min-width="200"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.reserve_service_role === 0">全部</span>
              <span
                v-else
                v-for="(item, index) in scope.row.role_rules"
                :key="index"
              >
                {{ item.name }}，
              </span>
            </template>
          </el-table-column>
          <el-table-column
            label="邀请函有效期"
            prop="status_time_str"
            show-overflow-tooltip
            min-width="340"
          >
            <template slot-scope="scope">
              <span>
                {{
                  scope.row.status_type === 0
                  ? "长期有效"
                  : scope.row.status_time_str
                }}
              </span>
            </template>
          </el-table-column>
          <el-table-column
            label="关联预约单数"
            prop="invite_order_num"
            show-overflow-tooltip
            min-width="140"
          >
          </el-table-column>

          <el-table-column
            label="更新时间"
            prop="updated_at"
            show-overflow-tooltip
            min-width="200"
          >
            <template slot-scope="scope">
              <span>
                {{ scope.row.updated_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
              </span>
            </template>
          </el-table-column>
          <el-table-column
            label="操作"
            fixed="right"
            width="230"
          >
            <template slot-scope="scope">
              <el-button
                @click="handleToView(scope.row.id)"
                type="text"
                size="small"
              >查看</el-button>
              <el-button
                @click="handleEdit(scope.row.id)"
                type="text"
                size="small"
                :disabled="scope.row.status_type == 1 &&
                  scope.row.end_time * 1000 <= Date.now()"
              >编辑</el-button>
              <el-button
                type="text"
                size="small"
                @click="handleToViewData(scope.row)"
              >查看数据</el-button>

              <el-dropdown
                class="del-dropdown"
                @command="handleDelete(scope.row.id)"
              >
                <span class="el-dropdown-link">
                  ...
                </span>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item>删除</el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>
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
            <div>
              {{ searchLoading ? "没有找到该邀请函" : "暂无邀请函" }}
            </div>
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
import { Component, Vue, Emit } from "vue-property-decorator";
import API from "@/api/customer/invitation-letter";
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
    search_text: "",
    reserve_service_id: ""
  };
  private searchLoading = false;

  /** 关联预约服务数据 */
  private associateServiceList = [];

  /** 生命周期 */

  async mounted() {
    await this.getAssociateServiceList();
    this.getTableDataFn();
  }

  /** methods */

  /** 获取数据 */
  async getTableDataFn() {
    const param: object = {
      page: this.cardTable.page,
      size: this.cardTable.size,
      keyword: this.searchGroup.search_text || "",
      reserve_service_id: this.searchGroup.reserve_service_id || 0
    };
    const res: any = await API.invitationLetterList(param);
    if (res.code == 200) {
      this.cardTable.table = res.data.list || [];
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

  /** 获取关联预约服务列表 */
  async getAssociateServiceList() {
    const res: any = await API.associateServiceList({
      page: 1,
      size: 1999,
      status: 1
    });
    if (res.code == 200) {
      this.associateServiceList = res.data.list || [];
    } else {
      this.$message.error(res.message);
    }
  }

  /** 搜索 */
  private searchFn() {
    this.cardTable.page = 1;
    this.getTableDataFn();
  }

  /** 重置 */
  private resetFn() {
    this.searchGroup.search_text = "";
    this.searchGroup.reserve_service_id = "";
    this.cardTable.page = 1;
    this.cardTable.size = 10;
    this.getTableDataFn();
  }

  /** 新增 */
  private createFn() {
    this.$router.push({
      name: "myapp_customer_invitation-letter_create"
    });
  }

  /** 查看详情 */
  private handleToView(id: any) {
    this.$router.push({
      name: "myapp_customer_invitation-letter_detail",
      query: { id: id }
    });
  }

  /** 编辑 */
  private handleEdit(id: any) {
    this.$router.push({
      name: "myapp_customer_invitation-letter_edit",
      query: { id: id }
    });
  }

  /** 删除 */
  private handleDelete(id: any) {
    const title = "是否确定删除当前邀请函？";
    this.$confirm(title, "删除邀请函", {
      confirmButtonText: "确认",
      cancelButtonText: "取消",
      type: "warning"
    }).then(async () => {
      API.invitationDataDelete({ id: id }).then((res: any) => {
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
    });
  }

  /** 查看数据 */
  @Emit("handleToViewData")
  private handleToViewData(data: any) {
    this.$emit("handleToViewData", data);
  }

  /** 每页显示数量改变 */
  private handleSizeChangeFn(val: number) {
    this.cardTable.page = 1;
    this.cardTable.size = val;
    this.getTableDataFn();
  }

  /** 页数改变 */
  private handlePageChangeFn(val: number) {
    this.cardTable.page = val;
    this.getTableDataFn();
  }
}
</script>
<style lang="scss" scoped>
.del-dropdown {
  cursor: pointer;
  margin-left: 15px;
}

.danger-txt {
  color: rgba($color: red, $alpha: 0.8) !important;
}

.welcome_container {
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
        margin: 20px 0 0;
        text-align: right;
      }
    }
  }
}</style>
