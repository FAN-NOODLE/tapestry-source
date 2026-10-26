<template>
  <div class="distribution-goods">
    <!-- 搜索部分 -->
    <div class="search-content">
      <search-item text="商品：">
        <el-input
          :placeholder="searchPlaceholder"
          class="search-input"
          size="medium"
          clearable
          v-model="searchValue"
          @keyup.enter.native="search"
        >
          <el-select
            slot="prepend"
            v-model="searchType"
            class="input-with-select"
          >
            <el-option
              v-for="item in searchTypeOptions"
              :value="item.key"
              :label="item.name"
              :key="item.key"
            ></el-option
          ></el-select>
        </el-input>
      </search-item>
      <el-button class="m-l-20" type="primary" size="small" @click="search">
        查询
      </el-button>
      <el-button class="m-l-20" size="small" @click="reset">重置</el-button>
    </div>

    <!-- 批量操作部分 -->
    <div class="operate-content">
      <el-button type="primary" size="small" @click="add">
        添加分销商品
      </el-button>
      <el-button size="small" @click="batchSetting">批量设置</el-button>
      <el-button size="small" @click="batchRemove">批量移除分销商品</el-button>
    </div>

    <!-- 表格部分 -->
    <div class="table-content">
      <el-table
        v-if="isLoaded"
        class="table-wrapper"
        ref="goodList"
        stripe
        header-row-class-name="header-row"
        :data="list"
        @selection-change="tableSelectionChange"
      >
        <el-table-column
          type="selection"
          class-name="selctRow"
          width="55"
        ></el-table-column>
        <el-table-column
          prop="sku_id"
          label="SKU ID"
          width="170"
          fixed="left"
        ></el-table-column>
        <el-table-column label="商品" width="320">
          <template slot-scope="{ row }">
            <div class="goods-warpper">
              <img :src="row.sku_img" alt="" class="goods-img" />
              <div class="goods-info">
                <p class="goods-name">{{ row.spu_name }}</p>
                <span class="goods-attr">{{ row.sku_attr }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          prop="sign_price"
          label="划线价（元）"
          min-width="130"
          align="right"
        >
          <template slot-scope="{ row }">{{ row.sign_price / 100 }}</template>
        </el-table-column>
        <el-table-column
          prop="sell_price"
          label="售卖价（元）"
          min-width="130"
          align="right"
        >
          <template slot-scope="{ row }">{{ row.sell_price / 100 }}</template>
        </el-table-column>
        <el-table-column label="推广佣金比例" min-width="230" align="right">
          <template slot-scope="{ row }">
            <span v-if="!row.isEditing">
              {{ row.promotion_commission_amount }}%
            </span>
            <el-input
              v-else
              placeholder="请输入推广佣金比例"
              v-model="row.promotion_commission_amount_new"
              size="small"
              @keyup.native="limitInput(row, 'promotion_commission_amount_new')"
            >
              <template slot="append">%</template>
            </el-input>
          </template>
        </el-table-column>
        <el-table-column
          v-if="USER.version !== 3"
          label="服务佣金比例"
          min-width="230"
          align="right"
        >
          <template slot-scope="{ row }">
            <span v-if="!row.isEditing">
              {{ row.service_commission_amount }}%
            </span>
            <el-input
              v-else
              placeholder="请输入服务佣金比例"
              v-model="row.service_commission_amount_new"
              size="small"
              @keyup.native="limitInput(row, 'service_commission_amount_new')"
            >
              <template slot="append">%</template>
            </el-input>
          </template>
        </el-table-column>
        <el-table-column label="分销期限" min-width="300">
          <template slot-scope="{ row }">
            <span v-if="!row.isEditing">
              {{ row.distribution_start | formatDate("yyyy-mm-dd") }} 至
              {{ row.distribution_end | formatDate("yyyy-mm-dd") }}
            </span>
            <el-date-picker
              v-else
              v-model="row.date_range_new"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="timestamp"
              size="small"
              style="width: 280px !important;"
              :picker-options="row.pickerOption"
              @change="row.handleDateChange()"
            ></el-date-picker>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template slot-scope="{ row }">
            <el-button
              v-if="!row.isEditing"
              type="text"
              @click="handleEditRow(row)"
            >
              编辑
            </el-button>
            <el-button
              v-if="row.isEditing"
              type="text"
              @click="handleSaveRow(row)"
            >
              保存
            </el-button>
            <el-button
              v-if="row.isEditing"
              size="small"
              type="text"
              @click="handleCancelRow(row)"
            >
              <span class="table-cancel-btn">取消</span>
            </el-button>
            <el-popconfirm
              class="mar-l-10"
              title="您确认移除该分销商品吗？"
              confirm-button-text="确认"
              @onConfirm="handleRemoveRow(row)"
            >
              <el-button class="m-l-20" slot="reference" type="text">
                <span class="table-remove-btn">移除</span>
              </el-button>
            </el-popconfirm>
          </template>
        </el-table-column>
        <div slot="empty">
          <layout-empty
            v-if="!searchValue"
            :message="'暂无分销商品'"
          ></layout-empty>
          <layout-empty
            v-if="searchValue"
            empty_type="tag"
            :message="'没有找到该分销商品'"
          ></layout-empty>
        </div>
      </el-table>
    </div>

    <!-- 分页部分 -->
    <div class="footer-content">
      <el-pagination
        v-if="isLoaded && list.length > 0"
        class="pagination-wrapper"
        layout="total, sizes, prev, pager, next, jumper"
        :page-sizes="[10, 20, 30, 40]"
        :total="meta.total"
        :page-size="meta.size"
        :current-page="meta.page"
        @current-change="handleCurrentPageChange"
        @size-change="handleSizeChange"
      ></el-pagination>
    </div>

    <!-- 批量设置弹框 -->
    <batch-setting
      v-if="batchSettingVisible"
      :visible="batchSettingVisible"
      :selectSkuNum="selectSkuNum"
      :display="batchSettingDisplay"
      @submit="handleBatchSettingDone"
      @close="closeBatchSetting"
    ></batch-setting>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import UserModule from "@/store/modules/user";
import BatchSetting from "../modal/batch-setting.vue";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
import DistributionApi from "@/api/wisdom/Distribution";
import { handleDisabledDate } from "../date-picker";
@Component({
  components: {
    BatchSetting,
    LayoutEmpty
  }
})
export default class DistributionGoods extends Vue {
  /*** data ***/
  private searchTypeOptions: Array<any> = [
    { name: "商品名称", key: 2 },
    { name: "SKU ID", key: 1 }
  ]; // 搜索下拉框选项
  private searchValue = ""; // 搜索输入框的值
  private searchType = 2; // 当前选中的搜索类型
  private list: Array<any> = []; // 列表
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private isLoaded = false; // 数据是否加载完成
  private selection: Array<any> = []; // 列表选中的数据
  private batchSettingVisible = false; // 控制批量设置弹框显示

  /*** 计算属性 ***/
  get USER() {
    return UserModule.getUserInfo;
  }
  get searchPlaceholder() {
    return this.searchType === 2 ? "请输入商品名称" : "请输入SKU ID";
  }
  get selectSkuNum() {
    return this.selection.length;
  }
  get batchSettingDisplay() {
    const display = ["promotion_commission_amount"];
    if (this.USER.version !== 3) {
      display.push("service_commission_amount");
    }
    return display;
  }

  /*** methods ***/
  // 获取列表数据
  async loadData() {
    const params = {
      page: this.meta.page,
      limit: this.meta.size,
      spu_name: this.searchType === 2 ? this.searchValue : undefined,
      sku_id: this.searchType === 1 ? this.searchValue : undefined
    };
    try {
      const { data } = await DistributionApi.getDistributionGoodsList(params);
      this.list = data.list.map((item: any) => ({
        ...item,
        promotion_commission_amount: item.promotion_commission_amount / 100,
        service_commission_amount: item.service_commission_amount / 100,
        sku_attr: (item.attribute_item_names || []).join(","),
        pickerOption: {
          disabledDate: handleDisabledDate(item.date_range),
          onPick({ minDate }: { minDate: any }) {
            this.disabledDate = handleDisabledDate(item.date_range, minDate);
          }
        },
        handleDateChange() {
          this.pickerOption.disabledDate = handleDisabledDate(item.date_range);
        },
        isEditing: false,
        promotion_commission_amount_new: item.promotion_commission_amount / 100,
        service_commission_amount_new: item.service_commission_amount / 100,
        date_range_new: [
          item.distribution_start * 1e3,
          item.distribution_end * 1e3
        ]
      }));
      this.meta.total = data.total;
    } catch (error) {
      console.log(error);
    }
    this.isLoaded = true;
  }

  // 搜索
  private search() {
    this.cancelCurrentEditing();
    this.meta.page = 1;
    this.loadData();
  }

  // 重置
  private reset() {
    this.cancelCurrentEditing();
    this.searchType = 2;
    this.searchValue = "";
    this.meta.page = 1;
    this.loadData();
  }

  // 当前页发生变化
  private handleCurrentPageChange(page: number) {
    this.cancelCurrentEditing();
    this.meta.page = page;
    this.loadData();
  }

  private handleSizeChange(size: number) {
    this.cancelCurrentEditing();
    this.meta.page = 1;
    this.meta.size = size;
    this.loadData();
  }

  // table选中项发生变化时
  private tableSelectionChange(selection: any) {
    this.selection = selection;
  }

  // 开启单行编辑
  private handleEditRow(row: any) {
    this.cancelCurrentEditing();
    this.list = this.list.map(item => {
      if (item.id === row.id) {
        item.isEditing = true;
      }
      return item;
    });
  }

  // 取消编辑
  private handleCancelRow() {
    this.cancelCurrentEditing();
  }

  // 保存
  async handleSaveRow(row: any) {
    if (
      row.promotion_commission_amount_new === "" ||
      row.promotion_commission_amount_new == undefined
    ) {
      this.$message.warning("请填写推广佣金比例");
      return false;
    }
    if (
      row.service_commission_amount_new === "" ||
      row.service_commission_amount_new == undefined
    ) {
      this.$message.warning("请填写服务佣金比例");
      return false;
    }
    if (!row.date_range_new?.every((item: any) => item > 0)) {
      this.$message.warning("请填写分销期限");
      return false;
    }
    const params = {
      distribution_goods_id: row.id,
      promotion_commission_amount: Math.round(
        row.promotion_commission_amount_new * 100
      ),
      service_commission_amount: Math.round(
        row.service_commission_amount_new * 100
      ),
      distribution_start: row.date_range_new[0] / 1e3,
      distribution_end: row.date_range_new[1] / 1e3
    };
    try {
      await DistributionApi.updateDistributionGoods(params);
      this.$message.success("保存成功");
      await this.loadData();
      return true;
    } catch (error) {
      console.log(error);
      return false;
    }
  }

  // 取消当前编辑
  cancelCurrentEditing() {
    for (let i = 0; i < this.list.length; i++) {
      const item = this.list[i];
      if (item.isEditing) {
        item.isEditing = false;
        item.promotion_commission_amount_new = item.promotion_commission_amount;
        item.service_commission_amount_new = item.service_commission_amount;
        item.date_range_new = [
          item.distribution_start * 1e3,
          item.distribution_end * 1e3
        ];
        return;
      }
    }
  }

  // 限制小数点位数与值范围
  limitInput(row: any, key: string) {
    let value = row[key];
    value = value > 100 ? 100 : value;
    value = value < 0 ? 0 : value;
    value = value + "";
    if (value.length > 1 && value[0] === "0" && value[1] !== ".") {
      value = value.substring(1);
    }
    row[key] = value
      .replace(/[^\d.]/g, "")
      .replace(/^(\d+)\.(\d\d).*$/, "$1.$2");
  }

  // 添加分销商品
  private add() {
    this.cancelCurrentEditing();
    this.$router.push({
      name: "myapp_wisdom_staff-distribution_distribution-setting_add-goods"
    });
  }

  // 批量设置
  private batchSetting() {
    this.cancelCurrentEditing();
    if (!this.selection.length) {
      this.$message.warning("请先选择分销商品!");
      return;
    }
    this.batchSettingVisible = true;
  }

  // 关闭批量设置弹窗
  private closeBatchSetting() {
    this.batchSettingVisible = false;
  }

  // 批量设置的回调
  private handleBatchSettingDone(form: any) {
    const params = {
      distribution_goods_ids: this.selection.map(item => item.id),
      service_commission_amount: Math.round(
        form.service_commission_amount * 100
      ),
      promotion_commission_amount: Math.round(
        form.promotion_commission_amount * 100
      )
    };
    try {
      DistributionApi.batchSettingDistributionGoods(params).then(() => {
        this.$message.success("设置成功");
        this.closeBatchSetting();
        this.loadData();
      });
    } catch (error) {
      console.log(error);
    }
  }
  // 批量删除
  private batchRemove() {
    this.cancelCurrentEditing();
    if (!this.selection.length) {
      this.$message.warning("请先选择要移除的商品!");
      return;
    }
    this.$confirm("你确认要移除分销商品？", "移除确认", {
      confirmButtonText: "确认",
      cancelButtonText: "取消"
    }).then(() => {
      this.batchDeleteDistributionGoods(this.selection);
    });
  }

  // 移除当前分销商品
  private handleRemoveRow(row: any) {
    // 如果当前是编辑状态直接删除
    if (row.isEditing) {
      this.batchDeleteDistributionGoods([row]);
      return;
    }
    this.cancelCurrentEditing();
    this.batchDeleteDistributionGoods([row]);
  }

  // 批量删除商品
  async batchDeleteDistributionGoods(goods: Array<any>) {
    const params = {
      distribution_goods_ids: goods.map((item: any) => item.id)
    };
    try {
      await DistributionApi.batchDeleteDistributionGoods(params);
      this.$message.success("移除成功");
      await this.loadData();
    } catch (error) {
      console.log(error);
    }
  }

  /*** 生命周期 ***/
  private mounted() {
    this.loadData();
  }
}
</script>

<style lang="scss" scoped>
.search-content {
  padding: 20px;
  margin-bottom: 20px;
  background-color: $color-background;
  border-radius: 4px;
  .search-item {
    display: inline-flex;
    align-items: center;
    margin-bottom: 0;
    margin-right: 0;
  }
}
.input-with-select {
  width: 110px;
}
.search-input {
  width: 280px !important;
  margin-bottom: 0px !important;
}
.operate-content {
  margin-bottom: 20px;
  color: $color-text-normal;
}
.footer-content {
  &:after {
    content: "";
    display: block;
    clear: both;
  }
  /deep/ {
    .el-pagination {
      float: right;
    }
  }
}
.el-select-dropdown__item.selected {
  font-weight: normal !important;
  height: auto !important;
}
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px;
}
.goods-warpper {
  display: flex;
  .goods-img {
    width: 55px;
    height: 55px;
    margin-right: 10px;
  }
  .goods-name {
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2;
    overflow: hidden;
    font-size: 14px;
    line-height: 18px;
    height: 36px;
  }
  .goods-attr {
    font-size: 12px;
    color: $color-text-second;
  }
}
.table-remove-btn {
  color: #e34b49;
}
.table-cancel-btn {
  color: $color-text-second;
}
/deep/.el-range-editor--small {
  width: 280px !important;
}
</style>
