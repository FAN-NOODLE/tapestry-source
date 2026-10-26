<template>
  <div class="select-member">
    <div class="wrap">
      <el-tabs
        class="tab-wrap"
        v-show="USER.version !== 3"
        v-model="select_tabs"
        tab-position="left"
        @tab-click="handleTabClick"
      >
        <el-tab-pane
          label="按门店所属区域选择"
          name="SelectMemberByArea"
        ></el-tab-pane>
        <el-tab-pane
          label="按导购所属部门选择"
          name="SelectMemberByDp"
        ></el-tab-pane>
      </el-tabs>
      <component
        :is="select_tabs"
        :ref="select_tabs"
        :selectArr="select_arr"
        :choice_name="select_name"
        :hasSelectedSearch="false"
      ></component>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from "vue-property-decorator";
import SelectMemberByArea from "@/components/sy-components/select-member-by-area.vue";
import SelectMemberByDp from "@/components/sy-components/select-member-by-dp.vue";
import UserModule from "@/store/modules/user";

@Component({
  components: {
    SelectMemberByArea,
    SelectMemberByDp,
  },
})
export default class extends Vue {
  @Prop({ default: () => [] }) selectArr?: Array<any>; // 选择的数据
  @Prop({ default: 2 }) selectedRangeType?: number; // 已选的tab
  @Prop({ default: false }) allowNull?: boolean; // 允许选择空保存
  @Prop({ default: "SelectMemberByArea" }) tabName?: string; // 默认选中tab名称

  private select_tabs = "SelectMemberByArea";
  private select_name = "已选管理范围";
  private select_arr: Array<any> = [];

  get USER() {
    return UserModule.getUserInfo;
  }

  /*** methods ***/
  // 切换tab栏
  private handleTabClick() {
    this.select_arr = [];
  }

  // 选择完成提交数据
  @Emit("submit")
  public save() {
    const type = this.select_tabs === "SelectMemberByDp" ? 4 : 2;
    const arr = (this.$refs[this.select_tabs] as any).select_arr;
    return { arr, type };
  }

  private created() {
    this.select_tabs = this.tabName || "SelectMemberByArea";
    this.select_arr = this.selectArr ? this.selectArr : [];
  }
}
</script>

<style lang="scss" scoped>
/deep/.el-tabs__item {
  padding-left: 0px !important;
  padding-right: 20px !important;
}

/deep/.el-tabs--left {
  .el-tabs__active-bar.is-left,
  .el-tabs__active-bar {
    width: 4px !important;
  }
}

.wrap {
  width: 100%;
  display: flex;

  .tab-wrap {
    width: 146px;
    /deep/.el-tabs--left {
      .el-tabs__active-bar.is-left,
      .el-tabs__active-bar {
        width: 4px !important;
      }
    }
  }
}
</style>
