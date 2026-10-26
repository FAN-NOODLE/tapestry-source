<!-- 类目规则 -->
<template>
  <el-form :model="formData" :rules="rules" :inline="true" ref="ruleForm">
    <!-- 次数规则 -->
    <div v-if="currentRule === 1">
      <span>导购每完成</span>
      <el-form-item prop="secondRule.rule_value" :rules="rules.rule_value">
        <el-input
          :disabled="disabled"
          class="bg-input"
          v-model="formData.secondRule.rule_value"
          oninput="value=value.replace(/^(0+)|[^\d]+/g,'')"
          placeholder="请输入大于0的正整数"
        ></el-input>
      </el-form-item>

      <span>次，增加</span>
      <el-form-item prop="secondRule.integral" :rules="rules.integral">
        <el-input
          :disabled="disabled"
          class="bg-input"
          v-model="formData.secondRule.integral"
          oninput="value=value.replace(/^(0+)|[^\d]+/g,'')"
          placeholder="请输入大于0的正整数"
        ></el-input>
        <span>积分</span>
      </el-form-item>
    </div>

    <!-- 条件规则 -->
    <div v-else-if="currentRule === 2">
      <div
        v-for="(item, index) in formData.conditionRule"
        :key="index"
        class="to-store-rate"
      >
        <span>当导购完成值</span>

        <!-- 高于/达到 -->
        <el-form-item
          :prop="`conditionRule.${index}.rule_if`"
          :rules="rules.rule_if"
        >
          <el-select
            v-model="item.rule_if"
            class="sm-select"
            :disabled="disabled"
          >
            <el-option
              v-for="item in conditions.condition"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>

        <!-- 平均值/百分比 -->
        <el-form-item
          :prop="`conditionRule.${index}.rule_if1`"
          :rules="rules.rule_if1"
        >
          <el-select
            v-model="item.rule_if1"
            class="sm-select"
            :disabled="disabled"
            @change="clearValue(item)"
          >
            <el-option
              v-for="item in conditions.condition1"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>

        <!-- 条件值 -->
        <template v-if="item.rule_if1 === 1">
          <el-form-item
            :prop="`conditionRule.${index}.rule_value`"
            :rules="rules.rule_value"
          >
            <el-input
              :disabled="disabled"
              class="bg-input"
              v-model="item.rule_value"
              oninput="value=value.replace(/^(0+)|[^\d]+/g,'')"
              placeholder="请输入大于0的正整数"
            ></el-input>
          </el-form-item>

          <span>%，增加</span>
        </template>
        <span v-else>，增加</span>

        <!-- 积分 -->
        <el-form-item
          :prop="`conditionRule.${index}.integral`"
          :rules="rules.integral"
        >
          <el-input
            :disabled="disabled"
            class="bg-input"
            v-model="item.integral"
            oninput="value=value.replace(/^(0+)|[^\d]+/g,'')"
            placeholder="请输入大于0的正整数"
          ></el-input>
          <span>积分</span>
        </el-form-item>

        <!-- 按钮操作新增、删除 -->
        <template v-if="!disabled">
          <i
            class="el-icon-remove icon disabled-icon"
            v-if="formData.conditionRule.length === 1"
          />
          <i
            class="el-icon-remove icon"
            v-else
            @click="delTerm('conditionRule', index)"
          />

          <i
            class="el-icon-circle-plus icon disabled-icon"
            v-if="formData.conditionRule.length === 5 && !index"
          />

          <i
            class="el-icon-circle-plus icon"
            v-else-if="!index"
            @click="addTerm('conditionRule')"
          />
        </template>
      </div>
    </div>

    <div v-else>--</div>
  </el-form>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import {
  condition,
  condition1,
  integralCategory,
  ItemType
} from "../../config";

@Component({})
export default class IntegralRule extends Vue {
  /** 奖励规则类型 1.邀约数（每次） 2.邀约到店率（每周） 3.邀约成功率（每周） 4.邀约成功（每次） 5.邀约占比（每周）6. 邀约到店确认（每次）*/
  @Prop() type!: number;

  /** 详情 */
  @Prop() info!: Array<ItemType>;

  /** 是否置灰 */
  @Prop({ default: false }) disabled?: boolean;
  /** todo: store_type:0默认 1奥莱 2零售 */
  @Prop({ default: 0 }) storeType?: number;
  /** 生成item */
  generateItem(key: string): ItemType {
    if (key === "conditionRule") {
      // integral_type 1.加积分 2.减积分  rule_if: 1:次数条件 2:高于3：达到
      return {
        rule_if: 2,
        rule_if1: 2,
        integral_type: 1,
        rule_value: "",
        integral: "",
        store_type: this.storeType
      };
    } else {
      return {
        rule_if: 1,
        rule_if1: 0,
        integral_type: 1,
        store_type: this.storeType
      };
    }
  }

  /** 生成默认表单数据 */
  generateData() {
    return {
      /** 次数规则 */
      secondRule: this.generateItem("secondRule"),
      /** 条件规则 */
      conditionRule: [this.generateItem("conditionRule")]
    };
  }

  @Watch("type", { immediate: true })
  changeType() {
    // 重置表单数据，防止数据缓存
    this.formData = this.generateData();
    // 获取当前类目
    const currentCategory = integralCategory.filter(
      item => item.value === this.type
    );
    this.currentRule = currentCategory[0].ruleType;

    // 先清空校验
    this.$nextTick(() => {
      (this.$refs["ruleForm"] as HTMLFormElement).clearValidate();
    });
  }

  /** 监听奖励规则类型 变化 */
  @Watch("info", { deep: true, immediate: true })
  changeInfo() {
    // 编辑或者查看 重新赋值
    if (this.currentRule === 1 && this.info && this.info.length > 0) {
      this.formData.secondRule = this.info[0];
    } else {
      this.formData.secondRule = this.generateItem("secondRule");
    }
    // 编辑或者查看 重新赋值
    if (this.currentRule === 2 && this.info && this.info.length > 0) {
      this.formData.conditionRule = this.info;
    } else {
      this.formData.conditionRule = [this.generateItem("conditionRule")];
    }
  }

  /** 表单数据 */
  private formData = this.generateData();

  /** 现有规则 */
  private currentRule = 0;

  /** 条件枚举 */
  private conditions = { condition, condition1 };

  /** 校验规则 */
  private rules = {
    integral: [{ required: true, message: "请输入积分", trigger: "blur" }],
    rule_value: [
      { required: true, validator: this.validRuleValue, trigger: "blur" }
    ],
    rule_if: [{ required: true, message: "请选择条件", trigger: "change" }],
    rule_if1: [
      { required: true, validator: this.validRuleIf1, trigger: "change" }
    ]
  };

  /** 清空指定值 */
  private clearValue(item: ItemType) {
    item.rule_value = "";
  }

  /** 校验条件 */
  private validRuleIf1(
    rule: any,
    value: { getTime: () => number }[],
    callback: (error?: Error) => void
  ) {
    if (!value) {
      return callback(new Error("请选择条件"));
    }
    if (this.currentRule === 2) {
      const index = Number(rule.fullField.split(".")[1]);
      const filterItem = this.formData.conditionRule[index];

      // 只有平均值需要校验
      if (filterItem.rule_if1 !== 2) return callback();

      // 判断值是否相等
      const filterRules = this.formData.conditionRule.filter(item => {
        return (
          item.rule_if === filterItem.rule_if &&
          item.rule_if1 === 2 &&
          filterItem.rule_if1 === 2
        );
      });
      if (filterRules.length > 1) {
        return callback(new Error("不能存在相同的规则"));
      }
    }

    callback();
  }

  /** 校验规则值 */
  private validRuleValue(
    rule: any,
    value: { getTime: () => number }[],
    callback: (error?: Error) => void
  ) {
    if (!value) {
      return callback(new Error("请输入大于0的正整数"));
    }

    if (this.currentRule === 2) {
      const index = Number(rule.fullField.split(".")[1]);
      const filterItem = this.formData.conditionRule[index];

      // 判断值是否相等
      const filterRules = this.formData.conditionRule.filter(item => {
        return (
          item.rule_if === filterItem.rule_if &&
          item.rule_if1 === filterItem.rule_if1 &&
          Number(item.rule_value) === Number(filterItem.rule_value)
        );
      });
      if (filterRules.length > 1) {
        return callback(new Error("不能存在相同的规则"));
      }
    }

    callback();
  }

  /** 供父组件获取提交数据 */
  public getData() {
    // 数据转换
    const dataMap = ["secondRule", "conditionRule"];
    let data = this.formData[dataMap[this.currentRule - 1]];

    // 次数规则转换成数组
    if (this.currentRule === 1) {
      data = [data];
    }

    // 格式化数组
    data = data.map((item: { rule_value: string; integral: string }) => ({
      ...item,
      rule_value: Number(item.rule_value),
      integral: Number(item.integral)
    }));
    return data;
  }

  /** 新增项 */
  private addTerm(key: string) {
    this.formData[key].push(this.generateItem(key));
  }

  /** 删除项 */
  private delTerm(key: string, index: number) {
    this.formData[key].splice(index, 1);
  }
}
</script>
<style lang="scss" scoped>
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

/deep/.el-form-item {
  margin-right: 0px !important;
}

.bg-input {
  width: 170px !important;
  margin: 0 10px;
}
.sm-select {
  margin: 0 5px;
  /deep/.el-input {
    width: 100px !important;
  }
}

.to-store-rate {
  margin-bottom: 10px;
}
</style>
