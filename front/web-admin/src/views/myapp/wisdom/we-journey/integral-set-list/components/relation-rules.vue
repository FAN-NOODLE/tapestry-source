<!-- 类目规则 -->
<template>
  <el-form :model="formData" :rules="rules" :inline="true" ref="ruleForm">
    <div
      v-for="(item, index) in formData.rule_conditions"
      :key="index"
      class="rule"
    >
      <span>{{ conditions.ruleTitle }}</span>

      <!-- 高于/达到 -->
      <el-form-item
        :prop="`rule_conditions.${index}.rule_condition`"
        :rules="rules.rule_condition"
      >
        <el-select v-model="item.rule_condition" class="sm-select">
          <el-option
            v-for="item in conditions.condition3"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          >
          </el-option>
        </el-select>
      </el-form-item>

      <!-- 条件值 -->
      <el-form-item
        :prop="`rule_conditions.${index}.rule_value`"
        :rules="rules.rule_value"
      >
        <el-input
          class="bg-input"
          v-model.number="item.rule_value"
          oninput="value=value.replace(/^(0+)|[^\d]+/g,'')"
          placeholder="请输入大于0的正整数"
        ></el-input>
      </el-form-item>

      <span>{{ calcUnit(item.rule_condition) }}，能力指标为</span>

      <!-- 能力指标 -->
      <el-form-item
        :prop="`rule_conditions.${index}.rule_give`"
        :rules="rules.rule_give"
      >
        <el-select
          v-model="item.rule_give"
          class="sm-select"
          @visible-change="changeGive"
        >
          <el-option
            v-for="item in conditions.condition4"
            :key="item.value"
            :label="item.label"
            :value="item.label"
            :disabled="item.disabled"
          >
          </el-option>
        </el-select>
      </el-form-item>

      <!-- 按钮操作新增、删除 -->
      <i
        class="el-icon-remove icon disabled-icon"
        v-if="rulesArray.length === 1"
      />
      <i class="el-icon-remove icon" v-else @click="delTerm(index)" />

      <i
        class="el-icon-circle-plus icon disabled-icon"
        v-if="rulesArray.length === 5 && !index"
      />

      <i
        class="el-icon-circle-plus icon"
        v-else-if="!index"
        @click="addTerm()"
      />
    </div>
  </el-form>
</template>
<script lang="ts">
import { Vue, Component, Prop, Watch } from "vue-property-decorator";
import { condition3, condition3_1, condition4, condition5 } from "../../config";

type relation_rule = {
  rule_condition: string | number;
  rule_value: string;
  rule_give: string | number;
};
@Component({})
export default class RelationRules extends Vue {
  // 16.奖励规则类型 线上销售业绩率（每月），17 线上销售金额（每月
  @Prop() type!: number;

  // 规则数组数据
  @Prop({ default: () => [] }) rulesArray!: relation_rule[];
  @Watch("rulesArray", { immediate: true, deep: true })
  watchRulesArray(val: relation_rule[]) {
    this.formData.rule_conditions = val;
  }

  // todo: store_type:0默认 1奥莱 2零售
  @Prop({ default: 0 }) storeType?: number;
  @Watch("type", { immediate: true })
  watchType(val: number) {
    this.conditions.ruleTitle = this.defaultText[val].ruleTitle;
    // todo:如果是3,4互动，培训维度使用top百分比
    if(val === 3 || val === 4) {
      this.conditions.condition3 = condition3_1
    } else {
      this.conditions.condition3 = condition3
    }
  }
  //默认标题
  defaultText = [
    { title: "", ruleTitle: "" },
    { title: "邀约能力指标配置", ruleTitle: "当导购邀约目标率" },
    { title: "任务能力指标配置", ruleTitle: "当导购任务完成率" },
    { title: "客户互动能力指标配置", ruleTitle: "当导购总互动积分排行" },
    { title: "培训能力指标配置", ruleTitle: "当导购总培训积分排行" },
    { title: "线上能力指标配置", ruleTitle: "当导购线上目标率" }
  ];
  formData = {
    rule_conditions: [] as relation_rule[]
  };
  conditions = {
    condition3,
    condition4,
    condition5,
    ruleTitle: ""
  };
  rules = {
    rule_condition: [{ required: true, message: "请选择", trigger: "change" }],
    rule_value: [
      { required: true, message: "请输入大于0的正整数", trigger: "blur" }
    ],
    rule_give: [
      { required: true, message: "请选择能力指标等级", trigger: "change" }
    ]
  };
  /**
   * @description: 计算单位
   * @param {*} value 1: top榜 2: 邀约维度 3: 客户互动维度
   */

  calcUnit(value: number | string) {
    // top 榜单位是名
    if (value === 1) {
      return "名";
    }
    // todo：5 top百分比
    if (value === 5) {
      return "%";
    }
    // 邀约维度、客户互动维度 除了【top榜】， 单位都是%
    if (/[125]/.test(this.type + "")) {
      return "%";
    }
    return "名";
  }
  /**
   * @description: 能力指标改变
   */
  changeGive(open: boolean) {
    this.$emit("changeGive", { open, storeType: this.storeType });
  }
  /**
   * @description: 新增
   */
  addTerm() {
    this.$emit("addTerm", this.storeType);
  }
  /**
   * @description: 删除
   */

  delTerm(index: number) {
    this.$emit("delTerm", { index, storeType: this.storeType });
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-form-item {
  margin-right: 0px !important;
}

.sm-select {
  margin: 0 5px;
  /deep/.el-input {
    width: 120px !important;
  }
}

.bg-input {
  width: 170px !important;
  margin: 0 10px;
}
.icon {
  cursor: pointer;
  font-size: 20px;
  vertical-align: middle;
  margin-left: 10px;
}

.disabled-icon {
  color: #909399;
  cursor: not-allowed;
}
.rule {
  margin-bottom: 10px;
}
</style>
