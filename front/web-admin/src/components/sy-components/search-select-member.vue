<template>
  <div class="search-select-member">
    <el-select
      class="select-type"
      v-if="hasPrepend"
      size="medium"
      v-model="memberType"
      @change="onMemberTypeChange"
    >
      <el-option
        v-for="item in memberTypeList"
        :key="item.key"
        :label="item.name"
        :value="item.key"
      ></el-option>
    </el-select>
    <el-select
      class="select-collapse"
      size="medium"
      multiple
      collapse-tags
      :placeholder="placeholder"
      v-model="selectedIds"
      @visible-change="onVisibleChange"
      @remove-tag="removeMember"
    >
      <el-option
        v-for="item in selectedList"
        :key="item.id"
        :label="item.name"
        :value="item.id"
      ></el-option>
    </el-select>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Watch, Emit } from "vue-property-decorator";

@Component
export default class SearchSelectMember extends Vue {
  @Prop({ default: false }) hasPrepend?: boolean; // 是否有前缀下拉搜索
  @Prop({ default: () => [] }) memberTypeList?: Array<any>; // 下拉选项
  @Prop({ default: "请选择导购" }) placeholder?: string;
  @Prop({ default: () => [] }) selectedList?: Array<any>; // 已选的数据

  private memberType: string | number = "";
  private selectedIds: Array<any> = [];

  @Watch("selectedList")
  onChangeSelectedList(list: Array<any>) {
    this.selectedIds = list.map((v) => v.id);
  }

  @Emit("member-type-change")
  onMemberTypeChange(val: string | number) {
    return val;
  }

  // 初始化 memberType
  public initMemberType() {
    if (this.memberTypeList && this.memberTypeList.length) {
      this.memberType = this.memberTypeList[0].key;
    }
  }

  // 隐藏选择导购下拉选项
  @Emit("visible-change")
  private onVisibleChange(is_show: boolean) {
    if (is_show) {
      const dom = document.getElementsByClassName("is-multiple")[0];
      (dom as any).style.opacity = 0;
    }
  }

  // 删除已选导购
  @Emit("remove-member")
  private removeMember(val: any) {
    return val;
  }

  private mounted() {
    this.initMemberType();
  }
}
</script>

<style lang="scss" scoped>
.search-select-member {
  width: 280px;
  display: flex;
  align-items: center;
  overflow: hidden;
  border: 1px solid $color-second-line;
  background-color: $color-white;
  border-radius: 0;
  /deep/.el-tabs__nav {
    /deep/.el-tabs__active-bar {
      width: 4px !important;
    }
  }
  .select-type {
    width: 110px !important;
    box-sizing: border-box;
    border-right: 1px solid $color-second-line;
  }
  /deep/.el-input__inner {
    border: none;
    border-radius: 0;
    color: $color-text-second;
  }
  .select-collapse {
    flex: 1;
    /deep/.el-input__suffix {
      display: none;
    }
  }
}
.el-select-dropdown__item.selected {
  font-weight: normal;
}
</style>
