<template>
  <div>
    <el-dialog
      :title="title"
      width="1000px"
      min-height="900px"
      :visible.sync="dialogVisible"
      :close-on-click-modal="false"
      @close="closeDialog"
    >
      <div>
        <search-item text="服务名称：">
          <cy-input
            width-size="small"
            size="small"
            prefix-icon="el-icon-search"
            placeholder="请输入服务名称"
            clearable
            v-model="searchGroup.service_name"
          ></cy-input>
        </search-item>
        <el-button
          type="primary"
          @click="handleSearch"
          size="small"
        >查询</el-button>
        <el-button
          @click="handleReset"
          size="small"
        >重置</el-button>
      </div>

      <el-table
        :data="dataList"
        row-key="id"
        height="333px"
      >
        <el-table-column
          label="服务ID"
          min-width="180"
          key="id"
          prop="id"
        >
          <template slot-scope="scope">
            <el-radio
              v-model="chooseData.id"
              :label="scope.row.id"
              @change="handleSelectionChange(scope.row)"
            ></el-radio>
          </template>
        </el-table-column>
        <el-table-column
          label="服务名称"
          prop="reserve_name"
        >
        </el-table-column>
        <el-table-column
          label="时长"
          prop="reserve_time"
        >
          <template slot-scope="scope">
            <span>{{ scope.row.reserve_time }}小时</span>
          </template>
        </el-table-column>
        <el-table-column
          label="图片"
          min-width="140"
        >
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
          label="有效期"
          prop="limit_time"
          min-width="160"
        >
          <template slot-scope="{ row }">
            <span v-if="row.service_type === 1">长期有效</span>
            <span v-else>
              {{ row.start_time | formatDate("yyyy-mm-dd") }}
              ~
              {{ row.end_time | formatDate("yyyy-mm-dd") }}
            </span>
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
          <div>{{ "暂无数据" }}</div>
        </div>
      </el-table>
      <div class="table_footer">
        <el-pagination
          v-if="cardTable.total > 0"
          :current-page="cardTable.page"
          :page-size="cardTable.size"
          layout="prev, pager, next"
          :total="cardTable.total"
        >
        </el-pagination>
      </div>
      <div
        slot="footer"
        class="dialog-footer"
      >
        <template>
          <el-button
            @click="handleCancel"
            size="small"
          >
            取消
          </el-button>
          <el-button
            @click="handleConfirm"
            size="small"
            type="primary"
          >
            确 认
          </el-button>
        </template>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from "vue-property-decorator";
import API from "@/api/customer/invitation-letter";

@Component
export default class extends Vue {
  @Prop({ default: "选择预约服务" }) title?: string;
  @Prop() selected?: any;

  /*************************** data start ***************************/

  /** 搜素数据 */
  private searchGroup = {
    service_name: ""
  };

  /** 表格数据 */
  private dataList = [];

  /** 分页数据 */
  private cardTable = {
    table: [],
    size: 10,
    page: 1,
    total: 1
  };

  /** 已选择数据 */
  private chooseData = {
    id: 0
  };

  /** 是否搜索 */
  private searchLoading = false;

  private dialogVisible = true;

  /*** methods START ***/

  private mounted() {
    this.selected && (this.chooseData = { ...this.selected });
    this.handleSearch();
  }

  /** 查询列表数据 */
  private handleSearch() {
    this.searchLoading = true;
    const params = {
      page: this.cardTable.page,
      size: this.cardTable.size,
      reserve_name: this.searchGroup.service_name || "",
      status: 1
    };
    API.associateServiceList(params).then((res: any) => {
      if (res.code === 200) {
        this.dataList = res.data.list || [];
        this.cardTable.total = res.data.total;
        const ids = this.dataList.map(item => item.id);
        console.log(ids);
      } else {
        this.$message.error(res.message);
      }
    });
  }

  /** 重置 */
  private handleReset() {
    this.searchGroup = {
      service_name: ""
    };
    this.handleSearch();
  }

  /** 表格数据选择 */
  private handleSelectionChange(val: any) {
    this.chooseData = { ...val };
  }

  /** 取消 */
  private handleCancel() {
    this.closeDialog();
  }

  /** 确认 */
  private handleConfirm() {
    this.$emit("submit", this.chooseData);
  }

  /** 关闭模态框 */
  @Emit("close")
  private closeDialog() {
    this.$emit("close");
  }

  /** methods END */
}
</script>

<style lang="scss" scoped>.table-empty-label {
  display: flex;
  height: 280px;
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
}</style>
