<template>
  <div>
    <el-table
      @sort-change="handleSortChange"
      @selection-change="handleSelectionChange"
      :data="tableConfig.data"
      :height="base.height"
      :maxHeight="base.maxHeight"
      style="width: auto"
      :cell-style="{ padding: base.tbPadding + 'px 0px' }"
      :header-cell-style="{ padding: '0px 0px' }"
      :row-key="base.rowKey || ''"
      stripe
      :default-sort="tableConfig.defaultSort"
    >
      <!--多选列表展示-->
      <el-table-column type="selection" width="55" v-if="base.checkbox.isShow">
      </el-table-column>
      <!--单选列表展示-->
      <el-table-column
        label=""
        width="54"
        class-name="radio"
        v-if="radio.isShow"
      >
        <template slot-scope="scope">
          <el-radio
            :label="scope.$index"
            v-model="radio.radioModel"
            @change.native="handleRadioChange(scope.row)"
            >{{ radio.radioText }}</el-radio
          >
        </template>
      </el-table-column>
      <!--基本列展示-->
      <el-table-column
        v-for="(columnItem, index) in columns"
        :key="index"
        :prop="columnItem.key"
        :label="columnItem.title"
        :fixed="columnItem.fixed"
        :width="columnItem.width"
        :sortable="columnItem.sortable"
        :type="columnItem.expand"
        :align="columnItem.align"
      >
        <template slot="header">
          <!--文本-->
          <span>{{ columnItem.title }}</span>
          <span class="column-icon" v-if="columnItem.isIcon">
            <i
              class="el-icon-question"
              :class="columnItem.iconName"
              aria-hidden="true"
              v-if="isSupportEl(columnItem.iconName)"
            ></i>
            <svg-icon :icon-class="columnItem.iconName" v-else></svg-icon>
          </span>
        </template>
        <template slot-scope="scope">
          <!--标签-->
          <slot
            :name="columnItem.key"
            v-if="columnItem.type === 'slot' || columnItem.expand === 'expand'"
            :row="scope.row"
          ></slot>
          <!--图片-->
          <el-image fit="cover"
            style="width: 50px; height: 50px"
            :src="scope.row[columnItem.key]"
            
            v-else-if="columnItem.type === 'image'"
          ></el-image>
          <!--文本-->
          <span v-else>
            <el-popover
              v-if="scope.row[columnItem.key].length > columnItem.splitLength"
              placement="bottom-end"
              trigger="hover"
              popper-class="cy-t-popper"
            >
              <span>{{ scope.row[columnItem.key] }}</span>
              <span slot="reference" style="cursor:pointer">{{
                scope.row[columnItem.key].substring(0, columnItem.splitLength) +
                  "..."
              }}</span>
            </el-popover>
            <span v-else>{{ scope.row[columnItem.key] }}</span>
          </span>
        </template>
      </el-table-column>
      <!--行排序-->
      <el-table-column label="排序" width="84" v-if="base.isRowSort">
        <template slot-scope="scope">
          <i
            class="el-icon-sort-down color-main pull-left"
            :class="{ 'color-second-line': scope.$index === 0 }"
            aria-hidden="true"
          ></i>
          <i class="line pull-left" aria-hidden="true"></i>
          <i
            class="el-icon-sort-up color-main pull-left"
            :class="{
              'color-second-line': scope.$index === tableConfig.data.length - 1,
            }"
            aria-hidden="true"
          ></i>
        </template>
      </el-table-column>
      <!--操作处理-->
      <el-table-column
        fixed="right"
        :label="operate.title || '操作'"
        :width="operate.width"
        v-if="operate.key.length > 0 || operate.more.length > 0"
      >
        <template slot-scope="scope">
          <el-button
            :class="keyItem.class"
            v-for="(keyItem, index) in operate.key"
            type="text"
            :key="index"
            @click="handleOperateKeyClick(scope.row, keyItem.key)"
            >{{ keyItem.text }}</el-button
          >
          <el-dropdown class="m-l-20" v-if="operate.more.length > 0">
            <el-button type="text"
              >更多<i
                class="el-icon-arrow-down el-icon--right"
                aria-hidden="true"
              ></i
            ></el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item
                v-for="moreItem in operate.more"
                :key="moreItem.key"
                >{{ moreItem.text }}</el-dropdown-item
              >
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
      <div slot="empty" class="table-empty">
        <img alt="" v-if="tableConfig.empty.img" :src="tableConfig.empty.img" />
        <p v-if="tableConfig.empty.text">{{ tableConfig.empty.text }}</p>
      </div>
    </el-table>
    <!--分页-->
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
/** 封装table组件
 * 1. 相关缩写
 * w/W : width
 * c/C : columns
 * d/D : data
 * o/O : operate
 *
 *
 * 2. 基本支持
 * * *1. 标题动态 title ok
 * * *2. 数据动态 key ok
 * * *4. 操作动态 operate ok
 * * *4. 列表宽度动态 width ok
 * * *5. 表格颜色间隔 默认 ok
 * * *6. 分页 pagination ok

 * 3. 关键支持
 * * *1. 单选 ok
 * * *2. 多选 ok
 * * *3. 更新排序 ok
 * * *4. 固定列 ok

 * 4. 其他支持
 * * *1. 树形数据 ok
 * * *2. 支持滚动 ok
 * * *3. 排序 ok
 * * *4. 支持列内容自定义
 **/
@Component({})
export default class CyTable extends Vue {
  @Prop() tableConfig!: any;
  // @Prop() columns!: Array<object>;
  // @Prop() data!: Array<object>;
  // @Prop() operate!: object;
  private base: any = {};
  private columns: Array<any> = [];
  private operate: any = {};
  private pagination: any = {};
  private radio: any = {
    isRadio: false,
    radioModel: "",
    radioText: "",
  };

  // 判断字体图标
  private isSupportEl(iconName: string) {
    return iconName.indexOf("el-") !== -1;
  }

  // 基础处理
  private handleBase() {
    // 列表滚动
    this.base.maxHeight = this.tableConfig["maxHeight"] || "auto";
    // 当数据为空的情况
    if (this.tableConfig.data.length === 0) {
      this.base.height = 500;
    } else {
      this.base.height = "auto";
    }
    /**
     * 1. 判断 是否启用树形列表
     * 2. 判断 是否启用多选功能
     * 备注：树形结构优先
     */
    if (this.tableConfig["tree"]) {
      this.base.tree = this.tableConfig.tree;
    } else {
      this.base.tree = {
        rowKey: "",
        isShow: false,
      };
    }
    // 处理列表多选功能
    if (this.tableConfig["checkbox"]) {
      this.base.checkbox = this.tableConfig.checkbox;
    } else {
      this.base.checkbox = {
        rowKey: "",
        isShow: false,
      };
    }
    if (this.base.tree.isTree) {
      this.base.rowKey = this.base.tree.rowKey;
    } else {
      if (this.base.checkbox.isShow) {
        this.base.rowKey = this.base.checkbox.rowKey;
      } else {
        this.base.rowKey = "";
      }
    }
    /**
     * 单选
     */
    if (this.tableConfig["radio"]) {
      this.radio.radioModel = this.tableConfig.radio.radioModel;
      this.radio.isShow = this.tableConfig.radio.isShow || false;
    }

    // 行更新排序
    this.base.isRowSort = this.tableConfig["isRowSort"] || false;
    // 分页
    this.pagination = {
      pageSizes: this.tableConfig.pagination["pageSizes"] || [10, 50, 100, 200],
      pageSize: this.tableConfig.pagination["pageSizes"]
        ? this.tableConfig.pagination.pageSizes[0]
        : 10,
      total: this.tableConfig.pagination.total,
      currentPage: this.tableConfig.pagination["currentPage"] || 1,
    };
    console.log(this.pagination);
  }
  // 处理-列
  private handleC() {
    console.log(JSON.stringify(this.tableConfig));
    this.columns = this.tableConfig.columns.map((item: any) => {
      // 列固定列处理
      item.fixed = item["fixed"] || false;
      if (item["type"] === "image") {
        this.base.tbPadding = 10;
      }
      // 列固定列处理
      item.splitLength = item["splitLength"] || 1000;
      return item;
    });
  }
  // 处理-操作
  private handleO() {
    console.log(JSON.stringify(this.tableConfig));
    const operate = this.tableConfig.operate;
    this.operate = operate;
    // 标题处理
    this.operate.title = operate["title"] || "操作";
  }
  // 生命周期
  private created() {
    this.handleC(); // 列处理
    this.handleO(); // 操作处理
    this.handleBase(); // 基本处理
  }
  private mounted() {
    console.log("mounted", this.tableConfig);
  }
  private destroyed() {
    console.log("destroyed");
  }

  /************************排序事件************************/
  // 点击列排序
  private handleSortChange(val: object) {
    this.tableConfig.columnSort(val);
  }
  /************************操作事件************************/
  private handleOperateKeyClick(row: object, key: string) {
    this.tableConfig.operate.operateKeyClick(row, key);
  }
  /************************分页事件************************/
  private handleSizeChange(val: string) {
    // this.pagination.pageSize = val;
    this.$emit("sizeChange", val);
  }
  private handleCurrentChange(val: string) {
    this.$emit("currentChange", val);
  }

  /************************单选事件************************/
  private handleRadioChange(row: object) {
    this.tableConfig.radio.radioChange(row);
  }

  /************************多选事件************************/
  private handleSelectionChange(row: any) {
    // this.$emit("selectionChange", row);
    this.tableConfig.checkbox.selectionChange(row);
  }
}
</script>
<style lang="scss" scoped>
// 公共颜色

.column-icon {
  margin-left: 5px;
}
.el-button {
  font-weight: normal;
}
.el-button + .el-button {
  margin-left: 20px;
}
.el-popover--plain {
  padding: 12px 10px !important;
}
.el-image {
  display: block;
}
.line {
  position: relative;
  width: 14px;
  height: 14px;
  display: inline-block;
  &:after {
    content: "";
    position: absolute;
    left: 50%;
    width: 1px;
    background: $color-second-line;
    height: 14px;
  }
}
/deep/ {
  .el-table__row {
    .cell {
      .tag {
        margin-right: 10px;
        margin-bottom: 10px;
      }
      .tag + .tag {
        margin-left: 0px !important;
      }
    }
  }
}
.table-empty {
  text-align: center;
  img {
    width: 100%;
  }
}
</style>
