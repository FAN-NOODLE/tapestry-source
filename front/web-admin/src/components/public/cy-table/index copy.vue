<template>
  <div>
    <el-table
      @cell-mouse-enter="isShowPopover"
      @cell-mouse-leave="handleMouseLeave"
      @sort-change="handleSortChange"
      @selection-change="handleSelectionChange"
      :data="tableConfig.data"
      :height="base.height"
      :maxHeight="base.maxHeight"
      style="width: auto"
      :row-key="base.rowkey"
      stripe
    >
      <el-table-column type="selection" width="55" v-if="base.checkbox.isShow">
      </el-table-column>
      <el-table-column label="" width="65" fixed v-if="radio.isShow">
        <template slot-scope="scope">
          <el-radio
            :label="scope.$index"
            v-model="radio.radioModel"
            @change.native="handleRadioChange(scope.row)"
          >
            {{ radio.radioText }}
          </el-radio>
        </template>
      </el-table-column>
      <el-table-column
        v-for="columnItem in columns"
        :fixed="columnItem.fixed"
        :key="columnItem.key"
        :prop="columnItem.key"
        :label="columnItem.title"
        :width="columnItem.width"
        :align="columnItem.align"
        :sortable="columnItem.sortable"
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
            v-if="columnItem.type === 'slot'"
            :row="scope.row"
          ></slot>
          <!--图片-->
          <el-image
            style="width: 50px; height: 50px"
            :src="scope.row[columnItem.key]"
            fit="fit"
            v-else-if="columnItem.type === 'image'"
          ></el-image>
          <!--文本-->
          <span v-else>
            {{ scope.row[columnItem.key] }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="排序" width="80" v-if="base.isRowSort">
        <template slot-scope="scope">
          <i
            class="el-icon-sort-down color-main"
            :class="{ 'color-warn': scope.$index === 0 }"
            aria-hidden="true"
          ></i>
          <i
            class="el-icon-sort-up color-main"
            :class="{
              'color-warn': scope.$index === tableConfig.data.length - 1
            }"
            aria-hidden="true"
          ></i>
        </template>
      </el-table-column>
      <el-table-column
        fixed="right"
        :label="operate.title || '操作'"
        :width="operate.width"
        v-if="operate.key.length > 0 || operate.more.length > 0"
      >
        <template slot-scope="scope">
          <el-button
            :class="keyItem.class"
            v-for="keyItem in operate.key"
            type="text"
            :key="keyItem.id"
            @click="handleOperateKeyClick(scope.row, keyItem.key)"
            >{{ keyItem.text }}</el-button
          >
          <el-dropdown class="m-l-10" v-if="operate.more.length > 0">
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
      <div solt="empty">
        暂无数据
      </div>
    </el-table>
    <!--分页-->
    <el-pagination
      class="pull-right"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="this.pagination.currentPage"
      :page-sizes="this.pagination.pageSizes"
      :page-size="this.pagination.pageSize"
      layout="total, sizes, prev, pager, next, jumper"
      :total="this.pagination.total"
    >
    </el-pagination>
    <div
      role="tooltip"
      id="el-popover-444"
      aria-hidden="true"
      class="el-popover el-popper font-12 el-popover--plain"
      tabindex="0"
      style="width: 200px; transform-origin: center top; z-index: 2008;"
      :style="{ left: offsetLeft + 'px', top: offsetTop + 'px' }"
      v-if="isTip"
      x-placement="bottom"
    >
      {{ msg }}
      <div x-arrow class="popper__arrow" style="left: 47.6562px;"></div>
    </div>
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
@Component
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
    radioText: ""
  };

  private msg = "";
  private isTip = false;
  private offsetLeft = 0;
  private offsetTop = 0;
  // 判断字体图标
  private isSupportEl(iconName: string) {
    return iconName.indexOf("el-") !== -1;
  }
  // 列表默认宽度
  private cDefaultW = 0;

  // 列表宽度计算
  private cWCount() {
    // let noWCount = 0; // 未定义的宽度计数
    // let sumW = 0; // 以定义宽度总和
    // let defaultW = 0; // 平均默认宽度
    // const columns = this.tableConfig.columns;
    // for (let i = 0; i < columns.length; i++) {
    //   if (columns["width"]) {
    //     sumW =
    //       sumW + (parseInt(columns.width.split("%").join("")) / 100) * 1440;
    //   } else {
    //     noWCount++;
    //   }
    // }
    // if (this.tableConfig.operate["width"]) {
    //   // 需要拿到操作的宽度
    //   sumW =
    //     sumW +
    //     (parseInt(this.tableConfig.operate.width.split("%").join("")) / 100) *
    //       1440;
    // } else {
    //   noWCount++;
    // }
    // if (noWCount !== 0) {
    //   defaultW = 1440 / noWCount;
    // }
    // return defaultW;
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
        rowkey: "",
        isShow: false
      };
    }
    // 处理列表多选功能
    if (this.tableConfig["checkbox"]) {
      this.base.checkbox = this.tableConfig.checkbox;
    } else {
      this.base.checkbox = {
        rowkey: "",
        isShow: false
      };
    }
    if (this.base.tree.isTree) {
      this.base.rowkey = this.base.tree.rowkey;
    } else {
      if (this.base.checkbox.isShow) {
        this.base.rowkey = this.base.checkbox.rowkey;
      } else {
        this.base.rowkey = "";
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
      currentPage: this.tableConfig.pagination["currentPage"] || 1
    };
    console.log(this.pagination);
  }
  // 处理-列
  private handleC() {
    console.log(JSON.stringify(this.tableConfig));
    this.columns = this.tableConfig.columns.map((item: any) => {
      // 列固定列处理
      item.fixed = item["fixed"] || false;
      // 列宽度处理
      if (item["width"]) {
        item.width = (parseInt(item.width.split("%").join("")) / 100) * 1440;
      } else {
        item.width = this.cDefaultW;
      }
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
    // 宽度处理
    if (this.tableConfig.operate["width"]) {
      this.operate.width =
        (parseInt(this.tableConfig.operate.width.split("%").join("")) / 100) *
        1440;
    } else {
      this.operate.width = this.cDefaultW;
    }
  }
  // 生命周期
  private created() {
    // console.log("created", JSON.stringify(this.tableConfig));
    // this.cDefaultW = this.cWCount(); // 列表默认宽度复制
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
  // 当单元格 hover 进入时会触发该事件
  private isShowPopover(row: any, column: any, cell: any, event: any) {
    const range = document.createRange();
    const cellChild = event.target.querySelector(".cell");
    range.setStart(cellChild, 0);
    range.setEnd(cellChild, cellChild.childNodes.length);
    const rangeWidth = range.getBoundingClientRect().width;
    if (
      rangeWidth > cellChild.offsetWidth ||
      cellChild.scrollWidth > cellChild.offsetWidth
    ) {
      this.msg = row[column.property];
      this.isTip = true;
      this.offsetLeft = cell.offsetLeft + cell.offsetWidth / 2;
      this.offsetTop = event.clientY - 21;
      console.log("111");
    } else {
      return false;
    }
  }
  private handleMouseLeave() {
    this.msg = "";
    this.isTip = false;
    this.offsetLeft = 0;
    this.offsetTop = 0;
  }

  /************************排序事件************************/
  // 点击列排序
  private handleSortChange(val: object) {
    this.$emit("columnSort", val);
  }
  /************************操作事件************************/
  private handleOperateKeyClick(row: object, key: string) {
    this.$emit("operateKeyClick", row, key);
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
    this.$emit("radioChange", row);
  }

  /************************多选事件************************/
  private handleSelectionChange(row: any) {
    this.$emit("selectionChange", row);
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
</style>
