<template>
  <div class="add-manual-label">
    <el-dialog
      title="选择商品"
      :visible="true"
      @close="closeDialogFn"
      :close-on-click-modal="false"
    >
      <!-- 搜索按钮 -->
      <div class="table-body">
        <div class="m-b-30">
          <div class="aligns-center">
            <div class="search-ipt pull-left">
              <cy-input
                size="small"
                width-size="medium"
                :placeholder="selectType == '2' ? '请输入商品名称' : '请输入SPU_ID'
                  "
                v-model.trim="searchInput"
                class="input-with-select"
              >
                <el-select
                  v-model="selectType"
                  slot="prepend"
                >
                  <el-option
                    v-for="item in selectDefault"
                    :label="item.label"
                    :value="item.value"
                    :key="item.value"
                  ></el-option>
                </el-select>
              </cy-input>
            </div>
            <div class="search-btn pull-left m-l-20">
              <el-button
                size="small"
                @click="searchFn()"
              >查询</el-button>
              <el-button
                size="small"
                @click="resetFn()"
              >重置</el-button>
            </div>
          </div>
        </div>
        <el-table
          ref="multipleTable"
          :data="tableData"
          tooltip-effect="dark"
          style="width: 100%"
          stripe
          height="340"
          max-height="340"
          row-key="id"
          @selection-change="handleSelectionChangeFn"
        >
          <el-table-column
            type="selection"
            :selectable="checkSelectableFn"
            :reserve-selection="true"
            width="55"
          >
          </el-table-column>
          <el-table-column
            prop="id"
            label="SPU ID"
          ></el-table-column>
          <el-table-column label="商品">
            <template slot-scope="scope">
              <div class="aligns-center">
                <!--图片-->
                <el-image fit="cover"
                  style="width: 50px; height: 50px"
                  :src="scope.row.spu_source.url"
                ></el-image>
                <el-tooltip
                  :disabled="scope.row.spu_name.length > 33 ? false : true"
                  :content="scope.row.spu_name"
                  placement="bottom"
                  effect="light"
                >
                  <span class="m-l-10 goods-name">{{
                    scope.row.spu_name
                  }}</span>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
          <div
            slot="empty"
            class="table-empty-label"
          >
            <img
              class="empty-img"
              :src="searchLoading
                    ? 'https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/empty-search.png'
                    : require('@/assets/images/empty-box.png')
                  "
              alt=""
            />
            <p>{{ searchLoading ? "查询无结果" : "暂无数据" }}</p>
          </div>
        </el-table>
      </div>
      <span
        slot="footer"
        class="dialog-footer"
      >
        <div class="m-b-20">
          <el-pagination
            @current-change="handleCurrentChangeFn"
            :current-page="meta.page"
            :page-size="meta.size"
            layout="total, prev, pager, next, jumper"
            :total="meta.total"
          >
          </el-pagination>
        </div>
        <el-button
          size="small"
          @click="closeDialogFn"
        >取 消</el-button>
        <el-button
          :disabled="selectList.length == 0"
          size="small"
          type="primary"
          @click="submitFn"
        >确 认</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import API from "@/api/customer";
@Component({})
export default class extends Vue {
  @Prop({ default: () => [] }) hasSelectArray?: Array<any>;
  /************data**************/
  // 搜索类型
  private selectDefault = [
    { label: "商品名称", value: "2" },
    { label: "SPU ID", value: "1" }
  ];

  private selectType = "2"; // 搜索类型默认选中
  private searchInput = ""; // 搜索内容
  private searchLoading = false; // 搜索框loading
  private tableData: any = []; // 列表数据
  private selectList: any = []; // 用户选择的数据
  private meta = {
    size: 5,
    page: 1,
    total: 0
  };
  /************methods**************/
  // 获取商品列表接口
  private getListFn() {
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      search_type: this.searchInput ? parseInt(this.selectType) : undefined,
      search: this.searchInput || undefined
    };
    API.tagGoodsSpuApi(params).then((res: any) => {
      if (res.code === 200) {
        const data = res.data.list;
        this.meta = {
          page: res.data.page,
          size: res.data.size,
          total: res.data.total
        };
        // 成功
        this.tableData = [...data];
      } else {
        this.$message.error(res.message);
      }
    });
  }

  // 页码选择回调
  handleCurrentChangeFn(val: number) {
    this.meta.page = val;
    this.getListFn();
  }
  // 勾选或点击权限
  private handleSelectionChangeFn(rows: any) {
    this.selectList = rows;
  }
  // 判断多选框是否可选
  private checkSelectableFn(row: any) {
    if (this.hasSelectArray && this.hasSelectArray.length > 0) {
      const arr = this.hasSelectArray.map(item => item.id);
      return !arr.includes(row.id);
    } else {
      return true;
    }
  }
  // 搜索框事件 搜索
  private async searchFn() {
    this.meta.page = 1;
    this.getListFn();
    this.searchLoading = true;
  }
  // 搜索框事件 重置
  private resetFn() {
    this.searchInput = "";
    this.selectType = "2";
    this.getListFn();
    this.searchLoading = false;
  }
  // 点击完成
  private submitFn() {
    this.$emit(
      "addShopsDialog",
      this.hasSelectArray && this.hasSelectArray.concat(this.selectList)
    );
  }
  // 关闭弹框
  private closeDialogFn() {
    this.$emit("addShopsDialog");
  }
  private created() {
    this.getListFn();
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-dialog {
  height: 620px;
}

/deep/.el-dialog__body {
  height: 400px;
}

/deep/.el-checkbox__input.is-disabled .el-checkbox__inner {
  background-color: #c0c4cc !important;
  border-color: #c0c4cc !important;
}

/deep/.gutter {
  width: 0px;
  display: none;
}

.table-body {
  display: flex;
  flex-direction: column;

  /deep/.el-select .el-input {
    width: 110px;
  }

  /deep/.el-table__body-wrapper {
    height: 304px !important;
  }

  .table-empty-label {
    display: flex;
    height: 340px;
    flex-direction: column;

    img {
      width: 124px;
      margin: 80px auto 0;
    }

    p {
      line-height: 20px;
      margin-top: 14px;
    }
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
}</style>
