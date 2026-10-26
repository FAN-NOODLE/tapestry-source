<template>
  <div class="add-goods">
    <layout-form>
      <template slot="content">
        <layout-list>
          <template slot="buttons">
            <el-steps class="steps-wrapper" :active="active" space="50%">
              <el-step title="选择分销商品"></el-step>
              <el-step title="设置分销期限佣金"></el-step>
            </el-steps>
          </template>

          <!-- 搜索部分 -->
          <template slot="search">
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
                  >
                  </el-option>
                </el-select>
              </el-input>
            </search-item>
            <el-button
              class="m-l-30"
              type="primary"
              size="small"
              @click="search"
            >
              查询
            </el-button>
            <el-button class="m-l-20" size="small" @click="reset"
              >重置</el-button
            >
          </template>

          <!-- 批量操作部分 -->
          <template slot="operate">
            <div v-show="active === 1">
              <el-dropdown
                size="medium"
                type="primary"
                @command="batchSetting"
                class="batch-setting"
              >
                <el-button type="primary" size="small">
                  <span class="m-r-5">批量设置</span>
                  <em class="el-icon-arrow-down"></em>
                </el-button>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item command="commission">
                    设置分销佣金
                  </el-dropdown-item>
                  <el-dropdown-item command="date"
                    >设置分销期限</el-dropdown-item
                  >
                </el-dropdown-menu>
              </el-dropdown>
              <el-button class="m-l-20" size="small" @click="batchRemove">
                批量移除
              </el-button>
            </div>
          </template>

          <!-- 表格部分 -->
          <template slot="table">
            <el-table
              v-if="isLoaded"
              class="table-wrapper"
              ref="table"
              stripe
              header-row-class-name="header-row"
              row-key="sku_id"
              :data="list"
              @selection-change="tableSelectionChange"
              @select="handleTableSelect"
              @select-all="handleTableSelectAll"
            >
              <el-table-column
                type="selection"
                width="55px"
                :reserve-selection="true"
              ></el-table-column>
              <el-table-column
                prop="sku_id"
                label="SKU ID"
                min-width="170"
                fixed="left"
              ></el-table-column>
              <el-table-column
                v-if="active !== 1"
                prop="spu_id"
                label="SPU ID"
                width="170"
              ></el-table-column>
              <el-table-column v-else label="" width="1"></el-table-column>
              <el-table-column label="商品" width="300">
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
                prop="market_price"
                label="划线价(元)"
                width="130"
                align="right"
              >
                <template slot-scope="{ row }">{{
                  row.market_price / 100
                }}</template>
              </el-table-column>
              <el-table-column
                v-if="active !== 1"
                prop="sell_price"
                label="售卖价(元)"
                align="right"
              >
                <template slot-scope="{ row }">{{
                  row.sell_price / 100
                }}</template>
              </el-table-column>
              <el-table-column
                v-else
                label="售卖价(元)"
                width="130"
                align="right"
              >
                <template slot-scope="{ row }">{{
                  row.sell_price / 100
                }}</template>
              </el-table-column>
              <el-table-column
                v-if="active === 1"
                prop="spu_name"
                label="推广佣金比例"
                width="240"
              >
                <template slot-scope="{ row }">
                  <el-input
                    placeholder="请输入推广佣金比例"
                    v-model="row.promotion_commission_amount"
                    size="small"
                    @keyup.native="
                      limitInput(row, 'promotion_commission_amount')
                    "
                  >
                    <template slot="append">%</template>
                  </el-input>
                </template>
              </el-table-column>
              <el-table-column
                v-if="active === 1 && USER.version !== 3"
                label="服务佣金比例"
                width="240"
              >
                <template slot-scope="{ row }">
                  <el-input
                    placeholder="请输入服务佣金比例"
                    v-model="row.service_commission_amount"
                    size="small"
                    @keyup.native="limitInput(row, 'service_commission_amount')"
                  >
                    <template slot="append">%</template>
                  </el-input>
                </template>
              </el-table-column>
              <el-table-column
                v-if="active === 1"
                label="分销期限"
                min-width="400"
              >
                <template slot-scope="{ row }">
                  <el-date-picker
                    v-model="row.date_range"
                    type="daterange"
                    range-separator="至"
                    start-placeholder="开始日期"
                    end-placeholder="结束日期"
                    size="small"
                    style="width: 280px !important;"
                    :picker-options="row.pickerOption"
                    @change="row.handleDateChange()"
                  ></el-date-picker>
                </template>
              </el-table-column>
              <div slot="empty">
                <layout-empty
                  v-if="!searchValue && active !== 1"
                  :message="'暂无可添加商品'"
                ></layout-empty>
                <layout-empty
                  v-if="!searchValue && active === 1"
                  :message="'暂无分销商品'"
                ></layout-empty>
                <layout-empty
                  v-if="searchValue"
                  empty_type="tag"
                  :message="'没有找到该分销商品'"
                ></layout-empty>
              </div>
            </el-table>

            <!-- 批量设置弹框 -->
            <batch-setting
              v-if="batchSettingVisible"
              :visible.sync="batchSettingVisible"
              :selectSkuNum="selectSkuNum"
              :display="batchSettingDisplay"
              @submit="handleBatchSettingDone"
              @close="closeBatchSetting"
            ></batch-setting>
          </template>

          <template slot="footer">
            <el-pagination
              v-if="isLoaded && list.length > 0 && active !== 1"
              layout="total, sizes, prev, pager, next, jumper"
              :page-sizes="[20, 40, 60, 80, 100]"
              :total="meta.total"
              :page-size="meta.size"
              :current-page="meta.page"
              @current-change="handleCurrentPageChange"
              @size-change="handleSizeChange"
            ></el-pagination>
          </template>
        </layout-list>
      </template>
      <template slot="footer">
        <div class="button-wrapper" v-if="isLoaded">
          <el-button
            v-if="active !== 1"
            size="small"
            type="primary"
            @click="next"
          >
            下一步
          </el-button>
          <el-button v-if="active !== 1" size="small" @click="cancel">
            取消
          </el-button>
          <el-button
            v-if="active === 1"
            size="small"
            type="primary"
            @click="save"
          >
            保存
          </el-button>
          <el-button v-if="active === 1" size="small" @click="prev">
            上一步
          </el-button>
        </div>
      </template>
    </layout-form>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import UserModule from "@/store/modules/user";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
import BatchSetting from "../modal/batch-setting.vue";
import DistributionApi from "@/api/wisdom/Distribution";
import { handleDisabledDate } from "../date-picker";
import { scrollTop } from "../scroll-top";
@Component({
  components: {
    BatchSetting,
    LayoutEmpty
  }
})
export default class extends Vue {
  /*** data ***/
  private active = 0; // 当前所在步骤
  private searchTypeOptions: Array<any> = [];
  private searchValue = ""; // 搜索输入框的值
  private searchType = 3; // 当前搜索类型
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 20 // 每页数量
  };
  private list: Array<any> = []; // 列表
  private isLoaded = false; // 数据是否加载完成
  private batchSettingDisplay: Array<any> = []; // 批量设置时批量组件显示的字段
  private batchSettingVisible = false; // 控制批量设置弹框显示
  private selection: Array<any> = []; // 列表中选中的数据
  private chooseGoods: Array<any> = []; // 第一步选中的商品

  /*** 计算属性 ***/
  get USER() {
    return UserModule.getUserInfo;
  }
  get searchPlaceholder() {
    const placeholders = {
      1: "请输入SKU ID",
      3: "请输入商品名称",
      4: "请输入SPU ID"
    };
    return placeholders[this.searchType];
  }
  get selectSkuNum() {
    return this.selection.length;
  }

  /*** methods ***/
  // 设置下拉搜索选项
  setSearchTypeOptions() {
    const searchTypeOptions = [
      { name: "商品名称", key: 3 },
      { name: "SKU ID", key: 1 },
      { name: "SPU ID", key: 4 }
    ];
    if (this.active === 0) {
      this.searchTypeOptions = [...searchTypeOptions];
    } else {
      this.searchTypeOptions = [...searchTypeOptions.slice(0, 2)];
    }
  }

  // 请求商品数据
  async loadGoods() {
    const params = {
      search_type: this.searchType,
      search_text: this.searchValue,
      page: this.meta.page,
      size: this.meta.size
    };
    try {
      const { data } = await DistributionApi.getSkuGoodList(params);
      this.list = [...data.list];
      this.meta.total = data.total;
    } catch (error) {
      console.log(error);
    }
    this.isLoaded = true;
  }

  // 当前页发生变化
  handleCurrentPageChange(page: number) {
    this.meta.page = page;
    this.loadGoods();
  }

  // 当前页数量发生变化
  handleSizeChange(size: number) {
    this.meta.page = 1;
    this.meta.size = size;
    this.loadGoods();
  }

  // table选中项发生变化时
  private tableSelectionChange(selection: any) {
    this.selection = selection;
  }

  // 处理单选，限制选中超过100
  handleTableSelect(selection: any) {
    this.handleTableSelectAll(selection);
  }

  // 处理全选，限制选中超过100
  handleTableSelectAll(selection: any) {
    const max = 100;
    if (selection.length <= max) {
      return;
    }
    this.$message.warning("单次最多选择100个商品");
    const arr = selection.slice(max);
    arr.forEach((row: any) => {
      (this.$refs.table as any).toggleRowSelection(row, false);
    });
  }

  // 限制小数点位数2与值范围0~100
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

  // 搜索
  search() {
    if (this.active === 0) {
      this.loadGoods();
      return;
    }
    if (!this.searchValue) {
      this.list = [...this.chooseGoods];
      return;
    }
    this.list = this.chooseGoods.filter(item => {
      if (this.searchType === 3) {
        return item.spu_name.includes(this.searchValue);
      }
      return item.sku_id == this.searchValue;
    });
  }

  // 重置
  reset() {
    this.searchType = 3;
    this.searchValue = "";
    if (this.active === 0) {
      this.meta.page = 1;
    }
    this.search();
  }

  // 批量移除
  batchRemove() {
    if (!this.selection.length) {
      this.$message.warning("请先选择要移除的商品");
      return;
    }
    const selection = this.selection.map(item => item.sku_id);
    this.chooseGoods = this.chooseGoods.filter(
      item => !selection.includes(item.sku_id)
    );
    this.search();
    (this.$refs.table as any).clearSelection();
    this.$message.warning("移除成功");
  }

  // 批量设置
  batchSetting(command: string) {
    if (!this.selection.length) {
      const text = command === "commission" ? "分销佣金" : "分销期限";
      this.$message.warning(`请先选择要设置${text}的商品`);
      return;
    }
    // 设置佣金
    if (command === "commission") {
      const display = ["promotion_commission_amount"];
      if (this.USER.version !== 3) {
        display.push("service_commission_amount");
      }
      this.batchSettingDisplay = [...display];
      this.batchSettingVisible = true;
      return;
    }
    // 设置分销期限
    // 仅可针对当日及之后未设置过分销期限的商品批量设置分销期限
    const current = new Date().setHours(23, 59, 59, 0) / 1000;
    if (
      this.selection.some(item =>
        item.disable_date_range.some(
          (date: any) => date.distribution_end > current
        )
      )
    ) {
      this.$message.warning(
        "仅可针对未设置过期限的商品执行此操作，请重新勾选商品再试"
      );
      return;
    }
    this.batchSettingDisplay = ["date_range"];
    this.batchSettingVisible = true;
  }

  handleBatchSettingDone(data: any) {
    this.chooseGoods = this.chooseGoods.map(item => {
      if (this.selection.find(val => val.sku_id === item.sku_id)) {
        return {
          ...item,
          ...data
        };
      }
      return item;
    });
    this.search();
    this.closeBatchSetting();
  }

  // 关闭批量设置弹窗
  private closeBatchSetting() {
    this.batchSettingVisible = false;
  }

  // 下一步
  async next() {
    if (!this.selection.length) {
      this.$message.warning("请选择商品");
      return;
    }
    this.active = 1;
    const list = [...this.selection];
    const params = {
      sku_ids: list.map(item => item.sku_id)
    };
    this.isLoaded = false;
    try {
      let { data } = await DistributionApi.repeatInspect(params);
      if (!Array.isArray(data)) {
        data = [];
      }
      console.log(this.chooseGoods);
      this.chooseGoods = list.map(item => {
        const temp = data.find((val: any) => val.sku_id === item.sku_id) || {
          date_range: []
        };
        // 查找点击上一步后，再次进入时原有数据
        const prev_edit_item = this.chooseGoods.find(
          val => val.sku_id === item.sku_id
        ) || {
          date_range: [],
          promotion_commission_amount: 0,
          service_commission_amount: 0
        };
        return {
          ...item,
          date_range: prev_edit_item.date_range,
          promotion_commission_amount:
            prev_edit_item.promotion_commission_amount,
          service_commission_amount: prev_edit_item.service_commission_amount,
          disable_date_range: temp.date_range,
          pickerOption: {
            disabledDate: handleDisabledDate(temp.date_range),
            onPick({ minDate }: { minDate: any }) {
              this.disabledDate = handleDisabledDate(temp.date_range, minDate);
            }
          },
          handleDateChange() {
            this.pickerOption.disabledDate = handleDisabledDate(
              temp.date_range
            );
          }
        };
      });
    } catch (error) {
      console.log(error);
    }
    this.isLoaded = true;
    // 清理之前的选中与筛选条件
    this.searchType = 3;
    this.searchValue = "";
    (this.$refs.table as any)?.clearSelection();
    this.setSearchTypeOptions();
    this.search();
  }

  // 取消
  cancel() {
    this.$router.go(-1);
  }

  // 同步数据，将 list 中填写的数据同步到chooseGoods
  synchronousData() {
    // 1.判断两者长度是否相等
    if (this.chooseGoods.length === this.list.length) {
      this.chooseGoods = [...this.list];
      return;
    }
    this.chooseGoods = this.chooseGoods.map(item => {
      const temp = this.list.find(val => item.sku_id === val.sku_id);
      return temp || item;
    });
  }

  scrollToErrorRow(i: number) {
    this.$nextTick(() => {
      const el = document.getElementsByClassName("el-main")[0];
      const sTop = el.scrollTop;
      const to = 260 + 80 * i;
      scrollTop(el, sTop, to, 100);
    });
  }

  // 将错误数据行滚动到可视范围
  showCheckErrorRow(index: number) {
    if (this.chooseGoods.length === this.list.length) {
      // 滚动
      this.scrollToErrorRow(index);
      return;
    }
    this.reset();
    // 滚动
    this.scrollToErrorRow(index);
  }

  // 校验数据
  checkData() {
    for (let i = 0; i < this.chooseGoods.length; i++) {
      const item = this.chooseGoods[i];
      if (
        item.promotion_commission_amount === undefined ||
        item.promotion_commission_amount === ""
      ) {
        this.$message.warning("请输入推广佣金比例");
        this.showCheckErrorRow(i);
        return false;
      }
      if (
        item.promotion_commission_amount > 100 ||
        item.promotion_commission_amount < 0
      ) {
        this.$message.warning("请输入合理的推广佣金比例");
        this.showCheckErrorRow(i);
        return false;
      }
      if (
        this.USER.version === 3 &&
        (item.service_commission_amount === undefined ||
          item.service_commission_amount === "")
      ) {
        this.$message.warning("请输入服务佣金比例");
        this.showCheckErrorRow(i);
        return false;
      }
      if (
        item.promotion_commission_amount > 100 ||
        item.promotion_commission_amount < 0
      ) {
        this.$message.warning("请输入合理的服务佣金比例");
        this.showCheckErrorRow(i);
        return false;
      }
      if (!(item.date_range && item.date_range.length === 2)) {
        this.$message.warning("请输入分销期限 ");
        this.showCheckErrorRow(i);
        return false;
      }
    }
    return true;
  }

  // 上一步
  prev() {
    this.active = 0;
    this.synchronousData();
    (this.$refs.table as any).clearSelection();
    this.setSearchTypeOptions();
    this.chooseGoods.forEach(row => {
      (this.$refs.table as any).toggleRowSelection(row, true);
    });
    this.loadGoods();
  }

  // 保存
  async save() {
    this.synchronousData();
    if (!this.checkData()) {
      return;
    }
    const params = this.chooseGoods.map(item => {
      const {
        spu_id,
        sku_id,
        promotion_commission_amount,
        service_commission_amount,
        date_range
      } = item;
      return {
        spu_id,
        sku_id,
        promotion_commission_amount: Math.round(
          promotion_commission_amount * 100
        ),
        service_commission_amount: Math.round(service_commission_amount * 100),
        distribution_start: date_range[0] / 1e3,
        distribution_end: date_range[1] / 1e3 + 86399
      };
    });
    try {
      const res: any = await DistributionApi.batchSaveDistributionGoods({
        params
      });
      if (res.code === 200) {
        this.$message.success("保存成功");
        this.cancel();
      }
    } catch (error) {
      console.log(error);
    }
  }

  /*** 生命周期 ***/
  private created() {
    this.setSearchTypeOptions();
  }
  private mounted() {
    sessionStorage.setItem("staff-setting-tab-active", "DistributionGoods");
    this.loadGoods();
  }
}
</script>

<style lang="scss" scoped>
.steps-wrapper {
  margin-bottom: 24px;
}

/deep/.layout-list {
  .slide-box {
    display: none;
  }
}

.search-item {
  display: inline-flex;
  align-items: center;
  margin-bottom: 0;
  margin-right: 0;
}

.input-with-select {
  width: 110px;
}

.search-input {
  width: 280px !important;
  margin-bottom: 0px !important;
}

.el-select-dropdown__item.selected {
  font-weight: normal !important;
  height: auto !important;
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
    color: #909399;
    font-size: 12px;
  }
}

.table-remove-btn {
  color: #e34b49;
}

.table-cancel-btn {
  color: #909399;
}

/deep/.el-range-editor--small {
  width: 280px !important;
}

/deep/.el-step__head.is-success {
  border-color: $color-main;
  .el-step__icon.is-text {
    color: $color-main;
  }
}
</style>
