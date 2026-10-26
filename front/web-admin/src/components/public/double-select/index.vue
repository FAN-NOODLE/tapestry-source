<template>
  <div
    class="double-select"
    :class="{
      'db-select-small': widthSize === 'small',
      'db-select-medium': widthSize === 'medium'
    }"
  >
    <!-- 位于主下拉框左侧的附属下拉框 -->
    <el-select
      class="short-select short-select-left"
      v-if="type === 'left'"
      v-model="short"
      :size="size"
      @change="selectChange"
    >
      <el-option
        v-for="item in secondaryList"
        :key="item.key"
        :label="item.name"
        :value="item.key"
      ></el-option>
    </el-select>
    <!-- 主下拉框 -->
    <el-select
      class="main-select"
      :class="{
        'main-type-left': type === 'left',
        'main-type-right': type === 'right'
      }"
      v-model="main"
      ref="select"
      :multiple="multiple"
      :disabled="disabled"
      :size="size"
      :clearable="clearable"
      :collapse-tags="collapseTags"
      :multiple-limit="multipleLimit"
      :name="name"
      :autocomplete="autocomplete"
      :placeholder="placeholder"
      :filterable="filterable"
      :allow-create="allowCreate"
      :filter-method="filterMethod"
      :loading="loading"
      :loading-text="loadingText"
      :no-match-text="noMatchText"
      :no-data-text="noDataText"
      :popper-class="popperClass"
      :reserve-keyword="reserveKeyword"
      :default-first-option="defaultFirstOption"
      :popper-append-to-body="popperAppendToBody"
      :automatic-dropdown="automaticDropdown"
      @change="onChangeSelect"
      @visible-change="onVisibleChangeSelect"
      @remove-tag="onRemoveTagSelect"
      @clear="onClearSelect"
      @blur="onBlurSelect"
      @focus="onFocusSelect"
    >
      <!-- 内置一个符合UI规范的筛选功能 -->
      <el-input
        v-if="searchInput"
        v-model="keywords"
        ref="select-search"
        class="select-search"
        :class="{
          'cy-search-small': widthSize === 'small',
          'cy-search-medium': widthSize === 'medium'
        }"
        placeholder="请输入"
        prefix-icon="el-icon-search"
        clearable
        :validate-event="false"
        @input="dropDownSearch"
      ></el-input>
      <template v-for="(item, index) in optionsList">
        <el-option
          v-if="
            mainList.some(
              option => option[options.value] === item[options.value]
            )
          "
          :value="item[options.value]"
          :label="item[options.label]"
          :key="index"
        ></el-option>
      </template>
      <template slot="prefix">
        <slot name="prefix"></slot>
      </template>
      <template slot="empty">
        <div>
          <el-input
            v-show="searchInput"
            v-model="keywords"
            ref="default-search"
            class="select-search"
            :class="{
              'cy-search-small': widthSize === 'small',
              'cy-search-medium': widthSize === 'medium'
            }"
            placeholder="请输入"
            prefix-icon="el-icon-search"
            clearable
            :validate-event="false"
            @input="dropDownSearch"
          ></el-input>
          <div class="select-empty">
            {{ noMatchText }}
          </div>
        </div>
        <slot name="empty"></slot>
      </template>
    </el-select>
    <!-- 位于主下拉框右侧的附属下拉框 -->
    <el-select
      class="short-select short-select-right"
      v-if="type === 'right'"
      v-model="short"
      :size="size"
      @change="selectChange"
    >
      <el-option
        v-for="item in secondaryList"
        :key="item.key"
        :label="item.name"
        :value="item.key"
      ></el-option>
    </el-select>
  </div>
</template>
<script lang="ts">
import {
  IDoubleType,
  ISize,
  IWidthSize,
  IOptionsItem,
  IListItem
} from "./type";
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
@Component({})
export default class DoubleSelect extends Vue {
  /**
   * @param {String} left
   * @description 附属下拉框位于主下拉框的左侧还是右侧 可选 left | right
   * */

  @Prop({ default: "left", type: String }) type?: IDoubleType;
  @Prop({ default: "small" }) widthSize?: IWidthSize; //非element-ui参数，由于UI定宽，仅有small宽度为280px和medium宽度为400px
  @Prop() value?: string | number | boolean; // 绑定值
  @Prop({ default: false, type: Boolean }) multiple!: boolean; //是否多选
  @Prop({ default: false, type: Boolean }) disabled!: boolean; //是否禁用
  @Prop({ default: "small" }) size!: ISize; //大小
  @Prop({ default: false, type: Boolean }) clearable!: boolean; //是否可以清空选项
  @Prop({ default: false, type: Boolean }) collapseTags!: boolean; //多选时是否将选中值按文字的形式展示
  @Prop({ default: 0 }) multipleLimit?: number; // 多选时用户最多可以选择的项目数，为 0 则不限制
  @Prop() name?: string; // select input 的 name 属性
  @Prop({ default: "off" }) autocomplete?: string; //select input 的 autocomplete 属性
  @Prop({ default: "请选择" }) placeholder?: string; //占位符
  @Prop({ default: false, type: Boolean }) filterable!: boolean; //是否可搜索
  @Prop({ default: false, type: Boolean }) allowCreate!: boolean; //是否允许用户创建新条目，需配合 filterable 使用
  @Prop({ default: Function }) filterMethod?: Function; //远程搜索方法
  @Prop({ default: false, type: Boolean }) loading!: boolean; //是否正在从远程获取数据
  @Prop({ default: "加载中" }) loadingText?: string; //远程加载时显示的文字
  @Prop({ default: "无匹配数据" }) noMatchText?: string; //搜索条件无匹配时显示的文字，也可以使用slot="empty"设置
  @Prop({ default: "无数据" }) noDataText?: string; //选项为空时显示的文字，也可以使用slot="empty"设置
  @Prop() popperClass?: string; //Select 下拉框的类名
  @Prop({ default: false, type: Boolean }) reserveKeyword!: boolean; //多选且可搜索时，是否在选中一个选项后保留当前的搜索关键词
  @Prop({ default: false, type: Boolean }) defaultFirstOption!: boolean; //在输入框按下回车，选择第一个匹配项。需配合 filterable 或 remote 使用
  @Prop({ default: true, type: Boolean }) popperAppendToBody!: boolean; //是否将弹出框插入至 body 元素。在弹出框的定位出现问题时，可将该属性设置为 false
  @Prop({ default: false, type: Boolean }) automaticDropdown!: boolean; //对于不可搜索的 Select，是否在输入框获得焦点后自动弹出选项菜单
  @Prop({ default: false, type: Boolean }) searchInput!: boolean; //开启输入框搜索
  @Prop({ default: () => [], type: Array }) private optionsList!: Array<any>;
  @Prop({
    default: () => {
      return { label: "label", value: "value" };
    }
  })
  options?: IOptionsItem;
  @Prop({ default: () => [], type: Array }) secondaryList!: Array<IListItem>;

  @Watch("value")
  onChangeValue(val: any) {
    this.main = val;
  }
  @Watch("main")
  onChangeMain(val: any) {
    this.$emit("input", val);
  }
  @Watch("optionsList")
  onchangeOptions() {
    this.slotUpdate();
  }
  @Watch("mainList")
  onchangeList(val: any) {
    /**
     * 保持input输入框有焦点
     * */
    if (val && this.searchInput) {
      if (this.mainList.length === 0) {
        const timer = setTimeout(() => {
          (this.$refs["default-search"] as any).focus();
          clearTimeout(timer);
        });
      } else {
        const timer = setTimeout(() => {
          (this.$refs["select-search"] as any).focus();
          clearTimeout(timer);
        });
      }
    }
  }
  private main = ""; //主下拉框的值 (长的)
  private short = this.secondaryList[0].key; //附属下拉框的值 (短的)
  private mainList: Array<any> = []; //主下拉框的选项list
  private keywords = ""; //筛选条件关键字
  /**
   * @name 方法
   * */
  private slotUpdate() {
    try {
      this.mainList = [...this.optionsList];
    } catch (error) {
      this.mainList = [];
    }
  }
  // 下拉搜索匹配option
  private dropDownSearch() {
    this.mainList = this.optionsList.filter(item =>
      item[(this.options as IOptionsItem).label].includes(this.keywords)
    );
  }
  // 选中值发生变化时触发
  private onChangeSelect(val: any) {
    const timeOut = setTimeout(() => {
      this.$emit("change", val);
      clearTimeout(timeOut);
    });
  }
  //下拉框出现/隐藏时触发
  private onVisibleChangeSelect(val: boolean) {
    if (this.searchInput && !val) {
      this.keywords = "";
      this.mainList = [...this.optionsList];
    }
    this.$emit("visible-change", val);
  }
  //多选模式下移除tag时触发
  private onRemoveTagSelect(val: any) {
    const timeOut = setTimeout(() => {
      this.$emit("remove-tag", val);
      clearTimeout(timeOut);
    });
  }

  //可清空的单选模式下用户点击清空按钮时触发
  private onClearSelect() {
    const timeOut = setTimeout(() => {
      this.$emit("clear");
      clearTimeout(timeOut);
    });
  }
  // 当失去焦点时触发
  private onBlurSelect(e: MouseEvent) {
    const timeOut = setTimeout(() => {
      this.$emit("blur", e);
      clearTimeout(timeOut);
    });
  }
  // 当获得焦点时触发
  private onFocusSelect(e: MouseEvent) {
    const timeOut = setTimeout(() => {
      this.$emit("focus", e);
      clearTimeout(timeOut);
    });
  }
  // 附属下拉框的change事件
  private selectChange() {
    this.$emit("secondaryChange", this.short);
  }
  // 将element自带的focus方法暴漏出去
  public focus() {
    (this.$refs["select"] as any).focus();
  }
  // 将element自带的blur方法暴漏出去
  public blur() {
    (this.$refs["select"] as any).blur();
  }
  public reset() {
    this.short = this.secondaryList[0].key;
  }
}
</script>
<style lang="scss" scoped>
.double-select {
  display: inline-block;
  .main-type-left {
    ::v-deep {
      .el-input__inner {
        border-radius: 0 4px 4px 0;
      }
    }
  }
  .main-type-right {
    ::v-deep {
      .el-input__inner {
        border-radius: 4px 0 0 4px;
      }
    }
  }
  .short-select {
    width: 82px !important;
    ::v-deep {
      .el-input__inner {
        padding-right: 6px !important;
      }
    }
  }
  .short-select-left {
    ::v-deep {
      .el-input__inner {
        border-radius: 4px 0 0 4px;
        border-right: none;
      }
    }
  }
  .short-select-right {
    ::v-deep {
      .el-input__inner {
        border-radius: 0 4px 4px 0;
        border-left: none;
      }
    }
  }
  ::v-deep.el-select {
    .el-input--mini {
      height: 32px !important;
      line-height: 32px !important;
      .el-input__inner {
        @extend .el-input--mini;
      }
    }
    .el-input--small {
      height: 36px !important;
      line-height: 36px !important;
      .el-input__inner {
        @extend .el-input--small;
      }
    }
    .el-input--medium {
      height: 36px !important;
      line-height: 36px !important;
      .el-input__inner {
        @extend .el-input--medium;
      }
    }
  }
}
/**select缺省样式**/
.select-empty {
  padding: 10px 0;
  margin: 0;
  text-align: center;
  color: #c0c4cc;
  font-size: 14px;
}
/**兼容两种长度的样式**/
.db-select-small {
  .main-select {
    ::v-deep {
      .el-input__inner {
        width: 219px;
      }
    }
  }
}
.db-select-medium {
  .main-select {
    ::v-deep {
      .el-input__inner {
        width: 289px;
      }
    }
  }
}
.select-search {
  margin-top: 6px;
  margin-bottom: 7px;
  margin-left: 20px;
  height: 34px !important;
  line-height: 34px !important;
  ::v-deepinput {
    height: 34px !important;
    line-height: 34px !important;
  }
}
.cy-search-small {
  width: 178px;
}
.cy-search-medium {
  width: 248px;
}
</style>
