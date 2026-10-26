<template>
  <div class="welcome_container">
    <div class="container_content">
      <div class="button-group">
        <el-button
          type="primary"
          @click="createFn"
          size="small"
        >
          新增渠道
        </el-button>
      </div>
      <div class="search_content">
        <div class="line">
          <search-item text="渠道名称：">
            <cy-input
              width-size="small"
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入渠道名称"
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
      <div class="section">
        <el-table
          :data="channelTable.table"
          stripe
          tooltip-effect="light"
        >
          <el-table-column
            label="渠道名称"
            min-width="300"
            prop="name"
            fixed="left"
            show-overflow-tooltip
          >
          </el-table-column>
          <el-table-column
            label="标签"
            min-width="300"
          >
            <template slot-scope="{ row }">
              <div v-if="!row.channel_labels || !row.channel_labels.length">
                -
              </div>
              <div v-else-if="row.channel_labels.length < 3">
                <el-tag
                  v-for="(item, index) in row.channel_labels"
                  :key="index"
                  type="info"
                  style="margin-right: 5px;"
                  effect="plain"
                >
                  {{ item }}
                </el-tag>
              </div>
              <div v-else>
                <el-tag
                  v-for="(item, index) in row.channel_labels.slice(0, 2)"
                  :key="index"
                  type="info"
                  style="margin-right: 5px;"
                  effect="plain"
                >
                  {{ item }}
                </el-tag>
                <el-popover
                  placement="top-start"
                  width="400"
                  trigger="hover"
                >
                  <el-tag
                    v-for="(item, index) in row.channel_labels"
                    :key="index"
                    size="small"
                    style="margin:5px;"
                  >
                    {{ item }}
                  </el-tag>
                  <el-button
                    type="text"
                    slot="reference"
                    icon="el-icon-more"
                    class="el-popover-ref"
                  ></el-button>
                </el-popover>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="created_at"
            label="创建时间"
            :formatter="formatDate"
            min-width="200"
          >
          </el-table-column>
          <el-table-column
            label="操作"
            width="72"
            fixed="right"
          >
            <template slot-scope="scope">
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
            <div>{{ searchLoading ? "没有找到该渠道" : "暂无渠道" }}</div>
          </div>
        </el-table>
        <div class="table_footer">
          <el-pagination
            v-if="channelTable.total > 0"
            @size-change="handleSizeChangeFn"
            @current-change="handlePageChangeFn"
            :current-page="channelTable.page"
            :page-sizes="[10, 20, 30, 40]"
            :page-size="channelTable.size"
            layout="total, sizes, prev, pager, next, jumper"
            :total="channelTable.total"
          >
          </el-pagination>
        </div>
      </div>
    </div>
    <create-or-edit-channel-dialog
      v-if="showCreateOrEditDialog"
      :selectedChannel="selectedChannel"
      :dialogType="dialogType"
      @channelDialog="channelDialogFn"
    ></create-or-edit-channel-dialog>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { dateFormat } from "../../../../customer/auto-market/util/index";
import createOrEditChannelDialog from "../../../components/create-or-eidt-channel-dialog.vue";
import API from "@/api/enterprise/drainage-coach";
@Component({
  components: {
    createOrEditChannelDialog
  }
})
export default class extends Vue {
  private channelTable = {
    table: [],
    size: 10,
    page: 1,
    total: 0
  };
  private selectedChannel = null;
  private showCreateOrEditDialog = false;
  private dialogType = "1";
  private searchGroup = {
    search_text: ""
  };
  private searchLoading = false;
  handleSizeChangeFn(val: number) {
    this.channelTable.page = 1;
    this.channelTable.size = val;
    this.getTableDataFn();
  }
  handlePageChangeFn(val: number) {
    this.channelTable.page = val;
    this.getTableDataFn();
  }
  formatDate(row: any) {
    return row.created_at
      ? dateFormat(row.created_at, "yyyy-MM-dd hh:mm:ss")
      : "";
  }

  async getTableDataFn() {
    const param: object = {
      page: this.channelTable.page,
      size: this.channelTable.size,
      channel_name: this.searchGroup.search_text || ""
    };
    const res: any = await API.channelsListApi(param);
    if (res.code == 200) {
      this.channelTable.table = res.data.list || [];
      this.channelTable.total = res.data.total;
      if (this.channelTable.table.length == 0) {
        this.searchLoading = true;
      } else {
        this.searchLoading = false;
      }
    } else {
      this.$message.error(res.message);
    }
  }
  searchFn() {
    this.channelTable.page = 1;
    this.getTableDataFn();
  }
  resetFn() {
    this.searchGroup.search_text = "";
    this.channelTable.page = 1;
    this.channelTable.size = 10;
    this.getTableDataFn();
  }
  editRowFn(row: any) {
    this.dialogType = "2";
    this.showCreateOrEditDialog = true;
    this.selectedChannel = row;
  }
  createFn() {
    // console.log()
    this.dialogType = "1";
    this.showCreateOrEditDialog = true;
  }
  channelDialogFn() {
    this.showCreateOrEditDialog = false;
    this.selectedChannel = null;
    this.getTableDataFn();
  }
  mounted() {
    this.getTableDataFn();
    //两个跳转其中之一是欢迎语
    //从欢迎语跳转的渠道列表，并打开新建渠道弹框
    const model: any = this.$route.query.showAddModal;
    if (model && model == 1) {
      this.createFn();
    }
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
