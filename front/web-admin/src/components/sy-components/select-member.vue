<template>
  <div class="select-member">
    <el-dialog
      :title="dialog_title"
      width="1000px"
      :visible.sync="dialogVisible"
      :close-on-click-modal="false"
      @close="closeDialog"
    >
      <div class="main-box">
        <div class="wrap clearfix">
          <el-tabs
            class="tab-wrap"
            v-if="USER.version !== 3"
            v-model="select_tabs"
            tab-position="left"
            @tab-click="handleClick"
          >
            <el-tab-pane
              v-for="item in tabList"
              :key="item.name"
              :label="item.label"
              :name="item.name"
            ></el-tab-pane>
          </el-tabs>
          <component
            :is="component_map[select_tabs]"
            :ref="select_tabs"
            :selectArr="select_arr"
            :filterArr="filterArr"
            :choice_name="choice_name"
            :is_show_guide="is_show_guide"
            :expandAll="expandAll"
            :worker_status="worker_status"
          ></component>
        </div>
      </div>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="closeDialog">取 消</el-button>
        <el-button size="small" type="primary" @click="save">确 认</el-button>
      </div>
    </el-dialog>
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
  @Prop({ default: "选择员工" }) dialog_title?: string; // 弹窗标题
  @Prop({ default: false }) visible?: boolean; // 控制模态框显隐
  @Prop({ default: () => [] }) selectArr?: Array<any>; // 选择的数据
  @Prop({ default: () => [] }) filterArr?: Array<any>; // 需要过滤的数据
  @Prop({ default: false }) allowNull?: boolean; // 允许选择空保存
  @Prop({ default: false }) isEdit?: boolean; // 是否是编辑状态
  @Prop({ default: 1 }) tabName?: number; // 默认选中tab名称
  @Prop() choice_name?: string; // 右侧选择栏显示名称
  @Prop({ default: false }) is_show_guide?: boolean; // 是否显示导购
  @Prop({ default: false }) reback_type?: boolean; //是否需要返回所选类型(区域/部门)
  @Prop({ default: false }) expandAll?: boolean;
  @Prop({ default: "" }) worker_status?: string | number;

  /*** data ***/
  private dialogVisible = true;
  private select_tabs = "3";
  private tabList: Array<any> = [
    { name: "1", label: "按门店所属区域选择" },
    { name: "3", label: "按导购所属部门选择" },
  ];
  private component_map = {
    "1": "SelectMemberByArea",
    "3": "SelectMemberByDp",
  };
  private select_arr: Array<any> = [];

  get USER() {
    return UserModule.getUserInfo;
  }

  // 切换tab栏
  private handleClick() {
    if (this.select_tabs === String(this.tabName)) {
      this.select_arr = this.selectArr ? this.selectArr : [];
    } else {
      this.select_arr = [];
    }
  }

  // 关闭模态框的回调
  @Emit("close")
  private closeDialog() {
    this.dialogVisible = false;
  }

  // 选择完成提交数据
  @Emit("submit")
  private save() {
    const type = Number(this.select_tabs);
    const list = (this.$refs[this.select_tabs] as any).select_arr;
    if (!list.length) {
      this.$message.warning(
        `请选择${this.is_show_guide ? "导购" : type === 3 ? "部门" : "门店"}`
      );
      return;
    }
    list.forEach((item: any) => {
      const idx = item.name.indexOf("(");
      if (idx > -1) {
        item.name = item.name.substring(0, idx);
      }
    });
    this.closeDialog();
    const temList = {
      type,
      list,
    };
    const reback_list = this.reback_type ? temList : list;
    return reback_list;
  }

  private mounted() {
    this.select_tabs = String(this.tabName);
    this.select_arr = this.selectArr ? this.selectArr : [];
  }
}
</script>

<style lang="scss" scoped>
/deep/.el-dialog__body {
  padding: 20px 30px;
}
/deep/.el-dialog__footer {
  padding: 0 30px 30px;
}
/deep/.el-tabs__item {
  padding-left: 0 !important;
  padding-right: 20px !important;
  font-weight: normal !important;
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
  }
}
</style>
